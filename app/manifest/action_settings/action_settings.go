package action_settings

import (
	"fmt"
	"wrench/app/manifest/action_settings/http_settings"
	"wrench/app/manifest/validation"
)

type ActionSettings struct {
	Id   string                     `yaml:"id"`
	Type ActionType                 `yaml:"type"`
	Http *http_settings.HttpSetting `yaml:"http"`
}

type ActionType string

const (
	ActionTypeHttpRequest     ActionType = "httpRequest"
	ActionTypeHttpRequestMock ActionType = "httpRequestMock"
)

func (setting ActionSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Id) == 0 {
		result.AddError("actions.id is required")
	}

	if len(setting.Type) == 0 {
		var msg = fmt.Sprintf("actions[%s].type is required", setting.Id)
		result.AddError(msg)
	} else {
		if (setting.Type == ActionTypeHttpRequest ||
			setting.Type == ActionTypeHttpRequestMock) == false {

			var msg = fmt.Sprintf("actions[%s].type should contain valid value", setting.Id)
			result.AddError(msg)
		}
	}

	if setting.Type == ActionTypeHttpRequest {
		setting.Http.ValidTypeActionTypeHttpRequest(&result)
	}

	if setting.Type == ActionTypeHttpRequestMock {
		setting.Http.ValidTypeActionTypeHttpRequestMock(&result)
	}

	return result
}
