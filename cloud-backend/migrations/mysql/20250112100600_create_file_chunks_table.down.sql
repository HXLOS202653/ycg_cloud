-- +migrate Down
-- 回滚迁移: 文件分片表
-- 版本: 20250112100600
-- 描述: 删除文件分片信息管理表及相关触发器、存储过程和视图
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:06:00

-- ============================================================================
-- 回滚文件分片表创建
-- ============================================================================

-- 删除视图
DROP VIEW IF EXISTS chunk_performance;
DROP VIEW IF EXISTS chunk_status_stats;

-- 删除存储过程
DROP PROCEDURE IF EXISTS ResetFailedChunks;
DROP PROCEDURE IF EXISTS GetUploadChunkStatus;
DROP PROCEDURE IF EXISTS CleanupExpiredChunkLocks;

-- 删除触发器
DROP TRIGGER IF EXISTS file_chunks_update_task;
DROP TRIGGER IF EXISTS file_chunks_lock_update;
DROP TRIGGER IF EXISTS file_chunks_status_update;

-- 删除文件分片表（会自动删除所有约束、索引和外键）
DROP TABLE IF EXISTS file_chunks;

-- 注意：由于使用了 IF EXISTS，即使对象不存在也不会报错
-- 这确保了回滚操作的安全性和幂等性