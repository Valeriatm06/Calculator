package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

// Index serves either an interactive HTML documentation page or a JSON discovery directory.
func (h *CalculatorHandler) Index(w http.ResponseWriter, r *http.Request) {
	// If requested from a web browser, serve human-friendly documentation
	if strings.Contains(r.Header.Get("Accept"), "text/html") {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Calculator REST API - Documentation</title>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f8fafc; color: #1e293b; margin: 0; padding: 2rem 1rem; }
    .container { max-width: 800px; margin: 0 auto; background: #ffffff; border-radius: 16px; padding: 2rem; box-shadow: 0 10px 25px rgba(0,0,0,0.05); border: 1px solid #e2e8f0; }
    h1 { margin-top: 0; display: flex; align-items: center; gap: 10px; font-size: 1.8rem; }
    .badge { background: #dcfce7; color: #15803d; font-size: 0.8rem; padding: 4px 10px; border-radius: 9999px; font-weight: 600; }
    .endpoint { background: #f1f5f9; border-radius: 8px; padding: 1rem; margin: 1rem 0; border-left: 4px solid #3b82f6; }
    .method { font-weight: 700; color: #2563eb; background: #dbeafe; padding: 2px 8px; border-radius: 4px; font-size: 0.85rem; margin-right: 8px; }
    .method.post { color: #16a34a; background: #dcfce7; }
    code, pre { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 0.9rem; }
    pre { background: #1e293b; color: #f8fafc; padding: 1rem; border-radius: 8px; overflow-x: auto; }
    a { color: #2563eb; text-decoration: none; font-weight: 500; }
    a:hover { text-decoration: underline; }
    .btn { display: inline-block; background: #2563eb; color: white; padding: 8px 16px; border-radius: 8px; text-decoration: none; font-weight: 500; margin-top: 10px; }
    .btn:hover { background: #1d4ed8; text-decoration: none; }
  </style>
</head>
<body>
  <div class="container">
    <h1><span>Calculator REST API</span> <span class="badge">Online (v1.0.0)</span></h1>
    <p>Welcome to the Go Calculator Microservice API. The following endpoints are available:</p>

    <div style="margin: 1.5rem 0 2rem 0;">
      <a href="/swagger" class="btn" style="background: #10b981; margin-right: 8px; font-weight: 600;">⚡ Open Swagger UI</a>
      <a href="/api/v1/health" class="btn" target="_blank">Check /health JSON</a>
    </div>

    <div class="endpoint">
      <div><span class="method post">POST</span> <code>/api/v1/calculate</code></div>
      <p>Primary calculation endpoint. Accepts operations: <code>add</code>, <code>subtract</code>, <code>multiply</code>, <code>divide</code>, <code>power</code>, <code>sqrt</code>, <code>percentage</code>.</p>
      <pre>curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "add", "a": 55, "b": 55}'</pre>
    </div>

    <div class="endpoint">
      <div><span class="method post">POST</span> <code>/api/v1/operations/{op}</code></div>
      <p>Dedicated sub-endpoints: <code>/api/v1/operations/add</code>, <code>/subtract</code>, <code>/multiply</code>, <code>/divide</code>, <code>/power</code>, <code>/sqrt</code>, <code>/percentage</code>.</p>
      <pre>curl -X POST http://localhost:8080/api/v1/operations/add \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "b": 20}'</pre>
    </div>

    <hr style="border: none; border-top: 1px solid #e2e8f0; margin: 2rem 0;">
    <p style="font-size: 0.85rem; color: #64748b;">Frontend UI is running at <a href="http://localhost:3000">http://localhost:3000</a></p>
  </div>
</body>
</html>`))
		return
	}

	// JSON response for API clients
	JSON(w, http.StatusOK, map[string]interface{}{
		"service": "calculator-api",
		"version": "1.0.0",
		"status":  "healthy",
		"endpoints": map[string]string{
			"swagger_ui":  "GET /swagger",
			"swagger_spec": "GET /swagger.json",
			"health":      "GET /api/v1/health",
			"calculate":   "POST /api/v1/calculate",
			"operations":  "POST /api/v1/operations/{add,subtract,multiply,divide,power,sqrt,percentage}",
			"frontend_ui": "http://localhost:3000",
		},
	})
}

// RegisterRoutes registers all calculator routes onto an http.ServeMux.
func (h *CalculatorHandler) RegisterRoutes(mux *http.ServeMux) {
	// Swagger UI & OpenAPI Specification routes
	mux.HandleFunc("GET /swagger", h.SwaggerUI)
	mux.HandleFunc("GET /swagger/", h.SwaggerUI)
	mux.HandleFunc("GET /docs", h.SwaggerUI)
	mux.HandleFunc("GET /swagger.json", h.SwaggerJSON)
	mux.HandleFunc("GET /api/v1/swagger.json", h.SwaggerJSON)

	// Root and API index documentation routes
	mux.HandleFunc("GET /{$}", h.Index)
	mux.HandleFunc("GET /api", h.Index)
	mux.HandleFunc("GET /api/", h.Index)
	mux.HandleFunc("GET /api/v1", h.Index)
	mux.HandleFunc("GET /api/v1/", h.Index)

	// Health and calculation endpoints
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
