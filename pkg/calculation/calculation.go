package calculation

import (
	"fmt"
	"strconv"
	"strings"
	"math"
)

func Calc(expression string) (float64, error) {
	if expression == "" {
		return 0, ErrEmptyInput
	}
	expression = strings.ReplaceAll(expression, " ", "")
	return evaluateExpression(expression)
}

func evaluateExpression(expression string) (float64, error) {
	var values []float64
	var ops []rune
	i := 0
	decimalPointFound := false

	for i < len(expression) {
		char := rune(expression[i])

		if char == '(' {
			ops = append(ops, char)
			i++
			continue
		} else if char == ')' {
			for len(ops) > 0 && ops[len(ops)-1] != '(' {
				if err := applyOperation(&values, &ops); err != nil {
					return 0, err
				}
			}
			if len(ops) == 0 || ops[len(ops)-1] != '(' {
				return 0, ErrMismatchedParentheses
			}
			ops = ops[:len(ops)-1]
			i++
			continue
		} else if isDigit(char) || (char == '-' && (i == 0 || expression[i-1] == '(' || isOperator(rune(expression[i-1])))) {
			start := i
			if char == '-' {
				i++
			}
			for i < len(expression) && (isDigit(rune(expression[i])) || expression[i] == '.') {
				if expression[i] == '.' {
					if decimalPointFound {
						return 0, ErrMultipleDecimalPoints
					}
					decimalPointFound = true
				}
				i++
			}
			num, err := strconv.ParseFloat(expression[start:i], 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number: %s: %w", expression[start:i], ErrInvalidNumber)
			}
			values = append(values, num)
			decimalPointFound = false
			continue
		} else if isOperator(char) {
			if i == len(expression)-1 {
				return 0, ErrOperatorAtEnd
			}
			for len(ops) > 0 && precedence(ops[len(ops)-1]) >= precedence(char) {
				if err := applyOperation(&values, &ops); err != nil {
					return 0, err
				}
			}
			ops = append(ops, char)
			i++
			continue
		} else {
			return 0, ErrInvalidExpression
		}
	}

	for len(ops) > 0 {
		if err := applyOperation(&values, &ops); err != nil {
			return 0, err
		}
	}

	if len(values) != 1 {
		return 0, ErrInvalidExpression
	}
	return values[0], nil
}

func applyOperation(values *[]float64, ops *[]rune) error {
	if len(*values) < 2 {
		return ErrNotEnoughValues
	}

	// Извлекаем два последних числа
	val2 := (*values)[len(*values)-1]
	val1 := (*values)[len(*values)-2]

	// Извлекаем последний оператор
	op := (*ops)[len(*ops)-1]

	// Удаляем последние значения и оператор
	*values = (*values)[:len(*values)-2]
	*ops = (*ops)[:len(*ops)-1]

	// Выполнение операции
	var result float64
	switch op {
	case '+':
		result = val1 + val2
	case '-':
		result = val1 - val2
	case '*':
		result = val1 * val2
	case '/':
		if val2 == 0 {
			return ErrDivisionByZero
		}
		result = val1 / val2
	case '^':
		result = math.Pow(val1, val2)
	default:
		return ErrInvalidOperator // Неподдерживаемый оператор
	}

	// Добавляем результат обратно в стек значений
	*values = append(*values, result)
	return nil
}

func precedence(op rune) int {
	switch op {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	case '^':
		return 3  // Возведение в степень имеет более высокий приоритет
	default:
		return 0
	}
}


func isOperator(c rune) bool {
	return c == '+' || c == '-' || c == '*' || c == '/' || c == '^'
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}
