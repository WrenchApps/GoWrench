package http_settings

import (
	"wrench/app/manifest/validation"
)

type HttpRequestMockSettings struct {
	Body        string            `yaml:"body"`
	ContentType string            `default:"application/json" yaml:"contentType"`
	Headers     map[string]string `yaml:"headers"`
	StatusCode  string            `default:"200" yaml:"statusCode"`
	MirrorBody  bool              `yaml:"mirrorBody"`
}

func (setting HttpRequestMockSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	return result
}
