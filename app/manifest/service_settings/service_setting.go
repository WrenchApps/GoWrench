package service_settings

import (
	"wrench/app/manifest/validation"
)

type ServiceSetting struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

func (setting ServiceSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Name) == 0 {
		result.AddError("service.name is required")
	}

	if len(setting.Version) == 0 {
		result.AddError("service.version is required")
	}

	return result
}
