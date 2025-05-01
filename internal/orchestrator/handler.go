package orchestrator

import (
	"context"
	"time"

	agentpb "github.com/Maraei/calculator-on-go/api/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Expression struct {
    ID        string    `gorm:"primaryKey"`  // id
    UserID    uint32    `json:"user_id"`     // user_id (с нижним подчеркиванием)
    Input     string    `json:"input"`       // input
    Status    string    `json:"status"`      // status
    Result    *float64 `json:"result"`      // result
    Error     *string   `json:"error"`       // error
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Task struct {
	ID           string    `gorm:"primaryKey"`
	ExpressionID string
	Arg1         float64
	Arg2         float64
	Operation    string
	Status       string // pending, completed, error
	Result       *float64
	Error        *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique"`
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}


type Handler struct {
	service *Service
	agentpb.UnimplementedTaskServiceServer
	agentpb.UnimplementedOrchestratorServiceServer
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetResult(ctx context.Context, req *agentpb.GetResultRequest) (*agentpb.GetResultResponse, error) {
	expr, err := h.service.GetExpressionByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "expression not found: %v", err)
	}

	var result float32
	if expr.Result != nil {
		result = float32(*expr.Result)
	}

	var errorMessage string
	if expr.Error != nil {
		errorMessage = *expr.Error
	}

	return &agentpb.GetResultResponse{
		Result: result,
		Status: expr.Status,
		Error:  errorMessage,  // передаем ошибку как строку, если она существует
	}, nil
}

func (h *Handler) AddExpression(ctx context.Context, req *agentpb.AddExpressionRequest) (*agentpb.AddExpressionResponse, error) {
    id, err := h.service.AddExpression(uint(req.UserId), req.Expression)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to add expression: %v", err)
    }
    return &agentpb.AddExpressionResponse{Id: id}, nil
}

func (h *Handler) GetExpressions(ctx context.Context, req *agentpb.GetExpressionsRequest) (*agentpb.GetExpressionsResponse, error) {
	// Заглушка: в будущем можно реализовать метод в сервисе и репозитории
	return &agentpb.GetExpressionsResponse{
		Expressions: []*agentpb.Expression{}, // Пока возвращаем пустой список
	}, nil
}

func (h *Handler) GetExpressionByID(ctx context.Context, req *agentpb.GetExpressionByIDRequest) (*agentpb.GetExpressionByIDResponse, error) {
    expr, err := h.service.GetExpressionByID(req.Id)
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "expression not found: %v", err)
    }

    var result float32
    if expr.Result != nil {
        result = float32(*expr.Result)
    }

    return &agentpb.GetExpressionByIDResponse{
        Expression: &agentpb.Expression{
            Id:     expr.ID,
            UserId: uint32(expr.UserID),  // используем uint32
            Input:  expr.Input,
            Result: result,
            Status: expr.Status,
            Error: func() string {
                if expr.Error != nil {
                    return *expr.Error
                }
                return ""  // возвращаем пустую строку, если ошибки нет
            }(),
        },
    }, nil
}

func (h *Handler) FetchTask(ctx context.Context, req *agentpb.FetchTaskRequest) (*agentpb.FetchTaskResponse, error) {
	task, err := h.service.GetNextTask()
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "no pending tasks: %v", err)
	}
	return &agentpb.FetchTaskResponse{
		TaskId:    task.ID,
		Arg1:      float32(task.Arg1),
		Arg2:      float32(task.Arg2),
		Operation: task.Operation,
	}, nil
}

func (h *Handler) SendResult(ctx context.Context, req *agentpb.SendResultRequest) (*agentpb.SendResultResponse, error) {
	if req.ErrorMessage != "" {
		// Если ErrorMessage не пустой, сохраняем ошибку в базу
		err := h.service.SubmitTaskError(req.TaskId, req.ErrorMessage)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to submit task error: %v", err)
		}
	} else {
		// Если результат получен, сохраняем его
		err := h.service.SubmitTaskResult(req.TaskId, float64(req.Result))
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to submit task result: %v", err)
		}
	}
	return &agentpb.SendResultResponse{Success: true}, nil
}
