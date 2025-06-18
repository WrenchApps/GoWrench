package action_settings

import "wrench/app/manifest/validation"

type BodyActionSettings struct {
	PreserveCurrentBody bool   `yaml:"preserveCurrentBody"`
	Use                 string `yaml:"use"`
}

func (setting BodyActionSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	return result
}
