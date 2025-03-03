package orchestrator

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type TaskManager struct {
	mu      sync.Mutex
	tasks   map[string]*Task
	results map[string]float64
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:   make(map[string]*Task),
		results: make(map[string]float64),
	}
}

// Генерация задач из выражения с параллельным вычислением
func (tm *TaskManager) GenerateTasks(expressionID, expression string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	re := regexp.MustCompile(`(\d+(\.\d+)?|\+|\-|\*|\/)`)
	tokens := re.FindAllString(expression, -1)
	log.Printf("Токены после разбиения регэкспом: %+v", tokens)

	var wg sync.WaitGroup
	results := make(chan *Task, 1) // Оставляем буфер на 1 задачу

	// В цикле будем создавать задачу только для первого оператора.
	for i := 1; i < len(tokens)-1; i += 2 {
		arg1, err1 := strconv.ParseFloat(tokens[i-1], 64)
		arg2, err2 := strconv.ParseFloat(tokens[i+1], 64)
		operation := tokens[i]

		if err1 != nil || err2 != nil {
			log.Printf("Ошибка парсинга аргументов: %v, %v", err1, err2)
			continue
		}

		// Используем переданный expressionID без добавления суффиксов
		taskID := expressionID

		wg.Add(1)
		go func(taskID string, arg1, arg2 float64, operation string) {
			defer wg.Done()

			result, err := performCalculation(arg1, arg2, operation)
			if err != nil {
				log.Printf("Ошибка выполнения задачи %s: %v", taskID, err)
				return
			}

			task := &Task{
				ID:            taskID, // Используем именно expressionID
				Arg1:          arg1,
				Arg2:          arg2,
				Operation:     operation,
				OperationTime: rand.Intn(3000) + 1000,
				Result:        result,
			}

			results <- task
		}(taskID, arg1, arg2, operation)

		// Если задача успешно создана для первого оператора, выходим из цикла
		break
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Сохраняем полученную задачу в мапу задач
	for task := range results {
		tm.tasks[task.ID] = task
		log.Printf("Добавлена задача: %s (%f %s %f) = %f", task.ID, task.Arg1, task.Operation, task.Arg2, task.Result)
	}
}


// Выполнение вычислений для одного задания
func performCalculation(arg1, arg2 float64, operation string) (float64, error) {
	switch operation {
	case "+":
		return arg1 + arg2, nil
	case "-":
		return arg1 - arg2, nil
	case "*":
		return arg1 * arg2, nil
	case "/":
		if arg2 == 0 {
			return 0, fmt.Errorf("деление на ноль")
		}
		return arg1 / arg2, nil
	default:
		return 0, fmt.Errorf("неизвестная операция: %s", operation)
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
