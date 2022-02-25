package api

import (
	"fmt"
)

type ListServiceProfileResponse struct {
	Result []ServiceProfile `json:"result"`
}

type CreateServiceProfileRequest struct {
	ServiceProfile `json:"serviceProfile"`
}

type ServiceProfile struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	OrganizationID     string `json:"organizationID"`  // nolint: tagliatelle
	NetworkServerID    string `json:"networkServerID"` // nolint: tagliatelle
	ULRate             int64  `json:"ulRate"`
	ULBucketSize       int64  `json:"ulBucketSize"`
	DLRate             int64  `json:"dlRate"`
	DLBucketSize       int64  `json:"dlBucketSize"`
	AddGatewayMetaData bool   `json:"addGWMetaData"` // nolint: tagliatelle
}

type CreateServiceProfileResponse struct {
	ID string `json:"id"`
}

func (a API) GetOrCreateServiceProfile(name, orgID, nsID string) (string, error) {
	var serviceProfiles ListServiceProfileResponse

	resp, err := a.Client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("limit", "100"). // in prduction it should be set with more caution
		SetQueryParam("organizationID", orgID).
		SetQueryParam("networkServerID", nsID).
		SetResult(&serviceProfiles).
		Get("/api/service-profiles")
	if err != nil {
		return "", fmt.Errorf("service profiles list request failed %w", err)
	}

	if resp.IsError() {
		// nolint: goerr113
		return "", fmt.Errorf("service profiles list failed with %d", resp.StatusCode())
	}

	for _, serviceProfile := range serviceProfiles.Result {
		if serviceProfile.Name == name {
			return serviceProfile.ID, nil
		}
	}

	return a.CreateServiceProfile(name, orgID, nsID)
}

func (a API) CreateServiceProfile(name, orgID, nsID string) (string, error) {
	var id CreateServiceProfileResponse

	resp, err := a.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(CreateServiceProfileRequest{
			ServiceProfile: ServiceProfile{
				ID:                 "",
				Name:               name,
				OrganizationID:     orgID,
				NetworkServerID:    nsID,
				ULRate:             0,
				ULBucketSize:       0,
				DLRate:             0,
				DLBucketSize:       0,
				AddGatewayMetaData: true,
			},
		}).
		SetResult(&id).
		Post("/api/service-profiles")
	if err != nil {
		return "", fmt.Errorf("service profile creation request failed %w", err)
	}

	if resp.IsSuccess() {
		return id.ID, nil
	}

	// nolint: goerr113
	return "", fmt.Errorf("service profile creation failed with %d", resp.StatusCode())
}
