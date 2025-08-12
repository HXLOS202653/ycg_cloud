// Package database provides database configuration and index management functionality.
// It includes index strategies, monitoring, and optimization features for database performance.
package database

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// IndexConfig 索引配置结构
type IndexConfig struct {
	// 基础配置
	EnableFullText     bool `yaml:"enable_fulltext" json:"enable_fulltext"`
	EnablePartitioning bool `yaml:"enable_partitioning" json:"enable_partitioning"`
	OptimizeForReads   bool `yaml:"optimize_for_reads" json:"optimize_for_reads"`

	// 全文索引配置
	FullTextConfig FullTextIndexConfig `yaml:"fulltext" json:"fulltext"`

	// 性能监控配置
	MonitoringConfig IndexMonitoringConfig `yaml:"monitoring" json:"monitoring"`

	// 维护配置
	MaintenanceConfig IndexMaintenanceConfig `yaml:"maintenance" json:"maintenance"`
}

// FullTextIndexConfig 全文索引配置
type FullTextIndexConfig struct {
	MinWordLength       int    `yaml:"min_word_length" json:"min_word_length"`
	NgramTokenSize      int    `yaml:"ngram_token_size" json:"ngram_token_size"`
	QueryExpansionLimit int    `yaml:"query_expansion_limit" json:"query_expansion_limit"`
	BooleanSyntax       string `yaml:"boolean_syntax" json:"boolean_syntax"`
}

// IndexMonitoringConfig 索引监控配置
type IndexMonitoringConfig struct {
	EnableSlowQueryLog bool          `yaml:"enable_slow_query_log" json:"enable_slow_query_log"`
	SlowQueryThreshold float64       `yaml:"slow_query_threshold" json:"slow_query_threshold"` // 秒
	LogUnusedIndexes   bool          `yaml:"log_unused_indexes" json:"log_unused_indexes"`
	MonitoringInterval time.Duration `yaml:"monitoring_interval" json:"monitoring_interval"`
}

// IndexMaintenanceConfig 索引维护配置
type IndexMaintenanceConfig struct {
	AutoAnalyze            bool          `yaml:"auto_analyze" json:"auto_analyze"`
	AnalyzeInterval        time.Duration `yaml:"analyze_interval" json:"analyze_interval"`
	AutoOptimize           bool          `yaml:"auto_optimize" json:"auto_optimize"`
	OptimizeInterval       time.Duration `yaml:"optimize_interval" json:"optimize_interval"`
	CleanupOldPartitions   bool          `yaml:"cleanup_old_partitions" json:"cleanup_old_partitions"`
	PartitionRetentionDays int           `yaml:"partition_retention_days" json:"partition_retention_days"`
}

// DefaultIndexConfig 返回默认索引配置
func DefaultIndexConfig() *IndexConfig {
	return &IndexConfig{
		EnableFullText:     true,
		EnablePartitioning: true,
		OptimizeForReads:   true,

		FullTextConfig: FullTextIndexConfig{
			MinWordLength:       1,
			NgramTokenSize:      2,
			QueryExpansionLimit: 20,
			BooleanSyntax:       "+ -><()~*:\"\"&|",
		},

		MonitoringConfig: IndexMonitoringConfig{
			EnableSlowQueryLog: true,
			SlowQueryThreshold: 1.0, // 1秒
			LogUnusedIndexes:   true,
			MonitoringInterval: 24 * time.Hour, // 每天检查一次
		},

		MaintenanceConfig: IndexMaintenanceConfig{
			AutoAnalyze:            true,
			AnalyzeInterval:        7 * 24 * time.Hour, // 每周分析一次
			AutoOptimize:           true,
			OptimizeInterval:       30 * 24 * time.Hour, // 每月优化一次
			CleanupOldPartitions:   true,
			PartitionRetentionDays: 365, // 保留1年数据
		},
	}
}

// IndexStrategy 索引策略接口
type IndexStrategy interface {
	// ApplyIndexes 应用索引策略
	ApplyIndexes(db *gorm.DB) error

	// AnalyzeIndexes 分析索引使用情况
	AnalyzeIndexes(db *gorm.DB) (*IndexAnalysisResult, error)

	// OptimizeIndexes 优化索引
	OptimizeIndexes(db *gorm.DB) error

	// GetTableIndexes 获取表的索引配置
	GetTableIndexes(tableName string) []IndexDefinition
}

// IndexDefinition 索引定义
type IndexDefinition struct {
	Name    string                 `json:"name"`
	Type    string                 `json:"type"` // PRIMARY, UNIQUE, INDEX, FULLTEXT
	Table   string                 `json:"table"`
	Columns []string               `json:"columns"`
	Comment string                 `json:"comment"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// IndexAnalysisResult 索引分析结果
type IndexAnalysisResult struct {
	TotalIndexes    int                   `json:"total_indexes"`
	UsedIndexes     int                   `json:"used_indexes"`
	UnusedIndexes   []string              `json:"unused_indexes"`
	SlowQueries     []SlowQueryInfo       `json:"slow_queries"`
	Recommendations []IndexRecommendation `json:"recommendations"`
	AnalysisTime    time.Time             `json:"analysis_time"`
}

// SlowQueryInfo 慢查询信息
type SlowQueryInfo struct {
	Query         string  `json:"query"`
	ExecutionTime float64 `json:"execution_time"` // 秒
	RowsExamined  int64   `json:"rows_examined"`
	RowsSent      int64   `json:"rows_sent"`
	Frequency     int     `json:"frequency"`
}

// IndexRecommendation 索引建议
type IndexRecommendation struct {
	Type       string  `json:"type"` // CREATE, DROP, MODIFY
	Table      string  `json:"table"`
	IndexName  string  `json:"index_name"`
	Columns    string  `json:"columns"`
	Reason     string  `json:"reason"`
	Impact     string  `json:"impact"`     // HIGH, MEDIUM, LOW
	Confidence float64 `json:"confidence"` // 0-1
}

// CoreTableIndexStrategy 核心表索引策略实现
type CoreTableIndexStrategy struct {
	config *IndexConfig
}

// NewCoreTableIndexStrategy 创建核心表索引策略
func NewCoreTableIndexStrategy(config *IndexConfig) *CoreTableIndexStrategy {
	if config == nil {
		config = DefaultIndexConfig()
	}
	return &CoreTableIndexStrategy{
		config: config,
	}
}

// ApplyIndexes 应用索引策略
func (s *CoreTableIndexStrategy) ApplyIndexes(db *gorm.DB) error {
	// 应用MySQL配置
	if err := s.applyMySQLConfig(db); err != nil {
		return fmt.Errorf("failed to apply MySQL config: %w", err)
	}

	// 应用表索引
	tables := []string{
		"users", "user_sessions", "user_settings",
		"files", "file_versions", "upload_tasks", "file_chunks",
		"folders", "system_configs", "file_shares", "share_access_logs",
		"file_permissions", "teams", "team_members", "team_invitations",
		"notifications", "operation_logs", "security_logs",
		"chat_rooms", "chat_room_members", "chat_messages",
		"file_contents", "search_logs",
	}

	for _, table := range tables {
		if err := s.applyTableIndexes(db, table); err != nil {
			return fmt.Errorf("failed to apply indexes for table %s: %w", table, err)
		}
	}

	return nil
}

// applyMySQLConfig 应用MySQL配置
func (s *CoreTableIndexStrategy) applyMySQLConfig(db *gorm.DB) error {
	// 关键配置 - 失败时应该返回错误
	criticalConfigs := []string{
		"SET GLOBAL innodb_stats_persistent = ON",
		"SET GLOBAL innodb_stats_auto_recalc = ON",
		"SET GLOBAL optimizer_switch = 'index_merge=on,index_merge_union=on,index_merge_sort_union=on,index_merge_intersection=on'",
	}

	// 可选配置 - 失败时只记录警告
	optionalConfigs := []string{
		fmt.Sprintf("SET GLOBAL ft_min_word_len = %d", s.config.FullTextConfig.MinWordLength),
		fmt.Sprintf("SET GLOBAL ngram_token_size = %d", s.config.FullTextConfig.NgramTokenSize),
		fmt.Sprintf("SET GLOBAL ft_query_expansion_limit = %d", s.config.FullTextConfig.QueryExpansionLimit),
		fmt.Sprintf("SET GLOBAL ft_boolean_syntax = '%s'", s.config.FullTextConfig.BooleanSyntax),
	}

	if s.config.MonitoringConfig.EnableSlowQueryLog {
		optionalConfigs = append(optionalConfigs, []string{
			"SET GLOBAL slow_query_log = ON",
			fmt.Sprintf("SET GLOBAL long_query_time = %.1f", s.config.MonitoringConfig.SlowQueryThreshold),
			"SET GLOBAL log_queries_not_using_indexes = ON",
		}...)
	}

	// 应用关键配置
	for _, config := range criticalConfigs {
		if err := db.Exec(config).Error; err != nil {
			return fmt.Errorf("failed to apply critical MySQL config '%s': %w", config, err)
		}
	}

	// 应用可选配置
	for _, config := range optionalConfigs {
		if err := db.Exec(config).Error; err != nil {
			// 某些配置可能需要特殊权限，记录警告但不中断
			fmt.Printf("Warning: Failed to apply optional config '%s': %v\n", config, err)
		}
	}

	return nil
}

// applyTableIndexes 应用表索引
func (s *CoreTableIndexStrategy) applyTableIndexes(db *gorm.DB, tableName string) error {
	indexes := s.GetTableIndexes(tableName)

	for _, index := range indexes {
		if err := s.createIndexIfNotExists(db, index); err != nil {
			return fmt.Errorf("failed to create index %s: %w", index.Name, err)
		}
	}

	return nil
}

// createIndexIfNotExists 创建索引（如果不存在）
func (s *CoreTableIndexStrategy) createIndexIfNotExists(db *gorm.DB, index IndexDefinition) error {
	// 检查索引是否存在
	var count int64
	query := `
		SELECT COUNT(*) FROM information_schema.statistics 
		WHERE table_schema = DATABASE() 
		AND table_name = ? 
		AND index_name = ?
	`

	if err := db.Raw(query, index.Table, index.Name).Scan(&count).Error; err != nil {
		return fmt.Errorf("failed to check index existence: %w", err)
	}

	if count > 0 {
		// 索引已存在
		return nil
	}

	// 构建创建索引的SQL
	sql := s.buildCreateIndexSQL(index)

	if err := db.Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to execute create index SQL: %w", err)
	}

	fmt.Printf("Created index: %s on table %s\n", index.Name, index.Table)
	return nil
}

// buildCreateIndexSQL 构建创建索引的SQL
func (s *CoreTableIndexStrategy) buildCreateIndexSQL(index IndexDefinition) string {
	var sql string
	columnsStr := fmt.Sprintf("(%s)", joinStrings(index.Columns, ", "))

	switch index.Type {
	case "FULLTEXT":
		sql = fmt.Sprintf("ALTER TABLE %s ADD FULLTEXT INDEX %s %s",
			index.Table, index.Name, columnsStr)
		if parser, ok := index.Options["parser"]; ok {
			sql += fmt.Sprintf(" WITH PARSER %s", parser)
		}
	case "UNIQUE":
		sql = fmt.Sprintf("ALTER TABLE %s ADD UNIQUE INDEX %s %s",
			index.Table, index.Name, columnsStr)
	default: // INDEX
		sql = fmt.Sprintf("ALTER TABLE %s ADD INDEX %s %s",
			index.Table, index.Name, columnsStr)
	}

	if index.Comment != "" {
		sql += fmt.Sprintf(" COMMENT '%s'", index.Comment)
	}

	return sql
}

// GetTableIndexes 获取表的索引配置
func (s *CoreTableIndexStrategy) GetTableIndexes(tableName string) []IndexDefinition {
	switch tableName {
	case "users":
		return s.getUsersIndexes()
	case "user_sessions":
		return s.getUserSessionsIndexes()
	case "files":
		return s.getFilesIndexes()
	case "file_versions":
		return s.getFileVersionsIndexes()
	case "upload_tasks":
		return s.getUploadTasksIndexes()
	case "file_chunks":
		return s.getFileChunksIndexes()
	case "file_shares":
		return s.getFileSharesIndexes()
	case "teams":
		return s.getTeamsIndexes()
	default:
		return []IndexDefinition{}
	}
}

// getUsersIndexes 用户表索引配置
func (s *CoreTableIndexStrategy) getUsersIndexes() []IndexDefinition {
	indexes := []IndexDefinition{
		{
			Name: "idx_users_username", Type: "INDEX", Table: "users",
			Columns: []string{"username"}, Comment: "用户名登录查询",
		},
		{
			Name: "idx_users_email", Type: "INDEX", Table: "users",
			Columns: []string{"email"}, Comment: "邮箱登录查询",
		},
		{
			Name: "idx_users_phone", Type: "INDEX", Table: "users",
			Columns: []string{"phone"}, Comment: "手机号登录查询",
		},
		{
			Name: "idx_users_status", Type: "INDEX", Table: "users",
			Columns: []string{"status"}, Comment: "用户状态筛选",
		},
		{
			Name: "idx_users_role", Type: "INDEX", Table: "users",
			Columns: []string{"role"}, Comment: "角色权限查询",
		},
		{
			Name: "idx_users_created_at", Type: "INDEX", Table: "users",
			Columns: []string{"created_at"}, Comment: "注册时间排序",
		},
		{
			Name: "idx_users_storage_used", Type: "INDEX", Table: "users",
			Columns: []string{"storage_used"}, Comment: "存储配额管理",
		},
	}

	return indexes
}

// getUserSessionsIndexes 用户会话表索引配置
func (s *CoreTableIndexStrategy) getUserSessionsIndexes() []IndexDefinition {
	return []IndexDefinition{
		{
			Name: "idx_user_sessions_user_id", Type: "INDEX", Table: "user_sessions",
			Columns: []string{"user_id"}, Comment: "用户会话查询",
		},
		{
			Name: "idx_user_sessions_session_token", Type: "INDEX", Table: "user_sessions",
			Columns: []string{"session_token"}, Comment: "Token验证",
		},
		{
			Name: "idx_user_sessions_expires_at", Type: "INDEX", Table: "user_sessions",
			Columns: []string{"expires_at"}, Comment: "过期会话清理",
		},
		{
			Name: "idx_user_sessions_is_active", Type: "INDEX", Table: "user_sessions",
			Columns: []string{"is_active"}, Comment: "活跃会话筛选",
		},
	}
}

// getFilesIndexes 文件表索引配置
func (s *CoreTableIndexStrategy) getFilesIndexes() []IndexDefinition {
	indexes := []IndexDefinition{
		{
			Name: "idx_files_user_parent_deleted", Type: "INDEX", Table: "files",
			Columns: []string{"user_id", "parent_id", "is_deleted"},
			Comment: "文件树查询核心索引",
		},
		{
			Name: "idx_files_filename", Type: "INDEX", Table: "files",
			Columns: []string{"filename"}, Comment: "文件名搜索",
		},
		{
			Name: "idx_files_file_type", Type: "INDEX", Table: "files",
			Columns: []string{"file_type"}, Comment: "文件类型筛选",
		},
		{
			Name: "idx_files_md5_hash", Type: "INDEX", Table: "files",
			Columns: []string{"md5_hash"}, Comment: "MD5去重查询",
		},
		{
			Name: "idx_files_is_deleted", Type: "INDEX", Table: "files",
			Columns: []string{"is_deleted"}, Comment: "软删除筛选",
		},
		{
			Name: "idx_files_created_at", Type: "INDEX", Table: "files",
			Columns: []string{"created_at"}, Comment: "创建时间排序",
		},
	}

	// 添加全文索引（如果启用）
	if s.config.EnableFullText {
		indexes = append(indexes, IndexDefinition{
			Name: "idx_files_fulltext_search", Type: "FULLTEXT", Table: "files",
			Columns: []string{"filename", "description"},
			Comment: "文件名描述全文搜索",
			Options: map[string]interface{}{"parser": "ngram"},
		})
	}

	return indexes
}

// getFileVersionsIndexes 文件版本表索引配置
func (s *CoreTableIndexStrategy) getFileVersionsIndexes() []IndexDefinition {
	return []IndexDefinition{
		{
			Name: "idx_file_versions_file_id", Type: "INDEX", Table: "file_versions",
			Columns: []string{"file_id"}, Comment: "文件版本查询",
		},
		{
			Name: "idx_file_versions_version_number", Type: "INDEX", Table: "file_versions",
			Columns: []string{"version_number"}, Comment: "版本号排序",
		},
		{
			Name: "idx_file_versions_is_active", Type: "INDEX", Table: "file_versions",
			Columns: []string{"is_active"}, Comment: "当前版本查询",
		},
		{
			Name: "idx_file_versions_created_by", Type: "INDEX", Table: "file_versions",
			Columns: []string{"created_by"}, Comment: "创建者查询",
		},
	}
}

// getUploadTasksIndexes 上传任务表索引配置
func (s *CoreTableIndexStrategy) getUploadTasksIndexes() []IndexDefinition {
	return []IndexDefinition{
		{
			Name: "idx_upload_tasks_upload_id", Type: "INDEX", Table: "upload_tasks",
			Columns: []string{"upload_id"}, Comment: "上传任务查询",
		},
		{
			Name: "idx_upload_tasks_user_id", Type: "INDEX", Table: "upload_tasks",
			Columns: []string{"user_id"}, Comment: "用户上传任务",
		},
		{
			Name: "idx_upload_tasks_status", Type: "INDEX", Table: "upload_tasks",
			Columns: []string{"status"}, Comment: "任务状态筛选",
		},
		{
			Name: "idx_upload_tasks_md5_hash", Type: "INDEX", Table: "upload_tasks",
			Columns: []string{"md5_hash"}, Comment: "秒传支持",
		},
	}
}

// getFileChunksIndexes 文件分片表索引配置
func (s *CoreTableIndexStrategy) getFileChunksIndexes() []IndexDefinition {
	return []IndexDefinition{
		{
			Name: "idx_file_chunks_upload_id", Type: "INDEX", Table: "file_chunks",
			Columns: []string{"upload_id"}, Comment: "上传任务分片查询",
		},
		{
			Name: "idx_file_chunks_chunk_number", Type: "INDEX", Table: "file_chunks",
			Columns: []string{"chunk_number"}, Comment: "分片序号排序",
		},
		{
			Name: "idx_file_chunks_status", Type: "INDEX", Table: "file_chunks",
			Columns: []string{"status"}, Comment: "分片状态筛选",
		},
	}
}

// getFileSharesIndexes 文件分享表索引配置
func (s *CoreTableIndexStrategy) getFileSharesIndexes() []IndexDefinition {
	indexes := []IndexDefinition{
		{
			Name: "idx_file_shares_share_code", Type: "INDEX", Table: "file_shares",
			Columns: []string{"share_code"}, Comment: "分享码访问查询",
		},
		{
			Name: "idx_file_shares_file_id", Type: "INDEX", Table: "file_shares",
			Columns: []string{"file_id"}, Comment: "文件分享查询",
		},
		{
			Name: "idx_file_shares_user_id", Type: "INDEX", Table: "file_shares",
			Columns: []string{"user_id"}, Comment: "用户分享管理",
		},
		{
			Name: "idx_file_shares_status", Type: "INDEX", Table: "file_shares",
			Columns: []string{"status"}, Comment: "分享状态管理",
		},
	}

	// 添加全文索引（如果启用）
	if s.config.EnableFullText {
		indexes = append(indexes, IndexDefinition{
			Name: "idx_file_shares_fulltext_search", Type: "FULLTEXT", Table: "file_shares",
			Columns: []string{"share_name", "share_description"},
			Comment: "分享内容全文搜索",
			Options: map[string]interface{}{"parser": "ngram"},
		})
	}

	return indexes
}

// getTeamsIndexes 团队表索引配置
func (s *CoreTableIndexStrategy) getTeamsIndexes() []IndexDefinition {
	indexes := []IndexDefinition{
		{
			Name: "idx_teams_name", Type: "INDEX", Table: "teams",
			Columns: []string{"name"}, Comment: "团队名称搜索",
		},
		{
			Name: "idx_teams_owner_id", Type: "INDEX", Table: "teams",
			Columns: []string{"owner_id"}, Comment: "团队所有者查询",
		},
		{
			Name: "idx_teams_status", Type: "INDEX", Table: "teams",
			Columns: []string{"status"}, Comment: "团队状态管理",
		},
		{
			Name: "idx_teams_is_active", Type: "INDEX", Table: "teams",
			Columns: []string{"is_active"}, Comment: "活跃团队筛选",
		},
	}

	// 添加全文索引（如果启用）
	if s.config.EnableFullText {
		indexes = append(indexes, IndexDefinition{
			Name: "idx_teams_fulltext_search", Type: "FULLTEXT", Table: "teams",
			Columns: []string{"name", "description"},
			Comment: "团队名称描述全文搜索",
			Options: map[string]interface{}{"parser": "ngram"},
		})
	}

	return indexes
}

// AnalyzeIndexes 分析索引使用情况
func (s *CoreTableIndexStrategy) AnalyzeIndexes(db *gorm.DB) (*IndexAnalysisResult, error) {
	result := &IndexAnalysisResult{
		AnalysisTime: time.Now(),
	}

	// 获取索引使用统计
	if err := s.analyzeIndexUsage(db, result); err != nil {
		return nil, fmt.Errorf("failed to analyze index usage: %w", err)
	}

	// 获取慢查询信息
	if err := s.analyzeSlowQueries(db, result); err != nil {
		return nil, fmt.Errorf("failed to analyze slow queries: %w", err)
	}

	// 生成优化建议
	s.generateRecommendations(result)

	return result, nil
}

// analyzeIndexUsage 分析索引使用情况
func (s *CoreTableIndexStrategy) analyzeIndexUsage(db *gorm.DB, result *IndexAnalysisResult) error {
	// 查询索引使用统计
	query := `
		SELECT 
		    t.table_name,
		    t.index_name,
		    COALESCE(p.count_read, 0) + COALESCE(p.count_write, 0) as usage_count
		FROM information_schema.statistics t
		LEFT JOIN performance_schema.table_io_waits_summary_by_index_usage p 
		    ON t.table_schema = p.object_schema 
		    AND t.table_name = p.object_name 
		    AND t.index_name = p.index_name
		WHERE t.table_schema = DATABASE()
		    AND t.index_name != 'PRIMARY'
		GROUP BY t.table_name, t.index_name
		ORDER BY usage_count DESC
	`

	var indexStats []struct {
		TableName  string `json:"table_name"`
		IndexName  string `json:"index_name"`
		UsageCount int64  `json:"usage_count"`
	}

	if err := db.Raw(query).Scan(&indexStats).Error; err != nil {
		return err
	}

	result.TotalIndexes = len(indexStats)

	for _, stat := range indexStats {
		if stat.UsageCount == 0 {
			result.UnusedIndexes = append(result.UnusedIndexes,
				fmt.Sprintf("%s.%s", stat.TableName, stat.IndexName))
		} else {
			result.UsedIndexes++
		}
	}

	return nil
}

// analyzeSlowQueries 分析慢查询
func (s *CoreTableIndexStrategy) analyzeSlowQueries(db *gorm.DB, result *IndexAnalysisResult) error {
	if !s.config.MonitoringConfig.EnableSlowQueryLog {
		return nil
	}

	query := `
		SELECT 
		    digest_text,
		    count_star as frequency,
		    avg_timer_wait/1000000000 as avg_time_sec,
		    sum_rows_examined/count_star as avg_rows_examined,
		    sum_rows_sent/count_star as avg_rows_sent
		FROM performance_schema.events_statements_summary_by_digest
		WHERE digest_text IS NOT NULL
		    AND avg_timer_wait > ?
		ORDER BY avg_timer_wait DESC
		LIMIT 10
	`

	threshold := s.config.MonitoringConfig.SlowQueryThreshold * 1000000000 // 转换为纳秒

	var slowQueries []SlowQueryInfo
	if err := db.Raw(query, threshold).Scan(&slowQueries).Error; err != nil {
		return err
	}

	result.SlowQueries = slowQueries
	return nil
}

// generateRecommendations 生成优化建议
func (s *CoreTableIndexStrategy) generateRecommendations(result *IndexAnalysisResult) {
	recommendations := make([]IndexRecommendation, 0, 10) // 预分配合理容量

	// 建议删除未使用的索引
	for _, unusedIndex := range result.UnusedIndexes {
		recommendations = append(recommendations, IndexRecommendation{
			Type:       "DROP",
			IndexName:  unusedIndex,
			Reason:     "索引未被使用，删除可以提高写入性能",
			Impact:     "LOW",
			Confidence: 0.8,
		})
	}

	// 基于慢查询的索引建议
	for _, slowQuery := range result.SlowQueries {
		if slowQuery.RowsExamined > slowQuery.RowsSent*100 {
			recommendations = append(recommendations, IndexRecommendation{
				Type:       "CREATE",
				Reason:     "查询扫描行数过多，可能缺少合适的索引",
				Impact:     "HIGH",
				Confidence: 0.6,
			})
		}
	}

	result.Recommendations = recommendations
}

// OptimizeIndexes 优化索引
func (s *CoreTableIndexStrategy) OptimizeIndexes(db *gorm.DB) error {
	if !s.config.MaintenanceConfig.AutoOptimize {
		return nil
	}

	// 更新索引统计信息
	if s.config.MaintenanceConfig.AutoAnalyze {
		tables := []string{"users", "files", "file_versions", "upload_tasks", "file_shares", "teams"}
		for _, table := range tables {
			if isValidSystemTableName(table) {
				if err := db.Exec(fmt.Sprintf("ANALYZE TABLE %s", table)).Error; err != nil {
					fmt.Printf("Warning: Failed to analyze table %s: %v\n", table, err)
				}
			}
		}
	}

	// 优化表（重建索引）
	optimizeTables := []string{"files", "file_shares", "teams"}
	for _, table := range optimizeTables {
		if isValidSystemTableName(table) {
			if err := db.Exec(fmt.Sprintf("OPTIMIZE TABLE %s", table)).Error; err != nil {
				fmt.Printf("Warning: Failed to optimize table %s: %v\n", table, err)
			}
		}
	}

	return nil
}

// joinStrings 连接字符串数组
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// isValidSystemTableName validates system table names to prevent SQL injection
func isValidSystemTableName(tableName string) bool {
	// Whitelist of allowed system table names
	allowedTables := map[string]bool{
		"users":          true,
		"files":          true,
		"file_versions":  true,
		"upload_tasks":   true,
		"file_shares":    true,
		"teams":          true,
		"team_members":   true,
		"notifications":  true,
		"operation_logs": true,
		"security_logs":  true,
	}

	return allowedTables[tableName]
}
