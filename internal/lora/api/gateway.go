package api

type Gateway struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	OrganizationID   string `json:"organizationID"` // nolint: tagliatelle
	DiscoveryEnabled bool   `json:"discoveryEnabled"`
	NetworkServerID  string `json:"networServerID"`   // nolint: tagliatelle
	GatewayProfileID string `json:"gatewayProfileID"` // nolint: tagliatelle
	ServiceProfileID string `json:"serviceProfileID"` // nolint: tagliatelle
}

func (a API) CreateGateway() error {
	return nil
}
