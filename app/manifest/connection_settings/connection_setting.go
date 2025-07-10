package connection_settings

import (
	"wrench/app/manifest/validation"
)

type ConnectionSettings struct {
	Nats  []*NatsConnectionSettings  `yaml:"nats"`
	Kafka []*KafkaConnectionSettings `yaml:"kafka"`
	Redis []*RedisConnectionSettings `yaml:"redis"`
}

func (settings *ConnectionSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(settings.Nats) > 0 {
		for _, validable := range settings.Nats {
			result.AppendValidable(validable)
		}
	}

	if len(settings.Kafka) > 0 {
		for _, validable := range settings.Kafka {
			result.AppendValidable(validable)
		}
	}

	if len(settings.Redis) > 0 {
		for _, validable := range settings.Redis {
			result.AppendValidable(validable)
		}
	}

	return result
}

func (settings *ConnectionSettings) Merge(toMerge *ConnectionSettings) error {

	if toMerge == nil {
		return nil
	}

	if len(toMerge.Nats) > 0 {
		if len(settings.Nats) == 0 {
			settings.Nats = toMerge.Nats
		} else {
			settings.Nats = append(settings.Nats, toMerge.Nats...)
		}
	}

	if len(toMerge.Redis) > 0 {
		if len(settings.Redis) == 0 {
			settings.Redis = toMerge.Redis
		} else {
			settings.Redis = append(settings.Redis, toMerge.Redis...)
		}
	}

	if len(toMerge.Kafka) > 0 {
		if len(settings.Kafka) == 0 {
			settings.Kafka = toMerge.Kafka
		} else {
			settings.Kafka = append(settings.Kafka, toMerge.Kafka...)
		}
	}

	return nil
}
