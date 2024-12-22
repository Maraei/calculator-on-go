package calculation

import (
	"errors"
	"testing"
)

func TestCalc(t *testing.T) {
	testCases := []struct {
		name           string
		expression     string
		expectedResult float64
		err            error
	}{
		{
			name:           "обычный случай",
			expression:     "1+1",
			expectedResult: 2,
			err:            nil,
		},
		{
			name:           "с возведением в степень",
			expression:     "2^3",
			expectedResult: 8,
			err:            nil,
		},
		{
			name:           "возведение в степень с отрицательным числом",
			expression:     "-2^3",
			expectedResult: -8,
			err:            nil,
		},
		{
			name:           "смешанное выражение с возведением в степень",
			expression:     "2+2^3",
			expectedResult: 10,
			err:            nil,
		},
		{
			name:           "возведение в степень с дробным числом",
			expression:     "4^0.5",
			expectedResult: 2,
			err:            nil,
		},
		{
			name:           "деление на ноль",
			expression:     "5/0",
			expectedResult: 0,
			err:            ErrDivisionByZero,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := Calc(testCase.expression)
			if err != nil && !errors.Is(err, testCase.err) {
				t.Fatalf("expected error %v, got %v", testCase.err, err)
			}

			if err == nil && result != testCase.expectedResult {
				t.Fatalf("expected result %f, got %f", testCase.expectedResult, result)
			}
		})
	}
}


func TestCalcErrors(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
		expectedErr error
	}{
		{
			name:        "деление на ноль",
			expression:  "10/0",
			expectedErr: ErrDivisionByZero,
		},
		{
			name:        "неверный символ",
			expression:  "not numbs",
			expectedErr: ErrInvalidExpression,
		},
		{
			name:        "неверные символы в выражении",
			expression:  "2r+10b",
			expectedErr: ErrInvalidExpression,
		},

		{
			name:        "пустое выражение",
			expression:  "",
			expectedErr: ErrEmptyInput,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := Calc(testCase.expression)
			if err == nil {
				t.Fatalf("expected error for expression %s, but no error occurred", testCase.expression)
			}

			// Проверяем только тип ошибки
			if !errors.Is(err, testCase.expectedErr) {
				t.Fatalf("expected error %v, but got %v", testCase.expectedErr, err)
			}
		})
	}
}