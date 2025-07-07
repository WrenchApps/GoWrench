package connection_settings

import "wrench/app/manifest/validation"

type RedisConnectionSettings struct {
	Id        string   `yaml:"id"`
	Address   string   `yaml:"address"`
	Addresses []string `yaml:"addresses"`
	IsCluster bool     `yaml:"isCluster"`
	Password  string   `yaml:"password"`
	Db        int      `yaml:"db"`
}

func (setting *RedisConnectionSettings) GetId() string {
	return setting.Id
}

func (settings RedisConnectionSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(settings.Id) == 0 {
		result.AddError("connections.redis.id is required")
	}

	if settings.IsCluster {
		if len(settings.Addresses) == 0 {
			result.AddError("the connections.redis.addresses is required when is cluster")
		}
	} else {
		if len(settings.Address) == 0 {
			result.AddError("the connections.redis.address is required when is not cluster")
		}
	}

	return result
}
