package models

// Expression представляет арифметическое выражение, отправленное пользователем
type Expression struct {
	ID     string   `json:"id"`
	Input  string   `json:"expression"`
	Status string   `json:"status"`  // pending, processing, completed
	Result *float64 `json:"result,omitempty"`
}

// Task представляет отдельную вычислительную задачу для агента
type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}
