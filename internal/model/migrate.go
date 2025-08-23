package model

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移所有模型
func AutoMigrate(db *gorm.DB) error {
	log.Println("开始数据库迁移...")

	// 定义迁移顺序，确保外键依赖正确
	models := []interface{}{
		// 基础模型（无外键依赖）
		&User{},
		&SystemConfig{},
		&PermissionTemplate{},
		&Role{},
		&Team{},

		// 依赖基础模型的模型
		&StorageConfig{},
		&configHistory{},
		&File{},
		&TeamMember{},
		&TeamFile{},
		&TeamRole{},
		&Conversation{},
		&RecycleItem{},
		&RecycleBin{},

		// 权限相关模型
		&templatePermission{},
		&userPermission{},
		&filePermission{},
		&userRole{},

		// 消息相关模型
		&ConversationMember{},
		&Message{},
		&MessageReadReceipt{},

		// 日志相关模型
		&OperationLog{},
		&SystemLog{},
		&SecurityLog{},

		// 回收站日志
		&RecycleLog{},
	}

	// 执行迁移
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("迁移模型 %T 失败: %w", model, err)
		}
		log.Printf("成功迁移模型: %T", model)
	}

	log.Println("数据库迁移完成")
	return nil
}

// CreateIndexes 创建额外的索引
func CreateIndexes(db *gorm.DB) error {
	log.Println("开始创建额外索引...")

	// 复合索引定义
	indexes := []struct {
		tableName string
		indexName string
		columns   []string
	}{
		// 文件表复合索引
		{"files", "idx_files_owner_status_created", []string{"owner_id", "status", "created_at"}},
		{"files", "idx_files_parent_status_type", []string{"parent_id", "status", "file_type"}},
		{"files", "idx_files_owner_name_type", []string{"owner_id", "name", "file_type"}},
		{"files", "idx_files_owner_md5", []string{"owner_id", "md5_hash"}},

		// 权限表复合索引
		{"user_permissions", "idx_user_permissions_user_expires", []string{"user_id", "expires_at"}},
		{"file_permissions", "idx_file_permissions_file_user_expires", []string{"file_id", "user_id", "expires_at"}},
		{"team_members", "idx_team_members_team_user_status", []string{"team_id", "user_id", "status"}},

		// 日志表复合索引
		{"operation_logs", "idx_operation_logs_user_created", []string{"user_id", "created_at"}},
		{"system_logs", "idx_system_logs_level_module_created", []string{"level", "module", "created_at"}},
		{"security_logs", "idx_security_logs_user_action_created", []string{"user_id", "action_type", "created_at"}},
	}

	// 创建复合索引
	for _, idx := range indexes {
		indexSQL := fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s (%s)",
			idx.indexName,
			idx.tableName,
			joinColumns(idx.columns))

		if err := db.Exec(indexSQL).Error; err != nil {
			log.Printf("创建索引 %s 失败: %v", idx.indexName, err)
			// 继续创建其他索引，不中断流程
		} else {
			log.Printf("成功创建索引: %s", idx.indexName)
		}
	}

	log.Println("索引创建完成")
	return nil
}

// joinColumns 连接列名
func joinColumns(columns []string) string {
	result := ""
	for i, col := range columns {
		if i > 0 {
			result += ", "
		}
		result += col
	}
	return result
}

// DropAllTables 删除所有表（用于测试）
func DropAllTables(db *gorm.DB) error {
	log.Println("开始删除所有表...")

	// 按相反顺序删除表，避免外键约束问题
	tables := []string{
		"recycle_logs",
		"security_logs",
		"system_logs",
		"operation_logs",
		"message_read_receipts",
		"messages",
		"conversation_members",
		"user_roles",
		"file_permissions",
		"user_permissions",
		"template_permissions",
		"recycle_bins",
		"recycle_items",
		"conversations",
		"team_roles",
		"team_files",
		"team_members",
		"files",
		"config_history",
		"storage_configs",
		"teams",
		"roles",
		"permission_templates",
		"system_configs",
		"users",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)).Error; err != nil {
			log.Printf("删除表 %s 失败: %v", table, err)
		} else {
			log.Printf("成功删除表: %s", table)
		}
	}

	log.Println("所有表删除完成")
	return nil
}

// ValidateModels 验证模型定义
func ValidateModels(db *gorm.DB) error {
	log.Println("开始验证模型定义...")

	// 检查表是否存在
	tables := []string{
		"users", "files", "teams", "team_members", "team_files",
		"permission_templates", "user_permissions", "file_permissions",
		"roles", "user_roles", "conversations", "messages",
		"recycle_items", "recycle_bins", "system_configs",
		"operation_logs", "system_logs", "security_logs",
	}

	for _, table := range tables {
		if !db.Migrator().HasTable(table) {
			return fmt.Errorf("表 %s 不存在", table)
		}
		log.Printf("表 %s 验证通过", table)
	}

	// 检查关键索引是否存在
	indexChecks := []struct {
		table string
		index string
	}{
		{"users", "idx_users_username"},
		{"users", "idx_users_email"},
		{"files", "idx_files_owner_id"},
		{"files", "idx_files_parent_id"},
		{"team_members", "idx_team_members_team_id"},
		{"team_members", "idx_team_members_user_id"},
	}

	for _, check := range indexChecks {
		if !db.Migrator().HasIndex(check.table, check.index) {
			log.Printf("警告: 表 %s 缺少索引 %s", check.table, check.index)
		} else {
			log.Printf("索引 %s.%s 验证通过", check.table, check.index)
		}
	}

	log.Println("模型验证完成")
	return nil
}
