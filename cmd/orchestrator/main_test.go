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

	// Добавляем тестовую задачу (теперь будет создана только одна задача с ID "expr_1")
	taskManager.GenerateTasks("expr_1", "2 + 2")

	// Создаём объект сервиса
	service := orchestrator.NewService(taskManager)

	// Создаём обработчик
	handler := orchestrator.NewHandler(service)

	// Создаём тестовый HTTP-запрос на путь "/internal/task"
	req, err := http.NewRequest("GET", "/internal/task", nil)
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}

	rr := httptest.NewRecorder()

	// Регистрируем обработчик на роутере
	router := mux.NewRouter()
	router.HandleFunc("/internal/task", handler.GetTask).Methods("GET")
	router.ServeHTTP(rr, req)

	// Проверяем код ответа
	if rr.Code != http.StatusOK {
		t.Fatalf("Ожидали код 200, но получили %v", rr.Code)
	}

	// Декодируем тело ответа
	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Ошибка при декодировании ответа: %v", err)
	}

	// Обработчик возвращает объект {"task": { ... }}, поэтому извлекаем вложенный объект
	taskData, ok := response["task"].(map[string]interface{})
	if !ok {
		t.Fatalf("Ожидали найти объект 'task' в ответе, но получили: %v", response)
	}

	// Проверяем, что поле "id" присутствует и не пустое
	id, exists := taskData["id"].(string)
	if !exists || id == "" {
		t.Fatalf("Ожидали получить ID задачи, но его нет или он пустой: %v", taskData)
	}

	// Дополнительная проверка: убеждаемся, что ID задачи соответствует "expr_1"
	if id != "expr_1" {
		t.Fatalf("Ожидали получить ID 'expr_1', но получили %v", id)
	}
}
