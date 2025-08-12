// Package database provides tests for connection pool functionality.
package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
)

// createTestDB creates a test database for testing.
func createTestDB(t testing.TB) *gorm.DB {
	cfg := &config.MySQLConfig{
		Host:            "localhost",
		Port:            3306,
		Username:        "test",
		Password:        "test",
		Database:        "test",
		Charset:         "utf8mb4",
		ParseTime:       true,
		Loc:             "Local",
		Timeout:         10 * time.Second,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		ConnMaxIdleTime: 10 * time.Minute,
	}

	// Note: In real tests, you might want to use SQLite or a test database
	// For this example, we'll create a mock or skip if no test DB is available
	db, err := NewMySQLConnection(cfg)
	if err != nil {
		t.Skip("Skipping test: test database not available")
	}

	return db
}

func TestPoolConfig_Default(t *testing.T) {
	config := DefaultPoolConfig()

	assert.Equal(t, 100, config.MaxOpenConns)
	assert.Equal(t, 10, config.MaxIdleConns)
	assert.Equal(t, time.Hour, config.ConnMaxLifetime)
	assert.Equal(t, 10*time.Minute, config.ConnMaxIdleTime)
	assert.Equal(t, 30*time.Second, config.ConnectionTimeout)
	assert.True(t, config.HealthCheckEnabled)
	assert.True(t, config.MetricsEnabled)
}

func TestConnectionPool_Creation(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	config := DefaultPoolConfig()
	config.HealthCheckEnabled = false // Disable for testing

	pool, err := NewConnectionPool(db, config)
	require.NoError(t, err)
	require.NotNil(t, pool)

	defer pool.Close()

	assert.Equal(t, db, pool.GetDB())
	assert.Equal(t, config, pool.GetConfig())
}

func TestConnectionPool_Metrics(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	config := DefaultPoolConfig()
	config.HealthCheckEnabled = false
	config.MetricsEnabled = true

	pool, err := NewConnectionPool(db, config)
	require.NoError(t, err)
	defer pool.Close()

	// Get initial metrics
	metrics := pool.GetMetrics()
	assert.Equal(t, "healthy", metrics.PoolStatus)
	assert.Equal(t, int64(0), metrics.TotalQueries)

	// Execute some operations to generate metrics
	err = pool.ExecuteWithMetrics(func(db *gorm.DB) error {
		// Simulate query
		time.Sleep(10 * time.Millisecond)
		return nil
	})
	require.NoError(t, err)

	// Check updated metrics
	metrics = pool.GetMetrics()
	assert.Equal(t, int64(1), metrics.TotalQueries)
	assert.True(t, metrics.AverageQueryTime > 0)
}

func TestConnectionPool_ExecuteWithTimeout(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	config := DefaultPoolConfig()
	config.HealthCheckEnabled = false

	pool, err := NewConnectionPool(db, config)
	require.NoError(t, err)
	defer pool.Close()

	// Test successful execution within timeout
	err = pool.ExecuteWithTimeout(1*time.Second, func(db *gorm.DB) error {
		time.Sleep(100 * time.Millisecond)
		return nil
	})
	assert.NoError(t, err)

	// Test timeout scenario
	err = pool.ExecuteWithTimeout(100*time.Millisecond, func(db *gorm.DB) error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestConnectionPool_ExecuteWithContext(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	config := DefaultPoolConfig()
	config.HealthCheckEnabled = false

	pool, err := NewConnectionPool(db, config)
	require.NoError(t, err)
	defer pool.Close()

	// Test context cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context after short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err = pool.ExecuteWithContext(ctx, func(db *gorm.DB) error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}

func TestConnectionPool_HealthCheck(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	config := DefaultPoolConfig()
	config.HealthCheckEnabled = true
	config.HealthCheckInterval = 100 * time.Millisecond

	pool, err := NewConnectionPool(db, config)
	require.NoError(t, err)
	defer pool.Close()

	// Initially should be healthy
	assert.True(t, pool.IsHealthy())

	// Wait for a health check cycle
	time.Sleep(200 * time.Millisecond)

	// Should still be healthy
	assert.True(t, pool.IsHealthy())

	metrics := pool.GetMetrics()
	assert.True(t, metrics.TotalHealthChecks > 0)
	assert.Equal(t, "healthy", metrics.HealthCheckStatus)
}

func TestConnectionPool_RetryLogic(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	config := DefaultPoolConfig()
	config.HealthCheckEnabled = false
	config.RetryAttempts = 3
	config.RetryDelay = 10 * time.Millisecond

	pool, err := NewConnectionPool(db, config)
	require.NoError(t, err)
	defer pool.Close()

	attempts := 0
	err = pool.ExecuteWithRetry(func(db *gorm.DB) error {
		attempts++
		if attempts < 3 {
			return assert.AnError // Simulate failure
		}
		return nil // Success on third attempt
	})

	assert.NoError(t, err)
	assert.Equal(t, 3, attempts)
}

func TestCreatePoolFromConfig(t *testing.T) {
	db := createTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	mysqlConfig := &config.MySQLConfig{
		MaxOpenConns:    50,
		MaxIdleConns:    20,
		ConnMaxLifetime: 2 * time.Hour,
		ConnMaxIdleTime: 30 * time.Minute,
		Timeout:         15 * time.Second,
		ReadTimeout:     45 * time.Second,
	}

	pool, err := CreatePoolFromConfig(db, mysqlConfig)
	require.NoError(t, err)
	require.NotNil(t, pool)
	defer pool.Close()

	poolConfig := pool.GetConfig()
	assert.Equal(t, 50, poolConfig.MaxOpenConns)
	assert.Equal(t, 20, poolConfig.MaxIdleConns)
	assert.Equal(t, 2*time.Hour, poolConfig.ConnMaxLifetime)
	assert.Equal(t, 30*time.Minute, poolConfig.ConnMaxIdleTime)
	assert.Equal(t, 15*time.Second, poolConfig.ConnectionTimeout)
	assert.Equal(t, 45*time.Second, poolConfig.QueryTimeout)
}

// Benchmark tests
func BenchmarkConnectionPool_ExecuteWithMetrics(b *testing.B) {
	db := createTestDB(b)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	config := DefaultPoolConfig()
	config.HealthCheckEnabled = false

	pool, err := NewConnectionPool(db, config)
	require.NoError(b, err)
	defer pool.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.ExecuteWithMetrics(func(db *gorm.DB) error {
				// Simulate lightweight operation
				return nil
			})
		}
	})
}

func BenchmarkConnectionPool_ExecuteWithRetry(b *testing.B) {
	db := createTestDB(b)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	config := DefaultPoolConfig()
	config.HealthCheckEnabled = false
	config.RetryAttempts = 1 // Minimize retries for benchmark

	pool, err := NewConnectionPool(db, config)
	require.NoError(b, err)
	defer pool.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.ExecuteWithRetry(func(db *gorm.DB) error {
				// Always succeed for benchmark
				return nil
			})
		}
	})
}
