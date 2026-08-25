package contexts

import (
	"encoding/json"
	"testing"

	"wrench/app/manifest/contract_settings/maps"
)

func TestParseValuesWhenEqualsSingleObjectKeepsPreviousBehavior(t *testing.T) {
	payload := []byte(`{ "address": { "type": "Mobile" }, "name": "John" }`)

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(payload, &jsonMap); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	result := ParseValues(jsonMap, &maps.ParseSettings{
		WhenEquals: []string{
			"{{bodyContext.address.type.Mobile:Movel}}",
		},
	})

	address := result["address"].(map[string]interface{})
	if address["type"] != "Movel" {
		t.Errorf("expected type 'Movel', got %v", address["type"])
	}

	if result["name"] != "John" {
		t.Errorf("expected name intact, got %v", result["name"])
	}
}

func TestParseValuesWhenEqualsForeachList(t *testing.T) {
	payload := []byte(`{
	  "Id": 1,
	  "BankAccounts": [
	    { "Id": 1, "Arranjo": "Visa" },
	    { "Id": 2, "Arranjo": "Mastercard" },
	    { "Id": 3, "Arranjo": "Elo" }
	  ]
	}`)

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(payload, &jsonMap); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	result := ParseValues(jsonMap, &maps.ParseSettings{
		WhenEquals: []string{
			"{{BankAccounts.Arranjo.Visa:1}}",
			"{{BankAccounts.Arranjo.Mastercard:2}}",
		},
	})

	accounts := result["BankAccounts"].([]interface{})
	expected := []string{"1", "2", "Elo"}
	for i, itemValue := range accounts {
		item := itemValue.(map[string]interface{})
		if item["Arranjo"] != expected[i] {
			t.Errorf("index %d: expected Arranjo %q, got %v", i, expected[i], item["Arranjo"])
		}
	}
}

func TestParseValuesWhenEqualsForeachListWithBodyContextPrefix(t *testing.T) {
	payload := []byte(`{
	  "BankAccounts": [
	    { "Id": 1, "Arranjo": "Visa" },
	    { "Id": 2, "Arranjo": "Mastercard" }
	  ]
	}`)

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(payload, &jsonMap); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	result := ParseValues(jsonMap, &maps.ParseSettings{
		WhenEquals: []string{
			"{{bodyContext.BankAccounts.Arranjo.Visa:1}}",
			"{{bodyContext.BankAccounts.Arranjo.Mastercard:2}}",
		},
	})

	accounts := result["BankAccounts"].([]interface{})
	expected := []string{"1", "2"}
	for i, itemValue := range accounts {
		item := itemValue.(map[string]interface{})
		if item["Arranjo"] != expected[i] {
			t.Errorf("index %d: expected Arranjo %q, got %v", i, expected[i], item["Arranjo"])
		}
	}
}

func TestCreatePropertiesInterpolationValueInsideListItems(t *testing.T) {
	payload := []byte(`{
	  "ProposalId": "124",
	  "BankAccounts": [
	    { "Arranjo": "Visa", "numero": "111" },
	    { "Arranjo": "Mastercard", "numero": "222" },
	    { "Arranjo": "Hipersom", "numero": "444" }
	  ]
	}`)

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(payload, &jsonMap); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	result := CreatePropertiesInterpolationValue(jsonMap, []string{"BankAccounts.Teste:Qualquercoisa"}, nil, nil, nil)

	accounts, ok := result["BankAccounts"].([]interface{})
	if !ok {
		t.Fatalf("expected BankAccounts to remain an array, got %T", result["BankAccounts"])
	}
	for i, itemValue := range accounts {
		item := itemValue.(map[string]interface{})
		if item["Teste"] != "Qualquercoisa" {
			t.Errorf("index %d: expected Teste 'Qualquercoisa', got %v", i, item["Teste"])
		}
	}
	if result["ProposalId"] != "124" {
		t.Errorf("expected ProposalId intact, got %v", result["ProposalId"])
	}
}
