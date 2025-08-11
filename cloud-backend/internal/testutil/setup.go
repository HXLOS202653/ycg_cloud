// Package testutil provides utilities for testing.
package testutil

import (
	"log"
	"os"
	"testing"
	"time"
)

// TestConfigData creates test configuration data without importing config package.
type TestConfigData struct {
	AppName        string
	AppVersion     string
	AppEnvironment string
	AppDebug       bool
	ServerHost     string
	ServerPort     string
	ServerTimeouts TestServerTimeouts
	Database       TestDatabaseConfig
	Auth           TestAuthConfig
	Log            TestLogConfig
	Monitor        TestMonitorConfig
}

// TestServerTimeouts represents test server timeout configuration.
type TestServerTimeouts struct {
	Read  int
	Write int
	Idle  int
}

// TestDatabaseConfig represents test database configuration.
type TestDatabaseConfig struct {
	MySQL   TestMySQLConfig
	MongoDB TestMongoDBConfig
	Redis   TestRedisConfig
}

// TestMySQLConfig represents test MySQL configuration.
type TestMySQLConfig struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	Charset         string
	ParseTime       bool
	Loc             string
	Timeout         time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// TestMongoDBConfig represents test MongoDB configuration.
type TestMongoDBConfig struct {
	URI              string
	Database         string
	Username         string
	Password         string
	ConnectTimeout   time.Duration
	SocketTimeout    time.Duration
	ServerSelTimeout time.Duration
	MaxPoolSize      uint64
	MinPoolSize      uint64
	MaxIdleTimeMS    uint64
}

// TestRedisConfig represents test Redis configuration.
type TestRedisConfig struct {
	Host          string
	Port          int
	Password      string
	DB            int
	Username      string
	PoolSize      int
	MinIdleConns  int
	MaxConnAge    time.Duration
	PoolTimeout   time.Duration
	IdleTimeout   time.Duration
	IdleCheckFreq time.Duration
	DialTimeout   time.Duration
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
}

// TestAuthConfig represents test authentication configuration.
type TestAuthConfig struct {
	JWTSecret     string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
	BcryptCost    int
}

// TestLogConfig represents test log configuration.
type TestLogConfig struct {
	Level  string
	Format string
	Output string
}

// TestMonitorConfig represents test monitoring configuration.
type TestMonitorConfig struct {
	Enabled    bool
	MetricsURL string
}

// GetTestConfig creates a test configuration.
func GetTestConfig() *TestConfigData {
	return &TestConfigData{
		AppName:        "ycg-cloud-test",
		AppVersion:     "1.0.0-test",
		AppEnvironment: "test",
		AppDebug:       true,
		ServerHost:     "localhost",
		ServerPort:     ":8081",
		ServerTimeouts: TestServerTimeouts{
			Read:  30,
			Write: 30,
			Idle:  120,
		},
		Database: TestDatabaseConfig{
			MySQL: TestMySQLConfig{
				Host:            "localhost",
				Port:            3306,
				Username:        "test",
				Password:        "test",
				Database:        "ycgcloud_test",
				Charset:         "utf8mb4",
				ParseTime:       true,
				Loc:             "Local",
				Timeout:         10 * time.Second,
				ReadTimeout:     30 * time.Second,
				WriteTimeout:    30 * time.Second,
				MaxOpenConns:    10,
				MaxIdleConns:    2,
				ConnMaxLifetime: time.Hour,
				ConnMaxIdleTime: 30 * time.Minute,
			},
			MongoDB: TestMongoDBConfig{
				URI:              "mongodb://localhost:27017",
				Database:         "ycgcloud_test",
				Username:         "",
				Password:         "",
				ConnectTimeout:   10 * time.Second,
				SocketTimeout:    30 * time.Second,
				ServerSelTimeout: 30 * time.Second,
				MaxPoolSize:      10,
				MinPoolSize:      2,
				MaxIdleTimeMS:    300000,
			},
			Redis: TestRedisConfig{
				Host:          "localhost",
				Port:          6379,
				Password:      "",
				DB:            1, // Use database 1 for tests
				Username:      "",
				PoolSize:      10,
				MinIdleConns:  2,
				MaxConnAge:    time.Hour,
				PoolTimeout:   4 * time.Second,
				IdleTimeout:   5 * time.Minute,
				IdleCheckFreq: time.Minute,
				DialTimeout:   5 * time.Second,
				ReadTimeout:   3 * time.Second,
				WriteTimeout:  3 * time.Second,
			},
		},
		Auth: TestAuthConfig{
			JWTSecret:     "test-secret-key-for-jwt",
			TokenExpiry:   time.Hour,
			RefreshExpiry: 24 * time.Hour,
			BcryptCost:    4, // Lower cost for faster tests
		},
		Log: TestLogConfig{
			Level:  "debug",
			Format: "text",
			Output: "stdout",
		},
		Monitor: TestMonitorConfig{
			Enabled:    true,
			MetricsURL: "/metrics",
		},
	}
}

// SetupTestEnvironment sets up environment variables for testing.
func SetupTestEnvironment() {
	if err := os.Setenv("APP_ENV", "test"); err != nil {
		log.Printf("Failed to set APP_ENV: %v", err)
	}
	if err := os.Setenv("APP_DEBUG", "true"); err != nil {
		log.Printf("Failed to set APP_DEBUG: %v", err)
	}
}

// CleanupTestEnvironment cleans up environment variables after testing.
func CleanupTestEnvironment() {
	if err := os.Unsetenv("APP_ENV"); err != nil {
		log.Printf("Failed to unset APP_ENV: %v", err)
	}
	if err := os.Unsetenv("APP_DEBUG"); err != nil {
		log.Printf("Failed to unset APP_DEBUG: %v", err)
	}
}

// SkipIfShort skips the test if running in short mode.
func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}
}

// MockDatabaseManager creates a mock database manager for testing.
func MockDatabaseManager() interface{} {
	// This is a placeholder for a mock database manager
	// In real tests, we would use testcontainers or in-memory databases
	// Returns interface{} to avoid circular imports
	return nil
}
