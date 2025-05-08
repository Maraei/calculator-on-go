package orchestrator_test

import (
	"context"
	"testing"
	"time"

	"github.com/Maraei/calculator-on-go/internal/orchestrator"
	agentpb "github.com/Maraei/calculator-on-go/api/api"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	err = orchestrator.Migrate(db)
	assert.NoError(t, err)

	return db
}

func TestAddExpressionAndGetResult(t *testing.T) {
	db := setupTestDB(t)
	repo := orchestrator.NewRepository(db)
	taskManager := orchestrator.NewTaskManager()
	service := orchestrator.NewService(repo, taskManager)
	handler := orchestrator.NewHandler(service)

	ctx := context.Background()

	addResp, err := handler.AddExpression(ctx, &agentpb.AddExpressionRequest{
		UserId:     1,
		Expression: "2 + 3",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, addResp.Id)

	expr, err := service.GetExpressionByID(addResp.Id)
	assert.NoError(t, err)
	assert.Equal(t, "pending", expr.Status)
	assert.Equal(t, "2 + 3", expr.Input)

	tasks, err := repo.GetPendingTasks(ctx)

	assert.NoError(t, err)
	assert.Len(t, tasks, 1)

	task := tasks[0]
	_, err = handler.SendResult(ctx, &agentpb.SendResultRequest{
		TaskId: task.ID,
		Result: 5.0,
	})	
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 50)

	resultResp, err := handler.GetResult(ctx, &agentpb.GetResultRequest{Id: addResp.Id})
	assert.NoError(t, err)
	assert.Equal(t, float32(5.0), resultResp.Result)
	assert.Equal(t, "completed", resultResp.Status)
	assert.Empty(t, resultResp.Error)
}

func TestFetchTask_NoTasks(t *testing.T) {
	db := setupTestDB(t)
	repo := orchestrator.NewRepository(db)
	taskManager := orchestrator.NewTaskManager()
	service := orchestrator.NewService(repo, taskManager)
	handler := orchestrator.NewHandler(service)

	ctx := context.Background()
	_, err := handler.FetchTask(ctx, &agentpb.FetchTaskRequest{})
	assert.Error(t, err)
}

func TestExpressionWithDivisionByZero(t *testing.T) {
	db := setupTestDB(t)
	repo := orchestrator.NewRepository(db)
	taskManager := orchestrator.NewTaskManager()
	service := orchestrator.NewService(repo, taskManager)
	handler := orchestrator.NewHandler(service)

	ctx := context.Background()
	addResp, err := handler.AddExpression(ctx, &agentpb.AddExpressionRequest{
		UserId:     1,
		Expression: "4 / 0",
	})
	assert.NoError(t, err)

	tasks, err := repo.GetPendingTasks(ctx)
	assert.NoError(t, err)
	assert.Len(t, tasks, 1)

	task := tasks[0]
	_, err = handler.SendResult(ctx, &agentpb.SendResultRequest{
		TaskId:       task.ID,
		ErrorMessage: "division by zero",
	})
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 50)

	resp, err := handler.GetResult(ctx, &agentpb.GetResultRequest{Id: addResp.Id})
	assert.NoError(t, err)
	assert.Equal(t, "error", resp.Status)
	assert.Contains(t, resp.Error, "division by zero")
}