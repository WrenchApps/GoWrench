package contexts

import (
	"encoding/json"
	"strings"
	"wrench/app/manifest/action_settings"
)

type BodyContext struct {
	CurrentBodyByteArray []byte
	BodyPreserved        map[string][]byte
	HttpStatusCode       int
	ContentType          string
	Headers              map[string]string
}

func (bodyContext *BodyContext) SetBodyPreserved(id string, body []byte) {
	if bodyContext.BodyPreserved == nil {
		bodyContext.BodyPreserved = make(map[string][]byte)
	}

	bodyContext.BodyPreserved[id] = body
}

func (bodyContext *BodyContext) SetBody(body []byte) {
	bodyContext.CurrentBodyByteArray = body
}

func (bodyContext *BodyContext) SetBodyAction(settings *action_settings.ActionSettings, body []byte) {
	if settings.ShouldPreserveBody() {
		bodyContext.SetBodyPreserved(settings.Id, body)
	} else {
		bodyContext.SetBody(body)
	}
}

func (bodyContext *BodyContext) GetBodyPreserved(id string) []byte {
	if bodyContext.BodyPreserved == nil {
		return nil
	}

	return bodyContext.BodyPreserved[id]
}

func (bodyContext *BodyContext) IsArray() bool {
	bodyText := string(bodyContext.CurrentBodyByteArray)
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
	jsonErr := json.Unmarshal(bodyContext.CurrentBodyByteArray, &jsonMap)

	if jsonErr != nil {
		return nil
	}
	return jsonMap
}

func (bodyContext *BodyContext) ParseBodyToMapObjectPreserved(actionId string) map[string]interface{} {
	var jsonMap map[string]interface{}
	bodyBytePreserved := bodyContext.GetBodyPreserved(actionId)
	jsonErr := json.Unmarshal(bodyBytePreserved, &jsonMap)

	if jsonErr != nil {
		return nil
	}
	return jsonMap
}

func (bodyContext *BodyContext) ParseBodyToMapObjectArray() []map[string]interface{} {
	var jsonMap []map[string]interface{}
	jsonErr := json.Unmarshal(bodyContext.CurrentBodyByteArray, &jsonMap)

	if jsonErr != nil {
		return nil
	}
	return jsonMap
}

func (bodyContext *BodyContext) SetMapObject(jsonMap map[string]interface{}) {
	jsonArray, _ := json.Marshal(jsonMap)
	bodyContext.CurrentBodyByteArray = jsonArray
}

func (bodyContext *BodyContext) ConvertMapToByteArray(jsonMap map[string]interface{}) ([]byte, error) {
	jsonArray, err := json.Marshal(jsonMap)
	if err != nil {
		return nil, err
	}

	return jsonArray, nil
}

func (bodyContext *BodyContext) SetArrayMapObject(arrayJsonMap []map[string]interface{}) {
	jsonArray, _ := json.Marshal(arrayJsonMap)
	bodyContext.CurrentBodyByteArray = jsonArray
}

func (bodyContext *BodyContext) GetBodyString() string {
	return string(bodyContext.CurrentBodyByteArray)
}

func (bodyContext *BodyContext) GetBody(settings *action_settings.ActionSettings) []byte {
	shouldUse, bodyRef := settings.ShouldUseBodyRef()

	if shouldUse {
		bodyRef = ReplaceCalculatedValue(bodyRef)
		bodyRef = ReplacePrefixBodyContextPreserved(bodyRef)
		return bodyContext.GetBodyPreserved(bodyRef)
	} else {
		return bodyContext.CurrentBodyByteArray
	}
}

func (bodyContext *BodyContext) GetCurrentBody() []byte {
	return bodyContext.CurrentBodyByteArray
}
