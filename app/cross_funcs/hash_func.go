package cross_funcs

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"wrench/app/manifest/types"
)

func GetHash(key string, hashType func() hash.Hash, value []byte) string {

	mac := hmac.New(hashType, []byte(fmt.Sprint(key)))

	mac.Write(value)
	expectedMAC := mac.Sum(nil)

	return hex.EncodeToString(expectedMAC)
}

func GetHashFunc(alg types.HashAlg) func() hash.Hash {
	switch alg {
	case types.HashAlgSHA1:
		return sha1.New
	case types.HashAlgSHA256:
		return sha256.New
	case types.HashAlgSHA512:
		return sha512.New
	case types.HashAlgMD5:
		return md5.New
	default:
		fmt.Println("Unsupported hash type")
		return nil
	}
}
