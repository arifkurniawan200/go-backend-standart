package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arifkurniawan200/go-backend-standart/config"
	"github.com/arifkurniawan200/go-backend-standart/internal/handler"
	"github.com/arifkurniawan200/go-backend-standart/internal/repository"
	"github.com/arifkurniawan200/go-backend-standart/internal/usecase"
	"github.com/arifkurniawan200/go-backend-standart/pkg/logger"
)

func main() {
	// Initialize logger
	zapLogger, err := logger.New()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer zapLogger.Sync()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		zapLogger.Fatal("Failed to load configuration", "error", err)
	}

	// Initialize dependencies (Clean Architecture - manually wiring)
	userRepo := repository.NewUserRepository()
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase, zapLogger)

	// Setup router
	mux := http.NewServeMux()
	userHandler.RegisterRoutes(mux)

	// Health check endpoint (for Traefik health checks)
	mux.HandleFunc("GET /health", healthCheck)

	// Create server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		zapLogger.Info("Starting server", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zapLogger.Fatal("Server failed", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zapLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zapLogger.Fatal("Server forced to shutdown", "error", err)
	}

	zapLogger.Info("Server exited gracefully")
}

// healthCheck returns 200 OK with JSON status for Traefik health checks
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"service": "go-backend-standart",
		"time": time.Now().Format(time.RFC3339),
	})
}
