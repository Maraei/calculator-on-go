package agent

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/Maraei/calculator-on-go/api/api"
	"github.com/stretchr/testify/assert"
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

func (m *MockAuthCalculatorServiceClient) Login(ctx context.Context, in *api.AuthRequest, opts ...grpc.CallOption) (*api.TokenResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.TokenResponse), args.Error(1)
}

func TestLogin_Success(t *testing.T) {
	os.Setenv("AGENT_LOGIN", "test")
	os.Setenv("AGENT_PASSWORD", "1234")

	mockClient := new(MockAuthCalculatorServiceClient)
	mockClient.On("Login", mock.Anything, mock.Anything).Return(&api.TokenResponse{
		Token: "valid_token",
	}, nil)

	assert.Equal(t, "valid_token", token)

	mockClient.AssertExpectations(t)
}

func TestLogin_Error(t *testing.T) {
	os.Setenv("AGENT_LOGIN", "test")
	os.Setenv("AGENT_PASSWORD", "1234")

	mockClient := new(MockAuthCalculatorServiceClient)
	mockClient.On("Login", mock.Anything, mock.Anything).Return(nil, errors.New("ошибка авторизации"))

	token = ""
	assert.Empty(t, token)

	mockClient.AssertExpectations(t)
}

func TestWorker(t *testing.T) {
	mockClient := new(MockTaskServiceClient)
	mockClient.On("FetchTask", mock.Anything, mock.Anything).Return(&api.FetchTaskResponse{
		TaskId:    "task_123",
		Arg1:      10,
		Arg2:      2,
		Operation: "+",
	}, nil).Once()

	mockClient.On("SendResult", mock.Anything, mock.Anything).Return(&api.SendResultResponse{
		Success: true,
	}, nil).Once()

	go worker(1, mockClient)

	mockClient.AssertExpectations(t)
}
