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
	taskManager := orchestrator.NewTaskManager()

	taskManager.GenerateTasks("expr_1", "2 + 2")

	service := orchestrator.NewService(taskManager)
	handler := orchestrator.NewHandler(service)

	req, err := http.NewRequest("GET", "/internal/task", nil)
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/internal/task", handler.GetTask).Methods("GET")
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Ожидали код 200, но получили %v", rr.Code)
	}

	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Ошибка при декодировании ответа: %v", err)
	}

	taskData, ok := response["task"].(map[string]interface{})
	if !ok {
		t.Fatalf("Ожидали найти объект 'task' в ответе, но получили: %v", response)
	}

	id, exists := taskData["id"].(string)
	if !exists || id == "" {
		t.Fatalf("Ожидали получить ID задачи, но его нет или он пустой: %v", taskData)
	}

	if id != "expr_1" {
		t.Fatalf("Ожидали получить ID 'expr_1', но получили %v", id)
	}
}
