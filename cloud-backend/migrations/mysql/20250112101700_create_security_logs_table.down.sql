-- +migrate Down
-- 回滚迁移: 安全日志表
-- 版本: 20250112101700
-- 描述: 删除安全日志表及相关触发器、存储过程和视图
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:17:00

-- ============================================================================
-- 回滚安全日志表创建
-- ============================================================================

-- 删除视图
DROP VIEW IF EXISTS unresolved_security_incidents;
DROP VIEW IF EXISTS high_risk_ip_monitor;
DROP VIEW IF EXISTS security_threat_overview;

-- 删除存储过程
DROP PROCEDURE IF EXISTS CleanupSecurityLogs;
DROP PROCEDURE IF EXISTS DetectAnomalousActivity;
DROP PROCEDURE IF EXISTS GetSecurityThreatStats;
DROP PROCEDURE IF EXISTS LogSecurityEvent;

-- 删除触发器
DROP TRIGGER IF EXISTS security_logs_aggregation;
DROP TRIGGER IF EXISTS security_logs_alert_trigger;
DROP TRIGGER IF EXISTS security_logs_auto_analysis;

-- 删除安全日志表（会自动删除所有约束、索引、外键和分区）
DROP TABLE IF EXISTS security_logs;

-- 注意：由于使用了 IF EXISTS，即使对象不存在也不会报错
-- 这确保了回滚操作的安全性和幂等性