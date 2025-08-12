// Package main provides a database connection test utility.
package main

import (
	"log"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/pkg/database"
)

func main() {
	log.Println("Starting database connection test...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %w", err)
	}

	log.Printf("Loaded configuration for environment: %s", cfg.App.Environment)

	// Test database connections
	log.Println("Testing database connections with real credentials...")

	// Initialize database manager
	dbManager, err := database.NewManager(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database manager: %w", err)
	}
	defer func() {
		if err := dbManager.Close(); err != nil {
			log.Printf("Failed to close database connections: %w", err)
		}
	}()

	// Perform health checks
	if err := dbManager.HealthCheck(); err != nil {
		log.Printf("Database health check failed: %w", err)
		// Let defer dbManager.Close() run naturally
		return
	}

	// Get connection statistics
	stats := dbManager.GetConnectionStats()
	log.Printf("Database connection statistics: %+v", stats)

	log.Println("✅ All database connections successful!")
	log.Println("Database connection test completed successfully.")

	// Connection details
	log.Printf("✅ MySQL connected to: %s:%d", cfg.Database.MySQL.Host, cfg.Database.MySQL.Port)
	log.Printf("✅ MongoDB connected to: %s", cfg.Database.MongoDB.URI)
	log.Printf("✅ Redis connected to: %s:%d", cfg.Database.Redis.Host, cfg.Database.Redis.Port)
}
