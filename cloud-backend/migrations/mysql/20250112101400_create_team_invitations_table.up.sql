-- +migrate Up
-- 创建迁移: 团队邀请表
-- 版本: 20250112101400
-- 描述: 创建团队邀请表，管理团队邀请流程和状态
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:14:00
-- 依赖: 20250112101200_create_teams_table, 20250112100000_create_users_table
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1
-- 说明: 团队邀请表管理邀请流程，支持邮箱邀请和用户邀请两种模式

-- ============================================================================
-- 团队邀请表 (team_invitations)
-- ============================================================================

CREATE TABLE team_invitations (
    -- 基础标识字段
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '邀请记录唯一标识',
    
    -- 关联的团队和邀请者
    team_id BIGINT UNSIGNED NOT NULL COMMENT '目标团队ID',
    inviter_id BIGINT UNSIGNED NOT NULL COMMENT '邀请发起者用户ID',
    
    -- 被邀请者信息
    invitee_email VARCHAR(100) NOT NULL COMMENT '被邀请者邮箱地址',
    invitee_id BIGINT UNSIGNED DEFAULT NULL COMMENT '被邀请者用户ID（如果已注册）',
    invitee_name VARCHAR(100) DEFAULT NULL COMMENT '被邀请者姓名',
    
    -- 邀请角色和权限
    role ENUM('admin', 'member', 'viewer', 'guest', 'custom') DEFAULT 'member' COMMENT '邀请的角色级别',
    custom_role_name VARCHAR(50) DEFAULT NULL COMMENT '自定义角色名称',
    permissions JSON DEFAULT NULL COMMENT '具体权限配置',
    
    -- 邀请令牌和验证
    invitation_token VARCHAR(64) NOT NULL UNIQUE COMMENT '邀请令牌，用于验证和访问',
    invitation_code VARCHAR(8) DEFAULT NULL COMMENT '短邀请码，便于分享',
    
    -- 邀请内容和个性化
    message TEXT DEFAULT NULL COMMENT '邀请附加消息',
    subject VARCHAR(200) DEFAULT NULL COMMENT '邀请邮件主题',
    custom_welcome_message TEXT DEFAULT NULL COMMENT '自定义欢迎消息',
    
    -- 邀请方式和渠道
    invitation_method ENUM('email', 'link', 'qr_code', 'bulk', 'api', 'referral') DEFAULT 'email' COMMENT '邀请方式',
    invitation_channel VARCHAR(50) DEFAULT NULL COMMENT '邀请渠道标识',
    
    -- 邀请状态和生命周期
    status ENUM('pending', 'sent', 'viewed', 'accepted', 'declined', 'expired', 'cancelled', 'failed') DEFAULT 'pending' COMMENT '邀请状态',
    
    -- 时间管理
    expires_at TIMESTAMP NOT NULL COMMENT '邀请过期时间',
    sent_at TIMESTAMP NULL DEFAULT NULL COMMENT '邀请发送时间',
    first_viewed_at TIMESTAMP NULL DEFAULT NULL COMMENT '首次查看时间',
    responded_at TIMESTAMP NULL DEFAULT NULL COMMENT '响应时间（接受或拒绝）',
    
    -- 邀请追踪和统计
    view_count INT UNSIGNED DEFAULT 0 COMMENT '邀请查看次数',
    reminder_count INT UNSIGNED DEFAULT 0 COMMENT '提醒次数',
    last_reminder_at TIMESTAMP NULL DEFAULT NULL COMMENT '最后提醒时间',
    
    -- 响应和处理信息
    decline_reason ENUM('not_interested', 'already_member', 'too_busy', 'privacy_concerns', 'other') DEFAULT NULL COMMENT '拒绝原因',
    decline_message TEXT DEFAULT NULL COMMENT '拒绝详细说明',
    response_ip VARCHAR(45) DEFAULT NULL COMMENT '响应IP地址',
    response_user_agent TEXT DEFAULT NULL COMMENT '响应时的用户代理',
    
    -- 邀请限制和条件
    max_uses INT UNSIGNED DEFAULT 1 COMMENT '最大使用次数',
    used_count INT UNSIGNED DEFAULT 0 COMMENT '已使用次数',
    ip_restrictions JSON DEFAULT NULL COMMENT 'IP访问限制',
    domain_restrictions JSON DEFAULT NULL COMMENT '邮箱域名限制',
    
    -- 自动处理设置
    auto_remind BOOLEAN DEFAULT TRUE COMMENT '是否自动发送提醒',
    auto_expire BOOLEAN DEFAULT TRUE COMMENT '是否自动过期',
    auto_accept BOOLEAN DEFAULT FALSE COMMENT '是否自动接受（内部邀请）',
    
    -- 邀请来源和上下文
    source_context JSON DEFAULT NULL COMMENT '邀请来源上下文信息',
    referrer_url VARCHAR(1000) DEFAULT NULL COMMENT '邀请页面来源URL',
    campaign_id VARCHAR(100) DEFAULT NULL COMMENT '营销活动ID',
    
    -- 安全和防滥用
    security_hash VARCHAR(64) DEFAULT NULL COMMENT '安全哈希，防止伪造',
    rate_limit_key VARCHAR(100) DEFAULT NULL COMMENT '限速键，防止频繁邀请',
    
    -- 批量邀请相关
    batch_id VARCHAR(36) DEFAULT NULL COMMENT '批量邀请批次ID',
    batch_sequence INT UNSIGNED DEFAULT NULL COMMENT '批次内序号',
    
    -- 扩展属性
    metadata JSON DEFAULT NULL COMMENT '扩展元数据',
    tags JSON DEFAULT NULL COMMENT '邀请标签',
    custom_fields JSON DEFAULT NULL COMMENT '自定义字段',
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '邀请创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '邀请更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',
    
    -- 业务索引设计
    INDEX idx_team_invitations_team_id (team_id) COMMENT '团队ID索引，查询团队邀请',
    INDEX idx_team_invitations_inviter_id (inviter_id) COMMENT '邀请者索引，查询发出的邀请',
    INDEX idx_team_invitations_invitee_email (invitee_email) COMMENT '被邀请者邮箱索引',
    INDEX idx_team_invitations_invitee_id (invitee_id) COMMENT '被邀请者用户ID索引',
    INDEX idx_team_invitations_invitation_token (invitation_token) COMMENT '邀请令牌索引，验证访问',
    INDEX idx_team_invitations_invitation_code (invitation_code) COMMENT '邀请码索引',
    INDEX idx_team_invitations_status (status) COMMENT '邀请状态索引',
    INDEX idx_team_invitations_expires_at (expires_at) COMMENT '过期时间索引，清理任务',
    INDEX idx_team_invitations_sent_at (sent_at) COMMENT '发送时间索引',
    INDEX idx_team_invitations_batch_id (batch_id) COMMENT '批量邀请批次索引',
    INDEX idx_team_invitations_campaign_id (campaign_id) COMMENT '营销活动索引',
    INDEX idx_team_invitations_created_at (created_at) COMMENT '创建时间索引',
    INDEX idx_team_invitations_deleted_at (deleted_at) COMMENT '软删除索引',
    
    -- 复合业务索引
    INDEX idx_team_invitations_team_status_expires (team_id, status, expires_at) COMMENT '团队状态过期复合索引',
    INDEX idx_team_invitations_inviter_status_time (inviter_id, status, created_at DESC) COMMENT '邀请者状态时间复合索引',
    INDEX idx_team_invitations_email_status (invitee_email, status, expires_at) COMMENT '邮箱状态过期复合索引',
    INDEX idx_team_invitations_pending_reminders (status, auto_remind, last_reminder_at) COMMENT '待提醒邀请索引',
    INDEX idx_team_invitations_batch_processing (batch_id, batch_sequence, status) COMMENT '批量处理优化索引',
    
    -- 全文搜索索引
    FULLTEXT idx_fulltext_search (message, subject, custom_welcome_message) COMMENT '邀请内容全文搜索',
    
    -- 外键约束
    CONSTRAINT fk_team_invitations_team_id 
        FOREIGN KEY (team_id) REFERENCES teams(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_team_invitations_inviter_id 
        FOREIGN KEY (inviter_id) REFERENCES users(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_team_invitations_invitee_id 
        FOREIGN KEY (invitee_id) REFERENCES users(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE
        
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='团队邀请表 - 管理团队邀请流程、状态跟踪和批量处理'
  ROW_FORMAT=DYNAMIC;

-- ============================================================================
-- 团队邀请表约束和检查
-- ============================================================================

-- 邮箱格式约束
ALTER TABLE team_invitations ADD CONSTRAINT chk_invitee_email_format 
CHECK (
    invitee_email REGEXP '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$'
);

-- 自定义角色约束
ALTER TABLE team_invitations ADD CONSTRAINT chk_custom_role_constraint 
CHECK (
    (role != 'custom') OR 
    (role = 'custom' AND custom_role_name IS NOT NULL)
);

-- 时间逻辑约束
ALTER TABLE team_invitations ADD CONSTRAINT chk_invitation_time_logic 
CHECK (
    expires_at > created_at AND
    (sent_at IS NULL OR sent_at >= created_at) AND
    (first_viewed_at IS NULL OR first_viewed_at >= created_at) AND
    (responded_at IS NULL OR responded_at >= created_at)
);

-- 使用次数约束
ALTER TABLE team_invitations ADD CONSTRAINT chk_invitation_usage_limits 
CHECK (
    max_uses > 0 AND
    used_count >= 0 AND
    used_count <= max_uses
);

-- 统计字段约束
ALTER TABLE team_invitations ADD CONSTRAINT chk_invitation_counters 
CHECK (
    view_count >= 0 AND
    reminder_count >= 0 AND
    (batch_sequence IS NULL OR batch_sequence > 0)
);

-- 响应状态逻辑约束
ALTER TABLE team_invitations ADD CONSTRAINT chk_response_logic 
CHECK (
    (status NOT IN ('accepted', 'declined') AND responded_at IS NULL) OR
    (status IN ('accepted', 'declined') AND responded_at IS NOT NULL)
);

-- 拒绝信息约束
ALTER TABLE team_invitations ADD CONSTRAINT chk_decline_info_logic 
CHECK (
    (status != 'declined') OR 
    (status = 'declined' AND decline_reason IS NOT NULL)
);

-- IP地址格式约束
ALTER TABLE team_invitations ADD CONSTRAINT chk_response_ip_format 
CHECK (
    response_ip IS NULL OR
    response_ip REGEXP '^([0-9]{1,3}\\.){3}[0-9]{1,3}$' OR  -- IPv4
    response_ip REGEXP '^[0-9a-fA-F:]+$'                    -- IPv6（简化）
);

-- 邀请码格式约束
ALTER TABLE team_invitations ADD CONSTRAINT chk_invitation_code_format 
CHECK (
    invitation_code IS NULL OR 
    invitation_code REGEXP '^[A-Z0-9]{6,8}$'
);

-- JSON字段验证
ALTER TABLE team_invitations ADD CONSTRAINT chk_team_invitations_json_valid 
CHECK (
    (permissions IS NULL OR JSON_VALID(permissions)) AND
    (ip_restrictions IS NULL OR JSON_VALID(ip_restrictions)) AND
    (domain_restrictions IS NULL OR JSON_VALID(domain_restrictions)) AND
    (source_context IS NULL OR JSON_VALID(source_context)) AND
    (metadata IS NULL OR JSON_VALID(metadata)) AND
    (tags IS NULL OR JSON_VALID(tags)) AND
    (custom_fields IS NULL OR JSON_VALID(custom_fields))
);

-- ============================================================================
-- 团队邀请管理触发器
-- ============================================================================

-- 邀请创建初始化触发器
DELIMITER //
CREATE TRIGGER team_invitations_initialize
BEFORE INSERT ON team_invitations
FOR EACH ROW
BEGIN
    -- 生成邀请令牌（如果未提供）
    IF NEW.invitation_token IS NULL THEN
        SET NEW.invitation_token = SHA2(CONCAT(NEW.team_id, NEW.invitee_email, UNIX_TIMESTAMP(), RAND()), 256);
    END IF;
    
    -- 生成短邀请码
    IF NEW.invitation_code IS NULL THEN
        SET NEW.invitation_code = UPPER(SUBSTRING(MD5(CONCAT(NEW.invitation_token, RAND())), 1, 8));
    END IF;
    
    -- 生成安全哈希
    SET NEW.security_hash = SHA2(CONCAT(NEW.invitation_token, 'security_salt', NEW.team_id), 256);
    
    -- 设置默认过期时间（如果未指定）
    IF NEW.expires_at IS NULL THEN
        SET NEW.expires_at = DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 7 DAY);
    END IF;
    
    -- 生成批次ID（如果是批量邀请）
    IF NEW.batch_id IS NULL AND NEW.invitation_method = 'bulk' THEN
        SET NEW.batch_id = UUID();
    END IF;
    
    -- 初始化默认权限（如果未指定）
    IF NEW.permissions IS NULL AND NEW.role != 'custom' THEN
        CASE NEW.role
            WHEN 'admin' THEN
                SET NEW.permissions = JSON_OBJECT(
                    'manage_members', true,
                    'manage_files', true,
                    'manage_settings', true,
                    'view_analytics', true
                );
            WHEN 'member' THEN
                SET NEW.permissions = JSON_OBJECT(
                    'upload_files', true,
                    'share_files', true,
                    'comment', true,
                    'collaborate', true
                );
            WHEN 'viewer' THEN
                SET NEW.permissions = JSON_OBJECT(
                    'view_files', true,
                    'download_files', true,
                    'comment', false
                );
            ELSE
                SET NEW.permissions = JSON_OBJECT('view_files', true);
        END CASE;
    END IF;
END//
DELIMITER ;

-- 邀请状态变更触发器
DELIMITER //
CREATE TRIGGER team_invitations_status_update
AFTER UPDATE ON team_invitations
FOR EACH ROW
BEGIN
    -- 当邀请被接受时，自动创建团队成员记录
    IF OLD.status != 'accepted' AND NEW.status = 'accepted' THEN
        INSERT INTO team_members (
            team_id, user_id, role, permissions, status, 
            invited_by, invitation_method, joined_at
        ) VALUES (
            NEW.team_id, 
            NEW.invitee_id, 
            NEW.role, 
            NEW.permissions, 
            'active',
            NEW.inviter_id, 
            NEW.invitation_method, 
            CURRENT_TIMESTAMP
        )
        ON DUPLICATE KEY UPDATE
            status = 'active',
            joined_at = CURRENT_TIMESTAMP,
            updated_at = CURRENT_TIMESTAMP;
    END IF;
    
    -- 更新查看统计
    IF NEW.view_count > OLD.view_count AND OLD.first_viewed_at IS NULL THEN
        UPDATE team_invitations 
        SET first_viewed_at = CURRENT_TIMESTAMP 
        WHERE id = NEW.id;
    END IF;
    
    -- 自动过期检查
    IF NEW.expires_at <= CURRENT_TIMESTAMP AND NEW.status = 'pending' THEN
        UPDATE team_invitations 
        SET status = 'expired' 
        WHERE id = NEW.id;
    END IF;
END//
DELIMITER ;

-- 邀请过期检查触发器
DELIMITER //
CREATE TRIGGER team_invitations_expiry_check
BEFORE UPDATE ON team_invitations
FOR EACH ROW
BEGIN
    -- 检查邀请是否过期
    IF NEW.expires_at <= CURRENT_TIMESTAMP AND OLD.status = 'pending' THEN
        SET NEW.status = 'expired';
    END IF;
    
    -- 防止过期邀请被接受
    IF NEW.status = 'accepted' AND NEW.expires_at <= CURRENT_TIMESTAMP THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '邀请已过期，无法接受';
    END IF;
    
    -- 检查邀请使用次数
    IF NEW.used_count >= NEW.max_uses AND NEW.status = 'pending' THEN
        SET NEW.status = 'expired';
    END IF;
END//
DELIMITER ;

-- ============================================================================
-- 团队邀请管理存储过程
-- ============================================================================

-- 创建团队邀请存储过程
DELIMITER //
CREATE PROCEDURE CreateTeamInvitation(
    IN team_id_param BIGINT UNSIGNED,
    IN inviter_id_param BIGINT UNSIGNED,
    IN invitee_email_param VARCHAR(100),
    IN role_param VARCHAR(20),
    IN message_param TEXT,
    IN expires_days INT,
    OUT invitation_id BIGINT UNSIGNED,
    OUT invitation_token VARCHAR(64),
    OUT success BOOLEAN,
    OUT error_message VARCHAR(500)
)
BEGIN
    DECLARE team_exists INT DEFAULT 0;
    DECLARE inviter_can_invite BOOLEAN DEFAULT FALSE;
    DECLARE already_invited INT DEFAULT 0;
    DECLARE already_member INT DEFAULT 0;
    DECLARE invitee_user_id BIGINT UNSIGNED DEFAULT NULL;
    
    SET success = FALSE;
    SET error_message = NULL;
    SET invitation_id = NULL;
    SET invitation_token = NULL;
    
    -- 检查团队是否存在且有效
    SELECT COUNT(*) INTO team_exists
    FROM teams
    WHERE id = team_id_param AND is_active = TRUE AND deleted_at IS NULL;
    
    IF team_exists = 0 THEN
        SET error_message = '团队不存在或已被禁用';
        LEAVE CreateTeamInvitation;
    END IF;
    
    -- 检查邀请者权限
    SELECT COUNT(*) > 0 INTO inviter_can_invite
    FROM team_members tm
    WHERE tm.team_id = team_id_param 
      AND tm.user_id = inviter_id_param 
      AND tm.status = 'active'
      AND tm.role IN ('owner', 'admin');
    
    IF NOT inviter_can_invite THEN
        SET error_message = '您没有邀请成员的权限';
        LEAVE CreateTeamInvitation;
    END IF;
    
    -- 检查是否已经是成员
    SELECT user_id INTO invitee_user_id
    FROM users 
    WHERE email = invitee_email_param AND status = 'active';
    
    IF invitee_user_id IS NOT NULL THEN
        SELECT COUNT(*) INTO already_member
        FROM team_members
        WHERE team_id = team_id_param 
          AND user_id = invitee_user_id 
          AND status IN ('active', 'pending');
        
        IF already_member > 0 THEN
            SET error_message = '该用户已经是团队成员';
            LEAVE CreateTeamInvitation;
        END IF;
    END IF;
    
    -- 检查是否已有有效邀请
    SELECT COUNT(*) INTO already_invited
    FROM team_invitations
    WHERE team_id = team_id_param 
      AND invitee_email = invitee_email_param 
      AND status = 'pending'
      AND expires_at > CURRENT_TIMESTAMP;
    
    IF already_invited > 0 THEN
        SET error_message = '该邮箱已有待处理的邀请';
        LEAVE CreateTeamInvitation;
    END IF;
    
    -- 创建邀请
    INSERT INTO team_invitations (
        team_id, inviter_id, invitee_email, invitee_id, role, message,
        expires_at, invitation_method
    ) VALUES (
        team_id_param, 
        inviter_id_param, 
        invitee_email_param, 
        invitee_user_id, 
        role_param, 
        message_param,
        DATE_ADD(CURRENT_TIMESTAMP, INTERVAL COALESCE(expires_days, 7) DAY),
        'email'
    );
    
    SET invitation_id = LAST_INSERT_ID();
    
    -- 获取生成的邀请令牌
    SELECT invitation_token INTO invitation_token
    FROM team_invitations
    WHERE id = invitation_id;
    
    SET success = TRUE;
    
END//
DELIMITER ;

-- 批量清理过期邀请存储过程
DELIMITER //
CREATE PROCEDURE CleanExpiredInvitations()
BEGIN
    DECLARE affected_rows INT DEFAULT 0;
    
    -- 标记过期邀请
    UPDATE team_invitations 
    SET status = 'expired',
        updated_at = CURRENT_TIMESTAMP
    WHERE status = 'pending' 
      AND expires_at <= CURRENT_TIMESTAMP;
    
    SET affected_rows = ROW_COUNT();
    
    -- 删除30天前的已过期邀请
    DELETE FROM team_invitations 
    WHERE status = 'expired' 
      AND expires_at <= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL 30 DAY);
    
    SELECT CONCAT('已处理 ', affected_rows, ' 个过期邀请，删除 ', ROW_COUNT(), ' 个历史邀请') as result;
END//
DELIMITER ;

-- 邀请统计分析存储过程
DELIMITER //
CREATE PROCEDURE GetInvitationStatistics(
    IN team_id_param BIGINT UNSIGNED,
    IN days_param INT DEFAULT 30
)
BEGIN
    SELECT 
        COUNT(*) as total_invitations,
        COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_invitations,
        COUNT(CASE WHEN status = 'accepted' THEN 1 END) as accepted_invitations,
        COUNT(CASE WHEN status = 'declined' THEN 1 END) as declined_invitations,
        COUNT(CASE WHEN status = 'expired' THEN 1 END) as expired_invitations,
        ROUND(COUNT(CASE WHEN status = 'accepted' THEN 1 END) * 100.0 / COUNT(*), 2) as acceptance_rate,
        AVG(view_count) as avg_views_per_invitation,
        AVG(TIMESTAMPDIFF(HOUR, created_at, COALESCE(responded_at, NOW()))) as avg_response_time_hours
    FROM team_invitations
    WHERE team_id = team_id_param
      AND created_at >= DATE_SUB(CURRENT_DATE, INTERVAL days_param DAY)
      AND deleted_at IS NULL;
END//
DELIMITER ;

-- ============================================================================
-- 团队邀请管理视图
-- ============================================================================

-- 活跃邀请概览视图
CREATE VIEW active_invitations_overview AS
SELECT 
    ti.id,
    ti.team_id,
    t.name as team_name,
    ti.invitee_email,
    u_invitee.username as invitee_username,
    ti.role,
    ti.status,
    u_inviter.username as inviter_username,
    ti.expires_at,
    TIMESTAMPDIFF(HOUR, NOW(), ti.expires_at) as hours_until_expire,
    ti.view_count,
    ti.created_at,
    ti.first_viewed_at
FROM team_invitations ti
JOIN teams t ON ti.team_id = t.id
JOIN users u_inviter ON ti.inviter_id = u_inviter.id
LEFT JOIN users u_invitee ON ti.invitee_id = u_invitee.id
WHERE ti.status = 'pending' 
  AND ti.expires_at > NOW()
  AND ti.deleted_at IS NULL;

-- 邀请转化率统计视图
CREATE VIEW invitation_conversion_stats AS
SELECT 
    t.id as team_id,
    t.name as team_name,
    COUNT(ti.id) as total_invitations,
    COUNT(CASE WHEN ti.status = 'accepted' THEN 1 END) as accepted_count,
    COUNT(CASE WHEN ti.status = 'declined' THEN 1 END) as declined_count,
    COUNT(CASE WHEN ti.status = 'expired' THEN 1 END) as expired_count,
    ROUND(COUNT(CASE WHEN ti.status = 'accepted' THEN 1 END) * 100.0 / COUNT(ti.id), 2) as acceptance_rate,
    ROUND(COUNT(CASE WHEN ti.status = 'declined' THEN 1 END) * 100.0 / COUNT(ti.id), 2) as decline_rate,
    AVG(ti.view_count) as avg_views,
    AVG(TIMESTAMPDIFF(HOUR, ti.created_at, COALESCE(ti.responded_at, NOW()))) as avg_response_hours
FROM teams t
LEFT JOIN team_invitations ti ON t.id = ti.team_id 
    AND ti.created_at >= DATE_SUB(NOW(), INTERVAL 90 DAY)
    AND ti.deleted_at IS NULL
WHERE t.deleted_at IS NULL
GROUP BY t.id, t.name
HAVING total_invitations > 0;

-- 即将过期邀请提醒视图
CREATE VIEW expiring_invitations_alert AS
SELECT 
    ti.id,
    ti.team_id,
    t.name as team_name,
    ti.invitee_email,
    ti.role,
    ti.expires_at,
    TIMESTAMPDIFF(HOUR, NOW(), ti.expires_at) as hours_until_expire,
    u_inviter.username as inviter_username,
    u_inviter.email as inviter_email
FROM team_invitations ti
JOIN teams t ON ti.team_id = t.id
JOIN users u_inviter ON ti.inviter_id = u_inviter.id
WHERE ti.status = 'pending'
  AND ti.expires_at > NOW()
  AND ti.expires_at <= DATE_ADD(NOW(), INTERVAL 24 HOUR)
  AND ti.deleted_at IS NULL
ORDER BY ti.expires_at ASC;

-- ============================================================================
-- 补充约束检查
-- ============================================================================

-- 邀请者不能邀请自己约束
ALTER TABLE team_invitations ADD CONSTRAINT chk_cannot_invite_self_team 
CHECK (
    (invitee_id IS NULL) OR 
    (inviter_id != invitee_id)
);

-- 邀请码格式约束
ALTER TABLE team_invitations ADD CONSTRAINT chk_invitation_code_format 
CHECK (
    invitation_code IS NULL OR 
    invitation_code REGEXP '^[A-Z0-9]{6,12}$'
);

-- 批次序号约束
ALTER TABLE team_invitations ADD CONSTRAINT chk_batch_sequence_positive 
CHECK (batch_sequence IS NULL OR batch_sequence > 0);

-- 过期时间逻辑约束
ALTER TABLE team_invitations ADD CONSTRAINT chk_invitation_time_logic 
CHECK (
    expires_at > created_at AND
    (sent_at IS NULL OR sent_at >= created_at) AND
    (responded_at IS NULL OR responded_at >= created_at)
);

-- 提醒间隔约束
ALTER TABLE team_invitations ADD CONSTRAINT chk_reminder_interval 
CHECK (reminder_interval_days IS NULL OR (reminder_interval_days >= 1 AND reminder_interval_days <= 30));