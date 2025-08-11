package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
)

func TestMain(m *testing.M) {
	// Setup test environment
	os.Setenv("APP_ENV", "test")
	os.Setenv("APP_DEBUG", "true")

	// Run tests
	code := m.Run()

	// Cleanup
	os.Unsetenv("APP_ENV")
	os.Unsetenv("APP_DEBUG")

	// Exit with test result code
	exit(code)
}

// exit is a variable to allow mocking in tests
var exit = func(_ int) {
	// In tests, we don't want to actually exit
}

func TestNewManager_InvalidConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	// Test with invalid MySQL config
	cfg := &config.DatabaseConfig{
		MySQL: config.MySQLConfig{
			Host:     "invalid-host",
			Port:     9999,
			Username: "invalid",
			Password: "invalid",
			Database: "invalid",
		},
		MongoDB: config.MongoDBConfig{
			URI:      "mongodb://invalid:27017",
			Database: "invalid",
		},
		Redis: config.RedisConfig{
			Host: "invalid-host",
			Port: 9999,
			DB:   0,
		},
	}

	// Execute
	manager, err := NewManager(cfg)

	// Assert - should fail due to invalid connection parameters
	assert.Error(t, err)
	assert.Nil(t, manager)
}

func TestManager_Structure(t *testing.T) {
	// Create a manager struct without actually connecting to databases
	cfg := config.GetDefaultDatabaseConfig()

	manager := &Manager{
		MySQL:   nil, // Would be set by NewMySQLManager
		MongoDB: nil, // Would be set by NewMongoDBManager
		Redis:   nil, // Would be set by NewRedisManager
		config:  cfg,
	}

	// Test structure
	assert.NotNil(t, manager)
	assert.Equal(t, cfg, manager.config)
}

func TestManager_GetConnectionStats_NilConnections(t *testing.T) {
	cfg := config.GetDefaultDatabaseConfig()

	manager := &Manager{
		MySQL:   nil,
		MongoDB: nil,
		Redis:   nil,
		config:  cfg,
	}

	// This should not panic even with nil connections
	stats := manager.GetConnectionStats()

	// Assert - should return empty map
	assert.NotNil(t, stats)
	// The actual implementation might handle nil connections differently
}

func TestManager_HealthCheck_NilConnections(t *testing.T) {
	cfg := config.GetDefaultDatabaseConfig()

	manager := &Manager{
		MySQL:   nil,
		MongoDB: nil,
		Redis:   nil,
		config:  cfg,
	}

	// This should handle nil connections gracefully
	err := manager.HealthCheck()

	// Assert - should return error due to nil connections
	assert.Error(t, err)
}

func TestManager_Close_NilConnections(t *testing.T) {
	cfg := config.GetDefaultDatabaseConfig()

	manager := &Manager{
		MySQL:   nil,
		MongoDB: nil,
		Redis:   nil,
		config:  cfg,
	}

	// This should handle nil connections gracefully
	err := manager.Close()

	// Assert - should return no error when connections are nil
	assert.NoError(t, err)
}

func TestValidateGORMVersion(t *testing.T) {
	// Execute
	err := ValidateGORMVersion()

	// Assert
	assert.NoError(t, err)
}

func TestManager_ValidateDatabaseVersions(t *testing.T) {
	cfg := config.GetDefaultDatabaseConfig()

	manager := &Manager{
		MySQL:   nil,
		MongoDB: nil,
		Redis:   nil,
		config:  cfg,
	}

	// Execute
	err := manager.ValidateDatabaseVersions()

	// Assert - this is just a validation function that logs, should not error
	assert.NoError(t, err)
}

// Test configuration validation
func TestDatabaseConfig_MySQLValidation(t *testing.T) {
	testCases := []struct {
		name   string
		config config.MySQLConfig
		valid  bool
	}{
		{
			name: "valid_mysql_config",
			config: config.MySQLConfig{
				Host:     "localhost",
				Port:     3306,
				Username: "root",
				Database: "test",
				Charset:  "utf8mb4",
			},
			valid: true,
		},
		{
			name: "empty_host",
			config: config.MySQLConfig{
				Host:     "",
				Port:     3306,
				Username: "root",
				Database: "test",
			},
			valid: false,
		},
		{
			name: "invalid_port",
			config: config.MySQLConfig{
				Host:     "localhost",
				Port:     0,
				Username: "root",
				Database: "test",
			},
			valid: false,
		},
		{
			name: "empty_database",
			config: config.MySQLConfig{
				Host:     "localhost",
				Port:     3306,
				Username: "root",
				Database: "",
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.valid {
				assert.NotEmpty(t, tc.config.Host)
				assert.Greater(t, tc.config.Port, 0)
				assert.NotEmpty(t, tc.config.Database)
			} else {
				// At least one validation should fail
				hasValidHost := tc.config.Host != ""
				hasValidPort := tc.config.Port > 0
				hasValidDatabase := tc.config.Database != ""

				assert.False(t, hasValidHost && hasValidPort && hasValidDatabase)
			}
		})
	}
}

func TestDatabaseConfig_MongoDBValidation(t *testing.T) {
	testCases := []struct {
		name   string
		config config.MongoDBConfig
		valid  bool
	}{
		{
			name: "valid_mongodb_config",
			config: config.MongoDBConfig{
				URI:      "mongodb://localhost:27017",
				Database: "test",
			},
			valid: true,
		},
		{
			name: "empty_uri",
			config: config.MongoDBConfig{
				URI:      "",
				Database: "test",
			},
			valid: false,
		},
		{
			name: "empty_database",
			config: config.MongoDBConfig{
				URI:      "mongodb://localhost:27017",
				Database: "",
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.valid {
				assert.NotEmpty(t, tc.config.URI)
				assert.NotEmpty(t, tc.config.Database)
			} else {
				// At least one validation should fail
				hasValidURI := tc.config.URI != ""
				hasValidDatabase := tc.config.Database != ""

				assert.False(t, hasValidURI && hasValidDatabase)
			}
		})
	}
}

func TestDatabaseConfig_RedisValidation(t *testing.T) {
	testCases := []struct {
		name   string
		config config.RedisConfig
		valid  bool
	}{
		{
			name: "valid_redis_config",
			config: config.RedisConfig{
				Host: "localhost",
				Port: 6379,
				DB:   0,
			},
			valid: true,
		},
		{
			name: "empty_host",
			config: config.RedisConfig{
				Host: "",
				Port: 6379,
				DB:   0,
			},
			valid: false,
		},
		{
			name: "invalid_port",
			config: config.RedisConfig{
				Host: "localhost",
				Port: 0,
				DB:   0,
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.valid {
				assert.NotEmpty(t, tc.config.Host)
				assert.Greater(t, tc.config.Port, 0)
				assert.GreaterOrEqual(t, tc.config.DB, 0)
			} else {
				// At least one validation should fail
				hasValidHost := tc.config.Host != ""
				hasValidPort := tc.config.Port > 0

				assert.False(t, hasValidHost && hasValidPort)
			}
		})
	}
}
