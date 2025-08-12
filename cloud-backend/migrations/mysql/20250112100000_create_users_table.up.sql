-- +migrate Up
-- 创建迁移: 用户表
-- 版本: 20250112100000
-- 描述: 创建用户基础信息表，包含用户认证、权限、配额等核心字段
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:00:00
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1

-- ============================================================================
-- 用户基础信息表 (users)
-- ============================================================================

CREATE TABLE users (
    -- 基础标识字段
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '用户唯一标识',
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名，用于登录，3-50字符',
    email VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱地址，用于登录和通知',
    email_verified BOOLEAN DEFAULT FALSE COMMENT '邮箱是否已验证',
    phone VARCHAR(20) DEFAULT NULL COMMENT '手机号码，支持国际格式',
    phone_verified BOOLEAN DEFAULT FALSE COMMENT '手机号是否已验证',
    
    -- 安全认证字段
    password_hash VARCHAR(255) NOT NULL COMMENT 'bcrypt加密后的密码，成本因子12',
    password_changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '密码最后修改时间',
    login_attempts INT UNSIGNED DEFAULT 0 COMMENT '连续登录失败次数',
    locked_until TIMESTAMP NULL DEFAULT NULL COMMENT '账户锁定截止时间',
    last_login_at TIMESTAMP NULL DEFAULT NULL COMMENT '最后登录时间',
    last_login_ip VARCHAR(45) DEFAULT NULL COMMENT '最后登录IP地址，支持IPv6',
    
    -- 双因子认证
    two_factor_enabled BOOLEAN DEFAULT FALSE COMMENT '是否启用双因子认证',
    two_factor_secret VARCHAR(32) DEFAULT NULL COMMENT '双因子认证密钥，Base32编码',
    
    -- 用户基本信息
    real_name VARCHAR(100) DEFAULT NULL COMMENT '真实姓名',
    avatar_url VARCHAR(500) DEFAULT NULL COMMENT '头像URL，支持OSS和CDN',
    
    -- 权限和状态
    role ENUM('user', 'vip', 'admin', 'super_admin') DEFAULT 'user' COMMENT '用户角色：普通用户、VIP、管理员、超级管理员',
    status ENUM('pending', 'active', 'disabled', 'banned') DEFAULT 'pending' COMMENT '用户状态：待激活、正常、禁用、封禁',
    
    -- 存储配额管理 (字节为单位)
    storage_quota BIGINT UNSIGNED DEFAULT 10737418240 COMMENT '存储配额，默认10GB（10*1024*1024*1024）',
    storage_used BIGINT UNSIGNED DEFAULT 0 COMMENT '已使用存储空间，字节',
    
    -- 带宽限制 (字节/秒)
    upload_bandwidth_limit INT UNSIGNED DEFAULT 10485760 COMMENT '上传带宽限制，默认10MB/s',
    download_bandwidth_limit INT UNSIGNED DEFAULT 10485760 COMMENT '下载带宽限制，默认10MB/s',
    max_file_size BIGINT UNSIGNED DEFAULT 10737418240 COMMENT '单文件最大大小，默认10GB',
    
    -- 文件类型限制 (JSON格式)
    allowed_file_types JSON DEFAULT NULL COMMENT '允许的文件类型数组，NULL表示无限制',
    forbidden_file_types JSON DEFAULT ('["exe", "bat", "com", "scr", "pif", "vbs", "js", "jar"]') COMMENT '禁止的文件类型数组，包含常见恶意文件扩展名',
    
    -- 用户偏好设置
    language VARCHAR(10) DEFAULT 'zh-CN' COMMENT '界面语言，ISO 639-1格式',
    timezone VARCHAR(50) DEFAULT 'Asia/Shanghai' COMMENT '时区设置，IANA时区标识符',
    theme ENUM('light', 'dark', 'auto') DEFAULT 'auto' COMMENT '主题设置：亮色、暗色、自动',
    
    -- 通知设置
    email_notifications BOOLEAN DEFAULT TRUE COMMENT '是否接收邮件通知',
    sms_notifications BOOLEAN DEFAULT FALSE COMMENT '是否接收短信通知',
    
    -- 时间戳字段
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '账户创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间，NULL表示未删除',
    
    -- 索引设计
    INDEX idx_users_username (username) COMMENT '用户名索引，用于登录验证',
    INDEX idx_users_email (email) COMMENT '邮箱索引，用于登录和重复检查',
    INDEX idx_users_phone (phone) COMMENT '手机号索引，用于登录和重复检查',
    INDEX idx_users_status (status) COMMENT '状态索引，用于筛选活跃用户',
    INDEX idx_users_role (role) COMMENT '角色索引，用于权限查询',
    INDEX idx_users_created_at (created_at) COMMENT '创建时间索引，用于统计和排序',
    INDEX idx_users_last_login_at (last_login_at) COMMENT '最后登录时间索引，用于活跃度统计',
    INDEX idx_users_deleted_at (deleted_at) COMMENT '软删除索引，配合软删除查询',
    INDEX idx_users_storage_used (storage_used) COMMENT '存储使用量索引，用于配额管理',
    INDEX idx_users_locked_until (locked_until) COMMENT '锁定时间索引，用于解锁查询'
    
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='用户基础信息表 - 存储用户核心数据，支持多种认证方式和精细化权限控制'
  ROW_FORMAT=DYNAMIC;

-- 创建存储使用量更新触发器
DELIMITER //
CREATE TRIGGER users_storage_quota_check 
BEFORE UPDATE ON users
FOR EACH ROW
BEGIN
    -- 检查存储使用量不能超过配额
    IF NEW.storage_used > NEW.storage_quota THEN
        SIGNAL SQLSTATE '45000' 
        SET MESSAGE_TEXT = '存储使用量不能超过配额限制';
    END IF;
    
    -- 检查存储配额不能小于已使用量
    IF NEW.storage_quota < OLD.storage_used THEN
        SIGNAL SQLSTATE '45000' 
        SET MESSAGE_TEXT = '存储配额不能小于当前已使用量';
    END IF;
END//
DELIMITER ;

-- 创建用户名格式检查约束
ALTER TABLE users ADD CONSTRAINT chk_username_format 
CHECK (username REGEXP '^[a-zA-Z0-9_]{3,50}$');

-- 创建邮箱格式检查约束  
ALTER TABLE users ADD CONSTRAINT chk_email_format 
CHECK (email REGEXP '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$');

-- 创建手机号格式检查约束（支持国际格式）
ALTER TABLE users ADD CONSTRAINT chk_phone_format 
CHECK (phone IS NULL OR phone REGEXP '^\\+?[1-9]\\d{1,14}$');

-- 添加存储配额范围约束
ALTER TABLE users ADD CONSTRAINT chk_storage_quota_range 
CHECK (storage_quota >= 0 AND storage_quota <= 10995116277760); -- 最大10TB

-- 添加带宽限制约束
ALTER TABLE users ADD CONSTRAINT chk_bandwidth_limits 
CHECK (
    upload_bandwidth_limit >= 1048576 AND upload_bandwidth_limit <= 1073741824 AND -- 1MB/s 到 1GB/s
    download_bandwidth_limit >= 1048576 AND download_bandwidth_limit <= 1073741824
);

-- 添加文件大小限制约束
ALTER TABLE users ADD CONSTRAINT chk_max_file_size_limit 
CHECK (max_file_size > 0 AND max_file_size <= 107374182400); -- 最大100GB

-- 添加登录失败次数约束
ALTER TABLE users ADD CONSTRAINT chk_login_attempts_limit 
CHECK (login_attempts >= 0 AND login_attempts <= 10);

-- 添加JSON字段验证约束
ALTER TABLE users ADD CONSTRAINT chk_users_json_valid 
CHECK (
    (allowed_file_types IS NULL OR JSON_VALID(allowed_file_types)) AND
    (forbidden_file_types IS NULL OR JSON_VALID(forbidden_file_types))
);

-- ============================================================================
-- 初始化系统管理员账户
-- ============================================================================

-- 插入超级管理员账户（密码: admin123，实际部署时需要修改）
INSERT INTO users (
    username, 
    email, 
    email_verified, 
    password_hash, 
    real_name, 
    role, 
    status,
    storage_quota,
    created_at
) VALUES (
    'admin',
    'admin@ycgcloud.com',
    TRUE,
    '$2a$12$rQ8Kd5K5K5K5K5K5K5K5Ku',  -- 需要在实际部署时生成真实的bcrypt哈希
    '系统管理员',
    'super_admin',
    'active',
    107374182400,  -- 100GB配额
    CURRENT_TIMESTAMP
);

-- 插入默认系统用户（用于系统操作）
INSERT INTO users (
    username, 
    email, 
    email_verified, 
    password_hash, 
    real_name, 
    role, 
    status,
    storage_quota,
    created_at
) VALUES (
    'system',
    'system@ycgcloud.com',
    TRUE,
    '$2a$12$sQ8Kd5K5K5K5K5K5K5K5Ku',  -- 需要在实际部署时生成真实的bcrypt哈希
    '系统用户',
    'admin',
    'active',
    0,  -- 系统用户无存储配额
    CURRENT_TIMESTAMP
);