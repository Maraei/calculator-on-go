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

	return s.CheckExpressionCompletion(task.ExpressionID)
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

func (s *Service) CheckExpressionCompletion(expressionID string) error {
	tasks, err := s.repo.GetTasksByExpression(expressionID)
	if err != nil {
		return fmt.Errorf("failed to get tasks: %w", err)
	}

	for _, task := range tasks {
		if task.Status != "completed" {
			return nil
		}
	}

	expr, err := s.repo.GetExpressionByID(expressionID)
	if err != nil {
		return fmt.Errorf("failed to get expression: %w", err)
	}

	rpnTokens, err := InfixToRPN(expr.Input)
	if err != nil {
		return s.repo.UpdateExpressionError(expressionID, fmt.Sprintf("failed to convert to RPN: %v", err))
	}

	result, err := EvaluateRPN(rpnTokens)
	if err != nil {
		return s.repo.UpdateExpressionError(expressionID, fmt.Sprintf("failed to evaluate RPN: %v", err))
	}

	return s.repo.UpdateExpressionResult(expressionID, result)
}
