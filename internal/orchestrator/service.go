package orchestrator

import (
	"sync"

	"github.com/google/uuid"
)

type Expression struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Result *float64 `json:"result,omitempty"`
}

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
	Result        float64 `json:"result,omitempty"`
}


type Service struct {
	mu          sync.Mutex
	expressions map[string]*Expression
	taskManager *TaskManager
}

func NewService(taskManager *TaskManager) *Service {
	return &Service{
		expressions: make(map[string]*Expression),
		taskManager: taskManager,
	}
}

// Добавление выражения
func (s *Service) AddExpression(expression string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	s.expressions[id] = &Expression{ID: id, Status: "pending"}

	s.taskManager.GenerateTasks(id, expression)

	return id, nil
}

// Получение списка выражений
func (s *Service) GetExpressions() []*Expression {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []*Expression
	for _, expr := range s.expressions {
		result = append(result, expr)
	}
	return result
}

// Получение выражения по ID
func (s *Service) GetExpressionByID(id string) (*Expression, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	expr, exists := s.expressions[id]
	return expr, exists
}

// Получение следующей задачи для агента
func (s *Service) GetNextTask() (*Task, bool) {
	return s.taskManager.GetNextTask()
}

// Обработка результата от агента
func (s *Service) SubmitTaskResult(id string, result float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.taskManager.CompleteTask(id, result); err != nil {
		return err
	}

	// Проверяем, завершены ли все вычисления выражения
	if exprID, allDone, finalResult := s.taskManager.CheckExpressionCompletion(id); allDone {
		s.expressions[exprID].Status = "completed"
		s.expressions[exprID].Result = &finalResult
	}

	return nil
}
