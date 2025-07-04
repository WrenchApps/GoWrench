package idemp_settings

import "wrench/app/manifest/validation"

type IdempSettings struct {
	Id                string `yaml:"id"`
	RedisConnectionId string `yaml:"redisConnectionId"`
	Key               string `yaml:"key"`
	TtlInSeconds      int    `yaml:"ttlInSeconds"`
}

func (setting *IdempSettings) GetId() string {
	return setting.Id
}

func (setting *IdempSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Id) == 0 {
		result.AddError("idemp.id is required")
	}

	if len(setting.RedisConnectionId) == 0 {
		result.AddError("idemp.redisConnectionId is required")
	}

	if len(setting.Key) == 0 {
		result.AddError("idemp.key is required")
	}

	if setting.TtlInSeconds < 60 {
		result.AddError("idemp.ttlInSeconds should be greater than 59")
	}

	return result
}
