package maps

import (
	"wrench/app/manifest/validation"
)

type ParseSettings struct {
	WhenEquals []string `yaml:"whenEquals"`
	ToArray    []string `yaml:"toArray"`
}

func (setting ParseSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.WhenEquals) == 0 &&
		len(setting.ToArray) == 0 {
		result.AddError("contract.maps.parse should configure whenEquals or toArray")
	}

	return result
}
