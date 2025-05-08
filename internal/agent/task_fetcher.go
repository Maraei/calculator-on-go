package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Maraei/calculator-on-go/api/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var token string

func GetAuthServerAddress() string {
	addr := os.Getenv("AUTH_SERVER_ADDRESS")
	if addr == "" {
		addr = "localhost:50051"
	}
	return addr
}

func GetOrchestratorAddress() string {
	addr := os.Getenv("ORCHESTRATOR_ADDRESS")
	if addr == "" {
		addr = "localhost:50052"
	}
	return addr
}

func Start(workerCount int) error {
	taskConn, err := grpc.Dial(GetOrchestratorAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("не удалось подключиться к оркестратору: %w", err)
	}
	defer taskConn.Close()

	taskClient := api.NewTaskServiceClient(taskConn)

	for i := 0; i < workerCount; i++ {
		go worker(i, taskClient)
	}

	select {}
}

func Login(authClient api.AuthCalculatorServiceClient) error {
	login := os.Getenv("AGENT_LOGIN")
	password := os.Getenv("AGENT_PASSWORD")
	if login == "" || password == "" {
		return fmt.Errorf("AGENT_LOGIN или AGENT_PASSWORD не установлены")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := authClient.Login(ctx, &api.AuthRequest{
		Username: login,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("не удалось залогиниться: %w", err)
	}

	token = resp.Token
	log.Println("Успешный вход агента, токен получен")
	return nil
}

func worker(id int, client api.TaskServiceClient) {
	for {
		task, err := FetchTask(client)
		if err != nil {
			log.Printf("[Worker %d] Ошибка при получении задачи: %v", id, err)
			time.Sleep(2 * time.Second)
			continue
		}

		if task == nil {
			log.Printf("[Worker %d] Нет задач, ждём...", id)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Printf("[Worker %d] Выполняем задачу: %f %s %f", id, task.Arg1, task.Operation, task.Arg2)

		result, err := Calculate(float64(task.Arg1), float64(task.Arg2), task.Operation)
		if err != nil {
			log.Printf("[Worker %d] Ошибка вычислений: %v", id, err)
			if err := SendResult(client, task.Id, 0, err.Error()); err != nil {
				log.Printf("[Worker %d] Ошибка отправки ошибки: %v", id, err)
			}
			continue
		}

		if err := SendResult(client, task.Id, result, ""); err != nil {
			log.Printf("[Worker %d] Ошибка отправки результата: %v", id, err)
		} else {
			log.Printf("[Worker %d] Результат отправлен успешно: %f", id, result)
		}
	}
}

func FetchTask(client api.TaskServiceClient) (*api.Task, error) {
	ctx := WithAuth(context.Background())
	resp, err := client.FetchTask(ctx, &api.FetchTaskRequest{})
	if err != nil {
		return nil, err
	}
	if resp.TaskId == "" {
		return nil, nil
	}
	return &api.Task{
		Id:        resp.TaskId,
		Arg1:      resp.Arg1,
		Arg2:      resp.Arg2,
		Operation: resp.Operation,
	}, nil
}

func SendResult(client api.TaskServiceClient, taskID string, result float64, errMsg string) error {
	ctx := WithAuth(context.Background())
	resp, err := client.SendResult(ctx, &api.SendResultRequest{
		TaskId:       taskID,
		Result:       float32(result),
		ErrorMessage: errMsg,
	})
	if err != nil {
		return err
	}
	if !resp.Success {
		return fmt.Errorf("сервер отклонил отправленный результат")
	}
	return nil
}

func WithAuth(ctx context.Context) context.Context {
	if token == "" {
		log.Println("Токен не установлен. Для авторизации используйте функцию login.")
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
}
