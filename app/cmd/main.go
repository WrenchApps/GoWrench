package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"wrench/app"
	"wrench/app/manifest/application_settings"
	"wrench/app/startup"
	"wrench/app/startup/token_credentials"

	"gopkg.in/yaml.v3"
)

func main() {
	startup.LoadEnvsFiles()

	byteArray, err := startup.LoadYamlFile(getFileConfigPath())
	if err != nil {
		log.Fatalf("Error loading YAML: %v", err)
	}
	byteArray = startup.EnvInterpolation(byteArray)

	applicationSetting, err := parseToApplicationSetting(byteArray)
	if err != nil {
		log.Fatalf("Error parse yaml: %v", err)
	}

	application_settings.ApplicationSettingsStatic = applicationSetting
	var result = applicationSetting.Valid()

	if result.HasError() == true {
		var errors = result.GetError()
		for _, error := range errors {
			fmt.Println(error)
		}
	} else {
		go token_credentials.LoadTokenCredentialAuthentication()
		var router = startup.LoadApplicationSettings(applicationSetting)
		port := getPort()
		http.ListenAndServe(port, router)
	}
}

func parseToApplicationSetting(data []byte) (*application_settings.ApplicationSettings, error) {

	applicationSettings := new(application_settings.ApplicationSettings)
	err := yaml.Unmarshal(data, applicationSettings)
	if err != nil {
		return nil, err
	}
	return applicationSettings, nil
}

func getFileConfigPath() string {
	configPath := os.Getenv(app.ENV_PATH_FILE_CONFIG)
	if len(configPath) == 0 {
		configPath = "../../configApp.yaml"
	}
	return configPath
}

func getPort() string {
	port := os.Getenv(app.ENV_PORT)
	if len(port) == 0 {
		port = ":9090"
	}

	if port[0] != ':' {
		port = ":" + port
	}

	return port
}
