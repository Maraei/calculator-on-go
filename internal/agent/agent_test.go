package agent_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Maraei/calculator-on-go/internal/agent"
	"github.com/Maraei/calculator-on-go/api/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockAuthClient struct {
	mock.Mock
}

func (m *MockAuthClient) Login(ctx context.Context, in *api.AuthRequest, opts ...grpc.CallOption) (*api.TokenResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.TokenResponse), args.Error(1)
}

func (m *MockAuthClient) Register(ctx context.Context, in *api.AuthRequest, opts ...grpc.CallOption) (*api.AuthResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.AuthResponse), args.Error(1)
}

type MockTaskClient struct {
	mock.Mock
}

func (m *MockTaskClient) FetchTask(ctx context.Context, in *api.FetchTaskRequest, opts ...grpc.CallOption) (*api.FetchTaskResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.FetchTaskResponse), args.Error(1)
}

func (m *MockTaskClient) SendResult(ctx context.Context, in *api.SendResultRequest, opts ...grpc.CallOption) (*api.SendResultResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.SendResultResponse), args.Error(1)
}
func TestCalculate(t *testing.T) {
	tests := []struct {
		arg1      float64
		arg2      float64
		operation string
		expected  float64
		err       error
	}{
		{2, 3, "+", 5, nil},
		{5, 3, "-", 2, nil},
		{2, 3, "*", 6, nil},
		{6, 2, "/", 3, nil},
		{6, 0, "/", 0, errors.New("деление на ноль")},
		{5, 3, "unknown", 0, errors.New("неизвестная операция: unknown")},
	}

	for _, tt := range tests {
		t.Run(tt.operation, func(t *testing.T) {
			result, err := agent.Calculate(tt.arg1, tt.arg2, tt.operation)
			if tt.err != nil {
				assert.NotNil(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
func TestLogin(t *testing.T) {
	mockAuthClient := new(MockAuthClient)

	mockAuthClient.On("Validate", mock.Anything, mock.Anything).Return(&api.ValidateResponse{
		Valid: true,
	}, nil)

	mockAuthClient.On("Login", mock.Anything, mock.Anything).Return(&api.TokenResponse{
		Token: "test-token",
	}, nil).Once()


	mockAuthClient.On("Login", mock.Anything, mock.Anything).Return(nil, errors.New("ошибка авторизации")).Once()
}
func TestWorker(t *testing.T) {
	mockTaskClient := new(MockTaskClient)

	mockTaskClient.On("FetchTask", mock.Anything, mock.Anything).Return(&api.FetchTaskResponse{
		TaskId:    "task-1",
		Arg1:      5,
		Arg2:      3,
		Operation: "+",
	}, nil).Once()

	mockTaskClient.On("SendResult", mock.Anything, mock.Anything).Return(&api.SendResultResponse{
		Success: true,
	}, nil).Once()

}
