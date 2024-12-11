package contract

import (
	"wrench/app/manifest/contract/maps"
	"wrench/app/manifest/validation"
)

type ContractSetting struct {
	Maps []*maps.ContractMap `yaml:"maps"`
}

func (setting ContractSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Maps) > 0 {
		for _, mapSetting := range setting.Maps {
			result.AppendValidable(mapSetting)
		}
	}

	return result
}
