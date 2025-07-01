package otel_settings

import "wrench/app/manifest/validation"

type OtelSettings struct {
	CollectorUrl        string `yaml:"collectorUrl"`
	Enable              bool   `yaml:"enable"`
	TraceConsoleExport  bool   `yaml:"traceConsoleExport"`
	MetricConsoleExport bool   `yaml:"metricConsoleExport"`
	LogConsoleExport    bool   `yaml:"logConsoleExport"`
}

func (setting OtelSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if setting.Enable && len(setting.CollectorUrl) == 0 {
		result.AddError("otel.collectorUrl is required")
	}

	return result
}
