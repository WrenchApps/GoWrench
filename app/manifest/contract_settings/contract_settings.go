package contract_settings

import (
	"wrench/app/manifest/contract_settings/maps"
	"wrench/app/manifest/validation"
)

type ContractSetting struct {
	Maps []*maps.ContractMapSetting `yaml:"maps"`
}

func (setting *ContractSetting) GetContractById(contractMapId string) *maps.ContractMapSetting {
	if setting.Maps == nil {
		return nil
	}

	var contractMap *maps.ContractMapSetting = nil

	for _, contract := range setting.Maps {
		if contract.Id == contractMapId {
			contractMap = contract
			break
		}
	}

	return contractMap
}

func (setting *ContractSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Maps) > 0 {
		for _, mapSetting := range setting.Maps {
			result.AppendValidable(mapSetting)
		}
	}

	return result
}

func (settings *ContractSetting) Merge(toMerge *ContractSetting) error {
	if len(toMerge.Maps) > 0 {
		if len(settings.Maps) == 0 {
			settings.Maps = toMerge.Maps
		} else {
			settings.Maps = append(settings.Maps, toMerge.Maps...)
		}
	}

	return nil
}
