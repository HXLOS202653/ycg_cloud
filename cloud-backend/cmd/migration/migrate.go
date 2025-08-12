package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// MigrationRecord 迁移记录结构
type MigrationRecord struct {
	ID        uint      `gorm:"primaryKey"`
	Version   string    `gorm:"uniqueIndex;size:20;not null" json:"version"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	UpSQL     string    `gorm:"type:text" json:"up_sql"`
	DownSQL   string    `gorm:"type:text" json:"down_sql"`
	AppliedAt time.Time `gorm:"not null" json:"applied_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (MigrationRecord) TableName() string {
	return "schema_migrations"
}

// MigrationFile 迁移文件结构
type MigrationFile struct {
	Version  string
	Name     string
	UpFile   string
	DownFile string
	Applied  bool
}

// Up 执行向上迁移
func (mt *MigrationTool) Up(targetVersion string, steps int) error {
	// 初始化迁移表
	if err := mt.initMigrationTable(); err != nil {
		return fmt.Errorf("初始化迁移表失败: %w", err)
	}

	// 获取待执行的迁移文件
	migrations, err := mt.getPendingMigrations()
	if err != nil {
		return fmt.Errorf("获取待迁移文件失败: %w", err)
	}

	if len(migrations) == 0 {
		fmt.Println("📋 没有待执行的迁移")
		return nil
	}

	// 过滤需要执行的迁移
	toApply := mt.filterMigrationsForUp(migrations, targetVersion, steps)

	if len(toApply) == 0 {
		fmt.Println("📋 没有符合条件的迁移需要执行")
		return nil
	}

	fmt.Printf("🚀 准备执行 %d 个迁移:\n", len(toApply))
	for _, migration := range toApply {
		fmt.Printf("  📄 %s_%s\n", migration.Version, migration.Name)
	}

	if mt.dryRun {
		fmt.Println("🔍 预览模式，不会实际执行迁移")
		return nil
	}

	// 执行迁移
	for i, migration := range toApply {
		fmt.Printf("\n⏳ [%d/%d] 执行迁移: %s_%s\n", i+1, len(toApply), migration.Version, migration.Name)

		if err := mt.applyMigration(migration); err != nil {
			return fmt.Errorf("执行迁移 %s_%s 失败: %w", migration.Version, migration.Name, err)
		}

		fmt.Printf("✅ 迁移 %s_%s 执行成功\n", migration.Version, migration.Name)
	}

	fmt.Printf("\n🎉 成功执行了 %d 个迁移\n", len(toApply))
	return nil
}

// Down 执行向下迁移（回滚）
func (mt *MigrationTool) Down(targetVersion string, steps int) error {
	// 初始化迁移表
	if err := mt.initMigrationTable(); err != nil {
		return fmt.Errorf("初始化迁移表失败: %w", err)
	}

	// 获取已应用的迁移
	appliedMigrations, err := mt.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("获取已应用迁移失败: %w", err)
	}

	if len(appliedMigrations) == 0 {
		fmt.Println("📋 没有已应用的迁移可以回滚")
		return nil
	}

	// 过滤需要回滚的迁移
	toRollback := mt.filterMigrationsForDown(appliedMigrations, targetVersion, steps)

	if len(toRollback) == 0 {
		fmt.Println("📋 没有符合条件的迁移需要回滚")
		return nil
	}

	fmt.Printf("⚠️  准备回滚 %d 个迁移:\n", len(toRollback))
	for _, migration := range toRollback {
		fmt.Printf("  📄 %s_%s\n", migration.Version, migration.Name)
	}

	if mt.dryRun {
		fmt.Println("🔍 预览模式，不会实际执行回滚")
		return nil
	}

	// 执行回滚
	for i, migration := range toRollback {
		fmt.Printf("\n⏳ [%d/%d] 回滚迁移: %s_%s\n", i+1, len(toRollback), migration.Version, migration.Name)

		if err := mt.rollbackMigration(migration); err != nil {
			return fmt.Errorf("回滚迁移 %s_%s 失败: %w", migration.Version, migration.Name, err)
		}

		fmt.Printf("✅ 迁移 %s_%s 回滚成功\n", migration.Version, migration.Name)
	}

	fmt.Printf("\n🎉 成功回滚了 %d 个迁移\n", len(toRollback))
	return nil
}

// Status 显示迁移状态
func (mt *MigrationTool) Status() error {
	// 初始化迁移表
	if err := mt.initMigrationTable(); err != nil {
		return fmt.Errorf("初始化迁移表失败: %w", err)
	}

	// 获取所有迁移文件
	allMigrations, err := mt.getAllMigrations()
	if err != nil {
		return fmt.Errorf("获取迁移文件失败: %w", err)
	}

	// 获取已应用的迁移记录
	appliedMap, err := mt.getAppliedMigrationsMap()
	if err != nil {
		return fmt.Errorf("获取已应用迁移记录失败: %w", err)
	}

	fmt.Println("📊 数据库迁移状态:")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%-8s %-10s %-30s %-20s\n", "版本", "状态", "名称", "应用时间")
	fmt.Println(strings.Repeat("-", 80))

	if len(allMigrations) == 0 {
		fmt.Println("没有找到迁移文件")
		return nil
	}

	appliedCount := 0
	for _, migration := range allMigrations {
		status := "待执行"
		appliedTime := ""

		if record, exists := appliedMap[migration.Version]; exists {
			status = "已应用"
			appliedTime = record.AppliedAt.Format("2006-01-02 15:04:05")
			appliedCount++
		}

		statusIcon := "⏳"
		if status == "已应用" {
			statusIcon = "✅"
		}

		fmt.Printf("%-8s %-10s %-30s %-20s\n",
			migration.Version,
			statusIcon+" "+status,
			migration.Name,
			appliedTime)
	}

	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("总计: %d 个迁移文件, %d 个已应用, %d 个待执行\n",
		len(allMigrations), appliedCount, len(allMigrations)-appliedCount)

	return nil
}

// Create 创建新的迁移文件
func (mt *MigrationTool) Create(name string) error {
	if err := validateMigrationName(name); err != nil {
		return err
	}

	// 生成版本号（基于时间戳）
	version := time.Now().Format("20060102150405")

	// 确保MySQL迁移目录存在
	mysqlDir := getMigrationsPath(mt.migrationsDir, "mysql")
	if err := ensureDir(mysqlDir); err != nil {
		return fmt.Errorf("创建MySQL迁移目录失败: %w", err)
	}

	// 生成文件名
	upFile := filepath.Join(mysqlDir, fmt.Sprintf("%s_%s.up.sql", version, name))
	downFile := filepath.Join(mysqlDir, fmt.Sprintf("%s_%s.down.sql", version, name))

	// 检查文件是否已存在
	if _, err := os.Stat(upFile); !os.IsNotExist(err) {
		return fmt.Errorf("迁移文件已存在: %s", upFile)
	}

	// 创建up文件
	upTemplate := fmt.Sprintf(`-- +migrate Up
-- 创建迁移: %s
-- 版本: %s
-- 创建时间: %s

-- 在此处添加你的向上迁移SQL语句
-- 例如:
-- CREATE TABLE example (
--     id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

`, name, version, time.Now().Format("2006-01-02 15:04:05"))

	if err := os.WriteFile(upFile, []byte(upTemplate), 0o600); err != nil {
		return fmt.Errorf("创建up文件失败: %w", err)
	}

	// 创建down文件
	downTemplate := fmt.Sprintf(`-- +migrate Down
-- 回滚迁移: %s
-- 版本: %s
-- 创建时间: %s

-- 在此处添加你的回滚SQL语句
-- 例如:
-- DROP TABLE IF EXISTS example;

`, name, version, time.Now().Format("2006-01-02 15:04:05"))

	if err := os.WriteFile(downFile, []byte(downTemplate), 0o600); err != nil {
		// 如果down文件创建失败，删除已创建的up文件
		_ = os.Remove(upFile)
		return fmt.Errorf("创建down文件失败: %w", err)
	}

	fmt.Printf("✅ 成功创建迁移文件:\n")
	fmt.Printf("  📄 Up:   %s\n", upFile)
	fmt.Printf("  📄 Down: %s\n", downFile)
	fmt.Printf("  🏷️  版本: %s\n", version)

	return nil
}

// 内部辅助方法

// initMigrationTable 初始化迁移记录表
func (mt *MigrationTool) initMigrationTable() error {
	if mt.dryRun {
		fmt.Println("🔍 预览模式: 跳过迁移表初始化")
		return nil
	}

	return mt.db.AutoMigrate(&MigrationRecord{})
}

// getPendingMigrations 获取待执行的迁移
func (mt *MigrationTool) getPendingMigrations() ([]*MigrationFile, error) {
	allMigrations, err := mt.getAllMigrations()
	if err != nil {
		return nil, err
	}

	appliedMap, err := mt.getAppliedMigrationsMap()
	if err != nil {
		return nil, err
	}

	var pending []*MigrationFile
	for _, migration := range allMigrations {
		if _, exists := appliedMap[migration.Version]; !exists {
			pending = append(pending, migration)
		}
	}

	return pending, nil
}

// getAppliedMigrations 获取已应用的迁移（按版本倒序）
func (mt *MigrationTool) getAppliedMigrations() ([]*MigrationRecord, error) {
	var records []*MigrationRecord
	err := mt.db.Order("version DESC").Find(&records).Error
	return records, err
}

// getAppliedMigrationsMap 获取已应用迁移的映射
func (mt *MigrationTool) getAppliedMigrationsMap() (map[string]*MigrationRecord, error) {
	records, err := mt.getAppliedMigrations()
	if err != nil {
		return nil, err
	}

	appliedMap := make(map[string]*MigrationRecord)
	for _, record := range records {
		appliedMap[record.Version] = record
	}

	return appliedMap, nil
}

// getAllMigrations 获取所有迁移文件
func (mt *MigrationTool) getAllMigrations() ([]*MigrationFile, error) {
	mysqlDir := getMigrationsPath(mt.migrationsDir, "mysql")

	// 检查目录是否存在
	if _, err := os.Stat(mysqlDir); os.IsNotExist(err) {
		return []*MigrationFile{}, nil
	}

	files, err := os.ReadDir(mysqlDir)
	if err != nil {
		return nil, fmt.Errorf("读取迁移目录失败: %w", err)
	}

	migrationsMap, err := mt.parseMigrationFiles(files, mysqlDir)
	if err != nil {
		return nil, err
	}

	migrations := mt.convertToSortedSlice(migrationsMap)
	return migrations, nil
}

// parseMigrationFiles 解析迁移文件并构建映射
func (mt *MigrationTool) parseMigrationFiles(files []os.DirEntry, mysqlDir string) (map[string]*MigrationFile, error) {
	migrationsMap := make(map[string]*MigrationFile)
	var invalidFiles []string

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		if !strings.HasSuffix(filename, ".sql") {
			continue
		}

		migrationInfo, err := mt.parseMigrationFilename(filename)
		if err != nil {
			invalidFiles = append(invalidFiles, filename)
			continue // 跳过无效文件，但记录
		}

		// 检查是否有重复的迁移版本和名称但方向不同的情况
		if existing, exists := migrationsMap[migrationInfo.version]; exists {
			if existing.Name != migrationInfo.name {
				return nil, fmt.Errorf("版本 %s 存在不同的迁移名称: '%s' 和 '%s'",
					migrationInfo.version, existing.Name, migrationInfo.name)
			}
		}

		migration := mt.getOrCreateMigration(migrationsMap, migrationInfo.version, migrationInfo.name)

		// 检查是否重复设置相同方向的文件
		filePath := filepath.Join(mysqlDir, filename)
		if err := mt.validateAndSetMigrationFile(migration, migrationInfo.direction, filePath); err != nil {
			return nil, fmt.Errorf("设置迁移文件失败 %s: %w", filename, err)
		}
	}

	// 如果有太多无效文件，可能表明目录配置错误
	if len(invalidFiles) > len(files)/2 {
		return nil, fmt.Errorf("发现过多无效的迁移文件 (%d/%d): %v",
			len(invalidFiles), len(files), invalidFiles)
	}

	// 记录无效文件但不中断
	if len(invalidFiles) > 0 && mt.verbose {
		fmt.Printf("Warning: 跳过 %d 个无效的迁移文件: %v\n", len(invalidFiles), invalidFiles)
	}

	return migrationsMap, nil
}

// migrationFileInfo 迁移文件信息
type migrationFileInfo struct {
	version   string
	name      string
	direction string
}

// parseMigrationFilename 解析迁移文件名
func (mt *MigrationTool) parseMigrationFilename(filename string) (*migrationFileInfo, error) {
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
		return nil, fmt.Errorf("no underscore found in filename")
	}

	version := nameWithVersion[:underscoreIndex]
	name := nameWithVersion[underscoreIndex+1:]

	// 验证版本号格式
	if !isValidVersion(version) {
		return nil, fmt.Errorf("invalid version format")
	}

	return &migrationFileInfo{
		version:   version,
		name:      name,
		direction: direction,
	}, nil
}

// getOrCreateMigration 获取或创建迁移对象
func (mt *MigrationTool) getOrCreateMigration(migrationsMap map[string]*MigrationFile, version, name string) *MigrationFile {
	migration, exists := migrationsMap[version]
	if !exists {
		migration = &MigrationFile{
			Version: version,
			Name:    name,
		}
		migrationsMap[version] = migration
	}
	return migration
}

// validateAndSetMigrationFile 验证并设置迁移文件路径
func (mt *MigrationTool) validateAndSetMigrationFile(migration *MigrationFile, direction, filePath string) error {
	switch direction {
	case "up":
		if migration.UpFile != "" {
			return fmt.Errorf("重复的up迁移文件: 已存在 %s，尝试设置 %s", migration.UpFile, filePath)
		}
		migration.UpFile = filePath
	case "down":
		if migration.DownFile != "" {
			return fmt.Errorf("重复的down迁移文件: 已存在 %s，尝试设置 %s", migration.DownFile, filePath)
		}
		migration.DownFile = filePath
	default:
		return fmt.Errorf("无效的迁移方向: %s", direction)
	}
	return nil
}

// convertToSortedSlice 转换为排序的切片
func (mt *MigrationTool) convertToSortedSlice(migrationsMap map[string]*MigrationFile) []*MigrationFile {
	var migrations []*MigrationFile
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

// filterMigrationsForUp 过滤需要向上迁移的文件
func (mt *MigrationTool) filterMigrationsForUp(migrations []*MigrationFile, targetVersion string, steps int) []*MigrationFile {
	// 预分配切片容量以提高性能
	capacity := len(migrations)
	if steps > 0 && steps < capacity {
		capacity = steps
	}
	filtered := make([]*MigrationFile, 0, capacity)

	for _, migration := range migrations {
		// 如果指定了目标版本，只迁移到该版本
		if targetVersion != "" && migration.Version > targetVersion {
			break
		}

		filtered = append(filtered, migration)

		// 如果指定了步数，限制迁移数量
		if steps > 0 && len(filtered) >= steps {
			break
		}
	}

	return filtered
}

// filterMigrationsForDown 过滤需要向下迁移的记录
func (mt *MigrationTool) filterMigrationsForDown(records []*MigrationRecord, targetVersion string, steps int) []*MigrationRecord {
	filtered := make([]*MigrationRecord, 0, len(records))

	for _, record := range records {
		// 如果指定了目标版本，只回滚到该版本之后的迁移
		if targetVersion != "" && record.Version <= targetVersion {
			break
		}

		filtered = append(filtered, record)

		// 如果指定了步数，限制回滚数量
		if steps > 0 && len(filtered) >= steps {
			break
		}
	}

	return filtered
}

// applyMigration 应用迁移
func (mt *MigrationTool) applyMigration(migration *MigrationFile) error {
	// 读取up文件内容
	upSQL, err := os.ReadFile(migration.UpFile)
	if err != nil {
		return fmt.Errorf("读取up文件失败: %w", err)
	}

	// 读取down文件内容
	downSQL, err := os.ReadFile(migration.DownFile)
	if err != nil {
		return fmt.Errorf("读取down文件失败: %w", err)
	}

	// 执行迁移SQL
	if err := mt.executeSQLStatements(string(upSQL)); err != nil {
		return fmt.Errorf("执行迁移SQL失败: %w", err)
	}

	// 记录迁移
	record := &MigrationRecord{
		Version:   migration.Version,
		Name:      migration.Name,
		UpSQL:     string(upSQL),
		DownSQL:   string(downSQL),
		AppliedAt: time.Now(),
	}

	if err := mt.db.Create(record).Error; err != nil {
		return fmt.Errorf("记录迁移失败: %w", err)
	}

	return nil
}

// rollbackMigration 回滚迁移
func (mt *MigrationTool) rollbackMigration(record *MigrationRecord) error {
	// 执行回滚SQL
	if err := mt.executeSQLStatements(record.DownSQL); err != nil {
		return fmt.Errorf("执行回滚SQL失败: %w", err)
	}

	// 删除迁移记录
	if err := mt.db.Delete(record).Error; err != nil {
		return fmt.Errorf("删除迁移记录失败: %w", err)
	}

	return nil
}

// executeSQLStatements 执行SQL语句
func (mt *MigrationTool) executeSQLStatements(sql string) error {
	// 分割SQL语句（简单实现，实际项目中可能需要更复杂的解析）
	statements := strings.Split(sql, ";")

	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" || strings.HasPrefix(statement, "--") {
			continue
		}

		if mt.verbose {
			fmt.Printf("🔧 执行SQL: %s\n", statement)
		}

		if err := mt.db.Exec(statement).Error; err != nil {
			return fmt.Errorf("执行SQL语句失败 '%s': %w", statement, err)
		}
	}

	return nil
}

// isValidVersion 验证版本号格式
func isValidVersion(version string) bool {
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
