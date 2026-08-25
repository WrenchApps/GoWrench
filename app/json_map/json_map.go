package json_map

import (
	"regexp"
	"strconv"
	"strings"
)

func GetValue(jsonMap map[string]interface{}, propertyName string, deleteProperty bool) (interface{}, map[string]interface{}) {
	var value interface{}

	var jsonMapCurrent map[string]interface{}
	jsonMapCurrent = jsonMap
	propertyNameSplitted := strings.Split(propertyName, ".")
	totalProperty := len(propertyNameSplitted)

	for index, property := range propertyNameSplitted {
		valueTemp, ok := jsonMapCurrent[property].(map[string]interface{})
		if ok {
			if index == totalProperty-1 {
				value = valueTemp
				if deleteProperty {
					delete(jsonMapCurrent, property)
				}
				break
			}
			jsonMapCurrent = valueTemp
			continue
		}

		match := regexp.MustCompile(`(\w+)\[(\d+)\]`).FindStringSubmatch(property)
		if match != nil {
			indexValueStringArray := strings.Split(property, "[")
			indexValueStringFirst := indexValueStringArray[0]
			indexValueStringLast := indexValueStringArray[1]
			indexValueStringCut := indexValueStringLast[:len(indexValueStringLast)-1]
			indexValue, _ := strconv.ParseInt(indexValueStringCut, 10, 0)
			valueTempString, ok := jsonMapCurrent[indexValueStringFirst].([]interface{})
			if ok {
				propertyNameToArray := strings.ReplaceAll(propertyName, property, "")
				item, ok2 := valueTempString[indexValue].(map[string]interface{})
				if ok2 {
					if len(propertyNameToArray) == 0 {
						value = item
					} else {
						if string(propertyNameToArray[0]) == "." {
							propertyNameToArray = propertyNameToArray[1:]
							return GetValue(item, propertyNameToArray, false)
						} else {
							value = jsonMapCurrent[property]
						}
					}
				}
				break
			}
		} else {

			value = jsonMapCurrent[property]

			if deleteProperty {
				delete(jsonMapCurrent, property)
			}

			break

		}
	}
	return value, jsonMap
}

func SetValue(jsonMap map[string]interface{}, propertyName string, newValue interface{}) map[string]interface{} {
	var jsonMapCurrent map[string]interface{}
	jsonMapCurrent = jsonMap
	propertyNameSplitted := strings.Split(propertyName, ".")
	total := len(propertyNameSplitted)

	for i, property := range propertyNameSplitted {
		valueTemp, ok := jsonMapCurrent[property].(map[string]interface{})
		if ok {
			jsonMapCurrent = valueTemp
			continue
		}

		if i+1 == total {
			jsonMapCurrent[property] = newValue
		}
	}

	return jsonMap
}

func CreateProperty(jsonMap map[string]interface{}, propertyName string, value interface{}) map[string]interface{} {
	propertyNameSplitted := strings.Split(propertyName, ".")
	createProperty(jsonMap, propertyNameSplitted, value)
	return jsonMap
}

func createProperty(valueCurrent interface{}, propertyNameSplitted []string, value interface{}) {
	if len(propertyNameSplitted) == 0 {
		return
	}

	jsonArray, ok := valueCurrent.([]interface{})
	if ok {
		for _, item := range jsonArray {
			createProperty(item, propertyNameSplitted, value)
		}
		return
	}

	jsonMapCurrent, ok := valueCurrent.(map[string]interface{})
	if !ok {
		return
	}

	property := propertyNameSplitted[0]

	if len(propertyNameSplitted) == 1 {
		jsonMapCurrent[property] = value
		return
	}

	switch next := jsonMapCurrent[property].(type) {
	case []interface{}:
		createProperty(next, propertyNameSplitted[1:], value)
	case map[string]interface{}:
		createProperty(next, propertyNameSplitted[1:], value)
	default:
		jsonMapNew := make(map[string]interface{})
		jsonMapCurrent[property] = jsonMapNew
		createProperty(jsonMapNew, propertyNameSplitted[1:], value)
	}
}

func RenameProperties(jsonMap map[string]interface{}, properties []string) map[string]interface{} {
	jsonValueCurrent := jsonMap
	for _, property := range properties {
		propertyNameSplitted := strings.Split(property, ":")
		propertyNameOld := propertyNameSplitted[0]
		propertyNameNew := propertyNameSplitted[1]
		jsonValueCurrent = RenameProperty(jsonValueCurrent, propertyNameOld, propertyNameNew)
	}
	return jsonValueCurrent
}

func DuplicatePropertiesValue(jsonMap map[string]interface{}, properties []string) map[string]interface{} {
	jsonValueCurrent := jsonMap
	for _, property := range properties {
		propertyNameSplitted := strings.Split(property, ":")
		propertyNameSource := propertyNameSplitted[0]
		propertyNameDestination := propertyNameSplitted[1]
		jsonValueCurrent = DuplicatePropertyValue(jsonValueCurrent, propertyNameSource, propertyNameDestination)
	}
	return jsonValueCurrent
}

func DuplicatePropertyValue(jsonMap map[string]interface{}, propertyNameSource string, propertyNameDestination string) map[string]interface{} {
	sourceSplitted := strings.Split(propertyNameSource, ".")
	destinationSplitted := strings.Split(propertyNameDestination, ".")

	handled, sawArray := walkListPath(jsonMap, sourceSplitted, destinationSplitted, func(jsonMapCurrent map[string]interface{}, sourceKey string, destinationKey string) bool {
		value, exists := jsonMapCurrent[sourceKey]
		if !exists {
			return false
		}
		jsonMapCurrent[destinationKey] = value
		return true
	})

	if sawArray || handled {
		return jsonMap
	}

	value, jsonValue := GetValue(jsonMap, propertyNameSource, false)
	return CreateProperty(jsonValue, propertyNameDestination, value)
}

func RenameProperty(jsonMap map[string]interface{}, propertyNameOld string, propertyNameNew string) map[string]interface{} {
	oldSplitted := strings.Split(propertyNameOld, ".")
	newSplitted := strings.Split(propertyNameNew, ".")

	handled, sawArray := walkListPath(jsonMap, oldSplitted, newSplitted, func(jsonMapCurrent map[string]interface{}, oldKey string, newKey string) bool {
		value, exists := jsonMapCurrent[oldKey]
		if !exists {
			return false
		}
		delete(jsonMapCurrent, oldKey)
		jsonMapCurrent[newKey] = value
		return true
	})

	if sawArray || handled {
		return jsonMap
	}

	value, jsonValue := GetValue(jsonMap, propertyNameOld, true)
	return CreateProperty(jsonValue, propertyNameNew, value)
}

func RemoveProperties(jsonMap map[string]interface{}, propertiesName []string) map[string]interface{} {
	if propertiesName == nil {
		return nil
	}

	currentJsonValue := jsonMap
	for _, property := range propertiesName {
		currentJsonValue = RemoveProperty(currentJsonValue, property)
	}

	return currentJsonValue
}

// walkListPath walks source and destination paths in parallel while they share
// the same structure, fanning out to every item whenever a list is found. The
// apply function runs at the leaf parent of each branch.
func walkListPath(valueCurrent interface{}, sourceSplitted []string, destinationSplitted []string, apply func(jsonMapCurrent map[string]interface{}, sourceKey string, destinationKey string) bool) (handled bool, sawArray bool) {
	if len(sourceSplitted) == 0 || len(sourceSplitted) != len(destinationSplitted) {
		return false, false
	}

	jsonArray, ok := valueCurrent.([]interface{})
	if ok {
		handled = len(jsonArray) > 0
		sawArray = true
		for _, item := range jsonArray {
			itemHandled, _ := walkListPath(item, sourceSplitted, destinationSplitted, apply)
			if !itemHandled {
				handled = false
			}
		}
		return handled, sawArray
	}

	jsonMapCurrent, ok := valueCurrent.(map[string]interface{})
	if !ok {
		return false, false
	}

	sourceKey := sourceSplitted[0]
	destinationKey := destinationSplitted[0]

	if len(sourceSplitted) == 1 {
		return apply(jsonMapCurrent, sourceKey, destinationKey), false
	}

	if sourceKey != destinationKey {
		return false, false
	}

	child, exists := jsonMapCurrent[sourceKey]
	if !exists {
		return false, false
	}

	return walkListPath(child, sourceSplitted[1:], destinationSplitted[1:], apply)
}

func SetValueWhenEquals(jsonMap map[string]interface{}, propertyName string, expectedValue string, newValue string) map[string]interface{} {
	propertyNameSplitted := strings.Split(propertyName, ".")
	setValueWhenEquals(jsonMap, propertyNameSplitted, expectedValue, newValue)
	return jsonMap
}

func setValueWhenEquals(valueCurrent interface{}, propertyNameSplitted []string, expectedValue string, newValue string) {
	if len(propertyNameSplitted) == 0 {
		return
	}

	jsonArray, ok := valueCurrent.([]interface{})
	if ok {
		for _, item := range jsonArray {
			setValueWhenEquals(item, propertyNameSplitted, expectedValue, newValue)
		}
		return
	}

	jsonMapCurrent, ok := valueCurrent.(map[string]interface{})
	if !ok {
		return
	}

	property := propertyNameSplitted[0]
	valueTemp := jsonMapCurrent[property]

	if len(propertyNameSplitted) == 1 {
		if currentValue, isString := valueTemp.(string); isString && currentValue == expectedValue {
			jsonMapCurrent[property] = newValue
		}
		return
	}

	setValueWhenEquals(valueTemp, propertyNameSplitted[1:], expectedValue, newValue)
}

func RemoveProperty(jsonMap map[string]interface{}, propertyName string) map[string]interface{} {
	propertyNameSplitted := strings.Split(propertyName, ".")

	walkListPath(jsonMap, propertyNameSplitted, propertyNameSplitted, func(jsonMapCurrent map[string]interface{}, sourceKey string, _ string) bool {
		delete(jsonMapCurrent, sourceKey)
		return true
	})

	return jsonMap
}
