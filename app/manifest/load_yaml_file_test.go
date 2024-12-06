package manifest

// import (
// 	"os"
// 	"testing"
// )

// func TestLoadYamlFile_Success(t *testing.T) {
// 	yamlContent := `
// version: 1

// service:
//   name: test-app
//   version: 1.0.0

// api:
//   endpoints:
//     - route: /api/mock
//       method: get
//       actionId: get_customer_mock

//     - route: /api/customer
//       method: get
//       actionId: get_customer

// actions:
//   - id: get_customer_mock
//     type: httpRequestMock
//     http:
//       mock:
//         body: "{ 'Say':'Hello' }"
//         contentType: "application/json"
//         statusCode: 201
//         headers:
//           "mock1": "value1"
//           "mock2": "value2"
//           "mock3": "value2"

//   - id: get_customer
//     type: httpRequest
//     http:
//       request:
//         method: get
//         url: 'http://localhost:8085/api/mock'
//         mapRequestHeaders:
//           - "mock1"
//           - "mock2:change2"
//         mapFixedHeaders:
//           "req1": "value1"
//           "req2": "value2"
//       response:
//         mapResponseHeaders:
//           - "mock1"
//           - "mock2:change2"
//         mapFixedHeaders:
//           "resp1": "value1"
//           "resp2": "value2"

// `
// 	tmpFile, err := os.CreateTemp("", "configApp*.yaml")
// 	if err != nil {
// 		t.Fatalf("Could not create temp file: %v", err)
// 	}
// 	defer os.Remove(tmpFile.Name())

// 	if _, err := tmpFile.Write([]byte(yamlContent)); err != nil {
// 		t.Fatalf("Failed to write to temp file: %v", err)
// 	}

// 	if err := tmpFile.Close(); err != nil {
// 		t.Fatalf("Failed to close temp file: %v", err)
// 	}

// 	originalFilePath := "configApp.yaml"
// 	os.Rename(tmpFile.Name(), originalFilePath)
// 	defer os.Rename(originalFilePath, tmpFile.Name())

// 	config, err := LoadYamlFile("configApp.yaml")
// 	if err != nil {
// 		t.Fatalf("LoadYamlFile failed: %v", err)
// 	}

// 	if config.Service.Name != "test-app" {
// 		t.Errorf("Expected service name 'test-app', got '%s'", config.Service.Name)
// 	}

// 	if config.Api.Endpoints[0].Route != "api/customer" {
// 		t.Errorf("Expected API route 'api/customer', got '%s'", config.Api.Endpoints[0].Route)
// 	}

// 	if config.Actions[0].Mock.Value != "Hi" {
// 		t.Errorf("Expected mock value 'Hi', got '%s'", config.Actions[0].Mock.Value)
// 	}
// }

// func TestLoadYamlFile_FileNotFound(t *testing.T) {
// 	originalFilePath := "configApp.yaml"
// 	os.Rename(originalFilePath, originalFilePath+".bak")
// 	defer os.Rename(originalFilePath+".bak", originalFilePath)

// 	_, err := LoadYamlFile("configApp.yaml")
// 	if err == nil {
// 		t.Fatalf("Expected error when file is not found, but got nil")
// 	}
// }

// func TestLoadYamlFile_InvalidYaml(t *testing.T) {
// 	yamlContent := `
// version: 1

// service:
//   name: test-app
//   version: 1.0.0

// api:
//   endpoints:
//     - route: /api/mock
//       method: get
//       actionId: get_customer_mock

//     - route: /api/customer
//       method: get
//       actionId: get_customer

// actions:
//   - id: get_customer_mock
//     type: httpRequestMock
//     http:
//       mock:
//         body: "{ 'Say':'Hello' }"
//         contentType: "application/json"
//         statusCode: 201
//         headers:
//           "mock1": "value1"
//           "mock2": "value2"
//           "mock3": "value2"

//   - id: get_customer
//     type: httpRequest
//     http:
//       request:
//         method: get
//         url: 'http://localhost:8085/api/mock'
//         mapRequestHeaders:
//           - "mock1"
//           - "mock2:change2"
//         mapFixedHeaders:
//           "req1": "value1"
//           "req2": "value2"
//       response:
//         mapResponseHeaders:
//           - "mock1"
//           - "mock2:change2"
//         mapFixedHeaders:
//           "resp1": "value1"
//           "resp2": "value2"
// `
// 	tmpFile, err := os.CreateTemp("", "configApp*.yaml")
// 	if err != nil {
// 		t.Fatalf("Could not create temp file: %v", err)
// 	}
// 	defer os.Remove(tmpFile.Name())

// 	if _, err := tmpFile.Write([]byte(yamlContent)); err != nil {
// 		t.Fatalf("Failed to write to temp file: %v", err)
// 	}

// 	if err := tmpFile.Close(); err != nil {
// 		t.Fatalf("Failed to close temp file: %v", err)
// 	}

// 	originalFilePath := "configApp.yaml"
// 	os.Rename(tmpFile.Name(), originalFilePath)
// 	defer os.Rename(originalFilePath, tmpFile.Name())

// 	_, err = LoadYamlFile("configApp.yaml")
// 	if err == nil {
// 		t.Fatalf("Expected error when file is not found, but got nil")
// 	}
// }
