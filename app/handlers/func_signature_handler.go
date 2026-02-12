package handlers

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
	"wrench/app/manifest/types"
	keys_load "wrench/app/startup/keys"
)

type FuncSignatureHandler struct {
	ActionSettings *settings.ActionSettings
	Next           Handler
}

func (handler *FuncSignatureHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	if (!wrenchContext.HasError || handler.ActionSettings.RunEvenIfFlowHasError) &&
		!wrenchContext.HasCache &&
		!wrenchContext.Unauthorized {

		ctxSpan, span := wrenchContext.GetSpan(ctx, *handler.ActionSettings)
		ctx = ctxSpan
		defer span.End()

		signSetting := handler.ActionSettings.Func.Sign
		priv, err := keys_load.GetPrivateKey(signSetting.KeyId)
		if err != nil {
			wrenchContext.SetHasError3(span, err.Error(), err, 500, bodyContext)
		} else {
			var sig []byte
			body, err := bodyContext.GetBody(handler.ActionSettings)

			if err != nil {
				wrenchContext.SetHasError3(span, err.Error(), err, 500, bodyContext)
			} else {
				if signSetting.Algorithm == types.HashAlgSHA256 {
					sig, err = handler.signBodyRSA_SHA256(priv, body)
				} else {
					wrenchContext.SetHasError3(span, fmt.Sprintf("action %s algorithm %s not supported", handler.ActionSettings.Id, types.HashAlgSHA256), err, 400, bodyContext)
				}

				if err != nil {
					wrenchContext.SetHasError3(span, err.Error(), err, 500, bodyContext)
				} else {
					bodyContext.SetBodyAction(handler.ActionSettings, sig)
				}
			}
		}
	}
	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handler *FuncSignatureHandler) SetNext(next Handler) {
	handler.Next = next
}

func (handler *FuncSignatureHandler) signBodyRSA_SHA256(priv *rsa.PrivateKey, body []byte) ([]byte, error) {
	// 1) Hash the body
	hash := sha256.Sum256(body)

	// 2) Sign the hash with RSA PKCS#1 v1.5 using SHA-256

	// NOTE: We must use RSASSA-PKCS1-v1_5 (rsa.SignPKCS1v15) here because the JWT
	// algorithm is RS256, which is defined as PKCS#1 v1.5 with SHA-256.
	// This is used for signatures only (not encryption), and is required
	// for compatibility with jwt.io and other RS256 consumers.
	// NOSONAR: S5542
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}
	// 3) Encode as base64 (good for HTTP headers)
	return sig, nil
}
