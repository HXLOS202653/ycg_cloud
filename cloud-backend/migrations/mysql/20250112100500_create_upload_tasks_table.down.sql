-- +migrate Down
-- 回滚迁移: 文件上传任务表
-- 版本: 20250112100500
-- 描述: 删除文件上传任务管理表及相关触发器、存储过程和视图
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:05:00

-- ============================================================================
-- 回滚上传任务表创建
-- ============================================================================

-- 删除视图
DROP VIEW IF EXISTS upload_performance;
DROP VIEW IF EXISTS upload_task_stats;

-- 删除存储过程
DROP PROCEDURE IF EXISTS GetUserUploadStats;
DROP PROCEDURE IF EXISTS CleanupExpiredUploads;

-- 删除触发器
DROP TRIGGER IF EXISTS upload_tasks_expiry_check;
DROP TRIGGER IF EXISTS upload_tasks_progress_update;

-- 删除上传任务表（会自动删除所有约束、索引和外键）
DROP TABLE IF EXISTS upload_tasks;

-- 注意：由于使用了 IF EXISTS，即使对象不存在也不会报错
-- 这确保了回滚操作的安全性和幂等性