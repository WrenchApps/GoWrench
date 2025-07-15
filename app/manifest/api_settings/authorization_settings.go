package api_settings

import "wrench/app/manifest/validation"

type AuthorizationSettings struct {
	Type      AuthorizationType `yaml:"type"`
	JwksUrl   string            `yaml:"jwksUrl"`
	Algorithm string            `yaml:"algorithm"`
	Kid       string            `yaml:"kid"`
}

type AuthorizationType string

const (
	JWKSAuthorizationType AuthorizationType = "jwks"
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

	if setting.Type != JWKSAuthorizationType {
		result.AddError("api.authorization.type should be a valid type (jwks)")
	}

	return result
}
