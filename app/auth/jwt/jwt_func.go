package auth_jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

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
