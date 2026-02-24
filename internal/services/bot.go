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
	Token string
}

func NewBotService(token string) *BotService {
	return &BotService{
		Token: token,
	}
}

func (s *BotService) ProcessUpdate(update models.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	userText := update.Message.Text

	if userText == "/start" {
		s.sendMessage(chatID, "Привет! Я ИИ-помощник архитектурного бюро. Спрашивай про наши проекты!")
		return
	}

	var voiceURL string

	if update.Message.Voice != nil {
		fmt.Println("🎙️ Получено голосовое сообщение! Обрабатываю...")
		url, err := s.getFileURL(update.Message.Voice.FileID)
		if err != nil {
			fmt.Printf("❌ Ошибка получения аудио: %v\n", err)
			s.sendMessage(chatID, "Не удалось загрузить голосовое сообщение 😔")
			return
		}
		voiceURL = url
	}

	aiReq := models.AIRequest{
		ChatID:   int64(chatID),
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
		fmt.Printf("❌ Ошибка декодирования: %v\n", err)
		return
	}

	if aiResp.Reply == "" {
		s.sendMessage(chatID, "Мне нечего сказать по этому поводу... 🤔")
		return
	}

	s.sendMessage(chatID, aiResp.Reply)
}

func (s *BotService) sendMessage(chatID int, text string) {
	// ✅ ИСПРАВЛЕНО: Здесь ВСЕГДА должен быть api.telegram.org!
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.Token)

	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		fmt.Printf("❌ Ошибка отправки в Telegram: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("❌ Telegram API Error: %s\n", string(body))
	} else {
		fmt.Printf("✅ Ответ улетел клиенту в чат %d\n", chatID)
	}
}

func (s *BotService) getFileURL(fileID string) (string, error) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", s.Token, fileID)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Ok     bool `json:"ok"`
		Result struct {
			FilePath string `json:"file_path"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || !result.Ok {
		return "", fmt.Errorf("ошибка API Telegram")
	}
	
	downloadURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", s.Token, result.Result.FilePath)
	return downloadURL, nil
}