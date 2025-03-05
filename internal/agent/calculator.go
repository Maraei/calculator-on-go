package agent

import (
	"os"
	"strconv"
	"fmt"
)

func Calculate(arg1, arg2 float64, operation string) (float64, error) {
    switch operation {
    case "+":
        return arg1 + arg2, nil
    case "-":
        return arg1 - arg2, nil
    case "*":
        return arg1 * arg2, nil
    case "/":
        if arg2 == 0 {
            return 0, fmt.Errorf("деление на ноль")
        }
        return arg1 / arg2, nil
    default:
        return 0, fmt.Errorf("неизвестная операция: %s", operation)
    }
}

func getOperationTime(envVar string, defaultTime int) int {
	value, err := strconv.Atoi(os.Getenv(envVar))
	if err != nil || value <= 0 {
		return defaultTime
	}
	return value
}
