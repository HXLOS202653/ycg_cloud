package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// Validator 迁移验证器
type Validator struct {
	migrationsDir string
	verbose       bool
}

// NewValidator 创建迁移验证器
func NewValidator(migrationsDir string, verbose bool) *Validator {
	return &Validator{
		migrationsDir: migrationsDir,
		verbose:       verbose,
	}
}

// ValidationResult 验证结果
type ValidationResult struct {
	IsValid  bool                `json:"is_valid"`
	Errors   []ValidationError   `json:"errors"`
	Warnings []ValidationWarning `json:"warnings"`
	Summary  ValidationSummary   `json:"summary"`
}

// ValidationError 验证错误
type ValidationError struct {
	Type       string `json:"type"`
	File       string `json:"file"`
	Line       int    `json:"line"`
	Message    string `json:"message"`
	Severity   string `json:"severity"`
	Suggestion string `json:"suggestion"`
}

// ValidationWarning 验证警告
type ValidationWarning struct {
	Type       string `json:"type"`
	File       string `json:"file"`
	Line       int    `json:"line"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
}

// ValidationSummary 验证摘要
type ValidationSummary struct {
	TotalFiles        int `json:"total_files"`
	ValidFiles        int `json:"valid_files"`
	FilesWithErrors   int `json:"files_with_errors"`
	FilesWithWarnings int `json:"files_with_warnings"`
	TotalErrors       int `json:"total_errors"`
	TotalWarnings     int `json:"total_warnings"`
}

// MigrationDirection 迁移方向
type MigrationDirection string

const (
	DirectionUp   MigrationDirection = "up"
	DirectionDown MigrationDirection = "down"
)

// ValidateAllMigrations 验证所有迁移文件
func (v *Validator) ValidateAllMigrations() (*ValidationResult, error) {
	fmt.Println("🔍 开始验证迁移文件...")

	result := &ValidationResult{
		IsValid:  true,
		Errors:   []ValidationError{},
		Warnings: []ValidationWarning{},
		Summary:  ValidationSummary{},
	}

	// 验证MySQL迁移
	if err := v.validateMySQLMigrations(result); err != nil {
		return nil, fmt.Errorf("验证MySQL迁移失败: %v", err)
	}

	// 验证MongoDB迁移
	if err := v.validateMongoDBMigrations(result); err != nil {
		return nil, fmt.Errorf("验证MongoDB迁移失败: %v", err)
	}

	// 计算摘要
	v.calculateSummary(result)

	// 如果有错误，标记为无效
	if len(result.Errors) > 0 {
		result.IsValid = false
	}

	v.printValidationResult(result)

	return result, nil
}

// validateMySQLMigrations 验证MySQL迁移文件
func (v *Validator) validateMySQLMigrations(result *ValidationResult) error {
	mysqlDir := filepath.Join(v.migrationsDir, "mysql")

	// 检查MySQL目录是否存在
	if _, err := os.Stat(mysqlDir); os.IsNotExist(err) {
		v.addWarning(result, "directory", "", 0,
			"MySQL迁移目录不存在",
			"创建migrations/mysql目录并添加迁移文件")
		return nil
	}

	// 获取所有MySQL迁移文件
	migrations, err := v.getAllMySQLMigrations()
	if err != nil {
		return err
	}

	if len(migrations) == 0 {
		v.addWarning(result, "empty_directory", mysqlDir, 0,
			"MySQL迁移目录为空",
			"添加数据库迁移文件")
		return nil
	}

	// 验证每个迁移文件
	for _, mig := range migrations {
		v.validateMySQLMigration(mig, result)
	}

	// 验证迁移版本序列
	v.validateMigrationSequence(migrations, result)

	return nil
}

// validateMongoDBMigrations 验证MongoDB迁移文件
func (v *Validator) validateMongoDBMigrations(result *ValidationResult) error {
	mongoDir := filepath.Join(v.migrationsDir, "mongodb")

	// 检查MongoDB目录是否存在
	if _, err := os.Stat(mongoDir); os.IsNotExist(err) {
		v.addWarning(result, "directory", "", 0,
			"MongoDB迁移目录不存在",
			"如果使用MongoDB，请创建migrations/mongodb目录")
		return nil
	}

	// 获取所有MongoDB迁移文件
	migrations, err := v.getAllMongoDBMigrations()
	if err != nil {
		return err
	}

	// 验证每个MongoDB迁移文件
	for _, mig := range migrations {
		v.validateMongoDBMigration(mig, result)
	}

	return nil
}

// validateMySQLMigration 验证单个MySQL迁移
func (v *Validator) validateMySQLMigration(mig *MySQLMigration, result *ValidationResult) {
	result.Summary.TotalFiles += 2 // up和down文件

	// 验证up文件
	if mig.UpFile != "" {
		v.validateSQLFile(mig.UpFile, DirectionUp, result)
	} else {
		v.addError(result, "missing_file", mig.Version+"_"+mig.Name, 0,
			"缺少up迁移文件", "high",
			"创建对应的.up.sql文件")
	}

	// 验证down文件
	if mig.DownFile != "" {
		v.validateSQLFile(mig.DownFile, DirectionDown, result)
	} else {
		v.addError(result, "missing_file", mig.Version+"_"+mig.Name, 0,
			"缺少down迁移文件", "high",
			"创建对应的.down.sql文件")
	}

	// 验证文件名格式
	v.validateMigrationNameFormat(mig, result)
}

// validateMongoDBMigration 验证单个MongoDB迁移
func (v *Validator) validateMongoDBMigration(mig *MongoDBMigration, result *ValidationResult) {
	result.Summary.TotalFiles++

	if mig.JSFile != "" {
		v.validateJavaScriptFile(mig.JSFile, result)
	} else {
		v.addError(result, "missing_file", mig.Version+"_"+mig.Name, 0,
			"缺少MongoDB迁移文件", "high",
			"创建对应的.js文件")
	}

	// 验证文件名格式
	v.validateMongoMigrationNameFormat(mig, result)
}

// validateSQLFile 验证SQL文件
func (v *Validator) validateSQLFile(filePath string, direction MigrationDirection, result *ValidationResult) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		v.addError(result, "file_read", filePath, 0,
			fmt.Sprintf("无法读取文件: %v", err), "high",
			"检查文件权限和路径")
		return
	}

	sql := string(content)
	lines := strings.Split(sql, "\n")

	// 基本SQL语法验证
	v.validateSQLSyntax(filePath, sql, lines, direction, result)

	// 检查危险操作
	v.validateSQLSafety(filePath, sql, lines, direction, result)

	// 检查编码规范
	v.validateSQLCodingStandards(filePath, sql, lines, result)

	result.Summary.ValidFiles++
}

// validateJavaScriptFile 验证JavaScript文件
func (v *Validator) validateJavaScriptFile(filePath string, result *ValidationResult) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		v.addError(result, "file_read", filePath, 0,
			fmt.Sprintf("无法读取文件: %v", err), "high",
			"检查文件权限和路径")
		return
	}

	js := string(content)
	lines := strings.Split(js, "\n")

	// 验证MongoDB JavaScript语法
	v.validateMongoJSSyntax(filePath, js, lines, result)

	// 检查MongoDB操作安全性
	v.validateMongoJSSafety(filePath, js, lines, result)

	result.Summary.ValidFiles++
}

// validateSQLSyntax 验证SQL语法
func (v *Validator) validateSQLSyntax(filePath, sql string, lines []string, _ MigrationDirection, result *ValidationResult) {
	// 检查是否为空文件
	sqlContent := strings.TrimSpace(regexp.MustCompile(`--.*`).ReplaceAllString(sql, ""))
	if sqlContent == "" {
		v.addWarning(result, "empty_file", filePath, 0,
			"迁移文件为空或只包含注释",
			"添加实际的SQL语句")
		return
	}

	// 检查SQL语句是否以分号结尾
	statements := v.extractSQLStatements(sql)
	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" && !strings.HasSuffix(stmt, ";") {
			v.addWarning(result, "missing_semicolon", filePath, i+1,
				"SQL语句缺少分号结尾",
				"在SQL语句末尾添加分号")
		}
	}

	// 检查关键字大小写
	v.validateSQLKeywordCase(filePath, lines, result)

	// 验证表名和字段名规范
	v.validateSQLNamingConventions(filePath, sql, result)
}

// validateSQLSafety 验证SQL安全性
func (v *Validator) validateSQLSafety(filePath, sql string, _ []string, direction MigrationDirection, result *ValidationResult) {
	// 检查危险操作
	dangerousOperations := map[string]string{
		"DROP DATABASE":      "删除数据库是极其危险的操作",
		"TRUNCATE":           "清空表数据可能导致数据丢失",
		"DELETE FROM.*WHERE": "删除操作应该有WHERE条件限制",
	}

	if direction == DirectionDown {
		// down迁移中的危险操作
		dangerousOperations["DROP TABLE"] = "回滚时删除表可能导致数据丢失"
		dangerousOperations["DROP COLUMN"] = "删除列可能导致数据丢失"
		dangerousOperations["DROP INDEX"] = "删除索引可能影响性能"
	}

	sqlUpper := strings.ToUpper(sql)
	for pattern, message := range dangerousOperations {
		if matched, _ := regexp.MatchString(pattern, sqlUpper); matched {
			severity := "medium"
			if strings.Contains(pattern, "DROP DATABASE") || strings.Contains(pattern, "TRUNCATE") {
				severity = "high"
			}

			v.addError(result, "dangerous_operation", filePath, 0,
				fmt.Sprintf("包含危险操作: %s - %s", pattern, message), severity,
				"仔细审查此操作的必要性，考虑添加备份策略")
		}
	}

	// 检查是否缺少WHERE条件的DELETE语句
	deletePattern := regexp.MustCompile(`(?i)DELETE\s+FROM\s+\w+\s*;`)
	if deletePattern.MatchString(sql) {
		v.addError(result, "unsafe_delete", filePath, 0,
			"DELETE语句缺少WHERE条件，可能删除所有数据", "high",
			"添加适当的WHERE条件限制删除范围")
	}
}

// validateSQLCodingStandards 验证SQL编码规范
func (v *Validator) validateSQLCodingStandards(filePath, sql string, _ []string, result *ValidationResult) {
	// 检查字符集设置
	if strings.Contains(strings.ToUpper(sql), "CREATE TABLE") {
		if !strings.Contains(strings.ToUpper(sql), "CHARSET=UTF8MB4") {
			v.addWarning(result, "charset", filePath, 0,
				"创建表时建议使用utf8mb4字符集",
				"添加 CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
		}
	}

	// 检查注释
	if strings.Contains(strings.ToUpper(sql), "CREATE TABLE") {
		if !strings.Contains(strings.ToUpper(sql), "COMMENT") {
			v.addWarning(result, "missing_comment", filePath, 0,
				"表和字段缺少注释",
				"为表和字段添加适当的COMMENT说明")
		}
	}

	// 检查索引命名规范
	indexPattern := regexp.MustCompile(`(?i)INDEX\s+(\w+)`)
	matches := indexPattern.FindAllStringSubmatch(sql, -1)
	for _, match := range matches {
		if len(match) > 1 {
			indexName := match[1]
			if !strings.HasPrefix(strings.ToLower(indexName), "idx_") {
				v.addWarning(result, "index_naming", filePath, 0,
					fmt.Sprintf("索引名称 '%s' 不符合命名规范", indexName),
					"索引名称建议以'idx_'开头")
			}
		}
	}
}

// validateMongoJSSyntax 验证MongoDB JavaScript语法
func (v *Validator) validateMongoJSSyntax(filePath, js string, _ []string, result *ValidationResult) {
	// 检查是否为空文件
	jsContent := strings.TrimSpace(regexp.MustCompile(`//.*`).ReplaceAllString(js, ""))
	if jsContent == "" {
		v.addWarning(result, "empty_file", filePath, 0,
			"MongoDB迁移文件为空或只包含注释",
			"添加实际的MongoDB操作")
		return
	}

	// 检查常见的MongoDB操作
	hasMongoOperations := false
	mongoOperations := []string{
		"db.", "createCollection", "createIndex", "insertOne", "insertMany",
		"updateOne", "updateMany", "deleteOne", "deleteMany", "drop",
	}

	for _, op := range mongoOperations {
		if strings.Contains(js, op) {
			hasMongoOperations = true
			break
		}
	}

	if !hasMongoOperations {
		v.addWarning(result, "no_mongo_operations", filePath, 0,
			"文件中没有检测到MongoDB操作",
			"确保文件包含有效的MongoDB JavaScript代码")
	}
}

// validateMongoJSSafety 验证MongoDB JavaScript安全性
func (v *Validator) validateMongoJSSafety(filePath, js string, _ []string, result *ValidationResult) {
	// 检查危险的MongoDB操作
	dangerousOperations := map[string]string{
		"db.dropDatabase()": "删除数据库是极其危险的操作",
		"drop()":            "删除集合可能导致数据丢失",
		"deleteMany({})":    "无条件删除所有文档可能导致数据丢失",
	}

	for pattern, message := range dangerousOperations {
		if strings.Contains(js, pattern) {
			v.addError(result, "dangerous_operation", filePath, 0,
				fmt.Sprintf("包含危险操作: %s - %s", pattern, message), "high",
				"仔细审查此操作的必要性，考虑添加备份策略")
		}
	}
}

// validateMigrationSequence 验证迁移版本序列
func (v *Validator) validateMigrationSequence(migrations []*MySQLMigration, result *ValidationResult) {
	if len(migrations) < 2 {
		return
	}

	// 按版本排序
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// 检查版本号是否连续且唯一
	versionMap := make(map[string]bool)
	for i, mig := range migrations {
		// 检查版本号唯一性
		if versionMap[mig.Version] {
			v.addError(result, "duplicate_version", mig.Version, 0,
				fmt.Sprintf("重复的迁移版本号: %s", mig.Version), "high",
				"确保每个迁移版本号都是唯一的")
		}
		versionMap[mig.Version] = true

		// 检查版本号格式
		if !v.isValidVersionFormat(mig.Version) {
			v.addError(result, "invalid_version_format", mig.Version, 0,
				fmt.Sprintf("无效的版本号格式: %s", mig.Version), "medium",
				"版本号应该是14位时间戳格式: YYYYMMDDHHMMSS")
		}

		// 检查时间顺序
		if i > 0 {
			prevTime, err1 := time.Parse("20060102150405", migrations[i-1].Version)
			currTime, err2 := time.Parse("20060102150405", mig.Version)

			if err1 == nil && err2 == nil {
				if currTime.Before(prevTime) {
					v.addWarning(result, "version_order", mig.Version, 0,
						"迁移版本时间顺序异常",
						"确保迁移版本号按时间顺序递增")
				}
			}
		}
	}
}

// validateMigrationNameFormat 验证MySQL迁移文件名格式
func (v *Validator) validateMigrationNameFormat(mig *MySQLMigration, result *ValidationResult) {
	// 验证文件名格式: {version}_{name}.{up|down}.sql
	upPattern := regexp.MustCompile(`^(\d{14})_([a-zA-Z0-9_]+)\.up\.sql$`)
	downPattern := regexp.MustCompile(`^(\d{14})_([a-zA-Z0-9_]+)\.down\.sql$`)

	if mig.UpFile != "" {
		baseName := filepath.Base(mig.UpFile)
		if !upPattern.MatchString(baseName) {
			v.addError(result, "invalid_filename", mig.UpFile, 0,
				"up文件名格式不正确", "medium",
				"文件名应为: {version}_{name}.up.sql")
		}
	}

	if mig.DownFile != "" {
		baseName := filepath.Base(mig.DownFile)
		if !downPattern.MatchString(baseName) {
			v.addError(result, "invalid_filename", mig.DownFile, 0,
				"down文件名格式不正确", "medium",
				"文件名应为: {version}_{name}.down.sql")
		}
	}
}

// validateMongoMigrationNameFormat 验证MongoDB迁移文件名格式
func (v *Validator) validateMongoMigrationNameFormat(mig *MongoDBMigration, result *ValidationResult) {
	// 验证文件名格式: {version}_{name}.js
	pattern := regexp.MustCompile(`^(\d{14})_([a-zA-Z0-9_]+)\.js$`)

	if mig.JSFile != "" {
		baseName := filepath.Base(mig.JSFile)
		if !pattern.MatchString(baseName) {
			v.addError(result, "invalid_filename", mig.JSFile, 0,
				"MongoDB迁移文件名格式不正确", "medium",
				"文件名应为: {version}_{name}.js")
		}
	}
}

// 辅助方法

// addError 添加验证错误
func (v *Validator) addError(result *ValidationResult, errType, file string, line int, message, severity, suggestion string) {
	result.Errors = append(result.Errors, ValidationError{
		Type:       errType,
		File:       file,
		Line:       line,
		Message:    message,
		Severity:   severity,
		Suggestion: suggestion,
	})
}

// addWarning 添加验证警告
func (v *Validator) addWarning(result *ValidationResult, warnType, file string, line int, message, suggestion string) {
	result.Warnings = append(result.Warnings, ValidationWarning{
		Type:       warnType,
		File:       file,
		Line:       line,
		Message:    message,
		Suggestion: suggestion,
	})
}

// calculateSummary 计算验证摘要
func (v *Validator) calculateSummary(result *ValidationResult) {
	result.Summary.TotalErrors = len(result.Errors)
	result.Summary.TotalWarnings = len(result.Warnings)

	// 统计有错误和警告的文件
	filesWithErrors := make(map[string]bool)
	filesWithWarnings := make(map[string]bool)

	for _, err := range result.Errors {
		filesWithErrors[err.File] = true
	}

	for _, warn := range result.Warnings {
		filesWithWarnings[warn.File] = true
	}

	result.Summary.FilesWithErrors = len(filesWithErrors)
	result.Summary.FilesWithWarnings = len(filesWithWarnings)
}

// printValidationResult 打印验证结果
func (v *Validator) printValidationResult(result *ValidationResult) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 迁移文件验证结果")
	fmt.Println(strings.Repeat("=", 80))

	// 打印摘要
	if result.IsValid {
		fmt.Println("✅ 验证通过")
	} else {
		fmt.Println("❌ 验证失败")
	}

	fmt.Printf("📁 文件统计: 总计 %d 个，有效 %d 个\n",
		result.Summary.TotalFiles, result.Summary.ValidFiles)
	fmt.Printf("🚨 错误: %d 个 (涉及 %d 个文件)\n",
		result.Summary.TotalErrors, result.Summary.FilesWithErrors)
	fmt.Printf("⚠️  警告: %d 个 (涉及 %d 个文件)\n",
		result.Summary.TotalWarnings, result.Summary.FilesWithWarnings)

	// 打印错误详情
	if len(result.Errors) > 0 {
		fmt.Println("\n🚨 错误详情:")
		fmt.Println(strings.Repeat("-", 80))
		for i, err := range result.Errors {
			fmt.Printf("%d. [%s] %s:%d\n", i+1, err.Severity, err.File, err.Line)
			fmt.Printf("   错误: %s\n", err.Message)
			if err.Suggestion != "" {
				fmt.Printf("   建议: %s\n", err.Suggestion)
			}
			fmt.Println()
		}
	}

	// 打印警告详情（仅在verbose模式下）
	if v.verbose && len(result.Warnings) > 0 {
		fmt.Println("\n⚠️  警告详情:")
		fmt.Println(strings.Repeat("-", 80))
		for i, warn := range result.Warnings {
			fmt.Printf("%d. %s:%d\n", i+1, warn.File, warn.Line)
			fmt.Printf("   警告: %s\n", warn.Message)
			if warn.Suggestion != "" {
				fmt.Printf("   建议: %s\n", warn.Suggestion)
			}
			fmt.Println()
		}
	}

	fmt.Println(strings.Repeat("=", 80))
}

// 工具方法

// getAllMySQLMigrations 获取所有MySQL迁移文件
func (v *Validator) getAllMySQLMigrations() ([]*MySQLMigration, error) {
	mysqlDir := filepath.Join(v.migrationsDir, "mysql")
	return scanMySQLMigrations(mysqlDir)
}

// getAllMongoDBMigrations 获取所有MongoDB迁移文件
func (v *Validator) getAllMongoDBMigrations() ([]*MongoDBMigration, error) {
	mongoDir := filepath.Join(v.migrationsDir, "mongodb")
	return scanMongoDBMigrations(mongoDir)
}

// extractSQLStatements 提取SQL语句
func (v *Validator) extractSQLStatements(sql string) []string {
	// 简单的SQL语句分割，实际项目中可能需要更复杂的解析
	statements := strings.Split(sql, ";")
	var result []string

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" && !strings.HasPrefix(stmt, "--") {
			result = append(result, stmt)
		}
	}

	return result
}

// validateSQLKeywordCase 验证SQL关键字大小写
func (v *Validator) validateSQLKeywordCase(filePath string, lines []string, result *ValidationResult) {
	keywords := []string{
		"CREATE", "ALTER", "DROP", "INSERT", "UPDATE", "DELETE", "SELECT",
		"TABLE", "INDEX", "PRIMARY", "KEY", "FOREIGN", "REFERENCES",
		"NOT", "NULL", "DEFAULT", "AUTO_INCREMENT", "ENGINE", "CHARSET",
	}

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}

		for _, keyword := range keywords {
			// 检查是否有小写的关键字
			pattern := regexp.MustCompile(`\b` + strings.ToLower(keyword) + `\b`)
			if pattern.MatchString(strings.ToLower(line)) {
				// 确认原文中确实是小写
				if strings.Contains(line, strings.ToLower(keyword)) {
					v.addWarning(result, "keyword_case", filePath, lineNum+1,
						fmt.Sprintf("SQL关键字建议使用大写: %s", keyword),
						"将SQL关键字改为大写以提高可读性")
				}
			}
		}
	}
}

// validateSQLNamingConventions 验证SQL命名规范
func (v *Validator) validateSQLNamingConventions(filePath, sql string, result *ValidationResult) {
	// 检查表名命名规范
	tablePattern := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(\w+)`)
	matches := tablePattern.FindAllStringSubmatch(sql, -1)
	for _, match := range matches {
		if len(match) > 1 {
			tableName := match[1]
			if !v.isValidTableName(tableName) {
				v.addWarning(result, "table_naming", filePath, 0,
					fmt.Sprintf("表名 '%s' 不符合命名规范", tableName),
					"表名应使用小写字母和下划线，如: user_sessions")
			}
		}
	}
}

// isValidVersionFormat 验证版本号格式
func (v *Validator) isValidVersionFormat(version string) bool {
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

// isValidTableName 验证表名是否有效
func (v *Validator) isValidTableName(name string) bool {
	// 表名应该是小写字母、数字和下划线的组合
	pattern := regexp.MustCompile(`^[a-z][a-z0-9_]*$`)
	return pattern.MatchString(name)
}
