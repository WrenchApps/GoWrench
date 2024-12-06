package main

import (
	"fmt"
	"log"
	"net/http"
	"wrench/app/manifest"
	"wrench/app/manifest/application_settings"
	startup "wrench/app/startup"
)

func main() {
	applicationSetting, err := manifest.LoadYamlFile("../../configApp.yaml")
	application_settings.ApplicationSettingsStatic = applicationSetting

	if err != nil {
		log.Fatalf("Error loading YAML: %v", err)
	}

	var result = applicationSetting.Valid()

	if result.HasError() == true {
		var errors = result.GetError()
		for _, error := range errors {
			fmt.Println(error)
		}
	} else {
		var router = startup.LoadApplicationSettings(applicationSetting)
		http.ListenAndServe(":8085", router)
	}
}
