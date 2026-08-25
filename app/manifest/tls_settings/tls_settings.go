package tls_settings

import (
	"fmt"
	"wrench/app/manifest/validation"
)

type CertificateFileType string

const (
	CertificateFileTypePem CertificateFileType = "pem"
	CertificateFileTypeCrt CertificateFileType = "crt"
)

type TlsSettings struct {
	Id                 string              `yaml:"id"`
	Type               CertificateFileType `yaml:"type"`
	CertificateBase64  string              `yaml:"certificateBase64"`
	IntermediateBase64 string              `yaml:"intermediateBase64"`
	CaBase64           string              `yaml:"caBase64"`
	ServerName         string              `yaml:"serverName"`
	KeyId              string              `yaml:"keyId"`
}

func (setting *TlsSettings) GetId() string {
	return setting.Id
}

func (settings *TlsSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(settings.Id) == 0 {
		result.AddError("tls.id is required")
	}

	if settings.Type != CertificateFileTypePem &&
		settings.Type != CertificateFileTypeCrt {
		result.AddError(fmt.Sprintf("tls[%s].type invalid type %s. Should be a valid value -> 'pem', 'crt' ", settings.Id, settings.Type))
	}

	if len(settings.CertificateBase64) == 0 {
		result.AddError(fmt.Sprintf("tls[%s].certificateBase64 is required", settings.Id))
	}

	if len(settings.KeyId) == 0 {
		result.AddError(fmt.Sprintf("tls[%s].keyId is required", settings.Id))
	}

	return result
}
