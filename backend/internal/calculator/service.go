package calculator

import (
	"math"
	"strings"
)

// Operation represents supported arithmetic operations.
type Operation string

const (
	OpAdd        Operation = "add"
	OpSubtract   Operation = "subtract"
	OpMultiply   Operation = "multiply"
	OpDivide     Operation = "divide"
	OpPower      Operation = "power"
	OpSquareRoot Operation = "sqrt"
	OpPercentage Operation = "percentage"
)

// Service defines calculator business logic.
type Service interface {
	Add(a, b float64) (float64, error)
	Subtract(a, b float64) (float64, error)
	Multiply(a, b float64) (float64, error)
	Divide(a, b float64) (float64, error)
	Power(base, exponent float64) (float64, error)
	SquareRoot(a float64) (float64, error)
	Percentage(a float64, b *float64) (float64, error)
	Calculate(op Operation, a float64, b *float64) (float64, error)
}

type calculatorService struct{}

// NewService instantiates a new calculator service.
func NewService() Service {
	return &calculatorService{}
}

// sanitizeFloat checks for NaN/Inf and trims micro floating-point representation artifacts.
func sanitizeFloat(val float64) (float64, error) {
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0, ErrInvalidResult
	}
	// Round to 12 decimal places to eliminate standard IEEE-754 float inaccuracies (e.g. 0.1 + 0.2)
	rounded := math.Round(val*1e12) / 1e12
	// Normalize negative zero (-0 -> 0)
	if rounded == 0 {
		return 0, nil
	}
	return rounded, nil
}

func (s *calculatorService) Add(a, b float64) (float64, error) {
	return sanitizeFloat(a + b)
}

func (s *calculatorService) Subtract(a, b float64) (float64, error) {
	return sanitizeFloat(a - b)
}

func (s *calculatorService) Multiply(a, b float64) (float64, error) {
	return sanitizeFloat(a * b)
}

func (s *calculatorService) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return sanitizeFloat(a / b)
}

func (s *calculatorService) Power(base, exponent float64) (float64, error) {
	res := math.Pow(base, exponent)
	return sanitizeFloat(res)
}

func (s *calculatorService) SquareRoot(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeSquareRoot
	}
	return sanitizeFloat(math.Sqrt(a))
}

func (s *calculatorService) Percentage(a float64, b *float64) (float64, error) {
	if b != nil {
		// e.g. 20% of 150 = (150 * 20) / 100
		return sanitizeFloat((a * (*b)) / 100.0)
	}
	// unary percentage: 25% = 0.25
	return sanitizeFloat(a / 100.0)
}

// NormalizeOperation maps common aliases and symbols to standard Operation enum.
func NormalizeOperation(raw string) (Operation, error) {
	cleaned := strings.ToLower(strings.TrimSpace(raw))
	switch cleaned {
	case "add", "+", "plus":
		return OpAdd, nil
	case "sub", "subtract", "-", "minus":
		return OpSubtract, nil
	case "mul", "multiply", "*", "times", "x":
		return OpMultiply, nil
	case "div", "divide", "/", "division":
		return OpDivide, nil
	case "pow", "power", "^", "exponent", "exponentiation":
		return OpPower, nil
	case "sqrt", "squareroot", "root", "√":
		return OpSquareRoot, nil
	case "percent", "percentage", "%":
		return OpPercentage, nil
	default:
		return "", ErrInvalidOperation
	}
}

func (s *calculatorService) Calculate(op Operation, a float64, b *float64) (float64, error) {
	switch op {
	case OpAdd:
		if b == nil {
			return 0, ErrMissingOperand
		}
		return s.Add(a, *b)
	case OpSubtract:
		if b == nil {
			return 0, ErrMissingOperand
		}
		return s.Subtract(a, *b)
	case OpMultiply:
		if b == nil {
			return 0, ErrMissingOperand
		}
		return s.Multiply(a, *b)
	case OpDivide:
		if b == nil {
			return 0, ErrMissingOperand
		}
		return s.Divide(a, *b)
	case OpPower:
		if b == nil {
			return 0, ErrMissingOperand
		}
		return s.Power(a, *b)
	case OpSquareRoot:
		return s.SquareRoot(a)
	case OpPercentage:
		return s.Percentage(a, b)
	default:
		return 0, ErrInvalidOperation
	}
}
