package connection_settings

import (
	"wrench/app/manifest/validation"
)

type KafkaConnectionType string

const (
	KafkaConnectionPlaintext KafkaConnectionType = "plaintext"
	KafkaConnectionSsl       KafkaConnectionType = "ssl"
)

type KafkaConnectionSettings struct {
	Id               string              `yaml:"id"`
	BootstrapServers string              `yaml:"bootstrapServers"`
	ConnectionType   KafkaConnectionType `yaml:"connectionType"`
}

func (setting *KafkaConnectionSettings) GetId() string {
	return setting.Id
}

func (setting KafkaConnectionSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Id) == 0 {
		result.AddError("connections.kafka.id is required")
	}

	if len(setting.BootstrapServers) == 0 {
		result.AddError("connections.kafka.bootstrapServers is required")
	}

	if len(setting.ConnectionType) == 0 {
		result.AddError("connections.kafka.connectionType is required")
	} else {
		if (setting.ConnectionType == KafkaConnectionPlaintext ||
			setting.ConnectionType == KafkaConnectionSsl) == false {
			result.AddError("connections.kafka.connectionType should be plaintext or ssl")
		}
	}

	return result
}
