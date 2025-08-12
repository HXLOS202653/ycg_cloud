-- +migrate Down
-- 回滚迁移: 团队成员表
-- 版本: 20250112101300
-- 描述: 删除团队成员表及相关触发器、存储过程和视图
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:13:00

-- ============================================================================
-- 回滚团队成员表创建
-- ============================================================================

-- 删除视图
DROP VIEW IF EXISTS member_activity_ranking;
DROP VIEW IF EXISTS team_role_statistics;
DROP VIEW IF EXISTS team_members_detail;

-- 删除存储过程
DROP PROCEDURE IF EXISTS CalculateMemberActivityScore;
DROP PROCEDURE IF EXISTS BulkImportTeamMembers;
DROP PROCEDURE IF EXISTS AddTeamMember;

-- 删除触发器
DROP TRIGGER IF EXISTS team_members_cleanup;
DROP TRIGGER IF EXISTS team_members_status_update;
DROP TRIGGER IF EXISTS team_members_initialize;

-- 删除团队成员表（会自动删除所有约束、索引和外键）
DROP TABLE IF EXISTS team_members;

-- 注意：由于使用了 IF EXISTS，即使对象不存在也不会报错
-- 这确保了回滚操作的安全性和幂等性