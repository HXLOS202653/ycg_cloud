// Package main provides the entry point for the cloud storage server.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/app"
	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %w", err)
	}

	// Create and start the server
	server, err := app.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %w", err)
	}

	// Configure HTTP server
	httpServer := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      server.Router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %w", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown the HTTP server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %w", err)
		// Let defer cancel() run naturally
		return
	}

	// Close database connections
	if err := server.Close(); err != nil {
		log.Printf("Error closing server resources: %w", err)
	}

	log.Println("Server exited")
}
