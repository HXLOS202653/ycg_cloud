// Package migration provides basic tests for migration functionality.
package migration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMySQLMigrationRecord_Fields(t *testing.T) {
	now := time.Now()
	record := MySQLMigrationRecord{
		ID:        1,
		Version:   "20240101120000",
		Name:      "create_users_table",
		UpSQL:     "CREATE TABLE users (id INT PRIMARY KEY);",
		DownSQL:   "DROP TABLE users;",
		AppliedAt: now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, uint(1), record.ID)
	assert.Equal(t, "20240101120000", record.Version)
	assert.Equal(t, "create_users_table", record.Name)
	assert.Equal(t, "CREATE TABLE users (id INT PRIMARY KEY);", record.UpSQL)
	assert.Equal(t, "DROP TABLE users;", record.DownSQL)
	assert.Equal(t, now, record.AppliedAt)
	assert.Equal(t, now, record.CreatedAt)
	assert.Equal(t, now, record.UpdatedAt)
}

func TestMySQLMigration_Fields(t *testing.T) {
	migration := MySQLMigration{
		Version:  "20240101120000",
		Name:     "create_users_table",
		UpFile:   "/path/to/20240101120000_create_users_table.up.sql",
		DownFile: "/path/to/20240101120000_create_users_table.down.sql",
		Applied:  false,
	}

	assert.Equal(t, "20240101120000", migration.Version)
	assert.Equal(t, "create_users_table", migration.Name)
	assert.Equal(t, "/path/to/20240101120000_create_users_table.up.sql", migration.UpFile)
	assert.Equal(t, "/path/to/20240101120000_create_users_table.down.sql", migration.DownFile)
	assert.False(t, migration.Applied)
}

func TestMongoDBMigrationRecord_Fields(t *testing.T) {
	now := time.Now()
	record := MongoDBMigrationRecord{
		ID:        "507f1f77bcf86cd799439011",
		Version:   "20240101120000",
		Name:      "create_users_collection",
		Script:    "db.users.createIndex({username: 1});",
		AppliedAt: now,
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.Equal(t, "507f1f77bcf86cd799439011", record.ID)
	assert.Equal(t, "20240101120000", record.Version)
	assert.Equal(t, "create_users_collection", record.Name)
	assert.Equal(t, "db.users.createIndex({username: 1});", record.Script)
	assert.Equal(t, now, record.AppliedAt)
	assert.Equal(t, now, record.CreatedAt)
	assert.Equal(t, now, record.UpdatedAt)
}

func TestMongoDBMigration_Fields(t *testing.T) {
	migration := MongoDBMigration{
		Version: "20240101120000",
		Name:    "create_users_collection",
		Applied: true,
	}

	assert.Equal(t, "20240101120000", migration.Version)
	assert.Equal(t, "create_users_collection", migration.Name)
	assert.True(t, migration.Applied)
}

func TestCreateTestMigration(t *testing.T) {
	testDB := SetupTestDatabase(t)
	defer testDB.Cleanup()

	version := "20240101120000"
	name := "test_migration"
	upSQL := "CREATE TABLE test (id INT PRIMARY KEY);"
	downSQL := "DROP TABLE test;"

	CreateTestMigration(t, testDB.TempDir, version, name, upSQL, downSQL)

	// Check files were created
	upFile := filepath.Join(testDB.TempDir, "mysql", version+"_"+name+".up.sql")
	downFile := filepath.Join(testDB.TempDir, "mysql", version+"_"+name+".down.sql")

	assert.FileExists(t, upFile)
	assert.FileExists(t, downFile)

	// Check file contents
	upContent, err := os.ReadFile(upFile)
	require.NoError(t, err)
	assert.Equal(t, upSQL, string(upContent))

	downContent, err := os.ReadFile(downFile)
	require.NoError(t, err)
	assert.Equal(t, downSQL, string(downContent))
}

func TestCreateTestMongoMigration(t *testing.T) {
	testDB := SetupTestDatabase(t)
	defer testDB.Cleanup()

	version := "20240101120000"
	name := "test_mongo_migration"
	script := "db.users.createIndex({username: 1});"

	CreateTestMongoMigration(t, testDB.TempDir, version, name, script)

	// Check file was created
	file := filepath.Join(testDB.TempDir, "mongodb", version+"_"+name+".js")
	assert.FileExists(t, file)

	// Check file content
	content, err := os.ReadFile(file)
	require.NoError(t, err)
	assert.Equal(t, script, string(content))
}

func TestGetBasicTestMigrations(t *testing.T) {
	testMigrations := GetBasicTestMigrations()

	assert.NotNil(t, testMigrations)
	assert.Len(t, testMigrations.Migrations, 3)

	// Check first migration
	first := testMigrations.Migrations[0]
	assert.Equal(t, "20240101120000", first.Version)
	assert.Equal(t, "create_users_table", first.Name)
	assert.Contains(t, first.UpSQL, "CREATE TABLE users")
	assert.Contains(t, first.DownSQL, "DROP TABLE")
}

func TestConfig_Fields(t *testing.T) {
	config := Config{
		MigrationsDir: "/path/to/migrations",
		MySQL: &MySQLMigrationConfig{
			Enabled:        true,
			DSN:            "user:pass@tcp(localhost:3306)/db",
			Database:       "testdb",
			TablePrefix:    "app_",
			MigrationTable: "schema_migrations",
			Timeout:        30,
		},
		MongoDB: &MongoDBMigrationConfig{
			Enabled:             true,
			URI:                 "mongodb://localhost:27017",
			Database:            "testdb",
			MigrationCollection: "migrations",
			Timeout:             30,
		},
	}

	assert.Equal(t, "/path/to/migrations", config.MigrationsDir)
	assert.True(t, config.MySQL.Enabled)
	assert.Equal(t, "testdb", config.MySQL.Database)
	assert.True(t, config.MongoDB.Enabled)
	assert.Equal(t, "testdb", config.MongoDB.Database)
}

func TestValidationConfig_Fields(t *testing.T) {
	config := ValidationConfig{
		Enabled:             true,
		StrictMode:          true,
		AllowedOperations:   []string{"CREATE", "ALTER", "DROP"},
		ForbiddenOperations: []string{"TRUNCATE", "DELETE"},
		RequireComments:     true,
		RequireBackup:       true,
	}

	assert.True(t, config.Enabled)
	assert.True(t, config.StrictMode)
	assert.Contains(t, config.AllowedOperations, "CREATE")
	assert.Contains(t, config.ForbiddenOperations, "TRUNCATE")
	assert.True(t, config.RequireComments)
	assert.True(t, config.RequireBackup)
}

func TestBackupConfig_Fields(t *testing.T) {
	config := BackupConfig{
		Enabled:       true,
		BackupDir:     "/backups",
		RetentionDays: 30,
		Compression:   true,
		AutoBackup:    true,
	}

	assert.True(t, config.Enabled)
	assert.Equal(t, "/backups", config.BackupDir)
	assert.Equal(t, 30, config.RetentionDays)
	assert.True(t, config.Compression)
	assert.True(t, config.AutoBackup)
}

func TestPlan_Fields(t *testing.T) {
	now := time.Now()
	plan := Plan{
		Direction:     "up",
		TargetVersion: "20240101120000",
		Steps:         3,
		EstimatedTime: 5 * time.Minute,
		CreatedAt:     now,
	}

	assert.Equal(t, Direction("up"), plan.Direction)
	assert.Equal(t, "20240101120000", plan.TargetVersion)
	assert.Equal(t, 3, plan.Steps)
	assert.Equal(t, 5*time.Minute, plan.EstimatedTime)
	assert.Equal(t, now, plan.CreatedAt)
}

func TestMySQLPlan_Fields(t *testing.T) {
	plan := MySQLPlan{
		Migrations: []*MySQLMigration{
			{Version: "20240101120000", Name: "create_users"},
			{Version: "20240101130000", Name: "create_posts"},
		},
		TotalSteps:    2,
		EstimatedTime: 2 * time.Minute,
	}

	assert.Len(t, plan.Migrations, 2)
	assert.Equal(t, 2, plan.TotalSteps)
	assert.Equal(t, 2*time.Minute, plan.EstimatedTime)
}

func TestMongoDBPlan_Fields(t *testing.T) {
	plan := MongoDBPlan{
		Migrations: []*MongoDBMigration{
			{Version: "20240101120000", Name: "create_users_index"},
		},
		TotalSteps:    1,
		EstimatedTime: 1 * time.Minute,
	}

	assert.Len(t, plan.Migrations, 1)
	assert.Equal(t, 1, plan.TotalSteps)
	assert.Equal(t, 1*time.Minute, plan.EstimatedTime)
}

// Benchmark tests
func BenchmarkMySQLMigrationRecord_TableName(b *testing.B) {
	record := MySQLMigrationRecord{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		record.TableName()
	}
}

func BenchmarkCreateTestMigration(b *testing.B) {
	testDB := SetupTestDatabase(&testing.T{})
	defer testDB.Cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		version := time.Now().Format("20060102150405")
		name := "benchmark_migration"
		CreateTestMigration(&testing.T{}, testDB.TempDir, version, name,
			"CREATE TABLE test (id INT);", "DROP TABLE test;")
	}
}
