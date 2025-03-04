package orchestrator

import (
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"errors"
)

type TaskManager struct {
	mu      sync.Mutex
	tasks   map[string]*Task
	results map[string]TaskResult
}

type TaskResult struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:   make(map[string]*Task),
		results: make(map[string]TaskResult),
	}
}

// Генерация задач из выражения
func (tm *TaskManager) GenerateTasks(expressionID, expression string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Валидация выражения
	if err := validateExpression(expression); err != nil {
		return err
	}

	// Остальная логика генерации задач...
	re := regexp.MustCompile(`(\d+(\.\d+)?|\+|\-|\*|\/)`)
	tokens := re.FindAllString(expression, -1)
	log.Printf("Токены после разбиения регэкспом: %+v", tokens)

	// Если токенов недостаточно для создания задачи
	if len(tokens) < 3 {
		return errors.New("недостаточно токенов для создания задачи")
	}

	// Генерация задач...
	for i := 1; i < len(tokens)-1; i += 2 {
		arg1, err1 := strconv.ParseFloat(tokens[i-1], 64)
		arg2, err2 := strconv.ParseFloat(tokens[i+1], 64)
		operation := tokens[i]

		if err1 != nil || err2 != nil {
			log.Printf("Ошибка парсинга аргументов: %v, %v", err1, err2)
			return errors.New("неверный формат чисел в выражении")
		}

		taskID := expressionID
		tm.tasks[taskID] = &Task{
			ID:            taskID,
			Arg1:          arg1,
			Arg2:          arg2,
			Operation:     operation,
			OperationTime: rand.Intn(3000) + 1000,
		}
	}

	return nil
}

// Получение следующей задачи
func (tm *TaskManager) GetNextTask() (*Task, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for id, task := range tm.tasks {
		delete(tm.tasks, id)
		return task, true
	}
	return nil, false
}

// Завершение задачи с результатом
func (tm *TaskManager) CompleteTask(taskID string, result float64) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.results[taskID] = TaskResult{Result: result}
	return nil
}

// Завершение задачи с ошибкой
func (tm *TaskManager) CompleteTaskWithError(taskID string, errorMsg string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.results[taskID] = TaskResult{Error: errorMsg}
	return nil
}

// Проверка завершения выражения
func (tm *TaskManager) CheckExpressionCompletion(taskID string) (string, bool, float64, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	exprID := strings.Split(taskID, "_")[0]
	hasError := false
	var finalResult float64

	// Проверяем результаты задач
	for id, res := range tm.results {
		if strings.HasPrefix(id, exprID) {
			if res.Error != "" {
				hasError = true
				break
			}
			finalResult += res.Result
		}
	}

	// Если есть ошибка, возвращаем true и флаг ошибки
	if hasError {
		return exprID, true, 0, true
	}

	// Проверяем, остались ли невыполненные задачи
	for id := range tm.tasks {
		if strings.HasPrefix(id, exprID) {
			return "", false, 0, false
		}
	}

	// Если все задачи выполнены и нет ошибок, возвращаем финальный результат
	return exprID, true, finalResult, false
}



func validateExpression(expression string) error {
	// Регулярное выражение для проверки допустимых символов
	validChars := regexp.MustCompile(`^[\d\s\+\-\*\/\.]+$`)
	if !validChars.MatchString(expression) {
		return errors.New("выражение содержит недопустимые символы")
	}

	// Проверка на пустое выражение
	if len(strings.TrimSpace(expression)) == 0 {
		return errors.New("выражение пустое")
	}

	// Проверка на наличие хотя бы одного оператора
	hasOperator := regexp.MustCompile(`[\+\-\*\/]`).MatchString(expression)
	if !hasOperator {
		return errors.New("выражение не содержит операторов")
	}

	return nil
}