package token_credentials

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type JwtData struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
	Scope            string `json:"scope"`

	jwtPaylodData map[string]interface{}
}

func (jwt *JwtData) LoadJwtPayload() {
	if len(jwt.AccessToken) > 0 {
		jwtArray := strings.Split(jwt.AccessToken, ".")
		payloadBase64 := jwtArray[1]
		jwt.jwtPaylodData = convertJwtPayloadBase64ToJwtPaylodData(payloadBase64)
	}
}

func (jwt *JwtData) IsExpired(lessTimeMinutes float64) bool {

	exp, ok := jwt.jwtPaylodData["exp"].(float64)
	if !ok {
		return true
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

func convertJwtPayloadBase64ToJwtPaylodData(jwtPayload string) map[string]interface{} {
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
