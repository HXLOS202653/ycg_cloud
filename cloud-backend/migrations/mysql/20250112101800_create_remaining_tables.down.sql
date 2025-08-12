-- +migrate Down
-- 回滚迁移: 剩余核心表批量删除
-- 版本: 20250112101800
-- 描述: 批量删除聊天系统和搜索系统表
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:18:00

-- ============================================================================
-- 回滚剩余表创建
-- ============================================================================

-- 删除视图
DROP VIEW IF EXISTS search_statistics;
DROP VIEW IF EXISTS chat_rooms_overview;

-- 删除表（按依赖顺序）
DROP TABLE IF EXISTS search_logs;
DROP TABLE IF EXISTS file_contents;
DROP TABLE IF EXISTS chat_messages;
DROP TABLE IF EXISTS chat_room_members;
DROP TABLE IF EXISTS chat_rooms;

-- 注意：由于使用了 IF EXISTS，即使对象不存在也不会报错
-- 这确保了回滚操作的安全性和幂等性