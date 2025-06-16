package func_settings

import "wrench/app/manifest/validation"

type FuncSettings struct {
	Hash    *FuncHashSettings `yaml:"hash"`
	Vars    map[string]string `yaml:"vars"`
	Concate []string          `yaml:"concate"`
}

func (setting FuncSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if setting.Hash != nil {
		result.AppendValidable(setting.Hash)
	}

	return result
}
