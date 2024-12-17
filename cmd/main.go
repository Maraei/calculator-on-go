package main

import (
	"github.com/Maraei/calculator-on-go/internal/application"
)

func main() {
	app := application.New()
	// app.Run()
	app.RunServer()
}