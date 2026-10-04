package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/valeriatocarruncho/calculator-backend/internal/middleware"
)

func TestCORS(t *testing.T) {
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	corsMiddleware := middleware.CORS(dummyHandler)

	t.Run("OPTIONS preflight returns 204", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/api/v1/calculate", nil)
		rr := httptest.NewRecorder()

		corsMiddleware.ServeHTTP(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Fatalf("expected 204 No Content for OPTIONS, got %d", rr.Code)
		}

		if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("missing or incorrect Access-Control-Allow-Origin header")
		}
	})

	t.Run("regular request forwards to handler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
		rr := httptest.NewRecorder()

		corsMiddleware.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d", rr.Code)
		}
		if rr.Body.String() != "ok" {
			t.Fatalf("expected body 'ok', got '%s'", rr.Body.String())
		}
	})
}

func TestLogger(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})

	loggedHandler := middleware.Logger(handler)
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	rr := httptest.NewRecorder()

	loggedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", rr.Code)
	}
}
