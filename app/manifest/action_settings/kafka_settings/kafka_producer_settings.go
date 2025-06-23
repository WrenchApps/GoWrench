package kafka_settings

import "wrench/app/manifest/validation"

type KafkaSettings struct {
	ConnectionId string            `yaml:"connectionId"`
	TopicName    string            `yaml:"topicName"`
	Headers      map[string]string `yaml:"headers"`
}

func (setting KafkaSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.ConnectionId) == 0 {
		result.AddError("actions.kafka.connectionId is required")
	}

	if len(setting.TopicName) == 0 {
		result.AddError("actions.kafka.topicName is required")
	}

	return result
}
