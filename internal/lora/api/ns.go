package api

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrDuplicateNS = errors.New("duplicate network-server")

type CreateNetworkServerRequest struct {
	NetworkServer `json:"networkServer"`
}

type NetworkServer struct {
	ID                          string `json:"id"`
	Name                        string `json:"name"`
	Server                      string `json:"server"`
	GatewayDiscoveryEnabled     bool   `json:"gatewayDiscoveryEnabled"`
	GatewayDiscoveryInterval    int64  `json:"gatewayDiscoveryInterval"`
	GatewayDiscoveryTXFrequency int64  `json:"gatewayDiscoveryTXFrequency"` // nolint: tagliatelle
	GatewayDiscoveryDR          int64  `json:"gatewayDiscoveryDR"`          // nolint: tagliatelle
}

func (a API) CreateNetworkServer(id, name, server string) error {
	resp, err := a.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(CreateNetworkServerRequest{
			NetworkServer: NetworkServer{
				ID:                          id,
				Name:                        name,
				Server:                      server,
				GatewayDiscoveryEnabled:     false,
				GatewayDiscoveryInterval:    0,
				GatewayDiscoveryTXFrequency: 0,
				GatewayDiscoveryDR:          0,
			},
		}).
		Post("/api/network-servers")
	if err != nil {
		return fmt.Errorf("network server creation request failed %w", err)
	}

	if resp.IsSuccess() {
		return nil
	}

	if resp.StatusCode() == http.StatusConflict {
		return ErrDuplicateNS
	}

	// nolint: goerr113
	return fmt.Errorf("network server creation failed with %d", resp.StatusCode())
}
