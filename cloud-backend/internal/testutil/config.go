// Package testutil provides test configuration utilities.
package testutil

import (
	"time"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
)

// TestConfig creates a test configuration.
func TestConfig() *config.Config {
	return &config.Config{
		App: config.AppConfig{
			Name:        "ycg-cloud-test",
			Version:     "1.0.0-test",
			Environment: "test",
			Debug:       true,
		},
		Server: config.ServerConfig{
			Host:         "localhost",
			Port:         ":8081",
			ReadTimeout:  30,
			WriteTimeout: 30,
			IdleTimeout:  120,
		},
		Database: config.DatabaseConfig{
			MySQL: config.MySQLConfig{
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
			MongoDB: config.MongoDBConfig{
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
			Redis: config.RedisConfig{
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
		Auth: config.AuthConfig{
			JWTSecret:     "test-secret-key-for-jwt",
			TokenExpiry:   time.Hour,
			RefreshExpiry: 24 * time.Hour,
			BCryptCost:    4, // Lower cost for faster tests
		},
		Log: config.LogConfig{
			Level:  "debug",
			Format: "text",
			Output: "stdout",
		},
		Monitor: config.MonitorConfig{
			Enabled:    true,
			MetricsURL: "/metrics",
		},
	}
}
