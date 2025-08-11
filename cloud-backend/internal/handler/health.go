// Package handler provides HTTP handlers for the application.
package handler

import (
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/pkg/database"
)

// HealthHandler handles health check and monitoring endpoints.
type HealthHandler struct {
	config *config.Config
	db     *database.Manager
}

// NewHealthHandler creates a new health check handler.
func NewHealthHandler(cfg *config.Config, db *database.Manager) *HealthHandler {
	return &HealthHandler{
		config: cfg,
		db:     db,
	}
}

// HealthResponse represents the structure of health check response.
type HealthResponse struct {
	Status      string                 `json:"status"`
	Version     string                 `json:"version"`
	Environment string                 `json:"environment"`
	Timestamp   string                 `json:"timestamp"`
	Uptime      string                 `json:"uptime"`
	Services    map[string]interface{} `json:"services"`
	System      SystemInfo             `json:"system"`
}

// SystemInfo contains system information.
type SystemInfo struct {
	Hostname    string `json:"hostname"`
	Platform    string `json:"platform"`
	GoVersion   string `json:"go_version"`
	Goroutines  int    `json:"goroutines"`
	MemoryUsage uint64 `json:"memory_usage_mb"`
	CPUCount    int    `json:"cpu_count"`
}

// MonitoringResponse represents the structure of monitoring response.
type MonitoringResponse struct {
	Application ApplicationMetrics `json:"application"`
	Database    DatabaseMetrics    `json:"database"`
	System      SystemMetrics      `json:"system"`
	Timestamp   string             `json:"timestamp"`
}

// ApplicationMetrics contains application-level metrics.
type ApplicationMetrics struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Uptime      string `json:"uptime"`
}

// DatabaseMetrics contains database connection metrics.
type DatabaseMetrics struct {
	MySQL   interface{} `json:"mysql"`
	MongoDB interface{} `json:"mongodb"`
	Redis   interface{} `json:"redis"`
}

// SystemMetrics contains system performance metrics.
type SystemMetrics struct {
	Hostname    string      `json:"hostname"`
	Platform    string      `json:"platform"`
	GoVersion   string      `json:"go_version"`
	Goroutines  int         `json:"goroutines"`
	MemoryStats MemoryStats `json:"memory"`
	CPUStats    CPUStats    `json:"cpu"`
}

// MemoryStats contains memory usage statistics.
type MemoryStats struct {
	Allocated    uint64 `json:"allocated_mb"`
	TotalAlloc   uint64 `json:"total_alloc_mb"`
	SystemMemory uint64 `json:"system_mb"`
	GCCycles     uint32 `json:"gc_cycles"`
}

// CPUStats contains CPU usage statistics.
type CPUStats struct {
	Count int `json:"count"`
}

// startTime records when the application started.
var startTime = time.Now()

// Health handles basic health check requests.
// @Summary Health check endpoint
// @Description Returns the health status of the application and its dependencies
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse "Service is healthy"
// @Failure 503 {object} HealthResponse "Service is unhealthy"
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	// Get system information
	hostname, _ := os.Hostname()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	systemInfo := SystemInfo{
		Hostname:    hostname,
		Platform:    runtime.GOOS + "/" + runtime.GOARCH,
		GoVersion:   runtime.Version(),
		Goroutines:  runtime.NumGoroutine(),
		MemoryUsage: m.Sys / 1024 / 1024, // Convert to MB
		CPUCount:    runtime.NumCPU(),
	}

	// Initialize response
	response := HealthResponse{
		Status:      "healthy",
		Version:     h.config.App.Version,
		Environment: h.config.App.Environment,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Uptime:      time.Since(startTime).String(),
		Services:    make(map[string]interface{}),
		System:      systemInfo,
	}

	// Check database health
	if h.db == nil {
		response.Status = "unhealthy"
		response.Services["database"] = map[string]interface{}{
			"status": "unhealthy",
			"error":  "Database manager not initialized",
		}
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	if err := h.db.HealthCheck(); err != nil {
		response.Status = "unhealthy"
		response.Services["database"] = map[string]interface{}{
			"status": "unhealthy",
			"error":  err.Error(),
		}
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	// Get database connection stats
	dbStats := h.db.GetConnectionStats()
	response.Services["database"] = map[string]interface{}{
		"status": "healthy",
		"stats":  dbStats,
	}

	c.JSON(http.StatusOK, response)
}

// HealthLive handles liveness probe requests (for Kubernetes).
// @Summary Liveness probe endpoint
// @Description Returns 200 if the application is running
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string "Application is alive"
// @Router /health/live [get]
func (h *HealthHandler) HealthLive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// HealthReady handles readiness probe requests (for Kubernetes).
// @Summary Readiness probe endpoint
// @Description Returns 200 if the application is ready to serve traffic
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string "Application is ready"
// @Failure 503 {object} map[string]string "Application is not ready"
// @Router /health/ready [get]
func (h *HealthHandler) HealthReady(c *gin.Context) {
	// Check if database is ready
	if h.db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "not ready",
			"error":     "Database not ready",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	if err := h.db.HealthCheck(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "not ready",
			"error":     "Database not ready",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// Metrics handles detailed monitoring and metrics requests.
// @Summary Monitoring metrics endpoint
// @Description Returns detailed application and system metrics
// @Tags Monitoring
// @Produce json
// @Success 200 {object} MonitoringResponse "Monitoring metrics"
// @Router /metrics [get]
func (h *HealthHandler) Metrics(c *gin.Context) {
	hostname, _ := os.Hostname()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Application metrics
	appMetrics := ApplicationMetrics{
		Name:        h.config.App.Name,
		Version:     h.config.App.Version,
		Environment: h.config.App.Environment,
		Uptime:      time.Since(startTime).String(),
	}

	// Database metrics
	dbMetrics := DatabaseMetrics{
		MySQL:   nil,
		MongoDB: nil,
		Redis:   nil,
	}

	// Get database connection stats if available
	if h.db != nil {
		stats := h.db.GetConnectionStats()
		if mysql, ok := stats["mysql"]; ok {
			dbMetrics.MySQL = mysql
		}
		if mongodb, ok := stats["mongodb"]; ok {
			dbMetrics.MongoDB = mongodb
		}
		if redis, ok := stats["redis"]; ok {
			dbMetrics.Redis = redis
		}
	}

	// System metrics
	systemMetrics := SystemMetrics{
		Hostname:   hostname,
		Platform:   runtime.GOOS + "/" + runtime.GOARCH,
		GoVersion:  runtime.Version(),
		Goroutines: runtime.NumGoroutine(),
		MemoryStats: MemoryStats{
			Allocated:    m.Alloc / 1024 / 1024,      // Convert to MB
			TotalAlloc:   m.TotalAlloc / 1024 / 1024, // Convert to MB
			SystemMemory: m.Sys / 1024 / 1024,        // Convert to MB
			GCCycles:     m.NumGC,
		},
		CPUStats: CPUStats{
			Count: runtime.NumCPU(),
		},
	}

	response := MonitoringResponse{
		Application: appMetrics,
		Database:    dbMetrics,
		System:      systemMetrics,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}
