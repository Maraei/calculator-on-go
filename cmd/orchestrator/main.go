package main

import (
	"log"
	"net"
	"github.com/Maraei/calculator-on-go/internal/orchestrator"
	"github.com/Maraei/calculator-on-go/internal/auth"
	agentpb "github.com/Maraei/calculator-on-go/api/api"
	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	HmacSampleSecret = "an7DkUH?L8iClxbVj5JZdbRVO2M$1Jc~D6CXsL@4"
)

func main() {
	db, err := gorm.Open(sqlite.Open("orchestrator.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось открыть базу данных: %v", err)
	}

	if err := orchestrator.Migrate(db); err != nil {
		log.Fatalf("Не удалось провести миграцию: %v", err)
	}

	repo := orchestrator.NewRepository(db)
	taskManager := orchestrator.NewTaskManager()
	service := orchestrator.NewService(repo, taskManager)
	handler := orchestrator.NewHandler(service)


	store, err := auth.NewStore("auth.db")
	if err != nil {
		log.Fatalf("Не удалось создать хранилище: %v", err)
	}
	authService := auth.NewAuthServer(store)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(auth.AuthMiddleware()),
	)
	agentpb.RegisterTaskServiceServer(grpcServer, handler)
	agentpb.RegisterAuthCalculatorServiceServer(grpcServer, authService)
	agentpb.RegisterOrchestratorServiceServer(grpcServer, handler)

	log.Println("Оркестратор запущен на порту 50052")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}
