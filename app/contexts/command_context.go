package contexts

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	auth_jwt "wrench/app/auth/jwt"
	"wrench/app/converts"
	"wrench/app/json_map"
	settings "wrench/app/manifest/action_settings"
	"wrench/app/manifest/action_settings/func_settings"
	"wrench/app/manifest/contract_settings/maps"

	"github.com/google/uuid"
)

const prefixWrenchContextRequest = "wrenchContext.request."
const prefixWrenchContextRequestUri = "wrenchContext.request.uri"
const prefixWrenchContextRequestUriParams = "wrenchContext.request.uri.params."
const prefixWrenchContextRequestTokenClaims = "wrenchContext.request.token.claims."
const prefixWrenchContextRequestHeaders = "wrenchContext.request.headers."

const prefixWrenchContextResponse = "wrenchContext.response."
const prefixWrenchContextResponseHttpStatusCode = "wrenchContext.response.httpStatusCode"
const prefixWrenchContextResponseHttpStatusCodeFamily = "wrenchContext.response.httpStatusCodeFamily"

const prefixBodyContext = "bodyContext."
const prefixBodyContextPreserved = "bodyContext.actions."
const prefixFunc = "func."

var layoutsDates = []string{
	// --- 1. Standard ISO / RFC / API Formats (Highly Recommended) ---
	time.RFC3339,                 // "2006-01-02T15:04:05Z07:00"
	time.RFC3339Nano,             // "2006-01-02T15:04:05.999999999Z07:00" (Catches almost all standard API payloads)
	"2006-01-02 15:04:05.999999", // MySQL / PostgreSQL timestamp with microseconds
	"2006-01-02 15:04:05",        // Standard SQL Date Time

	// --- 2. Slanted / European / Latin American Formats (DD/MM) ---
	"02/01/2006 15:04:05.999", // 31/12/2026 14:30:15.123
	"02/01/2006 15:04:05",     // 31/12/2026 14:30:15
	"02/01/2006 15:04",        // 31/12/2026 14:30
	"02/01/2006",              // 31/12/2026 (Pure Date)

	// --- 3. Dashed / European Formats (DD-MM) ---
	"02-01-2006 15:04:05.999", // 31-12-2026 14:30:15.123
	"02-01-2006 15:04:05",     // 31-12-2026 14:30:15
	"02-01-2006",              // 31-12-2026

	// --- 4. Pure Dates (No Time Components) ---
	"2006-01-02", // 2026-12-31 (Standard ISO Date)
	"20060102",   // 20261231 (Compact Date)

	// --- 5. US Standard Formats (MM/DD) *Keep at the bottom due to ambiguity* ---
	"01/02/2006 15:04:05.999", // 12/31/2026 14:30:15.123
	"01/02/2006 15:04:05",     // 12/31/2026 14:30:15
	"01/02/2006",              // 12/31/2026
}

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
	return strings.HasPrefix(command, prefixWrenchContextRequest) || strings.HasPrefix(command, prefixWrenchContextResponse)
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

func GetTokenClaims(wrenchContext *WrenchContext, claimName string) string {
	tokenString := wrenchContext.Request.Header.Get("Authorization")

	if len(tokenString) == 0 {
		return ""
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	tokenSplitted := strings.Split(tokenString, ".")
	tokenPayload := tokenSplitted[1]

	tokenPayloadMap := auth_jwt.ConvertJwtPayloadBase64ToJwtPaylodData(tokenPayload)
	claimTokenValue, _ := tokenPayloadMap[claimName].(string)

	return claimTokenValue
}

func GetValueWrenchContext(command string, wrenchContext *WrenchContext, bodyContext *BodyContext) string {

	if IsCalculatedValue(command) {
		command = ReplaceCalculatedValue(command)
	}

	if strings.HasPrefix(command, prefixWrenchContextRequest) {
		return getValueWrenchContextPrefixWrenchContextRequest(command, wrenchContext)
	} else if strings.HasPrefix(command, prefixWrenchContextResponse) {
		return getValueWrenchContextPrefixWrenchContextResponse(command, wrenchContext, bodyContext)
	}

	return ""
}

func getValueWrenchContextPrefixWrenchContextRequest(command string, wrenchContext *WrenchContext) string {
	if strings.HasPrefix(command, prefixWrenchContextRequestHeaders) {
		headerName := strings.ReplaceAll(command, prefixWrenchContextRequestHeaders, "")
		return wrenchContext.Request.Header.Get(headerName)
	}

	if strings.HasPrefix(command, prefixWrenchContextRequestUriParams) {
		parameterName := strings.ReplaceAll(command, prefixWrenchContextRequestUriParams, "")
		return GetRequestUriParams(wrenchContext, parameterName)
	}

	if strings.HasPrefix(command, prefixWrenchContextRequestTokenClaims) {
		parameterName := strings.ReplaceAll(command, prefixWrenchContextRequestTokenClaims, "")
		return GetTokenClaims(wrenchContext, parameterName)
	}

	if strings.HasPrefix(command, prefixWrenchContextRequestUri) {
		return wrenchContext.Request.RequestURI
	}
	return ""
}

func getValueWrenchContextPrefixWrenchContextResponse(command string, wrenchContext *WrenchContext, bodyContext *BodyContext) string {

	if strings.HasPrefix(command, prefixWrenchContextResponseHttpStatusCode) {
		return fmt.Sprint(bodyContext.HttpStatusCode)
	}

	if strings.HasPrefix(command, prefixWrenchContextResponseHttpStatusCodeFamily) {
		statusCode := bodyContext.HttpStatusCode
		if statusCode < 100 || statusCode > 599 {
			return "0"
		}
		return fmt.Sprint((statusCode / 100) * 100)
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

	if strings.Contains(command, "???") {
		commandNullOrEmptySpllited := strings.Split(command, "???")
		for _, commandNullOrEmpty := range commandNullOrEmptySpllited {

			value := GetCalculatedValue(commandNullOrEmpty, wrenchContext, bodyContext, action)
			if value != nil && len(fmt.Sprint(value)) > 0 {
				return value
			}
		}
	}

	if strings.Contains(command, "??") {
		commandNullOrEmptySpllited := strings.Split(command, "??")
		for _, commandNull := range commandNullOrEmptySpllited {

			value := GetCalculatedValue(commandNull, wrenchContext, bodyContext, action)
			if value != nil {
				return value
			}
		}
	}

	if IsCalculatedValue(command) {
		command = ReplaceCalculatedValue(command)
		if IsBodyContextCommand(command) {
			return GetValueBodyContext(command, bodyContext)
		} else if IsWrenchContextCommand(command) {
			return GetValueWrenchContext(command, wrenchContext, bodyContext)
		} else if IsFunc(command) {
			value, _ := GetFuncValue(func_settings.FuncGeneralType(command), wrenchContext, bodyContext, action)
			return value
		} else if command == "uuid" {
			return uuid.New().String()
		} else if strings.HasPrefix(command, "time") {
			timeFormat := strings.ReplaceAll(command, "time ", "")
			timeNow := time.Now()

			if len(timeFormat) > 0 {
				return timeNow.Format(timeFormat)
			} else {
				return timeNow.String()
			}
		} else {
			command = fmt.Sprintf("%v%v", prefixBodyContext, command)
			return GetValueBodyContext(command, bodyContext)
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
			bodyPreserved, _ := bodyContext.GetBodyPreserved(actionId)
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

func CreatePropertiesInterpolationValue(jsonMap map[string]interface{}, propertiesValues []string, wrenchContext *WrenchContext, bodyContext *BodyContext, action *settings.ActionSettings) map[string]interface{} {
	jsonValueCurrent := jsonMap
	for _, propertyValue := range propertiesValues {
		propertyValueSplitted := strings.Split(propertyValue, ":")
		propertyName := propertyValueSplitted[0]
		valueArray := propertyValueSplitted[1:]
		value := strings.Join(valueArray, ":")
		jsonValueCurrent = CreatePropertyInterpolationValue(jsonValueCurrent, propertyName, value, wrenchContext, bodyContext, action)
	}
	return jsonValueCurrent
}

func CreatePropertyInterpolationValue(jsonMap map[string]interface{}, propertyName string, value interface{}, wrenchContext *WrenchContext, bodyContext *BodyContext, action *settings.ActionSettings) map[string]interface{} {
	valueResult := value
	valueString := fmt.Sprint(valueResult)
	valueResult = GetCalculatedValue(valueString, wrenchContext, bodyContext, action)

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

	if len(parse.ToDate) > 0 {
		for _, toDate := range parse.ToDate {
			toDateSplitted := strings.Split(toDate, ":")

			originPropertyName := toDateSplitted[0]
			var destinyPropertyName string
			if len(toDateSplitted) == 1 {
				destinyPropertyName = originPropertyName
			} else {
				destinyPropertyName = toDateSplitted[1]
			}

			value, jsonMapResult := json_map.GetValue(jsonValueCurrent, originPropertyName, true)
			result, _ := parseStringToDateUsingLayouts(fmt.Sprint(value))
			jsonValueCurrent = json_map.CreateProperty(jsonMapResult, destinyPropertyName, result.Format(time.RFC3339Nano))
		}
	}

	return jsonValueCurrent
}

func parseStringToDateUsingLayouts(dateString string) (time.Time, error) {

	var sanitized = dateString
	if !strings.Contains(dateString, "Z") {
		sanitized = strings.Replace(fmt.Sprint(dateString), "T", " ", 1)
	}

	for _, layout := range layoutsDates {

		if t, err := time.Parse(layout, sanitized); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateString)
}

func formatDate(dateValue string, targetFormat string) (string, error) {

	dateParsed, err := parseStringToDateUsingLayouts(dateValue)

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
		"zz", "-07:00",
	)

	layout := replacer.Replace(targetFormat)

	return dateParsed.Format(layout), nil
}

func FormatValues(jsonMap map[string]interface{}, format *maps.FormatSettings) (map[string]interface{}, error) {
	jsonValueCurrent := jsonMap
	if len(format.Date) > 0 {
		for _, date := range format.Date {
			dateSplitted := strings.Split(date, ":")

			propertyName := dateSplitted[0]
			targetFormat := dateSplitted[1:]

			dateValue, jsonMapResult := json_map.GetValue(jsonValueCurrent, propertyName, true)

			strDate, ok := dateValue.(string)
			if !ok {
				continue
			}

			formatedDate, err := formatDate(strDate, strings.Join(targetFormat, ":"))
			if err != nil {
				return jsonValueCurrent, err
			}

			jsonValueCurrent = json_map.CreateProperty(jsonMapResult, propertyName, formatedDate)
		}
	}

	return jsonValueCurrent, nil
}

func ConcatenatePropertiesOrValues(jsonMap map[string]interface{}, concatenatePropertiesValues []string, wrenchContext *WrenchContext, bodyContext *BodyContext) map[string]interface{} {

	jsonValueCurrent := jsonMap
	for _, propertyValue := range concatenatePropertiesValues {
		propertyValueSplitted := strings.Split(propertyValue, ":")

		propertyName := propertyValueSplitted[0]
		valueToConcatenate := propertyValueSplitted[1:]
		currentValue, jsonValueCurrent := json_map.GetValue(jsonValueCurrent, propertyName, true)
		concatenatedValue := GetCalculatedValue(strings.Join(valueToConcatenate, ":"), wrenchContext, bodyContext, nil)
		finalValue := fmt.Sprintf("%v%v", currentValue, concatenatedValue)

		jsonValueCurrent = json_map.SetValue(jsonValueCurrent, propertyName, finalValue)
	}
	return jsonValueCurrent
}

func ApplyMathOperations(jsonMap map[string]interface{}, mapSettings *maps.MathSettings) (map[string]interface{}, error) {
	if mapSettings == nil {
		return jsonMap, nil
	}

	for _, expr := range *mapSettings {

		operatorIndex := strings.IndexAny(expr, "+-*/")

		operator := expr[operatorIndex : operatorIndex+1]
		fieldName := strings.TrimSpace(expr[:operatorIndex])
		rawFactor := strings.TrimSpace(expr[operatorIndex+1:])

		factor, _ := strconv.ParseFloat(rawFactor, 64)
		rawValue, jsonMapResult := json_map.GetValue(jsonMap, fieldName, true)
		if rawValue == nil {
			return nil, fmt.Errorf("field '%s' not found", fieldName)
		}

		numericValue, err := converts.ConvertToFloat(rawValue)

		if err != nil {
			return nil, err
		}

		var result float64
		switch operator {
		case "*":
			result = numericValue * factor
		case "/":
			result = numericValue / factor
		case "+":
			result = numericValue + factor
		case "-":
			result = numericValue - factor
		default:
			return nil, fmt.Errorf("invalid operation in math")
		}

		json_map.SetValue(jsonMapResult, fieldName, result)
	}

	return jsonMap, nil
}
