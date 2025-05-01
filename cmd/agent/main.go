package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/Maraei/calculator-on-go/internal/agent"
	"google.golang.org/grpc"
	"github.com/Maraei/calculator-on-go/api/api"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil || computingPower <= 0 {
		computingPower = 2
	}

	log.Printf("Агент запущен с %d вычислителями", computingPower)

	// Запуск агента без авторизации
	if err := agent.Start(computingPower); err != nil {
		log.Fatalf("Ошибка запуска агента: %v", err)
	}

	// Ожидание команды на авторизацию
	log.Println("Агент работает. Введите команду для авторизации...")
	for {
		var input string
		fmt.Scanln(&input)

		if input == "login" {
			// Вызов авторизации по запросу
			if err := authorizeAgent(); err != nil {
				log.Println("Ошибка авторизации:", err)
			} else {
				log.Println("Агент успешно авторизован")
			}
		} else {
			log.Println("Неизвестная команда. Для авторизации введите 'login'.")
		}

		// Даем возможность пользователю снова ввести команду
	}
}

// Функция авторизации по запросу
func authorizeAgent() error {
	// Подключаемся к Auth-серверу
	authConn, err := grpc.Dial(agent.GetAuthServerAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer authConn.Close()

	authClient := api.NewAuthCalculatorServiceClient(authConn)

	// Авторизация агента
	if err := agent.Login(authClient); err != nil {
		return err
	}
	log.Println("Агент успешно авторизован")
	return nil
}
