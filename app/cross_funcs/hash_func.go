package cross_funcs

import (
	"crypto/hmac"
	"encoding/hex"
	"fmt"
	"hash"
)

func GetHash(key string, hashType func() hash.Hash, value []byte) string {

	mac := hmac.New(hashType, []byte(fmt.Sprint(key)))

	mac.Write(value)
	expectedMAC := mac.Sum(nil)

	return hex.EncodeToString(expectedMAC)
}
