package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/valeriatocarruncho/calculator-backend/internal/calculator"
)

// CalculateRequest models incoming calculation payloads.
type CalculateRequest struct {
	Operation string   `json:"operation"`
	A         *float64 `json:"a"`
	B         *float64 `json:"b,omitempty"`
}

// BinaryOperationRequest models endpoints where operation is inferred from URL.
type BinaryOperationRequest struct {
	A *float64 `json:"a"`
	B *float64 `json:"b"`
}

// UnaryOperationRequest models single-operand requests (like sqrt, percentage).
type UnaryOperationRequest struct {
	A *float64 `json:"a"`
	B *float64 `json:"b,omitempty"`
}

// CalculatorHandler handles HTTP requests for calculator operations.
type CalculatorHandler struct {
	service calculator.Service
}

// NewCalculatorHandler creates a new CalculatorHandler.
func NewCalculatorHandler(service calculator.Service) *CalculatorHandler {
	return &CalculatorHandler{service: service}
}

// formatResult builds a human-readable mathematical equation string.
func formatResult(op calculator.Operation, a float64, b *float64, res float64) string {
	switch op {
	case calculator.OpAdd:
		return fmt.Sprintf("%v + %v = %v", a, *b, res)
	case calculator.OpSubtract:
		return fmt.Sprintf("%v - %v = %v", a, *b, res)
	case calculator.OpMultiply:
		return fmt.Sprintf("%v × %v = %v", a, *b, res)
	case calculator.OpDivide:
		return fmt.Sprintf("%v ÷ %v = %v", a, *b, res)
	case calculator.OpPower:
		return fmt.Sprintf("%v ^ %v = %v", a, *b, res)
	case calculator.OpSquareRoot:
		return fmt.Sprintf("√%v = %v", a, res)
	case calculator.OpPercentage:
		if b != nil {
			return fmt.Sprintf("%v%% of %v = %v", *b, a, res)
		}
		return fmt.Sprintf("%v%% = %v", a, res)
	default:
		return fmt.Sprintf("%v", res)
	}
}

// mapDomainError converts domain errors to appropriate HTTP status codes and response messages.
func mapDomainError(err error) (int, string, string) {
	switch {
	case errors.Is(err, calculator.ErrDivisionByZero):
		return http.StatusBadRequest, err.Error(), "DIVISION_BY_ZERO"
	case errors.Is(err, calculator.ErrNegativeSquareRoot):
		return http.StatusBadRequest, err.Error(), "NEGATIVE_SQUARE_ROOT"
	case errors.Is(err, calculator.ErrMissingOperand):
		return http.StatusBadRequest, err.Error(), "MISSING_OPERAND"
	case errors.Is(err, calculator.ErrInvalidOperation):
		return http.StatusBadRequest, err.Error(), "INVALID_OPERATION"
	case errors.Is(err, calculator.ErrInvalidResult):
		return http.StatusUnprocessableEntity, err.Error(), "NUMERIC_OVERFLOW"
	default:
		return http.StatusInternalServerError, "an unexpected error occurred", "INTERNAL_ERROR"
	}
}

// Calculate handles the unified POST /api/v1/calculate endpoint.
func (h *CalculatorHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		Error(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED")
		return
	}

	var req CalculateRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, fmt.Sprintf("invalid JSON payload: %v", err), "INVALID_PAYLOAD")
		return
	}

	if req.A == nil {
		Error(w, http.StatusBadRequest, "operand 'a' is required", "MISSING_OPERAND")
		return
	}

	op, err := calculator.NormalizeOperation(req.Operation)
	if err != nil {
		Error(w, http.StatusBadRequest, fmt.Sprintf("unsupported operation '%s'", req.Operation), "INVALID_OPERATION")
		return
	}

	result, err := h.service.Calculate(op, *req.A, req.B)
	if err != nil {
		status, msg, code := mapDomainError(err)
		Error(w, status, msg, code)
		return
	}

	JSON(w, http.StatusOK, CalculationResponse{
		Operation: string(op),
		A:         *req.A,
		B:         req.B,
		Result:    result,
		Formatted: formatResult(op, *req.A, req.B, result),
	})
}

// HandleOperation returns an HTTP handler for a specific operation endpoint (e.g. POST /api/v1/operations/add).
func (h *CalculatorHandler) HandleOperation(op calculator.Operation) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			Error(w, http.StatusMethodNotAllowed, "method not allowed", "METHOD_NOT_ALLOWED")
			return
		}

		var req CalculateRequest
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&req); err != nil {
			Error(w, http.StatusBadRequest, fmt.Sprintf("invalid JSON payload: %v", err), "INVALID_PAYLOAD")
			return
		}

		if req.A == nil {
			Error(w, http.StatusBadRequest, "operand 'a' is required", "MISSING_OPERAND")
			return
		}

		result, err := h.service.Calculate(op, *req.A, req.B)
		if err != nil {
			status, msg, code := mapDomainError(err)
			Error(w, status, msg, code)
			return
		}

		JSON(w, http.StatusOK, CalculationResponse{
			Operation: string(op),
			A:         *req.A,
			B:         req.B,
			Result:    result,
			Formatted: formatResult(op, *req.A, req.B, result),
		})
	}
}

// HealthCheck provides liveness/readiness probe.
func (h *CalculatorHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "calculator-api",
		"version": "1.0.0",
	})
}

// RegisterRoutes registers all calculator routes onto an http.ServeMux.
func (h *CalculatorHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/health", h.HealthCheck)
	mux.HandleFunc("POST /api/v1/calculate", h.Calculate)

	// Convenience sub-endpoints
	mux.HandleFunc("POST /api/v1/operations/add", h.HandleOperation(calculator.OpAdd))
	mux.HandleFunc("POST /api/v1/operations/subtract", h.HandleOperation(calculator.OpSubtract))
	mux.HandleFunc("POST /api/v1/operations/multiply", h.HandleOperation(calculator.OpMultiply))
	mux.HandleFunc("POST /api/v1/operations/divide", h.HandleOperation(calculator.OpDivide))
	mux.HandleFunc("POST /api/v1/operations/power", h.HandleOperation(calculator.OpPower))
	mux.HandleFunc("POST /api/v1/operations/sqrt", h.HandleOperation(calculator.OpSquareRoot))
	mux.HandleFunc("POST /api/v1/operations/percentage", h.HandleOperation(calculator.OpPercentage))
}
