package func_settings

import (
	"wrench/app/manifest/validation"
)

type FuncGeneralType string

const (
	FuncTypeTimestampMilli  FuncGeneralType = "func.timestamp(milli)"
	FuncTypeBase64Encode    FuncGeneralType = "func.base64(encode)"
	FuncTypeBase64UrlEncode FuncGeneralType = "func.base64Url(encode)"
)

type FuncSettings struct {
	Hash        *FuncHashSettings      `yaml:"hash"`
	Sign        *FuncSignatureSettings `yaml:"sign"`
	Vars        map[string]string      `yaml:"vars"`
	Concatenate []string               `yaml:"concatenate"`
	Command     FuncGeneralType        `yaml:"command"`
}

func (setting FuncSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if setting.Hash != nil {
		result.AppendValidable(setting.Hash)
	}

	if setting.Sign != nil {
		result.AppendValidable(setting.Sign)
	}

	if len(setting.Command) > 0 {
		if (string(setting.Command) == "{{"+string(FuncTypeTimestampMilli)+"}}" ||
			string(setting.Command) == "{{"+string(FuncTypeBase64Encode)+"}}" ||
			string(setting.Command) == "{{"+string(FuncTypeBase64UrlEncode)+"}}") == false {
			result.AddError("actions.func.command is invalid")
		}
	}

	return result
}
