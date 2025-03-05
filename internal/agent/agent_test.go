package agent

import (
	"testing"
	"os"
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

func TestGetOperationTime(t *testing.T) {
	os.Setenv("TIME_ADDITION_MS", "3000")
	defer os.Unsetenv("TIME_ADDITION_MS")

	time1 := getOperationTime("TIME_ADDITION_MS", 2000)
	if time1 != 3000 {
		t.Errorf("Ожидали 3000, но получили %v", time1)
	}

	time2 := getOperationTime("UNKNOWN_ENV", 4000)
	if time2 != 4000 {
		t.Errorf("Ожидали 4000, но получили %v", time2)
	}

	time3 := getOperationTime("", 5000)
	if time3 != 5000 {
		t.Errorf("Ожидали 5000, но получили %v", time3)
	}
}