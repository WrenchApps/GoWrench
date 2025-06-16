package contexts

import (
	"strconv"
	"time"
)

type GeneralFuncType string

const (
	FuncTypeTimestampMilli GeneralFuncType = "func.timestamp(milli)"
)

func GetFuncValue(funcType GeneralFuncType) string {
	switch funcType {
	case FuncTypeTimestampMilli:
		return getTimestamp()
	default:
		return ""
	}
}

func getTimestamp() string {
	return strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
}
