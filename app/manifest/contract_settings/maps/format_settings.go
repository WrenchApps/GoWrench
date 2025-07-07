package maps

import (
	"strings"
	"wrench/app/manifest/validation"
)

type FormatSettings struct {
	Date []string `yaml:"date"`
}

func (setting FormatSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Date) == 0 {
		result.AddError("contract.maps.format should configure Date")
	} else {
		for _, property := range setting.Date {
			errorSplitted := "contract.maps.format should be configured as 'datePropertyName:dateFormat' without spaces"
			if strings.Contains(property, " ") {
				result.AddError(errorSplitted)
			}

			propertySplitted := strings.Split(property, ":")
			if len(propertySplitted) != 2 {
				result.AddError(errorSplitted)
			}
		}
	}

	return result
}
