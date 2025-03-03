package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Maraei/calculator-on-go/internal/agent"
)

func main() {
	// Количество горутин для вычислений задается через переменную окружения
	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil || computingPower <= 0 {
		computingPower = 2 // Значение по умолчанию
	}

	log.Printf("Агент запущен с %d вычислителями", computingPower)

	// Запускаем агента с указанным числом вычислителей (горутин)
	agent.Start(computingPower)

	// Держим процесс живым, чтобы агент продолжал работу
	for {
		log.Println("Агент продолжает работать...")
		time.Sleep(10 * time.Second)
	}
}
