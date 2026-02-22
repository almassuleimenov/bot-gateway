package handlers

import (
	"encoding/json"
	"net/http"
	"bot-gateway/internal/services"
	"bot-gateway/internal/models"
)

func HandleWebhook(botservice *services.BotService) http.HandlerFunc{
	return func(w http.ResponseWriter , r *http.Request){

		var update models.Update

		

		err := json.NewDecoder(r.Body).Decode(&update)

		if err != nil{
			http.Error(w,"Invalid JSON",http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		go botservice.ProcessUpdate(update);

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}