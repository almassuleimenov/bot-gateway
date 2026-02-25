package services

import (
	"bot-gateway/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type BotService struct {
	ApiUrl   string
	IdInst   string
	ApiToken string
}

func NewBotService(apiUrl, idInst, apiToken string) *BotService {
	return &BotService{
		ApiUrl:   apiUrl,
		IdInst:   idInst,
		ApiToken: apiToken,
	}
}

func (s *BotService) ProcessUpdate(webhook models.GreenApiWebhook) {
	chatID := webhook.SenderData.ChatId
	msgType := webhook.MessageData.TypeMessage

	if chatID == "" {
		return
	}

	var userText string
	var voiceURL string

	if msgType == "textMessage" {
		userText = webhook.MessageData.TextMessageData.TextMessage
		fmt.Println("📩 Получен ТЕКСТ из WhatsApp:", userText)
	} else if msgType == "audioMessage" {
		voiceURL = webhook.MessageData.FileMessageData.DownloadUrl
		fmt.Println("🎙️ Получено ГОЛОСОВОЕ сообщение, ссылка:", voiceURL)
	}

	if userText == "" && voiceURL == "" {
		return
	}

	aiReq := models.AIRequest{
		ChatID:   chatID,
		UserText: userText,
		VoiceURL: voiceURL,
	}

	jsonData, _ := json.Marshal(aiReq)

	brainURL := "https://bot-brain-k9bb.onrender.com/generate-answer"
	resp, err := http.Post(brainURL, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Printf("❌ Питон оффлайн: %v\n", err)
		s.sendMessage(chatID, "Мой мозг сейчас обновляется, подождите минутку... 🧠🔄")
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
		fmt.Printf("❌ Ошибка декодирования ответа Питона: %v\n", err)
		return
	}

	if aiResp.Reply == "" {
		s.sendMessage(chatID, "Мне нечего сказать по этому поводу... 🤔")
		return
	}

	s.sendMessage(chatID, aiResp.Reply)
}

func (s *BotService) sendMessage(chatID string, text string) {
	sendUrl := fmt.Sprintf("%s/waInstance%s/sendMessage/%s", s.ApiUrl, s.IdInst, s.ApiToken)

	payload := models.GreenApiSendRequest{
		ChatId:  chatID,
		Message: text,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(sendUrl, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Printf("❌ Ошибка отправки в WhatsApp: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("✅ Ответ улетел клиенту в WhatsApp: %s\n", chatID)
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Ошибка API WhatsApp: %s\n", string(body))
	}
}