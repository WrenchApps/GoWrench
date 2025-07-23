package func_settings

import (
	"wrench/app/manifest/validation"
)

type FuncGeneralType string

const (
	FuncTypeTimestampMilli FuncGeneralType = "func.timestamp(milli)"
	FuncTypeBase64Encode   FuncGeneralType = "func.base64(encode)"
	FuncTypeCurrentDate    FuncGeneralType = "func.currentDate(utc)"
)

type FuncSettings struct {
	Hash        *FuncHashSettings `yaml:"hash"`
	Vars        map[string]string `yaml:"vars"`
	Concatenate []string          `yaml:"concatenate"`
	Command     FuncGeneralType   `yaml:"command"`
}

func (setting FuncSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if setting.Hash != nil {
		result.AppendValidable(setting.Hash)
	}

	if len(setting.Command) > 0 {
		if string(setting.Command) == "{{"+string(FuncTypeTimestampMilli)+"}}" ||
			string(setting.Command) == "{{"+string(FuncTypeCurrentDate)+"}}" ||
			string(setting.Command) == "{{"+string(FuncTypeBase64Encode)+"}}" == false {
			result.AddError("actions.func.command is invalid")
		}
	}

	return result
}
