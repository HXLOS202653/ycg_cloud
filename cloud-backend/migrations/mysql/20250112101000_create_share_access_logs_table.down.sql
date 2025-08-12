-- +migrate Down
-- 回滚迁移: 分享访问日志表
-- 版本: 20250112101000
-- 描述: 删除分享访问日志表及相关触发器、存储过程和视图
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:10:00

-- ============================================================================
-- 回滚分享访问日志表创建
-- ============================================================================

-- 删除视图
DROP VIEW IF EXISTS suspicious_access_monitor;
DROP VIEW IF EXISTS access_device_stats;
DROP VIEW IF EXISTS popular_share_access_stats;

-- 删除存储过程
DROP PROCEDURE IF EXISTS CleanOldAccessLogs;
DROP PROCEDURE IF EXISTS DetectSecurityThreats;
DROP PROCEDURE IF EXISTS GetAccessGeographyStats;
DROP PROCEDURE IF EXISTS GetShareAccessStatistics;

-- 删除触发器
DROP TRIGGER IF EXISTS share_access_logs_security_check;
DROP TRIGGER IF EXISTS share_access_logs_update_stats;

-- 删除分享访问日志表（会自动删除所有约束、索引、外键和分区）
DROP TABLE IF EXISTS share_access_logs;

-- 注意：由于使用了 IF EXISTS，即使对象不存在也不会报错
-- 这确保了回滚操作的安全性和幂等性