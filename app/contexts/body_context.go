package contexts

import (
	"encoding/json"
	"strings"
)

type BodyContext struct {
	BodyByteArray  []byte
	HttpStatusCode int
	ContentType    string
	Headers        map[string]string
}

func (bodyContext *BodyContext) IsArray() bool {
	bodyText := string(bodyContext.BodyByteArray)
	return strings.HasPrefix(bodyText, "[") && strings.HasSuffix(bodyText, "]")
}

func (bodyContext *BodyContext) SetHeaders(headers map[string]string) {
	if headers != nil {
		if bodyContext.Headers == nil {
			bodyContext.Headers = make(map[string]string)
		}

		for key, value := range headers {
			bodyContext.Headers[key] = value
		}
	}
}

func (bodyContext *BodyContext) SetHeader(key string, value string) {
	if len(key) > 0 {
		if bodyContext.Headers == nil {
			bodyContext.Headers = make(map[string]string)
		}

		bodyContext.Headers[key] = value
	}
}

func (bodyContext *BodyContext) ParseBodyToMapObject() map[string]interface{} {
	var jsonMap map[string]interface{}
	jsonErr := json.Unmarshal(bodyContext.BodyByteArray, &jsonMap)

	if jsonErr != nil {
		return nil
	}
	return jsonMap
}

func (bodyContext *BodyContext) ParseBodyToMapObjectArray() []map[string]interface{} {
	var jsonMap []map[string]interface{}
	jsonErr := json.Unmarshal(bodyContext.BodyByteArray, &jsonMap)

	if jsonErr != nil {
		return nil
	}
	return jsonMap
}

func (bodyContext *BodyContext) SetMapObject(jsonMap map[string]interface{}) {
	if jsonMap == nil {
		return
	}
	jsonArray, _ := json.Marshal(jsonMap)
	bodyContext.BodyByteArray = jsonArray
}

func (bodyContext *BodyContext) SetArrayMapObject(arrayJsonMap []map[string]interface{}) {
	if arrayJsonMap == nil {
		return
	}
	jsonArray, _ := json.Marshal(arrayJsonMap)
	bodyContext.BodyByteArray = jsonArray
}
