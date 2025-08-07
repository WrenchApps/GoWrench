package api_settings

import (
	"wrench/app/manifest/types"
	"wrench/app/manifest/validation"
)

type AuthorizationSettings struct {
	Type              AuthorizationType `yaml:"type"`
	JwksUrl           string            `yaml:"jwksUrl"`
	Algorithm         types.HashAlg     `yaml:"algorithm"`
	Kid               string            `yaml:"kid"`
	Key               string            `yaml:"key"`
	SignatureRef      string            `yaml:"signatureRef"`
	ConcatenateFields []string          `yaml:"concatenateFields"`
}

type AuthorizationType string

const (
	JWKSAuthorizationType AuthorizationType = "jwks"
	HMACAuthorizationType AuthorizationType = "hmac"
)

func (setting AuthorizationSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if setting.Type == JWKSAuthorizationType {
		if setting.JwksUrl == "" {
			result.AddError("api.authorization.jwksUrl is required when type is jwks")
		}

		if setting.Algorithm == "" {
			result.AddError("api.authorization.algorithm is required when type is jwks")
		}
	}

	if setting.Type == HMACAuthorizationType {
		if setting.Algorithm == "" {
			result.AddError("api.authorization.algorithm is required when type is hmac")
		}

		if setting.Key == "" {
			result.AddError("api.authorization.key is required when type is hmac")
		}

		if setting.SignatureRef == "" {
			result.AddError("api.authorization.signatureRef is required when type is hmac")
		}

		if len(setting.ConcatenateFields) == 0 {
			result.AddError("api.authorization.concatenateFields is required when type is hmac")
		}
	}

	if setting.Type != JWKSAuthorizationType &&
		setting.Type != HMACAuthorizationType {
		result.AddError("api.authorization.type should be a valid type (jwks or hmac)")
	}

	return result
}
