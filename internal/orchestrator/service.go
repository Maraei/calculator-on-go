package orchestrator

import (
	"errors"
	"fmt"
	"strings"
)

type Service struct {
	repo        *Repository
	taskManager *TaskManager
}

func NewService(repo *Repository, tm *TaskManager) *Service {
	return &Service{
		repo:        repo,
		taskManager: tm,
	}
}

func (s *Service) AddExpression(userID uint, input string) (string, error) {
	rpnTokens, err := InfixToRPN(input)
	if err != nil {
		return "", fmt.Errorf("не удалось разобрать выражение: %w", err)
	}
	rpn := strings.Join(rpnTokens, " ")

	expr, err := s.repo.CreateExpression(userID, input)
	if err != nil {
		return "", err
	}
	
	tasks, err := s.taskManager.GenerateTasks(expr.ID, rpn)
	if err != nil {
		return "", err
	}
	for _, t := range tasks {
		if err := s.repo.CreateTask(t); err != nil {
			return "", err
		}
	}
	return expr.ID, nil
}

func (s *Service) GetNextTask() (*Task, error) {
	return s.repo.GetPendingTask()
}

func (s *Service) SubmitTaskResult(taskID string, result float64) error {
	task, err := s.repo.GetPendingTask()
	if err != nil {
		return err
	}

	if task.ID != taskID {
		return errors.New("task ID mismatch")
	}

	if err := s.repo.CompleteTask(taskID, result); err != nil {
		return err
	}

	return s.checkExpressionCompletion(task.ExpressionID)
}

func (s *Service) SubmitTaskError(taskID string, errorMsg string) error {
	task, err := s.repo.GetPendingTask()
	if err != nil {
		return err
	}

	if task.ID != taskID {
		return errors.New("task ID mismatch")
	}

	if err := s.repo.FailTask(taskID, errorMsg); err != nil {
		return err
	}

	return s.repo.UpdateExpressionError(task.ExpressionID, errorMsg)
}

func (s *Service) GetExpressionByID(id string) (*Expression, error) {
	return s.repo.GetExpressionByID(id)
}

func (s *Service) checkExpressionCompletion(expressionID string) error {
	tasks, err := s.repo.GetTasksByExpression(expressionID)
	if err != nil {
		return err
	}

	allCompleted := true
	var total float64
	for _, task := range tasks {
		if task.Status == "error" {
			return s.repo.UpdateExpressionError(expressionID, *task.Error)
		}
		if task.Status != "completed" {
			allCompleted = false
		}
		if task.Result != nil {
			total += *task.Result
		}
	}

	if allCompleted {
		return s.repo.UpdateExpressionResult(expressionID, total)
	}

	return nil
}
