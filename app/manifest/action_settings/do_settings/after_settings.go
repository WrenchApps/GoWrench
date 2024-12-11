package do_settings

import (
	"fmt"
	"wrench/app/cross_cutting"
	"wrench/app/manifest/validation"
)

type AfterSetting struct {
	ContractMapId string `yaml:"contractMapId"`
}

func (setting AfterSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	contractMap := cross_cutting.GetContractById(setting.ContractMapId)

	if contractMap == nil {
		result.AddError(fmt.Sprintf("actions.http.request.do.after.contractMapId %s should be configured", setting.ContractMapId))
	}

	return result
}
