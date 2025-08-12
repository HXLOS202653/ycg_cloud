package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// MySQLMigrationRecord MySQL迁移记录
type MySQLMigrationRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Version   string    `gorm:"uniqueIndex;size:20;not null" json:"version"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	UpSQL     string    `gorm:"type:text" json:"up_sql"`
	DownSQL   string    `gorm:"type:text" json:"down_sql"`
	AppliedAt time.Time `gorm:"not null" json:"applied_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定MySQL迁移记录表名
func (MySQLMigrationRecord) TableName() string {
	return "schema_migrations"
}

// MongoDBMigrationRecord MongoDB迁移记录
type MongoDBMigrationRecord struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Version   string    `bson:"version" json:"version"`
	Name      string    `bson:"name" json:"name"`
	Script    string    `bson:"script" json:"script"`
	AppliedAt time.Time `bson:"applied_at" json:"applied_at"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// MySQLMigration MySQL迁移文件信息
type MySQLMigration struct {
	Version  string `json:"version"`
	Name     string `json:"name"`
	UpFile   string `json:"up_file"`
	DownFile string `json:"down_file"`
	Applied  bool   `json:"applied"`
}

// MongoDBMigration MongoDB迁移文件信息
type MongoDBMigration struct {
	Version string `json:"version"`
	Name    string `json:"name"`
	JSFile  string `json:"js_file"`
	Applied bool   `json:"applied"`
}

// Status 迁移状态
type Status struct {
	Database      string                  `json:"database"`
	MySQLStatus   *MySQLMigrationStatus   `json:"mysql_status"`
	MongoDBStatus *MongoDBMigrationStatus `json:"mongodb_status"`
	LastChecked   time.Time               `json:"last_checked"`
}

// MySQLMigrationStatus MySQL迁移状态
type MySQLMigrationStatus struct {
	TotalMigrations   int                   `json:"total_migrations"`
	AppliedMigrations int                   `json:"applied_migrations"`
	PendingMigrations int                   `json:"pending_migrations"`
	LastMigration     *MySQLMigrationRecord `json:"last_migration"`
	Migrations        []*MySQLMigration     `json:"migrations"`
}

// MongoDBMigrationStatus MongoDB迁移状态
type MongoDBMigrationStatus struct {
	TotalMigrations   int                     `json:"total_migrations"`
	AppliedMigrations int                     `json:"applied_migrations"`
	PendingMigrations int                     `json:"pending_migrations"`
	LastMigration     *MongoDBMigrationRecord `json:"last_migration"`
	Migrations        []*MongoDBMigration     `json:"migrations"`
}

// Config 迁移配置
type Config struct {
	MigrationsDir string                  `json:"migrations_dir" yaml:"migrations_dir"`
	MySQL         *MySQLMigrationConfig   `json:"mysql" yaml:"mysql"`
	MongoDB       *MongoDBMigrationConfig `json:"mongodb" yaml:"mongodb"`
	Validation    *ValidationConfig       `json:"validation" yaml:"validation"`
	Backup        *BackupConfig           `json:"backup" yaml:"backup"`
}

// MySQLMigrationConfig MySQL迁移配置
type MySQLMigrationConfig struct {
	Enabled        bool   `json:"enabled" yaml:"enabled"`
	DSN            string `json:"dsn" yaml:"dsn"`
	Database       string `json:"database" yaml:"database"`
	TablePrefix    string `json:"table_prefix" yaml:"table_prefix"`
	MigrationTable string `json:"migration_table" yaml:"migration_table"`
	Timeout        int    `json:"timeout" yaml:"timeout"` // 秒
}

// MongoDBMigrationConfig MongoDB迁移配置
type MongoDBMigrationConfig struct {
	Enabled             bool   `json:"enabled" yaml:"enabled"`
	URI                 string `json:"uri" yaml:"uri"`
	Database            string `json:"database" yaml:"database"`
	MigrationCollection string `json:"migration_collection" yaml:"migration_collection"`
	Timeout             int    `json:"timeout" yaml:"timeout"` // 秒
}

// ValidationConfig 验证配置
type ValidationConfig struct {
	Enabled             bool     `json:"enabled" yaml:"enabled"`
	StrictMode          bool     `json:"strict_mode" yaml:"strict_mode"`
	AllowedOperations   []string `json:"allowed_operations" yaml:"allowed_operations"`
	ForbiddenOperations []string `json:"forbidden_operations" yaml:"forbidden_operations"`
	RequireComments     bool     `json:"require_comments" yaml:"require_comments"`
	RequireBackup       bool     `json:"require_backup" yaml:"require_backup"`
}

// Direction 迁移方向
type Direction string

const (
	// DirectionUp 表示向上迁移方向（应用迁移）
	DirectionUp Direction = "up"
	// DirectionDown 表示向下迁移方向（回滚迁移）
	DirectionDown Direction = "down"
)

// BackupConfig 备份配置
type BackupConfig struct {
	Enabled       bool   `json:"enabled" yaml:"enabled"`
	BackupDir     string `json:"backup_dir" yaml:"backup_dir"`
	RetentionDays int    `json:"retention_days" yaml:"retention_days"`
	Compression   bool   `json:"compression" yaml:"compression"`
	AutoBackup    bool   `json:"auto_backup" yaml:"auto_backup"`
}

// Plan 迁移计划
type Plan struct {
	Direction     Direction     `json:"direction"`
	TargetVersion string        `json:"target_version"`
	Steps         int           `json:"steps"`
	MySQLPlan     *MySQLPlan    `json:"mysql_plan"`
	MongoDBPlan   *MongoDBPlan  `json:"mongodb_plan"`
	EstimatedTime time.Duration `json:"estimated_time"`
	CreatedAt     time.Time     `json:"created_at"`
}

// MySQLPlan MySQL迁移计划
type MySQLPlan struct {
	Migrations    []*MySQLMigration `json:"migrations"`
	TotalSteps    int               `json:"total_steps"`
	EstimatedTime time.Duration     `json:"estimated_time"`
}

// MongoDBPlan MongoDB迁移计划
type MongoDBPlan struct {
	Migrations    []*MongoDBMigration `json:"migrations"`
	TotalSteps    int                 `json:"total_steps"`
	EstimatedTime time.Duration       `json:"estimated_time"`
}

// 扫描函数

// scanMySQLMigrations 扫描MySQL迁移文件
func scanMySQLMigrations(dir string) ([]*MySQLMigration, error) {
	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return []*MySQLMigration{}, nil
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	migrationsMap, err := parseMySQLMigrationFiles(files, dir)
	if err != nil {
		return nil, err
	}

	return convertMySQLMigrationsToSlice(migrationsMap), nil
}

// parseMySQLMigrationFiles 解析MySQL迁移文件
func parseMySQLMigrationFiles(files []os.DirEntry, dir string) (map[string]*MySQLMigration, error) {
	migrationsMap := make(map[string]*MySQLMigration)

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		migration, err := parseSingleMySQLMigrationFile(file.Name(), dir)
		if err != nil {
			continue // 跳过无效文件
		}

		if migration != nil {
			if existing, exists := migrationsMap[migration.Version]; exists {
				// 合并up和down文件
				mergeMySQLMigrationFiles(existing, migration)
			} else {
				migrationsMap[migration.Version] = migration
			}
		}
	}

	return migrationsMap, nil
}

// parseSingleMySQLMigrationFile 解析单个MySQL迁移文件
func parseSingleMySQLMigrationFile(filename, dir string) (*MySQLMigration, error) {
	// 解析文件名: {version}_{name}.{up|down}.sql
	parts := strings.Split(filename, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid filename format")
	}

	direction := parts[1]       // up 或 down
	nameWithVersion := parts[0] // {version}_{name}

	// 分离版本号和名称
	underscoreIndex := strings.Index(nameWithVersion, "_")
	if underscoreIndex == -1 {
		return nil, fmt.Errorf("no underscore in filename")
	}

	version := nameWithVersion[:underscoreIndex]
	name := nameWithVersion[underscoreIndex+1:]

	// 验证版本号格式
	if !isValidMigrationVersion(version) {
		return nil, fmt.Errorf("invalid version format")
	}

	migration := &MySQLMigration{
		Version: version,
		Name:    name,
	}

	filePath := filepath.Join(dir, filename)
	switch direction {
	case "up":
		migration.UpFile = filePath
	case "down":
		migration.DownFile = filePath
	default:
		return nil, fmt.Errorf("invalid direction: %s", direction)
	}

	return migration, nil
}

// mergeMySQLMigrationFiles 合并MySQL迁移文件信息
func mergeMySQLMigrationFiles(existing, new *MySQLMigration) {
	if new.UpFile != "" {
		existing.UpFile = new.UpFile
	}
	if new.DownFile != "" {
		existing.DownFile = new.DownFile
	}
}

// convertMySQLMigrationsToSlice 转换MySQL迁移映射为排序切片
func convertMySQLMigrationsToSlice(migrationsMap map[string]*MySQLMigration) []*MySQLMigration {
	migrations := make([]*MySQLMigration, 0, len(migrationsMap))

	for _, migration := range migrationsMap {
		// 只包含同时有up和down文件的迁移
		if migration.UpFile != "" && migration.DownFile != "" {
			migrations = append(migrations, migration)
		}
	}

	// 按版本号排序
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations
}

// scanMongoDBMigrations 扫描MongoDB迁移文件
func scanMongoDBMigrations(dir string) ([]*MongoDBMigration, error) {
	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return []*MongoDBMigration{}, nil
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	migrations := make([]*MongoDBMigration, 0, len(files))

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		if !strings.HasSuffix(filename, ".js") {
			continue
		}

		// 解析文件名: {version}_{name}.js
		nameWithoutExt := strings.TrimSuffix(filename, ".js")
		underscoreIndex := strings.Index(nameWithoutExt, "_")
		if underscoreIndex == -1 {
			continue
		}

		version := nameWithoutExt[:underscoreIndex]
		name := nameWithoutExt[underscoreIndex+1:]

		// 验证版本号格式
		if !isValidMigrationVersion(version) {
			continue
		}

		migration := &MongoDBMigration{
			Version: version,
			Name:    name,
			JSFile:  filepath.Join(dir, filename),
		}

		migrations = append(migrations, migration)
	}

	// 按版本号排序
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// 工具函数

// isValidMigrationVersion 验证迁移版本号格式
func isValidMigrationVersion(version string) bool {
	// 版本号应该是14位数字的时间戳格式: YYYYMMDDHHMMSS
	if len(version) != 14 {
		return false
	}

	for _, char := range version {
		if char < '0' || char > '9' {
			return false
		}
	}

	// 验证时间戳是否有效
	_, err := time.Parse("20060102150405", version)
	return err == nil
}

// GetMigrationInfo 获取迁移信息
func (m *MySQLMigration) GetMigrationInfo() map[string]interface{} {
	return map[string]interface{}{
		"version":   m.Version,
		"name":      m.Name,
		"up_file":   m.UpFile,
		"down_file": m.DownFile,
		"applied":   m.Applied,
		"type":      "mysql",
	}
}

// GetMigrationInfo 获取MongoDB迁移信息
func (m *MongoDBMigration) GetMigrationInfo() map[string]interface{} {
	return map[string]interface{}{
		"version": m.Version,
		"name":    m.Name,
		"js_file": m.JSFile,
		"applied": m.Applied,
		"type":    "mongodb",
	}
}

// String 返回迁移的字符串表示
func (m *MySQLMigration) String() string {
	status := "pending"
	if m.Applied {
		status = "applied"
	}
	return m.Version + "_" + m.Name + " (" + status + ")"
}

// String 返回MongoDB迁移的字符串表示
func (m *MongoDBMigration) String() string {
	status := "pending"
	if m.Applied {
		status = "applied"
	}
	return m.Version + "_" + m.Name + " (" + status + ")"
}

// IsNewer 检查是否比指定版本更新
func (m *MySQLMigration) IsNewer(version string) bool {
	return m.Version > version
}

// IsOlder 检查是否比指定版本更旧
func (m *MySQLMigration) IsOlder(version string) bool {
	return m.Version < version
}

// IsNewer 检查MongoDB迁移是否比指定版本更新
func (m *MongoDBMigration) IsNewer(version string) bool {
	return m.Version > version
}

// IsOlder 检查MongoDB迁移是否比指定版本更旧
func (m *MongoDBMigration) IsOlder(version string) bool {
	return m.Version < version
}

// GetTimestamp 获取迁移的时间戳
func (m *MySQLMigration) GetTimestamp() (time.Time, error) {
	return time.Parse("20060102150405", m.Version)
}

// GetTimestamp 获取MongoDB迁移的时间戳
func (m *MongoDBMigration) GetTimestamp() (time.Time, error) {
	return time.Parse("20060102150405", m.Version)
}

// HasUpFile 检查是否有up文件
func (m *MySQLMigration) HasUpFile() bool {
	return m.UpFile != "" && fileExists(m.UpFile)
}

// HasDownFile 检查是否有down文件
func (m *MySQLMigration) HasDownFile() bool {
	return m.DownFile != "" && fileExists(m.DownFile)
}

// HasJSFile 检查是否有JavaScript文件
func (m *MongoDBMigration) HasJSFile() bool {
	return m.JSFile != "" && fileExists(m.JSFile)
}

// IsComplete 检查MySQL迁移是否完整（有up和down文件）
func (m *MySQLMigration) IsComplete() bool {
	return m.HasUpFile() && m.HasDownFile()
}

// IsComplete 检查MongoDB迁移是否完整（有JS文件）
func (m *MongoDBMigration) IsComplete() bool {
	return m.HasJSFile()
}

// fileExists 检查文件是否存在
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// GetDefaultConfig 获取默认迁移配置
func GetDefaultConfig() *Config {
	return &Config{
		MigrationsDir: "migrations",
		MySQL: &MySQLMigrationConfig{
			Enabled:        true,
			MigrationTable: "schema_migrations",
			Timeout:        300, // 5分钟
		},
		MongoDB: &MongoDBMigrationConfig{
			Enabled:             false,
			MigrationCollection: "migration_records",
			Timeout:             300, // 5分钟
		},
		Validation: &ValidationConfig{
			Enabled:         true,
			StrictMode:      false,
			RequireComments: true,
			RequireBackup:   true,
			AllowedOperations: []string{
				"CREATE TABLE", "ALTER TABLE", "CREATE INDEX", "DROP INDEX",
				"INSERT", "UPDATE",
			},
			ForbiddenOperations: []string{
				"DROP DATABASE", "TRUNCATE TABLE",
			},
		},
		Backup: &BackupConfig{
			Enabled:       true,
			BackupDir:     "backups",
			RetentionDays: 30,
			Compression:   true,
			AutoBackup:    true,
		},
	}
}
