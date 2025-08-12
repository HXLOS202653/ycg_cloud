-- ============================================================================
-- 种子数据管理脚本
-- 文件: seed_data_manager.sql
-- 描述: 生产环境安全的种子数据管理工具
-- ============================================================================

-- 设置安全模式
SET sql_mode = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';

-- ============================================================================
-- 存储过程：安全创建管理员用户
-- ============================================================================

DELIMITER //
CREATE PROCEDURE CreateAdminUser(
    IN p_username VARCHAR(50),
    IN p_email VARCHAR(100),
    IN p_real_name VARCHAR(100),
    IN p_password_hash VARCHAR(255),
    IN p_role ENUM('super_admin', 'admin', 'user', 'system'),
    IN p_storage_quota BIGINT
)
BEGIN
    DECLARE user_exists INT DEFAULT 0;
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;
    
    START TRANSACTION;
    
    -- 检查用户是否已存在
    SELECT COUNT(*) INTO user_exists 
    FROM users 
    WHERE username = p_username OR email = p_email;
    
    IF user_exists > 0 THEN
        SIGNAL SQLSTATE '45000' 
        SET MESSAGE_TEXT = '用户名或邮箱已存在';
    END IF;
    
    -- 创建用户
    INSERT INTO users (
        username, email, email_verified, password_hash, real_name, 
        role, status, storage_quota, upload_bandwidth_limit, 
        download_bandwidth_limit, max_file_size, timezone, language,
        created_at, updated_at
    ) VALUES (
        p_username, p_email, TRUE, p_password_hash, p_real_name,
        p_role, 'active', p_storage_quota, 104857600,
        104857600, 10737418240, 'Asia/Shanghai', 'zh-CN',
        CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
    );
    
    -- 获取新创建的用户ID
    SET @new_user_id = LAST_INSERT_ID();
    
    -- 创建用户设置
    INSERT INTO user_settings (
        user_id, theme, language, timezone, items_per_page,
        sync_enabled, sync_bandwidth_limit, auto_preview,
        trash_retention_days, share_default_expire, enable_notifications,
        notification_email, notification_browser, notification_mobile,
        auto_logout_minutes, enable_two_factor, cache_size_mb,
        concurrent_uploads, concurrent_downloads, created_at, updated_at
    ) VALUES (
        @new_user_id, 'light', 'zh-CN', 'Asia/Shanghai', 30,
        TRUE, 10485760, TRUE, 30, 7, TRUE,
        TRUE, TRUE, FALSE, 60, FALSE, 256,
        3, 5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
    );
    
    -- 创建根目录
    INSERT INTO files (
        filename, file_path, file_type, is_folder, user_id,
        parent_id, file_size, storage_location, created_at, updated_at
    ) VALUES (
        '根目录', '/', 'folder', TRUE, @new_user_id,
        NULL, 0, 'local', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
    );
    
    SET @root_folder_id = LAST_INSERT_ID();
    
    -- 创建文件夹扩展属性
    INSERT INTO folders (
        id, folder_type, folder_category, total_size, file_count,
        folder_count, total_file_count, total_folder_count,
        default_permission, max_file_size, max_file_versions,
        is_template, collaboration_enabled, discussion_enabled,
        task_management, search_weight, cache_ttl, stats_updated_at
    ) VALUES (
        @root_folder_id, 'root', 'personal', 0, 0,
        0, 0, 0, 'private', p_storage_quota, 50,
        FALSE, TRUE, TRUE, TRUE, 1.0, 3600, CURRENT_TIMESTAMP
    );
    
    COMMIT;
    
    SELECT 
        CONCAT('用户 ', p_username, ' 创建成功') as message,
        @new_user_id as user_id,
        @root_folder_id as root_folder_id;
        
END//
DELIMITER ;

-- ============================================================================
-- 存储过程：创建系统配置
-- ============================================================================

DELIMITER //
CREATE PROCEDURE CreateSystemConfig(
    IN p_config_key VARCHAR(100),
    IN p_config_value TEXT,
    IN p_value_type ENUM('string','integer','float','boolean','json','text'),
    IN p_category VARCHAR(50),
    IN p_subcategory VARCHAR(50),
    IN p_display_name VARCHAR(200),
    IN p_description TEXT,
    IN p_is_public BOOLEAN,
    IN p_is_editable BOOLEAN,
    IN p_is_sensitive BOOLEAN,
    IN p_access_level ENUM('public','user','admin','super_admin'),
    IN p_created_by INT
)
BEGIN
    DECLARE config_exists INT DEFAULT 0;
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;
    
    START TRANSACTION;
    
    -- 检查配置是否已存在
    SELECT COUNT(*) INTO config_exists 
    FROM system_configs 
    WHERE config_key = p_config_key;
    
    IF config_exists > 0 THEN
        SIGNAL SQLSTATE '45000' 
        SET MESSAGE_TEXT = '配置键已存在';
    END IF;
    
    -- 创建配置
    INSERT INTO system_configs (
        config_key, config_value, value_type, category, subcategory,
        config_group, display_name, description, is_public, is_editable,
        is_system, is_sensitive, access_level, is_active, environment,
        config_source, created_by, updated_by, created_at, updated_at
    ) VALUES (
        p_config_key, p_config_value, p_value_type, p_category, p_subcategory,
        'custom', p_display_name, p_description, p_is_public, p_is_editable,
        FALSE, p_is_sensitive, p_access_level, TRUE, 'production',
        'manual', p_created_by, p_created_by, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
    );
    
    COMMIT;
    
    SELECT 
        CONCAT('配置 ', p_config_key, ' 创建成功') as message,
        LAST_INSERT_ID() as config_id;
        
END//
DELIMITER ;

-- ============================================================================
-- 存储过程：生成随机密码哈希
-- ============================================================================

DELIMITER //
CREATE PROCEDURE GeneratePasswordHash(
    IN p_password VARCHAR(255)
)
BEGIN
    -- 注意：这只是示例，实际应该使用应用程序层的bcrypt
    -- 生产环境中应该在应用程序中生成哈希值
    SELECT 
        CONCAT('$2a$12$', SHA2(CONCAT(p_password, UNIX_TIMESTAMP(), RAND()), 256)) as password_hash,
        '警告：请在应用程序层使用bcrypt生成真实的哈希值' as warning_message;
END//
DELIMITER ;

-- ============================================================================
-- 存储过程：检查种子数据状态
-- ============================================================================

DELIMITER //
CREATE PROCEDURE CheckSeedDataStatus()
BEGIN
    SELECT 
        '种子数据状态检查' as title,
        (SELECT COUNT(*) FROM users WHERE role IN ('super_admin', 'admin')) as admin_users_count,
        (SELECT COUNT(*) FROM system_configs WHERE is_system = FALSE) as custom_configs_count,
        (SELECT COUNT(*) FROM teams WHERE team_type = 'organization') as organization_teams_count,
        (SELECT COUNT(*) FROM files WHERE is_folder = TRUE AND parent_id IS NULL) as root_folders_count,
        (SELECT COUNT(*) FROM notifications WHERE notification_category = 'welcome') as welcome_notifications_count;
        
    -- 检查关键配置
    SELECT 
        '关键系统配置状态' as config_check,
        config_key,
        config_value,
        is_active,
        access_level
    FROM system_configs 
    WHERE config_key IN (
        'system.app_name',
        'system.maintenance_mode',
        'security.password_min_length',
        'security.session_timeout',
        'file.max_file_size'
    )
    ORDER BY config_key;
    
    -- 检查管理员账户
    SELECT 
        '管理员账户状态' as admin_check,
        id,
        username,
        email,
        role,
        status,
        storage_quota,
        created_at
    FROM users 
    WHERE role IN ('super_admin', 'admin')
    ORDER BY role DESC, created_at;
END//
DELIMITER ;

-- ============================================================================
-- 存储过程：安全清理测试数据
-- ============================================================================

DELIMITER //
CREATE PROCEDURE CleanTestData()
BEGIN
    DECLARE done INT DEFAULT FALSE;
    DECLARE test_user_id INT;
    DECLARE test_users_cursor CURSOR FOR 
        SELECT id FROM users WHERE username LIKE 'test_%' OR email LIKE 'test%@%';
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;
    
    START TRANSACTION;
    
    -- 清理测试用户相关数据
    OPEN test_users_cursor;
    read_loop: LOOP
        FETCH test_users_cursor INTO test_user_id;
        IF done THEN
            LEAVE read_loop;
        END IF;
        
        -- 清理用户相关的所有数据
        DELETE FROM user_settings WHERE user_id = test_user_id;
        DELETE FROM user_sessions WHERE user_id = test_user_id;
        DELETE FROM notifications WHERE user_id = test_user_id OR sender_id = test_user_id;
        DELETE FROM team_members WHERE user_id = test_user_id;
        DELETE FROM team_invitations WHERE inviter_id = test_user_id OR invitee_id = test_user_id;
        DELETE FROM file_permissions WHERE user_id = test_user_id OR granted_by = test_user_id;
        DELETE FROM file_shares WHERE user_id = test_user_id;
        DELETE FROM upload_tasks WHERE user_id = test_user_id;
        DELETE FROM operation_logs WHERE user_id = test_user_id;
        DELETE FROM security_logs WHERE user_id = test_user_id;
        
        -- 清理用户文件（级联删除会处理相关表）
        DELETE FROM files WHERE user_id = test_user_id;
        
        -- 最后删除用户
        DELETE FROM users WHERE id = test_user_id;
        
    END LOOP;
    CLOSE test_users_cursor;
    
    -- 清理测试团队
    DELETE FROM teams WHERE name LIKE '测试%' OR name LIKE 'Test%';
    
    -- 清理测试配置
    DELETE FROM system_configs WHERE config_key LIKE 'test.%' OR display_name LIKE '测试%';
    
    COMMIT;
    
    SELECT '测试数据清理完成' as message;
END//
DELIMITER ;

-- ============================================================================
-- 使用示例和说明
-- ============================================================================

/*
使用示例：

1. 创建管理员用户：
CALL CreateAdminUser(
    'newadmin', 
    'newadmin@company.com', 
    '新管理员', 
    '$2a$12$hashedpassword...', 
    'admin', 
    107374182400
);

2. 创建系统配置：
CALL CreateSystemConfig(
    'custom.feature_flag',
    'true',
    'boolean',
    'feature',
    'flags',
    '自定义功能开关',
    '控制自定义功能的启用状态',
    FALSE,
    TRUE,
    FALSE,
    'admin',
    1
);

3. 生成密码哈希（仅供参考）：
CALL GeneratePasswordHash('mypassword123');

4. 检查种子数据状态：
CALL CheckSeedDataStatus();

5. 清理测试数据：
CALL CleanTestData();

注意事项：
- 生产环境中请使用强密码
- 密码哈希应在应用程序层生成
- 定期检查和更新系统配置
- 谨慎执行清理操作
*/