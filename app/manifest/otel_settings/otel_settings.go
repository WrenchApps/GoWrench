package otel_settings

import (
	"fmt"
	"wrench/app/manifest/validation"
)

type OtelSettings struct {
	CollectorUrl        string            `yaml:"collectorUrl"`
	Enable              bool              `yaml:"enable"`
	TraceConsoleExport  bool              `yaml:"traceConsoleExport"`
	MetricConsoleExport bool              `yaml:"metricConsoleExport"`
	TraceTags           map[string]string `yaml:"traceTags"`
}

func (setting OtelSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if setting.Enable && len(setting.CollectorUrl) == 0 {
		result.AddError("otel.collectorUrl is required")
	}

	if len((setting.TraceTags)) > 0 {
		for tagKey, tagValue := range setting.TraceTags {
			if len(tagValue) == 0 {
				result.AddError(fmt.Sprintf("otel.tags %v should contain value", tagKey))
			}
		}
	}

	return result
}
