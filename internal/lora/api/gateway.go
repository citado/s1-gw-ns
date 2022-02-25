package api

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrDuplicateGateway = errors.New("duplicate gateway")

type CreateGatewayRequest struct {
	Gateway `json:"gateway"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
}

type Gateway struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	OrganizationID   string   `json:"organizationID"` // nolint: tagliatelle
	DiscoveryEnabled bool     `json:"discoveryEnabled"`
	NetworkServerID  string   `json:"networkServerID"`  // nolint: tagliatelle
	ServiceProfileID string   `json:"serviceProfileID"` // nolint: tagliatelle
	Location         Location `json:"location"`
}

func (a API) CreateGateway(id, name, description, orgID, nsID, spID string) error {
	resp, err := a.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(CreateGatewayRequest{
			Gateway: Gateway{
				ID:               id,
				Name:             name,
				Description:      description,
				OrganizationID:   orgID,
				DiscoveryEnabled: false,
				NetworkServerID:  nsID,
				ServiceProfileID: spID,
				// nolint: gomnd
				Location: Location{
					Latitude:  35.723737,
					Longitude: 50.952981,
					Altitude:  0.0,
				},
			},
		}).
		Post("/api/gateways")
	if err != nil {
		return fmt.Errorf("gateway creation request failed %w", err)
	}

	if resp.IsSuccess() {
		return nil
	}

	if resp.StatusCode() == http.StatusConflict {
		return ErrDuplicateGateway
	}

	// nolint: goerr113
	return fmt.Errorf("gateway creation failed with %d", resp.StatusCode())
}
