-- ============================================================================
-- 种子数据回滚脚本 - 清理系统基础数据
-- 文件: 20250112120000_create_seed_data.down.sql
-- 描述: 回滚种子数据，清理创建的管理员用户、配置等数据
-- ============================================================================

-- 设置字符集和排序规则
SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 禁用外键约束检查，方便删除操作
SET FOREIGN_KEY_CHECKS = 0;

-- ============================================================================
-- 清理种子数据（按依赖关系逆序删除）
-- ============================================================================

-- 记录回滚操作日志
INSERT INTO operation_logs (
    user_id,
    session_id,
    operation_type,
    operation_category,
    resource_type,
    resource_id,
    action,
    result,
    ip_address,
    user_agent,
    description,
    correlation_id,
    created_at
) VALUES (
    1,
    'seed-data-rollback',
    'data_management',
    'rollback',
    'database',
    'seed_data',
    'delete',
    'success',
    '127.0.0.1',
    'MySQL/Migration',
    '回滚系统种子数据：清理管理员用户、默认配置、示例团队和文件夹',
    'seed-rollback-' || UNIX_TIMESTAMP(),
    CURRENT_TIMESTAMP
);

-- 1. 清理通知数据
DELETE FROM notifications WHERE id IN (1, 2, 3);

-- 2. 清理文件夹扩展属性
DELETE FROM folders WHERE id IN (1, 2, 3, 4, 5, 6, 7, 8, 9, 10);

-- 3. 清理文件数据（包括文件夹）
DELETE FROM files WHERE id IN (1, 2, 3, 4, 5, 6, 7, 8, 9, 10);

-- 4. 清理团队成员关系
DELETE FROM team_members WHERE team_id = 1;

-- 5. 清理团队数据
DELETE FROM teams WHERE id = 1;

-- 6. 清理系统配置数据
DELETE FROM system_configs WHERE id BETWEEN 1 AND 21;

-- 7. 清理用户设置数据
DELETE FROM user_settings WHERE user_id IN (1, 2, 3, 4, 5);

-- 8. 清理用户数据
DELETE FROM users WHERE id IN (1, 2, 3, 4, 5);

-- 9. 清理种子数据相关的操作日志
DELETE FROM operation_logs WHERE correlation_id LIKE 'seed-%';

-- ============================================================================
-- 重置自增ID
-- ============================================================================

-- 重置各表的自增ID（可选，根据需要）
-- ALTER TABLE users AUTO_INCREMENT = 1;
-- ALTER TABLE teams AUTO_INCREMENT = 1;
-- ALTER TABLE system_configs AUTO_INCREMENT = 1;
-- ALTER TABLE files AUTO_INCREMENT = 1;
-- ALTER TABLE notifications AUTO_INCREMENT = 1;

-- 重新启用外键约束检查
SET FOREIGN_KEY_CHECKS = 1;

-- ============================================================================
-- 验证清理结果
-- ============================================================================

-- 显示清理结果
SELECT 
    '种子数据回滚完成' as message,
    (SELECT COUNT(*) FROM users WHERE id IN (1,2,3,4,5)) as remaining_users,
    (SELECT COUNT(*) FROM teams WHERE id = 1) as remaining_teams,
    (SELECT COUNT(*) FROM system_configs WHERE id BETWEEN 1 AND 21) as remaining_configs,
    (SELECT COUNT(*) FROM files WHERE id BETWEEN 1 AND 10) as remaining_files,
    (SELECT COUNT(*) FROM notifications WHERE id BETWEEN 1 AND 3) as remaining_notifications;

-- 显示警告信息
SELECT 
    '警告：种子数据已完全清理' as warning_message,
    '请确保在生产环境中谨慎执行此操作' as production_warning;