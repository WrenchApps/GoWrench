package action_settings

import "wrench/app/manifest/validation"

type BodyActionSettings struct {
	PreserveCurrentBody bool   `yaml:"preserveCurrentBody"`
	UseValue            string `yaml:"useValue"`
}

func (setting BodyActionSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	return result
}
