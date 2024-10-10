package api_settings

import (
	"wrench/app/manifest/validation"
)

type ApiSetting struct {
	Endpoints []EndpointSetting `yaml:"endpoints"`
}

func (setting ApiSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Endpoints) == 0 {
		result.AddError("api.endpoints is required")
	} else {
		for _, validable := range setting.Endpoints {
			result.AppendValidable(validable)
		}
	}

	return result
}
