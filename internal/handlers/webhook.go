package handlers

import (
	"encoding/json"
	"net/http"
	"bot-gateway/internal/models"
	"bot-gateway/internal/services"
)

func HandleWebhook(botService *services.BotService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		var webhook models.GreenApiWebhook
		err := json.NewDecoder(r.Body).Decode(&webhook)
		if err != nil {
			return
		}
		defer r.Body.Close()

		if webhook.TypeWebhook != "incomingMessageReceived" {
			return
		}

		msgType := webhook.MessageData.TypeMessage
		if msgType != "textMessage" && msgType != "audioMessage" {
			return
		}

		go botService.ProcessUpdate(webhook)
	}
}