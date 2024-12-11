package token_credentials

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"
	client "wrench/app/clients/http"
	"wrench/app/manifest/application_settings"
	credential "wrench/app/manifest/token_credential"
)

var tokenCredentials map[string]*JwtData

func GetTokenCredentialById(tokenCredentialId string) *JwtData {
	if tokenCredentials == nil {
		return nil
	}

	return tokenCredentials[tokenCredentialId]
}

func LoadTokenCredentialAuthentication() {
	app_settings := application_settings.ApplicationSettingsStatic

	if len(app_settings.TokenCredentials) > 0 {
		if tokenCredentials == nil {
			tokenCredentials = make(map[string]*JwtData)
		}

		for {
			for _, setting := range app_settings.TokenCredentials {

				jwtData := GetTokenCredentialById(setting.Id)
				if jwtData != nil {
					if jwtData.IsExpired(5) == false {
						continue
					}
				}

				jwtData, err := authenticateClientCredentials(setting)
				if err != nil {
					continue
					// TODO setting error to unhealthy api
				}
				jwtData.LoadJwtPayload()
				tokenCredentials[setting.Id] = jwtData
			}

			time.Sleep(2 * time.Minute)
		}
	}
}

func authenticateClientCredentials(setting *credential.TokenCredentialSetting) (*JwtData, error) {
	request := new(client.HttpClientRequestData)
	data := url.Values{}
	data.Set("client_id", setting.ClientId)
	data.Set("client_secret", setting.ClientSecret)
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
			token := string(response.Body)
			tokenArray := []byte(token)
			jwtData := new(JwtData)

			jsonErr := json.Unmarshal(tokenArray, &jwtData)
			if jsonErr != nil {
				return nil, jsonErr
			}

			return jwtData, nil
		}

		return nil, errors.New(fmt.Sprintf("Can't get jwtToken response_status_code: %v", response.StatusCode))
	}
}
