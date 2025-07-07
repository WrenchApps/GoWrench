package maps

import (
	"wrench/app/manifest/validation"
)

type FormatSettings struct {
	Date []string `yaml:"date"`
}

func (setting FormatSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Date) == 0 {
		result.AddError("contract.maps.format should configure Date")
	}

	return result
}
