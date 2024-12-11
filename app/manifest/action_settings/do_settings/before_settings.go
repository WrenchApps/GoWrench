package do_settings

import (
	"fmt"
	"wrench/app/cross_cutting"
	"wrench/app/manifest/validation"
)

type BeforeSetting struct {
	ContractMapId string `yaml:"contractMapId"`
}

func (setting BeforeSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	contractMap := cross_cutting.GetContractById(setting.ContractMapId)

	if contractMap == nil {
		result.AddError(fmt.Sprintf("actions.http.request.do.before.contractMapId %s should be configured", setting.ContractMapId))
	}

	return result
}
