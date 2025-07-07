package handlers

import (
	action_settings "wrench/app/manifest/action_settings"
	settings "wrench/app/manifest/application_settings"
	"wrench/app/manifest_cross_funcs"
)

var ChainStatic *Chain = new(Chain)

type Chain struct {
	MapHandle map[string]Handler
}

func (chain *Chain) GetStatic() *Chain {
	return ChainStatic
}

func (chain *Chain) BuildChain(settings *settings.ApplicationSettings) {

	chain.MapHandle = make(map[string]Handler)
	if settings.Api == nil || settings.Api.Endpoints == nil {
		return
	}

	for _, endpoint := range settings.Api.Endpoints {
		var firstHandler = new(HttpFirstHandler)
		var currentHandler Handler
		currentHandler = firstHandler

		if len(endpoint.IdempId) > 0 {

			idempHandler := new(IdempHandler)
			idempHandler.EndpointSettings = &endpoint
			idempHandler.IdempSettings, _ = manifest_cross_funcs.GetIdempSettingById(endpoint.IdempId)
			idempHandler.RedisSettings, _ = manifest_cross_funcs.GetConnectionRedisSettingById(idempHandler.IdempSettings.RedisConnectionId)

			currentHandler.SetNext(idempHandler)
			currentHandler = idempHandler
		}

		if len(endpoint.ActionID) > 0 {
			action, _ := settings.GetActionById(endpoint.ActionID)
			currentHandler = buildChainToAction(currentHandler, settings, action)
		} else {
			for _, actionId := range endpoint.FlowActionID {
				action, _ := settings.GetActionById(actionId)
				currentHandler = buildChainToAction(currentHandler, settings, action)
			}
		}

		currentHandler.SetNext(new(HttpLastHandler))
		chain.MapHandle[endpoint.Route] = firstHandler
	}
}

func buildChainToAction(currentHandler Handler, settings *settings.ApplicationSettings, action *action_settings.ActionSettings) Handler {
	if action.Trigger != nil && action.Trigger.Before != nil {
		httpContractMapHandler := new(HttpContractMapHandler)

		contractMapId := action.Trigger.Before.ContractMapId
		httpContractMapHandler.ContractMap = settings.Contract.GetContractById(contractMapId)

		currentHandler.SetNext(httpContractMapHandler)
		currentHandler = httpContractMapHandler
	}

	if action.Type == action_settings.ActionTypeHttpRequest {
		httpRequestHadler := new(HttpRequestClientHandler)
		httpRequestHadler.ActionSettings = action
		currentHandler.SetNext(httpRequestHadler)
		currentHandler = httpRequestHadler
	}

	if action.Type == action_settings.ActionTypeHttpRequestMock {
		httpRequestMockHadler := new(HttpRequestClientMockHandler)
		httpRequestMockHadler.ActionSettings = action
		currentHandler.SetNext(httpRequestMockHadler)
		currentHandler = httpRequestMockHadler
	}

	if action.Type == action_settings.ActionTypeSnsPublish {
		snsPublishHandler := new(SnsPublishHandler)
		snsPublishHandler.ActionSettings = action
		currentHandler.SetNext(snsPublishHandler)
		currentHandler = snsPublishHandler
	}

	if action.Type == action_settings.ActionTypeFileReader {
		fileReaderHandler := new(FileReaderHandler)
		fileReaderHandler.ActionSettings = action
		currentHandler.SetNext(fileReaderHandler)
		currentHandler = fileReaderHandler
	}

	if action.Type == action_settings.ActionTypeNatsPublish {
		httpNatsPublishHandler := new(NatsPublishHandler)
		httpNatsPublishHandler.ActionSettings = action
		currentHandler.SetNext(httpNatsPublishHandler)
		currentHandler = httpNatsPublishHandler
	}

	if action.Type == action_settings.ActionTypeFuncHash {
		funcHashHandler := new(FuncHashHandler)
		funcHashHandler.ActionSettings = action
		currentHandler.SetNext(funcHashHandler)
		currentHandler = funcHashHandler
	}

	if action.Type == action_settings.ActionTypeFuncVarContext {
		funcVarHandler := new(FuncVarContextHandler)
		funcVarHandler.ActionSettings = action
		currentHandler.SetNext(funcVarHandler)
		currentHandler = funcVarHandler
	}

	if action.Type == action_settings.ActionTypeFuncStringConcatenate {
		funcStringConcateHandler := new(FuncStringConcatenateHandler)
		funcStringConcateHandler.ActionSettings = action
		currentHandler.SetNext(funcStringConcateHandler)
		currentHandler = funcStringConcateHandler
	}

	if action.Type == action_settings.ActionTypeFuncGeneral {
		funcGeneralHandler := new(FuncGeneralHandler)
		funcGeneralHandler.ActionSettings = action
		currentHandler.SetNext(funcGeneralHandler)
		currentHandler = funcGeneralHandler
	}

	if action.Type == action_settings.ActionTypeKafkaProducer {
		kafkaProducerHandler := new(KafkaProducerHandler)
		kafkaProducerHandler.ActionSettings = action
		currentHandler.SetNext(kafkaProducerHandler)
		currentHandler = kafkaProducerHandler
	}

	if action.Trigger != nil && action.Trigger.After != nil {
		httpContractMapHandler := new(HttpContractMapHandler)

		contractMapId := action.Trigger.After.ContractMapId
		httpContractMapHandler.ContractMap = settings.Contract.GetContractById(contractMapId)

		currentHandler.SetNext(httpContractMapHandler)
		currentHandler = httpContractMapHandler
	}

	return currentHandler
}

func (chain *Chain) GetHandler(key string) Handler {
	return chain.MapHandle[key]
}
