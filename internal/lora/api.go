package lora

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// API for accessing chirpstack.
type API struct {
	Client *resty.Client
}

func NewAPI(baseURL string) API {
	client := resty.New()

	client.SetBaseURL(baseURL)

	return API{
		Client: client,
	}
}

type ActivationDeviceRequest struct {
	DeviceActivation `json:"deviceActivation"`
}

type DeviceActivation struct {
	DevEUI          string `json:"devEUI"` // nolint: tagliatelle
	DevAddr         string `json:"devAddr"`
	ApplicationSKey string `json:"appSKey"`
	NetworkSKey     string `json:"nwkSEncKey"`
}

func (a API) Activate(devEUI, devAddr, applicationSKey, networkSKey string) error {
	_, err := a.Client.R().
		SetHeader("Content-Type", "application/json").
		SetPathParam("devEUI", devEUI).
		SetBody(ActivationDeviceRequest{
			DeviceActivation{
				DevEUI:          devEUI,
				DevAddr:         devAddr,
				ApplicationSKey: applicationSKey,
				NetworkSKey:     networkSKey,
			},
		}).
		Post("/api/devices/{devEUI}/activate")
	if err != nil {
		return fmt.Errorf("activation request failed %w", err)
	}

	return nil
}
