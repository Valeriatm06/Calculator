package calculator

import "errors"

var (
	// ErrDivisionByZero occurs when attempting to divide by zero.
	ErrDivisionByZero = errors.New("division by zero is undefined")

	// ErrNegativeSquareRoot occurs when taking square root of a negative number.
	ErrNegativeSquareRoot = errors.New("square root of negative number is undefined for real numbers")

	// ErrInvalidOperation occurs when an unsupported operation is requested.
	ErrInvalidOperation = errors.New("invalid or unsupported operation")

	// ErrMissingOperand occurs when a required operand is not provided.
	ErrMissingOperand = errors.New("missing required operand")

	// ErrInvalidResult occurs when the calculation results in NaN or Infinity.
	ErrInvalidResult = errors.New("calculation resulted in an invalid number or overflow")
)
