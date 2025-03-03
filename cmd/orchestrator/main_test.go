package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/Maraei/calculator-on-go/internal/orchestrator"
)

func TestGetTask(t *testing.T) {
	
	// Создаём объект TaskManager
	taskManager := orchestrator.NewTaskManager()

	// Добавляем тестовую задачу
	taskManager.GenerateTasks("expr_1", "2 + 2")

	// Создаём объект сервиса
	service := orchestrator.NewService(taskManager)

	// Создаём обработчик
	handler := orchestrator.NewHandler(service)

	// Создаём тестовый HTTP-запрос
	req, err := http.NewRequest("GET", "/api/v1/task", nil)
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}

	// Создаём записывающий ответ
	rr := httptest.NewRecorder()

	// Создаём роутер и связываем его с обработчиком
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/task", handler.GetTask).Methods("GET")
	router.ServeHTTP(rr, req)

	// Проверяем код ответа
	if rr.Code != http.StatusOK {
		t.Fatalf("Ожидали код 200, но получили %v", rr.Code)
	}

	// Проверяем тело ответа
	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Ошибка при декодировании ответа: %v", err)
	}

	// Проверяем, что ID задачи присутствует в ответе
	if response["id"] == nil {
		t.Fatalf("Ожидали получить ID задачи, но его нет")
	}
}
