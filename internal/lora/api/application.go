package api

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrDuplicateApp = errors.New("duplicate application")

type CreateApplicationRequest struct {
	Application `json:"application"`
}

type Application struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	OrganizationID   string `json:"organizationID"`   // nolint: tagliatelle
	ServiceProfileID string `json:"serviceProfileID"` // nolint: tagliatelle
}

func (a API) CreateApplication(name, description, orgID, spID string) error {
	resp, err := a.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(CreateApplicationRequest{
			Application: Application{
				Name:             name,
				Description:      description,
				OrganizationID:   orgID,
				ServiceProfileID: spID,
			},
		}).
		Post("/api/applications")
	if err != nil {
		return fmt.Errorf("application creation request failed %w", err)
	}

	if resp.IsSuccess() {
		return nil
	}

	if resp.StatusCode() == http.StatusConflict {
		return ErrDuplicateApp
	}

	// nolint: goerr113
	return fmt.Errorf("application creation failed with %d", resp.StatusCode())
}
