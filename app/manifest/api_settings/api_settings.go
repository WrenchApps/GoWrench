package api_settings

import (
	"errors"
	"wrench/app/manifest/validation"
)

type ApiSettings struct {
	Endpoints     []EndpointSettings     `yaml:"endpoints"`
	Authorization *AuthorizationSettings `yaml:"authorization"`
	Cors          *CorsSettings          `yaml:"cors"`
}

func (setting *ApiSettings) HasAuthorization() bool {
	return setting.Authorization != nil
}

func (setting *ApiSettings) GetEndpointByRoute(route string) (*EndpointSettings, error) {
	for _, endpoint := range setting.Endpoints {
		if endpoint.Route == route {
			return &endpoint, nil
		}
	}

	return nil, errors.New("endpoint not found")
}

func (settings *ApiSettings) Merge(toMerge *ApiSettings) error {

	if toMerge == nil {
		return nil
	}

	if settings.Authorization != nil {
		return errors.New("should configure only once api.authorization")
	}

	if len(toMerge.Endpoints) > 0 {
		if len(settings.Endpoints) == 0 {
			settings.Endpoints = toMerge.Endpoints
		} else {
			settings.Endpoints = append(settings.Endpoints, toMerge.Endpoints...)
		}
	}

	if settings.Cors == nil && toMerge.Cors != nil {
		settings.Cors = &CorsSettings{}
	}
	if settings.Cors != nil {
		if err := settings.Cors.Merge(toMerge.Cors); err != nil {
			return err
		}
	}

	return nil
}

func (setting *ApiSettings) Valid() validation.ValidateResult {
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

	if setting.Cors != nil {
		result.AppendValidable(setting.Cors)
	}

	return result
}
