package keys_load

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"strings"
	"wrench/app"
	"wrench/app/manifest/application_settings"
)

var privateKeys map[string]*rsa.PrivateKey
var ErrorLoadKeys []error

func LoadKeys() {
	settings := application_settings.ApplicationSettingsStatic

	if settings.Keys == nil {
		return
	}

	for _, key := range settings.Keys {
		_, err := LoadEncryptedPrivateKey(key.Id, key.PrivateRsaKeyDERBase64)
		addIfErrorKey(err)
	}
}

func addIfErrorKey(err error) {
	if err != nil {
		app.LogError2("Error connections: %v", err)
		ErrorLoadKeys = append(ErrorLoadKeys, err)
	}
}

func LoadEncryptedPrivateKey(keyId, privateRsakeyDERBase64 string) (*rsa.PrivateKey, error) {

	if privateKeys == nil {
		privateKeys = make(map[string]*rsa.PrivateKey)
	}

	derBytes, err := base64.StdEncoding.DecodeString(privateRsakeyDERBase64)
	if err != nil {
		return nil, fmt.Errorf("read key: %w", err)
	}

	key, err := x509.ParsePKCS8PrivateKey(derBytes)
	if err != nil {

		if strings.Contains(err.Error(), "use ParsePKCS1PrivateKey") {
			key, err = x509.ParsePKCS1PrivateKey(derBytes)
		}

		if err != nil {
			return nil, fmt.Errorf("parse key: %w", err)
		}
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not RSA")
	}

	privateKeys[keyId] = rsaKey

	return rsaKey, nil
}

func GetPrivateKey(keyId string) (*rsa.PrivateKey, error) {
	key, ok := privateKeys[keyId]
	if !ok {
		return nil, fmt.Errorf("key not found: %s", keyId)
	}
	return key, nil
}
