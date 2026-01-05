package cross_validation

import (
	"fmt"
	"wrench/app/manifest/application_settings"
	"wrench/app/manifest/rate_limit_settings"
	"wrench/app/manifest/validation"
	"wrench/app/manifest_cross_funcs"
)

func rateLimitCrossValidation(appSetting *application_settings.ApplicationSettings) validation.ValidateResult {
	var result validation.ValidateResult

	rateLimits := appSetting.RateLimits

	result.Append(rateLimitDuplicatedId(rateLimits))
	result.Append(rateLimitRedisConnectionId(rateLimits))

	return result
}

func rateLimitDuplicatedId(rateLimits []*rate_limit_settings.RateLimitSettings) validation.ValidateResult {
	var result validation.ValidateResult

	if len(rateLimits) > 0 {
		hasIds := toHasIdSlice(rateLimits)
		duplicateIds := duplicateIdsValid(hasIds)

		for _, id := range duplicateIds {
			result.AddError(fmt.Sprintf("rateLimits.id %v duplicated", id))
		}
	}

	return result
}

func rateLimitRedisConnectionId(rateLimits []*rate_limit_settings.RateLimitSettings) validation.ValidateResult {
	var result validation.ValidateResult

	if len(rateLimits) > 0 {
		for _, rateLimit := range rateLimits {
			if len(rateLimit.RedisConnectionId) > 0 {
				_, err := manifest_cross_funcs.GetConnectionRedisSettingById(rateLimit.RedisConnectionId)

				if err != nil {
					result.AddError(fmt.Sprintf("rateLimits[%v].redisConnectionId %v don't exist in connections.redis", rateLimit.Id, rateLimit.RedisConnectionId))
				}
			}
		}
	}

	return result
}
