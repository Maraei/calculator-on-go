package main

import (
	"log"
	"net/http"

	"github.com/Maraei/calculator-on-go/internal/orchestrator"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env, используются переменные окружения по умолчанию")
	}
}

func main() {
	r := mux.NewRouter()

	taskManager := orchestrator.NewTaskManager()
	service := orchestrator.NewService(taskManager)
	handler := orchestrator.NewHandler(service)

	handler.RegisterRoutes(r)

	serverAddr := ":8080"
	log.Printf("Оркестратор запущен на %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}