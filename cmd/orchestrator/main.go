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
	"google.golang.org/grpc/reflection"
)

func main() {
	serverOptions := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(50 * 1024 * 1024), // Увеличиваем максимальный размер получаемого сообщения до 50 MB
	}
	// Открываем базу данных SQLite
	db, err := gorm.Open(sqlite.Open("orchestrator.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось открыть базу данных: %v", err)
	}

	// Миграция таблиц
	if err := orchestrator.Migrate(db); err != nil {
		log.Fatalf("Не удалось провести миграцию: %v", err)
	}
	server := grpc.NewServer(serverOptions...)
	// Создаем репозиторий, сервис и обработчик оркестратора
	repo := orchestrator.NewRepository(db)
	taskManager := orchestrator.NewTaskManager()
	service := orchestrator.NewService(repo, taskManager)
	handler := orchestrator.NewHandler(service)
	reflection.Register(server)
	// Создаем хранилище и сервис для аутентификации
	store, err := auth.NewStore("auth.db")
	if err != nil {
		log.Fatalf("Не удалось создать хранилище: %v", err)
	}
	authService := auth.NewAuthServer(store)

	// Настраиваем и запускаем gRPC сервер
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}

	grpcServer := grpc.NewServer()
	agentpb.RegisterTaskServiceServer(grpcServer, handler)
	agentpb.RegisterAuthCalculatorServiceServer(grpcServer, authService)
	agentpb.RegisterOrchestratorServiceServer(grpcServer, handler)

	log.Println("Оркестратор запущен на порту 50052")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка сервера: %v", err)
	}
}
