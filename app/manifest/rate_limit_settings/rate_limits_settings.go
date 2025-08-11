package rate_limit_settings

import "wrench/app/manifest/validation"

type RateLimitSettings struct {
	Id                string   `yaml:"id"`
	RedisConnectionId string   `yaml:"redisConnectionId"`
	RouteEnabled      bool     `yaml:"routeEnabled"`
	Keys              []string `yaml:"keys"`
	RequestsPerSecond int      `yaml:"requestsPerSecond"`
	RequestsPerMinute int      `yaml:"requestsPerMinute"`
	BurstLimit        int      `yaml:"burstLimit"`
}

func (setting *RateLimitSettings) GetId() string {
	return setting.Id
}

func (setting *RateLimitSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Id) == 0 {
		result.AddError("rateLimits.id is required")
	}

	if len(setting.Keys) == 0 &&
		!setting.RouteEnabled {
		result.AddError("should set at least one of keys or enable routeEnabled")
	}

	if setting.RedisConnectionId == "" {
		result.AddError("should set redisConnectionId")
	}

	if setting.RequestsPerSecond >= 0 && setting.RequestsPerMinute >= 0 {
		result.AddError("should set requestsPerSecond or requestsPerMinute, not both")
	}

	return result
}
