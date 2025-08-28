package auth

import (
	"fmt"
	"wrench/app/contexts"
	"wrench/app/cross_funcs"
	"wrench/app/manifest/api_settings"
)

func HMACValidate(wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext, authorizationSettings *api_settings.AuthorizationSettings) bool {
	var data = ""

	for _, item := range authorizationSettings.ConcatenateFields {
		data += fmt.Sprint(contexts.GetCalculatedValue(item, wrenchContext, bodyContext, nil))
	}

	expectedHash := contexts.GetCalculatedValue(authorizationSettings.SignatureRef, wrenchContext, bodyContext, nil)

	hashFn := cross_funcs.GetHashFunc(authorizationSettings.Algorithm)
	actualHash := cross_funcs.GetHash(authorizationSettings.Key, hashFn, []byte(data))

	return expectedHash == actualHash
}
