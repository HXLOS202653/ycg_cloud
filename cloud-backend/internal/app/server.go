// Package app provides the main application setup and configuration.
package app

import (
	"github.com/gin-gonic/gin"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/pkg/database"
)

// Server represents the main application server.
type Server struct {
	Router *gin.Engine
	DB     *database.Manager
	Config *config.Config
}

// NewServer creates a new server instance with the given configuration.
func NewServer(cfg *config.Config) (*Server, error) {
	// Initialize database connections
	dbManager, err := database.NewManager(&cfg.Database)
	if err != nil {
		return nil, err
	}

	// Test database connections
	if err := dbManager.HealthCheck(); err != nil {
		return nil, err
	}

	// Create server instance
	server := &Server{
		DB:     dbManager,
		Config: cfg,
	}

	// Setup router
	server.setupRouter()

	return server, nil
}

// setupRouter configures the Gin router with middlewares and routes.
func (s *Server) setupRouter() {
	// Set Gin mode based on environment
	if s.Config.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create router
	r := gin.New()

	// Add middlewares
	s.addMiddlewares(r)

	// Add routes
	s.addRoutes(r)

	s.Router = r
}

// Close closes all server resources including database connections.
func (s *Server) Close() error {
	if s.DB != nil {
		return s.DB.Close()
	}
	return nil
}
