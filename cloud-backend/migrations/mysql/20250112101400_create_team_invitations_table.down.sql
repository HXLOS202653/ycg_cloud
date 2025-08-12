-- +migrate Down
-- 回滚迁移: 团队邀请表
-- 版本: 20250112101400
-- 描述: 删除团队邀请表及相关触发器、存储过程和视图
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:14:00

-- ============================================================================
-- 回滚团队邀请表创建
-- ============================================================================

-- 删除视图
DROP VIEW IF EXISTS expiring_invitations_alert;
DROP VIEW IF EXISTS invitation_conversion_stats;
DROP VIEW IF EXISTS active_invitations_overview;

-- 删除存储过程
DROP PROCEDURE IF EXISTS GetInvitationStatistics;
DROP PROCEDURE IF EXISTS CleanExpiredInvitations;
DROP PROCEDURE IF EXISTS CreateTeamInvitation;

-- 删除触发器
DROP TRIGGER IF EXISTS team_invitations_expiry_check;
DROP TRIGGER IF EXISTS team_invitations_status_update;
DROP TRIGGER IF EXISTS team_invitations_initialize;

-- 删除团队邀请表（会自动删除所有约束、索引和外键）
DROP TABLE IF EXISTS team_invitations;

-- 注意：由于使用了 IF EXISTS，即使对象不存在也不会报错
-- 这确保了回滚操作的安全性和幂等性