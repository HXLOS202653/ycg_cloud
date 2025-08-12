-- +migrate Down
-- 回滚迁移: 文件版本表
-- 版本: 20250112100400
-- 描述: 删除文件版本历史管理表及相关触发器、存储过程和视图
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:04:00

-- ============================================================================
-- 回滚文件版本表创建
-- ============================================================================

-- 删除视图
DROP VIEW IF EXISTS version_activity;
DROP VIEW IF EXISTS file_version_stats;

-- 删除存储过程
DROP PROCEDURE IF EXISTS GetFileVersionHistory;
DROP PROCEDURE IF EXISTS CleanupExpiredVersions;

-- 删除触发器
DROP TRIGGER IF EXISTS file_versions_cleanup;
DROP TRIGGER IF EXISTS file_versions_update_files;
DROP TRIGGER IF EXISTS file_versions_auto_version;

-- 删除文件版本表（会自动删除所有约束、索引和外键）
DROP TABLE IF EXISTS file_versions;

-- 注意：由于使用了 IF EXISTS，即使对象不存在也不会报错
-- 这确保了回滚操作的安全性和幂等性