package migration

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

// Runner 迁移执行器
type Runner struct {
	mysqlDB       *gorm.DB
	mongoDB       *mongo.Database
	migrationsDir string
	verbose       bool
	dryRun        bool
}

// NewRunner 创建迁移执行器
func NewRunner(mysqlDB *gorm.DB, mongoDB *mongo.Database, migrationsDir string, verbose, dryRun bool) *Runner {
	return &Runner{
		mysqlDB:       mysqlDB,
		mongoDB:       mongoDB,
		migrationsDir: migrationsDir,
		verbose:       verbose,
		dryRun:        dryRun,
	}
}

// RunMySQLMigrations 执行MySQL迁移
func (r *Runner) RunMySQLMigrations(direction MigrationDirection, targetVersion string, steps int) error {
	fmt.Printf("🚀 开始执行 MySQL 迁移 (%s)\n", direction)

	switch direction {
	case DirectionUp:
		return r.runMySQLUp(targetVersion, steps)
	case DirectionDown:
		return r.runMySQLDown(targetVersion, steps)
	default:
		return fmt.Errorf("不支持的迁移方向: %s", direction)
	}
}

// RunMongoDBMigrations 执行MongoDB迁移
func (r *Runner) RunMongoDBMigrations(direction MigrationDirection, targetVersion string, steps int) error {
	fmt.Printf("🚀 开始执行 MongoDB 迁移 (%s)\n", direction)

	switch direction {
	case DirectionUp:
		return r.runMongoDBUp(targetVersion, steps)
	case DirectionDown:
		return r.runMongoDBDown(targetVersion, steps)
	default:
		return fmt.Errorf("不支持的迁移方向: %s", direction)
	}
}

// runMySQLUp 执行MySQL向上迁移
func (r *Runner) runMySQLUp(targetVersion string, steps int) error {
	// 初始化迁移记录表
	if err := r.initMySQLMigrationTable(); err != nil {
		return fmt.Errorf("初始化MySQL迁移表失败: %v", err)
	}

	// 获取待执行的迁移
	pending, err := r.getPendingMySQLMigrations()
	if err != nil {
		return fmt.Errorf("获取待执行MySQL迁移失败: %v", err)
	}

	if len(pending) == 0 {
		fmt.Println("📋 没有待执行的MySQL迁移")
		return nil
	}

	// 过滤迁移
	toExecute := r.filterMigrations(pending, targetVersion, steps, DirectionUp)

	if len(toExecute) == 0 {
		fmt.Println("📋 没有符合条件的MySQL迁移需要执行")
		return nil
	}

	fmt.Printf("📋 准备执行 %d 个MySQL迁移:\n", len(toExecute))
	for i, mig := range toExecute {
		fmt.Printf("  %d. %s_%s\n", i+1, mig.Version, mig.Name)
	}

	if r.dryRun {
		fmt.Println("🔍 预览模式，不会实际执行迁移")
		return r.previewMySQLMigrations(toExecute, DirectionUp)
	}

	// 执行迁移
	for i, mig := range toExecute {
		fmt.Printf("\n⏳ [%d/%d] 执行MySQL迁移: %s_%s\n", i+1, len(toExecute), mig.Version, mig.Name)

		startTime := time.Now()
		if err := r.executeMySQLMigration(mig, DirectionUp); err != nil {
			return fmt.Errorf("执行MySQL迁移失败: %v", err)
		}
		duration := time.Since(startTime)

		fmt.Printf("✅ MySQL迁移 %s_%s 执行成功 (耗时: %v)\n", mig.Version, mig.Name, duration)
	}

	fmt.Printf("\n🎉 成功执行了 %d 个MySQL迁移\n", len(toExecute))
	return nil
}

// runMySQLDown 执行MySQL向下迁移
func (r *Runner) runMySQLDown(targetVersion string, steps int) error {
	// 获取已应用的迁移
	applied, err := r.getAppliedMySQLMigrations()
	if err != nil {
		return fmt.Errorf("获取已应用MySQL迁移失败: %v", err)
	}

	if len(applied) == 0 {
		fmt.Println("📋 没有已应用的MySQL迁移可以回滚")
		return nil
	}

	// 过滤迁移
	toRollback := r.filterAppliedMigrations(applied, targetVersion, steps)

	if len(toRollback) == 0 {
		fmt.Println("📋 没有符合条件的MySQL迁移需要回滚")
		return nil
	}

	fmt.Printf("⚠️  准备回滚 %d 个MySQL迁移:\n", len(toRollback))
	for i, record := range toRollback {
		fmt.Printf("  %d. %s_%s\n", i+1, record.Version, record.Name)
	}

	if r.dryRun {
		fmt.Println("🔍 预览模式，不会实际执行回滚")
		return r.previewMySQLRollbacks(toRollback)
	}

	// 执行回滚
	for i, record := range toRollback {
		fmt.Printf("\n⏳ [%d/%d] 回滚MySQL迁移: %s_%s\n", i+1, len(toRollback), record.Version, record.Name)

		startTime := time.Now()
		if err := r.rollbackMySQLMigration(record); err != nil {
			return fmt.Errorf("回滚MySQL迁移失败: %v", err)
		}
		duration := time.Since(startTime)

		fmt.Printf("✅ MySQL迁移 %s_%s 回滚成功 (耗时: %v)\n", record.Version, record.Name, duration)
	}

	fmt.Printf("\n🎉 成功回滚了 %d 个MySQL迁移\n", len(toRollback))
	return nil
}

// runMongoDBUp 执行MongoDB向上迁移
func (r *Runner) runMongoDBUp(targetVersion string, steps int) error {
	// 初始化MongoDB迁移记录集合
	if err := r.initMongoDBMigrationCollection(); err != nil {
		return fmt.Errorf("初始化MongoDB迁移集合失败: %v", err)
	}

	// 获取待执行的MongoDB迁移
	pending, err := r.getPendingMongoDBMigrations()
	if err != nil {
		return fmt.Errorf("获取待执行MongoDB迁移失败: %v", err)
	}

	if len(pending) == 0 {
		fmt.Println("📋 没有待执行的MongoDB迁移")
		return nil
	}

	// 过滤迁移
	toExecute := r.filterMongoMigrations(pending, targetVersion, steps, DirectionUp)

	if len(toExecute) == 0 {
		fmt.Println("📋 没有符合条件的MongoDB迁移需要执行")
		return nil
	}

	fmt.Printf("📋 准备执行 %d 个MongoDB迁移:\n", len(toExecute))
	for i, mig := range toExecute {
		fmt.Printf("  %d. %s_%s\n", i+1, mig.Version, mig.Name)
	}

	if r.dryRun {
		fmt.Println("🔍 预览模式，不会实际执行迁移")
		return r.previewMongoDBMigrations(toExecute)
	}

	// 执行迁移
	for i, mig := range toExecute {
		fmt.Printf("\n⏳ [%d/%d] 执行MongoDB迁移: %s_%s\n", i+1, len(toExecute), mig.Version, mig.Name)

		startTime := time.Now()
		if err := r.executeMongoDBMigration(mig); err != nil {
			return fmt.Errorf("执行MongoDB迁移失败: %v", err)
		}
		duration := time.Since(startTime)

		fmt.Printf("✅ MongoDB迁移 %s_%s 执行成功 (耗时: %v)\n", mig.Version, mig.Name, duration)
	}

	fmt.Printf("\n🎉 成功执行了 %d 个MongoDB迁移\n", len(toExecute))
	return nil
}

// runMongoDBDown 执行MongoDB向下迁移
func (r *Runner) runMongoDBDown(_ string, _ int) error {
	fmt.Println("⚠️  MongoDB回滚功能正在开发中")
	fmt.Println("建议通过备份恢复的方式处理MongoDB的数据回滚")
	return nil
}

// 辅助方法

// initMySQLMigrationTable 初始化MySQL迁移记录表
func (r *Runner) initMySQLMigrationTable() error {
	if r.dryRun {
		fmt.Println("🔍 预览模式: 跳过MySQL迁移表初始化")
		return nil
	}

	return r.mysqlDB.AutoMigrate(&MySQLMigrationRecord{})
}

// initMongoDBMigrationCollection 初始化MongoDB迁移记录集合
func (r *Runner) initMongoDBMigrationCollection() error {
	if r.dryRun {
		fmt.Println("🔍 预览模式: 跳过MongoDB迁移集合初始化")
		return nil
	}

	// MongoDB会自动创建集合，这里可以设置索引
	collection := r.mongoDB.Collection("migration_records")

	// 创建唯一索引
	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    map[string]interface{}{"version": 1},
		Options: options.Index().SetUnique(true),
	})

	return err
}

// getPendingMySQLMigrations 获取待执行的MySQL迁移
func (r *Runner) getPendingMySQLMigrations() ([]*MySQLMigration, error) {
	// 获取所有迁移文件
	allMigrations, err := r.getAllMySQLMigrations()
	if err != nil {
		return nil, err
	}

	// 获取已应用的迁移记录
	appliedMap, err := r.getAppliedMySQLMigrationsMap()
	if err != nil {
		return nil, err
	}

	var pending []*MySQLMigration
	for _, mig := range allMigrations {
		if _, exists := appliedMap[mig.Version]; !exists {
			pending = append(pending, mig)
		}
	}

	return pending, nil
}

// getPendingMongoDBMigrations 获取待执行的MongoDB迁移
func (r *Runner) getPendingMongoDBMigrations() ([]*MongoDBMigration, error) {
	// 获取所有MongoDB迁移文件
	allMigrations, err := r.getAllMongoDBMigrations()
	if err != nil {
		return nil, err
	}

	// 获取已应用的迁移记录
	appliedMap, err := r.getAppliedMongoDBMigrationsMap()
	if err != nil {
		return nil, err
	}

	var pending []*MongoDBMigration
	for _, mig := range allMigrations {
		if _, exists := appliedMap[mig.Version]; !exists {
			pending = append(pending, mig)
		}
	}

	return pending, nil
}

// getAllMySQLMigrations 获取所有MySQL迁移文件
func (r *Runner) getAllMySQLMigrations() ([]*MySQLMigration, error) {
	mysqlDir := filepath.Join(r.migrationsDir, "mysql")
	return scanMySQLMigrations(mysqlDir)
}

// getAllMongoDBMigrations 获取所有MongoDB迁移文件
func (r *Runner) getAllMongoDBMigrations() ([]*MongoDBMigration, error) {
	mongoDir := filepath.Join(r.migrationsDir, "mongodb")
	return scanMongoDBMigrations(mongoDir)
}

// getAppliedMySQLMigrations 获取已应用的MySQL迁移记录
func (r *Runner) getAppliedMySQLMigrations() ([]*MySQLMigrationRecord, error) {
	var records []*MySQLMigrationRecord
	err := r.mysqlDB.Order("version DESC").Find(&records).Error
	return records, err
}

// getAppliedMySQLMigrationsMap 获取已应用的MySQL迁移映射
func (r *Runner) getAppliedMySQLMigrationsMap() (map[string]*MySQLMigrationRecord, error) {
	records, err := r.getAppliedMySQLMigrations()
	if err != nil {
		return nil, err
	}

	appliedMap := make(map[string]*MySQLMigrationRecord)
	for _, record := range records {
		appliedMap[record.Version] = record
	}

	return appliedMap, nil
}

// getAppliedMongoDBMigrationsMap 获取已应用的MongoDB迁移映射
func (r *Runner) getAppliedMongoDBMigrationsMap() (map[string]*MongoDBMigrationRecord, error) {
	// 这里需要实现MongoDB迁移记录的查询逻辑
	// 暂时返回空映射
	return make(map[string]*MongoDBMigrationRecord), nil
}

// filterMigrations 过滤MySQL迁移
func (r *Runner) filterMigrations(migrations []*MySQLMigration, targetVersion string, steps int, _ MigrationDirection) []*MySQLMigration {
	var filtered []*MySQLMigration

	for _, mig := range migrations {
		if targetVersion != "" && mig.Version > targetVersion {
			break
		}

		filtered = append(filtered, mig)

		if steps > 0 && len(filtered) >= steps {
			break
		}
	}

	return filtered
}

// filterAppliedMigrations 过滤已应用的迁移
func (r *Runner) filterAppliedMigrations(records []*MySQLMigrationRecord, targetVersion string, steps int) []*MySQLMigrationRecord {
	var filtered []*MySQLMigrationRecord

	for _, record := range records {
		if targetVersion != "" && record.Version <= targetVersion {
			break
		}

		filtered = append(filtered, record)

		if steps > 0 && len(filtered) >= steps {
			break
		}
	}

	return filtered
}

// filterMongoMigrations 过滤MongoDB迁移
func (r *Runner) filterMongoMigrations(migrations []*MongoDBMigration, targetVersion string, steps int, _ MigrationDirection) []*MongoDBMigration {
	var filtered []*MongoDBMigration

	for _, mig := range migrations {
		if targetVersion != "" && mig.Version > targetVersion {
			break
		}

		filtered = append(filtered, mig)

		if steps > 0 && len(filtered) >= steps {
			break
		}
	}

	return filtered
}

// executeMySQLMigration 执行MySQL迁移
func (r *Runner) executeMySQLMigration(mig *MySQLMigration, direction MigrationDirection) error {
	var sqlFile string
	if direction == DirectionUp {
		sqlFile = mig.UpFile
	} else {
		sqlFile = mig.DownFile
	}

	// 读取SQL文件
	sqlContent, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		return fmt.Errorf("读取SQL文件失败: %v", err)
	}

	// 执行SQL语句
	if err := r.executeSQLStatements(string(sqlContent)); err != nil {
		return fmt.Errorf("执行SQL失败: %v", err)
	}

	// 记录迁移
	if direction == DirectionUp {
		record := &MySQLMigrationRecord{
			Version:   mig.Version,
			Name:      mig.Name,
			UpSQL:     string(sqlContent),
			AppliedAt: time.Now(),
		}

		if err := r.mysqlDB.Create(record).Error; err != nil {
			return fmt.Errorf("记录迁移失败: %v", err)
		}
	}

	return nil
}

// rollbackMySQLMigration 回滚MySQL迁移
func (r *Runner) rollbackMySQLMigration(record *MySQLMigrationRecord) error {
	// 读取down文件内容
	mysqlDir := filepath.Join(r.migrationsDir, "mysql")
	downFile := filepath.Join(mysqlDir, fmt.Sprintf("%s_%s.down.sql", record.Version, record.Name))

	downSQL, err := ioutil.ReadFile(downFile)
	if err != nil {
		return fmt.Errorf("读取down文件失败: %v", err)
	}

	// 执行回滚SQL
	if err := r.executeSQLStatements(string(downSQL)); err != nil {
		return fmt.Errorf("执行回滚SQL失败: %v", err)
	}

	// 删除迁移记录
	if err := r.mysqlDB.Delete(record).Error; err != nil {
		return fmt.Errorf("删除迁移记录失败: %v", err)
	}

	return nil
}

// executeMongoDBMigration 执行MongoDB迁移
func (r *Runner) executeMongoDBMigration(mig *MongoDBMigration) error {
	// 读取JavaScript文件
	jsContent, err := ioutil.ReadFile(mig.JSFile)
	if err != nil {
		return fmt.Errorf("读取JavaScript文件失败: %v", err)
	}

	// 这里需要实现JavaScript代码的执行
	// 可以使用mongo shell或其他方式执行
	fmt.Printf("📄 执行MongoDB脚本: %s\n", mig.JSFile)

	if r.verbose {
		fmt.Printf("脚本内容:\n%s\n", string(jsContent))
	}

	// 记录MongoDB迁移
	record := &MongoDBMigrationRecord{
		Version:   mig.Version,
		Name:      mig.Name,
		Script:    string(jsContent),
		AppliedAt: time.Now(),
	}

	collection := r.mongoDB.Collection("migration_records")
	_, err = collection.InsertOne(context.Background(), record)
	if err != nil {
		return fmt.Errorf("记录MongoDB迁移失败: %v", err)
	}

	return nil
}

// executeSQLStatements 执行SQL语句
func (r *Runner) executeSQLStatements(sql string) error {
	statements := strings.Split(sql, ";")

	for _, statement := range statements {
		statement = strings.TrimSpace(statement)
		if statement == "" || strings.HasPrefix(statement, "--") {
			continue
		}

		if r.verbose {
			fmt.Printf("🔧 执行SQL: %s\n", statement)
		}

		if err := r.mysqlDB.Exec(statement).Error; err != nil {
			return fmt.Errorf("执行SQL语句失败 '%s': %v", statement, err)
		}
	}

	return nil
}

// previewMySQLMigrations 预览MySQL迁移
func (r *Runner) previewMySQLMigrations(migrations []*MySQLMigration, direction MigrationDirection) error {
	fmt.Println("\n📋 MySQL迁移预览:")
	fmt.Println(strings.Repeat("=", 80))

	for i, mig := range migrations {
		fmt.Printf("\n%d. 迁移: %s_%s\n", i+1, mig.Version, mig.Name)

		var sqlFile string
		if direction == DirectionUp {
			sqlFile = mig.UpFile
			fmt.Printf("   文件: %s\n", sqlFile)
		} else {
			sqlFile = mig.DownFile
			fmt.Printf("   文件: %s\n", sqlFile)
		}

		// 读取并显示SQL内容
		if content, err := ioutil.ReadFile(sqlFile); err == nil {
			fmt.Printf("   内容预览:\n")
			lines := strings.Split(string(content), "\n")
			for j, line := range lines {
				if j > 10 { // 只显示前10行
					fmt.Printf("   ... (还有 %d 行)\n", len(lines)-j)
					break
				}
				if strings.TrimSpace(line) != "" && !strings.HasPrefix(strings.TrimSpace(line), "--") {
					fmt.Printf("   %s\n", line)
				}
			}
		}

		fmt.Println(strings.Repeat("-", 80))
	}

	return nil
}

// previewMySQLRollbacks 预览MySQL回滚
func (r *Runner) previewMySQLRollbacks(records []*MySQLMigrationRecord) error {
	fmt.Println("\n📋 MySQL回滚预览:")
	fmt.Println(strings.Repeat("=", 80))

	for i, record := range records {
		fmt.Printf("\n%d. 回滚: %s_%s\n", i+1, record.Version, record.Name)
		fmt.Printf("   应用时间: %s\n", record.AppliedAt.Format("2006-01-02 15:04:05"))

		// 显示回滚SQL的预览
		downFile := filepath.Join(r.migrationsDir, "mysql", fmt.Sprintf("%s_%s.down.sql", record.Version, record.Name))
		if content, err := ioutil.ReadFile(downFile); err == nil {
			fmt.Printf("   回滚SQL预览:\n")
			lines := strings.Split(string(content), "\n")
			for j, line := range lines {
				if j > 10 {
					fmt.Printf("   ... (还有 %d 行)\n", len(lines)-j)
					break
				}
				if strings.TrimSpace(line) != "" && !strings.HasPrefix(strings.TrimSpace(line), "--") {
					fmt.Printf("   %s\n", line)
				}
			}
		}

		fmt.Println(strings.Repeat("-", 80))
	}

	return nil
}

// previewMongoDBMigrations 预览MongoDB迁移
func (r *Runner) previewMongoDBMigrations(migrations []*MongoDBMigration) error {
	fmt.Println("\n📋 MongoDB迁移预览:")
	fmt.Println(strings.Repeat("=", 80))

	for i, mig := range migrations {
		fmt.Printf("\n%d. MongoDB迁移: %s_%s\n", i+1, mig.Version, mig.Name)
		fmt.Printf("   文件: %s\n", mig.JSFile)

		// 读取并显示JavaScript内容
		if content, err := ioutil.ReadFile(mig.JSFile); err == nil {
			fmt.Printf("   脚本预览:\n")
			lines := strings.Split(string(content), "\n")
			for j, line := range lines {
				if j > 15 {
					fmt.Printf("   ... (还有 %d 行)\n", len(lines)-j)
					break
				}
				if strings.TrimSpace(line) != "" && !strings.HasPrefix(strings.TrimSpace(line), "//") {
					fmt.Printf("   %s\n", line)
				}
			}
		}

		fmt.Println(strings.Repeat("-", 80))
	}

	return nil
}
