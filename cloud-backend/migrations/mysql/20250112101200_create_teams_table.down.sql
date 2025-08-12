-- +migrate Down
-- 回滚迁移: 团队表
-- 版本: 20250112101200
-- 描述: 删除团队表及相关触发器、存储过程和视图
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:12:00

-- ============================================================================
-- 回滚团队表创建
-- ============================================================================

-- 删除视图
DROP VIEW IF EXISTS teams_activity_ranking;
DROP VIEW IF EXISTS teams_storage_monitor;
DROP VIEW IF EXISTS teams_overview_stats;

-- 删除存储过程
DROP PROCEDURE IF EXISTS UpdateTeamActivityStats;
DROP PROCEDURE IF EXISTS CheckTeamQuotas;
DROP PROCEDURE IF EXISTS GetTeamStatistics;

-- 删除触发器
DROP TRIGGER IF EXISTS teams_deletion_check;
DROP TRIGGER IF EXISTS teams_update_stats;
DROP TRIGGER IF EXISTS teams_initialize_settings;

-- 删除团队表（会自动删除所有约束、索引和外键）
DROP TABLE IF EXISTS teams;

-- 注意：由于使用了 IF EXISTS，即使对象不存在也不会报错
-- 这确保了回滚操作的安全性和幂等性