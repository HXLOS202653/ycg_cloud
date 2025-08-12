// Package database provides enhanced connection pool management.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
)

// PoolConfig contains advanced connection pool configuration.
type PoolConfig struct {
	// Basic pool settings
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`

	// Advanced settings
	ConnectionTimeout time.Duration `json:"connection_timeout"`
	QueryTimeout      time.Duration `json:"query_timeout"`
	PingInterval      time.Duration `json:"ping_interval"`
	RetryAttempts     int           `json:"retry_attempts"`
	RetryDelay        time.Duration `json:"retry_delay"`

	// Health check settings
	HealthCheckEnabled  bool          `json:"health_check_enabled"`
	HealthCheckInterval time.Duration `json:"health_check_interval"`
	HealthCheckTimeout  time.Duration `json:"health_check_timeout"`

	// Monitoring settings
	MetricsEnabled     bool          `json:"metrics_enabled"`
	SlowQueryThreshold time.Duration `json:"slow_query_threshold"`
}

// DefaultPoolConfig returns default pool configuration.
func DefaultPoolConfig() *PoolConfig {
	return &PoolConfig{
		MaxOpenConns:        100,
		MaxIdleConns:        10,
		ConnMaxLifetime:     time.Hour,
		ConnMaxIdleTime:     10 * time.Minute,
		ConnectionTimeout:   30 * time.Second,
		QueryTimeout:        30 * time.Second,
		PingInterval:        time.Minute,
		RetryAttempts:       3,
		RetryDelay:          time.Second,
		HealthCheckEnabled:  true,
		HealthCheckInterval: 30 * time.Second,
		HealthCheckTimeout:  5 * time.Second,
		MetricsEnabled:      true,
		SlowQueryThreshold:  time.Second,
	}
}

// PoolMetrics contains connection pool metrics.
type PoolMetrics struct {
	// Connection statistics
	OpenConnections   int `json:"open_connections"`
	InUseConnections  int `json:"in_use_connections"`
	IdleConnections   int `json:"idle_connections"`
	WaitCount         int `json:"wait_count"`
	WaitDuration      int `json:"wait_duration_ms"`
	MaxIdleClosed     int `json:"max_idle_closed"`
	MaxIdleTimeClosed int `json:"max_idle_time_closed"`
	MaxLifetimeClosed int `json:"max_lifetime_closed"`

	// Query statistics
	TotalQueries     int64         `json:"total_queries"`
	SlowQueries      int64         `json:"slow_queries"`
	FailedQueries    int64         `json:"failed_queries"`
	AverageQueryTime time.Duration `json:"average_query_time"`

	// Health check statistics
	HealthCheckStatus   string    `json:"health_check_status"`
	LastHealthCheck     time.Time `json:"last_health_check"`
	ConsecutiveFailures int       `json:"consecutive_failures"`
	TotalHealthChecks   int64     `json:"total_health_checks"`
	FailedHealthChecks  int64     `json:"failed_health_checks"`

	// Pool status
	PoolStatus string    `json:"pool_status"`
	LastError  string    `json:"last_error,omitempty"`
	Uptime     time.Time `json:"uptime"`
}

// ConnectionPool manages an enhanced database connection pool.
type ConnectionPool struct {
	db     *gorm.DB
	sqlDB  *sql.DB
	config *PoolConfig

	// Metrics and monitoring
	metrics     *PoolMetrics
	metricsLock sync.RWMutex

	// Health monitoring
	healthTicker   *time.Ticker
	healthStopChan chan struct{}

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc
}

// NewConnectionPool creates a new enhanced connection pool.
func NewConnectionPool(db *gorm.DB, config *PoolConfig) (*ConnectionPool, error) {
	if config == nil {
		config = DefaultPoolConfig()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure basic pool settings
	if err := configurePoolSettings(sqlDB, config); err != nil {
		return nil, fmt.Errorf("failed to configure pool settings: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	pool := &ConnectionPool{
		db:             db,
		sqlDB:          sqlDB,
		config:         config,
		metrics:        &PoolMetrics{Uptime: time.Now(), PoolStatus: "healthy"},
		healthStopChan: make(chan struct{}),
		ctx:            ctx,
		cancel:         cancel,
	}

	// Start health monitoring if enabled
	if config.HealthCheckEnabled {
		pool.startHealthMonitoring()
	}

	log.Printf("Enhanced connection pool initialized with config: %+v", config)
	return pool, nil
}

// configurePoolSettings applies pool configuration to sql.DB.
func configurePoolSettings(sqlDB *sql.DB, config *PoolConfig) error {
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), config.ConnectionTimeout)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Connection pool configured: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%v, MaxIdleTime=%v",
		config.MaxOpenConns, config.MaxIdleConns, config.ConnMaxLifetime, config.ConnMaxIdleTime)

	return nil
}

// startHealthMonitoring starts the health check monitoring goroutine.
func (cp *ConnectionPool) startHealthMonitoring() {
	cp.healthTicker = time.NewTicker(cp.config.HealthCheckInterval)

	go func() {
		defer cp.healthTicker.Stop()

		for {
			select {
			case <-cp.healthTicker.C:
				cp.performHealthCheck()
			case <-cp.healthStopChan:
				return
			case <-cp.ctx.Done():
				return
			}
		}
	}()

	log.Printf("Health monitoring started with interval: %v", cp.config.HealthCheckInterval)
}

// performHealthCheck performs a health check on the connection pool.
func (cp *ConnectionPool) performHealthCheck() {
	cp.metricsLock.Lock()
	defer cp.metricsLock.Unlock()

	cp.metrics.TotalHealthChecks++
	cp.metrics.LastHealthCheck = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), cp.config.HealthCheckTimeout)
	defer cancel()

	if err := cp.sqlDB.PingContext(ctx); err != nil {
		cp.metrics.FailedHealthChecks++
		cp.metrics.ConsecutiveFailures++
		cp.metrics.HealthCheckStatus = "unhealthy"
		cp.metrics.LastError = err.Error()

		if cp.metrics.ConsecutiveFailures >= 3 {
			cp.metrics.PoolStatus = "degraded"
		}

		log.Printf("Health check failed: %v (consecutive failures: %d)", err, cp.metrics.ConsecutiveFailures)
	} else {
		cp.metrics.ConsecutiveFailures = 0
		cp.metrics.HealthCheckStatus = "healthy"
		cp.metrics.PoolStatus = "healthy"
		cp.metrics.LastError = ""
	}
}

// GetMetrics returns current pool metrics.
func (cp *ConnectionPool) GetMetrics() PoolMetrics {
	cp.metricsLock.RLock()
	defer cp.metricsLock.RUnlock()

	// Update current connection statistics
	stats := cp.sqlDB.Stats()
	cp.metrics.OpenConnections = stats.OpenConnections
	cp.metrics.InUseConnections = stats.InUse
	cp.metrics.IdleConnections = stats.Idle
	cp.metrics.WaitCount = int(stats.WaitCount)
	cp.metrics.WaitDuration = int(stats.WaitDuration.Milliseconds())
	cp.metrics.MaxIdleClosed = int(stats.MaxIdleClosed)
	cp.metrics.MaxIdleTimeClosed = int(stats.MaxIdleTimeClosed)
	cp.metrics.MaxLifetimeClosed = int(stats.MaxLifetimeClosed)

	return *cp.metrics
}

// UpdateQueryMetrics updates query performance metrics.
func (cp *ConnectionPool) UpdateQueryMetrics(duration time.Duration, failed bool) {
	cp.metricsLock.Lock()
	defer cp.metricsLock.Unlock()

	cp.metrics.TotalQueries++

	if failed {
		cp.metrics.FailedQueries++
	}

	if duration > cp.config.SlowQueryThreshold {
		cp.metrics.SlowQueries++
	}

	// Update average query time (simplified moving average)
	if cp.metrics.TotalQueries == 1 {
		cp.metrics.AverageQueryTime = duration
	} else {
		cp.metrics.AverageQueryTime = (cp.metrics.AverageQueryTime + duration) / 2
	}
}

// ExecuteWithMetrics executes a function with query metrics tracking.
func (cp *ConnectionPool) ExecuteWithMetrics(fn func(*gorm.DB) error) error {
	if !cp.config.MetricsEnabled {
		return fn(cp.db)
	}

	start := time.Now()
	err := fn(cp.db)
	duration := time.Since(start)

	cp.UpdateQueryMetrics(duration, err != nil)

	return err
}

// ExecuteWithTimeout executes a function with timeout.
func (cp *ConnectionPool) ExecuteWithTimeout(timeout time.Duration, fn func(*gorm.DB) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return cp.ExecuteWithContext(ctx, fn)
}

// ExecuteWithContext executes a function with context.
func (cp *ConnectionPool) ExecuteWithContext(ctx context.Context, fn func(*gorm.DB) error) error {
	return fn(cp.db.WithContext(ctx))
}

// ExecuteWithRetry executes a function with retry logic.
func (cp *ConnectionPool) ExecuteWithRetry(fn func(*gorm.DB) error) error {
	var lastErr error

	for attempt := 0; attempt <= cp.config.RetryAttempts; attempt++ {
		if attempt > 0 {
			// Wait before retry
			select {
			case <-time.After(cp.config.RetryDelay * time.Duration(attempt)):
			case <-cp.ctx.Done():
				return cp.ctx.Err()
			}
		}

		start := time.Now()
		err := fn(cp.db)
		duration := time.Since(start)

		if cp.config.MetricsEnabled {
			cp.UpdateQueryMetrics(duration, err != nil)
		}

		if err == nil {
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if !isRetryableError(err) {
			break
		}

		log.Printf("Database operation failed (attempt %d/%d): %v", attempt+1, cp.config.RetryAttempts+1, err)
	}

	return fmt.Errorf("operation failed after %d attempts: %w", cp.config.RetryAttempts+1, lastErr)
}

// isRetryableError determines if an error is retryable.
func isRetryableError(err error) bool {
	// Add logic to determine if error is retryable
	// For now, assume network-related errors are retryable
	return err != nil
}

// Close closes the connection pool and stops monitoring.
func (cp *ConnectionPool) Close() error {
	cp.cancel()

	if cp.healthTicker != nil {
		cp.healthTicker.Stop()
		close(cp.healthStopChan)
	}

	return cp.sqlDB.Close()
}

// GetDB returns the GORM database instance.
func (cp *ConnectionPool) GetDB() *gorm.DB {
	return cp.db
}

// GetSQLDB returns the underlying sql.DB instance.
func (cp *ConnectionPool) GetSQLDB() *sql.DB {
	return cp.sqlDB
}

// GetConfig returns the pool configuration.
func (cp *ConnectionPool) GetConfig() *PoolConfig {
	return cp.config
}

// IsHealthy returns whether the pool is healthy.
func (cp *ConnectionPool) IsHealthy() bool {
	cp.metricsLock.RLock()
	defer cp.metricsLock.RUnlock()
	return cp.metrics.PoolStatus == "healthy"
}

// CreatePoolFromConfig creates a connection pool from MySQL config.
func CreatePoolFromConfig(db *gorm.DB, mysqlConfig *config.MySQLConfig) (*ConnectionPool, error) {
	poolConfig := &PoolConfig{
		MaxOpenConns:        mysqlConfig.MaxOpenConns,
		MaxIdleConns:        mysqlConfig.MaxIdleConns,
		ConnMaxLifetime:     mysqlConfig.ConnMaxLifetime,
		ConnMaxIdleTime:     mysqlConfig.ConnMaxIdleTime,
		ConnectionTimeout:   mysqlConfig.Timeout,
		QueryTimeout:        mysqlConfig.ReadTimeout,
		PingInterval:        time.Minute,
		RetryAttempts:       3,
		RetryDelay:          time.Second,
		HealthCheckEnabled:  true,
		HealthCheckInterval: 30 * time.Second,
		HealthCheckTimeout:  5 * time.Second,
		MetricsEnabled:      true,
		SlowQueryThreshold:  time.Second,
	}

	return NewConnectionPool(db, poolConfig)
}
