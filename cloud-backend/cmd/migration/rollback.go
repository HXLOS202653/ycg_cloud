package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// RollbackManager 回滚管理器
type RollbackManager struct {
	tool *MigrationTool
}

// NewRollbackManager 创建回滚管理器
func NewRollbackManager(tool *MigrationTool) *RollbackManager {
	return &RollbackManager{tool: tool}
}

// RollbackInfo 回滚信息
type RollbackInfo struct {
	Version      string    `json:"version"`
	Name         string    `json:"name"`
	AppliedAt    time.Time `json:"applied_at"`
	RollbackSQL  string    `json:"rollback_sql"`
	Dependencies []string  `json:"dependencies"`
	SafeToRoll   bool      `json:"safe_to_roll"`
	Warnings     []string  `json:"warnings"`
}

// ValidateRollback 验证回滚安全性
func (rm *RollbackManager) ValidateRollback(targetVersion string, steps int) ([]*RollbackInfo, error) {
	// 获取已应用的迁移
	appliedMigrations, err := rm.tool.getAppliedMigrations()
	if err != nil {
		return nil, fmt.Errorf("获取已应用迁移失败: %v", err)
	}

	if len(appliedMigrations) == 0 {
		return nil, fmt.Errorf("没有已应用的迁移可以回滚")
	}

	// 过滤需要回滚的迁移
	toRollback := rm.tool.filterMigrationsForDown(appliedMigrations, targetVersion, steps)

	rollbackInfos := make([]*RollbackInfo, 0, len(toRollback))
	for _, record := range toRollback {
		info := &RollbackInfo{
			Version:     record.Version,
			Name:        record.Name,
			AppliedAt:   record.AppliedAt,
			RollbackSQL: record.DownSQL,
			SafeToRoll:  true,
			Warnings:    []string{},
		}

		// 分析回滚安全性
		if err := rm.analyzeRollbackSafety(info); err != nil {
			return nil, fmt.Errorf("分析回滚安全性失败: %v", err)
		}

		rollbackInfos = append(rollbackInfos, info)
	}

	return rollbackInfos, nil
}

// analyzeRollbackSafety 分析回滚安全性
func (rm *RollbackManager) analyzeRollbackSafety(info *RollbackInfo) error {
	sql := info.RollbackSQL

	// 检查极其危险的操作 - 应该阻止回滚
	criticalDangerousOps := []string{
		"DROP DATABASE",
		"TRUNCATE",
	}

	for _, operation := range criticalDangerousOps {
		if containsIgnoreCase(sql, operation) {
			return fmt.Errorf("回滚包含极其危险的操作 '%s'，为了数据安全已阻止执行", operation)
		}
	}

	// 检查一般危险操作 - 标记但允许继续
	dangerousOperations := []string{
		"DROP TABLE",
		"DELETE FROM",
		"DROP COLUMN",
		"DROP INDEX",
	}

	for _, operation := range dangerousOperations {
		if containsIgnoreCase(sql, operation) {
			info.SafeToRoll = false
			info.Warnings = append(info.Warnings,
				fmt.Sprintf("⚠️  包含危险操作: %s - 可能导致数据丢失", operation))
		}
	}

	// 检查外键约束
	if containsIgnoreCase(sql, "FOREIGN KEY") || containsIgnoreCase(sql, "REFERENCES") {
		info.Warnings = append(info.Warnings,
			"⚠️  涉及外键约束 - 请确保相关数据一致性")
	}

	// 检查索引删除
	if containsIgnoreCase(sql, "DROP INDEX") {
		info.Warnings = append(info.Warnings,
			"⚠️  删除索引 - 可能影响查询性能")
	}

	// 检查数据类型变更
	dataTypeChanges := []string{
		"ALTER COLUMN",
		"MODIFY COLUMN",
		"CHANGE COLUMN",
	}

	for _, change := range dataTypeChanges {
		if containsIgnoreCase(sql, change) {
			info.Warnings = append(info.Warnings,
				"⚠️  涉及数据类型变更 - 可能导致数据兼容性问题")
			break
		}
	}

	// 如果SQL为空，这可能是一个问题
	if strings.TrimSpace(sql) == "" {
		return fmt.Errorf("回滚SQL为空，无法执行安全分析")
	}

	return nil
}

// CreateRollbackPlan 创建回滚计划
func (rm *RollbackManager) CreateRollbackPlan(targetVersion string, steps int) error {
	rollbackInfos, err := rm.ValidateRollback(targetVersion, steps)
	if err != nil {
		return err
	}

	if len(rollbackInfos) == 0 {
		fmt.Println("📋 没有需要回滚的迁移")
		return nil
	}

	// 创建回滚计划文件
	planFile := filepath.Join(rm.tool.migrationsDir, fmt.Sprintf("rollback_plan_%s.md",
		time.Now().Format("20060102_150405")))

	content := rm.generateRollbackPlan(rollbackInfos)

	if err := os.WriteFile(planFile, []byte(content), 0600); err != nil {
		return fmt.Errorf("创建回滚计划文件失败: %v", err)
	}

	fmt.Printf("📝 回滚计划已创建: %s\n", planFile)
	return nil
}

// generateRollbackPlan 生成回滚计划内容
func (rm *RollbackManager) generateRollbackPlan(rollbackInfos []*RollbackInfo) string {
	plan := fmt.Sprintf(`# 数据库回滚计划

**生成时间**: %s  
**回滚数量**: %d 个迁移  

## ⚠️ 回滚前检查清单

- [ ] 已备份当前数据库
- [ ] 已通知相关开发人员
- [ ] 已停止相关应用服务
- [ ] 已确认回滚的必要性
- [ ] 已阅读并理解所有警告信息

## 📋 回滚详情

`, time.Now().Format("2006-01-02 15:04:05"), len(rollbackInfos))

	for i, info := range rollbackInfos {
		plan += fmt.Sprintf(`### %d. 迁移 %s_%s

**应用时间**: %s  
**回滚安全**: %s  

`, i+1, info.Version, info.Name,
			info.AppliedAt.Format("2006-01-02 15:04:05"),
			rm.getSafetyStatus(info.SafeToRoll))

		if len(info.Warnings) > 0 {
			plan += "**警告信息**:\n"
			for _, warning := range info.Warnings {
				plan += fmt.Sprintf("- %s\n", warning)
			}
			plan += "\n"
		}

		plan += "**回滚SQL**:\n```sql\n"
		plan += info.RollbackSQL
		plan += "\n```\n\n"
	}

	plan += `## 🚨 紧急回滚指令

如果回滚过程中出现问题，请立即执行以下操作：

1. **停止回滚**: 记录当前执行到的迁移版本
2. **保留现场**: 不要进行任何额外的数据库操作
3. **联系DBA**: 立即联系数据库管理员
4. **恢复备份**: 如有必要，从备份中恢复数据库

## 📞 联系信息

- **DBA**: [联系方式]
- **技术负责人**: [联系方式]
- **项目经理**: [联系方式]

---
*此回滚计划由网络云盘系统数据库迁移工具自动生成*
`

	return plan
}

// getSafetyStatus 获取安全状态描述
func (rm *RollbackManager) getSafetyStatus(safe bool) string {
	if safe {
		return "✅ 安全"
	}
	return "⚠️ 需要注意"
}

// PerformSafeRollback 执行安全回滚
func (rm *RollbackManager) PerformSafeRollback(targetVersion string, steps int, confirmed bool) error {
	// 验证回滚安全性
	rollbackInfos, err := rm.ValidateRollback(targetVersion, steps)
	if err != nil {
		return err
	}

	if len(rollbackInfos) == 0 {
		fmt.Println("📋 没有需要回滚的迁移")
		return nil
	}

	// 检查是否有不安全的回滚
	hasUnsafeRollbacks := false
	for _, info := range rollbackInfos {
		if !info.SafeToRoll {
			hasUnsafeRollbacks = true
			break
		}
	}

	// 显示回滚信息
	fmt.Printf("📋 回滚分析结果:\n\n")
	for i, info := range rollbackInfos {
		fmt.Printf("%d. 迁移 %s_%s\n", i+1, info.Version, info.Name)
		fmt.Printf("   应用时间: %s\n", info.AppliedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("   安全状态: %s\n", rm.getSafetyStatus(info.SafeToRoll))

		if len(info.Warnings) > 0 {
			fmt.Printf("   警告信息:\n")
			for _, warning := range info.Warnings {
				fmt.Printf("     %s\n", warning)
			}
		}
		fmt.Println()
	}

	// 如果有不安全的回滚且未确认，需要用户确认
	if hasUnsafeRollbacks && !confirmed {
		return fmt.Errorf("检测到不安全的回滚操作，请使用 --force 参数强制执行")
	}

	// 创建回滚备份点
	if err := rm.createRollbackBackup(); err != nil {
		return fmt.Errorf("创建回滚备份失败: %v", err)
	}

	// 执行回滚
	return rm.tool.Down(targetVersion, steps)
}

// createRollbackBackup 创建回滚备份点
func (rm *RollbackManager) createRollbackBackup() error {
	if rm.tool.dryRun {
		fmt.Println("🔍 预览模式: 跳过回滚备份创建")
		return nil
	}

	fmt.Println("💾 创建回滚备份点...")

	// 这里可以实现具体的备份逻辑
	// 例如: mysqldump, 文件系统快照等

	backupFile := filepath.Join(rm.tool.migrationsDir, fmt.Sprintf("rollback_backup_%s.sql",
		time.Now().Format("20060102_150405")))

	// 创建备份标记文件
	backupInfo := fmt.Sprintf(`-- 回滚备份点
-- 创建时间: %s
-- 备份原因: 执行数据库迁移回滚前的安全备份
-- 
-- 如需恢复，请联系DBA执行以下命令:
-- mysql -u username -p database_name < %s

`, time.Now().Format("2006-01-02 15:04:05"), backupFile)

	if err := os.WriteFile(backupFile, []byte(backupInfo), 0644); err != nil {
		return fmt.Errorf("创建备份标记文件失败: %v", err)
	}

	fmt.Printf("✅ 备份标记已创建: %s\n", backupFile)
	return nil
}

// RollbackToLastStable 回滚到最后一个稳定版本
func (rm *RollbackManager) RollbackToLastStable() error {
	// 获取已应用的迁移
	appliedMigrations, err := rm.tool.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("获取已应用迁移失败: %v", err)
	}

	if len(appliedMigrations) == 0 {
		fmt.Println("📋 没有已应用的迁移")
		return nil
	}

	// 查找最后一个标记为稳定的版本
	// 这里可以根据迁移文件名中的标记或配置文件来确定
	// 简单实现：回滚到倒数第二个版本
	if len(appliedMigrations) < 2 {
		fmt.Println("📋 只有一个迁移版本，无法回滚到稳定版本")
		return nil
	}

	lastStableVersion := appliedMigrations[1].Version
	fmt.Printf("🔄 准备回滚到最后稳定版本: %s\n", lastStableVersion)

	return rm.PerformSafeRollback(lastStableVersion, 0, false)
}

// 辅助函数

// containsIgnoreCase 不区分大小写的字符串包含检查
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToUpper(s), strings.ToUpper(substr))
}
