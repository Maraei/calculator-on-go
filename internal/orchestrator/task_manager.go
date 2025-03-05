package orchestrator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"math/rand"
	"log"
    
	"github.com/Maraei/calculator-on-go/internal/agent"
)

type TaskManager struct {
	mu           sync.Mutex
	tasks        map[string]*Task
	results      map[string]TaskResult
	expressions  map[string]*Expression
}

type TaskResult struct {
	Result float64 `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:       make(map[string]*Task),
		results:     make(map[string]TaskResult),
		expressions: make(map[string]*Expression),
	}
}

var precedence = map[string]int{
	"+": 1, "-": 1, "*": 2, "/": 2,
}

func toRPN(expression string) ([]string, error) {
    tokens, err := tokenizeExpression(expression)
    if err != nil {
        return nil, fmt.Errorf("%v", err)
    }
    var output []string
    var stack []string

    for _, token := range tokens {
        if _, err := strconv.ParseFloat(token, 64); err == nil {
            output = append(output, token)
        } else if token == "(" {
            stack = append(stack, token)
        } else if token == ")" {
            for len(stack) > 0 && stack[len(stack)-1] != "(" {
                output = append(output, stack[len(stack)-1])
                stack = stack[:len(stack)-1]
            }
            if len(stack) > 0 && stack[len(stack)-1] == "(" {
                stack = stack[:len(stack)-1]
            }
        } else {
            for len(stack) > 0 && precedence[stack[len(stack)-1]] >= precedence[token] {
                output = append(output, stack[len(stack)-1])
                stack = stack[:len(stack)-1]
            }
            stack = append(stack, token)
        }
    }

    for len(stack) > 0 {
        output = append(output, stack[len(stack)-1])
        stack = stack[:len(stack)-1]
    }

    return output, nil
}

func tokenizeExpression(expression string) ([]string, error) {
    re := regexp.MustCompile(`(\d+(\.\d*)?|\+|\-|\*|\/|\(|\))`)
    matches := re.FindAllString(expression, -1)
    if matches == nil {
        return nil, fmt.Errorf("разрешены только числа и ( ) + - * /")
    }
    return matches, nil
}

func (tm *TaskManager) GenerateTasks(expressionID, expression string) ([]string, error) {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    rpn, err := toRPN(expression)
    if err != nil {
        return nil, fmt.Errorf("выражение неверно: %v", err)
    }

    tm.expressions[expressionID] = &Expression{
        ID:     expressionID,
        Input:  expression,
        Status: "pending",
    }

    var stack []string
    taskIDs := []string{}

    for _, token := range rpn {
        if _, err := strconv.ParseFloat(token, 64); err == nil {
            stack = append(stack, token)
        } else {
            if len(stack) < 2 {
                return nil, fmt.Errorf("недостаточно операндов для операции %s", token)
            }

            arg2 := stack[len(stack)-1]
            arg1 := stack[len(stack)-2]
            stack = stack[:len(stack)-2]
            taskID := expressionID
            task := &Task{
                ID:            taskID,
                Arg1:          parseFloat(arg1),
                Arg2:          parseFloat(arg2),
                Operation:     token,
                OperationTime: rand.Intn(3000) + 1000,
            }
            tm.tasks[taskID] = task
            taskIDs = append(taskIDs, taskID)

            stack = append(stack, taskID)
        }
    }

    return taskIDs, nil
}

func parseFloat(str string) float64 {
    val, _ := strconv.ParseFloat(str, 64)
    return val
}

func (tm *TaskManager) GetNextTask() (*Task, bool) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	for id, task := range tm.tasks {
		delete(tm.tasks, id)
		return task, true
	}
	return nil, false
}

func (tm *TaskManager) CompleteTask(taskID string, result float64) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.results[taskID] = TaskResult{Result: result}
	return nil
}

func (tm *TaskManager) CompleteTaskWithError(taskID string, errorMsg string) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.results[taskID] = TaskResult{Error: errorMsg}
	return nil
}

func (tm *TaskManager) CheckExpressionCompletion(taskID string) (string, bool, float64, bool) {
    tm.mu.Lock()
    defer tm.mu.Unlock()

    exprID := strings.Split(taskID, "_")[0]
    if exprID == "" {
        log.Println("Ошибка: пустой exprID")
        return "", false, 0, false
    }

    expr, exists := tm.expressions[exprID]
    if !exists {
        log.Printf("Ошибка: выражение с ID %s не найдено", exprID)
        return exprID, true, 0, true
    }

    hasError := false
    results := make(map[string]float64)

    for id, res := range tm.results {
        if strings.HasPrefix(id, exprID) {
            if res.Error != "" {
                hasError = true
                break
            }
            results[id] = res.Result
        }
    }

    if hasError {
        log.Printf("Ошибка в задачах для выражения %s", exprID)
        return exprID, true, 0, true
    }

    for id := range tm.tasks {
        if strings.HasPrefix(id, exprID) {
            log.Printf("Невыполненные задачи для выражения %s", exprID)
            return "", false, 0, false
        }
    }

    var stack []float64
    rpn, _ := toRPN(expr.Input)
    for _, token := range rpn {
        if num, err := strconv.ParseFloat(token, 64); err == nil {
            stack = append(stack, num)
        } else {
            if len(stack) < 2 {
                log.Printf("Ошибка: недостаточно операндов для операции %s", token)
                return exprID, true, 0, true
            }
            arg2 := stack[len(stack)-1]
            arg1 := stack[len(stack)-2]
            stack = stack[:len(stack)-2]

            result, err := agent.Calculate(arg1, arg2, token)
            if err != nil {
                log.Printf("Ошибка вычисления: %v", err)
                return exprID, true, 0, true
            }
            stack = append(stack, result)
        }
    }

    if len(stack) != 1 {
        log.Printf("Ошибка: некорректный результат вычислений")
        return exprID, true, 0, true
    }

    log.Printf("Выражение %s успешно завершено, результат: %v", exprID, stack[0])
    return exprID, true, stack[0], false
}
