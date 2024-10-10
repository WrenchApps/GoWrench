package manifest

import (
	"io"
	"os"

	"wrench/app/manifest/application_settings"

	"gopkg.in/yaml.v3"
)

func LoadYamlFile(pathFile string) (*application_settings.ApplicationSetting, error) {
	file, err := os.Open(pathFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	applicationSetting := new(application_settings.ApplicationSetting)
	err = yaml.Unmarshal(data, applicationSetting)
	if err != nil {
		return nil, err
	}

	return applicationSetting, nil
}
