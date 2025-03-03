package agent

import (
	"errors"
	"os"
	"strconv"
	"time"
)

// Выполняет операцию в зависимости от типа
func Calculate(arg1, arg2 float64, operation string) (float64, error) {
	// Время выполнения операции задается переменными окружения
	var duration int
	switch operation {
	case "+":
		duration = getOperationTime("TIME_ADDITION_MS", 2000)
	case "-":
		duration = getOperationTime("TIME_SUBTRACTION_MS", 2000)
	case "*":
		duration = getOperationTime("TIME_MULTIPLICATIONS_MS", 3000)
	case "/":
		duration = getOperationTime("TIME_DIVISIONS_MS", 3000)
	default:
		return 0, errors.New("неизвестная операция: " + operation)
	}

	// Симуляция длительного вычисления
	time.Sleep(time.Duration(duration) * time.Millisecond)

	// Выполнение самой операции
	switch operation {
	case "+":
		return arg1 + arg2, nil
	case "-":
		return arg1 - arg2, nil
	case "*":
		return arg1 * arg2, nil
	case "/":
		if arg2 == 0 {
			return 0, errors.New("деление на ноль")
		}
		return arg1 / arg2, nil
	}
	return 0, errors.New("неизвестная операция: " + operation)
}

// Получение времени выполнения операции из переменных окружения
func getOperationTime(envVar string, defaultTime int) int {
	value, err := strconv.Atoi(os.Getenv(envVar))
	if err != nil || value <= 0 {
		return defaultTime
	}
	return value
}
