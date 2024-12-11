package maps

import (
	"strings"
	"wrench/app/manifest/validation"
)

type ContractMap struct {
	Id         string   `yaml:"id"`
	Properties []string `yaml:"properties"`
	Remove     []string `yaml:"remove"`
}

func (setting ContractMap) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Id) <= 0 {
		result.AddError("contract.maps.id is required")
	}

	if len(setting.Properties) > 0 {
		errorSplitted := "contract.maps.properties should be configured looks like 'propertySource:propertyDestination' without space"
		for _, property := range setting.Properties {

			if strings.Contains(property, " ") {
				result.AddError(errorSplitted)
			}

			propertySplitted := strings.Split(property, ":")
			if len(propertySplitted) != 2 {
				result.AddError(errorSplitted)
			}
		}
	}

	if len(setting.Remove) > 0 {
		for _, remove := range setting.Remove {

			if strings.Contains(remove, " ") {
				result.AddError("contract.maps.remove can't contain space")
			}
		}
	}

	return result
}
