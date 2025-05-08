package integration

import (
	"context"
	"testing"

	agentpb "github.com/Maraei/calculator-on-go/api/api"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestAgentTaskLifecycle(t *testing.T) {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	authClient := agentpb.NewAuthCalculatorServiceClient(conn)
	client := agentpb.NewOrchestratorServiceClient(conn)

	regResp, err := authClient.Register(context.Background(), &agentpb.AuthRequest{
		Username: "testuser",
		Password: "testpass",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, regResp.Message)

	loginResp, err := authClient.Login(context.Background(), &agentpb.AuthRequest{
		Username: "testuser",
		Password: "testpass",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, loginResp.Token)

	md := metadata.New(map[string]string{
		"authorization": "Bearer " + loginResp.Token,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	addResp, err := client.AddExpression(ctx, &agentpb.AddExpressionRequest{
		UserId:    0,
		Expression: "4 * 5",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, addResp.Id)

	taskResp, err := client.GetTask(ctx, &agentpb.GetTaskRequest{})
	assert.NoError(t, err)
	assert.NotEmpty(t, taskResp.Task.Id)

	_, err = client.SubmitResult(ctx, &agentpb.SubmitResultRequest{
		Id:     taskResp.Task.Id,
		Result: 20,
		Error:  "",
	})
	assert.NoError(t, err)

	resultResp, err := client.GetResult(ctx, &agentpb.GetResultRequest{
		Id: addResp.Id,
	})
	assert.NoError(t, err)
	assert.Equal(t, float32(20), resultResp.Result)
	assert.Equal(t, "done", resultResp.Status)
}
