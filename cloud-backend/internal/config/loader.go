// Package config provides configuration loading and management.
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config represents the complete application configuration.
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Log      LogConfig      `mapstructure:"log"`
	Monitor  MonitorConfig  `mapstructure:"monitor"`
}

// AppConfig contains application-level configuration.
type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	Debug       bool   `mapstructure:"debug"`
}

// ServerConfig contains HTTP server configuration.
type ServerConfig struct {
	Host         string `mapstructure:"host" json:"host"`
	Port         string `mapstructure:"port" json:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout" json:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout" json:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout" json:"idle_timeout"`
}

// StorageConfig contains file storage configuration.
type StorageConfig struct {
	Type      string      `mapstructure:"type"`
	Local     LocalConfig `mapstructure:"local"`
	ChunkSize int64       `mapstructure:"chunk_size"`
}

// LocalConfig contains local storage configuration.
type LocalConfig struct {
	Path string `mapstructure:"path"`
}

// AuthConfig contains authentication configuration.
type AuthConfig struct {
	JWTSecret     string        `mapstructure:"jwt_secret"`
	TokenExpiry   time.Duration `mapstructure:"token_expiry"`
	RefreshExpiry time.Duration `mapstructure:"refresh_expiry"`
	BCryptCost    int           `mapstructure:"bcrypt_cost"`
}

// LogConfig contains logging configuration.
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// MonitorConfig contains monitoring configuration.
type MonitorConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	MetricsURL string `mapstructure:"metrics_url"`
}

// LoadConfig loads configuration from multiple sources with priority:
// 1. Environment variables
// 2. Configuration files (development.yaml, production.yaml, etc.)
// 3. Default values
func LoadConfig() (*Config, error) {
	v := viper.New()

	// Set configuration file paths and names
	v.SetConfigName("development") // default config file
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath("../configs")
	v.AddConfigPath(".")

	// Override config file based on environment
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	switch env {
	case "production":
		v.SetConfigName("production")
	case "testing":
		v.SetConfigName("testing")
	case "staging":
		v.SetConfigName("staging")
	default:
		v.SetConfigName("development")
	}

	// Enable environment variable binding
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values
	setDefaultValues(v)

	// Read configuration file
	if err := v.ReadInConfig(); err != nil {
		// If config file is not found, use environment variables and defaults
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Parse configuration into struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// setDefaultValues sets default configuration values.
func setDefaultValues(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "ycg-cloud")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.debug", true)

	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", ":8080")
	v.SetDefault("server.read_timeout", 30)
	v.SetDefault("server.write_timeout", 30)
	v.SetDefault("server.idle_timeout", 120)

	// Database defaults - using provided connection information
	v.SetDefault("database.mysql.host", "dbconn.sealosbja.site")
	v.SetDefault("database.mysql.port", 36973)
	v.SetDefault("database.mysql.username", "root")
	v.SetDefault("database.mysql.password", "hxl202653")
	v.SetDefault("database.mysql.database", "ycgcloud")
	v.SetDefault("database.mysql.charset", "utf8mb4")
	v.SetDefault("database.mysql.parse_time", true)
	v.SetDefault("database.mysql.loc", "Local")
	v.SetDefault("database.mysql.timeout", "10s")
	v.SetDefault("database.mysql.read_timeout", "30s")
	v.SetDefault("database.mysql.write_timeout", "30s")
	v.SetDefault("database.mysql.max_open_conns", 100)
	v.SetDefault("database.mysql.max_idle_conns", 10)
	v.SetDefault("database.mysql.conn_max_lifetime", "1h")
	v.SetDefault("database.mysql.conn_max_idle_time", "10m")

	// MongoDB defaults
	v.SetDefault("database.mongodb.uri", "mongodb://root:hxl202653@dbconn.sealosbja.site:42033/?directConnection=true")
	v.SetDefault("database.mongodb.database", "ycgcloud")
	v.SetDefault("database.mongodb.username", "root")
	v.SetDefault("database.mongodb.password", "hxl202653")
	v.SetDefault("database.mongodb.connect_timeout", "10s")
	v.SetDefault("database.mongodb.socket_timeout", "30s")
	v.SetDefault("database.mongodb.server_sel_timeout", "30s")
	v.SetDefault("database.mongodb.max_pool_size", 100)
	v.SetDefault("database.mongodb.min_pool_size", 5)
	v.SetDefault("database.mongodb.max_idle_time_ms", 300000)

	// Redis defaults
	v.SetDefault("database.redis.host", "dbconn.sealosbja.site")
	v.SetDefault("database.redis.port", 49169)
	v.SetDefault("database.redis.password", "jdsj6j67")
	v.SetDefault("database.redis.db", 0)
	v.SetDefault("database.redis.username", "default")
	v.SetDefault("database.redis.pool_size", 20)
	v.SetDefault("database.redis.min_idle_conns", 5)
	v.SetDefault("database.redis.max_conn_age", "1h")
	v.SetDefault("database.redis.pool_timeout", "4s")
	v.SetDefault("database.redis.idle_timeout", "5m")
	v.SetDefault("database.redis.idle_check_freq", "1m")
	v.SetDefault("database.redis.dial_timeout", "5s")
	v.SetDefault("database.redis.read_timeout", "3s")
	v.SetDefault("database.redis.write_timeout", "3s")

	// Storage defaults
	v.SetDefault("storage.type", "local")
	v.SetDefault("storage.chunk_size", 5242880)
	v.SetDefault("storage.local.path", "./storage")

	// Auth defaults
	v.SetDefault("auth.jwt_secret", "your-secret-key-change-in-production")
	v.SetDefault("auth.token_expiry", "24h")
	v.SetDefault("auth.refresh_expiry", "168h")
	v.SetDefault("auth.bcrypt_cost", 12)

	// Log defaults
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output", "stdout")

	// Monitor defaults
	v.SetDefault("monitor.enabled", true)
	v.SetDefault("monitor.metrics_url", "/metrics")
}

// validateConfig validates the loaded configuration.
func validateConfig(config *Config) error {
	// Validate required fields
	if config.App.Name == "" {
		return fmt.Errorf("app.name is required")
	}

	if config.Database.MySQL.Host == "" {
		return fmt.Errorf("database.mysql.host is required")
	}

	if config.Database.MongoDB.URI == "" {
		return fmt.Errorf("database.mongodb.uri is required")
	}

	if config.Database.Redis.Host == "" {
		return fmt.Errorf("database.redis.host is required")
	}

	if config.Auth.JWTSecret == "" || config.Auth.JWTSecret == "your-secret-key-change-in-production" {
		if config.App.Environment == "production" {
			return fmt.Errorf("auth.jwt_secret must be set for production environment")
		}
	}

	// Validate database connections can be established
	// This would be done in the database managers during initialization

	return nil
}

// GetDatabaseConfig returns only the database configuration portion.
func (c *Config) GetDatabaseConfig() DatabaseConfig {
	return c.Database
}
