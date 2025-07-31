package application_settings

import (
	"errors"
	"log"
	"wrench/app/manifest/action_settings"
	"wrench/app/manifest/api_settings"

	"wrench/app/manifest/connection_settings"
	"wrench/app/manifest/contract_settings"
	"wrench/app/manifest/idemp_settings"
	"wrench/app/manifest/service_settings"
	credential "wrench/app/manifest/token_credential_settings"
	"wrench/app/manifest/validation"

	"gopkg.in/yaml.v3"
)

var ApplicationSettingsStatic *ApplicationSettings

type ApplicationSettings struct {
	Connections      *connection_settings.ConnectionSettings `yaml:"connections"`
	Api              *api_settings.ApiSettings               `yaml:"api"`
	Service          *service_settings.ServiceSettings       `yaml:"service"`
	Contract         *contract_settings.ContractSetting      `yaml:"contract"`
	Actions          []*action_settings.ActionSettings       `yaml:"actions"`
	TokenCredentials []*credential.TokenCredentialSetting    `yaml:"tokenCredentials"`
	Idemps           []*idemp_settings.IdempSettings         `yaml:"idemps"`
}

func (settings *ApplicationSettings) GetActionById(actionId string) (*action_settings.ActionSettings, error) {
	for _, action := range settings.Actions {
		if action.Id == actionId {
			return action, nil
		}
	}

	return nil, errors.New("action not found")
}

func (settings *ApplicationSettings) GetEndpointByActionId(actionId string) (*api_settings.EndpointSettings, error) {
	for _, endpoint := range settings.Api.Endpoints {
		if endpoint.ActionID == actionId {
			return &endpoint, nil
		}
	}

	return nil, errors.New("endpoint not found")
}

func (settings *ApplicationSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if settings.Connections != nil {
		result.AppendValidable(settings.Connections)
	}

	if settings.Api != nil {
		result.AppendValidable(settings.Api)
	}

	if settings.Actions != nil {
		for _, validable := range settings.Actions {
			result.AppendValidable(validable)
			result.Append(actionValidation(validable))
		}
	}

	if settings.Api != nil {
		result.AppendValidable(settings.Api)
		result.Append(apiEndpointsValidation())
	}

	if settings.TokenCredentials != nil {
		for _, validable := range settings.TokenCredentials {
			result.AppendValidable(validable)
		}
	}

	if settings.Contract != nil {
		result.AppendValidable(settings.Contract)
	}

	if settings.Idemps != nil {
		for _, validable := range settings.Idemps {
			result.AppendValidable(validable)
		}
	}

	return result
}

func (settings *ApplicationSettings) Merge(toMerge *ApplicationSettings) error {

	if settings.Service != nil && toMerge.Service != nil {
		return errors.New("should be informed only once service")
	} else {
		if settings.Service == nil {
			settings.Service = toMerge.Service
		}
	}

	if settings.Connections == nil && toMerge.Connections != nil {
		settings.Connections = &connection_settings.ConnectionSettings{}
	}
	if settings.Connections != nil && toMerge.Connections != nil {
		if err := settings.Connections.Merge(toMerge.Connections); err != nil {
			return err
		}
	}

	if settings.Api == nil && toMerge.Api != nil {
		settings.Api = &api_settings.ApiSettings{}
	}
	if settings.Api != nil && toMerge.Api != nil {
		if err := settings.Api.Merge(toMerge.Api); err != nil {
			return err
		}
	}

	if settings.Contract == nil && toMerge.Contract != nil {
		settings.Contract = &contract_settings.ContractSetting{}
	}
	if settings.Contract != nil && toMerge.Contract != nil {
		if err := settings.Contract.Merge(toMerge.Contract); err != nil {
			return err
		}
	}

	if len(toMerge.Actions) > 0 {
		if len(settings.Actions) == 0 {
			settings.Actions = toMerge.Actions
		} else {
			settings.Actions = append(settings.Actions, toMerge.Actions...)
		}
	}

	if len(toMerge.TokenCredentials) > 0 {
		if len(settings.TokenCredentials) == 0 {
			settings.TokenCredentials = toMerge.TokenCredentials
		} else {
			settings.TokenCredentials = append(settings.TokenCredentials, toMerge.TokenCredentials...)
		}
	}

	if len(toMerge.Idemps) > 0 {
		if len(settings.Idemps) == 0 {
			settings.Idemps = toMerge.Idemps
		} else {
			settings.Idemps = append(settings.Idemps, toMerge.Idemps...)
		}
	}

	return nil
}

func ParseMapToApplicationSetting(datas map[string][]byte) (*ApplicationSettings, error) {

	applicationSettings := new(ApplicationSettings)

	for key, data := range datas {
		toMerge, err := ParseToApplicationSetting(data)

		if err != nil {
			return nil, err
		}

		if err2 := applicationSettings.Merge(toMerge); err2 != nil {
			return nil, err2
		}
		log.Printf("Done config file %s", key)
	}

	return applicationSettings, nil
}

func ParseToApplicationSetting(data []byte) (*ApplicationSettings, error) {

	applicationSettings := new(ApplicationSettings)

	err := yaml.Unmarshal(data, applicationSettings)
	if err != nil {
		return nil, err
	}
	return applicationSettings, nil

}
