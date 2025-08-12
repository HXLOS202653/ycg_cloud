-- +migrate Up
-- 创建迁移: 文件权限表
-- 版本: 20250112101100
-- 描述: 创建文件权限表，管理文件和文件夹的细粒度权限控制
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:11:00
-- 依赖: 20250112100300_create_files_table, 20250112100000_create_users_table
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1
-- 说明: 文件权限表实现细粒度的访问控制，支持用户、团队和公开权限

-- ============================================================================
-- 文件权限表 (file_permissions)
-- ============================================================================

CREATE TABLE file_permissions (
    -- 基础标识字段
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '权限记录唯一标识',
    
    -- 权限目标对象
    file_id BIGINT UNSIGNED NOT NULL COMMENT '关联的文件或文件夹ID',
    
    -- 权限授予对象（互斥关系）
    user_id BIGINT UNSIGNED DEFAULT NULL COMMENT '被授权的用户ID，与team_id互斥',
    team_id BIGINT UNSIGNED DEFAULT NULL COMMENT '被授权的团队ID，与user_id互斥',
    role_name VARCHAR(50) DEFAULT NULL COMMENT '角色名称，用于基于角色的权限控制',
    
    -- 权限分类和类型
    permission_type ENUM('user', 'team', 'role', 'public', 'inherit') NOT NULL COMMENT '权限类型',
    permission_level ENUM('owner', 'editor', 'viewer', 'uploader', 'commenter', 'custom') DEFAULT 'viewer' COMMENT '权限级别',
    
    -- 详细权限控制（JSON格式）
    permissions JSON NOT NULL COMMENT '具体权限配置',
    
    -- 权限来源和管理
    granted_by BIGINT UNSIGNED NOT NULL COMMENT '权限授予者用户ID',
    permission_source ENUM('direct', 'inherited', 'team', 'share', 'system') DEFAULT 'direct' COMMENT '权限来源',
    inherit_from_parent BOOLEAN DEFAULT TRUE COMMENT '是否从父文件夹继承权限',
    
    -- 权限生效和过期
    effective_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '权限生效时间',
    expires_at TIMESTAMP NULL DEFAULT NULL COMMENT '权限过期时间，NULL表示永久有效',
    auto_extend BOOLEAN DEFAULT FALSE COMMENT '是否自动延期（临时权限）',
    
    -- 权限状态和控制
    is_active BOOLEAN DEFAULT TRUE COMMENT '权限是否有效',
    is_temporary BOOLEAN DEFAULT FALSE COMMENT '是否为临时权限',
    can_delegate BOOLEAN DEFAULT FALSE COMMENT '是否可以将权限委托给其他用户',
    
    -- 访问限制条件
    access_conditions JSON DEFAULT NULL COMMENT '访问条件限制，如IP范围、时间段等',
    ip_whitelist JSON DEFAULT NULL COMMENT 'IP白名单',
    ip_blacklist JSON DEFAULT NULL COMMENT 'IP黑名单',
    time_restrictions JSON DEFAULT NULL COMMENT '时间访问限制',
    device_restrictions JSON DEFAULT NULL COMMENT '设备访问限制',
    
    -- 权限使用统计
    usage_count INT UNSIGNED DEFAULT 0 COMMENT '权限使用次数',
    last_used_at TIMESTAMP NULL DEFAULT NULL COMMENT '最后使用时间',
    max_usage_count INT UNSIGNED DEFAULT NULL COMMENT '最大使用次数限制',
    
    -- 审批和工作流
    approval_required BOOLEAN DEFAULT FALSE COMMENT '是否需要审批才能生效',
    approved_by BIGINT UNSIGNED DEFAULT NULL COMMENT '审批者用户ID',
    approved_at TIMESTAMP NULL DEFAULT NULL COMMENT '审批时间',
    approval_notes TEXT DEFAULT NULL COMMENT '审批备注',
    
    -- 权限描述和备注
    permission_name VARCHAR(200) DEFAULT NULL COMMENT '权限名称或标题',
    description TEXT DEFAULT NULL COMMENT '权限描述和说明',
    internal_notes TEXT DEFAULT NULL COMMENT '内部备注（仅管理员可见）',
    
    -- 权限继承链
    parent_permission_id BIGINT UNSIGNED DEFAULT NULL COMMENT '父权限ID，用于权限继承链',
    inheritance_depth INT UNSIGNED DEFAULT 0 COMMENT '继承深度级别',
    
    -- 扩展属性
    tags JSON DEFAULT NULL COMMENT '权限标签，便于分类管理',
    metadata JSON DEFAULT NULL COMMENT '扩展元数据',
    custom_fields JSON DEFAULT NULL COMMENT '自定义字段',
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '权限创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '权限更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',
    
    -- 业务索引设计
    INDEX idx_file_permissions_file_id (file_id) COMMENT '文件ID索引，查询文件的权限',
    INDEX idx_file_permissions_user_id (user_id) COMMENT '用户ID索引，查询用户权限',
    INDEX idx_file_permissions_team_id (team_id) COMMENT '团队ID索引，查询团队权限',
    INDEX idx_file_permissions_permission_type (permission_type) COMMENT '权限类型索引',
    INDEX idx_file_permissions_permission_level (permission_level) COMMENT '权限级别索引',
    INDEX idx_file_permissions_granted_by (granted_by) COMMENT '授权者索引',
    INDEX idx_file_permissions_permission_source (permission_source) COMMENT '权限来源索引',
    INDEX idx_file_permissions_is_active (is_active) COMMENT '有效状态索引',
    INDEX idx_file_permissions_expires_at (expires_at) COMMENT '过期时间索引，清理任务使用',
    INDEX idx_file_permissions_effective_at (effective_at) COMMENT '生效时间索引',
    INDEX idx_file_permissions_last_used_at (last_used_at) COMMENT '最后使用时间索引',
    INDEX idx_file_permissions_approval_required (approval_required) COMMENT '审批状态索引',
    INDEX idx_file_permissions_parent_permission_id (parent_permission_id) COMMENT '父权限索引',
    INDEX idx_file_permissions_deleted_at (deleted_at) COMMENT '软删除索引',
    
    -- 复合业务索引
    INDEX idx_file_permissions_file_user_active (file_id, user_id, is_active, expires_at) COMMENT '文件用户权限查询优化',
    INDEX idx_file_permissions_file_team_active (file_id, team_id, is_active, expires_at) COMMENT '文件团队权限查询优化',
    INDEX idx_file_permissions_user_file_permissions (user_id, file_id, permission_type, is_active) COMMENT '用户文件权限查询',
    INDEX idx_file_permissions_type_level_active (permission_type, permission_level, is_active) COMMENT '权限类型级别查询',
    INDEX idx_file_permissions_granted_by_time (granted_by, created_at DESC) COMMENT '授权者时间查询',
    INDEX idx_file_permissions_expiring_permissions (expires_at, is_active, approval_required) COMMENT '即将过期权限查询',
    
    -- 权限继承链索引
    INDEX idx_file_permissions_inheritance_chain (parent_permission_id, inheritance_depth, is_active) COMMENT '权限继承链查询',
    
    -- 全文搜索索引
    FULLTEXT idx_fulltext_search (permission_name, description) COMMENT '权限名称描述搜索',
    
    -- 外键约束
    CONSTRAINT fk_file_permissions_file_id 
        FOREIGN KEY (file_id) REFERENCES files(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_file_permissions_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_file_permissions_granted_by 
        FOREIGN KEY (granted_by) REFERENCES users(id) 
        ON DELETE RESTRICT 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_file_permissions_approved_by 
        FOREIGN KEY (approved_by) REFERENCES users(id) 
        ON DELETE SET NULL 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_file_permissions_parent 
        FOREIGN KEY (parent_permission_id) REFERENCES file_permissions(id) 
        ON DELETE SET NULL 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_file_permissions_team_id 
        FOREIGN KEY (team_id) REFERENCES teams(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE
        
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='文件权限表 - 管理文件和文件夹的细粒度权限控制，支持继承和委托'
  ROW_FORMAT=DYNAMIC;

-- ============================================================================
-- 文件权限表约束和检查
-- ============================================================================

-- 权限授予对象互斥约束：user_id、team_id、role_name只能有一个非空（public类型除外）
ALTER TABLE file_permissions ADD CONSTRAINT chk_permission_target_exclusive 
CHECK (
    (permission_type = 'public' AND user_id IS NULL AND team_id IS NULL) OR
    (permission_type = 'user' AND user_id IS NOT NULL AND team_id IS NULL AND role_name IS NULL) OR
    (permission_type = 'team' AND user_id IS NULL AND team_id IS NOT NULL AND role_name IS NULL) OR
    (permission_type = 'role' AND user_id IS NULL AND team_id IS NULL AND role_name IS NOT NULL) OR
    (permission_type = 'inherit')
);

-- 权限时间逻辑约束
ALTER TABLE file_permissions ADD CONSTRAINT chk_permission_time_logic 
CHECK (
    effective_at <= COALESCE(expires_at, '2099-12-31 23:59:59')
);

-- 使用次数约束
ALTER TABLE file_permissions ADD CONSTRAINT chk_usage_count_limits 
CHECK (
    usage_count >= 0 AND
    (max_usage_count IS NULL OR max_usage_count > 0) AND
    (max_usage_count IS NULL OR usage_count <= max_usage_count)
);

-- 继承深度约束
ALTER TABLE file_permissions ADD CONSTRAINT chk_inheritance_depth 
CHECK (inheritance_depth >= 0 AND inheritance_depth <= 10);

-- 审批逻辑约束
ALTER TABLE file_permissions ADD CONSTRAINT chk_approval_logic 
CHECK (
    (approval_required = FALSE) OR
    (approval_required = TRUE AND approved_by IS NOT NULL AND approved_at IS NOT NULL)
);

-- 角色名称格式约束
ALTER TABLE file_permissions ADD CONSTRAINT chk_role_name_format 
CHECK (
    role_name IS NULL OR 
    (role_name REGEXP '^[a-zA-Z][a-zA-Z0-9_-]*$' AND LENGTH(role_name) <= 50)
);

-- JSON字段验证
ALTER TABLE file_permissions ADD CONSTRAINT chk_file_permissions_json_valid 
CHECK (
    JSON_VALID(permissions) AND
    (access_conditions IS NULL OR JSON_VALID(access_conditions)) AND
    (ip_whitelist IS NULL OR JSON_VALID(ip_whitelist)) AND
    (ip_blacklist IS NULL OR JSON_VALID(ip_blacklist)) AND
    (time_restrictions IS NULL OR JSON_VALID(time_restrictions)) AND
    (device_restrictions IS NULL OR JSON_VALID(device_restrictions)) AND
    (tags IS NULL OR JSON_VALID(tags)) AND
    (metadata IS NULL OR JSON_VALID(metadata)) AND
    (custom_fields IS NULL OR JSON_VALID(custom_fields))
);

-- 权限JSON结构验证
ALTER TABLE file_permissions ADD CONSTRAINT chk_permissions_structure 
CHECK (
    JSON_CONTAINS_PATH(permissions, 'one', '$.read') OR
    JSON_CONTAINS_PATH(permissions, 'one', '$.write') OR
    JSON_CONTAINS_PATH(permissions, 'one', '$.delete') OR
    JSON_CONTAINS_PATH(permissions, 'one', '$.share')
);

-- ============================================================================
-- 文件权限管理触发器
-- ============================================================================

-- 权限生效检查触发器
DELIMITER //
CREATE TRIGGER file_permissions_activation_check
BEFORE UPDATE ON file_permissions
FOR EACH ROW
BEGIN
    -- 检查权限是否过期
    IF NEW.expires_at IS NOT NULL AND NEW.expires_at <= CURRENT_TIMESTAMP THEN
        SET NEW.is_active = FALSE;
    END IF;
    
    -- 检查权限是否已到生效时间
    IF NEW.effective_at > CURRENT_TIMESTAMP THEN
        SET NEW.is_active = FALSE;
    END IF;
    
    -- 检查使用次数限制
    IF NEW.max_usage_count IS NOT NULL AND NEW.usage_count >= NEW.max_usage_count THEN
        SET NEW.is_active = FALSE;
    END IF;
    
    -- 更新最后使用时间（如果使用次数增加）
    IF NEW.usage_count > OLD.usage_count THEN
        SET NEW.last_used_at = CURRENT_TIMESTAMP;
    END IF;
END//
DELIMITER ;

-- 权限继承管理触发器
DELIMITER //
CREATE TRIGGER file_permissions_inheritance_update
AFTER INSERT ON file_permissions
FOR EACH ROW
BEGIN
    -- 如果是启用继承的权限，自动为子文件/文件夹创建继承权限
    IF NEW.inherit_from_parent = TRUE AND NEW.permission_type IN ('user', 'team') THEN
        INSERT INTO file_permissions (
            file_id, user_id, team_id, permission_type, permissions,
            granted_by, permission_source, inherit_from_parent,
            parent_permission_id, inheritance_depth, is_active
        )
        SELECT 
            f.id, NEW.user_id, NEW.team_id, NEW.permission_type, NEW.permissions,
            NEW.granted_by, 'inherited', TRUE,
            NEW.id, NEW.inheritance_depth + 1, NEW.is_active
        FROM files f
        WHERE f.parent_id = NEW.file_id 
          AND f.is_deleted = FALSE
          AND NEW.inheritance_depth < 5; -- 限制继承深度防止无限递归
    END IF;
END//
DELIMITER ;

-- 权限冲突检测触发器
DELIMITER //
CREATE TRIGGER file_permissions_conflict_check
BEFORE INSERT ON file_permissions
FOR EACH ROW
BEGIN
    DECLARE existing_count INT DEFAULT 0;
    
    -- 检查是否存在相同的权限配置（防止重复）
    SELECT COUNT(*) INTO existing_count
    FROM file_permissions
    WHERE file_id = NEW.file_id
      AND permission_type = NEW.permission_type
      AND (
          (NEW.user_id IS NOT NULL AND user_id = NEW.user_id) OR
          (NEW.team_id IS NOT NULL AND team_id = NEW.team_id) OR
          (NEW.role_name IS NOT NULL AND role_name = NEW.role_name)
      )
      AND is_active = TRUE
      AND deleted_at IS NULL;
    
    -- 如果存在重复权限，标记为需要审批
    IF existing_count > 0 THEN
        SET NEW.approval_required = TRUE;
        SET NEW.is_active = FALSE;
    END IF;
END//
DELIMITER ;

-- ============================================================================
-- 权限管理存储过程
-- ============================================================================

-- 检查用户对文件的权限
DELIMITER //
CREATE PROCEDURE CheckUserFilePermission(
    IN user_id_param BIGINT UNSIGNED,
    IN file_id_param BIGINT UNSIGNED,
    IN permission_action VARCHAR(20),
    OUT has_permission BOOLEAN,
    OUT permission_source VARCHAR(50)
)
BEGIN
    DECLARE direct_permission BOOLEAN DEFAULT FALSE;
    DECLARE inherited_permission BOOLEAN DEFAULT FALSE;
    DECLARE team_permission BOOLEAN DEFAULT FALSE;
    DECLARE public_permission BOOLEAN DEFAULT FALSE;
    DECLARE file_owner_id BIGINT UNSIGNED;
    
    SET has_permission = FALSE;
    SET permission_source = 'none';
    
    -- 检查文件是否存在和获取所有者
    SELECT user_id INTO file_owner_id
    FROM files
    WHERE id = file_id_param AND is_deleted = FALSE;
    
    -- 如果文件不存在
    IF file_owner_id IS NULL THEN
        LEAVE CheckUserFilePermission;
    END IF;
    
    -- 1. 检查是否为文件所有者
    IF file_owner_id = user_id_param THEN
        SET has_permission = TRUE;
        SET permission_source = 'owner';
        LEAVE CheckUserFilePermission;
    END IF;
    
    -- 2. 检查直接用户权限
    SELECT JSON_EXTRACT(permissions, CONCAT('$.', permission_action)) = 'true'
    INTO direct_permission
    FROM file_permissions
    WHERE file_id = file_id_param
      AND user_id = user_id_param
      AND permission_type = 'user'
      AND is_active = TRUE
      AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
      AND effective_at <= CURRENT_TIMESTAMP
      AND deleted_at IS NULL
    LIMIT 1;
    
    IF direct_permission = TRUE THEN
        SET has_permission = TRUE;
        SET permission_source = 'direct';
        LEAVE CheckUserFilePermission;
    END IF;
    
    -- 3. 检查团队权限（这里简化处理，实际需要查询用户的团队成员关系）
    SELECT JSON_EXTRACT(fp.permissions, CONCAT('$.', permission_action)) = 'true'
    INTO team_permission
    FROM file_permissions fp
    WHERE fp.file_id = file_id_param
      AND fp.permission_type = 'team'
      AND fp.is_active = TRUE
      AND (fp.expires_at IS NULL OR fp.expires_at > CURRENT_TIMESTAMP)
      AND fp.effective_at <= CURRENT_TIMESTAMP
      AND fp.deleted_at IS NULL
    LIMIT 1;
    
    IF team_permission = TRUE THEN
        SET has_permission = TRUE;
        SET permission_source = 'team';
        LEAVE CheckUserFilePermission;
    END IF;
    
    -- 4. 检查公开权限
    SELECT JSON_EXTRACT(permissions, CONCAT('$.', permission_action)) = 'true'
    INTO public_permission
    FROM file_permissions
    WHERE file_id = file_id_param
      AND permission_type = 'public'
      AND is_active = TRUE
      AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
      AND effective_at <= CURRENT_TIMESTAMP
      AND deleted_at IS NULL
    LIMIT 1;
    
    IF public_permission = TRUE THEN
        SET has_permission = TRUE;
        SET permission_source = 'public';
        LEAVE CheckUserFilePermission;
    END IF;
    
    -- 5. 检查继承权限（简化实现，实际需要递归查询父文件夹）
    -- TODO: 实现完整的权限继承逻辑
    
END//
DELIMITER ;

-- 批量授权存储过程
DELIMITER //
CREATE PROCEDURE BatchGrantPermissions(
    IN file_ids JSON,
    IN target_user_id BIGINT UNSIGNED,
    IN target_team_id BIGINT UNSIGNED,
    IN permission_config JSON,
    IN granted_by_user_id BIGINT UNSIGNED,
    OUT affected_count INT
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE file_count INT DEFAULT 0;
    DECLARE current_file_id BIGINT UNSIGNED;
    
    SET affected_count = 0;
    SET file_count = JSON_LENGTH(file_ids);
    
    -- 遍历文件ID列表
    WHILE i < file_count DO
        SET current_file_id = JSON_UNQUOTE(JSON_EXTRACT(file_ids, CONCAT('$[', i, ']')));
        
        -- 插入权限记录
        INSERT INTO file_permissions (
            file_id, user_id, team_id, permission_type, permissions, granted_by
        ) VALUES (
            current_file_id,
            target_user_id,
            target_team_id,
            CASE WHEN target_user_id IS NOT NULL THEN 'user' ELSE 'team' END,
            permission_config,
            granted_by_user_id
        )
        ON DUPLICATE KEY UPDATE
            permissions = VALUES(permissions),
            updated_at = CURRENT_TIMESTAMP;
        
        SET affected_count = affected_count + 1;
        SET i = i + 1;
    END WHILE;
    
END//
DELIMITER ;

-- 清理过期权限存储过程
DELIMITER //
CREATE PROCEDURE CleanExpiredPermissions()
BEGIN
    DECLARE affected_rows INT DEFAULT 0;
    
    -- 标记过期权限为无效
    UPDATE file_permissions 
    SET is_active = FALSE,
        updated_at = CURRENT_TIMESTAMP
    WHERE expires_at IS NOT NULL 
      AND expires_at <= CURRENT_TIMESTAMP 
      AND is_active = TRUE;
    
    SET affected_rows = ROW_COUNT();
    
    -- 删除30天前的已过期临时权限
    DELETE FROM file_permissions 
    WHERE is_temporary = TRUE 
      AND is_active = FALSE 
      AND expires_at <= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL 30 DAY);
    
    SELECT CONCAT('已处理 ', affected_rows, ' 个过期权限，删除 ', ROW_COUNT(), ' 个临时权限') as result;
END//
DELIMITER ;

-- ============================================================================
-- 权限管理视图
-- ============================================================================

-- 用户权限概览视图
CREATE VIEW user_permissions_overview AS
SELECT 
    u.id as user_id,
    u.username,
    COUNT(fp.id) as total_permissions,
    COUNT(CASE WHEN fp.permission_level = 'owner' THEN 1 END) as owner_permissions,
    COUNT(CASE WHEN fp.permission_level = 'editor' THEN 1 END) as editor_permissions,
    COUNT(CASE WHEN fp.permission_level = 'viewer' THEN 1 END) as viewer_permissions,
    COUNT(CASE WHEN fp.is_temporary = TRUE THEN 1 END) as temporary_permissions,
    COUNT(CASE WHEN fp.expires_at IS NOT NULL AND fp.expires_at <= DATE_ADD(NOW(), INTERVAL 7 DAY) THEN 1 END) as expiring_soon,
    MAX(fp.last_used_at) as last_permission_used
FROM users u
LEFT JOIN file_permissions fp ON u.id = fp.user_id AND fp.is_active = TRUE AND fp.deleted_at IS NULL
GROUP BY u.id, u.username;

-- 文件权限详情视图
CREATE VIEW file_permissions_detail AS
SELECT 
    f.id as file_id,
    f.filename,
    f.file_path,
    fo.username as file_owner,
    fp.id as permission_id,
    fp.permission_type,
    fp.permission_level,
    CASE 
        WHEN fp.user_id IS NOT NULL THEN u.username
        WHEN fp.team_id IS NOT NULL THEN CONCAT('Team:', fp.team_id)
        WHEN fp.role_name IS NOT NULL THEN CONCAT('Role:', fp.role_name)
        ELSE 'Public'
    END as permission_target,
    fp.permissions,
    fp.is_active,
    fp.expires_at,
    gb.username as granted_by_user,
    fp.created_at
FROM files f
JOIN users fo ON f.user_id = fo.id
LEFT JOIN file_permissions fp ON f.id = fp.file_id AND fp.deleted_at IS NULL
LEFT JOIN users u ON fp.user_id = u.id
LEFT JOIN users gb ON fp.granted_by = gb.id
WHERE f.is_deleted = FALSE;

-- 即将过期权限提醒视图
CREATE VIEW expiring_permissions_alert AS
SELECT 
    fp.id,
    f.filename,
    f.file_path,
    CASE 
        WHEN fp.user_id IS NOT NULL THEN u.username
        WHEN fp.team_id IS NOT NULL THEN CONCAT('Team:', fp.team_id)
        ELSE 'Public'
    END as permission_holder,
    u.email as user_email,
    fp.permission_level,
    fp.expires_at,
    TIMESTAMPDIFF(HOUR, NOW(), fp.expires_at) as hours_until_expire,
    gb.username as granted_by
FROM file_permissions fp
JOIN files f ON fp.file_id = f.id
LEFT JOIN users u ON fp.user_id = u.id
JOIN users gb ON fp.granted_by = gb.id
WHERE fp.is_active = TRUE
  AND fp.expires_at IS NOT NULL
  AND fp.expires_at > NOW()
  AND fp.expires_at <= DATE_ADD(NOW(), INTERVAL 168 HOUR) -- 7天内过期
ORDER BY fp.expires_at ASC;
