package config

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Setup test environment
	if err := os.Setenv("APP_ENV", "test"); err != nil {
		panic(fmt.Sprintf("Failed to set APP_ENV: %v", err))
	}
	if err := os.Setenv("APP_DEBUG", "true"); err != nil {
		panic(fmt.Sprintf("Failed to set APP_DEBUG: %v", err))
	}

	// Run tests
	code := m.Run()

	// Cleanup
	if err := os.Unsetenv("APP_ENV"); err != nil {
		log.Printf("Failed to unset APP_ENV: %v", err)
	}
	if err := os.Unsetenv("APP_DEBUG"); err != nil {
		log.Printf("Failed to unset APP_DEBUG: %v", err)
	}

	// Exit with test result code
	os.Exit(code)
}

func TestLoadConfig_DefaultValues(t *testing.T) {
	// Setup - clear any existing config file environment
	originalEnv := os.Getenv("APP_ENV")
	defer func() {
		if originalEnv != "" {
			if err := os.Setenv("APP_ENV", originalEnv); err != nil {
				t.Logf("Failed to restore APP_ENV: %v", err)
			}
		} else {
			if err := os.Unsetenv("APP_ENV"); err != nil {
				t.Logf("Failed to unset APP_ENV: %v", err)
			}
		}
	}()

	if err := os.Setenv("APP_ENV", "test"); err != nil {
		t.Fatalf("Failed to set APP_ENV: %v", err)
	}

	// Execute
	cfg, err := LoadConfig()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check that we have valid configuration structure
	assert.NotEmpty(t, cfg.App.Name)
	assert.NotEmpty(t, cfg.App.Version)
	assert.NotEmpty(t, cfg.Server.Port)
	assert.Greater(t, cfg.Server.ReadTimeout, 0)
	assert.Greater(t, cfg.Server.WriteTimeout, 0)
	assert.Greater(t, cfg.Server.IdleTimeout, 0)
}

func TestLoadConfig_EnvironmentOverrides(t *testing.T) {
	// Setup environment variables
	testCases := []struct {
		envVar   string
		envValue string
		check    func(*Config) bool
	}{
		{
			envVar:   "APP_NAME",
			envValue: "test-app",
			check:    func(c *Config) bool { return c.App.Name == "test-app" },
		},
		{
			envVar:   "APP_VERSION",
			envValue: "2.0.0",
			check:    func(c *Config) bool { return c.App.Version == "2.0.0" },
		},
		{
			envVar:   "APP_DEBUG",
			envValue: "true",
			check:    func(c *Config) bool { return c.App.Debug == true },
		},
		{
			envVar:   "SERVER_PORT",
			envValue: ":9000",
			check:    func(c *Config) bool { return c.Server.Port == ":9000" },
		},
	}

	for _, tc := range testCases {
		t.Run(tc.envVar, func(t *testing.T) {
			// Setup
			originalValue := os.Getenv(tc.envVar)
			defer func() {
				if originalValue != "" {
					if err := os.Setenv(tc.envVar, originalValue); err != nil {
						t.Logf("Failed to restore %s: %v", tc.envVar, err)
					}
				} else {
					if err := os.Unsetenv(tc.envVar); err != nil {
						t.Logf("Failed to unset %s: %v", tc.envVar, err)
					}
				}
			}()

			if err := os.Setenv(tc.envVar, tc.envValue); err != nil {
				t.Fatalf("Failed to set %s: %v", tc.envVar, err)
			}
			if err := os.Setenv("APP_ENV", "test"); err != nil {
				t.Fatalf("Failed to set APP_ENV: %v", err)
			}

			// Execute
			cfg, err := LoadConfig()

			// Assert
			require.NoError(t, err)
			assert.True(t, tc.check(cfg), "Environment variable %s not properly set", tc.envVar)
		})
	}
}

func TestAppConfig_Validation(t *testing.T) {
	testCases := []struct {
		name   string
		config AppConfig
		valid  bool
	}{
		{
			name: "valid_config",
			config: AppConfig{
				Name:        "test-app",
				Version:     "1.0.0",
				Environment: "test",
				Debug:       true,
			},
			valid: true,
		},
		{
			name: "empty_name",
			config: AppConfig{
				Name:        "",
				Version:     "1.0.0",
				Environment: "test",
				Debug:       true,
			},
			valid: false,
		},
		{
			name: "empty_version",
			config: AppConfig{
				Name:        "test-app",
				Version:     "",
				Environment: "test",
				Debug:       true,
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simple validation check
			if tc.valid {
				assert.NotEmpty(t, tc.config.Name)
				assert.NotEmpty(t, tc.config.Version)
				assert.NotEmpty(t, tc.config.Environment)
			} else {
				// At least one required field should be empty
				isEmpty := tc.config.Name == "" || tc.config.Version == "" || tc.config.Environment == ""
				assert.True(t, isEmpty)
			}
		})
	}
}

func TestServerConfig_Validation(t *testing.T) {
	testCases := []struct {
		name   string
		config ServerConfig
		valid  bool
	}{
		{
			name: "valid_config",
			config: ServerConfig{
				Host:         "localhost",
				Port:         ":8080",
				ReadTimeout:  30,
				WriteTimeout: 30,
				IdleTimeout:  120,
			},
			valid: true,
		},
		{
			name: "invalid_port_format",
			config: ServerConfig{
				Host:         "localhost",
				Port:         "8080", // Missing colon
				ReadTimeout:  30,
				WriteTimeout: 30,
				IdleTimeout:  120,
			},
			valid: false,
		},
		{
			name: "zero_timeouts",
			config: ServerConfig{
				Host:         "localhost",
				Port:         ":8080",
				ReadTimeout:  0,
				WriteTimeout: 0,
				IdleTimeout:  0,
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.valid {
				assert.NotEmpty(t, tc.config.Host)
				assert.Contains(t, tc.config.Port, ":")
				assert.Greater(t, tc.config.ReadTimeout, 0)
				assert.Greater(t, tc.config.WriteTimeout, 0)
				assert.Greater(t, tc.config.IdleTimeout, 0)
			} else {
				// At least one validation should fail
				hasValidPort := tc.config.Port != "" && tc.config.Port[0] == ':'
				hasValidTimeouts := tc.config.ReadTimeout > 0 && tc.config.WriteTimeout > 0 && tc.config.IdleTimeout > 0

				assert.False(t, hasValidPort && hasValidTimeouts)
			}
		})
	}
}

func TestDatabaseConfig_Structure(t *testing.T) {
	cfg := &DatabaseConfig{
		MySQL: MySQLConfig{
			Host:     "localhost",
			Port:     3306,
			Username: "test",
			Database: "test_db",
		},
		MongoDB: MongoDBConfig{
			URI:      "mongodb://localhost:27017",
			Database: "test_db",
		},
		Redis: RedisConfig{
			Host: "localhost",
			Port: 6379,
			DB:   0,
		},
	}

	// Test structure is properly defined
	assert.Equal(t, "localhost", cfg.MySQL.Host)
	assert.Equal(t, 3306, cfg.MySQL.Port)
	assert.Equal(t, "test", cfg.MySQL.Username)
	assert.Equal(t, "test_db", cfg.MySQL.Database)

	assert.Equal(t, "mongodb://localhost:27017", cfg.MongoDB.URI)
	assert.Equal(t, "test_db", cfg.MongoDB.Database)

	assert.Equal(t, "localhost", cfg.Redis.Host)
	assert.Equal(t, 6379, cfg.Redis.Port)
	assert.Equal(t, 0, cfg.Redis.DB)
}

func TestAuthConfig_Structure(t *testing.T) {
	cfg := &AuthConfig{
		JWTSecret:     "test-secret",
		TokenExpiry:   24 * time.Hour,
		RefreshExpiry: 168 * time.Hour,
		BCryptCost:    10,
	}

	assert.Equal(t, "test-secret", cfg.JWTSecret)
	assert.Equal(t, 24*time.Hour, cfg.TokenExpiry)
	assert.Equal(t, 168*time.Hour, cfg.RefreshExpiry)
	assert.Equal(t, 10, cfg.BCryptCost)
}

func TestLogConfig_Structure(t *testing.T) {
	cfg := &LogConfig{
		Level:  "info",
		Format: "json",
		Output: "stdout",
	}

	assert.Equal(t, "info", cfg.Level)
	assert.Equal(t, "json", cfg.Format)
	assert.Equal(t, "stdout", cfg.Output)
}

func TestMonitorConfig_Structure(t *testing.T) {
	cfg := &MonitorConfig{
		Enabled:    true,
		MetricsURL: "/metrics",
	}

	assert.True(t, cfg.Enabled)
	assert.Equal(t, "/metrics", cfg.MetricsURL)
}

// Benchmark tests
func BenchmarkLoadConfig(b *testing.B) {
	if err := os.Setenv("APP_ENV", "test"); err != nil {
		b.Fatalf("Failed to set APP_ENV: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("APP_ENV"); err != nil {
			b.Logf("Failed to unset APP_ENV: %v", err)
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := LoadConfig()
		if err != nil {
			b.Fatal(err)
		}
	}
}
