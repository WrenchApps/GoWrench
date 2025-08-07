package auth

import (
	"context"
	"crypto/sha256"
	"fmt"
	contexts "wrench/app/contexts"
	"wrench/app/cross_funcs"
	"wrench/app/manifest/api_settings"
)

func HMACValidate(ctx context.Context, authorizationSettings *api_settings.AuthorizationSettings) bool {
	var data = ""
	bodyContext := new(contexts.BodyContext)
	wrenchContext := new(contexts.WrenchContext)
	for _, item := range authorizationSettings.ConcatenateFields {
		data += fmt.Sprint(contexts.GetCalculatedValue(item, wrenchContext, bodyContext, nil))
	}

	hash := cross_funcs.GetHash(authorizationSettings.Key, sha256.New, []byte(data))

	return hash == authorizationSettings.SignatureRef
}
