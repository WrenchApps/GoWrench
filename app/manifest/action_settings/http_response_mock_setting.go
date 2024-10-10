package action_settings

import (
	"wrench/app/manifest/validation"
)

type HttpResponseMockSetting struct {
	Value       string `yaml:"value"`
	ContentType string `yaml:"contentType"`
	Method      string `yaml:"method"`
}

func (setting HttpResponseMockSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Value) == 0 {
		result.AddError("actions.mock.value is required")
	}

	if len(setting.ContentType) == 0 {
		result.AddError("actions.mock.contentType is required")
	}

	if len(setting.Method) == 0 {
		result.AddError("actions.mock.method is required")
	}

	return result
}
