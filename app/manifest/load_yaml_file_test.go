package manifest

import (
	"os"
	"testing"
)

func TestLoadYamlFile_Success(t *testing.T) {
	yamlContent := `
version: 1
service:
  name: test-app
  version: 1.0.0
api:
  endpoints:
    - route: "api/customer"
      method: get
      actionId: get_customer
actions:
  - id: get_customer
    type: httpRequestMock
    mock:
      value: "Hi"
      contentType: "application/json"
      method: get

`
	tmpFile, err := os.CreateTemp("", "configApp*.yaml")
	if err != nil {
		t.Fatalf("Could not create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	originalFilePath := "configApp.yaml"
	os.Rename(tmpFile.Name(), originalFilePath)
	defer os.Rename(originalFilePath, tmpFile.Name())

	config, err := LoadYamlFile("configApp.yaml")
	if err != nil {
		t.Fatalf("LoadYamlFile failed: %v", err)
	}

	if config.Service.Name != "test-app" {
		t.Errorf("Expected service name 'test-app', got '%s'", config.Service.Name)
	}

	if config.Api.Endpoints[0].Route != "api/customer" {
		t.Errorf("Expected API route 'api/customer', got '%s'", config.Api.Endpoints[0].Route)
	}

	if config.Actions[0].Mock.Value != "Hi" {
		t.Errorf("Expected mock value 'Hi', got '%s'", config.Actions[0].Mock.Value)
	}
}

func TestLoadYamlFile_FileNotFound(t *testing.T) {
	originalFilePath := "configApp.yaml"
	os.Rename(originalFilePath, originalFilePath+".bak")
	defer os.Rename(originalFilePath+".bak", originalFilePath)

	_, err := LoadYamlFile("configApp.yaml")
	if err == nil {
		t.Fatalf("Expected error when file is not found, but got nil")
	}
}

func TestLoadYamlFile_InvalidYaml(t *testing.T) {
	yamlContent := `
version: 1
service:
	name: test-app
	version: 1.0.0
api:
	endpoints:
		- route: "api/customer"
			method: get
			actionId: get_customer
actions:
	- id: get_customer
		type: httpRequestMock
		mock:
			value: "Hi"
			contentType: "application/json"
			method: get
`
	tmpFile, err := os.CreateTemp("", "configApp*.yaml")
	if err != nil {
		t.Fatalf("Could not create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	originalFilePath := "configApp.yaml"
	os.Rename(tmpFile.Name(), originalFilePath)
	defer os.Rename(originalFilePath, tmpFile.Name())

	_, err = LoadYamlFile("configApp.yaml")
	if err == nil {
		t.Fatalf("Expected error when file is not found, but got nil")
	}
}
