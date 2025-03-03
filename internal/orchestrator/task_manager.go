package orchestrator

import (
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"log"
)

type TaskManager struct {
	mu    sync.Mutex
	tasks map[string]*Task
	results map[string]float64
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:   make(map[string]*Task),
		results: make(map[string]float64),
	}
}
// AddTask добавляет одну задачу в список
func (tm *TaskManager) AddTask(taskID string, task *Task) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.tasks[taskID] = task
}

// Генерация задач из выражения
func (tm *TaskManager) GenerateTasks(expressionID, expression string) {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    tokens := strings.Split(expression, " ")
    for i := 1; i < len(tokens)-1; i += 2 {
        arg1, err1 := strconv.ParseFloat(tokens[i-1], 64)
        arg2, err2 := strconv.ParseFloat(tokens[i+1], 64)
        operation := tokens[i]

        if err1 != nil || err2 != nil {
            log.Printf("Ошибка парсинга аргументов: %v, %v", err1, err2)
            continue
        }

        taskID := expressionID + "_" + strconv.Itoa(i)

        tm.tasks[taskID] = &Task{
            ID:            taskID,
            Arg1:          arg1,
            Arg2:          arg2,
            Operation:     operation,
            OperationTime: rand.Intn(3000) + 1000,
        }

        log.Printf("Добавлена задача: %s (%f %s %f)", taskID, arg1, operation, arg2)
    }
}


// Получение следующей задачи
func (tm *TaskManager) GetNextTask() (*Task, bool) {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    log.Println("Все текущие задачи:", tm.tasks)

    for id, task := range tm.tasks {
        log.Println("Выдаётся задача:", id)
        delete(tm.tasks, id)
        return task, true
    }
    return nil, false
}


// Завершение задачи
func (tm *TaskManager) CompleteTask(taskID string, result float64) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.results[taskID] = result
	return nil
}

// Проверка завершения выражения
func (tm *TaskManager) CheckExpressionCompletion(taskID string) (string, bool, float64) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	exprID := strings.Split(taskID, "_")[0]

	for id := range tm.tasks {
		if strings.HasPrefix(id, exprID) {
			return "", false, 0
		}
	}

	var finalResult float64
	for id, res := range tm.results {
		if strings.HasPrefix(id, exprID) {
			finalResult += res
		}
	}
	return exprID, true, finalResult
}
