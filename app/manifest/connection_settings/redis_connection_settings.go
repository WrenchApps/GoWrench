package connection_settings

import "wrench/app/manifest/validation"

type RedisConnectionSettings struct {
	Id       string `yaml:"id"`
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
}

func (setting *RedisConnectionSettings) GetId() string {
	return setting.Id
}

func (settings RedisConnectionSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(settings.Id) == 0 {
		result.AddError("connections.redis.id is required")
	}

	if len(settings.Address) == 0 {
		result.AddError("the connections.redis.address is required")
	}

	return result
}
