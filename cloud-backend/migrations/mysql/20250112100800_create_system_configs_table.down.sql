-- +migrate Down
-- 回滚迁移: 系统配置表
-- 版本: 20250112100800
-- 描述: 删除系统配置表及相关触发器、存储过程和视图
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:08:00

-- ============================================================================
-- 回滚系统配置表创建
-- ============================================================================

-- 删除视图
DROP VIEW IF EXISTS recently_modified_configs;
DROP VIEW IF EXISTS config_category_stats;
DROP VIEW IF EXISTS public_system_configs;

-- 删除存储过程
DROP PROCEDURE IF EXISTS ValidateConfigValue;
DROP PROCEDURE IF EXISTS ImportSystemConfigs;

-- 删除触发器
DROP TRIGGER IF EXISTS system_configs_cache_invalidate;
DROP TRIGGER IF EXISTS system_configs_audit_update;

-- 删除辅助表
DROP TABLE IF EXISTS cache_invalidation_queue;
DROP TABLE IF EXISTS system_logs;

-- 删除系统配置表（会自动删除所有约束、索引和外键）
DROP TABLE IF EXISTS system_configs;

-- 注意：由于使用了 IF EXISTS，即使对象不存在也不会报错
-- 这确保了回滚操作的安全性和幂等性