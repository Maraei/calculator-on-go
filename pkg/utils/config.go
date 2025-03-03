package utils

import (
	"os"
	"strconv"
)

// Config хранит настройки приложения
type Config struct {
	ServerPort         string
	OrchestratorURL    string
	ComputingPower     int
	TimeAdditionMS     int
	TimeSubtractionMS  int
	TimeMultiplicationMS int
	TimeDivisionMS     int
}

// LoadConfig загружает переменные окружения в конфиг
func LoadConfig() Config {
	return Config{
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		OrchestratorURL:    getEnv("ORCHESTRATOR_URL", "http://localhost"),
		ComputingPower:     getEnvAsInt("COMPUTING_POWER", 2),
		TimeAdditionMS:     getEnvAsInt("TIME_ADDITION_MS", 2000),
		TimeSubtractionMS:  getEnvAsInt("TIME_SUBTRACTION_MS", 2000),
		TimeMultiplicationMS: getEnvAsInt("TIME_MULTIPLICATIONS_MS", 3000),
		TimeDivisionMS:     getEnvAsInt("TIME_DIVISIONS_MS", 3000),
	}
}

// getEnv получает строковую переменную окружения или использует значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает числовую переменную окружения или использует значение по умолчанию
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
