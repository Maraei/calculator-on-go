package agent

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"errors"
)

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

func Start(workerCount int) {
	for i := 0; i < workerCount; i++ {
		go worker(i)
	}
}

func worker(id int) {
	serverURL := os.Getenv("ORCHESTRATOR_URL")
	if serverURL == "" {
		serverURL = "http://localhost"
		log.Println("ORCHESTRATOR_URL не задан, используется значение по умолчанию:", serverURL)
	}
	taskEndpoint := serverURL + "/internal/task"
	
	for {
		task, err := fetchTask(taskEndpoint)
		if err != nil {
			log.Printf("[Worker %d] Ошибка при получении задачи: %v", id, err)
			time.Sleep(2 * time.Second)
			continue
		}

		if task == nil {
			log.Printf("[Worker %d] Нет новых задач, ждем...", id)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Printf("[Worker %d] Получена задача: %v", id, task)

		result, err := Calculate(task.Arg1, task.Arg2, task.Operation)
		if err != nil {
			if err := sendResult(taskEndpoint, task.ID, 0, err.Error()); err != nil {
				log.Printf("[Worker %d] Ошибка отправки ошибки: %v", id, err)
			} else {
				log.Printf("[Worker %d] Ошибка вычисления: %v", id, err)
			}
			continue
		}

		if err := sendResult(taskEndpoint, task.ID, result, ""); err != nil {
			log.Printf("[Worker %d] Ошибка отправки результата: %v", id, err)
		} else {
			log.Printf("[Worker %d] Успешно отправлен результат: %v", id, result)
		}
	}
}

func fetchTask(url string) (*Task, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("не удалось получить задачу: статус " + resp.Status)
	}

	var response struct {
		Task Task `json:"task"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response.Task, nil
}

func sendResult(url, taskID string, result float64, errMsg string) error {
	payload := make(map[string]interface{})
	payload["id"] = taskID
	if errMsg != "" {
		payload["error"] = errMsg
	} else {
		payload["result"] = result
	}

	jsonPayload, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("не удалось отправить результат: статус " + resp.Status)
	}
	return nil
}