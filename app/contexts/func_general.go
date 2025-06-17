package contexts

import (
	"encoding/base64"
	"strconv"
	"time"
	settings "wrench/app/manifest/action_settings"
)

type FuncGeneralType string

const (
	FuncTypeTimestampMilli FuncGeneralType = "func.timestamp(milli)"
	FuncTypeBase64Encode   FuncGeneralType = "func.base64(encode)"
)

func GetFuncValue(funcType FuncGeneralType, wrenchContext *WrenchContext, bodyContext *BodyContext, action *settings.ActionSettings) string {
	switch funcType {
	case FuncTypeTimestampMilli:
		return getTimestamp()
	case FuncTypeBase64Encode:
		bodyArray := bodyContext.GetBody(action)
		return base64.StdEncoding.EncodeToString(bodyArray)
	default:
		return ""
	}
}

func getTimestamp() string {
	return strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
}
