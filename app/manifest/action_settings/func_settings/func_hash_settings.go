package func_settings

import (
	"wrench/app/manifest/types"
	"wrench/app/manifest/validation"
)

type FuncHashSettings struct {
	Alg types.HashAlg `yaml:"alg"`
	Key string        `yaml:"key"`
}

func (setting FuncHashSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Alg) == 0 {
		result.AddError("actions.func.hash.alg is required")
	} else {
		if (setting.Alg == types.HashAlgSHA512 ||
			setting.Alg == types.HashAlgSHA256 ||
			setting.Alg == types.HashAlgSHA1 ||
			setting.Alg == types.HashAlgMD5) == false {
			result.AddError("actions.func.hash.alg should contain valid values SHA-512, SHA-256, SHA-1 or MD5")
		}
	}

	if len(setting.Key) == 0 {
		result.AddError("actions.func.hash.key is required")
	}

	return result
}
