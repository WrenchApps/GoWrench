package json_map

import "testing"

func TestSetValueWhenEqualsSingleObject(t *testing.T) {
	jsonMap := map[string]interface{}{
		"address": map[string]interface{}{
			"type": "Mobile",
		},
		"name": "John",
	}

	result := SetValueWhenEquals(jsonMap, "address.type", "Mobile", "Móvel")

	address := result["address"].(map[string]interface{})
	if address["type"] != "Móvel" {
		t.Errorf("expected type 'Móvel', got %v", address["type"])
	}

	if _, exists := result["name"]; !exists {
		t.Error("expected unrelated property intact")
	}
}

func TestSetValueWhenEqualsSingleObjectNoMatch(t *testing.T) {
	jsonMap := map[string]interface{}{
		"address": map[string]interface{}{
			"type": "Residential",
		},
	}

	result := SetValueWhenEquals(jsonMap, "address.type", "Mobile", "Móvel")

	address := result["address"].(map[string]interface{})
	if address["type"] != "Residential" {
		t.Errorf("expected type untouched, got %v", address["type"])
	}
}

func TestSetValueWhenEqualsForeachListOnlyMatching(t *testing.T) {
	jsonMap := map[string]interface{}{
		"BankAccounts": []interface{}{
			map[string]interface{}{"Id": float64(1), "Arranjo": "Visa"},
			map[string]interface{}{"Id": float64(2), "Arranjo": "Mastercard"},
			map[string]interface{}{"Id": float64(3), "Arranjo": "Elo"},
		},
	}

	result := SetValueWhenEquals(jsonMap, "BankAccounts.Arranjo", "Visa", "1")

	pairs := []string{"1", "Mastercard", "Elo"}
	accounts := result["BankAccounts"].([]interface{})
	for i, itemValue := range accounts {
		item := itemValue.(map[string]interface{})
		if item["Arranjo"] != pairs[i] {
			t.Errorf("index %d: expected Arranjo %q, got %v", i, pairs[i], item["Arranjo"])
		}
		if item["Id"] != float64(i+1) {
			t.Errorf("index %d: expected Id intact, got %v", i, item["Id"])
		}
	}
}

func TestSetValueWhenEqualsNestedList(t *testing.T) {
	jsonMap := map[string]interface{}{
		"company": map[string]interface{}{
			"partners": []interface{}{
				map[string]interface{}{"Name": "EDUARDA", "Pep": "false"},
				map[string]interface{}{"Name": "DANIELA", "Pep": "false"},
			},
		},
	}

	result := SetValueWhenEquals(jsonMap, "company.partners.Pep", "false", "true")

	partners := result["company"].(map[string]interface{})["partners"].([]interface{})
	for i, itemValue := range partners {
		item := itemValue.(map[string]interface{})
		if item["Pep"] != "true" {
			t.Errorf("partner %d: expected Pep 'true', got %v", i, item["Pep"])
		}
	}
}

func TestSetValueWhenEqualsKeepsPreviousBehaviorNonStringNotMatched(t *testing.T) {
	jsonMap := map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"Code": float64(51)},
			map[string]interface{}{"Code": true},
		},
	}

	result := SetValueWhenEquals(jsonMap, "items.Code", "51", "replaced")

	items := result["items"].([]interface{})
	for i, itemValue := range items {
		item := itemValue.(map[string]interface{})
		if item["Code"] != float64(51) && item["Code"] != true {
			t.Errorf("index %d: expected value untouched, got %v", i, item["Code"])
		}
	}
}
