-- +migrate Up
-- 创建迁移: 系统配置表
-- 版本: 20250112100800
-- 描述: 创建系统配置表，管理系统全局配置参数
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:08:00
-- 依赖: 20250112100000_create_users_table
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1
-- 说明: 系统配置表用于存储系统的动态配置参数，支持多种数据类型和权限控制

-- ============================================================================
-- 系统配置表 (system_configs)
-- ============================================================================

CREATE TABLE system_configs (
    -- 基础标识字段
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '配置项ID',
    
    -- 配置键值核心字段
    config_key VARCHAR(100) NOT NULL UNIQUE COMMENT '配置键名，使用点分命名法，如: app.name, storage.max_size',
    config_value TEXT NOT NULL COMMENT '配置值，根据类型存储不同格式的数据',
    config_type ENUM('string', 'number', 'boolean', 'json', 'array', 'encrypted') DEFAULT 'string' COMMENT '配置值数据类型',
    
    -- 配置分类和组织
    category VARCHAR(50) NOT NULL COMMENT '配置分类：app、storage、security、notification、ui、system、performance',
    subcategory VARCHAR(50) DEFAULT NULL COMMENT '配置子分类，便于细分管理',
    config_group VARCHAR(50) DEFAULT NULL COMMENT '配置组，用于界面分组显示',
    display_order INT UNSIGNED DEFAULT 0 COMMENT '显示顺序，用于界面排序',
    
    -- 配置描述和文档
    display_name VARCHAR(200) DEFAULT NULL COMMENT '配置项显示名称，用于界面展示',
    description TEXT DEFAULT NULL COMMENT '配置项详细描述',
    example_value TEXT DEFAULT NULL COMMENT '配置值示例',
    help_text TEXT DEFAULT NULL COMMENT '帮助说明文本',
    
    -- 权限和安全控制
    is_public BOOLEAN DEFAULT FALSE COMMENT '是否为公开配置（前端可访问）',
    is_editable BOOLEAN DEFAULT TRUE COMMENT '是否允许编辑',
    is_system BOOLEAN DEFAULT FALSE COMMENT '是否为系统级配置（谨慎修改）',
    is_sensitive BOOLEAN DEFAULT FALSE COMMENT '是否为敏感配置（需要特殊权限）',
    is_encrypted BOOLEAN DEFAULT FALSE COMMENT '是否加密存储',
    access_level ENUM('public', 'user', 'admin', 'system') DEFAULT 'admin' COMMENT '访问级别',
    
    -- 配置验证和约束
    validation_rule TEXT DEFAULT NULL COMMENT '验证规则（正则表达式或JSON格式规则）',
    validation_message VARCHAR(500) DEFAULT NULL COMMENT '验证失败提示信息',
    min_value DECIMAL(20,4) DEFAULT NULL COMMENT '数值类型最小值',
    max_value DECIMAL(20,4) DEFAULT NULL COMMENT '数值类型最大值',
    allowed_values JSON DEFAULT NULL COMMENT '允许的枚举值列表',
    default_value TEXT DEFAULT NULL COMMENT '默认值',
    
    -- 配置状态和生命周期
    is_active BOOLEAN DEFAULT TRUE COMMENT '配置是否有效',
    is_deprecated BOOLEAN DEFAULT FALSE COMMENT '是否已废弃',
    deprecation_reason TEXT DEFAULT NULL COMMENT '废弃原因说明',
    replacement_key VARCHAR(100) DEFAULT NULL COMMENT '替代配置键',
    
    -- 配置变更和缓存控制
    requires_restart BOOLEAN DEFAULT FALSE COMMENT '修改后是否需要重启服务',
    cache_enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用缓存',
    cache_ttl INT UNSIGNED DEFAULT 3600 COMMENT '缓存生存时间（秒）',
    last_accessed_at TIMESTAMP NULL DEFAULT NULL COMMENT '最后访问时间',
    access_count INT UNSIGNED DEFAULT 0 COMMENT '访问次数统计',
    
    -- 配置版本和历史
    version VARCHAR(20) DEFAULT '1.0.0' COMMENT '配置版本号',
    config_schema JSON DEFAULT NULL COMMENT '配置模式定义（用于复杂配置的结构验证）',
    migration_notes TEXT DEFAULT NULL COMMENT '迁移说明',
    
    -- 环境和部署相关
    environment ENUM('development', 'testing', 'staging', 'production', 'all') DEFAULT 'all' COMMENT '适用环境',
    feature_flag BOOLEAN DEFAULT TRUE COMMENT '功能开关（用于A/B测试或灰度发布）',
    rollout_percentage INT UNSIGNED DEFAULT 100 COMMENT '推出百分比（0-100）',
    
    -- 监控和告警
    monitor_changes BOOLEAN DEFAULT FALSE COMMENT '是否监控配置变更',
    alert_on_change BOOLEAN DEFAULT FALSE COMMENT '变更时是否发送告警',
    alert_recipients JSON DEFAULT NULL COMMENT '告警接收者列表',
    
    -- 配置来源和同步
    config_source ENUM('manual', 'import', 'sync', 'migration', 'api') DEFAULT 'manual' COMMENT '配置来源',
    external_sync BOOLEAN DEFAULT FALSE COMMENT '是否需要外部同步',
    sync_endpoint VARCHAR(500) DEFAULT NULL COMMENT '同步端点URL',
    last_sync_at TIMESTAMP NULL DEFAULT NULL COMMENT '最后同步时间',
    
    -- 操作审计字段
    created_by BIGINT UNSIGNED DEFAULT NULL COMMENT '创建者ID',
    updated_by BIGINT UNSIGNED DEFAULT NULL COMMENT '最后更新者ID',
    approved_by BIGINT UNSIGNED DEFAULT NULL COMMENT '审批者ID（敏感配置需要审批）',
    approved_at TIMESTAMP NULL DEFAULT NULL COMMENT '审批时间',
    
    -- 扩展和自定义字段
    tags JSON DEFAULT NULL COMMENT '配置标签，便于分类和搜索',
    custom_attributes JSON DEFAULT NULL COMMENT '自定义属性',
    related_configs JSON DEFAULT NULL COMMENT '关联的配置项',
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    -- 业务索引设计
    INDEX idx_system_configs_config_key (config_key) COMMENT '配置键索引',
    INDEX idx_system_configs_category (category) COMMENT '配置分类索引',
    INDEX idx_system_configs_subcategory (subcategory) COMMENT '配置子分类索引',
    INDEX idx_system_configs_config_group (config_group) COMMENT '配置组索引',
    INDEX idx_system_configs_is_public (is_public) COMMENT '公开配置索引',
    INDEX idx_system_configs_is_editable (is_editable) COMMENT '可编辑配置索引',
    INDEX idx_system_configs_is_system (is_system) COMMENT '系统配置索引',
    INDEX idx_system_configs_is_sensitive (is_sensitive) COMMENT '敏感配置索引',
    INDEX idx_system_configs_access_level (access_level) COMMENT '访问级别索引',
    INDEX idx_system_configs_is_active (is_active) COMMENT '有效配置索引',
    INDEX idx_system_configs_environment (environment) COMMENT '环境索引',
    INDEX idx_system_configs_feature_flag (feature_flag) COMMENT '功能开关索引',
    INDEX idx_system_configs_config_source (config_source) COMMENT '配置来源索引',
    INDEX idx_system_configs_updated_at (updated_at) COMMENT '更新时间索引',
    INDEX idx_system_configs_created_by (created_by) COMMENT '创建者索引',
    INDEX idx_system_configs_updated_by (updated_by) COMMENT '更新者索引',
    
    -- 复合业务索引
    INDEX idx_system_configs_category_group (category, config_group) COMMENT '分类组合索引',
    INDEX idx_system_configs_access_active (access_level, is_active) COMMENT '访问级别活跃状态复合索引',
    INDEX idx_system_configs_public_active (is_public, is_active) COMMENT '公开活跃状态复合索引',
    INDEX idx_system_configs_system_sensitive (is_system, is_sensitive) COMMENT '系统敏感配置复合索引',
    INDEX idx_system_configs_env_flag (environment, feature_flag) COMMENT '环境功能开关复合索引',
    
    -- 全文搜索索引
    FULLTEXT idx_fulltext_search (display_name, description, help_text) COMMENT '全文搜索索引',
    
    -- 外键约束
    CONSTRAINT fk_system_configs_created_by 
        FOREIGN KEY (created_by) REFERENCES users(id) 
        ON DELETE SET NULL 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_system_configs_updated_by 
        FOREIGN KEY (updated_by) REFERENCES users(id) 
        ON DELETE SET NULL 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_system_configs_approved_by 
        FOREIGN KEY (approved_by) REFERENCES users(id) 
        ON DELETE SET NULL 
        ON UPDATE CASCADE
        
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='系统配置表 - 存储系统全局配置参数，支持多种数据类型和权限控制'
  ROW_FORMAT=DYNAMIC;

-- ============================================================================
-- 系统配置表约束和检查
-- ============================================================================

-- 配置键命名规范约束
ALTER TABLE system_configs ADD CONSTRAINT chk_config_key_format 
CHECK (
    config_key REGEXP '^[a-z][a-z0-9_]*(\.[a-z][a-z0-9_]*)*$' AND
    LENGTH(config_key) >= 3 AND
    LENGTH(config_key) <= 100
);

-- 数值范围约束
ALTER TABLE system_configs ADD CONSTRAINT chk_numeric_range 
CHECK (
    (min_value IS NULL OR max_value IS NULL OR min_value <= max_value) AND
    (config_type != 'number' OR (min_value IS NOT NULL OR max_value IS NOT NULL))
);

-- 显示顺序约束
ALTER TABLE system_configs ADD CONSTRAINT chk_display_order 
CHECK (display_order >= 0 AND display_order <= 9999);

-- 缓存TTL约束
ALTER TABLE system_configs ADD CONSTRAINT chk_cache_ttl 
CHECK (cache_ttl > 0 AND cache_ttl <= 604800); -- 最长7天

-- 推出百分比约束
ALTER TABLE system_configs ADD CONSTRAINT chk_rollout_percentage 
CHECK (rollout_percentage >= 0 AND rollout_percentage <= 100);

-- 访问次数约束
ALTER TABLE system_configs ADD CONSTRAINT chk_access_count 
CHECK (access_count >= 0);

-- JSON字段验证
ALTER TABLE system_configs ADD CONSTRAINT chk_system_configs_json_valid 
CHECK (
    (allowed_values IS NULL OR JSON_VALID(allowed_values)) AND
    (config_schema IS NULL OR JSON_VALID(config_schema)) AND
    (alert_recipients IS NULL OR JSON_VALID(alert_recipients)) AND
    (tags IS NULL OR JSON_VALID(tags)) AND
    (custom_attributes IS NULL OR JSON_VALID(custom_attributes)) AND
    (related_configs IS NULL OR JSON_VALID(related_configs))
);

-- 敏感配置必须审批约束
ALTER TABLE system_configs ADD CONSTRAINT chk_sensitive_approval 
CHECK (
    NOT (is_sensitive = TRUE AND is_active = TRUE AND approved_by IS NULL)
);

-- 废弃配置约束
ALTER TABLE system_configs ADD CONSTRAINT chk_deprecated_config 
CHECK (
    NOT (is_deprecated = TRUE AND replacement_key IS NULL AND deprecation_reason IS NULL)
);

-- ============================================================================
-- 系统配置管理触发器
-- ============================================================================

-- 配置变更审计触发器
DELIMITER //
CREATE TRIGGER system_configs_audit_update
AFTER UPDATE ON system_configs
FOR EACH ROW
BEGIN
    -- 记录敏感配置的变更
    IF NEW.is_sensitive = TRUE OR OLD.config_value != NEW.config_value THEN
        INSERT INTO system_logs (
            log_type, 
            module, 
            action, 
            user_id, 
            details, 
            ip_address, 
            created_at
        ) VALUES (
            'config_change',
            'system_config',
            'update',
            NEW.updated_by,
            JSON_OBJECT(
                'config_key', NEW.config_key,
                'old_value', CASE WHEN OLD.is_sensitive THEN '***SENSITIVE***' ELSE OLD.config_value END,
                'new_value', CASE WHEN NEW.is_sensitive THEN '***SENSITIVE***' ELSE NEW.config_value END,
                'category', NEW.category,
                'is_sensitive', NEW.is_sensitive
            ),
            '127.0.0.1', -- 在实际应用中应该从会话中获取
            CURRENT_TIMESTAMP
        );
    END IF;
    
    -- 更新访问统计
    IF OLD.config_value != NEW.config_value THEN
        UPDATE system_configs 
        SET access_count = access_count + 1,
            last_accessed_at = CURRENT_TIMESTAMP
        WHERE id = NEW.id;
    END IF;
END//
DELIMITER ;

-- 配置缓存失效触发器
DELIMITER //
CREATE TRIGGER system_configs_cache_invalidate
AFTER UPDATE ON system_configs
FOR EACH ROW
BEGIN
    -- 如果配置值发生变化，标记需要清除缓存
    IF OLD.config_value != NEW.config_value OR OLD.is_active != NEW.is_active THEN
        -- 在实际应用中，这里会触发缓存清除操作
        -- 这里使用一个临时表来记录需要清除的缓存键
        INSERT INTO cache_invalidation_queue (cache_key, created_at)
        VALUES (CONCAT('config:', NEW.config_key), CURRENT_TIMESTAMP)
        ON DUPLICATE KEY UPDATE created_at = CURRENT_TIMESTAMP;
    END IF;
END//
DELIMITER ;

-- ============================================================================
-- 系统配置管理存储过程
-- ============================================================================

-- 批量导入配置的存储过程
DELIMITER //
CREATE PROCEDURE ImportSystemConfigs(
    IN config_data JSON,
    IN import_user_id BIGINT UNSIGNED,
    IN override_existing BOOLEAN DEFAULT FALSE
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE config_count INT DEFAULT 0;
    DECLARE current_config JSON;
    DECLARE config_key_val VARCHAR(100);
    DECLARE config_exists INT DEFAULT 0;
    
    -- 获取配置数组的长度
    SET config_count = JSON_LENGTH(config_data);
    
    -- 遍历配置数组
    WHILE i < config_count DO
        SET current_config = JSON_EXTRACT(config_data, CONCAT('$[', i, ']'));
        SET config_key_val = JSON_UNQUOTE(JSON_EXTRACT(current_config, '$.config_key'));
        
        -- 检查配置是否已存在
        SELECT COUNT(*) INTO config_exists 
        FROM system_configs 
        WHERE config_key = config_key_val;
        
        -- 如果不存在或允许覆盖，则插入或更新
        IF config_exists = 0 THEN
            INSERT INTO system_configs (
                config_key,
                config_value,
                config_type,
                category,
                description,
                is_public,
                is_editable,
                default_value,
                created_by,
                config_source
            ) VALUES (
                config_key_val,
                JSON_UNQUOTE(JSON_EXTRACT(current_config, '$.config_value')),
                JSON_UNQUOTE(JSON_EXTRACT(current_config, '$.config_type')),
                JSON_UNQUOTE(JSON_EXTRACT(current_config, '$.category')),
                JSON_UNQUOTE(JSON_EXTRACT(current_config, '$.description')),
                JSON_EXTRACT(current_config, '$.is_public'),
                JSON_EXTRACT(current_config, '$.is_editable'),
                JSON_UNQUOTE(JSON_EXTRACT(current_config, '$.default_value')),
                import_user_id,
                'import'
            );
        ELSEIF override_existing = TRUE THEN
            UPDATE system_configs 
            SET config_value = JSON_UNQUOTE(JSON_EXTRACT(current_config, '$.config_value')),
                description = JSON_UNQUOTE(JSON_EXTRACT(current_config, '$.description')),
                updated_by = import_user_id,
                config_source = 'import'
            WHERE config_key = config_key_val;
        END IF;
        
        SET i = i + 1;
    END WHILE;
    
    SELECT CONCAT('成功导入 ', config_count, ' 个配置项') as result;
END//
DELIMITER ;

-- 配置验证存储过程
DELIMITER //
CREATE PROCEDURE ValidateConfigValue(
    IN config_key_param VARCHAR(100),
    IN config_value_param TEXT,
    OUT is_valid BOOLEAN,
    OUT error_message VARCHAR(500)
)
BEGIN
    DECLARE config_type_val ENUM('string', 'number', 'boolean', 'json', 'array', 'encrypted');
    DECLARE validation_rule_val TEXT;
    DECLARE min_val DECIMAL(20,4);
    DECLARE max_val DECIMAL(20,4);
    DECLARE allowed_vals JSON;
    DECLARE numeric_val DECIMAL(20,4);
    
    SET is_valid = TRUE;
    SET error_message = NULL;
    
    -- 获取配置的验证规则
    SELECT config_type, validation_rule, min_value, max_value, allowed_values
    INTO config_type_val, validation_rule_val, min_val, max_val, allowed_vals
    FROM system_configs
    WHERE config_key = config_key_param;
    
    -- 如果配置不存在
    IF config_type_val IS NULL THEN
        SET is_valid = FALSE;
        SET error_message = '配置项不存在';
        LEAVE ValidateConfigValue;
    END IF;
    
    -- 根据类型进行验证
    CASE config_type_val
        WHEN 'boolean' THEN
            IF config_value_param NOT IN ('true', 'false', '1', '0') THEN
                SET is_valid = FALSE;
                SET error_message = '布尔值必须是 true, false, 1, 或 0';
            END IF;
            
        WHEN 'number' THEN
            SET numeric_val = CAST(config_value_param AS DECIMAL(20,4));
            IF numeric_val IS NULL THEN
                SET is_valid = FALSE;
                SET error_message = '数值格式不正确';
            ELSEIF min_val IS NOT NULL AND numeric_val < min_val THEN
                SET is_valid = FALSE;
                SET error_message = CONCAT('数值不能小于 ', min_val);
            ELSEIF max_val IS NOT NULL AND numeric_val > max_val THEN
                SET is_valid = FALSE;
                SET error_message = CONCAT('数值不能大于 ', max_val);
            END IF;
            
        WHEN 'json' THEN
            IF NOT JSON_VALID(config_value_param) THEN
                SET is_valid = FALSE;
                SET error_message = 'JSON格式不正确';
            END IF;
    END CASE;
    
    -- 验证允许的枚举值
    IF is_valid = TRUE AND allowed_vals IS NOT NULL THEN
        IF NOT JSON_CONTAINS(allowed_vals, JSON_QUOTE(config_value_param)) THEN
            SET is_valid = FALSE;
            SET error_message = CONCAT('值必须是以下之一: ', JSON_UNQUOTE(allowed_vals));
        END IF;
    END IF;
    
    -- 自定义验证规则（正则表达式）
    IF is_valid = TRUE AND validation_rule_val IS NOT NULL THEN
        IF NOT config_value_param REGEXP validation_rule_val THEN
            SET is_valid = FALSE;
            SET error_message = '值不符合验证规则';
        END IF;
    END IF;
    
END//
DELIMITER ;

-- ============================================================================
-- 系统配置管理视图
-- ============================================================================

-- 公开配置视图（前端可访问）
CREATE VIEW public_system_configs AS
SELECT 
    config_key,
    config_value,
    config_type,
    category,
    display_name,
    description,
    updated_at
FROM system_configs
WHERE is_public = TRUE 
  AND is_active = TRUE 
  AND feature_flag = TRUE
  AND (environment = 'all' OR environment = 'production');

-- 配置分类统计视图
CREATE VIEW config_category_stats AS
SELECT 
    category,
    COUNT(*) as total_configs,
    COUNT(CASE WHEN is_active = TRUE THEN 1 END) as active_configs,
    COUNT(CASE WHEN is_public = TRUE THEN 1 END) as public_configs,
    COUNT(CASE WHEN is_sensitive = TRUE THEN 1 END) as sensitive_configs,
    COUNT(CASE WHEN is_system = TRUE THEN 1 END) as system_configs,
    MAX(updated_at) as last_updated
FROM system_configs
GROUP BY category
ORDER BY category;

-- 最近修改的配置视图
CREATE VIEW recently_modified_configs AS
SELECT 
    sc.config_key,
    sc.config_value,
    sc.category,
    sc.is_sensitive,
    sc.updated_at,
    u.username as updated_by_user,
    sc.validation_rule,
    sc.description
FROM system_configs sc
LEFT JOIN users u ON sc.updated_by = u.id
WHERE sc.updated_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
ORDER BY sc.updated_at DESC;

-- ============================================================================
-- 初始化默认系统配置
-- ============================================================================

-- 插入默认系统配置
INSERT INTO system_configs (
    config_key, config_value, config_type, category, subcategory, 
    display_name, description, is_public, is_editable, is_system,
    default_value, environment, created_by
) VALUES
-- 应用基础配置
('app.name', '网络云盘系统', 'string', 'app', 'basic', '应用名称', '系统显示名称', TRUE, TRUE, TRUE, '网络云盘系统', 'all', 1),
('app.version', '1.0.0', 'string', 'app', 'basic', '应用版本', '当前系统版本号', TRUE, FALSE, TRUE, '1.0.0', 'all', 1),
('app.description', '高性能云存储解决方案', 'string', 'app', 'basic', '应用描述', '系统功能描述', TRUE, TRUE, FALSE, '高性能云存储解决方案', 'all', 1),
('app.timezone', 'Asia/Shanghai', 'string', 'app', 'basic', '默认时区', '系统默认时区设置', FALSE, TRUE, FALSE, 'Asia/Shanghai', 'all', 1),
('app.language', 'zh-CN', 'string', 'app', 'basic', '默认语言', '系统默认语言', TRUE, TRUE, FALSE, 'zh-CN', 'all', 1),

-- 存储配置
('storage.default_quota', '10737418240', 'number', 'storage', 'quota', '默认用户配额', '新用户默认存储配额（字节），10GB', FALSE, TRUE, FALSE, '10737418240', 'all', 1),
('storage.max_file_size', '1073741824', 'number', 'storage', 'limits', '单文件最大大小', '单个文件最大上传大小（字节），1GB', FALSE, TRUE, FALSE, '1073741824', 'all', 1),
('storage.allowed_extensions', '["jpg","jpeg","png","gif","pdf","doc","docx","xls","xlsx","ppt","pptx","txt","zip","rar"]', 'json', 'storage', 'limits', '允许的文件扩展名', '允许上传的文件类型列表', FALSE, TRUE, FALSE, '[]', 'all', 1),
('storage.enable_versioning', 'true', 'boolean', 'storage', 'features', '启用文件版本控制', '是否启用文件版本管理功能', FALSE, TRUE, FALSE, 'true', 'all', 1),
('storage.max_versions', '10', 'number', 'storage', 'features', '最大版本数', '每个文件保留的最大版本数', FALSE, TRUE, FALSE, '10', 'all', 1),

-- 安全配置
('security.password_min_length', '8', 'number', 'security', 'password', '密码最小长度', '用户密码最小长度要求', FALSE, TRUE, TRUE, '8', 'all', 1),
('security.password_require_special', 'true', 'boolean', 'security', 'password', '密码需要特殊字符', '密码是否必须包含特殊字符', FALSE, TRUE, TRUE, 'true', 'all', 1),
('security.session_timeout', '7200', 'number', 'security', 'session', '会话超时时间', '用户会话超时时间（秒），2小时', FALSE, TRUE, TRUE, '7200', 'all', 1),
('security.max_login_attempts', '5', 'number', 'security', 'login', '最大登录尝试次数', '账户锁定前的最大失败登录次数', FALSE, TRUE, TRUE, '5', 'all', 1),
('security.lockout_duration', '900', 'number', 'security', 'login', '账户锁定时间', '账户被锁定的时间（秒），15分钟', FALSE, TRUE, TRUE, '900', 'all', 1),

-- 通知配置
('notification.email_enabled', 'true', 'boolean', 'notification', 'email', '启用邮件通知', '是否启用邮件通知功能', FALSE, TRUE, FALSE, 'true', 'all', 1),
('notification.sms_enabled', 'false', 'boolean', 'notification', 'sms', '启用短信通知', '是否启用短信通知功能', FALSE, TRUE, FALSE, 'false', 'all', 1),
('notification.push_enabled', 'true', 'boolean', 'notification', 'push', '启用推送通知', '是否启用浏览器推送通知', TRUE, TRUE, FALSE, 'true', 'all', 1),

-- 界面配置
('ui.theme', 'light', 'string', 'ui', 'appearance', '默认主题', '系统默认主题', TRUE, TRUE, FALSE, 'light', 'all', 1),
('ui.items_per_page', '20', 'number', 'ui', 'pagination', '每页显示数量', '列表页面默认每页显示的项目数量', TRUE, TRUE, FALSE, '20', 'all', 1),
('ui.enable_animations', 'true', 'boolean', 'ui', 'effects', '启用动画效果', '是否启用界面动画效果', TRUE, TRUE, FALSE, 'true', 'all', 1),
('ui.sidebar_collapsed', 'false', 'boolean', 'ui', 'layout', '侧边栏默认折叠', '侧边栏是否默认折叠状态', TRUE, TRUE, FALSE, 'false', 'all', 1),

-- 系统性能配置
('system.cache_enabled', 'true', 'boolean', 'system', 'performance', '启用缓存', '是否启用系统缓存', FALSE, TRUE, TRUE, 'true', 'all', 1),
('system.cache_ttl', '3600', 'number', 'system', 'performance', '缓存有效期', '默认缓存有效期（秒），1小时', FALSE, TRUE, TRUE, '3600', 'all', 1),
('system.enable_compression', 'true', 'boolean', 'system', 'performance', '启用压缩', '是否启用文件压缩', FALSE, TRUE, FALSE, 'true', 'all', 1),
('system.log_level', 'info', 'string', 'system', 'logging', '日志级别', '系统日志记录级别', FALSE, TRUE, TRUE, 'info', 'all', 1),
('system.max_concurrent_uploads', '5', 'number', 'system', 'performance', '最大并发上传数', '单用户最大并发上传文件数', FALSE, TRUE, FALSE, '5', 'all', 1);

-- 为导入的配置设置适当的约束
UPDATE system_configs SET 
    min_value = 1, max_value = 50 WHERE config_key IN ('storage.max_versions', 'ui.items_per_page', 'system.max_concurrent_uploads');
UPDATE system_configs SET 
    min_value = 6, max_value = 32 WHERE config_key = 'security.password_min_length';
UPDATE system_configs SET 
    min_value = 1, max_value = 20 WHERE config_key = 'security.max_login_attempts';
UPDATE system_configs SET 
    allowed_values = '["light","dark","auto"]' WHERE config_key = 'ui.theme';
UPDATE system_configs SET 
    allowed_values = '["debug","info","warn","error"]' WHERE config_key = 'system.log_level';

-- 创建缓存失效队列表（用于触发器）
CREATE TABLE IF NOT EXISTS cache_invalidation_queue (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    cache_key VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='缓存失效队列';

-- 创建系统日志表（用于触发器）
CREATE TABLE IF NOT EXISTS system_logs (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    log_type VARCHAR(50) NOT NULL,
    module VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    user_id BIGINT UNSIGNED DEFAULT NULL,
    details JSON DEFAULT NULL,
    ip_address VARCHAR(45) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_log_type (log_type),
    INDEX idx_module (module),
    INDEX idx_user_id (user_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统操作日志表';