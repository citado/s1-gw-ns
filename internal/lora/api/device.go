package api

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrDuplicateDevice = errors.New("duplicate device")

type ActivationDeviceRequest struct {
	DeviceActivation `json:"deviceActivation"`
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
	Device `json:"device"`
}

type Device struct {
	DevEUI            string            `json:"devEUI"` // nolint: tagliatelle
	Name              string            `json:"name"`
	ApplicationID     int64             `json:"applicationID"` // nolint: tagliatelle
	Description       string            `json:"description"`
	DeviceProfileID   string            `json:"deviceProfileID"` // nolint: tagliatelle
	SkipFCntCheck     bool              `json:"skipFCntCheck"`
	ReferenceAltitude float64           `json:"referenceAltitude"`
	Variables         map[string]string `json:"variables"`
	Tags              map[string]string `json:"tags"`
	IsDisabled        bool              `json:"isDisabled"`
}

func (a API) CreateDevice(devEUI string, name string, applicationID int64, description string,
	deviceProfileID string, skipFCntCheck bool, referenceAlltitude float64, isDisabled bool) error {
	resp, err := a.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(CreateDeviceRequest{
			Device: Device{
				DevEUI:            devEUI,
				Name:              name,
				ApplicationID:     applicationID,
				Description:       description,
				DeviceProfileID:   deviceProfileID,
				SkipFCntCheck:     skipFCntCheck,
				ReferenceAltitude: referenceAlltitude,
				Variables:         map[string]string{},
				Tags:              map[string]string{},
				IsDisabled:        isDisabled,
			},
		}).
		Post("/api/devices")
	if err != nil {
		return fmt.Errorf("device creation request failed %w", err)
	}

	if resp.IsSuccess() {
		return nil
	}

	if resp.StatusCode() == http.StatusConflict {
		return ErrDuplicateDevice
	}

	// nolint: goerr113
	return fmt.Errorf("device creation failed with %d", resp.StatusCode())
}
