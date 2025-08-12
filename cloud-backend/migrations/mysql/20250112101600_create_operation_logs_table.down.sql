-- +migrate Down
-- 回滚迁移: 操作日志表
-- 版本: 20250112101600
-- 描述: 删除操作日志表及相关触发器、存储过程和视图
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:16:00

-- ============================================================================
-- 回滚操作日志表创建
-- ============================================================================

-- 删除视图
DROP VIEW IF EXISTS system_performance_monitor;
DROP VIEW IF EXISTS suspicious_operations_monitor;
DROP VIEW IF EXISTS operation_stats_overview;

-- 删除存储过程
DROP PROCEDURE IF EXISTS CleanupOperationLogs;
DROP PROCEDURE IF EXISTS GetOperationTrends;
DROP PROCEDURE IF EXISTS GetUserOperationStats;
DROP PROCEDURE IF EXISTS LogOperation;

-- 删除触发器
DROP TRIGGER IF EXISTS operation_logs_stats_update;
DROP TRIGGER IF EXISTS operation_logs_risk_assessment;

-- 删除操作日志表（会自动删除所有约束、索引、外键和分区）
DROP TABLE IF EXISTS operation_logs;

-- 注意：由于使用了 IF EXISTS，即使对象不存在也不会报错
-- 这确保了回滚操作的安全性和幂等性