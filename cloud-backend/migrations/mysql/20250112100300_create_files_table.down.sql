-- +migrate Down
-- 回滚迁移: 文件表
-- 版本: 20250112100300
-- 描述: 删除文件和文件夹管理表及相关触发器和视图
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:03:00

-- ============================================================================
-- 回滚文件表创建
-- ============================================================================

-- 删除视图
DROP VIEW IF EXISTS popular_files;
DROP VIEW IF EXISTS file_stats;

-- 删除触发器
DROP TRIGGER IF EXISTS files_stats_delete;
DROP TRIGGER IF EXISTS files_stats_update;

-- 删除文件表（会自动删除所有约束、索引和外键）
DROP TABLE IF EXISTS files;

-- 注意：由于使用了 IF EXISTS，即使对象不存在也不会报错
-- 这确保了回滚操作的安全性和幂等性