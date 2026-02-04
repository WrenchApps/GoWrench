package action_settings

import (
	"fmt"
	"wrench/app/manifest/action_settings/dynamodb_settings"
	"wrench/app/manifest/action_settings/file_settings"
	"wrench/app/manifest/action_settings/func_settings"
	"wrench/app/manifest/action_settings/http_settings"
	"wrench/app/manifest/action_settings/kafka_settings"
	"wrench/app/manifest/action_settings/nats_settings"
	"wrench/app/manifest/action_settings/sns_settings"
	"wrench/app/manifest/action_settings/trigger_settings"
	"wrench/app/manifest/validation"
)

type ActionSettings struct {
	Id                      string                              `yaml:"id"`
	Type                    ActionType                          `yaml:"type"`
	Http                    *http_settings.HttpSettings         `yaml:"http"`
	SNS                     *sns_settings.SnsSettings           `yaml:"sns"`
	Trigger                 *trigger_settings.TriggerSetting    `yaml:"trigger"`
	File                    *file_settings.FileSettings         `yaml:"file"`
	Nats                    *nats_settings.NatsSettings         `yaml:"nats"`
	Kafka                   *kafka_settings.KafkaSettings       `yaml:"kafka"`
	Func                    *func_settings.FuncSettings         `yaml:"func"`
	DynamoDb                *dynamodb_settings.DynamoDbSettings `yaml:"dynamodb"`
	Body                    *BodyActionSettings                 `yaml:"body"`
	RunEvenThereIsFlowError bool                                `yaml:"runEvenThereIsFlowError"`
}

func (setting *ActionSettings) GetId() string {
	return setting.Id
}

type ActionType string

const (
	ActionTypeHttpRequest           ActionType = "httpRequest"
	ActionTypeHttpRequestMock       ActionType = "httpRequestMock"
	ActionTypeSnsPublish            ActionType = "snsPublish"
	ActionTypeFileReader            ActionType = "fileReader"
	ActionTypeNatsPublish           ActionType = "natsPublish"
	ActionTypeKafkaProducer         ActionType = "kafkaProducer"
	ActionTypeFuncHash              ActionType = "funcHash"
	ActionTypeFuncSignature         ActionType = "funcSignature"
	ActionTypeFuncVarContext        ActionType = "funcVarContext"
	ActionTypeFuncStringConcatenate ActionType = "funcStringConcatenate"
	ActionTypeFuncGeneral           ActionType = "funcGeneral"
	ActionTypeDynamoDb              ActionType = "dynamodb"
)

func (setting *ActionSettings) Valid() validation.ValidateResult {
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
			setting.Type == ActionTypeKafkaProducer ||
			setting.Type == ActionTypeFuncHash ||
			setting.Type == ActionTypeFuncSignature ||
			setting.Type == ActionTypeFuncVarContext ||
			setting.Type == ActionTypeFuncStringConcatenate ||
			setting.Type == ActionTypeFuncGeneral ||
			setting.Type == ActionTypeDynamoDb) == false {

			var msg = fmt.Sprintf("actions[%s].type should contain valid value", setting.Id)
			result.AddError(msg)
		}
	}

	if setting.Http != nil {
		result.AppendValidable(setting.Http)
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

	if setting.DynamoDb != nil {
		result.AppendValidable(setting.DynamoDb)
	}

	result.Append(setting.ActionTypeKafkaProducerValid())
	result.Append(setting.checkTypes())

	return result
}

func (setting *ActionSettings) ShouldPreserveBody() bool {
	return setting.Body != nil && setting.Body.PreserveCurrentBody
}

func (setting *ActionSettings) ShouldUseBodyRef() (shouldUse bool, valueRef string) {
	bodyConfig := setting.Body
	if bodyConfig == nil {
		return false, ""
	} else {
		return len(bodyConfig.Use) > 0, bodyConfig.Use
	}
}

func (setting *ActionSettings) ActionTypeKafkaProducerValid() validation.ValidateResult {
	var result validation.ValidateResult

	if setting.Type == ActionTypeKafkaProducer {
		if setting.Kafka == nil {
			result.AddError("action.kafka is required when action type is kafkaProducer")
		} else {
			result.AppendValidable(setting.Kafka)
		}

	}

	return result
}

func (setting *ActionSettings) checkTypes() validation.ValidateResult {
	var result validation.ValidateResult

	if setting.Type == ActionTypeHttpRequest && setting.Http == nil {
		result.AddError(fmt.Sprintf("actions[%v].http is required when type is %v", setting.Id, setting.Type))
	}

	if setting.Type == ActionTypeSnsPublish && setting.SNS == nil {
		result.AddError(fmt.Sprintf("actions[%v].sns is required when type is %v", setting.Id, setting.Type))
	}

	if setting.Type == ActionTypeFileReader && setting.File == nil {
		result.AddError(fmt.Sprintf("actions[%v].file is required when type is %v", setting.Id, setting.Type))
	}

	if setting.Type == ActionTypeNatsPublish && setting.Nats == nil {
		result.AddError(fmt.Sprintf("actions[%v].nats is required when type is %v", setting.Id, setting.Type))
	}

	if setting.Type == ActionTypeDynamoDb && setting.DynamoDb == nil {
		result.AddError(fmt.Sprintf("actions[%v].dynamodb is required when type is %v", setting.Id, setting.Type))
	}

	if (setting.Type == ActionTypeFuncVarContext ||
		setting.Type == ActionTypeFuncStringConcatenate ||
		setting.Type == ActionTypeFuncHash ||
		setting.Type == ActionTypeFuncGeneral) && setting.Func == nil {

		if setting.Func == nil {
			result.AddError(fmt.Sprintf("actions[%v].func is required when type is %v", setting.Id, setting.Type))
		} else {
			if setting.Type == ActionTypeFuncVarContext && len(setting.Func.Vars) == 0 {
				result.AddError(fmt.Sprintf("actions[%v].func.vars is required when type is %v", setting.Id, setting.Type))
			} else if setting.Type == ActionTypeFuncStringConcatenate && len(setting.Func.Concatenate) == 0 {
				result.AddError(fmt.Sprintf("actions[%v].func.concatenate is required when type is %v", setting.Id, setting.Type))
			} else if setting.Type == ActionTypeFuncHash && setting.Func.Hash == nil {
				result.AddError(fmt.Sprintf("actions[%v].func.hash is required when type is %v", setting.Id, setting.Type))
			} else if setting.Type == ActionTypeFuncSignature && setting.Func.Sign == nil {
				result.AddError(fmt.Sprintf("actions[%v].func.sign is required when type is %v", setting.Id, setting.Type))
			} else if setting.Type == ActionTypeFuncGeneral && len(setting.Func.Command) == 0 {
				result.AddError(fmt.Sprintf("actions[%v].func.command is required when type is %v", setting.Id, setting.Type))
			}
		}
	}

	return result
}
