package handlers

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
	"wrench/app/manifest/action_settings/func_settings"
)

type FuncHashHandler struct {
	ActionSettings *settings.ActionSettings
	Next           Handler
}

func (handler *FuncHashHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {
	if !wrenchContext.HasError {
		key := handler.ActionSettings.Hash.Key

		hashType := handler.getHashFunc(handler.ActionSettings.Hash.Alg)

		mac := hmac.New(hashType, []byte(key))
		currentBody := bodyContext.BodyByteArray
		mac.Write(currentBody)

		expectedMAC := mac.Sum(nil)
		hashEncode := hex.EncodeToString(expectedMAC)

		if handler.ActionSettings.PreserveCurrentBody {
			bodyContext.SetBodyPreserved(handler.ActionSettings.Id, []byte(hashEncode))
		} else {
			bodyContext.BodyByteArray = []byte(hashEncode)
		}
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}

}

func (handler *FuncHashHandler) SetNext(next Handler) {
	handler.Next = next
}

func (handler *FuncHashHandler) getHashFunc(alg func_settings.FuncHashAlg) func() hash.Hash {
	switch alg {
	case func_settings.FuncHashAlgSHA1:
		return sha1.New
	case func_settings.FuncHashAlgSHA256:
		return sha256.New
	case func_settings.FuncHashAlgSHA512:
		return sha512.New
	case func_settings.FuncHashAlgMD5:
		return md5.New
	default:
		fmt.Println("Unsupported hash type, using SHA-256 as default.")
		return nil
	}
}
