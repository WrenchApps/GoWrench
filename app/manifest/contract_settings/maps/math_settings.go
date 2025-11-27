package maps

import (
	"fmt"
	"regexp"
	"strings"
	"wrench/app/manifest/validation"
)

var mathExprRe = regexp.MustCompile(`^[a-zA-Z0-9_.]+(?:[+\-/][0-9]+|\*-?[0-9]+)$`)

type MathSettings []string

func (setting MathSettings) Valid() validation.ValidateResult {
    var result validation.ValidateResult

    if len(setting) == 0 {
        result.AddError("contract.maps.math must configure at least one expression")
        return result
    }

    for _, expr := range setting {
        if strings.Contains(expr, " ") {
            result.AddError("contract.maps.math expressions can't contain spaces")
            continue
        }
        if !mathExprRe.MatchString(expr) {
            result.AddError(fmt.Sprintf("contract.maps.math expression invalid: '%s'", expr))
        }
    }

    return result
}