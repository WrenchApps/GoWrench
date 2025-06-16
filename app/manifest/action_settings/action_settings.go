package action_settings

import (
	"fmt"
	"wrench/app/manifest/action_settings/file_settings"
	"wrench/app/manifest/action_settings/func_settings"
	"wrench/app/manifest/action_settings/http_settings"
	"wrench/app/manifest/action_settings/nats_settings"
	"wrench/app/manifest/action_settings/sns_settings"
	"wrench/app/manifest/action_settings/trigger_settings"
	"wrench/app/manifest/validation"
)

type ActionSettings struct {
	Id      string                           `yaml:"id"`
	Type    ActionType                       `yaml:"type"`
	Http    *http_settings.HttpSetting       `yaml:"http"`
	SNS     *sns_settings.SnsSettings        `yaml:"sns"`
	Trigger *trigger_settings.TriggerSetting `yaml:"trigger"`
	File    *file_settings.FileSettings      `yaml:"file"`
	Nats    *nats_settings.NatsSettings      `yaml:"nats"`
	Func    *func_settings.FuncSettings      `yaml:"func"`
	Body    *BodyActionSettings              `yaml:"body"`
}

type ActionType string

const (
	ActionTypeHttpRequest       ActionType = "httpRequest"
	ActionTypeHttpRequestMock   ActionType = "httpRequestMock"
	ActionTypeSnsPublish        ActionType = "snsPublish"
	ActionTypeFileReader        ActionType = "fileReader"
	ActionTypeNatsPublish       ActionType = "natsPublish"
	ActionTypeFuncHash          ActionType = "funcHash"
	ActionTypeFuncVarContext    ActionType = "funcVarContext"
	ActionTypeFuncStringConcate ActionType = "funcStringConcate"
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
			setting.Type == ActionTypeHttpRequestMock ||
			setting.Type == ActionTypeSnsPublish ||
			setting.Type == ActionTypeFileReader ||
			setting.Type == ActionTypeNatsPublish ||
			setting.Type == ActionTypeFuncHash ||
			setting.Type == ActionTypeFuncStringConcate) == false {

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

	if setting.SNS != nil {
		result.AppendValidable(setting.SNS)
	}

	if setting.Trigger != nil {
		result.AppendValidable(setting.Trigger)
	}

	if setting.File != nil {
		result.AppendValidable(setting.File)
	}

	if setting.Nats != nil {
		result.AppendValidable(setting.Nats)
	}

	if setting.Func != nil {
		result.AppendValidable(setting.Func)
	}

	return result
}

func (setting ActionSettings) ShouldPreserveBody() bool {
	return setting.Body != nil && setting.Body.PreserveCurrentBody
}

func (setting ActionSettings) ShouldUseBodyRef() (shouldUse bool, valueRef string) {
	bodyConfig := setting.Body
	if bodyConfig == nil {
		return false, ""
	} else {
		return len(bodyConfig.UseValue) > 0, bodyConfig.UseValue
	}
}
