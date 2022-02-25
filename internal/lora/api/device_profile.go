package api

import "fmt"

type ListDeviceProfileResponse struct {
	Result []DeviceProfile `json:"result"`
}

type CreateDeviceProfileRequest struct {
	DeviceProfile `json:"deviceProfile"`
}

type DeviceProfile struct {
	ID                         string `json:"id"`
	Name                       string `json:"name"`
	OrganizationID             string `json:"organizationID"`  // nolint: tagliatelle
	NetworkServerID            string `json:"networkServerID"` // nolint: tagliatelle
	MACVersion                 string `json:"macVersion"`
	RegionalParametersRevision string `json:"regParamsRevision"`
	ADRAlgorithmID             string `json:"adrAlgorithmID"` // nolint: tagliatelle
	MaxEIRP                    int64  `json:"maxEIRP"`        // nolint: tagliatelle
	UplinekInterval            string `json:"uplinkInterval"`
}

type CreateDeviceProfileResponse struct {
	ID string `json:"id"`
}

func (a API) GetOrCreateDeviceProfile(name, orgID, nsID string) (string, error) {
	var deviceProfiles ListDeviceProfileResponse

	resp, err := a.Client.R().
		SetHeader("Content-Type", "application/json").
		SetQueryParam("limit", "100"). // in prduction it should be set with more caution
		SetQueryParam("organizationID", orgID).
		SetQueryParam("networkServerID", nsID).
		SetResult(&deviceProfiles).
		Get("/api/device-profiles")
	if err != nil {
		return "", fmt.Errorf("device profiles list request failed %w", err)
	}

	if resp.IsError() {
		// nolint: goerr113
		return "", fmt.Errorf("device profiles list failed with %d", resp.StatusCode())
	}

	for _, deviceProfile := range deviceProfiles.Result {
		if deviceProfile.Name == name {
			return deviceProfile.ID, nil
		}
	}

	return a.CreateDeviceProfile(name, orgID, nsID)
}

func (a API) CreateDeviceProfile(name, orgID, nsID string) (string, error) {
	var id CreateDeviceProfileResponse

	resp, err := a.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(CreateDeviceProfileRequest{
			DeviceProfile: DeviceProfile{
				ID:                         "",
				Name:                       name,
				OrganizationID:             orgID,
				NetworkServerID:            nsID,
				MACVersion:                 "1.0.2",
				RegionalParametersRevision: "A",
				ADRAlgorithmID:             "",
				UplinekInterval:            "0",
				MaxEIRP:                    0,
			},
		}).
		SetResult(&id).
		Post("/api/device-profiles")
	if err != nil {
		return "", fmt.Errorf("device profile creation request failed %w", err)
	}

	if resp.IsSuccess() {
		return id.ID, nil
	}

	// nolint: goerr113
	return "", fmt.Errorf("device profile creation failed with %d", resp.StatusCode())
}
