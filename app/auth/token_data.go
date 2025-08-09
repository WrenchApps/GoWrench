package auth

import (
	"fmt"
	"strings"
	"time"
	auth_jwt "wrench/app/auth/jwt"
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
		token.jwtPaylodData = auth_jwt.ConvertJwtPayloadBase64ToJwtPaylodData(payloadBase64)
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
	token.AccessToken = fmt.Sprint(accessToken)
	token.TokenType = tokenType
	token.HeaderName = headerName
}
