package maps

import (
	"wrench/app/manifest/validation"
)

type Operator string

const (
	OperatorToArray  Operator = "to_array"
	OperatorGetFirst Operator = "get_first"
)

type ParseSettings struct {
	WhenEquals []string   `yaml:"whenEquals"`
	Operator   []Operator `yaml:"operator"`
}

func (setting ParseSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	return result
}
