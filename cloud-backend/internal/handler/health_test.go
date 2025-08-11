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

func TestMain(m *testing.M) {
	// Setup test environment
	testutil.SetupTestEnvironment()

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Run tests
	code := m.Run()

	// Cleanup
	testutil.CleanupTestEnvironment()

	// Exit with test result code
	exit(code)
}

// exit is a variable to allow mocking in tests
var exit = func(_ int) {
	// In tests, we don't want to actually exit
}

func TestHealthHandler_HealthLive(t *testing.T) {
	// Setup
	cfg := testutil.TestConfig()
	handler := NewHealthHandler(cfg, nil)

	// Create test request
	req, err := http.NewRequest("GET", "/health/live", http.NoBody)
	require.NoError(t, err)

	// Create response recorder
	w := httptest.NewRecorder()

	// Create Gin context
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.HealthLive(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "alive", response["status"])
	assert.Contains(t, response, "timestamp")
}

func TestHealthHandler_HealthReady_WithoutDB(t *testing.T) {
	// Setup - handler without database manager
	cfg := testutil.TestConfig()
	handler := NewHealthHandler(cfg, nil)

	// Create test request
	req, err := http.NewRequest("GET", "/health/ready", http.NoBody)
	require.NoError(t, err)

	// Create response recorder
	w := httptest.NewRecorder()

	// Create Gin context
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.HealthReady(c)

	// Assert
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "not ready", response["status"])
	assert.Equal(t, "Database not ready", response["error"])
	assert.Contains(t, response, "timestamp")
}

func TestHealthHandler_Health_WithoutDB(t *testing.T) {
	// Setup - handler without database manager
	cfg := testutil.TestConfig()
	handler := NewHealthHandler(cfg, nil)

	// Create test request
	req, err := http.NewRequest("GET", "/health", http.NoBody)
	require.NoError(t, err)

	// Create response recorder
	w := httptest.NewRecorder()

	// Create Gin context
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.Health(c)

	// Assert
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)

	var response HealthResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "unhealthy", response.Status)
	assert.Equal(t, cfg.App.Version, response.Version)
	assert.Equal(t, cfg.App.Environment, response.Environment)
	assert.Contains(t, response.Services, "database")

	dbService, ok := response.Services["database"].(map[string]interface{})
	require.True(t, ok, "database service should be a map[string]interface{}")
	assert.Equal(t, "unhealthy", dbService["status"])
}

func TestHealthHandler_Metrics_WithoutDB(t *testing.T) {
	// Setup
	cfg := testutil.TestConfig()
	handler := NewHealthHandler(cfg, nil)

	// Create test request
	req, err := http.NewRequest("GET", "/metrics", http.NoBody)
	require.NoError(t, err)

	// Create response recorder
	w := httptest.NewRecorder()

	// Create Gin context
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	handler.Metrics(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response MonitoringResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	// Check application metrics
	assert.Equal(t, cfg.App.Name, response.Application.Name)
	assert.Equal(t, cfg.App.Version, response.Application.Version)
	assert.Equal(t, cfg.App.Environment, response.Application.Environment)
	assert.NotEmpty(t, response.Application.Uptime)

	// Check system metrics
	assert.NotEmpty(t, response.System.Hostname)
	assert.NotEmpty(t, response.System.Platform)
	assert.NotEmpty(t, response.System.GoVersion)
	assert.Greater(t, response.System.Goroutines, 0)
	assert.Greater(t, response.System.CPUStats.Count, 0)

	// Check memory stats
	assert.GreaterOrEqual(t, response.System.MemoryStats.Allocated, uint64(0))
	assert.GreaterOrEqual(t, response.System.MemoryStats.SystemMemory, uint64(0))

	// Check timestamp
	assert.NotEmpty(t, response.Timestamp)
}

func TestNewHealthHandler(t *testing.T) {
	cfg := testutil.TestConfig()

	handler := NewHealthHandler(cfg, nil)

	assert.NotNil(t, handler)
	assert.Equal(t, cfg, handler.config)
	assert.Nil(t, handler.db)
}

// Benchmark tests
func BenchmarkHealthHandler_HealthLive(b *testing.B) {
	cfg := testutil.TestConfig()
	handler := NewHealthHandler(cfg, nil)

	gin.SetMode(gin.TestMode)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/health/live", http.NoBody)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.HealthLive(c)
	}
}

func BenchmarkHealthHandler_Metrics(b *testing.B) {
	cfg := testutil.TestConfig()
	handler := NewHealthHandler(cfg, nil)

	gin.SetMode(gin.TestMode)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := http.NewRequest("GET", "/metrics", http.NoBody)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.Metrics(c)
	}
}
