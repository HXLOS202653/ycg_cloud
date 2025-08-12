-- +migrate Down
-- 回滚迁移: 用户会话表
-- 版本: 20250112100100
-- 描述: 删除用户会话管理表及相关触发器、存储过程和视图
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:01:00

-- ============================================================================
-- 回滚用户会话表创建
-- ============================================================================

-- 删除视图
DROP VIEW IF EXISTS user_session_stats;

-- 删除存储过程
DROP PROCEDURE IF EXISTS CleanExpiredSessions;

-- 删除触发器
DROP TRIGGER IF EXISTS user_sessions_expires_check;
DROP TRIGGER IF EXISTS user_sessions_update_check;

-- 删除用户会话表（会自动删除所有约束、索引和外键）
DROP TABLE IF EXISTS user_sessions;

-- 注意：由于使用了 IF EXISTS，即使对象不存在也不会报错
-- 这确保了回滚操作的安全性和幂等性