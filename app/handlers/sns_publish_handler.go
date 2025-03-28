package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
	"wrench/app/manifest/action_settings/sns_settings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

type SnsPublishHandler struct {
	Next           Handler
	ActionSettings *settings.ActionSettings
}

type SnsActions struct {
	SnsClient *sns.Client
}

func (snsActions *SnsActions) Load() {
	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println("Couldn't load default configuration. Have you set up your AWS account?")
		fmt.Println(err)
		return
	}
	snsActions.SnsClient = sns.NewFromConfig(sdkConfig)
}

func (handler *SnsPublishHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {
	if !wrenchContext.HasError {
		settings := handler.ActionSettings.SNS
		message := bodyContext.GetBodyString()
		actor := new(SnsActions)
		actor.Load()

		publishInput := sns.PublishInput{TopicArn: aws.String(settings.TopicArn), Message: aws.String(message)}

		if settings.IsFifo() {
			groupId := getCalculatedValue(settings.Fifo.GroupId, wrenchContext, bodyContext)
			if groupId != "" {
				publishInput.MessageGroupId = aws.String(groupId)
			}

			dedupId := getCalculatedValue(settings.Fifo.DeduplicationId, wrenchContext, bodyContext)
			if dedupId != "" {
				publishInput.MessageDeduplicationId = aws.String(dedupId)
			}
		}

		if len(settings.Filters) > 0 {
			publishInput.MessageAttributes = getSnsFilter(settings, wrenchContext, bodyContext)
		}

		_, err := actor.SnsClient.Publish(ctx, &publishInput)
		if err != nil {
			msg := fmt.Sprintf("Couldn't publish message to topic %v. Here's why: %v", settings.TopicArn, err)
			log.Print(msg)
			bodyContext.HttpStatusCode = 500
			bodyContext.BodyByteArray = []byte(msg)
			bodyContext.ContentType = "text/plain"
			wrenchContext.SetHasError()
		} else {
			bodyContext.HttpStatusCode = 202
			bodyContext.BodyByteArray = []byte("{ 'success': 'true' }")
		}
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handler *SnsPublishHandler) SetNext(next Handler) {
	handler.Next = next
}

func getCalculatedValue(getCalculatedValue string, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) string {
	var value string
	if len(getCalculatedValue) > 0 {
		if contexts.IsCalculatedValue(getCalculatedValue) {
			command := contexts.ReplaceCalculatedValue(getCalculatedValue)
			if contexts.IsWrenchContextCommand(command) {
				value = contexts.GetValueWrenchContext(command, wrenchContext)
			} else {
				if contexts.IsBodyContextCommand(command) {
					propertyName := contexts.ReplacePrefixBodyContext(command)
					jsonMap := bodyContext.ParseBodyToMapObject()
					if jsonMap != nil {
						jsonMapCurrent := jsonMap
						propertyNameSplitted := strings.Split(propertyName, ".")
						
						for i, property := range propertyNameSplitted {
							if i == len(propertyNameSplitted)-1 {
								if val, ok := jsonMapCurrent[property]; ok {
									switch v := val.(type) {
									case string:
										value = v
									case int:
										value = strconv.Itoa(v)
									default:
										if jsonBytes, err := json.Marshal(v); err == nil {
											value = string(jsonBytes)
										}
									}
								}
								break
							} else {
								if mapVal, ok := jsonMapCurrent[property].(map[string]interface{}); ok {
									jsonMapCurrent = mapVal
								} else {
									break
								}
							}
						}
					}
				} else {
					value = contexts.GetValueBodyContext(command, bodyContext)
				}
			}
		} else {
			value = getCalculatedValue
		}
	}

	return value
}

func getSnsFilter(snsSettings *sns_settings.SnsSettings, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) map[string]types.MessageAttributeValue {
	mapAttributes := map[string]types.MessageAttributeValue{}
	for _, filter := range snsSettings.Filters {
		filterSplitted := strings.Split(filter, ":")
		filterKey := filterSplitted[0]
		filterValue := filterSplitted[1]
		filterValue = getCalculatedValue(filterValue, wrenchContext, bodyContext)

		if filterKey != "" && filterValue != "" {
			if intValue, err := strconv.Atoi(filterValue); err == nil {
				mapAttributes[filterKey] = types.MessageAttributeValue{
					DataType:    aws.String("Number"),
					StringValue: aws.String(strconv.Itoa(intValue)),
				}
			} else {
				mapAttributes[filterKey] = types.MessageAttributeValue{
					DataType:    aws.String("String"),
					StringValue: aws.String(filterValue),
				}
			}
		}
	}

	return mapAttributes
}
