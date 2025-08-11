package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/testutil"
)

func TestSystemHandler_SystemStats(t *testing.T) {
	// Setup
	cfg := testutil.TestConfig()
	handler := NewSystemHandler(cfg, nil)

	// Create test request
	req, err := http.NewRequest("GET", "/admin/stats", http.NoBody)
	require.NoError(t, err)

	// Create response recorder
	w := httptest.NewRecorder()

	// Create Gin context
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.SystemStats(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response SystemStatsResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Check application stats
	assert.Equal(t, cfg.App.Name, response.Application.Name)
	assert.Equal(t, cfg.App.Version, response.Application.Version)
	assert.Equal(t, cfg.App.Environment, response.Application.Environment)
	assert.NotEmpty(t, response.Application.Uptime)
	assert.NotEmpty(t, response.Application.StartTime)

	// Check system stats
	assert.NotEmpty(t, response.System.Hostname)
	assert.NotEmpty(t, response.System.Platform)
	assert.NotEmpty(t, response.System.Architecture)
	assert.NotEmpty(t, response.System.GoVersion)
	assert.Greater(t, response.System.CPUInfo.Count, 0)
	assert.Greater(t, response.System.GoroutineInfo.Count, 0)

	// Check memory info
	assert.GreaterOrEqual(t, response.System.MemoryInfo.AllocatedMB, uint64(0))
	assert.GreaterOrEqual(t, response.System.MemoryInfo.SystemMB, uint64(0))
	assert.GreaterOrEqual(t, response.System.MemoryInfo.HeapAllocMB, uint64(0))

	// Check GC info
	assert.GreaterOrEqual(t, response.System.GCInfo.NumGC, uint32(0))
	assert.NotEmpty(t, response.System.GCInfo.LastGC)
	assert.GreaterOrEqual(t, response.System.GCInfo.NextGC, uint64(0))

	// Check database stats structure
	assert.NotNil(t, response.Database.Health)

	// Check timestamp
	assert.NotEmpty(t, response.Timestamp)
}

func TestSystemHandler_checkMySQLHealth_WithoutDB(t *testing.T) {
	// Setup - handler without database manager
	cfg := testutil.TestConfig()
	handler := NewSystemHandler(cfg, nil)

	// Execute
	result := handler.checkMySQLHealth()

	// Assert
	assert.False(t, result)
}

func TestSystemHandler_checkMongoDBHealth_WithoutDB(t *testing.T) {
	// Setup - handler without database manager
	cfg := testutil.TestConfig()
	handler := NewSystemHandler(cfg, nil)

	// Execute
	result := handler.checkMongoDBHealth()

	// Assert
	assert.False(t, result)
}

func TestSystemHandler_checkRedisHealth_WithoutDB(t *testing.T) {
	// Setup - handler without database manager
	cfg := testutil.TestConfig()
	handler := NewSystemHandler(cfg, nil)

	// Execute
	result := handler.checkRedisHealth()

	// Assert
	assert.False(t, result)
}

func TestNewSystemHandler(t *testing.T) {
	cfg := testutil.TestConfig()

	handler := NewSystemHandler(cfg, nil)

	assert.NotNil(t, handler)
	assert.Equal(t, cfg, handler.config)
	assert.Nil(t, handler.db)
}

// Test detailed memory info structure
func TestDetailedMemoryInfo_Structure(t *testing.T) {
	cfg := testutil.TestConfig()
	handler := NewSystemHandler(cfg, nil)

	req, err := http.NewRequest("GET", "/admin/stats", http.NoBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.SystemStats(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response SystemStatsResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	memInfo := response.System.MemoryInfo

	// Check all memory fields are present and properly typed
	assert.GreaterOrEqual(t, memInfo.Allocated, uint64(0))
	assert.GreaterOrEqual(t, memInfo.AllocatedMB, uint64(0))
	assert.GreaterOrEqual(t, memInfo.TotalAlloc, uint64(0))
	assert.GreaterOrEqual(t, memInfo.TotalAllocMB, uint64(0))
	assert.GreaterOrEqual(t, memInfo.System, uint64(0))
	assert.GreaterOrEqual(t, memInfo.SystemMB, uint64(0))
	assert.GreaterOrEqual(t, memInfo.Lookups, uint64(0))
	assert.GreaterOrEqual(t, memInfo.Mallocs, uint64(0))
	assert.GreaterOrEqual(t, memInfo.Frees, uint64(0))
	assert.GreaterOrEqual(t, memInfo.HeapAlloc, uint64(0))
	assert.GreaterOrEqual(t, memInfo.HeapAllocMB, uint64(0))
	assert.GreaterOrEqual(t, memInfo.HeapSys, uint64(0))
	assert.GreaterOrEqual(t, memInfo.HeapSysMB, uint64(0))
	assert.GreaterOrEqual(t, memInfo.HeapIdle, uint64(0))
	assert.GreaterOrEqual(t, memInfo.HeapIdleMB, uint64(0))
	assert.GreaterOrEqual(t, memInfo.HeapInuse, uint64(0))
	assert.GreaterOrEqual(t, memInfo.HeapInuseMB, uint64(0))
	assert.GreaterOrEqual(t, memInfo.HeapReleased, uint64(0))
	assert.GreaterOrEqual(t, memInfo.HeapReleasedMB, uint64(0))
	assert.GreaterOrEqual(t, memInfo.HeapObjects, uint64(0))
}

// Benchmark tests
func BenchmarkSystemHandler_SystemStats(b *testing.B) {
	cfg := testutil.TestConfig()
	handler := NewSystemHandler(cfg, nil)

	gin.SetMode(gin.TestMode)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/admin/stats", http.NoBody)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.SystemStats(c)
	}
}
