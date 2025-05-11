package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gorilla/mux"
	"github.com/Maraei/calculator-on-go/internal/orchestrator"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGetTask(t *testing.T) {
	// Открытие базы данных в памяти
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	// Миграция базы данных
	if err := orchestrator.Migrate(db); err != nil {
		t.Fatalf("Ошибка при миграции базы данных: %v", err)
	}

	// Инициализация TaskManager
	taskManager := orchestrator.NewTaskManager()

	// Генерация задач для выражения
	tasks, err := taskManager.GenerateTasks("expr_1", "2 + 2")
	if err != nil {
		t.Fatalf("Ошибка при генерации задачи: %v", err)
	}
	t.Logf("Задачи успешно сгенерированы: %+v", tasks)

	// Создание GET-запроса
	req, err := http.NewRequest("GET", "/internal/task/expr_1", nil)
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}

	// Рекордер для записи ответа
	rr := httptest.NewRecorder()

	// Создание маршрутизатора
	router := mux.NewRouter()
	router.ServeHTTP(rr, req)

	// Проверка кода ответа
	if rr.Code != http.StatusOK {
		t.Fatalf("Ожидали код 200, но получили %v", rr.Code)
	}

	// Декодирование ответа
	var response map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Ошибка при декодировании ответа: %v", err)
	}

	// Проверка наличия задачи в ответе
	taskData, ok := response["task"].(map[string]interface{})
	if !ok {
		t.Fatalf("Ожидали найти объект 'task' в ответе, но получили: %v", response)
	}

	// Проверка ID задачи
	id, exists := taskData["id"].(string)
	if !exists || id == "" {
		t.Fatalf("Ожидали получить ID задачи, но его нет или он пустой: %v", taskData)
	}

	// Проверка правильности ID
	if id != "expr_1" {
		t.Fatalf("Ожидали получить ID 'expr_1', но получили %v", id)
	}
}
