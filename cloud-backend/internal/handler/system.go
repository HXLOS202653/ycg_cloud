// Package handler provides system monitoring handlers.
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

// SystemHandler handles system monitoring and statistics endpoints.
type SystemHandler struct {
	config *config.Config
	db     *database.Manager
}

// NewSystemHandler creates a new system monitoring handler.
func NewSystemHandler(cfg *config.Config, db *database.Manager) *SystemHandler {
	return &SystemHandler{
		config: cfg,
		db:     db,
	}
}

// SystemStatsResponse represents system statistics response.
type SystemStatsResponse struct {
	Application ApplicationStats      `json:"application"`
	System      DetailedSystemStats   `json:"system"`
	Database    DetailedDatabaseStats `json:"database"`
	Timestamp   string                `json:"timestamp"`
}

// ApplicationStats contains application statistics.
type ApplicationStats struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Environment string `json:"environment"`
	Uptime      string `json:"uptime"`
	StartTime   string `json:"start_time"`
}

// DetailedSystemStats contains detailed system statistics.
type DetailedSystemStats struct {
	Hostname      string             `json:"hostname"`
	Platform      string             `json:"platform"`
	Architecture  string             `json:"architecture"`
	GoVersion     string             `json:"go_version"`
	CPUInfo       CPUInfo            `json:"cpu"`
	MemoryInfo    DetailedMemoryInfo `json:"memory"`
	GoroutineInfo GoroutineInfo      `json:"goroutines"`
	GCInfo        GCInfo             `json:"garbage_collection"`
}

// CPUInfo contains CPU information.
type CPUInfo struct {
	Count int `json:"count"`
}

// DetailedMemoryInfo contains detailed memory information.
type DetailedMemoryInfo struct {
	Allocated      uint64 `json:"allocated_bytes"`
	AllocatedMB    uint64 `json:"allocated_mb"`
	TotalAlloc     uint64 `json:"total_alloc_bytes"`
	TotalAllocMB   uint64 `json:"total_alloc_mb"`
	System         uint64 `json:"system_bytes"`
	SystemMB       uint64 `json:"system_mb"`
	Lookups        uint64 `json:"lookups"`
	Mallocs        uint64 `json:"mallocs"`
	Frees          uint64 `json:"frees"`
	HeapAlloc      uint64 `json:"heap_alloc_bytes"`
	HeapAllocMB    uint64 `json:"heap_alloc_mb"`
	HeapSys        uint64 `json:"heap_sys_bytes"`
	HeapSysMB      uint64 `json:"heap_sys_mb"`
	HeapIdle       uint64 `json:"heap_idle_bytes"`
	HeapIdleMB     uint64 `json:"heap_idle_mb"`
	HeapInuse      uint64 `json:"heap_inuse_bytes"`
	HeapInuseMB    uint64 `json:"heap_inuse_mb"`
	HeapReleased   uint64 `json:"heap_released_bytes"`
	HeapReleasedMB uint64 `json:"heap_released_mb"`
	HeapObjects    uint64 `json:"heap_objects"`
}

// GoroutineInfo contains goroutine information.
type GoroutineInfo struct {
	Count int `json:"count"`
}

// GCInfo contains garbage collection information.
type GCInfo struct {
	NumGC         uint32   `json:"num_gc"`
	NumForcedGC   uint32   `json:"num_forced_gc"`
	GCCPUFraction float64  `json:"gc_cpu_fraction"`
	EnableGC      bool     `json:"enable_gc"`
	DebugGC       bool     `json:"debug_gc"`
	LastGC        string   `json:"last_gc"`
	NextGC        uint64   `json:"next_gc_mb"`
	PauseTotal    string   `json:"pause_total"`
	PauseNs       []uint64 `json:"pause_ns"`
}

// DetailedDatabaseStats contains detailed database statistics.
type DetailedDatabaseStats struct {
	MySQL   interface{}     `json:"mysql"`
	MongoDB interface{}     `json:"mongodb"`
	Redis   interface{}     `json:"redis"`
	Health  map[string]bool `json:"health_status"`
}

// SystemStats handles detailed system statistics requests.
// @Summary System statistics endpoint
// @Description Returns detailed system and application statistics
// @Tags System
// @Produce json
// @Success 200 {object} SystemStatsResponse "System statistics"
// @Router /admin/stats [get]
func (s *SystemHandler) SystemStats(c *gin.Context) {
	hostname, _ := os.Hostname()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Application stats
	appStats := ApplicationStats{
		Name:        s.config.App.Name,
		Version:     s.config.App.Version,
		Environment: s.config.App.Environment,
		Uptime:      time.Since(startTime).String(),
		StartTime:   startTime.Format(time.RFC3339),
	}

	// Detailed memory information
	memoryInfo := DetailedMemoryInfo{
		Allocated:      m.Alloc,
		AllocatedMB:    m.Alloc / 1024 / 1024,
		TotalAlloc:     m.TotalAlloc,
		TotalAllocMB:   m.TotalAlloc / 1024 / 1024,
		System:         m.Sys,
		SystemMB:       m.Sys / 1024 / 1024,
		Lookups:        m.Lookups,
		Mallocs:        m.Mallocs,
		Frees:          m.Frees,
		HeapAlloc:      m.HeapAlloc,
		HeapAllocMB:    m.HeapAlloc / 1024 / 1024,
		HeapSys:        m.HeapSys,
		HeapSysMB:      m.HeapSys / 1024 / 1024,
		HeapIdle:       m.HeapIdle,
		HeapIdleMB:     m.HeapIdle / 1024 / 1024,
		HeapInuse:      m.HeapInuse,
		HeapInuseMB:    m.HeapInuse / 1024 / 1024,
		HeapReleased:   m.HeapReleased,
		HeapReleasedMB: m.HeapReleased / 1024 / 1024,
		HeapObjects:    m.HeapObjects,
	}

	// GC information
	gcInfo := GCInfo{
		NumGC:         m.NumGC,
		NumForcedGC:   m.NumForcedGC,
		GCCPUFraction: m.GCCPUFraction,
		EnableGC:      m.EnableGC,
		DebugGC:       m.DebugGC,
		NextGC:        m.NextGC / 1024 / 1024, // Convert to MB
		PauseTotal:    time.Duration(m.PauseTotalNs).String(),
		PauseNs:       m.PauseNs[:],
	}

	if m.NumGC > 0 {
		gcInfo.LastGC = time.Unix(0, int64(m.LastGC)).Format(time.RFC3339)
	} else {
		gcInfo.LastGC = "never"
	}

	// System stats
	systemStats := DetailedSystemStats{
		Hostname:     hostname,
		Platform:     runtime.GOOS,
		Architecture: runtime.GOARCH,
		GoVersion:    runtime.Version(),
		CPUInfo: CPUInfo{
			Count: runtime.NumCPU(),
		},
		MemoryInfo: memoryInfo,
		GoroutineInfo: GoroutineInfo{
			Count: runtime.NumGoroutine(),
		},
		GCInfo: gcInfo,
	}

	// Database stats
	dbStats := DetailedDatabaseStats{
		MySQL:   nil,
		MongoDB: nil,
		Redis:   nil,
		Health:  make(map[string]bool),
	}

	if s.db != nil {
		// Get database connection stats
		stats := s.db.GetConnectionStats()
		if mysql, ok := stats["mysql"]; ok {
			dbStats.MySQL = mysql
		}
		if mongodb, ok := stats["mongodb"]; ok {
			dbStats.MongoDB = mongodb
		}
		if redis, ok := stats["redis"]; ok {
			dbStats.Redis = redis
		}

		// Check individual database health
		dbStats.Health["mysql"] = s.checkMySQLHealth()
		dbStats.Health["mongodb"] = s.checkMongoDBHealth()
		dbStats.Health["redis"] = s.checkRedisHealth()
	}

	response := SystemStatsResponse{
		Application: appStats,
		System:      systemStats,
		Database:    dbStats,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}

// checkMySQLHealth checks MySQL database health.
func (s *SystemHandler) checkMySQLHealth() bool {
	if s.db == nil {
		return false
	}

	return s.db.MySQL.Health() == nil
}

// checkMongoDBHealth checks MongoDB database health.
func (s *SystemHandler) checkMongoDBHealth() bool {
	if s.db == nil {
		return false
	}

	return s.db.MongoDB.Health() == nil
}

// checkRedisHealth checks Redis database health.
func (s *SystemHandler) checkRedisHealth() bool {
	if s.db == nil {
		return false
	}

	return s.db.Redis.Health() == nil
}
