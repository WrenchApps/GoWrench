package rate_limit_settings

import "wrench/app/manifest/validation"

type RateLimitSetting struct {
	Id                string   `yaml:"id"`
	RedisConnectionId string   `yaml:"redisConnectionId"`
	RouteEnabled      bool     `yaml:"routeEnabled"`
	BodyFields        []string `yaml:"bodyFields"`
	Headers           []string `yaml:"headers"`
	Claims            []string `yaml:"claims"`
	RequestsPerSecond int      `yaml:"requestsPerSecond"`
	BurstLimit        int      `yaml:"burstLimit"`
}

func (setting *RateLimitSetting) GetId() string {
	return setting.Id
}

func (setting *RateLimitSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Id) == 0 {
		result.AddError("idemp.id is required")
	}

	if len(setting.BodyFields) == 0 &&
		len(setting.Headers) == 0 &&
		len(setting.Claims) == 0 &&
		!setting.RouteEnabled {
		result.AddError("should set at least one of bodyFields, headers, claims or enable routeEnabled")
	}

	if setting.RequestsPerSecond >= 0 {
		result.AddError("should set requestsPerSecond")
	}

	return result
}
