package json_map

import (
	"strings"
	"time"
	"wrench/app/contexts"
	"wrench/app/manifest/contract_settings/maps"

	"github.com/google/uuid"
)

func GetValue(jsonMap map[string]interface{}, propertyName string, deleteProperty bool) (interface{}, map[string]interface{}) {
	var value interface{}

	var jsonMapCurrent map[string]interface{}
	jsonMapCurrent = jsonMap
	propertyNameSplitted := strings.Split(propertyName, ".")
	index := 0
	for i, property := range propertyNameSplitted {
		index++
		valueObject, ok := jsonMapCurrent[property].(map[string]interface{})
		if ok {
			jsonMapCurrent = valueObject
			if i == index {
				value = jsonMapCurrent
				break
			}
			continue
		}

		valueTempString, ok := jsonMapCurrent[property].(string)
		if ok {
			value = valueTempString

			if deleteProperty {
				delete(jsonMapCurrent, property)
			}

			break
		}
	}
	return value, jsonMap
}

func getArrayValue(jsonMap map[string]interface{}, propertyName string) []map[string]interface{} {
	var value []map[string]interface{}

	jsonMapCurrent := jsonMap
	propertyNameSplitted := strings.Split(propertyName, ".")
	index := 0
	for i, property := range propertyNameSplitted {

		valueObject, ok := jsonMapCurrent[property].(map[string]interface{})
		if ok {
			jsonMapCurrent = valueObject
			continue
		}

		valueArray, ok := jsonMapCurrent[property].([]map[string]interface{})
		if ok {
			if i == index {
				value = valueArray
				break
			}
		}
		index++
	}
	return value
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
		valueTemp, ok := jsonMapCurrent[property].(map[string]interface{})

		if ok && i+1 == total {
			delete(jsonMapCurrent, property)
			break
		} else {
			jsonMapCurrent = valueTemp
		}

		if i+1 == total {
			delete(jsonMapCurrent, property)
		}
	}

	return jsonMap
}

func CreatePropertiesInterpolationValue(jsonMap map[string]interface{}, propertiesValues []string, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) map[string]interface{} {
	jsonValueCurrent := jsonMap
	for _, propertyValue := range propertiesValues {
		propertyValueSplitted := strings.Split(propertyValue, ":")
		propertyName := propertyValueSplitted[0]
		valueArray := propertyValueSplitted[1:]
		value := strings.Join(valueArray, ":")
		jsonValueCurrent = CreatePropertyInterpolationValue(jsonValueCurrent, propertyName, value, wrenchContext, bodyContext)
	}
	return jsonValueCurrent
}

func calculatedValue(value string) bool {
	return strings.HasPrefix(value, "{{") && strings.HasSuffix(value, "}}")
}

func replaceCalculatedValue(command string) string {
	return strings.ReplaceAll(strings.ReplaceAll(command, "{{", ""), "}}", "")
}

func CreatePropertyInterpolationValue(jsonMap map[string]interface{}, propertyName string, value string, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) map[string]interface{} {
	valueResult := value

	if calculatedValue(value) {
		rawValue := replaceCalculatedValue(value)

		if rawValue == "uuid" {
			valueResult = uuid.New().String()
		} else if strings.HasPrefix(rawValue, "time") {
			timeFormat := strings.ReplaceAll(rawValue, "time ", "")
			timeNow := time.Now()

			if len(timeFormat) > 0 {
				valueResult = timeNow.Format(timeFormat)
			} else {
				valueResult = timeNow.String()
			}
		} else if strings.HasPrefix(rawValue, "wrenchContext") {
			valueResult = getValueWrenchContext(rawValue, wrenchContext)
		}

	}

	return CreateProperty(jsonMap, propertyName, valueResult)
}

const wrenchContextRequestHeaders = "wrenchContext.request.headers."

func getValueWrenchContext(command string, wrenchContext *contexts.WrenchContext) string {

	if strings.HasPrefix(command, wrenchContextRequestHeaders) {
		headerName := strings.ReplaceAll(command, wrenchContextRequestHeaders, "")
		return wrenchContext.Request.Header.Get(headerName)
	}

	return ""
}

const bodyContext = "bodyContext."

func ParseValues(jsonMap map[string]interface{}, parse *maps.ParseSettings) map[string]interface{} {

	jsonValueCurrent := jsonMap

	if len(parse.WhenEquals) > 0 {
		jsonValueCurrent = parseWhenEqualValue(jsonValueCurrent, parse.WhenEquals)
	} else if len(parse.Operator) > 0 {
		jsonValueCurrent = parseOperatorValues(jsonValueCurrent, parse.Operator)
	}

	return jsonValueCurrent
}

func parseWhenEqualValue(jsonMap map[string]interface{}, whenEqual []string) map[string]interface{} {
	jsonValueCurrent := jsonMap

	for _, whenEqual := range whenEqual {
		if calculatedValue(whenEqual) {
			whenEqual = strings.ReplaceAll(whenEqual, bodyContext, "")
			rawWhenEqual := replaceCalculatedValue(whenEqual)

			whenEqualSplitted := strings.Split(rawWhenEqual, ":")
			propertyNameWithEqualValue := whenEqualSplitted[0]
			propertyNameWithEqualValueSplitted := strings.Split(propertyNameWithEqualValue, ".")

			lenWithEqual := len(propertyNameWithEqualValueSplitted)

			valueArray := propertyNameWithEqualValueSplitted[:lenWithEqual-1]

			propertyName := strings.Join(valueArray, ".")
			equalValue := propertyNameWithEqualValueSplitted[lenWithEqual-1] // value to compare

			parseToValue := whenEqualSplitted[1] // value if equals should be used

			valueCurrent, _ := GetValue(jsonMap, propertyName, false)

			if valueCurrent == equalValue {
				jsonValueCurrent = SetValue(jsonValueCurrent, propertyName, parseToValue)
			}
		}
	}

	return jsonValueCurrent
}

func parseOperatorValues(jsonMap map[string]interface{}, operator []maps.Operator) map[string]interface{} {
	jsonValueCurrent := jsonMap

	for _, operatorValue := range operator {

		operatorValueSplitted := strings.Split(string(operatorValue), ".")
		operatorValue := operatorValueSplitted[0]
		operatorCommand := operatorValueSplitted[1]

		if calculatedValue(operatorValue) {
			operatorValue = replaceCalculatedValue(operatorValue)

			if operatorValue == "bodyContext" {
				if operatorCommand == "to_array" {
					arrayJsonMap := make([]map[string]interface{}, 1)
					arrayJsonMap[0] = jsonMap

					jsonMapTemp := map[string]interface{}{
						"root": arrayJsonMap,
					}
					jsonValueCurrent = jsonMapTemp
				}
			}

		} else {
			propertyName := operatorValue
			value, jsonValueCurrent := GetValue(jsonMap, propertyName, true)

			if operatorCommand == "to_array" {
				arrayJsonMap := make([]interface{}, 1)
				arrayJsonMap[0] = value
				jsonValueCurrent = SetValue(jsonValueCurrent, propertyName, arrayJsonMap)
			} else if operatorCommand == "get_first" {
				arrayValue := json_map.Get
			}
		}
	}

	return jsonValueCurrent
}

func SetObjectRoot(jsonMap map[string]interface{}, propertyName string) map[string]interface{} {
	value, _ := GetValue(jsonMap, propertyName, false)

	objectValue, isObject := value.(map[string]interface{})
	if isObject {
		return objectValue
	}
	return nil
}

func SetArrayRoot(jsonMap map[string]interface{}, propertyName string) []map[string]interface{} {
	return getArrayValue(jsonMap, propertyName)
}
