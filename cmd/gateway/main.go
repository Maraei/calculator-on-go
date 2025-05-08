package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Maraei/calculator-on-go/api/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := api.RegisterAuthCalculatorServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
		log.Fatalf("Ошибка регистрации Auth-сервиса: %v", err)
	}

	if err := api.RegisterOrchestratorServiceHandlerFromEndpoint(ctx, mux, "localhost:50052", opts); err != nil {
		log.Fatalf("Ошибка регистрации Orchestrator-сервиса: %v", err)
	}

	log.Println("HTTP Gateway запущен на :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Ошибка запуска HTTP сервера: %v", err)
	}
}
