// Package config provides database configuration for the cloud storage system.
package config

import (
	"time"
)

// DatabaseConfig contains all database configurations.
type DatabaseConfig struct {
	MySQL    MySQLConfig   `mapstructure:"mysql" json:"mysql"`
	MongoDB  MongoDBConfig `mapstructure:"mongodb" json:"mongodb"`
	Redis    RedisConfig   `mapstructure:"redis" json:"redis"`
	MaxConns int           `mapstructure:"max_conns" json:"max_conns"`
	MaxIdle  int           `mapstructure:"max_idle" json:"max_idle"`
}

// MySQLConfig contains MySQL database configuration.
type MySQLConfig struct {
	Host         string        `mapstructure:"host" json:"host"`
	Port         int           `mapstructure:"port" json:"port"`
	Username     string        `mapstructure:"username" json:"username"`
	Password     string        `mapstructure:"password" json:"password"`
	Database     string        `mapstructure:"database" json:"database"`
	Charset      string        `mapstructure:"charset" json:"charset"`
	ParseTime    bool          `mapstructure:"parse_time" json:"parse_time"`
	Loc          string        `mapstructure:"loc" json:"loc"`
	Timeout      time.Duration `mapstructure:"timeout" json:"timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout"`

	// Connection pool settings
	MaxOpenConns    int           `mapstructure:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" json:"conn_max_idle_time"`
}

// MongoDBConfig contains MongoDB configuration.
type MongoDBConfig struct {
	URI              string        `mapstructure:"uri" json:"uri"`
	Database         string        `mapstructure:"database" json:"database"`
	Username         string        `mapstructure:"username" json:"username"`
	Password         string        `mapstructure:"password" json:"password"`
	ConnectTimeout   time.Duration `mapstructure:"connect_timeout" json:"connect_timeout"`
	SocketTimeout    time.Duration `mapstructure:"socket_timeout" json:"socket_timeout"`
	ServerSelTimeout time.Duration `mapstructure:"server_sel_timeout" json:"server_sel_timeout"`

	// Connection pool settings
	MaxPoolSize   uint64 `mapstructure:"max_pool_size" json:"max_pool_size"`
	MinPoolSize   uint64 `mapstructure:"min_pool_size" json:"min_pool_size"`
	MaxIdleTimeMS uint64 `mapstructure:"max_idle_time_ms" json:"max_idle_time_ms"`
}

// RedisConfig contains Redis configuration.
type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
	Username string `mapstructure:"username" json:"username"`

	// Connection pool settings
	PoolSize      int           `mapstructure:"pool_size" json:"pool_size"`
	MinIdleConns  int           `mapstructure:"min_idle_conns" json:"min_idle_conns"`
	MaxConnAge    time.Duration `mapstructure:"max_conn_age" json:"max_conn_age"`
	PoolTimeout   time.Duration `mapstructure:"pool_timeout" json:"pool_timeout"`
	IdleTimeout   time.Duration `mapstructure:"idle_timeout" json:"idle_timeout"`
	IdleCheckFreq time.Duration `mapstructure:"idle_check_freq" json:"idle_check_freq"`

	// Timeouts
	DialTimeout  time.Duration `mapstructure:"dial_timeout" json:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout"`
}

// GetDefaultDatabaseConfig returns default database configuration.
func GetDefaultDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		MySQL: MySQLConfig{
			Host:            "localhost",
			Port:            3306,
			Username:        "root",
			Password:        "",
			Database:        "ycgcloud",
			Charset:         "utf8mb4",
			ParseTime:       true,
			Loc:             "Local",
			Timeout:         10 * time.Second,
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			MaxOpenConns:    100,
			MaxIdleConns:    10,
			ConnMaxLifetime: time.Hour,
			ConnMaxIdleTime: 10 * time.Minute,
		},
		MongoDB: MongoDBConfig{
			URI:              "mongodb://localhost:27017",
			Database:         "ycgcloud",
			Username:         "",
			Password:         "",
			ConnectTimeout:   10 * time.Second,
			SocketTimeout:    30 * time.Second,
			ServerSelTimeout: 30 * time.Second,
			MaxPoolSize:      100,
			MinPoolSize:      5,
			MaxIdleTimeMS:    300000, // 5 minutes
		},
		Redis: RedisConfig{
			Host:          "localhost",
			Port:          6379,
			Password:      "",
			DB:            0,
			Username:      "",
			PoolSize:      20,
			MinIdleConns:  5,
			MaxConnAge:    time.Hour,
			PoolTimeout:   4 * time.Second,
			IdleTimeout:   5 * time.Minute,
			IdleCheckFreq: time.Minute,
			DialTimeout:   5 * time.Second,
			ReadTimeout:   3 * time.Second,
			WriteTimeout:  3 * time.Second,
		},
		MaxConns: 100,
		MaxIdle:  10,
	}
}
