-- +migrate Down
-- 回滚迁移: 文件权限表
-- 版本: 20250112101100
-- 描述: 删除文件权限表及相关触发器、存储过程和视图
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:11:00

-- ============================================================================
-- 回滚文件权限表创建
-- ============================================================================

-- 删除视图
DROP VIEW IF EXISTS expiring_permissions_alert;
DROP VIEW IF EXISTS file_permissions_detail;
DROP VIEW IF EXISTS user_permissions_overview;

-- 删除存储过程
DROP PROCEDURE IF EXISTS CleanExpiredPermissions;
DROP PROCEDURE IF EXISTS BatchGrantPermissions;
DROP PROCEDURE IF EXISTS CheckUserFilePermission;

-- 删除触发器
DROP TRIGGER IF EXISTS file_permissions_conflict_check;
DROP TRIGGER IF EXISTS file_permissions_inheritance_update;
DROP TRIGGER IF EXISTS file_permissions_activation_check;

-- 删除文件权限表（会自动删除所有约束、索引和外键）
DROP TABLE IF EXISTS file_permissions;

-- 注意：由于使用了 IF EXISTS，即使对象不存在也不会报错
-- 这确保了回滚操作的安全性和幂等性