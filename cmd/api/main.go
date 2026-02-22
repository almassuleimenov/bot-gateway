package main


import (
	"fmt"
	"log"
	"os"
	"net/http"

	"bot-gateway/internal/handlers"
	"bot-gateway/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

)

func main(){
	err := godotenv.Load()
	if err != nil{
		log.Println("Предупреждение: Файл .env не найден,теперь будем использовать системные переменные")
	}
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if botToken == ""{
		log.Fatal("TELEGRAM_BOT_TOKEN не задан БРОО!!!")
	}

	botService := services.NewBotService(botToken)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/webhook",handlers.HandleWebhook(botService))

	port := ":8080"

	fmt.Printf("Api Cateway запущен на порту %s\n",port)

	err = http.ListenAndServe(port,r)
	if err != nil{
		fmt.Printf("Ошибка при запуске сервера: %v\n",err)
	}

}