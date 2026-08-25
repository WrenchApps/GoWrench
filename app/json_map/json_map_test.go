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

func TestCreatePropertyInsideArrayItems(t *testing.T) {
	jsonMap := map[string]interface{}{
		"ProposalId": "124",
		"BankAccounts": []interface{}{
			map[string]interface{}{"Arranjo": "Visa", "numero": "111"},
			map[string]interface{}{"Arranjo": "Mastercard", "numero": "222"},
		},
	}

	result := CreateProperty(jsonMap, "BankAccounts.Teste", "Qualquercoisa")

	accounts, ok := result["BankAccounts"].([]interface{})
	if !ok {
		t.Fatalf("expected BankAccounts to remain an array")
	}
	if len(accounts) != 2 {
		t.Fatalf("expected 2 items, got %d", len(accounts))
	}
	for i, itemValue := range accounts {
		item := itemValue.(map[string]interface{})
		if item["Teste"] != "Qualquercoisa" {
			t.Errorf("index %d: expected Teste 'Qualquercoisa', got %v", i, item["Teste"])
		}
		if item["numero"] == nil {
			t.Errorf("index %d: expected existing properties intact", i)
		}
	}
	if result["ProposalId"] != "124" {
		t.Errorf("expected ProposalId intact, got %v", result["ProposalId"])
	}
}

func TestCreatePropertyNestedArray(t *testing.T) {
	jsonMap := map[string]interface{}{
		"company": map[string]interface{}{
			"partners": []interface{}{
				map[string]interface{}{"Name": "A"},
				map[string]interface{}{"Name": "B"},
			},
		},
	}

	result := CreateProperty(jsonMap, "company.partners.Country", "BR")

	partners := result["company"].(map[string]interface{})["partners"].([]interface{})
	for i, itemValue := range partners {
		item := itemValue.(map[string]interface{})
		if item["Country"] != "BR" {
			t.Errorf("partner %d: expected Country 'BR', got %v", i, item["Country"])
		}
	}
}

func TestCreatePropertyCreatesIntermediateMaps(t *testing.T) {
	jsonMap := map[string]interface{}{
		"name": "John",
	}

	result := CreateProperty(jsonMap, "address.city.name", "São Paulo")

	address := result["address"].(map[string]interface{})
	city := address["city"].(map[string]interface{})
	if city["name"] != "São Paulo" {
		t.Errorf("expected name 'São Paulo', got %v", city["name"])
	}
	if result["name"] != "John" {
		t.Errorf("expected name intact, got %v", result["name"])
	}
}

func TestCreatePropertyReplacesScalarIntermediate(t *testing.T) {
	jsonMap := map[string]interface{}{
		"address": "some string",
	}

	result := CreateProperty(jsonMap, "address.city.name", "São Paulo")

	city := result["address"].(map[string]interface{})["city"].(map[string]interface{})
	if city["name"] != "São Paulo" {
		t.Errorf("expected name 'São Paulo', got %v", city["name"])
	}
}

func listFixture() map[string]interface{} {
	return map[string]interface{}{
		"ProposalId": "124",
		"BankAccounts": []interface{}{
			map[string]interface{}{"Arranjo": "Visa", "numero": "111"},
			map[string]interface{}{"Arranjo": "Mastercard", "numero": "222"},
			map[string]interface{}{"Arranjo": "Hipersom", "numero": "444"},
		},
	}
}

func TestRenamePropertyInsideListItems(t *testing.T) {
	result := RenameProperty(listFixture(), "BankAccounts.numero", "BankAccounts.Numbers")

	accounts, ok := result["BankAccounts"].([]interface{})
	if !ok {
		t.Fatalf("expected BankAccounts to remain an array")
	}
	for i, itemValue := range accounts {
		item := itemValue.(map[string]interface{})
		if _, exists := item["numero"]; exists {
			t.Errorf("index %d: expected numero removed", i)
		}
		if item["Numbers"] == nil {
			t.Errorf("index %d: expected Numbers set", i)
		}
		if item["Arranjo"] == nil {
			t.Errorf("index %d: expected Arranjo intact", i)
		}
	}
	if result["ProposalId"] != "124" {
		t.Errorf("expected ProposalId intact, got %v", result["ProposalId"])
	}
}

func TestRenameTopLevelList(t *testing.T) {
	result := RenameProperty(listFixture(), "BankAccounts", "Accounts")

	if result["BankAccounts"] != nil {
		t.Error("expected BankAccounts removed")
	}
	accounts, ok := result["Accounts"].([]interface{})
	if !ok || len(accounts) != 3 {
		t.Fatalf("expected Accounts array with 3 items")
	}
}

func TestRenamePropertyNestedObjectNonList(t *testing.T) {
	jsonMap := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "value",
			},
		},
	}

	result := RenameProperty(jsonMap, "a.b.c", "a.b.d")

	b := result["a"].(map[string]interface{})["b"].(map[string]interface{})
	if b["d"] != "value" {
		t.Errorf("expected d 'value', got %v", b["d"])
	}
	if _, exists := b["c"]; exists {
		t.Error("expected c removed")
	}
}

func TestDuplicatePropertyInsideListItems(t *testing.T) {
	result := DuplicatePropertyValue(listFixture(), "BankAccounts.numero", "BankAccounts.NumeroCopy")

	accounts := result["BankAccounts"].([]interface{})
	for i, itemValue := range accounts {
		item := itemValue.(map[string]interface{})
		if item["NumeroCopy"] == nil {
			t.Errorf("index %d: expected NumeroCopy set", i)
		}
		if item["numero"] == nil {
			t.Errorf("index %d: expected original numero intact", i)
		}
	}
}

func TestRemovePropertyInsideListItems(t *testing.T) {
	result := RemoveProperty(listFixture(), "BankAccounts.numero")

	accounts, ok := result["BankAccounts"].([]interface{})
	if !ok {
		t.Fatalf("expected BankAccounts to remain an array")
	}
	for i, itemValue := range accounts {
		item := itemValue.(map[string]interface{})
		if _, exists := item["numero"]; exists {
			t.Errorf("index %d: expected numero removed", i)
		}
		if item["Arranjo"] == nil {
			t.Errorf("index %d: expected Arranjo intact", i)
		}
	}
}
