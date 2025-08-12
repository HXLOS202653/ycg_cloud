-- +migrate Down
-- 回滚迁移: 通知表
-- 版本: 20250112101500
-- 描述: 删除通知表及相关触发器、存储过程和视图
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:15:00

-- ============================================================================
-- 回滚通知表创建
-- ============================================================================

-- 删除视图
DROP VIEW IF EXISTS urgent_unread_notifications;
DROP VIEW IF EXISTS notification_type_stats;
DROP VIEW IF EXISTS user_notification_stats;

-- 删除存储过程
DROP PROCEDURE IF EXISTS CleanExpiredNotifications;
DROP PROCEDURE IF EXISTS MarkNotificationsAsRead;
DROP PROCEDURE IF EXISTS CreateNotification;

-- 删除触发器
DROP TRIGGER IF EXISTS notifications_cleanup;
DROP TRIGGER IF EXISTS notifications_initialize;
DROP TRIGGER IF EXISTS notifications_status_update;

-- 删除通知表（会自动删除所有约束、索引、外键和分区）
DROP TABLE IF EXISTS notifications;

-- 注意：由于使用了 IF EXISTS，即使对象不存在也不会报错
-- 这确保了回滚操作的安全性和幂等性