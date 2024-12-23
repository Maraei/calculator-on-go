package application

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Maraei/calculator-on-go/pkg/calculation"
)

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	// Проверка метода запроса
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"Invalid Body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var request struct {
		Expression string `json:"expression"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil || request.Expression == "" {
		http.Error(w, `{"error":"Invalid Body"}`, http.StatusBadRequest)
		return
	}

	// Выполнение вычислений
	result, err := calculation.Calc(request.Expression)
	if err != nil {
		var errorMsg string
		statusCode := http.StatusUnprocessableEntity

		// Обработка ошибок
		switch err {
		case calculation.ErrInvalidExpression:
			errorMsg = "Error calculation"
		case calculation.ErrDivisionByZero:
			errorMsg = "Division by zero"
		case calculation.ErrMismatchedParentheses:
			errorMsg = "Mismatched parentheses"
		case calculation.ErrInvalidNumber:
			errorMsg = "Invalid number"
		case calculation.ErrUnexpectedToken:
			errorMsg = "Unexpected token"
		case calculation.ErrNotEnoughValues:
			errorMsg = "Not enough values in expression"
		case calculation.ErrInvalidOperator:
			errorMsg = "Invalid operator"
		case calculation.ErrOperatorAtEnd:
			errorMsg = "Operator at the end"
		case calculation.ErrMultipleDecimalPoints:
			errorMsg = "Multiple decimal points"
		case calculation.ErrEmptyInput:
			errorMsg = "Empty expression"
		default:
			errorMsg = "Error calculation"
			statusCode = http.StatusUnprocessableEntity
		}

		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, errorMsg), statusCode)
		return
	}

	// Формирование успешного ответа
	response := struct {
		Result string `json:"result"`
	}{
		Result: fmt.Sprintf("%v", result),
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		http.Error(w, `{"error":"Unknown error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJson)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
