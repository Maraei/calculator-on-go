package orchestrator

import (
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type TaskManager struct{}

func NewTaskManager() *TaskManager {
	return &TaskManager{}
}

func (tm *TaskManager) GenerateTasks(expressionID, expression string) ([]*Task, error) {
	tokens := strings.Fields(expression)
	stack := []string{}
	var tasks []*Task

	for _, token := range tokens {
		if isOperator(token) {
			if len(stack) < 2 {
				return nil, errors.New("некорректное выражение")
			}
			arg2, _ := strconv.ParseFloat(stack[len(stack)-1], 64)
			arg1, _ := strconv.ParseFloat(stack[len(stack)-2], 64)
			stack = stack[:len(stack)-2]

			task := &Task{
				ID:           uuid.New().String(),
				ExpressionID: expressionID,
				Arg1:         arg1,
				Arg2:         arg2,
				Operation:    token,
				Status:       "pending",
			}
			stack = append(stack, "0") // временная замена результата
			tasks = append(tasks, task)
		} else {
			stack = append(stack, token)
		}
	}
	return tasks, nil
}

func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}
