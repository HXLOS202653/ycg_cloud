// Package database provides MySQL database connection with GORM 1.25.12.
package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
)

// MySQLManager manages MySQL database connections.
type MySQLManager struct {
	db     *gorm.DB
	config config.MySQLConfig
}

// NewMySQLConnection creates a new MySQL connection with GORM 1.25.12.
func NewMySQLConnection(cfg *config.MySQLConfig) (*gorm.DB, error) {
	// Build DSN (Data Source Name)
	dsn := buildMySQLDSN(cfg)

	// Configure GORM logger
	gormLogger := logger.Default
	if cfg.Database != "" {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger:                                   gormLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   false,
		PrepareStmt:                              true,

		// Naming strategy for table and column names
		NamingStrategy: nil, // Use default naming strategy
	}

	// Open database connection
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool settings
	configureConnectionPool(sqlDB, cfg)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL database: %w", err)
	}

	log.Printf("MySQL database connected successfully to %s:%d/%s", cfg.Host, cfg.Port, cfg.Database)
	return db, nil
}

// buildMySQLDSN builds MySQL Data Source Name.
func buildMySQLDSN(cfg *config.MySQLConfig) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
		cfg.ParseTime,
		cfg.Loc,
	)

	// Add timeout parameters
	if cfg.Timeout > 0 {
		dsn += fmt.Sprintf("&timeout=%s", cfg.Timeout.String())
	}
	if cfg.ReadTimeout > 0 {
		dsn += fmt.Sprintf("&readTimeout=%s", cfg.ReadTimeout.String())
	}
	if cfg.WriteTimeout > 0 {
		dsn += fmt.Sprintf("&writeTimeout=%s", cfg.WriteTimeout.String())
	}

	return dsn
}

// configureConnectionPool configures MySQL connection pool settings.
func configureConnectionPool(sqlDB *sql.DB, cfg *config.MySQLConfig) {
	// Set maximum number of open connections
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	} else {
		sqlDB.SetMaxOpenConns(100) // Default value
	}

	// Set maximum number of idle connections
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	} else {
		sqlDB.SetMaxIdleConns(10) // Default value
	}

	// Set maximum lifetime of connections
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	} else {
		sqlDB.SetConnMaxLifetime(time.Hour) // Default 1 hour
	}

	// Set maximum idle time of connections
	if cfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	} else {
		sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Default 10 minutes
	}

	log.Printf("MySQL connection pool configured: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%v, MaxIdleTime=%v",
		cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.ConnMaxLifetime, cfg.ConnMaxIdleTime)
}

// NewMySQLManager creates a new MySQL manager.
func NewMySQLManager(cfg *config.MySQLConfig) (*MySQLManager, error) {
	db, err := NewMySQLConnection(cfg)
	if err != nil {
		return nil, err
	}

	return &MySQLManager{
		db:     db,
		config: *cfg,
	}, nil
}

// GetDB returns the GORM database instance.
func (m *MySQLManager) GetDB() *gorm.DB {
	return m.db
}

// Close closes the database connection.
func (m *MySQLManager) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Health checks the database health.
func (m *MySQLManager) Health() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// GetStats returns database connection statistics.
func (m *MySQLManager) GetStats() interface{} {
	sqlDB, err := m.db.DB()
	if err != nil {
		return nil
	}
	return sqlDB.Stats()
}
