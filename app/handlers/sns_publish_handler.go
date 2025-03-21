package handlers

import (
	"context"
	"fmt"
	"log"
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
			groupId := getConculatedValue(settings.Fifo.GroupId, wrenchContext, bodyContext)
			if groupId != "" {
				publishInput.MessageGroupId = aws.String(groupId)
			}

			dedupId := getConculatedValue(settings.Fifo.DeduplicationId, wrenchContext, bodyContext)
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

func getConculatedValue(getConculatedValue string, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) string {
	var value string
	if len(getConculatedValue) > 0 {
		if contexts.IsCalculatedValue(getConculatedValue) {
			command := contexts.ReplaceCalculatedValue(getConculatedValue)
			if contexts.IsWrenchContextCommand(command) {
				value = contexts.GetValueWrenchContext(command, wrenchContext)
			} else {
				value = contexts.GetValueBodyContext(command, bodyContext)
			}
		} else {
			value = getConculatedValue
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
		filterValue = getConculatedValue(filterValue, wrenchContext, bodyContext)

		if filterKey != "" && filterValue != "" {
			mapAttributes[filterKey] = types.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String(filterValue)}
		}
	}

	return mapAttributes
}
