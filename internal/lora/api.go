package lora

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/pterm/pterm"
)

type APIConfig struct {
	URL      string
	Username string
	Password string
}

// API for accessing chirpstack.
type API struct {
	Client   *resty.Client
	Username string
	Password string
	Token    string
}

func NewAPI(cfg APIConfig) API {
	client := resty.New()

	client.SetBaseURL(cfg.URL)

	return API{
		Client:   client,
		Username: cfg.Username,
		Password: cfg.Password,
		Token:    "",
	}
}

type ActivationDeviceRequest struct {
	DeviceActivation `json:"deviceActivation"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	JWT string `json:"jwt"`
}

type DeviceActivation struct {
	DevEUI                      string `json:"devEUI"` // nolint: tagliatelle
	DevAddr                     string `json:"devAddr"`
	ApplicationSKey             string `json:"appSKey"`
	NetworkSEncKey              string `json:"nwkSEncKey"`
	ServingNetworkSIntKey       string `json:"sNwkSIntKey"`
	ForwardingNetworkSIntKey    string `json:"fNwkSIntKey"`
	UplinkFrameCounter          int    `json:"fCntUp"`
	DownlinkNetworkFrameCounter int    `json:"nFCntDown"`
	DownlinkAppFrameCounter     int    `json:"aFCntDown"`
}

type CreateDeviceRequest struct {
	APIDevice `json:"device"`
}

type APIDevice struct {
	DevEUI            string  `json:"devEUI"` // nolint: tagliatelle
	Name              string  `json:"name"`
	ApplicationID     int64   `json:"applicationID"` // nolint: tagliatelle
	Description       string  `json:"description"`
	DeviceProfileID   string  `json:"deviceProfileID"` // nolint: tagliatelle
	SkipFCntCheck     bool    `json:"skipFCntCheck"`
	ReferenceAltitude float64 `json:"referenceAltitude"`
	IsDisabled        bool    `json:"isDisabled"`
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

	a.Token = jwt.JWT
}

func (a API) CreateDevice(devEUI string, name string, applicationID int64, description string,
	deviceProfileID string, skipFCntCheck bool, referenceAlltitude float64, isDisabled bool) error {
	resp, err := a.Client.R().
		SetAuthToken(a.Token).
		SetHeader("Content-Type", "application/json").
		SetBody(CreateDeviceRequest{
			APIDevice: APIDevice{
				DevEUI:            "",
				Name:              "",
				ApplicationID:     0,
				Description:       "",
				DeviceProfileID:   "",
				SkipFCntCheck:     false,
				ReferenceAltitude: 0.0,
				IsDisabled:        false,
			},
		}).
		Post("/api/device")
	if err != nil {
		return fmt.Errorf("activation request failed %w", err)
	}

	if resp.IsSuccess() {
		return nil
	}

	return nil
}

func (a API) Activate(devEUI, devAddr, applicationSKey, networkSKey string) error {
	resp, err := a.Client.R().
		SetAuthToken(a.Token).
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

	if resp.StatusCode() == http.StatusUnauthorized {
		a.Login()

		return a.Activate(devEUI, devAddr, applicationSKey, networkSKey)
	}

	return nil
}
