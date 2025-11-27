package cross_funcs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

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