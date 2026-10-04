package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/valeriatocarruncho/calculator-backend/internal/calculator"
	"github.com/valeriatocarruncho/calculator-backend/internal/handler"
)

func setupTestServer() *http.ServeMux {
	svc := calculator.NewService()
	h := handler.NewCalculatorHandler(svc)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}

func TestHealthCheck(t *testing.T) {
	mux := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got %v", body["status"])
	}
}

func TestCalculateEndpoint_Success(t *testing.T) {
	mux := setupTestServer()

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedResult float64
		expectedOp     string
	}{
		{
			name:           "add operation",
			payload:        `{"operation": "add", "a": 15.5, "b": 4.5}`,
			expectedStatus: http.StatusOK,
			expectedResult: 20.0,
			expectedOp:     "add",
		},
		{
			name:           "subtract operation with symbol",
			payload:        `{"operation": "-", "a": 10, "b": 3}`,
			expectedStatus: http.StatusOK,
			expectedResult: 7.0,
			expectedOp:     "subtract",
		},
		{
			name:           "multiply operation",
			payload:        `{"operation": "multiply", "a": 6, "b": 7}`,
			expectedStatus: http.StatusOK,
			expectedResult: 42.0,
			expectedOp:     "multiply",
		},
		{
			name:           "divide operation",
			payload:        `{"operation": "/", "a": 100, "b": 4}`,
			expectedStatus: http.StatusOK,
			expectedResult: 25.0,
			expectedOp:     "divide",
		},
		{
			name:           "power operation",
			payload:        `{"operation": "^", "a": 2, "b": 8}`,
			expectedStatus: http.StatusOK,
			expectedResult: 256.0,
			expectedOp:     "power",
		},
		{
			name:           "sqrt operation",
			payload:        `{"operation": "sqrt", "a": 49}`,
			expectedStatus: http.StatusOK,
			expectedResult: 7.0,
			expectedOp:     "sqrt",
		},
		{
			name:           "unary percentage",
			payload:        `{"operation": "%", "a": 15}`,
			expectedStatus: http.StatusOK,
			expectedResult: 0.15,
			expectedOp:     "percentage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d. Body: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}

			var resp handler.CalculationResponse
			if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resp.Result != tt.expectedResult {
				t.Errorf("expected result %v, got %v", tt.expectedResult, resp.Result)
			}
			if resp.Operation != tt.expectedOp {
				t.Errorf("expected operation %s, got %s", tt.expectedOp, resp.Operation)
			}
		})
	}
}

func TestCalculateEndpoint_Errors(t *testing.T) {
	mux := setupTestServer()

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "division by zero",
			payload:        `{"operation": "divide", "a": 10, "b": 0}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "DIVISION_BY_ZERO",
		},
		{
			name:           "negative square root",
			payload:        `{"operation": "sqrt", "a": -16}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "NEGATIVE_SQUARE_ROOT",
		},
		{
			name:           "missing operand a",
			payload:        `{"operation": "add", "b": 10}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_OPERAND",
		},
		{
			name:           "missing operand b for binary operation",
			payload:        `{"operation": "add", "a": 10}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "MISSING_OPERAND",
		},
		{
			name:           "unknown operation",
			payload:        `{"operation": "unknown", "a": 10, "b": 2}`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_OPERATION",
		},
		{
			name:           "malformed JSON",
			payload:        `{"operation": "add", a: 10`,
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "INVALID_PAYLOAD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d. Body: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}

			var errResp handler.ErrorResponse
			if err := json.Unmarshal(rr.Body.Bytes(), &errResp); err != nil {
				t.Fatalf("failed to parse error response: %v", err)
			}

			if errResp.Code != tt.expectedCode {
				t.Errorf("expected code %s, got %s", tt.expectedCode, errResp.Code)
			}
		})
	}
}

func TestOperationsSubEndpoints(t *testing.T) {
	mux := setupTestServer()

	t.Run("POST /api/v1/operations/add success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/operations/add", bytes.NewBufferString(`{"a": 20, "b": 30}`))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", rr.Code)
		}

		var resp handler.CalculationResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp.Result != 50 {
			t.Errorf("expected result 50, got %v", resp.Result)
		}
	})

	t.Run("POST /api/v1/operations/divide by zero", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/operations/divide", bytes.NewBufferString(`{"a": 20, "b": 0}`))
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", rr.Code)
		}

		var errResp handler.ErrorResponse
		_ = json.Unmarshal(rr.Body.Bytes(), &errResp)
		if errResp.Code != "DIVISION_BY_ZERO" {
			t.Errorf("expected DIVISION_BY_ZERO, got %s", errResp.Code)
		}
	})
}
