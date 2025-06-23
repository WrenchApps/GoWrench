package handlers

import settings "wrench/app/manifest/action_settings"

type KafkaProducerHandler struct {
	Next           Handler
	ActionSettings *settings.ActionSettings
}
