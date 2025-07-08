package api_settings

import "wrench/app/manifest/validation"

type CorsSettings struct {
	Origins []string `yaml:"origins"`
	Methods []string `yaml:"methods"`
	Headers []string `yaml:"headers"`
}

func (setting *CorsSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	return result
}

func (settings *CorsSettings) Merge(toMerge *CorsSettings) error {

	if toMerge == nil {
		return nil
	}

	if len(toMerge.Origins) > 0 {
		if len(settings.Origins) == 0 {
			settings.Origins = toMerge.Origins
		} else {
			settings.Origins = append(settings.Origins, toMerge.Origins...)
		}
	}

	if len(toMerge.Methods) > 0 {
		if len(settings.Methods) == 0 {
			settings.Methods = toMerge.Methods
		} else {
			settings.Methods = append(settings.Methods, toMerge.Methods...)
		}
	}

	if len(toMerge.Headers) > 0 {
		if len(settings.Headers) == 0 {
			settings.Headers = toMerge.Headers
		} else {
			settings.Headers = append(settings.Headers, toMerge.Headers...)
		}
	}

	return nil
}
