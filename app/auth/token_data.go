package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"wrench/app/json_map"
)

type TokenData struct {
	AccessToken      string  `json:"access_token"`
	ExpiresIn        float64 `json:"expires_in"`
	RefreshExpiresIn int     `json:"refresh_expires_in"`
	TokenType        string  `json:"token_type"`
	Scope            string  `json:"scope"`

	jwtPaylodData map[string]interface{}
	CustomToken   map[string]interface{}

	ForceReloadSeconds int64
	IsNotJwt           bool
	HeaderName         string
}

func (token *TokenData) LoadJwtPayload() {
	if len(token.AccessToken) > 0 && !token.IsNotJwt {
		jwtArray := strings.Split(token.AccessToken, ".")
		payloadBase64 := jwtArray[1]
		token.jwtPaylodData = ConvertJwtPayloadBase64ToJwtPaylodData(payloadBase64)
	}
}

func (token *TokenData) IsExpired(lessTimeMinutes float64, isOpaque bool) bool {
	var exp float64
	var ok bool

	if isOpaque || token.IsNotJwt {
		exp = token.ExpiresIn
	} else {
		exp, ok = token.jwtPaylodData["exp"].(float64)
		if !ok {
			return true
		}
	}

	lessTimes := -time.Duration(lessTimeMinutes) * time.Minute
	expireIn := time.Unix(int64(exp), 0).Add(lessTimes).Unix()

	currentTime := time.Now().Unix()

	if expireIn < currentTime {
		return true
	} else {
		return false
	}
}

func (token *TokenData) LoadCustomToken(forceReloadSeconds int64, accessTokenPropertyName string, tokenType string, headerName string) {
	token.IsNotJwt = true
	token.ForceReloadSeconds = forceReloadSeconds
	var now = time.Now().UTC().Add(time.Second * time.Duration(token.ForceReloadSeconds))
	token.ExpiresIn = float64(now.Unix())
	accessToken, _ := json_map.GetValue(token.CustomToken, accessTokenPropertyName, false)
	token.AccessToken = accessToken
	token.TokenType = tokenType
	token.HeaderName = headerName
}

func ConvertJwtPayloadBase64ToJwtPaylodData(jwtPayload string) map[string]interface{} {
	jwtPayload = strings.ReplaceAll(jwtPayload, "-", "+")
	jwtPayload = strings.ReplaceAll(jwtPayload, "_", "/")
	switch len(jwtPayload) % 4 {
	case 2:
		jwtPayload += "=="
	case 3:
		jwtPayload += "="
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(jwtPayload)
	if err != nil {
		fmt.Println("Error decoding Base64:", err)
		return nil
	}

	var jwtPaylodData map[string]interface{}
	jsonErr := json.Unmarshal(decodedBytes, &jwtPaylodData)
	if jsonErr != nil {
		return nil
	}
	return jwtPaylodData
}
