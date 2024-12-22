package main

import (
	"log"
	"net/http"
	"github.com/Maraei/calculator-on-go/internal/application"
)

func main() {
	http.HandleFunc("/api/v1/calculate", application.CalcHandler)
	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
} 
