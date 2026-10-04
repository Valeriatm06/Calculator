package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/valeriatocarruncho/calculator-backend/internal/calculator"
	"github.com/valeriatocarruncho/calculator-backend/internal/handler"
	"github.com/valeriatocarruncho/calculator-backend/internal/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize domain service & handlers
	calcService := calculator.NewService()
	calcHandler := handler.NewCalculatorHandler(calcService)

	mux := http.NewServeMux()
	calcHandler.RegisterRoutes(mux)

	// Wrap router with CORS and Logger middlewares
	rootHandler := middleware.Logger(middleware.CORS(mux))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      rootHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Server run context for graceful shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for OS signals to stop server
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		log.Println("Shutting down calculator API service gracefully...")

		// Shutdown signal with grace period of 15 seconds
		shutdownCtx, shutdownCancel := context.WithTimeout(serverCtx, 15*time.Second)
		defer shutdownCancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("Graceful shutdown timed out.. forcing exit.")
			}
		}()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}
		serverStopCtx()
	}()

	log.Printf("🚀 Calculator API server started on port %s", port)
	log.Printf("Health check: http://localhost:%s/api/v1/health", port)
	log.Printf("Calculate endpoint: http://localhost:%s/api/v1/calculate", port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	log.Println("Calculator API service stopped.")
}
