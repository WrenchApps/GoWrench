package token_credentials

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
	"wrench/app/auth"
	client "wrench/app/clients/http"
	"wrench/app/json_map"
	"wrench/app/manifest/application_settings"
	credential "wrench/app/manifest/token_credential_settings"
)

var tokenCredentials map[string]*auth.TokenData

func GetTokenCredentialById(tokenCredentialId string) *auth.TokenData {
	if tokenCredentials == nil {
		return nil
	}

	return tokenCredentials[tokenCredentialId]
}

func LoadTokenCredentialAuthentication() {
	app_settings := application_settings.ApplicationSettingsStatic

	if len(app_settings.TokenCredentials) > 0 {
		if tokenCredentials == nil {
			tokenCredentials = make(map[string]*auth.TokenData)
		}

		for {
			for _, setting := range app_settings.TokenCredentials {

				tokenData := GetTokenCredentialById(setting.Id)
				if tokenData != nil {
					if !tokenData.IsExpired(5, setting.IsOpaque) {
						continue
					}
				}

				var err error
				if setting.Type == credential.TokenCredentialClientCredential {
					tokenData, err = authenticateClientCredentials(setting)
				} else if setting.Type == credential.TokenCredentialBasicCredential {
					tokenData, err = basicCredentials(setting)
				} else if setting.Type == credential.TokenCredentialCustomAuthentication {
					tokenData, err = customAuthentication(setting)
				}

				if err != nil {
					continue
					// TODO setting error to unhealthy api
				}

				if !setting.IsOpaque {
					tokenData.LoadJwtPayload()
				}

				tokenCredentials[setting.Id] = tokenData
			}

			time.Sleep(2 * time.Minute)
		}
	}
}

func authenticateClientCredentials(setting *credential.TokenCredentialSetting) (*auth.TokenData, error) {
	request := new(client.HttpClientRequestData)
	data := url.Values{}
	data.Set("client_id", setting.ClientCredential.ClientId)
	data.Set("client_secret", setting.ClientCredential.ClientSecret)
	data.Set("grant_type", "client_credentials")

	request.Body = []byte(data.Encode())
	request.Method = "POST"
	request.Url = setting.AuthEndpoint

	request.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	ctx := context.Background()

	response, err := client.HttpClientDo(ctx, request)

	if err != nil {
		// TODO add error flow
		return nil, err
	} else {

		if response.StatusCodeSuccess() {
			tokenData := new(auth.TokenData)
			jsonErr := json.Unmarshal(response.Body, &tokenData)
			if jsonErr != nil {
				return nil, jsonErr
			}

			return tokenData, nil
		}

		return nil, fmt.Errorf("Can't get jwtToken response_status_code: %v", response.StatusCode)
	}
}

func basicCredentials(setting *credential.TokenCredentialSetting) (*auth.TokenData, error) {
	request := new(client.HttpClientRequestData)
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	request.Body = []byte(data.Encode())
	request.Method = "POST"
	request.Url = setting.AuthEndpoint

	credential := fmt.Sprintf("%s:%s", setting.Basic.Username, setting.Basic.Password)
	credentialEncoded := base64.StdEncoding.EncodeToString([]byte(credential))

	request.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	request.SetHeader("Authorization", fmt.Sprintf("Basic %s", credentialEncoded))

	ctx := context.Background()

	response, err := client.HttpClientDo(ctx, request)

	if err != nil {
		// TODO add error flow
		return nil, err
	} else {

		if response.StatusCodeSuccess() {
			tokenData := new(auth.TokenData)
			jsonErr := json.Unmarshal(response.Body, &tokenData)
			if jsonErr != nil {
				return nil, jsonErr
			}

			return tokenData, nil
		}

		return nil, fmt.Errorf("Can't get jwtToken response_status_code: %v", response.StatusCode)
	}
}

func customAuthentication(setting *credential.TokenCredentialSetting) (*auth.TokenData, error) {
	request := new(client.HttpClientRequestData)
	isFormUrlencoded := false
	if len(setting.Custom.RequestHeaders) > 0 {
		for key, value := range setting.Custom.RequestHeaders {
			request.SetHeader(key, value)
		}
		value := setting.Custom.RequestHeaders["Content-Type"]
		isFormUrlencoded = value == "application/x-www-form-urlencoded"
	}

	if isFormUrlencoded {
		data := url.Values{}

		for prop, value := range setting.Custom.RequestBody {
			data.Set(prop, value)
		}

		request.Body = []byte(data.Encode())
	} else {
		if len(setting.Custom.RequestBody) > 0 {
			request.Body = getBodyCustomAuthentication(setting.Custom.RequestBody)
		}
	}

	request.Method = string(setting.Custom.Method)
	request.Url = setting.AuthEndpoint

	ctx := context.Background()
	response, err := client.HttpClientDo(ctx, request)

	if err != nil {
		// TODO add error flow
		return nil, err
	} else {

		if response.StatusCodeSuccess() {
			tokenData := new(auth.TokenData)

			jsonErr := json.Unmarshal(response.Body, &tokenData.CustomToken)
			if jsonErr != nil {
				return nil, jsonErr
			}

			tokenData.LoadCustomToken(setting.GetForceReloadTimeSecondsValue(),
				setting.Custom.Configs.AccessTokenPropertyName,
				setting.Custom.Configs.TokenType,
				setting.Custom.Configs.HeaderName)

			return tokenData, nil
		}

		return nil, fmt.Errorf("Can't get jwtToken response_status_code: %v", response.StatusCode)
	}
}

func getBodyCustomAuthentication(customBody map[string]string) []byte {
	jsonMap := make(map[string]interface{})
	for key, value := range customBody {
		json_map.CreateProperty(jsonMap, key, value)
	}

	jsonArray, _ := json.Marshal(jsonMap)
	return jsonArray
}
