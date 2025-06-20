package connection_settings

import "wrench/app/manifest/validation"

type ConnectionSettings struct {
	Nats []*NatsConnectionSettings `yaml:"nats"`
}

func (settings ConnectionSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(settings.Nats) > 0 {
		for _, validable := range settings.Nats {
			result.AppendValidable(validable)
		}
	}

	return result
}
