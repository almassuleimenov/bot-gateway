package services

import (
	"bot-gateway/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)


type BotService struct{
	Token string
}

func NewBotService(token string) *BotService {
	return &BotService{
		Token : token,
	}	
}

func (s *BotService) ProcessUpdate(update models.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	if update.Message.Text == "/start" {
		s.sendMessage(update.Message.Chat.ID, "Привет! Я ИИ-помощник архитектурного бюро. Спрашивай про наши проекты!")
		return
	}

	chatID := update.Message.Chat.ID
	userText := update.Message.Text

	aiReq := models.AIRequest{
		ChatID:   int64(chatID),
		UserText: userText,
	}
	jsonData, _ := json.Marshal(aiReq)

	resp, err := http.Post("http://127.0.0.1:8000/generate-answer", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Питон оффлайн: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Ошибка Питона (%d): %s\n", resp.StatusCode, string(body))
		s.sendMessage(chatID, "ИИ запутался в данных... 😵")
		return
	}

	var aiResp models.AIResponse
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		fmt.Printf("❌ Ошибка декодирования: %v\n", err)
		return
	}

	if aiResp.Reply == "" {
		fmt.Println("⚠️ ИИ прислал пустой ответ")
		s.sendMessage(chatID, "Мне нечего сказать по этому поводу... 🤔")
		return
	}

	s.sendMessage(chatID, aiResp.Reply)
}

func (s *BotService) sendMessage(chatID int,text string){
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.Token)

	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,

	}

	jsonData , _ := json.Marshal(payload)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))


	if err != nil{
		fmt.Printf("❌ Критическая ошибка HTTP: %v\n", err)
		return 
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Telegram API Error: %s\n", string(body))
	} else {
		fmt.Printf("✅ Сообщение успешно улетело в чат %d\n", chatID)
	}

}
