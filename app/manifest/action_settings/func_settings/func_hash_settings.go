package func_settings

import (
	"wrench/app/manifest/validation"
)

type FuncHashAlg string

const (
	FuncHashAlgSHA512 FuncHashAlg = "SHA-512"
	FuncHashAlgSHA256 FuncHashAlg = "SHA-256"
	FuncHashAlgSHA1   FuncHashAlg = "SHA-1"
	FuncHashAlgMD5    FuncHashAlg = "MD5"
)

type FuncHashSettings struct {
	Alg FuncHashAlg `yaml:"alg"`
	Key string      `yaml:"key"`
}

func (setting FuncHashSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Alg) == 0 {
		result.AddError("actions.hash.alg is required")
	} else {
		if (setting.Alg == FuncHashAlgSHA512 ||
			setting.Alg == FuncHashAlgSHA1 ||
			setting.Alg == FuncHashAlgMD5) == false {
			result.AddError("actions.hash.alg should contain valid values SHA-512, SHA-1 or MD5")
		}
	}

	if len(setting.Key) == 0 {
		result.AddError("actions.hash.key is required")
	}

	return result
}
