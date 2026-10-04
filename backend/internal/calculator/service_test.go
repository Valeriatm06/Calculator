package calculator_test

import (
	"math"
	"testing"

	"github.com/valeriatocarruncho/calculator-backend/internal/calculator"
)

func ptr(v float64) *float64 {
	return &v
}

func TestCalculatorService_Add(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive integers", 5, 3, 8},
		{"positive and negative", 10, -4, 6},
		{"two negatives", -7, -3, -10},
		{"floats precision", 0.1, 0.2, 0.3},
		{"zero value", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.Add(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, res)
			}
		})
	}
}

func TestCalculatorService_Subtract(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"basic subtraction", 10, 4, 6},
		{"resulting negative", 4, 10, -6},
		{"subtract negative", 5, -5, 10},
		{"floats precision", 0.3, 0.1, 0.2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.Subtract(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, res)
			}
		})
	}
}

func TestCalculatorService_Multiply(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 6, 7, 42},
		{"by zero", 99, 0, 0},
		{"negative numbers", -4, -5, 20},
		{"positive and negative", -4, 5, -20},
		{"decimals", 2.5, 4, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.Multiply(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, res)
			}
		})
	}
}

func TestCalculatorService_Divide(t *testing.T) {
	svc := calculator.NewService()

	t.Run("valid divisions", func(t *testing.T) {
		tests := []struct {
			name     string
			a, b     float64
			expected float64
		}{
			{"standard division", 20, 4, 5},
			{"floating result", 7, 2, 3.5},
			{"negative divisor", 15, -3, -5},
			{"zero numerator", 0, 5, 0},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res, err := svc.Divide(tt.a, tt.b)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if res != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, res)
				}
			})
		}
	})

	t.Run("division by zero", func(t *testing.T) {
		_, err := svc.Divide(10, 0)
		if err != calculator.ErrDivisionByZero {
			t.Fatalf("expected ErrDivisionByZero, got %v", err)
		}
	})
}

func TestCalculatorService_Power(t *testing.T) {
	svc := calculator.NewService()

	tests := []struct {
		name     string
		base, exp float64
		expected float64
	}{
		{"square", 3, 2, 9},
		{"cube", 2, 3, 8},
		{"power of zero", 5, 0, 1},
		{"zero to power", 0, 3, 0},
		{"negative exponent", 2, -1, 0.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := svc.Power(tt.base, tt.exp)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, res)
			}
		})
	}
}

func TestCalculatorService_SquareRoot(t *testing.T) {
	svc := calculator.NewService()

	t.Run("valid square roots", func(t *testing.T) {
		tests := []struct {
			name     string
			a        float64
			expected float64
		}{
			{"perfect square", 16, 4},
			{"zero", 0, 0},
			{"decimal square", 6.25, 2.5},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res, err := svc.SquareRoot(tt.a)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if res != tt.expected {
					t.Errorf("expected %v, got %v", tt.expected, res)
				}
			})
		}
	})

	t.Run("negative square root error", func(t *testing.T) {
		_, err := svc.SquareRoot(-9)
		if err != calculator.ErrNegativeSquareRoot {
			t.Fatalf("expected ErrNegativeSquareRoot, got %v", err)
		}
	})
}

func TestCalculatorService_Percentage(t *testing.T) {
	svc := calculator.NewService()

	t.Run("unary percentage", func(t *testing.T) {
		res, err := svc.Percentage(75, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != 0.75 {
			t.Errorf("expected 0.75, got %v", res)
		}
	})

	t.Run("binary percentage", func(t *testing.T) {
		// 20% of 250 -> 250 * 20 / 100 = 50
		res, err := svc.Percentage(250, ptr(20))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != 50 {
			t.Errorf("expected 50, got %v", res)
		}
	})
}

func TestCalculatorService_Calculate(t *testing.T) {
	svc := calculator.NewService()

	t.Run("valid dispatch operations", func(t *testing.T) {
		res, err := svc.Calculate(calculator.OpAdd, 10, ptr(5))
		if err != nil || res != 15 {
			t.Fatalf("Add failed: res=%v, err=%v", res, err)
		}

		res, err = svc.Calculate(calculator.OpSubtract, 10, ptr(5))
		if err != nil || res != 5 {
			t.Fatalf("Subtract failed: res=%v, err=%v", res, err)
		}

		res, err = svc.Calculate(calculator.OpMultiply, 10, ptr(5))
		if err != nil || res != 50 {
			t.Fatalf("Multiply failed: res=%v, err=%v", res, err)
		}

		res, err = svc.Calculate(calculator.OpDivide, 10, ptr(5))
		if err != nil || res != 2 {
			t.Fatalf("Divide failed: res=%v, err=%v", res, err)
		}

		res, err = svc.Calculate(calculator.OpPower, 2, ptr(3))
		if err != nil || res != 8 {
			t.Fatalf("Power failed: res=%v, err=%v", res, err)
		}

		res, err = svc.Calculate(calculator.OpSquareRoot, 25, nil)
		if err != nil || res != 5 {
			t.Fatalf("SquareRoot failed: res=%v, err=%v", res, err)
		}

		res, err = svc.Calculate(calculator.OpPercentage, 50, nil)
		if err != nil || res != 0.5 {
			t.Fatalf("Percentage failed: res=%v, err=%v", res, err)
		}
	})

	t.Run("missing operand errors", func(t *testing.T) {
		ops := []calculator.Operation{
			calculator.OpAdd,
			calculator.OpSubtract,
			calculator.OpMultiply,
			calculator.OpDivide,
			calculator.OpPower,
		}

		for _, op := range ops {
			_, err := svc.Calculate(op, 10, nil)
			if err != calculator.ErrMissingOperand {
				t.Errorf("op %s: expected ErrMissingOperand, got %v", op, err)
			}
		}
	})

	t.Run("invalid operation", func(t *testing.T) {
		_, err := svc.Calculate("invalid_op", 10, ptr(5))
		if err != calculator.ErrInvalidOperation {
			t.Fatalf("expected ErrInvalidOperation, got %v", err)
		}
	})
}

func TestNormalizeOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected calculator.Operation
		hasErr   bool
	}{
		{"+", calculator.OpAdd, false},
		{"ADD", calculator.OpAdd, false},
		{"-", calculator.OpSubtract, false},
		{"subtract", calculator.OpSubtract, false},
		{"*", calculator.OpMultiply, false},
		{"x", calculator.OpMultiply, false},
		{"/", calculator.OpDivide, false},
		{"^", calculator.OpPower, false},
		{"sqrt", calculator.OpSquareRoot, false},
		{"√", calculator.OpSquareRoot, false},
		{"%", calculator.OpPercentage, false},
		{"unknown", "", true},
	}

	for _, tt := range tests {
		got, err := calculator.NormalizeOperation(tt.input)
		if (err != nil) != tt.hasErr {
			t.Errorf("input %s: unexpected error status: %v", tt.input, err)
		}
		if got != tt.expected {
			t.Errorf("input %s: expected %s, got %s", tt.input, tt.expected, got)
		}
	}
}

func TestSanitizeFloat_Overflow(t *testing.T) {
	svc := calculator.NewService()
	// Large exponentiation causing overflow
	_, err := svc.Power(1e200, 1e200)
	if err != calculator.ErrInvalidResult && !math.IsInf(1e200, 0) {
		t.Fatalf("expected ErrInvalidResult on overflow, got %v", err)
	}
}
