// Package migration provides test utilities for migration testing.
package migration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

// TestDatabase holds test database connections.
type TestDatabase struct {
	MySQL   *gorm.DB
	MongoDB *mongo.Database
	TempDir string
	cleanup func()
}

// SetupTestDatabase creates a test database setup.
func SetupTestDatabase(t *testing.T) *TestDatabase {
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "migration_test_*")
	require.NoError(t, err)

	// For testing, we'll use a mock or skip actual database connections
	// This is a simplified test setup without actual database connections
	var sqliteDB *gorm.DB // Will be nil for simplified testing

	// Setup MongoDB test connection
	var mongoDB *mongo.Database
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err == nil {
		// Test if MongoDB is available
		err = mongoClient.Ping(context.Background(), nil)
		if err == nil {
			mongoDB = mongoClient.Database("migration_test_" + fmt.Sprintf("%d", time.Now().UnixNano()))
		}
	}

	cleanup := func() {
		// Clean up MongoDB
		if mongoDB != nil && mongoClient != nil {
			_ = mongoDB.Drop(context.Background())
			_ = mongoClient.Disconnect(context.Background())
		}

		// Clean up temp directory
		_ = os.RemoveAll(tempDir)
	}

	return &TestDatabase{
		MySQL:   sqliteDB,
		MongoDB: mongoDB,
		TempDir: tempDir,
		cleanup: cleanup,
	}
}

// Cleanup cleans up test resources.
func (td *TestDatabase) Cleanup() {
	if td.cleanup != nil {
		td.cleanup()
	}
}

// CreateTestMigration creates a test migration file.
func CreateTestMigration(t *testing.T, dir, version, name, upSQL, downSQL string) {
	// Create migration directory if it doesn't exist
	migrationDir := filepath.Join(dir, "mysql")
	err := os.MkdirAll(migrationDir, 0o750)
	require.NoError(t, err)

	// Create up migration file
	upFile := filepath.Join(migrationDir, fmt.Sprintf("%s_%s.up.sql", version, name))
	err = os.WriteFile(upFile, []byte(upSQL), 0o600)
	require.NoError(t, err)

	// Create down migration file
	downFile := filepath.Join(migrationDir, fmt.Sprintf("%s_%s.down.sql", version, name))
	err = os.WriteFile(downFile, []byte(downSQL), 0o600)
	require.NoError(t, err)
}

// CreateTestMongoMigration creates a test MongoDB migration file.
func CreateTestMongoMigration(t *testing.T, dir, version, name, script string) {
	// Create migration directory if it doesn't exist
	migrationDir := filepath.Join(dir, "mongodb")
	err := os.MkdirAll(migrationDir, 0o750)
	require.NoError(t, err)

	// Create migration file
	migrationFile := filepath.Join(migrationDir, fmt.Sprintf("%s_%s.js", version, name))
	err = os.WriteFile(migrationFile, []byte(script), 0600)
	require.NoError(t, err)
}

// TestMigrationSet provides a set of test migrations.
type TestMigrationSet struct {
	Migrations []TestMigrationData
}

// TestMigrationData represents test migration data.
type TestMigrationData struct {
	Version string
	Name    string
	UpSQL   string
	DownSQL string
}

// GetBasicTestMigrations returns a basic set of test migrations.
func GetBasicTestMigrations() *TestMigrationSet {
	return &TestMigrationSet{
		Migrations: []TestMigrationData{
			{
				Version: "20240101120000",
				Name:    "create_users_table",
				UpSQL: `
CREATE TABLE users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username VARCHAR(50) UNIQUE NOT NULL,
	email VARCHAR(100) UNIQUE NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
`,
				DownSQL: `
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP TABLE IF EXISTS users;
`,
			},
			{
				Version: "20240101130000",
				Name:    "create_posts_table",
				UpSQL: `
CREATE TABLE posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	title VARCHAR(255) NOT NULL,
	content TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_created_at ON posts(created_at);
`,
				DownSQL: `
DROP INDEX IF EXISTS idx_posts_created_at;
DROP INDEX IF EXISTS idx_posts_user_id;
DROP TABLE IF EXISTS posts;
`,
			},
			{
				Version: "20240101140000",
				Name:    "add_users_status_column",
				UpSQL: `
ALTER TABLE users ADD COLUMN status VARCHAR(20) DEFAULT 'active';
CREATE INDEX idx_users_status ON users(status);
`,
				DownSQL: `
DROP INDEX IF EXISTS idx_users_status;
ALTER TABLE users DROP COLUMN status;
`,
			},
		},
	}
}

// CreateTestMigrations creates test migrations in the specified directory.
func (tms *TestMigrationSet) CreateTestMigrations(t *testing.T, dir string) {
	for _, migration := range tms.Migrations {
		CreateTestMigration(t, dir, migration.Version, migration.Name, migration.UpSQL, migration.DownSQL)
	}
}

// TableExists checks if a table exists in the database.
func TableExists(db *gorm.DB, tableName string) bool {
	var count int64
	// For SQLite
	db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", tableName).Scan(&count)
	return count > 0
}

// ColumnExists checks if a column exists in a table.
func ColumnExists(db *gorm.DB, tableName, columnName string) bool {
	var count int64
	// For SQLite
	db.Raw("SELECT COUNT(*) FROM pragma_table_info(?) WHERE name=?", tableName, columnName).Scan(&count)
	return count > 0
}

// IndexExists checks if an index exists.
func IndexExists(db *gorm.DB, indexName string) bool {
	var count int64
	// For SQLite
	db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name=?", indexName).Scan(&count)
	return count > 0
}

// GetAppliedMigrations returns the list of applied migrations.
func GetAppliedMigrations(db *gorm.DB) ([]MySQLMigrationRecord, error) {
	var migrations []MySQLMigrationRecord
	err := db.Order("version ASC").Find(&migrations).Error
	return migrations, err
}

// ClearMigrationTable clears the migration table.
func ClearMigrationTable(db *gorm.DB) error {
	return db.Exec("DELETE FROM schema_migrations").Error
}

// AssertTableStructure asserts that a table has the expected structure.
func AssertTableStructure(t *testing.T, db *gorm.DB, tableName string, expectedColumns []string) {
	require.True(t, TableExists(db, tableName), "Table %s should exist", tableName)

	for _, column := range expectedColumns {
		require.True(t, ColumnExists(db, tableName, column), "Column %s should exist in table %s", column, tableName)
	}
}

// AssertMigrationRecord asserts that a migration record exists.
func AssertMigrationRecord(t *testing.T, db *gorm.DB, version, name string) {
	var record MySQLMigrationRecord
	err := db.Where("version = ? AND name = ?", version, name).First(&record).Error
	require.NoError(t, err, "Migration record should exist for version %s", version)
	require.Equal(t, version, record.Version)
	require.Equal(t, name, record.Name)
}

// AssertMigrationNotExists asserts that a migration record does not exist.
func AssertMigrationNotExists(t *testing.T, db *gorm.DB, version string) {
	var count int64
	err := db.Model(&MySQLMigrationRecord{}).Where("version = ?", version).Count(&count).Error
	require.NoError(t, err)
	require.Equal(t, int64(0), count, "Migration record should not exist for version %s", version)
}

// WaitForDatabase waits for database to be ready.
func WaitForDatabase(db *gorm.DB, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}

		if err := sqlDB.Ping(); err == nil {
			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("database not ready within timeout")
}

// MockMigrationFile represents a mock migration file for testing.
type MockMigrationFile struct {
	Path    string
	Content string
}

// CreateMockMigrationFiles creates mock migration files for testing.
func CreateMockMigrationFiles(t *testing.T, baseDir string, files []MockMigrationFile) {
	for _, file := range files {
		fullPath := filepath.Join(baseDir, file.Path)

		// Create directory if needed
		dir := filepath.Dir(fullPath)
		err := os.MkdirAll(dir, 0o750)
		require.NoError(t, err)

		// Write file
		err = os.WriteFile(fullPath, []byte(file.Content), 0o600)
		require.NoError(t, err)
	}
}

// TestRunnerConfig holds configuration for test runner.
type TestRunnerConfig struct {
	Verbose bool
	DryRun  bool
}

// DefaultTestRunnerConfig returns default test runner configuration.
func DefaultTestRunnerConfig() *TestRunnerConfig {
	return &TestRunnerConfig{
		Verbose: false,
		DryRun:  false,
	}
}

// ExecuteSQL executes raw SQL and returns the result.
func ExecuteSQL(db *gorm.DB, query string, args ...interface{}) error {
	return db.Exec(query, args...).Error
}

// GetTableCount returns the number of rows in a table.
func GetTableCount(db *gorm.DB, tableName string) (int64, error) {
	// Validate table name to prevent SQL injection
	if !isValidTableName(tableName) {
		return 0, fmt.Errorf("invalid table name: %s", tableName)
	}

	var count int64
	// Note: Table names cannot be parameterized in SQL, so we validate above
	err := db.Raw("SELECT COUNT(*) FROM " + tableName).Scan(&count).Error
	return count, err
}

// isValidTableName validates table name to prevent SQL injection
func isValidTableName(tableName string) bool {
	// Allow only alphanumeric characters and underscores
	for _, char := range tableName {
		if (char < 'a' || char > 'z') &&
			(char < 'A' || char > 'Z') &&
			(char < '0' || char > '9') &&
			char != '_' {
			return false
		}
	}
	return tableName != "" && len(tableName) <= 64 // MySQL table name limit
}

// BackupDatabase creates a backup of the test database.
func BackupDatabase(_ *gorm.DB, _ string) error {
	// For testing purposes, this is a mock implementation
	return fmt.Errorf("backup not implemented for test database")
}

// RestoreDatabase restores the test database from backup.
func RestoreDatabase(_ *gorm.DB, _ string) error {
	// For testing purposes, this is a mock implementation
	return fmt.Errorf("restore not implemented for test database")
}
