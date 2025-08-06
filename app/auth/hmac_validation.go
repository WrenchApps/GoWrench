package auth

import (
	"context"
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

	hashFn := cross_funcs.GetHashFunc(authorizationSettings.Algorithm)
	hash := cross_funcs.GetHash(authorizationSettings.Key, hashFn, []byte(data))

	return hash == authorizationSettings.SignatureRef
}
