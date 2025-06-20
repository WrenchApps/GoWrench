package http_settings

import (
	"wrench/app/manifest/validation"
)

type HttpRequestMockSetting struct {
	Body        string            `yaml:"body"`
	ContentType string            `default:"application/json" yaml:"contentType"`
	Headers     map[string]string `yaml:"headers"`
	StatusCode  int               `default:"200" yaml:"statusCode"`
	MirrorBody  bool              `yaml:"mirrorBody"`
}

func (setting HttpRequestMockSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Body) == 0 && setting.MirrorBody == false {
		result.AddError("actions.http.mock.body is required or setting.mirrorBody equals true")
	}

	if len(setting.ContentType) == 0 {
		result.AddError("actions.http.mock.contentType is required")
	}

	return result
}
