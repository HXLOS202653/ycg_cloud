-- +migrate Up
-- 创建迁移: 用户会话表
-- 版本: 20250112100100
-- 描述: 创建用户会话管理表，支持多设备登录、JWT token管理和会话安全控制
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:01:00
-- 依赖: 20250112100000_create_users_table
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1

-- ============================================================================
-- 用户会话管理表 (user_sessions)
-- ============================================================================

CREATE TABLE user_sessions (
    -- 基础标识字段
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '会话ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID，关联users表',
    
    -- JWT Token 管理
    session_token VARCHAR(128) NOT NULL UNIQUE COMMENT '会话令牌，JWT Access Token的JTI（唯一标识符）',
    refresh_token_hash VARCHAR(255) NOT NULL COMMENT 'Refresh Token的SHA256哈希值，确保安全存储',
    
    -- 设备和环境信息
    device_info JSON DEFAULT NULL COMMENT '设备信息，包含设备类型、操作系统、浏览器等',
    user_agent TEXT DEFAULT NULL COMMENT '完整的用户代理字符串',
    ip_address VARCHAR(45) NOT NULL COMMENT '登录IP地址，支持IPv4和IPv6',
    location VARCHAR(100) DEFAULT NULL COMMENT '登录地理位置，格式：城市,省份,国家',
    
    -- 会话状态和安全
    is_active BOOLEAN DEFAULT TRUE COMMENT '会话是否有效，false表示已注销或被强制下线',
    session_type ENUM('web', 'mobile', 'desktop', 'api') DEFAULT 'web' COMMENT '会话类型：网页、移动端、桌面端、API',
    
    -- 时间管理
    expires_at TIMESTAMP NOT NULL COMMENT '会话过期时间，基于JWT的exp声明',
    last_activity_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '最后活动时间，用于判断空闲超时',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '会话创建时间',
    
    -- 安全审计字段
    login_method ENUM('password', 'oauth', 'sso', 'api_key', 'two_factor') DEFAULT 'password' COMMENT '登录方式',
    risk_score INT UNSIGNED DEFAULT 0 COMMENT '风险评分，0-100，用于异常检测',
    country_code VARCHAR(2) DEFAULT NULL COMMENT '登录国家代码，ISO 3166-1 alpha-2',
    
    -- 外键约束
    CONSTRAINT fk_user_sessions_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
    
    -- 索引设计
    INDEX idx_user_sessions_user_id (user_id) COMMENT '用户ID索引，用于查询用户的所有会话',
    INDEX idx_user_sessions_session_token (session_token) COMMENT '会话令牌索引，用于token验证',
    INDEX idx_user_sessions_expires_at (expires_at) COMMENT '过期时间索引，用于清理过期会话',
    INDEX idx_user_sessions_is_active (is_active) COMMENT '活跃状态索引，用于筛选有效会话',
    INDEX idx_user_sessions_last_activity_at (last_activity_at) COMMENT '最后活动时间索引，用于空闲检测',
    INDEX idx_user_sessions_ip_address (ip_address) COMMENT 'IP地址索引，用于安全分析',
    INDEX idx_user_sessions_user_active (user_id, is_active, expires_at) COMMENT '复合索引，用于查询用户的有效会话',
    INDEX idx_user_sessions_created_at (created_at) COMMENT '创建时间索引，用于统计和分析'
    
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='用户会话管理表 - 支持JWT token管理、多设备登录和会话安全控制'
  ROW_FORMAT=DYNAMIC;

-- ============================================================================
-- 会话管理触发器和约束
-- ============================================================================

-- 创建会话过期检查触发器
DELIMITER //
CREATE TRIGGER user_sessions_expires_check 
BEFORE INSERT ON user_sessions
FOR EACH ROW
BEGIN
    -- 检查过期时间必须大于当前时间
    IF NEW.expires_at <= CURRENT_TIMESTAMP THEN
        SIGNAL SQLSTATE '45000' 
        SET MESSAGE_TEXT = '会话过期时间必须大于当前时间';
    END IF;
    
    -- 设置默认的最后活动时间
    IF NEW.last_activity_at IS NULL THEN
        SET NEW.last_activity_at = CURRENT_TIMESTAMP;
    END IF;
END//
DELIMITER ;

-- 创建会话更新触发器
DELIMITER //
CREATE TRIGGER user_sessions_update_check 
BEFORE UPDATE ON user_sessions
FOR EACH ROW
BEGIN
    -- 如果会话被标记为非活跃，记录时间
    IF OLD.is_active = TRUE AND NEW.is_active = FALSE THEN
        SET NEW.last_activity_at = CURRENT_TIMESTAMP;
    END IF;
END//
DELIMITER ;

-- 添加设备信息JSON结构约束
ALTER TABLE user_sessions ADD CONSTRAINT chk_device_info_format 
CHECK (
    device_info IS NULL OR 
    JSON_VALID(device_info)
);

-- 添加IP地址格式约束
ALTER TABLE user_sessions ADD CONSTRAINT chk_ip_address_format 
CHECK (
    ip_address REGEXP '^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$' OR
    ip_address REGEXP '^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$' OR
    ip_address REGEXP '^::1$' OR
    ip_address REGEXP '^::(ffff:)?((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$'
);

-- 添加风险评分范围约束
ALTER TABLE user_sessions ADD CONSTRAINT chk_risk_score_range 
CHECK (risk_score >= 0 AND risk_score <= 100);

-- ============================================================================
-- 会话清理存储过程
-- ============================================================================

-- 创建清理过期会话的存储过程
DELIMITER //
CREATE PROCEDURE CleanExpiredSessions()
BEGIN
    DECLARE affected_rows INT DEFAULT 0;
    
    -- 清理过期的会话
    DELETE FROM user_sessions 
    WHERE expires_at < CURRENT_TIMESTAMP 
       OR (is_active = FALSE AND last_activity_at < DATE_SUB(CURRENT_TIMESTAMP, INTERVAL 7 DAY));
    
    SET affected_rows = ROW_COUNT();
    
    -- 记录清理结果
    INSERT INTO operation_logs (
        table_name, 
        operation_type, 
        affected_rows, 
        operation_desc, 
        created_at
    ) VALUES (
        'user_sessions',
        'DELETE',
        affected_rows,
        CONCAT('清理过期会话，删除 ', affected_rows, ' 条记录'),
        CURRENT_TIMESTAMP
    ) ON DUPLICATE KEY UPDATE operation_desc = operation_desc;
    
    SELECT CONCAT('已清理 ', affected_rows, ' 个过期会话') as result;
END//
DELIMITER ;

-- 创建用户会话统计视图
CREATE VIEW user_session_stats AS
SELECT 
    u.id as user_id,
    u.username,
    COUNT(s.id) as total_sessions,
    COUNT(CASE WHEN s.is_active = TRUE AND s.expires_at > CURRENT_TIMESTAMP THEN 1 END) as active_sessions,
    COUNT(CASE WHEN s.session_type = 'web' THEN 1 END) as web_sessions,
    COUNT(CASE WHEN s.session_type = 'mobile' THEN 1 END) as mobile_sessions,
    COUNT(CASE WHEN s.session_type = 'desktop' THEN 1 END) as desktop_sessions,
    MAX(s.last_activity_at) as last_activity,
    MAX(s.created_at) as last_login
FROM users u
LEFT JOIN user_sessions s ON u.id = s.user_id
WHERE u.deleted_at IS NULL
GROUP BY u.id, u.username;