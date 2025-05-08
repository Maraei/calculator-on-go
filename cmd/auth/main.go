package main

import (
	"log"
	"net"

	"github.com/Maraei/calculator-on-go/internal/auth"
	authpb "github.com/Maraei/calculator-on-go/api/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {	
	store, err := auth.NewStore("users.db")
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	serverOptions := []grpc.ServerOption{
		grpc.MaxRecvMsgSize(50 * 1024 * 1024),
	}

	grpcServer := grpc.NewServer(serverOptions...)
	authService := auth.NewAuthServer(store)

	authpb.RegisterAuthCalculatorServiceServer(grpcServer, authService)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка прослушивания порта: %v", err)
	}

	log.Println("Auth gRPC сервер запущен на порту :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
