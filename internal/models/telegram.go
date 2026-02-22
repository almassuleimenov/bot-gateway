package models

type Update struct {
	UpdateID int  `json:"update_id"`
	Message *Message `json:"message"`
}


type Message struct {
	Text string `json:"text"`
	Chat Chat `json:"chat"`
}

type Chat struct{
	ID int `json:"id"`
}
type SendMessagePayload struct{
	ChatID int `json:"chat_id"`
	Text string `json:"text"`
}

type AIRequest struct{
	ChatID int64 `json:"chat_id"`
	UserText string `json:"user_text"`
}

type AIResponse struct{
	Reply string `json:"reply"`
}