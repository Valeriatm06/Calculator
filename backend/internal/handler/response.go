package handler

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse defines standard JSON wrapper for successful requests.
type SuccessResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
}

// ErrorResponse defines standard JSON wrapper for error responses.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// CalculationResponse is the payload returned for calculator operations.
type CalculationResponse struct {
	Operation string   `json:"operation"`
	A         float64  `json:"a"`
	B         *float64 `json:"b,omitempty"`
	Result    float64  `json:"result"`
	Formatted string   `json:"formatted"`
}

// JSON writes a JSON response with status code.
func JSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if payload != nil {
		_ = json.NewEncoder(w).Encode(payload)
	}
}

// Error sends a formatted JSON error.
func Error(w http.ResponseWriter, status int, message, code string) {
	JSON(w, status, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
	})
}
