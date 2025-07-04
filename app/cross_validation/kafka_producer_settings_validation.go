package cross_validation

import (
	"fmt"
	"wrench/app/manifest/action_settings"
	"wrench/app/manifest/application_settings"
	"wrench/app/manifest/validation"
	"wrench/app/manifest_cross_funcs"
)

func kafkaProducerCronsValidation(appSetting *application_settings.ApplicationSettings) validation.ValidateResult {
	var result validation.ValidateResult

	actionKafkaRequest := getActionsByType(appSetting.Actions, action_settings.ActionTypeKafkaProducer)

	if len(actionKafkaRequest) > 0 {
		for _, action := range actionKafkaRequest {
			if len(action.Kafka.ConnectionId) > 0 {
				_, err := manifest_cross_funcs.GetConnectionKafkaSettingById(action.Kafka.ConnectionId)

				if err != nil {
					result.AddError(fmt.Sprintf("actions.kafka.connectionId %v don't exist in connections.kafka", action.Kafka.ConnectionId))
				}
			}
		}
	}

	return result
}
