package token_credential_settings

import (
	"strconv"
	"strings"
	"wrench/app/manifest/validation"
)

type TokenCredentialSetting struct {
	Id               string                   `yaml:"id"`
	Type             TokenCredentialType      `yaml:"type"`
	AuthEndpoint     string                   `yaml:"authEndpoint"`
	IsOpaque         bool                     `yaml:"isOpaque"`
	ClientCredential *ClientCredentialSetting `yaml:"clientCredential"`
	Basic            *BasicSetting            `yaml:"basic"`
	ForceReload      string                   `yaml:"forceReload"`
	Custom           *CustomAuthentication    `yaml:"custom"`
}

func (setting *TokenCredentialSetting) GetId() string {
	return setting.Id
}

type TokenCredentialType string

const (
	TokenCredentialClientCredential     TokenCredentialType = "client_credentials"
	TokenCredentialBasicCredential      TokenCredentialType = "basic"
	TokenCredentialCustomAuthentication TokenCredentialType = "custom_authentication"
)

var forceReloadTimeSelector string = ""
var forceReloadTimeValue int64 = 0

func (setting TokenCredentialSetting) GetForceReloadTimeSelector() string {
	if forceReloadTimeSelector != "" {
		return forceReloadTimeSelector
	}
	sizeForceReload := len(setting.ForceReload)
	forceReloadTimeSelector = string(setting.ForceReload[sizeForceReload-1])
	return forceReloadTimeSelector
}

func (setting TokenCredentialSetting) GetForceReloadTimeValue() int64 {
	if forceReloadTimeValue > 0 {
		return forceReloadTimeValue
	}
	timeSelector := setting.GetForceReloadTimeSelector()
	timeValueString := strings.ReplaceAll(setting.ForceReload, timeSelector, "")
	timeValueInt, err := strconv.ParseInt(timeValueString, 10, 0)
	if err == nil {
		forceReloadTimeValue = timeValueInt
	} else {
		forceReloadTimeValue = -1
	}
	return forceReloadTimeValue
}

func (setting TokenCredentialSetting) GetForceReloadTimeSecondsValue() int64 {
	if setting.GetForceReloadTimeSelector() == "s" {
		return setting.GetForceReloadTimeValue()
	} else if setting.GetForceReloadTimeSelector() == "m" {
		return setting.GetForceReloadTimeValue() * 60
	} else {
		return 3600
	}
}

func (setting TokenCredentialSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Id) == 0 {
		result.AddError("tokenCredentials.id is required")
	} else {
		if strings.ContainsAny(setting.Id, " ") {
			result.AddError("tokenCredentials.id can't contains space")
		}
	}

	if len(setting.AuthEndpoint) == 0 {
		result.AddError("tokenCredentials.authEndpoint is required")
	}

	if setting.Type == TokenCredentialClientCredential {

		if setting.ClientCredential == nil {
			result.AddError("tokenCredentials.ClientCredential is required")
		} else {
			result.AppendValidable(setting.ClientCredential)
		}
	}

	if setting.Type == TokenCredentialBasicCredential {

		if setting.Basic == nil {
			result.AddError("tokenCredentials.Basic is required")
		} else {
			result.AppendValidable(setting.Basic)
		}
	}

	if setting.Type == TokenCredentialCustomAuthentication {
		if len(setting.ForceReload) == 0 {
			result.AddError("tokenCredentials.ForceReload is required when type is custom_authentication")
		} else {
			timeSelector := setting.GetForceReloadTimeSelector()
			if timeSelector != "s" && timeSelector != "m" {
				result.AddError("tokenCredentials.ForceReload should use 's' to seconds or 'm' to minutes. Ex: 60s or 1m")
			}

			timeValueInt := setting.GetForceReloadTimeSecondsValue()
			if timeValueInt < 0 {
				result.AddError("tokenCredentials.ForceReload should inform the int value to refresh token after the selector time. Ex: 60x or 1m")
			} else {
				if timeValueInt < 600 {
					result.AddError("tokenCredentials.ForceReload should be greater than 600s or 10m")
				}
			}
			result.AppendValidable(setting.Custom)
		}
	}

	return result
}
