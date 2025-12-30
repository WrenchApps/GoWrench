package tls_load

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"wrench/app/manifest/application_settings"
	keys_load "wrench/app/startup/keys"
)

var ErrorLoadTls []error
var tlsConfig map[string]*tls.Config

func LoadTls(ctx context.Context) {
	tlsConfig = make(map[string]*tls.Config)
	settings := application_settings.ApplicationSettingsStatic

	if len(settings.Tls) == 0 {
		return
	}

	for _, setting := range settings.Tls {

		var certDERs [][]byte
		certPEM, _ := base64.StdEncoding.DecodeString(setting.CertificateBase64)
		for {
			b, r := pem.Decode(certPEM)
			if b == nil {
				break
			}
			certPEM = r
			if b.Type == "CERTIFICATE" {
				certDERs = append(certDERs, b.Bytes)
			}
		}

		if len(setting.IntermediateBase64) > 0 {
			interPEM, _ := base64.StdEncoding.DecodeString(setting.IntermediateBase64)
			for {
				b, r := pem.Decode(interPEM)
				if b == nil {
					break
				}
				interPEM = r
				if b.Type == "CERTIFICATE" {
					certDERs = append(certDERs, b.Bytes)
				}
			}
		}

		if len(certDERs) == 0 {
			ErrorLoadTls = append(ErrorLoadTls, fmt.Errorf("None certificates found"))
		} else {

			key, errKey := keys_load.GetPrivateKey(setting.KeyId)

			if errKey != nil {
				ErrorLoadTls = append(ErrorLoadTls, errKey)
			}

			tlsCert := tls.Certificate{
				Certificate: certDERs,
				PrivateKey:  key,
			}

			var roots *x509.CertPool = nil
			if len(setting.CaBase64) > 0 {
				caCertPEM, _ := base64.StdEncoding.DecodeString(setting.CaBase64)

				roots, err := x509.SystemCertPool()
				if err != nil {
					roots = x509.NewCertPool()
				}

				if ok := roots.AppendCertsFromPEM(caCertPEM); !ok {
					ErrorLoadTls = append(ErrorLoadTls, fmt.Errorf("failed to append CA PEM"))
				}
			}

			tlsCfg := &tls.Config{
				MinVersion:   tls.VersionTLS12,
				Certificates: []tls.Certificate{tlsCert},
				RootCAs:      roots,
				ServerName:   setting.ServerName,
			}

			leaf, err := x509.ParseCertificate(certDERs[0])
			if err != nil {
				ErrorLoadTls = append(ErrorLoadTls, err)
			}

			rsaPub, ok := leaf.PublicKey.(*rsa.PublicKey)
			if !ok {
				ErrorLoadTls = append(ErrorLoadTls, fmt.Errorf("cert public key is not RSA: %T", leaf.PublicKey))
			}

			if rsaPub.N.Cmp(key.N) != 0 {
				ErrorLoadTls = append(ErrorLoadTls, fmt.Errorf("CERT/KEY MISMATCH: cert public modulus != private key modulus"))
			}

			tlsConfig[setting.Id] = tlsCfg
		}
	}
}

func GetTlsConfigById(tlsId string) (*tls.Config, error) {
	tlsConfig, ok := tlsConfig[tlsId]
	if !ok {
		return nil, fmt.Errorf("TlsConfig not found: %s", tlsId)
	}
	return tlsConfig, nil
}
