package action_settings

import (
	"fmt"
	"wrench/app/manifest/validation"
)

type ActionSetting struct {
	Id   string                   `yaml:"id"`
	Type ActionType               `yaml:"type"`
	Mock *HttpResponseMockSetting `yaml:"mock"`
}

type ActionType string

const (
	ActionTypeHttpRequest ActionType = "httpRequest"
)

func (setting ActionSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if setting.Mock != nil {
		result.AppendValidable(setting.Mock)
	}

	if len(setting.Id) == 0 {
		result.AddError("actions.id is required")
	}

	if len(setting.Type) == 0 {
		var msg = fmt.Sprintf("actions[%s].type is required", setting.Id)
		result.AddError(msg)
	} else {
		if (setting.Type == ActionTypeHttpRequest) == false {

			var msg = fmt.Sprintf("actions[%s].type should contain valid value (httpRequest)", setting.Id)
			result.AddError(msg)
		}
	}

	return result
}
