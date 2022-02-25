package api

type CreateApplicationRequest struct {
	Application `json:"application"`
}

type Application struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	OrganizationID   string `json:"organizationID"`   // nolint: tagliatelle
	ServiceProfileID string `json:"serviceProfileID"` // nolint: tagliatelle
}
