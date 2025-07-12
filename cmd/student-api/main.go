package main

import (
	"context"
	"fmt"
	"github/black-spidera/student-api/internal/config"
	"github/black-spidera/student-api/internal/http/handlers/students"
	"github/black-spidera/student-api/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize the application
	config := config.ConfigLoader()

	// DB connection and other initializations can be done here
	storage, err := sqlite.New(config)
	if err != nil {
		slog.Error("Failed to connect to the database", "error", err)
		os.Exit(1)
	}
	// Route handlers can be set up here
	routers := http.NewServeMux()
	routers.HandleFunc("POST /v1/api/students", students.New(storage))
	routers.HandleFunc("GET /v1/api/students/{id}", students.GetById(storage))
	// Start the HTTP server with the configuration
	server := &http.Server{
		Addr:    config.HTTPServer.Addr,
		Handler: routers,
	}
	fmt.Printf("Starting server on http://%s\n", config.HTTPServer.Addr)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
			return
		}
	}()

	// Wait for termination signal
	<-done
	slog.Info("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server gracefully", "error", err)
	} else {
		slog.Info("Server shutdown successfully")
	}
}
