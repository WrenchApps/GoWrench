package manifest_cross_funcs

import (
	"errors"
	"wrench/app/manifest/application_settings"
	"wrench/app/manifest/connection_settings"
	"wrench/app/manifest/token_credential_settings"
)

func GetTokenCredentialSettingById(id string) (*token_credential_settings.TokenCredentialSetting, error) {
	appSetting := application_settings.ApplicationSettingsStatic

	if len(appSetting.TokenCredentials) > 0 {
		for _, token := range appSetting.TokenCredentials {
			if token.Id == id {
				return token, nil
			}
		}
	}

	return nil, errors.New("not found")
}

func GetConnectionKafkaSettingById(kafkaId string) (*connection_settings.KafkaConnectionSettings, error) {
	appSetting := application_settings.ApplicationSettingsStatic

	if appSetting.Connections != nil && len(appSetting.Connections.Kafka) > 0 {
		for _, kafka := range appSetting.Connections.Kafka {
			if kafka.Id == kafkaId {
				return kafka, nil
			}
		}
	}

	return nil, errors.New("not found")
}
