package maps

import (
	"fmt"
	"slices"
	"strings"
	"wrench/app/manifest/validation"
)

var funcValids = []string{"rename", "new", "remove", "duplicate", "parse", "format"}

type ContractMapSetting struct {
	Id        string          `yaml:"id"`
	Rename    []string        `yaml:"rename"`
	Remove    []string        `yaml:"remove"`
	Sequence  []string        `yaml:"sequence"`
	New       []string        `yaml:"new"`
	Duplicate []string        `yaml:"duplicate"`
	Parse     *ParseSettings  `yaml:"parse"`
	Format    *FormatSettings `yaml:"format"`
}

func (setting ContractMapSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult
	totalMapConfigured := 0

	if len(setting.Id) <= 0 {
		result.AddError("contract.maps.id is required")
	}

	if len(setting.Rename) > 0 {
		totalMapConfigured++
		errorSplitted := "contract.maps.rename should be configured looks like 'propertySource:propertyDestination' without space"
		for _, property := range setting.Rename {

			if strings.Contains(property, " ") {
				result.AddError(errorSplitted)
			}

			propertySplitted := strings.Split(property, ":")
			if len(propertySplitted) != 2 {
				result.AddError(errorSplitted)
			}
		}
	}

	if len(setting.Duplicate) > 0 {
		totalMapConfigured++
		errorSplitted := "contract.maps.duplicate should be configured looks like 'propertySource:propertyDestination' without space"
		for _, property := range setting.Duplicate {

			if strings.Contains(property, " ") {
				result.AddError(errorSplitted)
			}

			propertySplitted := strings.Split(property, ":")
			if len(propertySplitted) != 2 {
				result.AddError(errorSplitted)
			}
		}
	}

	if len(setting.Remove) > 0 {
		totalMapConfigured++
		for _, remove := range setting.Remove {

			if strings.Contains(remove, " ") {
				result.AddError("contract.maps.remove can't contain space")
			}
		}
	}

	if len(setting.New) > 0 {
		totalMapConfigured++
	}

	if setting.Parse != nil {
		totalMapConfigured++
		result.AppendValidable(setting.Parse)
	}

	if setting.Format != nil {
		totalMapConfigured++
		result.AppendValidable(setting.Format)
	}

	if len(setting.Sequence) > 0 {

		if totalMapConfigured != len(setting.Sequence) {
			result.AddError("When sequence is configured should be informed all maps configured")
		}

		for _, s := range setting.Sequence {
			if slices.Contains(funcValids, s) == false {
				result.AddError(fmt.Sprintf("contract.maps.sequence should contain valid values. The value %s is not valid", s))
			}

			if s == "rename" && setting.Rename == nil {
				result.AddError("contract.maps.sequence rename not configured")
			} else if s == "new" && setting.New == nil {
				result.AddError("contract.maps.new rename not configured")
			} else if s == "remove" && setting.Remove == nil {
				result.AddError("contract.maps.sequence remove not configured")
			} else if s == "parse" && setting.Parse == nil {
				result.AddError("contract.maps.sequence parse not configured")
			} else if s == "format" && setting.Format == nil {
				result.AddError("contract.maps.sequence format not configured")
			}
		}
	}

	return result
}
