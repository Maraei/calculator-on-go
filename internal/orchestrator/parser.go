package orchestrator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var precedence = map[string]int{
	"+": 1,
	"-": 1,
	"*": 2,
	"/": 2,
}

func isOperatorToken(token string) bool {
	_, ok := precedence[token]
	return ok
}

func isNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

func tokenize(expr string) ([]string, error) {
	var tokens []string
	var current strings.Builder

	for _, r := range expr {
		switch {
		case unicode.IsDigit(r) || r == '.':
			current.WriteRune(r)

		case unicode.IsSpace(r):
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}

		case strings.ContainsRune("+-*/()", r):
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(r))

		default:
			return nil, fmt.Errorf("invalid character: %c", r)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens, nil
}

func InfixToRPN(expr string) ([]string, error) {
	tokens, err := tokenize(expr)
	if err != nil {
		return nil, err
	}

	var output []string
	var stack []string

	for _, token := range tokens {
		switch {
		case isNumber(token):
			output = append(output, token)

		case isOperatorToken(token):
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if isOperatorToken(top) && precedence[top] >= precedence[token] {
					output = append(output, top)
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			stack = append(stack, token)

		case token == "(":
			stack = append(stack, token)

		case token == ")":
			foundLeftParen := false
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top == "(" {
					foundLeftParen = true
					break
				}
				output = append(output, top)
			}
			if !foundLeftParen {
				return nil, fmt.Errorf("mismatched parentheses")
			}

		default:
			return nil, fmt.Errorf("unknown token: %s", token)
		}
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if top == "(" || top == ")" {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		output = append(output, top)
	}

	return output, nil
}

func EvaluateRPN(tokens []string) (float64, error) {
	var stack []float64

	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			if len(stack) < 2 {
				return 0, errors.New("invalid expression")
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var res float64
			switch token {
			case "+":
				res = a + b
			case "-":
				res = a - b
			case "*":
				res = a * b
			case "/":
				if b == 0 {
					return 0, errors.New("division by zero")
				}
				res = a / b
			}
			stack = append(stack, res)

		default:
			val, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, val)
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("invalid expression")
	}

	return stack[0], nil
}
