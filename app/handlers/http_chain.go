package handlers

import (
	action_settings "wrench/app/manifest/action_settings"
	settings "wrench/app/manifest/application_settings"
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

	for _, action := range settings.Actions {

		var firstHandler = new(HttpFirstHandler)
		var currentHandler Handler
		currentHandler = firstHandler

		if action.Trigger != nil && action.Trigger.Before != nil {
			httpContractMapHandler := new(HttpContractMapHandler)

			contractMapId := action.Trigger.Before.ContractMapId
			httpContractMapHandler.ContractMap = settings.Contract.GetContractById(contractMapId)

			currentHandler.SetNext(httpContractMapHandler)
			currentHandler = httpContractMapHandler
		}

		if action.Type == action_settings.ActionTypeHttpRequest {
			httpRequestHadler := new(HttpRequestClientHandler)
			httpRequestHadler.ActionSettings = &action
			currentHandler.SetNext(httpRequestHadler)
			currentHandler = httpRequestHadler
		}

		if action.Type == action_settings.ActionTypeHttpRequestMock {
			httpRequestMockHadler := new(HttpRequestClientMockHandler)
			httpRequestMockHadler.ActionSettings = &action
			currentHandler.SetNext(httpRequestMockHadler)
			currentHandler = httpRequestMockHadler
		}

		if action.Type == action_settings.ActionTypeSnsPublish {
			snsPublishHandler := new(SnsPublishHandler)
			snsPublishHandler.ActionSettings = &action
			currentHandler.SetNext(snsPublishHandler)
			currentHandler = snsPublishHandler
		}

		if action.Type == action_settings.ActionTypeFileReader {
			fileReaderHandler := new(FileReaderHandler)
			fileReaderHandler.ActionSettings = &action
			currentHandler.SetNext(fileReaderHandler)
			currentHandler = fileReaderHandler
		}

		if action.Trigger != nil && action.Trigger.After != nil {
			httpContractMapHandler := new(HttpContractMapHandler)

			contractMapId := action.Trigger.After.ContractMapId
			httpContractMapHandler.ContractMap = settings.Contract.GetContractById(contractMapId)

			currentHandler.SetNext(httpContractMapHandler)
			currentHandler = httpContractMapHandler
		}

		currentHandler.SetNext(new(HttpLastHandler))

		chain.MapHandle[action.Id] = firstHandler
	}
}

func (chain *Chain) GetByActionId(actionId string) Handler {
	return chain.MapHandle[actionId]
}
