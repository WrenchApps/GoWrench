package cross_validation

import (
	"wrench/app/manifest/validation"
)

func Valid() validation.ValidateResult {
	var result validation.ValidateResult

	result.Append(httpRequestCrossValid())

	return result
}
