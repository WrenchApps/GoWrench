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
					if string(propertyNameToArray[0]) == "." {
						propertyNameToArray = propertyNameToArray[1:]
						return GetValue(item, propertyNameToArray, false)
					} else {
						value = jsonMapCurrent[property]
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

func SetValue(jsonMap map[string]interface{}, propertyName string, newValue string) map[string]interface{} {
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

	var jsonMapCurrent map[string]interface{}
	jsonMapCurrent = jsonMap
	propertyNameSplitted := strings.Split(propertyName, ".")
	total := len(propertyNameSplitted)

	for i, property := range propertyNameSplitted {
		valueTemp, ok := jsonMapCurrent[property].(map[string]interface{})
		if ok {
			jsonMapCurrent = valueTemp
		} else {
			if i+1 < total {
				jsonMapNew := make(map[string]interface{})
				jsonMapCurrent[property] = jsonMapNew
				jsonMapCurrent = jsonMapNew
			}
		}

		if i+1 == total {
			jsonMapCurrent[property] = value
		}
	}
	return jsonMap
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
	value, jsonValue := GetValue(jsonMap, propertyNameSource, false)
	return CreateProperty(jsonValue, propertyNameDestination, value)
}

func RenameProperty(jsonMap map[string]interface{}, propertyNameOld string, propertyNameNew string) map[string]interface{} {
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

func RemoveProperty(jsonMap map[string]interface{}, propertyName string) map[string]interface{} {
	var jsonMapCurrent map[string]interface{}
	jsonMapCurrent = jsonMap

	propertyNameSplitted := strings.Split(propertyName, ".")
	total := len(propertyNameSplitted)

	for i, property := range propertyNameSplitted {
		if i == total-1 {
			delete(jsonMapCurrent, property)
			break
		}

		valueTemp, ok := jsonMapCurrent[property].(map[string]interface{})

		if !ok {
			break
		}

		jsonMapCurrent = valueTemp
	}

	return jsonMap
}
