package api

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/pterm/pterm"
)

const (
	DevEUILen  = 8
	DevAddrLen = 4
)

// API for accessing chirpstack.
type API struct {
	Client   *resty.Client
	Username string
	Password string
}

func New(cfg Config) API {
	client := resty.New()

	a := API{
		Client:   client,
		Username: cfg.Username,
		Password: cfg.Password,
	}

	client.SetBaseURL(cfg.URL)
	client.SetRetryCount(1)
	client.AddRetryCondition(func(r *resty.Response, e error) bool {
		return r.StatusCode() == http.StatusUnauthorized
	})
	client.AddRetryHook(func(r *resty.Response, e error) {
		if r.StatusCode() == http.StatusUnauthorized {
			a.Login()
		}
	})

	return a
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	JWT string `json:"jwt"`
}

func (a *API) Login() {
	var jwt LoginResponse

	resp, err := a.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(LoginRequest{
			Email:    a.Username,
			Password: a.Password,
		}).
		SetResult(&jwt).
		Post("/api/internal/login")
	if err != nil {
		pterm.Fatal.Printf("cannot login into loraserver %s\n", err)
	}

	if resp.IsError() {
		pterm.Fatal.Printf("cannot login into loraserver %d\n", resp.StatusCode())
	}

	a.Client.SetAuthToken(jwt.JWT)
}

func (a API) Activate(devEUI, devAddr, applicationSKey, networkSKey string) error {
	resp, err := a.Client.R().
		SetHeader("Content-Type", "application/json").
		SetPathParam("devEUI", devEUI).
		SetBody(ActivationDeviceRequest{
			DeviceActivation{
				DevEUI:                      devEUI,
				DevAddr:                     devAddr,
				ApplicationSKey:             applicationSKey,
				NetworkSEncKey:              networkSKey,
				ServingNetworkSIntKey:       networkSKey,
				ForwardingNetworkSIntKey:    networkSKey,
				UplinkFrameCounter:          0,
				DownlinkNetworkFrameCounter: 0,
				DownlinkAppFrameCounter:     0,
			},
		}).
		Post("/api/devices/{devEUI}/activate")
	if err != nil {
		return fmt.Errorf("activation request failed %w", err)
	}

	if resp.IsSuccess() {
		return nil
	}

	// nolint: goerr113
	return fmt.Errorf("activation request failed with %d", resp.StatusCode())
}
