-- +migrate Up
-- 创建迁移: 团队表
-- 版本: 20250112101200
-- 描述: 创建团队表，管理团队信息和协作设置
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:12:00
-- 依赖: 20250112100000_create_users_table
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1
-- 说明: 团队表是协作功能的核心，管理团队基本信息、成员配额和存储分配

-- ============================================================================
-- 团队表 (teams)
-- ============================================================================

CREATE TABLE teams (
    -- 基础标识字段
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '团队唯一标识',
    
    -- 团队基本信息
    name VARCHAR(100) NOT NULL COMMENT '团队名称，显示名称',
    slug VARCHAR(100) NOT NULL UNIQUE COMMENT '团队标识符，用于URL和API',
    description TEXT DEFAULT NULL COMMENT '团队描述和简介',
    
    -- 团队品牌和展示
    avatar_url VARCHAR(500) DEFAULT NULL COMMENT '团队头像URL',
    banner_url VARCHAR(500) DEFAULT NULL COMMENT '团队横幅图片URL',
    website_url VARCHAR(500) DEFAULT NULL COMMENT '团队官网或主页URL',
    
    -- 团队管理和所有权
    owner_id BIGINT UNSIGNED NOT NULL COMMENT '团队所有者用户ID，具有最高权限',
    
    -- 团队类型和分类
    team_type ENUM('personal', 'business', 'education', 'nonprofit', 'enterprise') DEFAULT 'business' COMMENT '团队类型',
    industry VARCHAR(100) DEFAULT NULL COMMENT '所属行业',
    organization_size ENUM('1-10', '11-50', '51-200', '201-1000', '1000+') DEFAULT '1-10' COMMENT '组织规模',
    
    -- 成员管理
    max_members INT UNSIGNED DEFAULT 50 COMMENT '最大成员数限制',
    member_count INT UNSIGNED DEFAULT 1 COMMENT '当前成员数，包含所有者',
    active_member_count INT UNSIGNED DEFAULT 1 COMMENT '活跃成员数（30天内有活动）',
    
    -- 存储配额管理
    storage_quota BIGINT UNSIGNED DEFAULT 107374182400 COMMENT '团队存储配额，默认100GB（字节）',
    storage_used BIGINT UNSIGNED DEFAULT 0 COMMENT '已使用存储空间（字节）',
    storage_warning_threshold DECIMAL(3,2) DEFAULT 0.80 COMMENT '存储警告阈值（80%）',
    
    -- 功能配额管理
    max_projects INT UNSIGNED DEFAULT 10 COMMENT '最大项目数限制',
    max_shared_links INT UNSIGNED DEFAULT 100 COMMENT '最大分享链接数',
    max_file_size BIGINT UNSIGNED DEFAULT 5368709120 COMMENT '单文件最大大小，默认5GB',
    
    -- 团队状态和生命周期
    status ENUM('active', 'suspended', 'inactive', 'deleted') DEFAULT 'active' COMMENT '团队状态',
    is_active BOOLEAN DEFAULT TRUE COMMENT '团队是否有效',
    is_verified BOOLEAN DEFAULT FALSE COMMENT '团队是否已验证',
    verification_level ENUM('none', 'basic', 'business', 'enterprise') DEFAULT 'none' COMMENT '验证级别',
    
    -- 订阅和计费信息
    plan_type ENUM('free', 'basic', 'professional', 'enterprise', 'custom') DEFAULT 'free' COMMENT '订阅计划类型',
    plan_expires_at TIMESTAMP NULL DEFAULT NULL COMMENT '计划过期时间',
    billing_email VARCHAR(100) DEFAULT NULL COMMENT '账单邮箱',
    
    -- 团队设置和配置
    settings JSON DEFAULT NULL COMMENT '团队详细设置配置',
    permissions JSON DEFAULT NULL COMMENT '团队级别权限配置',
    integrations JSON DEFAULT NULL COMMENT '第三方集成配置',
    
    -- 安全和合规设置
    security_settings JSON DEFAULT NULL COMMENT '安全设置：2FA要求、IP限制等',
    compliance_settings JSON DEFAULT NULL COMMENT '合规设置：数据保留、审计等',
    
    -- 地理位置和本地化
    timezone VARCHAR(50) DEFAULT 'UTC' COMMENT '团队默认时区',
    locale VARCHAR(10) DEFAULT 'en-US' COMMENT '团队默认语言地区',
    country_code VARCHAR(2) DEFAULT NULL COMMENT '团队所在国家代码',
    
    -- 联系信息
    contact_email VARCHAR(100) DEFAULT NULL COMMENT '团队联系邮箱',
    contact_phone VARCHAR(20) DEFAULT NULL COMMENT '团队联系电话',
    address TEXT DEFAULT NULL COMMENT '团队地址信息',
    
    -- 活动和统计
    last_activity_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '团队最后活动时间',
    total_files_created INT UNSIGNED DEFAULT 0 COMMENT '团队创建的文件总数',
    total_shares_created INT UNSIGNED DEFAULT 0 COMMENT '团队创建的分享总数',
    
    -- 标签和元数据
    tags JSON DEFAULT NULL COMMENT '团队标签，便于分类和搜索',
    metadata JSON DEFAULT NULL COMMENT '扩展元数据',
    custom_fields JSON DEFAULT NULL COMMENT '自定义字段',
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '团队创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '团队信息更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',
    
    -- 业务索引设计
    INDEX idx_teams_name (name) COMMENT '团队名称索引，搜索和排序',
    INDEX idx_teams_slug (slug) COMMENT '团队标识符索引，URL访问',
    INDEX idx_teams_owner_id (owner_id) COMMENT '所有者索引，查询用户拥有的团队',
    INDEX idx_teams_team_type (team_type) COMMENT '团队类型索引',
    INDEX idx_teams_status (status) COMMENT '团队状态索引',
    INDEX idx_teams_is_active (is_active) COMMENT '有效状态索引',
    INDEX idx_teams_plan_type (plan_type) COMMENT '订阅计划索引',
    INDEX idx_teams_verification_level (verification_level) COMMENT '验证级别索引',
    INDEX idx_teams_created_at (created_at) COMMENT '创建时间索引',
    INDEX idx_teams_last_activity_at (last_activity_at) COMMENT '最后活动时间索引',
    INDEX idx_teams_storage_used (storage_used) COMMENT '存储使用量索引',
    INDEX idx_teams_member_count (member_count) COMMENT '成员数量索引',
    INDEX idx_teams_deleted_at (deleted_at) COMMENT '软删除索引',
    INDEX idx_teams_country_code (country_code) COMMENT '国家代码索引，地理统计',
    
    -- 复合业务索引
    INDEX idx_teams_owner_status (owner_id, status, is_active) COMMENT '所有者状态复合索引',
    INDEX idx_teams_type_plan (team_type, plan_type, is_active) COMMENT '类型计划复合索引',
    INDEX idx_teams_activity_stats (last_activity_at DESC, member_count, storage_used) COMMENT '活跃度统计索引',
    INDEX idx_teams_storage_quota_usage (storage_used, storage_quota, storage_warning_threshold) COMMENT '存储使用情况索引',
    
    -- 全文搜索索引
    FULLTEXT idx_fulltext_search (name, description) COMMENT '团队名称描述全文搜索',
    
    -- 外键约束
    CONSTRAINT fk_teams_owner_id 
        FOREIGN KEY (owner_id) REFERENCES users(id) 
        ON DELETE RESTRICT 
        ON UPDATE CASCADE
        
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='团队信息表 - 管理团队基本信息、配额和协作设置'
  ROW_FORMAT=DYNAMIC;

-- ============================================================================
-- 团队表约束和检查
-- ============================================================================

-- 团队名称长度约束
ALTER TABLE teams ADD CONSTRAINT chk_team_name_length 
CHECK (LENGTH(TRIM(name)) >= 2 AND LENGTH(name) <= 100);

-- 团队标识符格式约束
ALTER TABLE teams ADD CONSTRAINT chk_team_slug_format 
CHECK (
    slug REGEXP '^[a-z0-9][a-z0-9-]*[a-z0-9]$' AND
    LENGTH(slug) >= 3 AND
    LENGTH(slug) <= 100
);

-- 成员数量约束
ALTER TABLE teams ADD CONSTRAINT chk_member_count_limits 
CHECK (
    member_count >= 1 AND
    member_count <= max_members AND
    active_member_count <= member_count AND
    max_members > 0 AND
    max_members <= 10000
);

-- 存储配额约束
ALTER TABLE teams ADD CONSTRAINT chk_storage_quotas 
CHECK (
    storage_quota > 0 AND
    storage_used >= 0 AND
    storage_used <= storage_quota AND
    storage_warning_threshold > 0 AND
    storage_warning_threshold <= 1
);

-- 功能配额约束
ALTER TABLE teams ADD CONSTRAINT chk_feature_quotas 
CHECK (
    max_projects >= 0 AND
    max_shared_links >= 0 AND
    max_file_size > 0
);

-- 统计字段约束
ALTER TABLE teams ADD CONSTRAINT chk_statistics_counts 
CHECK (
    total_files_created >= 0 AND
    total_shares_created >= 0
);

-- 国家代码格式约束
ALTER TABLE teams ADD CONSTRAINT chk_country_code_format 
CHECK (
    country_code IS NULL OR 
    country_code REGEXP '^[A-Z]{2}$'
);

-- 邮箱格式约束
ALTER TABLE teams ADD CONSTRAINT chk_email_formats 
CHECK (
    (billing_email IS NULL OR billing_email REGEXP '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$') AND
    (contact_email IS NULL OR contact_email REGEXP '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$')
);

-- JSON字段验证
ALTER TABLE teams ADD CONSTRAINT chk_teams_json_valid 
CHECK (
    (settings IS NULL OR JSON_VALID(settings)) AND
    (permissions IS NULL OR JSON_VALID(permissions)) AND
    (integrations IS NULL OR JSON_VALID(integrations)) AND
    (security_settings IS NULL OR JSON_VALID(security_settings)) AND
    (compliance_settings IS NULL OR JSON_VALID(compliance_settings)) AND
    (tags IS NULL OR JSON_VALID(tags)) AND
    (metadata IS NULL OR JSON_VALID(metadata)) AND
    (custom_fields IS NULL OR JSON_VALID(custom_fields))
);

-- ============================================================================
-- 团队管理触发器
-- ============================================================================

-- 团队创建初始化触发器
DELIMITER //
CREATE TRIGGER teams_initialize_settings
BEFORE INSERT ON teams
FOR EACH ROW
BEGIN
    -- 自动生成团队标识符（如果未提供）
    IF NEW.slug IS NULL OR NEW.slug = '' THEN
        SET NEW.slug = LOWER(REPLACE(REPLACE(NEW.name, ' ', '-'), '_', '-'));
        -- 确保标识符唯一性（简化处理）
        SET NEW.slug = CONCAT(NEW.slug, '-', UNIX_TIMESTAMP());
    END IF;
    
    -- 初始化默认设置
    IF NEW.settings IS NULL THEN
        SET NEW.settings = JSON_OBJECT(
            'file_versioning', true,
            'auto_save', true,
            'collaboration_mode', 'open',
            'default_permissions', JSON_OBJECT(
                'view', true,
                'comment', true,
                'edit', false,
                'admin', false
            )
        );
    END IF;
    
    -- 初始化安全设置
    IF NEW.security_settings IS NULL THEN
        SET NEW.security_settings = JSON_OBJECT(
            'require_2fa', false,
            'session_timeout', 7200,
            'ip_restrictions_enabled', false,
            'allowed_domains', JSON_ARRAY()
        );
    END IF;
    
    -- 根据计划类型设置默认配额
    IF NEW.plan_type = 'free' THEN
        SET NEW.max_members = LEAST(NEW.max_members, 5);
        SET NEW.storage_quota = LEAST(NEW.storage_quota, 5368709120); -- 5GB
        SET NEW.max_projects = LEAST(NEW.max_projects, 3);
    ELSEIF NEW.plan_type = 'basic' THEN
        SET NEW.max_members = LEAST(NEW.max_members, 25);
        SET NEW.storage_quota = LEAST(NEW.storage_quota, 107374182400); -- 100GB
        SET NEW.max_projects = LEAST(NEW.max_projects, 10);
    END IF;
END//
DELIMITER ;

-- 团队统计更新触发器
DELIMITER //
CREATE TRIGGER teams_update_stats
BEFORE UPDATE ON teams
FOR EACH ROW
BEGIN
    -- 更新最后活动时间（如果有重要字段变化）
    IF OLD.storage_used != NEW.storage_used OR 
       OLD.member_count != NEW.member_count OR
       OLD.settings != NEW.settings THEN
        SET NEW.last_activity_at = CURRENT_TIMESTAMP;
    END IF;
    
    -- 存储警告检查
    IF NEW.storage_used >= (NEW.storage_quota * NEW.storage_warning_threshold) THEN
        -- 这里可以触发存储警告通知
        -- 实际实现中应该通过应用层处理
        SET NEW.metadata = JSON_SET(
            COALESCE(NEW.metadata, JSON_OBJECT()),
            '$.storage_warning_triggered_at',
            CURRENT_TIMESTAMP
        );
    END IF;
    
    -- 计划过期检查
    IF NEW.plan_expires_at IS NOT NULL AND NEW.plan_expires_at <= CURRENT_TIMESTAMP THEN
        SET NEW.plan_type = 'free';
        -- 重置为免费计划限制
        SET NEW.max_members = LEAST(NEW.max_members, 5);
        SET NEW.max_projects = LEAST(NEW.max_projects, 3);
    END IF;
END//
DELIMITER ;

-- 团队删除检查触发器
DELIMITER //
CREATE TRIGGER teams_deletion_check
BEFORE UPDATE ON teams
FOR EACH ROW
BEGIN
    -- 软删除时更新相关状态
    IF OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL THEN
        SET NEW.status = 'deleted';
        SET NEW.is_active = FALSE;
    END IF;
    
    -- 恢复时重置状态
    IF OLD.deleted_at IS NOT NULL AND NEW.deleted_at IS NULL THEN
        SET NEW.status = 'active';
        SET NEW.is_active = TRUE;
    END IF;
END//
DELIMITER ;

-- ============================================================================
-- 团队管理存储过程
-- ============================================================================

-- 团队统计信息汇总存储过程
DELIMITER //
CREATE PROCEDURE GetTeamStatistics(
    IN team_id_param BIGINT UNSIGNED
)
BEGIN
    SELECT 
        t.id,
        t.name,
        t.member_count,
        t.active_member_count,
        t.storage_used,
        t.storage_quota,
        ROUND(t.storage_used * 100.0 / t.storage_quota, 2) as storage_usage_percentage,
        t.total_files_created,
        t.total_shares_created,
        t.last_activity_at,
        DATEDIFF(NOW(), t.created_at) as days_since_creation,
        -- 成员活跃度统计
        (SELECT COUNT(*) FROM team_members tm WHERE tm.team_id = t.id AND tm.status = 'active') as active_members,
        -- 最近30天创建的文件数（需要关联文件表）
        (SELECT COUNT(*) FROM files f 
         JOIN team_members tm ON f.user_id = tm.user_id 
         WHERE tm.team_id = t.id 
           AND f.created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)) as recent_files_count
    FROM teams t
    WHERE t.id = team_id_param AND t.deleted_at IS NULL;
END//
DELIMITER ;

-- 团队配额检查存储过程
DELIMITER //
CREATE PROCEDURE CheckTeamQuotas(
    IN team_id_param BIGINT UNSIGNED,
    OUT storage_available BIGINT,
    OUT members_available INT,
    OUT projects_available INT,
    OUT within_limits BOOLEAN
)
BEGIN
    DECLARE current_storage BIGINT DEFAULT 0;
    DECLARE current_members INT DEFAULT 0;
    DECLARE current_projects INT DEFAULT 0;
    DECLARE max_storage BIGINT DEFAULT 0;
    DECLARE max_members_limit INT DEFAULT 0;
    DECLARE max_projects_limit INT DEFAULT 0;
    
    -- 获取当前使用情况
    SELECT storage_used, member_count, max_members, storage_quota, max_projects
    INTO current_storage, current_members, max_members_limit, max_storage, max_projects_limit
    FROM teams
    WHERE id = team_id_param AND deleted_at IS NULL;
    
    -- 计算可用配额
    SET storage_available = max_storage - current_storage;
    SET members_available = max_members_limit - current_members;
    SET projects_available = max_projects_limit - current_projects;
    
    -- 检查是否在限制范围内
    SET within_limits = (
        storage_available >= 0 AND
        members_available >= 0 AND
        projects_available >= 0
    );
    
    -- 返回详细信息
    SELECT 
        current_storage,
        max_storage,
        storage_available,
        current_members,
        max_members_limit,
        members_available,
        current_projects,
        max_projects_limit,
        projects_available,
        within_limits;
        
END//
DELIMITER ;

-- 批量更新团队活跃度存储过程
DELIMITER //
CREATE PROCEDURE UpdateTeamActivityStats()
BEGIN
    DECLARE affected_teams INT DEFAULT 0;
    
    -- 更新活跃成员数（30天内有活动的成员）
    UPDATE teams t
    SET active_member_count = (
        SELECT COUNT(DISTINCT tm.user_id)
        FROM team_members tm
        JOIN users u ON tm.user_id = u.id
        WHERE tm.team_id = t.id
          AND tm.status = 'active'
          AND u.last_login_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
    ),
    last_activity_at = (
        SELECT MAX(GREATEST(
            COALESCE(u.last_login_at, '1970-01-01'),
            COALESCE(tm.updated_at, '1970-01-01')
        ))
        FROM team_members tm
        JOIN users u ON tm.user_id = u.id
        WHERE tm.team_id = t.id AND tm.status = 'active'
    )
    WHERE t.deleted_at IS NULL;
    
    SET affected_teams = ROW_COUNT();
    
    SELECT CONCAT('已更新 ', affected_teams, ' 个团队的活跃度统计') as result;
END//
DELIMITER ;

-- ============================================================================
-- 团队管理视图
-- ============================================================================

-- 团队概览统计视图
CREATE VIEW teams_overview_stats AS
SELECT 
    t.id,
    t.name,
    t.slug,
    t.team_type,
    t.plan_type,
    u.username as owner_name,
    u.email as owner_email,
    t.member_count,
    t.active_member_count,
    t.storage_used,
    t.storage_quota,
    ROUND(t.storage_used * 100.0 / t.storage_quota, 2) as storage_usage_percentage,
    t.status,
    t.is_active,
    t.verification_level,
    t.created_at,
    t.last_activity_at,
    DATEDIFF(NOW(), t.last_activity_at) as days_since_last_activity
FROM teams t
JOIN users u ON t.owner_id = u.id
WHERE t.deleted_at IS NULL;

-- 存储使用情况监控视图
CREATE VIEW teams_storage_monitor AS
SELECT 
    t.id,
    t.name,
    t.plan_type,
    t.storage_used,
    t.storage_quota,
    ROUND(t.storage_used * 100.0 / t.storage_quota, 2) as usage_percentage,
    t.storage_warning_threshold * 100 as warning_threshold_percentage,
    CASE 
        WHEN t.storage_used >= t.storage_quota THEN 'exceeded'
        WHEN t.storage_used >= (t.storage_quota * t.storage_warning_threshold) THEN 'warning'
        ELSE 'normal'
    END as storage_status,
    t.storage_quota - t.storage_used as available_storage,
    u.username as owner_name,
    u.email as owner_email
FROM teams t
JOIN users u ON t.owner_id = u.id
WHERE t.deleted_at IS NULL AND t.is_active = TRUE;

-- 团队活跃度排行视图
CREATE VIEW teams_activity_ranking AS
SELECT 
    t.id,
    t.name,
    t.member_count,
    t.active_member_count,
    ROUND(t.active_member_count * 100.0 / t.member_count, 2) as activity_rate,
    t.total_files_created,
    t.total_shares_created,
    t.last_activity_at,
    DATEDIFF(NOW(), t.created_at) as team_age_days,
    RANK() OVER (ORDER BY t.active_member_count DESC) as activity_rank,
    RANK() OVER (ORDER BY t.total_files_created DESC) as productivity_rank
FROM teams t
WHERE t.deleted_at IS NULL 
  AND t.is_active = TRUE 
  AND t.member_count > 0
ORDER BY activity_rate DESC, t.active_member_count DESC;