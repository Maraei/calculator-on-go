package main

import (
	"fmt"
	"github.com/Maraei/calculator-on-go/internal/application"
)

func main() {
	app := application.New()
	fmt.Println("RunServer")
	// app.Run()
	app.RunServer()
}