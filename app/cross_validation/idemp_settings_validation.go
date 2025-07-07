package cross_validation

import (
	"fmt"
	"wrench/app/manifest/application_settings"
	"wrench/app/manifest/validation"
	"wrench/app/manifest_cross_funcs"
)

func idempCrossValidation(appSetting *application_settings.ApplicationSettings) validation.ValidateResult {
	var result validation.ValidateResult

	idemps := appSetting.Idemps

	if len(idemps) > 0 {
		for _, idemp := range idemps {
			if len(idemp.RedisConnectionId) > 0 {
				_, err := manifest_cross_funcs.GetConnectionRedisSettingById(idemp.RedisConnectionId)

				if err != nil {
					result.AddError(fmt.Sprintf("idemps[%v].redisConnectionId %v don't exist in connections.redis", idemp.Id, idemp.RedisConnectionId))
				}
			}
		}
	}

	return result
}
