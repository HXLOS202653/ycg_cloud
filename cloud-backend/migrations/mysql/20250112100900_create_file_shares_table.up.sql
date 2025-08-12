-- +migrate Up
-- 创建迁移: 文件分享表
-- 版本: 20250112100900
-- 描述: 创建文件分享表，管理文件和文件夹的分享功能
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:09:00
-- 依赖: 20250112100300_create_files_table, 20250112100000_create_users_table
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1
-- 说明: 文件分享表用于管理文件和文件夹的对外分享，支持多种访问控制和权限设置

-- ============================================================================
-- 文件分享表 (file_shares)
-- ============================================================================

CREATE TABLE file_shares (
    -- 基础标识字段
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '分享记录唯一标识',
    
    -- 分享标识和访问控制
    share_code VARCHAR(8) NOT NULL UNIQUE COMMENT '分享码，8位随机字符串，用于访问分享链接',
    share_token VARCHAR(64) DEFAULT NULL COMMENT '分享令牌，用于API访问验证',
    
    -- 分享对象关联（互斥关系）
    file_id BIGINT UNSIGNED DEFAULT NULL COMMENT '分享的单个文件ID，与folder_id互斥',
    folder_id BIGINT UNSIGNED DEFAULT NULL COMMENT '分享的文件夹ID，与file_id互斥',
    shared_items JSON DEFAULT NULL COMMENT '批量分享的文件/文件夹ID列表，格式：{"files": [1,2], "folders": [3,4]}',
    
    -- 分享创建者和基本信息
    user_id BIGINT UNSIGNED NOT NULL COMMENT '分享创建者用户ID',
    share_name VARCHAR(255) NOT NULL COMMENT '分享名称，用户自定义的分享标题',
    share_description TEXT DEFAULT NULL COMMENT '分享描述，详细说明分享内容',
    
    -- 分享类型和模式
    share_type ENUM('file', 'folder', 'multiple') NOT NULL COMMENT '分享类型：单文件、单文件夹、多个项目',
    share_mode ENUM('link', 'email', 'qr_code', 'embed') DEFAULT 'link' COMMENT '分享方式：链接、邮件、二维码、嵌入',
    
    -- 访问控制和安全
    access_type ENUM('public', 'password', 'private', 'whitelist') DEFAULT 'public' COMMENT '访问类型：公开、密码保护、私密、白名单',
    password_hash VARCHAR(255) DEFAULT NULL COMMENT '访问密码的bcrypt哈希值',
    access_whitelist JSON DEFAULT NULL COMMENT '访问白名单，存储允许访问的用户ID或邮箱列表',
    
    -- 权限控制开关
    download_enabled BOOLEAN DEFAULT TRUE COMMENT '是否允许下载文件',
    preview_enabled BOOLEAN DEFAULT TRUE COMMENT '是否允许在线预览',
    comment_enabled BOOLEAN DEFAULT FALSE COMMENT '是否允许访问者添加评论',
    upload_enabled BOOLEAN DEFAULT FALSE COMMENT '是否允许访问者上传文件（文件夹分享）',
    edit_enabled BOOLEAN DEFAULT FALSE COMMENT '是否允许在线编辑（支持的文档类型）',
    print_enabled BOOLEAN DEFAULT TRUE COMMENT '是否允许打印文档',
    watermark_enabled BOOLEAN DEFAULT FALSE COMMENT '是否添加水印保护',
    
    -- 下载限制和统计
    max_downloads INT UNSIGNED DEFAULT NULL COMMENT '最大下载次数限制，NULL表示无限制',
    download_count INT UNSIGNED DEFAULT 0 COMMENT '累计下载次数统计',
    max_views INT UNSIGNED DEFAULT NULL COMMENT '最大查看次数限制，NULL表示无限制',
    view_count INT UNSIGNED DEFAULT 0 COMMENT '累计查看次数统计',
    unique_visitors INT UNSIGNED DEFAULT 0 COMMENT '独立访客数统计',
    
    -- 时间控制
    expires_at TIMESTAMP NULL DEFAULT NULL COMMENT '分享过期时间，NULL表示永久有效',
    auto_delete_after_expire BOOLEAN DEFAULT FALSE COMMENT '过期后是否自动删除分享记录',
    
    -- 分享状态和生命周期
    status ENUM('active', 'expired', 'disabled', 'deleted') DEFAULT 'active' COMMENT '分享状态',
    is_active BOOLEAN DEFAULT TRUE COMMENT '分享是否有效，可手动禁用',
    
    -- 外观和展示设置
    theme ENUM('light', 'dark', 'custom') DEFAULT 'light' COMMENT '分享页面主题',
    custom_css TEXT DEFAULT NULL COMMENT '自定义CSS样式',
    logo_url VARCHAR(1000) DEFAULT NULL COMMENT '自定义Logo URL',
    background_url VARCHAR(1000) DEFAULT NULL COMMENT '背景图片URL',
    
    -- 高级功能设置
    require_email BOOLEAN DEFAULT FALSE COMMENT '是否要求访问者填写邮箱',
    collect_user_info BOOLEAN DEFAULT FALSE COMMENT '是否收集访问者信息',
    notify_on_access BOOLEAN DEFAULT FALSE COMMENT '有人访问时是否通知创建者',
    notify_on_download BOOLEAN DEFAULT FALSE COMMENT '有人下载时是否通知创建者',
    
    -- 嵌入和集成选项
    iframe_enabled BOOLEAN DEFAULT FALSE COMMENT '是否允许iframe嵌入',
    api_access_enabled BOOLEAN DEFAULT FALSE COMMENT '是否允许API访问',
    
    -- 地理位置限制
    geo_restriction_enabled BOOLEAN DEFAULT FALSE COMMENT '是否启用地理位置限制',
    allowed_countries JSON DEFAULT NULL COMMENT '允许访问的国家代码列表',
    blocked_countries JSON DEFAULT NULL COMMENT '禁止访问的国家代码列表',
    
    -- 设备限制
    device_restriction JSON DEFAULT NULL COMMENT '设备访问限制配置',
    
    -- 访问来源控制
    referrer_restriction JSON DEFAULT NULL COMMENT '来源网站限制配置',
    
    -- QR码和短链接
    qr_code_url VARCHAR(1000) DEFAULT NULL COMMENT 'QR码图片URL',
    short_url VARCHAR(100) DEFAULT NULL COMMENT '短链接URL',
    
    -- 统计和分析
    last_accessed_at TIMESTAMP NULL DEFAULT NULL COMMENT '最后访问时间',
    last_downloaded_at TIMESTAMP NULL DEFAULT NULL COMMENT '最后下载时间',
    peak_concurrent_users INT UNSIGNED DEFAULT 0 COMMENT '最高并发访问用户数',
    
    -- 扩展属性
    metadata JSON DEFAULT NULL COMMENT '分享元数据，如来源应用、分享渠道等',
    tags JSON DEFAULT NULL COMMENT '分享标签，便于分类和搜索',
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '分享创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '分享信息更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',
    
    -- 业务索引设计
    INDEX idx_file_shares_share_code (share_code) COMMENT '分享码索引，分享访问查询',
    INDEX idx_file_shares_file_id (file_id) COMMENT '文件ID索引，查询文件的分享',
    INDEX idx_file_shares_folder_id (folder_id) COMMENT '文件夹ID索引，查询文件夹的分享',
    INDEX idx_file_shares_user_id (user_id) COMMENT '用户ID索引，查询用户创建的分享',
    INDEX idx_file_shares_share_type (share_type) COMMENT '分享类型索引',
    INDEX idx_file_shares_access_type (access_type) COMMENT '访问类型索引',
    INDEX idx_file_shares_status (status) COMMENT '分享状态索引',
    INDEX idx_file_shares_is_active (is_active) COMMENT '有效状态索引',
    INDEX idx_file_shares_expires_at (expires_at) COMMENT '过期时间索引，清理任务使用',
    INDEX idx_file_shares_created_at (created_at) COMMENT '创建时间索引，时间排序',
    INDEX idx_file_shares_last_accessed_at (last_accessed_at) COMMENT '最后访问时间索引',
    INDEX idx_file_shares_download_count (download_count) COMMENT '下载次数索引，热门分享统计',
    INDEX idx_file_shares_view_count (view_count) COMMENT '查看次数索引，访问统计',
    INDEX idx_file_shares_deleted_at (deleted_at) COMMENT '软删除索引',
    
    -- 复合业务索引
    INDEX idx_file_shares_user_status (user_id, status, created_at DESC) COMMENT '用户分享状态复合索引',
    INDEX idx_file_shares_type_status (share_type, status, expires_at) COMMENT '类型状态过期时间复合索引',
    INDEX idx_file_shares_active_expires (is_active, expires_at, last_accessed_at) COMMENT '活跃分享过期查询优化',
    
    -- 全文搜索索引
    FULLTEXT idx_fulltext_search (share_name, share_description) COMMENT '分享名称描述全文搜索',
    
    -- 外键约束
    CONSTRAINT fk_file_shares_file_id 
        FOREIGN KEY (file_id) REFERENCES files(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_file_shares_folder_id 
        FOREIGN KEY (folder_id) REFERENCES files(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_file_shares_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE
        
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='文件分享表 - 管理文件和文件夹的对外分享功能，支持多种访问控制'
  ROW_FORMAT=DYNAMIC;

-- ============================================================================
-- 文件分享表约束和检查
-- ============================================================================

-- 分享对象互斥约束：file_id、folder_id、shared_items只能有一个非空
ALTER TABLE file_shares ADD CONSTRAINT chk_share_target_exclusive 
CHECK (
    (file_id IS NOT NULL AND folder_id IS NULL AND shared_items IS NULL) OR
    (file_id IS NULL AND folder_id IS NOT NULL AND shared_items IS NULL) OR
    (file_id IS NULL AND folder_id IS NULL AND shared_items IS NOT NULL)
);

-- 分享类型与对象一致性约束
ALTER TABLE file_shares ADD CONSTRAINT chk_share_type_consistency 
CHECK (
    (share_type = 'file' AND file_id IS NOT NULL) OR
    (share_type = 'folder' AND folder_id IS NOT NULL) OR
    (share_type = 'multiple' AND shared_items IS NOT NULL)
);

-- 密码保护访问类型约束
ALTER TABLE file_shares ADD CONSTRAINT chk_password_access_type 
CHECK (
    (access_type = 'password' AND password_hash IS NOT NULL) OR
    (access_type != 'password')
);

-- 分享码格式约束
ALTER TABLE file_shares ADD CONSTRAINT chk_share_code_format 
CHECK (
    share_code REGEXP '^[A-Za-z0-9]{8}$'
);

-- 下载限制约束
ALTER TABLE file_shares ADD CONSTRAINT chk_download_limits 
CHECK (
    (max_downloads IS NULL OR max_downloads > 0) AND
    download_count >= 0 AND
    (max_downloads IS NULL OR download_count <= max_downloads)
);

-- 查看限制约束
ALTER TABLE file_shares ADD CONSTRAINT chk_view_limits 
CHECK (
    (max_views IS NULL OR max_views > 0) AND
    view_count >= 0 AND
    unique_visitors >= 0 AND
    unique_visitors <= view_count
);

-- 并发用户数约束
ALTER TABLE file_shares ADD CONSTRAINT chk_concurrent_users 
CHECK (peak_concurrent_users >= 0);

-- JSON字段验证
ALTER TABLE file_shares ADD CONSTRAINT chk_file_shares_json_valid 
CHECK (
    (shared_items IS NULL OR JSON_VALID(shared_items)) AND
    (access_whitelist IS NULL OR JSON_VALID(access_whitelist)) AND
    (allowed_countries IS NULL OR JSON_VALID(allowed_countries)) AND
    (blocked_countries IS NULL OR JSON_VALID(blocked_countries)) AND
    (device_restriction IS NULL OR JSON_VALID(device_restriction)) AND
    (referrer_restriction IS NULL OR JSON_VALID(referrer_restriction)) AND
    (metadata IS NULL OR JSON_VALID(metadata)) AND
    (tags IS NULL OR JSON_VALID(tags))
);

-- ============================================================================
-- 文件分享管理触发器
-- ============================================================================

-- 分享码生成触发器
DELIMITER //
CREATE TRIGGER file_shares_generate_code
BEFORE INSERT ON file_shares
FOR EACH ROW
BEGIN
    DECLARE code_exists INT DEFAULT 1;
    DECLARE new_code VARCHAR(8);
    
    -- 如果没有提供分享码，自动生成
    IF NEW.share_code IS NULL OR NEW.share_code = '' THEN
        WHILE code_exists > 0 DO
            -- 生成8位随机分享码（字母数字混合）
            SET new_code = UPPER(SUBSTRING(MD5(CONCAT(UNIX_TIMESTAMP(), RAND())), 1, 8));
            
            -- 检查分享码是否已存在
            SELECT COUNT(*) INTO code_exists 
            FROM file_shares 
            WHERE share_code = new_code;
        END WHILE;
        
        SET NEW.share_code = new_code;
    END IF;
    
    -- 自动生成分享令牌
    IF NEW.share_token IS NULL THEN
        SET NEW.share_token = SHA2(CONCAT(NEW.share_code, UNIX_TIMESTAMP(), RAND()), 256);
    END IF;
END//
DELIMITER ;

-- 分享访问统计更新触发器
DELIMITER //
CREATE TRIGGER file_shares_update_stats
BEFORE UPDATE ON file_shares
FOR EACH ROW
BEGIN
    -- 更新最后访问时间
    IF NEW.view_count > OLD.view_count THEN
        SET NEW.last_accessed_at = CURRENT_TIMESTAMP;
    END IF;
    
    -- 更新最后下载时间
    IF NEW.download_count > OLD.download_count THEN
        SET NEW.last_downloaded_at = CURRENT_TIMESTAMP;
    END IF;
    
    -- 自动更新分享状态
    IF NEW.expires_at IS NOT NULL AND NEW.expires_at <= CURRENT_TIMESTAMP THEN
        SET NEW.status = 'expired';
        SET NEW.is_active = FALSE;
    END IF;
END//
DELIMITER ;

-- 分享文件删除检查触发器
DELIMITER //
CREATE TRIGGER file_shares_check_file_deleted
BEFORE UPDATE ON files
FOR EACH ROW
BEGIN
    -- 当文件被软删除时，自动禁用相关分享
    IF OLD.is_deleted = FALSE AND NEW.is_deleted = TRUE THEN
        UPDATE file_shares 
        SET status = 'disabled', 
            is_active = FALSE,
            updated_at = CURRENT_TIMESTAMP
        WHERE (file_id = NEW.id OR folder_id = NEW.id) 
          AND status = 'active';
    END IF;
    
    -- 当文件恢复时，可以选择性地重新启用分享（需要手动操作）
END//
DELIMITER ;

-- ============================================================================
-- 分享管理存储过程
-- ============================================================================

-- 批量清理过期分享的存储过程
DELIMITER //
CREATE PROCEDURE CleanExpiredShares()
BEGIN
    DECLARE affected_rows INT DEFAULT 0;
    
    -- 更新过期分享状态
    UPDATE file_shares 
    SET status = 'expired', 
        is_active = FALSE,
        updated_at = CURRENT_TIMESTAMP
    WHERE expires_at IS NOT NULL 
      AND expires_at <= CURRENT_TIMESTAMP 
      AND status = 'active';
    
    SET affected_rows = ROW_COUNT();
    
    -- 删除设置了自动删除的过期分享
    DELETE FROM file_shares 
    WHERE status = 'expired' 
      AND auto_delete_after_expire = TRUE 
      AND expires_at <= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL 1 DAY);
    
    SELECT CONCAT('已处理 ', affected_rows, ' 个过期分享，删除 ', ROW_COUNT(), ' 个自动删除分享') as result;
END//
DELIMITER ;

-- 获取分享统计信息的存储过程
DELIMITER //
CREATE PROCEDURE GetShareStatistics(
    IN user_id_param BIGINT UNSIGNED,
    IN days_param INT DEFAULT 30
)
BEGIN
    SELECT 
        COUNT(*) as total_shares,
        COUNT(CASE WHEN status = 'active' THEN 1 END) as active_shares,
        COUNT(CASE WHEN status = 'expired' THEN 1 END) as expired_shares,
        SUM(download_count) as total_downloads,
        SUM(view_count) as total_views,
        SUM(unique_visitors) as total_visitors,
        AVG(download_count) as avg_downloads_per_share,
        AVG(view_count) as avg_views_per_share,
        COUNT(CASE WHEN created_at >= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL days_param DAY) THEN 1 END) as recent_shares
    FROM file_shares
    WHERE user_id = user_id_param
      AND deleted_at IS NULL;
END//
DELIMITER ;

-- 验证分享访问权限的存储过程
DELIMITER //
CREATE PROCEDURE ValidateShareAccess(
    IN share_code_param VARCHAR(8),
    IN access_password VARCHAR(255),
    IN visitor_ip VARCHAR(45),
    IN visitor_country VARCHAR(2),
    OUT access_granted BOOLEAN,
    OUT error_message VARCHAR(500)
)
BEGIN
    DECLARE share_exists INT DEFAULT 0;
    DECLARE share_active BOOLEAN DEFAULT FALSE;
    DECLARE share_expired BOOLEAN DEFAULT FALSE;
    DECLARE password_required BOOLEAN DEFAULT FALSE;
    DECLARE password_correct BOOLEAN DEFAULT FALSE;
    DECLARE geo_restricted BOOLEAN DEFAULT FALSE;
    DECLARE country_allowed BOOLEAN DEFAULT TRUE;
    
    SET access_granted = FALSE;
    SET error_message = NULL;
    
    -- 检查分享是否存在
    SELECT COUNT(*), 
           status = 'active',
           expires_at IS NOT NULL AND expires_at <= CURRENT_TIMESTAMP,
           access_type = 'password',
           geo_restriction_enabled
    INTO share_exists, share_active, share_expired, password_required, geo_restricted
    FROM file_shares
    WHERE share_code = share_code_param AND deleted_at IS NULL;
    
    -- 分享不存在
    IF share_exists = 0 THEN
        SET error_message = '分享不存在或已被删除';
        LEAVE ValidateShareAccess;
    END IF;
    
    -- 分享已过期
    IF share_expired THEN
        SET error_message = '分享已过期';
        LEAVE ValidateShareAccess;
    END IF;
    
    -- 分享未激活
    IF NOT share_active THEN
        SET error_message = '分享已被禁用';
        LEAVE ValidateShareAccess;
    END IF;
    
    -- 检查密码（如果需要）
    IF password_required THEN
        SELECT password_hash = SHA2(CONCAT('share_password_salt_', access_password), 256)
        INTO password_correct
        FROM file_shares
        WHERE share_code = share_code_param;
        
        IF NOT password_correct THEN
            SET error_message = '访问密码错误';
            LEAVE ValidateShareAccess;
        END IF;
    END IF;
    
    -- 检查地理位置限制（如果启用）
    IF geo_restricted AND visitor_country IS NOT NULL THEN
        -- 这里需要根据allowed_countries和blocked_countries进行检查
        -- 简化实现，实际应用中需要解析JSON字段
        SET country_allowed = TRUE; -- 简化处理
        
        IF NOT country_allowed THEN
            SET error_message = '您所在的地区无法访问此分享';
            LEAVE ValidateShareAccess;
        END IF;
    END IF;
    
    -- 所有检查通过
    SET access_granted = TRUE;
    
    -- 更新访问统计
    UPDATE file_shares 
    SET view_count = view_count + 1,
        last_accessed_at = CURRENT_TIMESTAMP
    WHERE share_code = share_code_param;
    
END//
DELIMITER ;

-- ============================================================================
-- 分享管理视图
-- ============================================================================

-- 活跃分享统计视图
CREATE VIEW active_shares_stats AS
SELECT 
    fs.user_id,
    u.username,
    COUNT(*) as total_active_shares,
    COUNT(CASE WHEN fs.share_type = 'file' THEN 1 END) as file_shares,
    COUNT(CASE WHEN fs.share_type = 'folder' THEN 1 END) as folder_shares,
    COUNT(CASE WHEN fs.share_type = 'multiple' THEN 1 END) as multiple_shares,
    SUM(fs.download_count) as total_downloads,
    SUM(fs.view_count) as total_views,
    MAX(fs.last_accessed_at) as last_share_access,
    AVG(fs.download_count) as avg_downloads_per_share
FROM file_shares fs
JOIN users u ON fs.user_id = u.id
WHERE fs.status = 'active' 
  AND fs.is_active = TRUE 
  AND fs.deleted_at IS NULL
GROUP BY fs.user_id, u.username;

-- 热门分享排行视图
CREATE VIEW popular_shares_ranking AS
SELECT 
    fs.id,
    fs.share_code,
    fs.share_name,
    fs.share_type,
    u.username as creator,
    fs.download_count,
    fs.view_count,
    fs.unique_visitors,
    fs.created_at,
    fs.last_accessed_at,
    RANK() OVER (ORDER BY fs.view_count DESC) as view_rank,
    RANK() OVER (ORDER BY fs.download_count DESC) as download_rank,
    RANK() OVER (ORDER BY fs.unique_visitors DESC) as visitor_rank
FROM file_shares fs
JOIN users u ON fs.user_id = u.id
WHERE fs.status = 'active' 
  AND fs.is_active = TRUE 
  AND fs.deleted_at IS NULL
  AND fs.view_count > 0
ORDER BY fs.view_count DESC;

-- 即将过期分享提醒视图
CREATE VIEW expiring_shares_alert AS
SELECT 
    fs.id,
    fs.share_code,
    fs.share_name,
    fs.user_id,
    u.username,
    u.email,
    fs.expires_at,
    TIMESTAMPDIFF(HOUR, CURRENT_TIMESTAMP, fs.expires_at) as hours_until_expire
FROM file_shares fs
JOIN users u ON fs.user_id = u.id
WHERE fs.status = 'active' 
  AND fs.is_active = TRUE 
  AND fs.expires_at IS NOT NULL
  AND fs.expires_at > CURRENT_TIMESTAMP
  AND fs.expires_at <= DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 24 HOUR)
ORDER BY fs.expires_at ASC;