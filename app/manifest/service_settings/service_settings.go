package service_settings

import (
	"wrench/app/manifest/aws_settings"
	"wrench/app/manifest/otel_settings"
	"wrench/app/manifest/validation"
)

type ServiceSettings struct {
	Name    string                      `yaml:"name"`
	Version string                      `yaml:"version"`
	Otel    *otel_settings.OtelSettings `yaml:"otel"`
	Aws     *aws_settings.AwsSettings   `yaml:"aws"`
}

func (setting ServiceSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Name) == 0 {
		result.AddError("service.name is required")
	}

	if len(setting.Version) == 0 {
		result.AddError("service.version is required")
	}

	if setting.Otel != nil {
		result.AppendValidable(setting.Otel)
	}

	if setting.Aws != nil {
		result.AppendValidable(setting.Aws)
	}

	return result
}
