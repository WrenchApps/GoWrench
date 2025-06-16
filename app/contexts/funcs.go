package contexts

import (
	"strconv"
	"time"
)

type FuncType string

const (
	FuncTypeTimestampMilli FuncType = "func.timestamp(milli)"
)

func GetFuncValue(funcType FuncType) string {
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
