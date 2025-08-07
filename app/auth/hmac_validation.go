package auth

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	contexts "wrench/app/contexts"
	"wrench/app/cross_funcs"
	"wrench/app/manifest/api_settings"
)

func HMACValidate(r *http.Request, authorizationSettings *api_settings.AuthorizationSettings) bool {
	var data = ""
	bodyContext := new(contexts.BodyContext)
	wrenchContext := new(contexts.WrenchContext)
	wrenchContext.Request = r

	body, err := readAndReplaceBody(r)
	if err != nil {
		return false
	}
	bodyContext.SetBody(body)

	for _, item := range authorizationSettings.ConcatenateFields {
		data += fmt.Sprint(contexts.GetCalculatedValue(item, wrenchContext, bodyContext, nil))
	}

	expectedHash := contexts.GetCalculatedValue(authorizationSettings.SignatureRef, wrenchContext, bodyContext, nil)

	hashFn := cross_funcs.GetHashFunc(authorizationSettings.Algorithm)
	actualHash := cross_funcs.GetHash(authorizationSettings.Key, hashFn, []byte(data))

	return expectedHash == actualHash
}

func readAndReplaceBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes, nil
}
