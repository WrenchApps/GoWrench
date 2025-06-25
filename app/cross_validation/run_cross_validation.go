package cross_validation

import (
	"fmt"
	"wrench/app/manifest/application_settings"
	"wrench/app/manifest/validation"
)

func Valid() validation.ValidateResult {
	appSetting := application_settings.ApplicationSettingsStatic

	var result validation.ValidateResult

	result.Append(httpRequestCrossValid(appSetting))
	result.Append(kafkaProducerCronsValidation(appSetting))

	if len(appSetting.Actions) > 0 {
		hasIds := toHasIdSlice(appSetting.Actions)
		duplicateIds := duplicateIdsValid(hasIds)

		for _, id := range duplicateIds {
			result.AddError(fmt.Sprintf("actions.id %v duplicated", id))
		}
	}

	if len(appSetting.TokenCredentials) > 0 {
		hasIds := toHasIdSlice(appSetting.TokenCredentials)
		duplicateIds := duplicateIdsValid(hasIds)

		for _, id := range duplicateIds {
			result.AddError(fmt.Sprintf("tokenCredentials.id %v duplicated", id))
		}
	}

	if appSetting.Connections != nil && len(appSetting.Connections.Nats) > 0 {
		hasIds := toHasIdSlice(appSetting.Connections.Nats)
		duplicateIds := duplicateIdsValid(hasIds)

		for _, id := range duplicateIds {
			result.AddError(fmt.Sprintf("connections.nats.id %v duplicated", id))
		}
	}

	if appSetting.Connections != nil && len(appSetting.Connections.Kafka) > 0 {
		hasIds := toHasIdSlice(appSetting.Connections.Kafka)
		duplicateIds := duplicateIdsValid(hasIds)

		for _, id := range duplicateIds {
			result.AddError(fmt.Sprintf("connections.kafka.id %v duplicated", id))
		}

	}

	return result
}
