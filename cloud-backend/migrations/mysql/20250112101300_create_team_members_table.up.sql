-- +migrate Up
-- 创建迁移: 团队成员表
-- 版本: 20250112101300
-- 描述: 创建团队成员表，管理团队成员关系和权限
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:13:00
-- 依赖: 20250112101200_create_teams_table, 20250112100000_create_users_table
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1
-- 说明: 团队成员表管理用户与团队的关系，包括角色、权限和状态

-- ============================================================================
-- 团队成员表 (team_members)
-- ============================================================================

CREATE TABLE team_members (
    -- 基础标识字段
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '团队成员关系唯一标识',
    
    -- 关联关系（复合唯一约束）
    team_id BIGINT UNSIGNED NOT NULL COMMENT '所属团队ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '成员用户ID',
    
    -- 成员角色和权限
    role ENUM('owner', 'admin', 'member', 'viewer', 'guest', 'custom') DEFAULT 'member' COMMENT '成员角色级别',
    role_name VARCHAR(50) DEFAULT NULL COMMENT '自定义角色名称（当role为custom时）',
    permissions JSON DEFAULT NULL COMMENT '详细权限配置，JSON格式',
    
    -- 成员状态和生命周期
    status ENUM('active', 'inactive', 'pending', 'suspended', 'left', 'removed') DEFAULT 'pending' COMMENT '成员状态',
    
    -- 邀请和加入信息
    invited_by BIGINT UNSIGNED DEFAULT NULL COMMENT '邀请者用户ID',
    invitation_method ENUM('email', 'link', 'direct', 'bulk_import', 'api') DEFAULT 'email' COMMENT '邀请方式',
    invitation_token VARCHAR(128) DEFAULT NULL COMMENT '邀请令牌',
    
    -- 时间管理
    invited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '邀请时间',
    joined_at TIMESTAMP NULL DEFAULT NULL COMMENT '实际加入时间（接受邀请时间）',
    last_active_at TIMESTAMP NULL DEFAULT NULL COMMENT '最后活跃时间',
    
    -- 成员个性化设置
    display_name VARCHAR(100) DEFAULT NULL COMMENT '在团队中的显示名称',
    bio TEXT DEFAULT NULL COMMENT '成员简介或职位描述',
    department VARCHAR(100) DEFAULT NULL COMMENT '所属部门',
    job_title VARCHAR(100) DEFAULT NULL COMMENT '职位头衔',
    
    -- 通知和偏好设置
    notification_settings JSON DEFAULT NULL COMMENT '通知偏好设置',
    team_email_enabled BOOLEAN DEFAULT TRUE COMMENT '是否接收团队邮件通知',
    mention_notifications BOOLEAN DEFAULT TRUE COMMENT '是否接收@提及通知',
    
    -- 访问控制和限制
    access_level ENUM('full', 'limited', 'readonly', 'restricted') DEFAULT 'full' COMMENT '访问级别',
    ip_restrictions JSON DEFAULT NULL COMMENT 'IP访问限制',
    time_restrictions JSON DEFAULT NULL COMMENT '时间访问限制',
    
    -- 配额和使用限制
    storage_quota BIGINT UNSIGNED DEFAULT NULL COMMENT '个人存储配额（字节），NULL表示无限制',
    storage_used BIGINT UNSIGNED DEFAULT 0 COMMENT '已使用存储空间',
    max_file_size BIGINT UNSIGNED DEFAULT NULL COMMENT '单文件最大大小限制',
    max_daily_uploads INT UNSIGNED DEFAULT NULL COMMENT '每日最大上传文件数',
    
    -- 团队贡献统计
    files_uploaded INT UNSIGNED DEFAULT 0 COMMENT '上传文件总数',
    files_shared INT UNSIGNED DEFAULT 0 COMMENT '分享文件总数',
    comments_made INT UNSIGNED DEFAULT 0 COMMENT '评论总数',
    collaborations_count INT UNSIGNED DEFAULT 0 COMMENT '协作项目数',
    
    -- 行为分析
    login_count INT UNSIGNED DEFAULT 0 COMMENT '登录团队次数',
    last_login_ip VARCHAR(45) DEFAULT NULL COMMENT '最后登录IP',
    activity_score DECIMAL(5,2) DEFAULT 0.00 COMMENT '活跃度评分（0-100）',
    
    -- 成员标签和分类
    tags JSON DEFAULT NULL COMMENT '成员标签，便于分组管理',
    groups JSON DEFAULT NULL COMMENT '所属工作组或项目组',
    
    -- 外部集成和同步
    external_id VARCHAR(100) DEFAULT NULL COMMENT '外部系统用户ID（如LDAP、AD等）',
    sync_source VARCHAR(50) DEFAULT NULL COMMENT '同步来源系统',
    last_synced_at TIMESTAMP NULL DEFAULT NULL COMMENT '最后同步时间',
    
    -- 离开和删除相关
    left_reason ENUM('voluntary', 'removed', 'inactive', 'violation', 'other') DEFAULT NULL COMMENT '离开原因',
    left_notes TEXT DEFAULT NULL COMMENT '离开备注说明',
    removed_by BIGINT UNSIGNED DEFAULT NULL COMMENT '移除操作者ID',
    
    -- 扩展属性
    metadata JSON DEFAULT NULL COMMENT '扩展元数据',
    custom_fields JSON DEFAULT NULL COMMENT '自定义字段',
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',
    
    -- 业务唯一约束
    UNIQUE KEY uk_team_user (team_id, user_id) COMMENT '团队用户唯一约束，防止重复加入',
    
    -- 业务索引设计
    INDEX idx_team_members_team_id (team_id) COMMENT '团队ID索引，查询团队成员',
    INDEX idx_team_members_user_id (user_id) COMMENT '用户ID索引，查询用户加入的团队',
    INDEX idx_team_members_role (role) COMMENT '角色索引，按角色筛选',
    INDEX idx_team_members_status (status) COMMENT '状态索引，筛选活跃成员',
    INDEX idx_team_members_invited_by (invited_by) COMMENT '邀请者索引，查询邀请记录',
    INDEX idx_team_members_joined_at (joined_at) COMMENT '加入时间索引，时间排序',
    INDEX idx_team_members_last_active_at (last_active_at) COMMENT '最后活跃时间索引',
    INDEX idx_team_members_access_level (access_level) COMMENT '访问级别索引',
    INDEX idx_team_members_department (department) COMMENT '部门索引，组织架构查询',
    INDEX idx_team_members_external_id (external_id) COMMENT '外部ID索引，系统集成',
    INDEX idx_team_members_deleted_at (deleted_at) COMMENT '软删除索引',
    
    -- 复合业务索引
    INDEX idx_team_members_team_status_role (team_id, status, role) COMMENT '团队状态角色复合索引',
    INDEX idx_team_members_team_active_members (team_id, status, last_active_at DESC) COMMENT '团队活跃成员查询优化',
    INDEX idx_team_members_user_teams_active (user_id, status, joined_at DESC) COMMENT '用户团队活跃状态查询',
    INDEX idx_team_members_invitation_tracking (invited_by, status, invited_at DESC) COMMENT '邀请跟踪复合索引',
    INDEX idx_team_members_activity_analysis (team_id, activity_score DESC, last_active_at) COMMENT '活跃度分析索引',
    INDEX idx_team_members_storage_usage (team_id, storage_used DESC, storage_quota) COMMENT '存储使用情况索引',
    
    -- 全文搜索索引
    FULLTEXT idx_fulltext_search (display_name, bio, job_title) COMMENT '成员信息全文搜索',
    
    -- 外键约束
    CONSTRAINT fk_team_members_team_id 
        FOREIGN KEY (team_id) REFERENCES teams(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_team_members_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_team_members_invited_by 
        FOREIGN KEY (invited_by) REFERENCES users(id) 
        ON DELETE SET NULL 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_team_members_removed_by 
        FOREIGN KEY (removed_by) REFERENCES users(id) 
        ON DELETE SET NULL 
        ON UPDATE CASCADE
        
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='团队成员关系表 - 管理用户与团队的关系、角色和权限'
  ROW_FORMAT=DYNAMIC;

-- ============================================================================
-- 团队成员表约束和检查
-- ============================================================================

-- 角色名称约束（自定义角色时必须提供角色名称）
ALTER TABLE team_members ADD CONSTRAINT chk_custom_role_name 
CHECK (
    (role != 'custom') OR 
    (role = 'custom' AND role_name IS NOT NULL AND LENGTH(TRIM(role_name)) >= 2)
);

-- 加入时间逻辑约束
ALTER TABLE team_members ADD CONSTRAINT chk_join_time_logic 
CHECK (
    (status = 'pending' AND joined_at IS NULL) OR
    (status != 'pending' AND joined_at IS NOT NULL AND joined_at >= invited_at)
);

-- 存储配额约束
ALTER TABLE team_members ADD CONSTRAINT chk_member_storage_quotas 
CHECK (
    (storage_quota IS NULL OR storage_quota > 0) AND
    storage_used >= 0 AND
    (storage_quota IS NULL OR storage_used <= storage_quota)
);

-- 统计字段约束
ALTER TABLE team_members ADD CONSTRAINT chk_member_statistics 
CHECK (
    files_uploaded >= 0 AND
    files_shared >= 0 AND
    comments_made >= 0 AND
    collaborations_count >= 0 AND
    login_count >= 0
);

-- 活跃度评分约束
ALTER TABLE team_members ADD CONSTRAINT chk_activity_score_range 
CHECK (activity_score >= 0.00 AND activity_score <= 100.00);

-- IP地址格式约束
ALTER TABLE team_members ADD CONSTRAINT chk_last_login_ip_format 
CHECK (
    last_login_ip IS NULL OR
    last_login_ip REGEXP '^([0-9]{1,3}\\.){3}[0-9]{1,3}$' OR  -- IPv4
    last_login_ip REGEXP '^[0-9a-fA-F:]+$'                    -- IPv6（简化）
);

-- JSON字段验证
ALTER TABLE team_members ADD CONSTRAINT chk_team_members_json_valid 
CHECK (
    (permissions IS NULL OR JSON_VALID(permissions)) AND
    (notification_settings IS NULL OR JSON_VALID(notification_settings)) AND
    (ip_restrictions IS NULL OR JSON_VALID(ip_restrictions)) AND
    (time_restrictions IS NULL OR JSON_VALID(time_restrictions)) AND
    (tags IS NULL OR JSON_VALID(tags)) AND
    (groups IS NULL OR JSON_VALID(groups)) AND
    (metadata IS NULL OR JSON_VALID(metadata)) AND
    (custom_fields IS NULL OR JSON_VALID(custom_fields))
);

-- ============================================================================
-- 团队成员管理触发器
-- ============================================================================

-- 成员加入初始化触发器
DELIMITER //
CREATE TRIGGER team_members_initialize
BEFORE INSERT ON team_members
FOR EACH ROW
BEGIN
    -- 生成邀请令牌（如果未提供）
    IF NEW.invitation_token IS NULL THEN
        SET NEW.invitation_token = SHA2(CONCAT(NEW.team_id, NEW.user_id, UNIX_TIMESTAMP(), RAND()), 256);
    END IF;
    
    -- 根据团队设置初始化成员权限
    IF NEW.permissions IS NULL THEN
        SET NEW.permissions = (
            SELECT JSON_EXTRACT(settings, '$.default_permissions')
            FROM teams 
            WHERE id = NEW.team_id
        );
    END IF;
    
    -- 初始化通知设置
    IF NEW.notification_settings IS NULL THEN
        SET NEW.notification_settings = JSON_OBJECT(
            'email_notifications', true,
            'push_notifications', true,
            'mention_notifications', true,
            'team_updates', true,
            'file_shared', true
        );
    END IF;
    
    -- 设置显示名称（如果未提供）
    IF NEW.display_name IS NULL THEN
        SET NEW.display_name = (
            SELECT COALESCE(real_name, username)
            FROM users 
            WHERE id = NEW.user_id
        );
    END IF;
END//
DELIMITER ;

-- 成员状态变更触发器
DELIMITER //
CREATE TRIGGER team_members_status_update
AFTER UPDATE ON team_members
FOR EACH ROW
BEGIN
    -- 当成员从pending变为active时，设置加入时间
    IF OLD.status = 'pending' AND NEW.status = 'active' AND NEW.joined_at IS NULL THEN
        UPDATE team_members 
        SET joined_at = CURRENT_TIMESTAMP 
        WHERE id = NEW.id;
    END IF;
    
    -- 更新团队成员计数
    IF OLD.status != NEW.status THEN
        -- 计算活跃成员数变化
        IF (OLD.status != 'active' AND NEW.status = 'active') THEN
            UPDATE teams 
            SET member_count = member_count + 1,
                last_activity_at = CURRENT_TIMESTAMP
            WHERE id = NEW.team_id;
        ELSEIF (OLD.status = 'active' AND NEW.status != 'active') THEN
            UPDATE teams 
            SET member_count = GREATEST(member_count - 1, 0),
                last_activity_at = CURRENT_TIMESTAMP
            WHERE id = NEW.team_id;
        END IF;
    END IF;
    
    -- 更新最后活跃时间
    IF NEW.status = 'active' AND (NEW.login_count > OLD.login_count OR NEW.files_uploaded > OLD.files_uploaded) THEN
        UPDATE team_members 
        SET last_active_at = CURRENT_TIMESTAMP 
        WHERE id = NEW.id;
    END IF;
END//
DELIMITER ;

-- 成员删除清理触发器
DELIMITER //
CREATE TRIGGER team_members_cleanup
AFTER DELETE ON team_members
FOR EACH ROW
BEGIN
    -- 更新团队成员计数
    UPDATE teams 
    SET member_count = GREATEST(member_count - 1, 0),
        last_activity_at = CURRENT_TIMESTAMP
    WHERE id = OLD.team_id;
    
    -- 如果删除的是所有者，需要特殊处理（这里简化处理）
    IF OLD.role = 'owner' THEN
        -- 实际应用中应该有复杂的所有权转移逻辑
        -- 这里只是标记需要处理
        UPDATE teams 
        SET metadata = JSON_SET(
            COALESCE(metadata, JSON_OBJECT()),
            '$.owner_removal_alert',
            CURRENT_TIMESTAMP
        )
        WHERE id = OLD.team_id;
    END IF;
END//
DELIMITER ;

-- ============================================================================
-- 团队成员管理存储过程
-- ============================================================================

-- 添加团队成员存储过程
DELIMITER //
CREATE PROCEDURE AddTeamMember(
    IN team_id_param BIGINT UNSIGNED,
    IN user_id_param BIGINT UNSIGNED,
    IN role_param VARCHAR(20),
    IN invited_by_param BIGINT UNSIGNED,
    OUT member_id BIGINT UNSIGNED,
    OUT success BOOLEAN,
    OUT error_message VARCHAR(500)
)
BEGIN
    DECLARE team_exists INT DEFAULT 0;
    DECLARE user_exists INT DEFAULT 0;
    DECLARE already_member INT DEFAULT 0;
    DECLARE team_member_limit INT DEFAULT 0;
    DECLARE current_member_count INT DEFAULT 0;
    
    SET success = FALSE;
    SET error_message = NULL;
    SET member_id = NULL;
    
    -- 检查团队是否存在且有效
    SELECT COUNT(*), max_members, member_count
    INTO team_exists, team_member_limit, current_member_count
    FROM teams
    WHERE id = team_id_param AND is_active = TRUE AND deleted_at IS NULL;
    
    IF team_exists = 0 THEN
        SET error_message = '团队不存在或已被禁用';
        LEAVE AddTeamMember;
    END IF;
    
    -- 检查用户是否存在
    SELECT COUNT(*) INTO user_exists
    FROM users
    WHERE id = user_id_param AND status = 'active';
    
    IF user_exists = 0 THEN
        SET error_message = '用户不存在或已被禁用';
        LEAVE AddTeamMember;
    END IF;
    
    -- 检查是否已经是成员
    SELECT COUNT(*) INTO already_member
    FROM team_members
    WHERE team_id = team_id_param AND user_id = user_id_param AND deleted_at IS NULL;
    
    IF already_member > 0 THEN
        SET error_message = '用户已经是团队成员';
        LEAVE AddTeamMember;
    END IF;
    
    -- 检查成员数量限制
    IF current_member_count >= team_member_limit THEN
        SET error_message = '团队成员数量已达上限';
        LEAVE AddTeamMember;
    END IF;
    
    -- 添加团队成员
    INSERT INTO team_members (
        team_id, user_id, role, invited_by, invitation_method, status
    ) VALUES (
        team_id_param, user_id_param, role_param, invited_by_param, 'direct', 'active'
    );
    
    SET member_id = LAST_INSERT_ID();
    SET success = TRUE;
    
END//
DELIMITER ;

-- 批量导入团队成员存储过程
DELIMITER //
CREATE PROCEDURE BulkImportTeamMembers(
    IN team_id_param BIGINT UNSIGNED,
    IN user_ids JSON,
    IN default_role VARCHAR(20),
    IN invited_by_param BIGINT UNSIGNED,
    OUT success_count INT,
    OUT failed_count INT,
    OUT error_details JSON
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE user_count INT DEFAULT 0;
    DECLARE current_user_id BIGINT UNSIGNED;
    DECLARE member_id BIGINT UNSIGNED;
    DECLARE operation_success BOOLEAN;
    DECLARE operation_error VARCHAR(500);
    DECLARE errors_array JSON DEFAULT JSON_ARRAY();
    
    SET success_count = 0;
    SET failed_count = 0;
    SET user_count = JSON_LENGTH(user_ids);
    
    -- 遍历用户ID列表
    WHILE i < user_count DO
        SET current_user_id = JSON_UNQUOTE(JSON_EXTRACT(user_ids, CONCAT('$[', i, ']')));
        
        -- 调用单个添加成员的存储过程
        CALL AddTeamMember(
            team_id_param, 
            current_user_id, 
            default_role, 
            invited_by_param,
            member_id,
            operation_success,
            operation_error
        );
        
        IF operation_success THEN
            SET success_count = success_count + 1;
        ELSE
            SET failed_count = failed_count + 1;
            SET errors_array = JSON_ARRAY_APPEND(
                errors_array, 
                '$', 
                JSON_OBJECT(
                    'user_id', current_user_id,
                    'error', operation_error
                )
            );
        END IF;
        
        SET i = i + 1;
    END WHILE;
    
    SET error_details = errors_array;
    
END//
DELIMITER ;

-- 计算成员活跃度评分存储过程
DELIMITER //
CREATE PROCEDURE CalculateMemberActivityScore(
    IN member_id_param BIGINT UNSIGNED
)
BEGIN
    DECLARE login_score DECIMAL(5,2) DEFAULT 0;
    DECLARE file_score DECIMAL(5,2) DEFAULT 0;
    DECLARE collaboration_score DECIMAL(5,2) DEFAULT 0;
    DECLARE recency_score DECIMAL(5,2) DEFAULT 0;
    DECLARE total_score DECIMAL(5,2) DEFAULT 0;
    DECLARE days_since_last_active INT DEFAULT 0;
    
    -- 获取成员活动数据
    SELECT 
        LEAST(login_count * 0.5, 20),  -- 登录频次评分（最多20分）
        LEAST(files_uploaded * 0.3, 25), -- 文件上传评分（最多25分）
        LEAST(collaborations_count * 2, 30), -- 协作评分（最多30分）
        COALESCE(DATEDIFF(NOW(), last_active_at), 999) -- 最后活跃距今天数
    INTO login_score, file_score, collaboration_score, days_since_last_active
    FROM team_members
    WHERE id = member_id_param;
    
    -- 计算时间衰减评分
    IF days_since_last_active <= 7 THEN
        SET recency_score = 25;  -- 7天内活跃满分
    ELSEIF days_since_last_active <= 30 THEN
        SET recency_score = 20;  -- 30天内活跃20分
    ELSEIF days_since_last_active <= 90 THEN
        SET recency_score = 10;  -- 90天内活跃10分
    ELSE
        SET recency_score = 0;   -- 超过90天无活跃0分
    END IF;
    
    -- 计算总评分
    SET total_score = login_score + file_score + collaboration_score + recency_score;
    SET total_score = LEAST(total_score, 100.00); -- 确保不超过100分
    
    -- 更新活跃度评分
    UPDATE team_members 
    SET activity_score = total_score,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = member_id_param;
    
    SELECT total_score as calculated_score;
    
END//
DELIMITER ;

-- ============================================================================
-- 团队成员管理视图
-- ============================================================================

-- 团队成员详情视图
CREATE VIEW team_members_detail AS
SELECT 
    tm.id,
    tm.team_id,
    t.name as team_name,
    tm.user_id,
    u.username,
    u.email,
    u.real_name,
    tm.role,
    tm.role_name,
    tm.status,
    tm.display_name,
    tm.job_title,
    tm.department,
    tm.joined_at,
    tm.last_active_at,
    tm.activity_score,
    CASE 
        WHEN tm.last_active_at >= DATE_SUB(NOW(), INTERVAL 7 DAY) THEN 'highly_active'
        WHEN tm.last_active_at >= DATE_SUB(NOW(), INTERVAL 30 DAY) THEN 'active'
        WHEN tm.last_active_at >= DATE_SUB(NOW(), INTERVAL 90 DAY) THEN 'moderately_active'
        ELSE 'inactive'
    END as activity_level,
    ib.username as invited_by_username,
    tm.files_uploaded,
    tm.files_shared,
    tm.storage_used,
    tm.storage_quota
FROM team_members tm
JOIN teams t ON tm.team_id = t.id
JOIN users u ON tm.user_id = u.id
LEFT JOIN users ib ON tm.invited_by = ib.id
WHERE tm.deleted_at IS NULL;

-- 团队角色统计视图
CREATE VIEW team_role_statistics AS
SELECT 
    t.id as team_id,
    t.name as team_name,
    COUNT(tm.id) as total_members,
    COUNT(CASE WHEN tm.role = 'owner' THEN 1 END) as owners,
    COUNT(CASE WHEN tm.role = 'admin' THEN 1 END) as admins,
    COUNT(CASE WHEN tm.role = 'member' THEN 1 END) as members,
    COUNT(CASE WHEN tm.role = 'viewer' THEN 1 END) as viewers,
    COUNT(CASE WHEN tm.role = 'guest' THEN 1 END) as guests,
    COUNT(CASE WHEN tm.status = 'active' THEN 1 END) as active_members,
    COUNT(CASE WHEN tm.status = 'pending' THEN 1 END) as pending_members,
    AVG(tm.activity_score) as avg_activity_score
FROM teams t
LEFT JOIN team_members tm ON t.id = tm.team_id AND tm.deleted_at IS NULL
WHERE t.deleted_at IS NULL
GROUP BY t.id, t.name;

-- 成员活跃度排行视图
CREATE VIEW member_activity_ranking AS
SELECT 
    tm.id,
    tm.team_id,
    t.name as team_name,
    u.username,
    tm.display_name,
    tm.activity_score,
    tm.files_uploaded,
    tm.files_shared,
    tm.collaborations_count,
    tm.last_active_at,
    RANK() OVER (PARTITION BY tm.team_id ORDER BY tm.activity_score DESC) as team_rank,
    RANK() OVER (ORDER BY tm.activity_score DESC) as global_rank
FROM team_members tm
JOIN teams t ON tm.team_id = t.id
JOIN users u ON tm.user_id = u.id
WHERE tm.status = 'active' 
  AND tm.deleted_at IS NULL 
  AND t.deleted_at IS NULL
ORDER BY tm.activity_score DESC;

-- ============================================================================
-- 补充约束检查
-- ============================================================================

-- 部门名称长度约束
ALTER TABLE team_members ADD CONSTRAINT chk_department_length 
CHECK (department IS NULL OR (LENGTH(TRIM(department)) >= 2 AND LENGTH(department) <= 100));

-- 职位名称长度约束
ALTER TABLE team_members ADD CONSTRAINT chk_job_title_length 
CHECK (job_title IS NULL OR (LENGTH(TRIM(job_title)) >= 2 AND LENGTH(job_title) <= 100));

-- 显示名称约束
ALTER TABLE team_members ADD CONSTRAINT chk_display_name_length 
CHECK (display_name IS NULL OR (LENGTH(TRIM(display_name)) >= 1 AND LENGTH(display_name) <= 100));

-- 邀请者不能是自己约束
ALTER TABLE team_members ADD CONSTRAINT chk_cannot_invite_self 
CHECK (user_id != invited_by);

-- 移除者逻辑约束（自己退出或被管理员移除）
ALTER TABLE team_members ADD CONSTRAINT chk_removal_logic 
CHECK (
    (user_id = removed_by AND status = 'left') OR
    (user_id != removed_by) OR
    (removed_by IS NULL)
);