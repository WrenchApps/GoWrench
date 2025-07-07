package handlers

import (
	"context"
	"fmt"
	contexts "wrench/app/contexts"
	"wrench/app/json_map"
	"wrench/app/manifest/contract_settings/maps"
)

type HttpContractMapHandler struct {
	Next        Handler
	ContractMap *maps.ContractMapSetting
}

func (handler *HttpContractMapHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	if !wrenchContext.HasError {
		spanDisplay := fmt.Sprintf("contract.maps.%v", handler.ContractMap.Id)
		ctxSpan, span := wrenchContext.GetSpan2(ctx, spanDisplay)
		ctx = ctxSpan
		defer span.End()

		isArray := bodyContext.IsArray()

		if isArray {
			currentBodyContextArray := bodyContext.ParseBodyToMapObjectArray()
			lenArrayBody := len(currentBodyContextArray)
			if lenArrayBody > 0 {
				resultCurrentBodyContext := make([]map[string]interface{}, lenArrayBody)
				for i, currentBodyContext := range currentBodyContextArray {
					if len(handler.ContractMap.Sequence) > 0 {
						currentBodyContext = handler.doSequency(wrenchContext, bodyContext, currentBodyContext)
					} else {
						currentBodyContext = handler.doDefault(wrenchContext, bodyContext, currentBodyContext)
					}
					resultCurrentBodyContext[i] = currentBodyContext
				}
				bodyContext.SetArrayMapObject(resultCurrentBodyContext)
			}

		} else {
			currentBodyContext := bodyContext.ParseBodyToMapObject()

			if len(handler.ContractMap.Sequence) > 0 {
				currentBodyContext = handler.doSequency(wrenchContext, bodyContext, currentBodyContext)
			} else {
				currentBodyContext = handler.doDefault(wrenchContext, bodyContext, currentBodyContext)
			}
			bodyContext.SetMapObject(currentBodyContext)
		}
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handler *HttpContractMapHandler) doDefault(wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext, currentBodyContext map[string]interface{}) map[string]interface{} {

	if handler.ContractMap.Rename != nil {
		currentBodyContext = json_map.RenameProperties(currentBodyContext, handler.ContractMap.Rename)
	}

	if handler.ContractMap.New != nil {
		currentBodyContext = contexts.CreatePropertiesInterpolationValue(
			currentBodyContext,
			handler.ContractMap.New,
			wrenchContext,
			bodyContext)
	}

	if handler.ContractMap.Duplicate != nil {
		currentBodyContext = json_map.DuplicatePropertiesValue(currentBodyContext, handler.ContractMap.Duplicate)
	}

	if handler.ContractMap.Remove != nil {
		currentBodyContext = json_map.RemoveProperties(currentBodyContext, handler.ContractMap.Remove)
	}

	if handler.ContractMap.Parse != nil {
		currentBodyContext = contexts.ParseValues(currentBodyContext, handler.ContractMap.Parse)
	}

	if handler.ContractMap.Format != nil {
		currentBodyContext = contexts.FormatValues(currentBodyContext, handler.ContractMap.Format)
	}

	return currentBodyContext
}

func (handler *HttpContractMapHandler) doSequency(wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext, currentBodyContext map[string]interface{}) map[string]interface{} {
	for _, action := range handler.ContractMap.Sequence {
		if action == "rename" {
			currentBodyContext = json_map.RenameProperties(currentBodyContext, handler.ContractMap.Rename)
		} else if action == "new" {
			currentBodyContext = contexts.CreatePropertiesInterpolationValue(
				currentBodyContext,
				handler.ContractMap.New,
				wrenchContext,
				bodyContext)
		} else if action == "remove" {
			currentBodyContext = json_map.RemoveProperties(currentBodyContext, handler.ContractMap.Remove)
		} else if action == "duplicate" {
			currentBodyContext = json_map.DuplicatePropertiesValue(currentBodyContext, handler.ContractMap.Duplicate)
		} else if action == "parse" {
			currentBodyContext = contexts.ParseValues(currentBodyContext, handler.ContractMap.Parse)
		} else if action == "format" {
			currentBodyContext = contexts.FormatValues(currentBodyContext, handler.ContractMap.Format)
		}
	}

	return currentBodyContext
}

func (handler *HttpContractMapHandler) SetNext(next Handler) {
	handler.Next = next
}
