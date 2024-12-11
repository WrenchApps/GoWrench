package do_settings

import (
	"wrench/app/manifest/validation"
)

type DoSetting struct {
	Before *BeforeSetting `yaml:"before"`
	After  *AfterSetting  `yaml:"after"`
}

func (setting DoSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if setting.After != nil {
		result.AppendValidable(setting.After)
	}

	if setting.Before != nil {
		result.AppendValidable(setting.Before)
	}

	return result
}
