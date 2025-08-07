package manifest_cross_funcs

import (
	"errors"
	"wrench/app/manifest/api_settings"
	"wrench/app/manifest/application_settings"
	"wrench/app/manifest/connection_settings"
	"wrench/app/manifest/idemp_settings"
	"wrench/app/manifest/service_settings"
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

	return nil, errors.New("token credential not found")
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

	return nil, errors.New("kafka not found")
}

func GetConnectionRedisSettingById(redisConnectionId string) (*connection_settings.RedisConnectionSettings, error) {
	appSetting := application_settings.ApplicationSettingsStatic

	if appSetting.Connections != nil && len(appSetting.Connections.Redis) > 0 {
		for _, redis := range appSetting.Connections.Redis {
			if redis.Id == redisConnectionId {
				return redis, nil
			}
		}
	}

	return nil, errors.New("redis not found")
}

func GetIdempSettingById(idempId string) (*idemp_settings.IdempSettings, error) {
	appSetting := application_settings.ApplicationSettingsStatic

	if len(appSetting.Idemps) > 0 {
		for _, idemp := range appSetting.Idemps {
			if idemp.Id == idempId {
				return idemp, nil
			}
		}
	}

	return nil, errors.New("idemp not found")
}

func GetService() *service_settings.ServiceSettings {
	appSetting := application_settings.ApplicationSettingsStatic
	return appSetting.Service
}

func GetAuthorizationSettings() *api_settings.AuthorizationSettings {
	appSetting := application_settings.ApplicationSettingsStatic

	if appSetting == nil || appSetting.Api == nil || appSetting.Api.Authorization == nil {
		return nil
	}

	return appSetting.Api.Authorization
}
