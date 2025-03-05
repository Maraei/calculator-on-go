package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Maraei/calculator-on-go/internal/agent"
)

func main() {
	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil || computingPower <= 0 {
		computingPower = 2
	}

	log.Printf("Агент запущен с %d вычислителями", computingPower)

	agent.Start(computingPower)

	for {
		log.Println("Агент продолжает работать...")
		time.Sleep(10 * time.Second)
	}
}
