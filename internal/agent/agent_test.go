package agent

import (
	"testing"

)

func TestCalculate(t *testing.T) {
	tests := []struct {
		arg1       float64
		arg2       float64
		operation  string
		expected   float64
		expectErr  bool
	}{
		{2, 3, "+", 5, false},
		{5, 2, "-", 3, false},
		{4, 3, "*", 12, false},
		{10, 2, "/", 5, false},
		{10, 0, "/", 0, true},
		{1, 1, "%", 0, true},
	}

	for _, test := range tests {
		result, err := Calculate(test.arg1, test.arg2, test.operation)
		if (err != nil) != test.expectErr {
			t.Errorf("Calculate(%v, %v, %v) ожидал ошибку: %v, но получил %v", 
				test.arg1, test.arg2, test.operation, test.expectErr, err)
		}
		if !test.expectErr && result != test.expected {
			t.Errorf("Calculate(%v, %v, %v) ожидал %v, но получил %v", 
				test.arg1, test.arg2, test.operation, test.expected, result)
		}
	}
}

