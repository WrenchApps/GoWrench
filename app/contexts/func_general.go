package contexts

import (
	"encoding/base64"
	"strconv"
	"time"
	settings "wrench/app/manifest/action_settings"
	"wrench/app/manifest/action_settings/func_settings"
)

func GetFuncValue(funcType func_settings.FuncGeneralType, wrenchContext *WrenchContext, bodyContext *BodyContext, action *settings.ActionSettings) string {
	switch funcType {
	case func_settings.FuncTypeTimestampMilli:
		return getTimestamp()
	case func_settings.FuncTypeBase64Encode:
		bodyArray := bodyContext.GetBody(action)
		return base64.StdEncoding.EncodeToString(bodyArray)
	default:
		return ""
	}
}

func getTimestamp() string {
	return strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
}
