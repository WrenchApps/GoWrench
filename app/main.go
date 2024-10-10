package main

import (
	"log"
	"wrench/app/manifest"
)

func main() {
	applicationSetting, err := manifest.LoadYamlFile("../configApp.yaml")

	if err != nil {
		log.Fatalf("Error loading YAML: %v", err)
	}

	var result = applicationSetting.Valid()
	log.Printf("Service: %v", result.HasError())

}
