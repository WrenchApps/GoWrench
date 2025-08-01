package contexts

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"wrench/app/json_map"
	settings "wrench/app/manifest/action_settings"
	"wrench/app/manifest/action_settings/func_settings"
	"wrench/app/manifest/contract_settings/maps"

	"github.com/google/uuid"
)

const prefixWrenchContextRequest = "wrenchContext.request."
const prefixWrenchContextRequestUri = "wrenchContext.request.uri"
const prefixWrenchContextRequestUriParams = "wrenchContext.request.uri.params."
const prefixWrenchContextRequestHeaders = "wrenchContext.request.headers."
const prefixBodyContext = "bodyContext."
const prefixBodyContextPreserved = "bodyContext.actions."
const prefixFunc = "func."

func IsCalculatedValue(value string) bool {
	return strings.HasPrefix(value, "{{") && strings.HasSuffix(value, "}}")
}

func ReplaceCalculatedValue(command string) string {
	return strings.ReplaceAll(strings.ReplaceAll(command, "{{", ""), "}}", "")
}

func ReplacePrefixBodyContextPreserved(command string) string {
	return strings.ReplaceAll(command, prefixBodyContextPreserved, "")
}

func IsWrenchContextCommand(command string) bool {
	return strings.HasPrefix(command, prefixWrenchContextRequest)
}

func IsBodyContextCommand(command string) bool {
	return strings.HasPrefix(command, prefixBodyContext)
}

func IsFunc(command string) bool {
	return strings.HasPrefix(command, prefixFunc)
}

func GetRequestUriParams(wrenchContext *WrenchContext, parameterName string) string {
	uriSplited := strings.Split(wrenchContext.Request.RequestURI, "/")
	routeSplited := strings.Split(wrenchContext.Endpoint.Route, "/")

	for i, routeValue := range routeSplited {
		if routeValue == fmt.Sprintf("{%s}", parameterName) {
			return uriSplited[i]
		}
	}

	return ""
}

func GetValueWrenchContext(command string, wrenchContext *WrenchContext) string {

	if IsCalculatedValue(command) {
		command = ReplaceCalculatedValue(command)
	}

	if strings.HasPrefix(command, prefixWrenchContextRequestHeaders) {
		headerName := strings.ReplaceAll(command, prefixWrenchContextRequestHeaders, "")
		return wrenchContext.Request.Header.Get(headerName)
	}

	if strings.HasPrefix(command, prefixWrenchContextRequestUriParams) {
		parameterName := strings.ReplaceAll(command, prefixWrenchContextRequestUriParams, "")
		return GetRequestUriParams(wrenchContext, parameterName)
	}

	if strings.HasPrefix(command, prefixWrenchContextRequestUri) {
		return wrenchContext.Request.RequestURI
	}

	return ""
}

func ReplacePrefixBodyContext(command string) string {
	if strings.HasPrefix(command, prefixBodyContext) {
		command = strings.ReplaceAll(command, prefixBodyContext, "")
	}
	return command
}

func GetCalculatedValue(command string, wrenchContext *WrenchContext, bodyContext *BodyContext, action *settings.ActionSettings) interface{} {
	if IsCalculatedValue(command) {
		command = ReplaceCalculatedValue(command)
		if IsBodyContextCommand(command) {
			return GetValueBodyContext(command, bodyContext)
		} else if IsWrenchContextCommand(command) {
			return GetValueWrenchContext(command, wrenchContext)
		} else if IsFunc(command) {
			return GetFuncValue(func_settings.FuncGeneralType(command), wrenchContext, bodyContext, action)
		} else {
			return command
		}
	} else {
		return command
	}
}

func GetValueBodyContext(command string, bodyContext *BodyContext) interface{} {

	if IsCalculatedValue(command) {
		command = ReplaceCalculatedValue(command)
	}

	if strings.HasPrefix(command, prefixBodyContextPreserved) {
		bodyPreservedMap := strings.ReplaceAll(command, prefixBodyContextPreserved, "")
		bodyPreservedMapSplitted := strings.Split(bodyPreservedMap, ".")
		actionId := bodyPreservedMapSplitted[0]
		if len(bodyPreservedMapSplitted) == 1 {
			bodyPreserved := bodyContext.GetBodyPreserved(actionId)
			return string(bodyPreserved)
		} else {
			jsonMap := bodyContext.ParseBodyToMapObjectPreserved(actionId)
			propertyName := strings.ReplaceAll(bodyPreservedMap, actionId+".", "")
			value, _ := json_map.GetValue(jsonMap, propertyName, false)
			return value
		}

	} else if strings.HasPrefix(command, prefixBodyContext) {
		propertyName := strings.ReplaceAll(command, prefixBodyContext, "")
		jsonMap := bodyContext.ParseBodyToMapObject()
		value, _ := json_map.GetValue(jsonMap, propertyName, false)
		if (value == nil || len(fmt.Sprint(value)) == 0) && propertyName == "currentBody" {
			value = bodyContext.GetBodyString()
		}
		return value
	}

	return ""
}

func GetCalculatedMap(mapConfigured map[string]string, wrenchContext *WrenchContext, bodyContext *BodyContext, action *settings.ActionSettings) map[string]interface{} {
	if mapConfigured == nil {
		return nil
	}
	mapResult := make(map[string]interface{})

	for key, value := range mapConfigured {
		mapResult[key] = GetCalculatedValue(value, wrenchContext, bodyContext, action)
	}

	return mapResult
}

func CreatePropertiesInterpolationValue(jsonMap map[string]interface{}, propertiesValues []string, wrenchContext *WrenchContext, bodyContext *BodyContext) map[string]interface{} {
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

func CreatePropertyInterpolationValue(jsonMap map[string]interface{}, propertyName string, value interface{}, wrenchContext *WrenchContext, bodyContext *BodyContext) map[string]interface{} {
	valueResult := value
	valueString := fmt.Sprint(valueResult)

	if IsCalculatedValue(valueString) {

		rawValue := ReplaceCalculatedValue(valueString)

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
			valueResult = GetValueWrenchContext(rawValue, wrenchContext)
		} else if strings.HasPrefix(rawValue, "bodyContext") {
			valueResult = GetValueBodyContext(rawValue, bodyContext)
		}
	}

	return json_map.CreateProperty(jsonMap, propertyName, valueResult)
}

func ParseValues(jsonMap map[string]interface{}, parse *maps.ParseSettings) map[string]interface{} {
	jsonValueCurrent := jsonMap
	if parse.WhenEquals != nil {
		for _, whenEqual := range parse.WhenEquals {
			if IsCalculatedValue(whenEqual) {
				whenEqual = ReplacePrefixBodyContext(whenEqual)
				rawWhenEqual := ReplaceCalculatedValue(whenEqual)

				whenEqualSplitted := strings.Split(rawWhenEqual, ":")
				propertyNameWithEqualValue := whenEqualSplitted[0]
				propertyNameWithEqualValueSplitted := strings.Split(propertyNameWithEqualValue, ".")

				lenWithEqual := len(propertyNameWithEqualValueSplitted)

				valueArray := propertyNameWithEqualValueSplitted[:lenWithEqual-1]

				propertyName := strings.Join(valueArray, ".")
				equalValue := propertyNameWithEqualValueSplitted[lenWithEqual-1] // value to compare

				parseToValue := whenEqualSplitted[1] // value if equals should be used

				valueCurrent, _ := json_map.GetValue(jsonMap, propertyName, false)

				if valueCurrent == equalValue {
					jsonValueCurrent = json_map.SetValue(jsonValueCurrent, propertyName, parseToValue)
				}
			}
		}
	}

	if len(parse.ToArray) > 0 {
		for _, toArray := range parse.ToArray {
			toArraySplitted := strings.Split(toArray, ":")

			originPropertyName := toArraySplitted[0]
			var destinyPropertyName string
			if len(toArraySplitted) == 1 {
				destinyPropertyName = originPropertyName
			} else {
				destinyPropertyName = toArraySplitted[1]
			}

			value, jsonMapResult := json_map.GetValue(jsonValueCurrent, originPropertyName, true)

			var arrayValue = [1]interface{}{value}
			jsonValueCurrent = json_map.CreateProperty(jsonMapResult, destinyPropertyName, arrayValue)
		}
	}

	if len(parse.ToMap) > 0 {
		for _, ToMap := range parse.ToMap {
			ToMapSplitted := strings.Split(ToMap, ":")

			originPropertyName := ToMapSplitted[0]
			var destinyPropertyName string
			if len(ToMapSplitted) == 1 {
				destinyPropertyName = originPropertyName
			} else {
				destinyPropertyName = ToMapSplitted[1]
			}

			value, jsonMapResult := json_map.GetValue(jsonValueCurrent, originPropertyName, true)

			var result map[string]interface{}
			json.Unmarshal([]byte(fmt.Sprint(value)), &result)

			jsonValueCurrent = json_map.CreateProperty(jsonMapResult, destinyPropertyName, result)
		}
	}

	return jsonValueCurrent
}

func formatDate(dateValue string, targetFormat string) (string, error) {
	t, err := time.Parse(time.RFC3339Nano, dateValue)
	if err != nil {
		return "", err
	}

	replacer := strings.NewReplacer(
		"yyyy", "2006",
		"yy", "06",
		"MM", "01",
		"dd", "02",
		"HH", "15",
		"hh", "03",
		"mm", "04",
		"ss", "05",
		"tt", "PM",
	)

	layout := replacer.Replace(targetFormat)

	return t.Format(layout), nil
}

func FormatValues(jsonMap map[string]interface{}, format *maps.FormatSettings) (map[string]interface{}, error) {
	jsonValueCurrent := jsonMap
	if len(format.Date) > 0 {
		for _, Date := range format.Date {
			DateSplitted := strings.Split(Date, ":")

			propertyName := DateSplitted[0]
			targetFormat := DateSplitted[1]

			dateValue, jsonMapResult := json_map.GetValue(jsonValueCurrent, propertyName, true)

			strDate, ok := dateValue.(string)
			if !ok {
				continue
			}

			formatedDate, err := formatDate(strDate, targetFormat)
			if err != nil {
				return jsonValueCurrent, err
			}

			jsonValueCurrent = json_map.CreateProperty(jsonMapResult, propertyName, formatedDate)
		}
	}

	return jsonValueCurrent, nil
}
