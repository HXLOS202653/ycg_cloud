-- +migrate Down
-- 回滚迁移: 文件夹扩展属性表
-- 版本: 20250112100700
-- 描述: 删除文件夹扩展属性表及相关触发器、存储过程和视图
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:07:00

-- ============================================================================
-- 回滚文件夹扩展属性表创建
-- ============================================================================

-- 删除视图
DROP VIEW IF EXISTS collaborative_folders;
DROP VIEW IF EXISTS folder_size_ranking;
DROP VIEW IF EXISTS folder_detailed_stats;

-- 删除存储过程
DROP PROCEDURE IF EXISTS InitializeFolderRecords;
DROP PROCEDURE IF EXISTS RecalculateFolderStats;
DROP PROCEDURE IF EXISTS UpdateParentFolderStats;

-- 删除触发器
DROP TRIGGER IF EXISTS folders_permission_inherit;
DROP TRIGGER IF EXISTS folders_stats_delete;
DROP TRIGGER IF EXISTS folders_stats_update;

-- 删除文件夹扩展属性表（会自动删除所有约束、索引和外键）
DROP TABLE IF EXISTS folders;

-- 注意：由于使用了 IF EXISTS，即使对象不存在也不会报错
-- 这确保了回滚操作的安全性和幂等性