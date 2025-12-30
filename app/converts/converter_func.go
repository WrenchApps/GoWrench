package converts

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const StringToBool = "{{convertStringToBool}}"
const StringToFloat = "{{convertStringToFloat}}"

func ConvertToTypes(value any, targetType string) (any, error) {
	switch targetType {
	case StringToFloat:
		return ConvertToFloat(value)
	case StringToBool:
		return ConvertStringToBool(value)
	default:
		return nil, fmt.Errorf("unsupported target type: %s", targetType)
	}
}

func ConvertToFloat(value any) (float64, error) {
	switch x := value.(type) {
	case float64:
		return x, nil
	case float32:
		return float64(x), nil
	case int:
		return float64(x), nil
	case int64:
		return float64(x), nil
	case json.Number:
		return x.Float64()
	case string:
		return strconv.ParseFloat(strings.TrimSpace(x), 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

func ConvertStringToBool(value any) (bool, error) {
	switch x := value.(type) {
	case bool:
		return x, nil
	case string:
		trimmed := strings.TrimSpace(strings.ToLower(x))
		if trimmed == "true" || trimmed == "\"true\"" || trimmed == "1" || trimmed == "yes" {
			return true, nil
		} else if trimmed == "false" || trimmed == "\"false\"" || trimmed == "0" || trimmed == "no" {
			return false, nil
		}
		return false, fmt.Errorf("cannot convert string '%s' to bool", x)
	default:
		return false, fmt.Errorf("cannot convert %T to bool", value)
	}
}
