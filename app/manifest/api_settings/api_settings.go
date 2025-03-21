package api_settings

import (
	"errors"
	"wrench/app/manifest/validation"
)

type ApiSettings struct {
	Endpoints     []EndpointSettings     `yaml:"endpoints"`
	Authorization *AuthorizationSettings `yaml:"authorization"`
}

func (setting ApiSettings) HasAuthorization() bool {
	return setting.Authorization != nil
}

func (setting ApiSettings) GetEndpointByRoute(route string) (*EndpointSettings, error) {
	for _, endpoint := range setting.Endpoints {
		if endpoint.Route == route {
			return &endpoint, nil
		}
	}

	return nil, errors.New("endpoint not found")
}

func (setting ApiSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Endpoints) == 0 {
		result.AddError("api.endpoints is required")
	} else {
		for _, validable := range setting.Endpoints {
			result.AppendValidable(validable)
		}
	}

	if setting.Authorization != nil {
		result.AppendValidable(setting.Authorization)
	}

	return result
}
