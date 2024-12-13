package handlers

import (
	"context"
	contexts "wrench/app/contexts"
	"wrench/app/json_map"
	"wrench/app/manifest/contract_settings/maps"
)

type HttpContractMapHandler struct {
	Next        Handler
	ContractMap *maps.ContractMapSetting
}

func (handler *HttpContractMapHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	if wrenchContext.HasError == false {
		currentBodyContext := bodyContext.BodyArray

		if len(handler.ContractMap.Sequence) > 0 {
			currentBodyContext = handler.doSequency(wrenchContext, bodyContext)
		} else {

			if handler.ContractMap.Rename != nil {
				currentBodyContext = json_map.RenameProperties(currentBodyContext, handler.ContractMap.Rename)
			}

			if handler.ContractMap.New != nil {
				currentBodyContext = json_map.CreatePropertiesInterpolationValue(
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
		}

		bodyContext.BodyArray = currentBodyContext
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handler *HttpContractMapHandler) doSequency(wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) []byte {
	currentBodyContext := bodyContext.BodyArray

	for _, action := range handler.ContractMap.Sequence {
		if action == "rename" {
			currentBodyContext = json_map.RenameProperties(currentBodyContext, handler.ContractMap.Rename)
		} else if action == "new" {
			currentBodyContext = json_map.CreatePropertiesInterpolationValue(
				currentBodyContext,
				handler.ContractMap.New,
				wrenchContext,
				bodyContext)
		} else if action == "remove" {
			currentBodyContext = json_map.RemoveProperties(currentBodyContext, handler.ContractMap.Remove)
		} else if action == "duplicate" {
			currentBodyContext = json_map.DuplicatePropertiesValue(currentBodyContext, handler.ContractMap.Duplicate)
		}
	}

	return currentBodyContext
}

func (handler *HttpContractMapHandler) SetNext(next Handler) {
	handler.Next = next
}