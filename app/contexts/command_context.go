package contexts

import (
	"strings"
)

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

func IsWrenchContextCommand(command string) bool {
	return strings.HasPrefix(command, prefixWrenchContextRequestHeaders)
}

func IsBodyContextCommand(command string) bool {
	return strings.HasPrefix(command, prefixBodyContext)
}

func IsFunc(command string) bool {
	return strings.HasPrefix(command, prefixFunc)
}

func GetValueWrenchContext(command string, wrenchContext *WrenchContext) string {

	if IsCalculatedValue(command) {
		command = ReplaceCalculatedValue(command)
	}

	if strings.HasPrefix(command, prefixWrenchContextRequestHeaders) {
		headerName := strings.ReplaceAll(command, prefixWrenchContextRequestHeaders, "")
		return wrenchContext.Request.Header.Get(headerName)
	}

	return ""
}

func ReplacePrefixBodyContext(command string) string {
	if strings.HasPrefix(command, prefixBodyContext) {
		command = strings.ReplaceAll(command, prefixBodyContext, "")
	}
	return command
}

func GetCalculatedValue(command string, wrenchContext *WrenchContext, bodyContext *BodyContext) string {
	command = ReplaceCalculatedValue(command)

	if IsBodyContextCommand(command) {
		return GetValueBodyContext(command, bodyContext)
	} else if IsWrenchContextCommand(command) {
		return GetValueWrenchContext(command, wrenchContext)
	} else if IsFunc(command) {
		return GetFuncValue(FuncType(command))
	} else {
		return ""
	}
}

func GetValueBodyContext(command string, bodyContext *BodyContext) string {

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
			return getBodyValue(jsonMap, propertyName)
		}

	} else if strings.HasPrefix(command, prefixBodyContext) {
		propertyName := strings.ReplaceAll(command, prefixBodyContext, "")
		jsonMap := bodyContext.ParseBodyToMapObject()
		return getBodyValue(jsonMap, propertyName)
	}

	return ""
}

func getBodyValue(jsonMap map[string]interface{}, propertyName string) string {
	value := ""

	var jsonMapCurrent map[string]interface{}
	jsonMapCurrent = jsonMap
	propertyNameSplitted := strings.Split(propertyName, ".")

	for _, property := range propertyNameSplitted {
		valueTemp, ok := jsonMapCurrent[property].(map[string]interface{})
		if ok {
			jsonMapCurrent = valueTemp
			continue
		}

		valueTempString, ok := jsonMapCurrent[property].(string)
		if ok {
			value = valueTempString
			break
		}
	}
	return value
}

func GetCalculatedMap(mapConfigured map[string]string, wrenchContext *WrenchContext, bodyContext *BodyContext) map[string]interface{} {
	if mapConfigured == nil {
		return nil
	}
	mapResult := make(map[string]interface{})

	for key, value := range mapConfigured {
		var finalValue string
		if IsCalculatedValue(value) {
			finalValue = GetCalculatedValue(value, wrenchContext, bodyContext)
		} else {
			finalValue = value
		}

		mapResult[key] = finalValue
	}

	return mapResult
}
