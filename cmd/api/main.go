package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"bot-gateway/internal/handlers"
	"bot-gateway/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Предупреждение: Файл .env не найден, используем системные переменные")
	}

	apiUrl := os.Getenv("GREEN_API_URL")
	idInst := os.Getenv("GREEN_API_ID")
	apiToken := os.Getenv("GREEN_API_TOKEN")

	if apiUrl == "" || idInst == "" || apiToken == "" {
		log.Fatal("Ключи GREEN_API не заданы БРОО!!!")
	}

	botService := services.NewBotService(apiUrl, idInst, apiToken)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/webhook", handlers.HandleWebhook(botService))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("API Gateway запущен на порту %s\n", port)

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %v\n", err)
	}
}