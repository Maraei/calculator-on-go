package integration

import (
	"context"
	"testing"

	"github.com/Maraei/calculator-on-go/api/api"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockTaskServiceClient struct {
	mock.Mock
}

func (m *MockTaskServiceClient) FetchTask(ctx context.Context, in *api.FetchTaskRequest, opts ...grpc.CallOption) (*api.FetchTaskResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.FetchTaskResponse), args.Error(1)
}

func (m *MockTaskServiceClient) SendResult(ctx context.Context, in *api.SendResultRequest, opts ...grpc.CallOption) (*api.SendResultResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.SendResultResponse), args.Error(1)
}

type MockAuthCalculatorServiceClient struct {
	mock.Mock
}

func (m *MockAuthCalculatorServiceClient) Login(ctx context.Context, in *api.AuthRequest, opts ...grpc.CallOption) (*api.AuthResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.AuthResponse), args.Error(1)
}

func (m *MockAuthCalculatorServiceClient) Register(ctx context.Context, in *api.AuthRequest, opts ...grpc.CallOption) (*api.AuthResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.AuthResponse), args.Error(1)
}

func TestAgentTaskLifecycle(t *testing.T) {
	authClient := new(MockAuthCalculatorServiceClient)
	taskClient := new(MockTaskServiceClient)

	authClient.On("Register", mock.Anything, mock.Anything).Return(&api.AuthResponse{
		Message: "Registration successful",
	}, nil)

	taskClient.On("FetchTask", mock.Anything, mock.Anything).Return(&api.FetchTaskResponse{
		TaskId:    "task1",
		Arg1:      4,
		Arg2:      5,
		Operation: "*",
	}, nil)
	taskClient.On("SendResult", mock.Anything, mock.Anything).Return(&api.SendResultResponse{
		Success: true,
	}, nil)
}
