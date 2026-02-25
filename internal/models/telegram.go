package models

type GreenApiWebhook struct {
	TypeWebhook string `json:"typeWebhook"`
	SenderData  struct {
		ChatId string `json:"chatId"`
	} `json:"senderData"`
	MessageData struct {
		TypeMessage     string `json:"typeMessage"`
		TextMessageData struct {
			TextMessage string `json:"textMessage"`
		} `json:"textMessageData"`
	} `json:"messageData"`
}

type AIRequest struct {
	ChatID   string `json:"chat_id"` 
	UserText string `json:"user_text"`
	VoiceURL string `json:"voice_url"`
}

type AIResponse struct {
	Reply string `json:"reply"` 
}

type GreenApiSendRequest struct {
	ChatId  string `json:"chatId"`
	Message string `json:"message"`
}