-- +migrate Up
-- 创建迁移: 操作日志表
-- 版本: 20250112101600
-- 描述: 创建操作日志表，记录用户的所有操作行为
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:16:00
-- 依赖: 20250112100000_create_users_table
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1
-- 说明: 操作日志表用于审计追踪，记录系统中所有重要的用户操作

-- ============================================================================
-- 操作日志表 (operation_logs)
-- ============================================================================

CREATE TABLE operation_logs (
    -- 基础标识字段
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '操作日志唯一标识',
    
    -- 操作主体信息
    user_id BIGINT UNSIGNED DEFAULT NULL COMMENT '操作用户ID，系统操作为NULL',
    session_id VARCHAR(128) DEFAULT NULL COMMENT '用户会话ID',
    
    -- 操作分类和类型
    operation_type ENUM(
        'login', 'logout', 'register', 'password_change', 'profile_update',
        'file_upload', 'file_download', 'file_view', 'file_delete', 'file_restore',
        'file_share', 'file_unshare', 'file_move', 'file_copy', 'file_rename',
        'folder_create', 'folder_delete', 'folder_move', 'folder_rename',
        'team_create', 'team_update', 'team_delete', 'team_join', 'team_leave',
        'member_invite', 'member_remove', 'member_role_change',
        'permission_grant', 'permission_revoke', 'permission_update',
        'share_create', 'share_access', 'share_download', 'share_expire',
        'notification_send', 'notification_read',
        'settings_update', 'quota_change', 'backup_create', 'backup_restore',
        'system_maintenance', 'data_export', 'data_import'
    ) NOT NULL COMMENT '操作类型',
    
    operation_category ENUM('auth', 'file', 'team', 'permission', 'share', 'system', 'admin') NOT NULL COMMENT '操作类别',
    
    -- 操作目标资源
    resource_type ENUM('user', 'file', 'folder', 'team', 'share', 'permission', 'notification', 'system', 'session') NOT NULL COMMENT '资源类型',
    resource_id BIGINT UNSIGNED DEFAULT NULL COMMENT '资源ID',
    resource_name VARCHAR(255) DEFAULT NULL COMMENT '资源名称，冗余存储便于查看',
    resource_path VARCHAR(1000) DEFAULT NULL COMMENT '资源路径（文件/文件夹）',
    
    -- 操作详情和上下文
    action VARCHAR(100) NOT NULL COMMENT '具体操作动作',
    action_detail TEXT DEFAULT NULL COMMENT '操作详细描述',
    action_parameters JSON DEFAULT NULL COMMENT '操作参数',
    
    -- 操作结果和状态
    result ENUM('success', 'failure', 'partial', 'pending', 'cancelled') DEFAULT 'success' COMMENT '操作结果',
    error_code VARCHAR(50) DEFAULT NULL COMMENT '错误代码',
    error_message TEXT DEFAULT NULL COMMENT '错误详细信息',
    
    -- 性能和监控信息
    duration_ms INT UNSIGNED DEFAULT NULL COMMENT '操作耗时（毫秒）',
    cpu_usage DECIMAL(5,2) DEFAULT NULL COMMENT 'CPU使用率（%）',
    memory_usage BIGINT UNSIGNED DEFAULT NULL COMMENT '内存使用量（字节）',
    
    -- 网络和设备信息
    ip_address VARCHAR(45) DEFAULT NULL COMMENT '客户端IP地址',
    user_agent TEXT DEFAULT NULL COMMENT '用户代理字符串',
    device_info JSON DEFAULT NULL COMMENT '设备详细信息',
    location_info JSON DEFAULT NULL COMMENT '地理位置信息',
    
    -- 请求和响应信息
    request_method VARCHAR(10) DEFAULT NULL COMMENT 'HTTP请求方法',
    request_url VARCHAR(1000) DEFAULT NULL COMMENT '请求URL',
    request_headers JSON DEFAULT NULL COMMENT '请求头信息',
    response_code INT UNSIGNED DEFAULT NULL COMMENT 'HTTP响应状态码',
    response_size BIGINT UNSIGNED DEFAULT NULL COMMENT '响应数据大小（字节）',
    
    -- 业务上下文信息
    business_context JSON DEFAULT NULL COMMENT '业务上下文数据',
    correlation_id VARCHAR(64) DEFAULT NULL COMMENT '关联ID，用于跟踪相关操作',
    trace_id VARCHAR(64) DEFAULT NULL COMMENT '链路追踪ID',
    
    -- 变更前后状态
    before_state JSON DEFAULT NULL COMMENT '操作前状态',
    after_state JSON DEFAULT NULL COMMENT '操作后状态',
    change_summary TEXT DEFAULT NULL COMMENT '变更摘要',
    
    -- 安全和风险评估
    risk_score INT UNSIGNED DEFAULT 0 COMMENT '风险评分（0-100）',
    is_suspicious BOOLEAN DEFAULT FALSE COMMENT '是否被标记为可疑操作',
    security_flags JSON DEFAULT NULL COMMENT '安全标记',
    
    -- 审计和合规
    compliance_level ENUM('public', 'internal', 'confidential', 'restricted') DEFAULT 'internal' COMMENT '合规级别',
    audit_required BOOLEAN DEFAULT FALSE COMMENT '是否需要审计',
    retention_period INT UNSIGNED DEFAULT 2555 COMMENT '日志保留期（天），默认7年',
    
    -- API和集成信息
    api_version VARCHAR(20) DEFAULT NULL COMMENT 'API版本',
    api_key_id VARCHAR(64) DEFAULT NULL COMMENT 'API密钥ID',
    client_application VARCHAR(100) DEFAULT NULL COMMENT '客户端应用标识',
    integration_source VARCHAR(100) DEFAULT NULL COMMENT '集成来源系统',
    
    -- 批量操作相关
    batch_id VARCHAR(36) DEFAULT NULL COMMENT '批量操作批次ID',
    batch_sequence INT UNSIGNED DEFAULT NULL COMMENT '批次内序号',
    batch_total INT UNSIGNED DEFAULT NULL COMMENT '批次总数',
    
    -- 扩展属性
    metadata JSON DEFAULT NULL COMMENT '操作扩展元数据',
    tags JSON DEFAULT NULL COMMENT '操作标签',
    custom_fields JSON DEFAULT NULL COMMENT '自定义字段',
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '操作发生时间',
    processed_at TIMESTAMP NULL DEFAULT NULL COMMENT '日志处理时间',
    
    -- 业务索引设计
    INDEX idx_operation_logs_user_id (user_id) COMMENT '用户ID索引，查询用户操作记录',
    INDEX idx_operation_logs_session_id (session_id) COMMENT '会话ID索引，会话追踪',
    INDEX idx_operation_logs_operation_type (operation_type) COMMENT '操作类型索引',
    INDEX idx_operation_logs_operation_category (operation_category) COMMENT '操作类别索引',
    INDEX idx_operation_logs_resource_type (resource_type) COMMENT '资源类型索引',
    INDEX idx_operation_logs_resource_id (resource_id) COMMENT '资源ID索引',
    INDEX idx_operation_logs_action (action) COMMENT '操作动作索引',
    INDEX idx_operation_logs_result (result) COMMENT '操作结果索引',
    INDEX idx_operation_logs_ip_address (ip_address) COMMENT 'IP地址索引，安全分析',
    INDEX idx_operation_logs_created_at (created_at) COMMENT '创建时间索引',
    INDEX idx_operation_logs_correlation_id (correlation_id) COMMENT '关联ID索引',
    INDEX idx_operation_logs_trace_id (trace_id) COMMENT '链路追踪索引',
    INDEX idx_operation_logs_batch_id (batch_id) COMMENT '批次ID索引',
    INDEX idx_operation_logs_is_suspicious (is_suspicious) COMMENT '可疑操作索引',
    INDEX idx_operation_logs_audit_required (audit_required) COMMENT '审计需求索引',
    INDEX idx_operation_logs_api_key_id (api_key_id) COMMENT 'API密钥索引',
    
    -- 复合业务索引
    INDEX idx_operation_logs_user_type_time (user_id, operation_type, created_at DESC) COMMENT '用户操作类型时间复合索引',
    INDEX idx_operation_logs_resource_action_result (resource_type, action, result, created_at) COMMENT '资源操作结果复合索引',
    INDEX idx_operation_logs_category_result_time (operation_category, result, created_at DESC) COMMENT '类别结果时间复合索引',
    INDEX idx_operation_logs_ip_suspicious_time (ip_address, is_suspicious, created_at DESC) COMMENT 'IP可疑操作时间复合索引',
    INDEX idx_operation_logs_error_analysis (result, error_code, created_at DESC) COMMENT '错误分析复合索引',
    INDEX idx_operation_logs_performance_analysis (operation_type, duration_ms, created_at) COMMENT '性能分析复合索引',
    
    -- 分区优化索引（基于时间）
    INDEX idx_operation_logs_daily_stats (DATE(created_at), operation_category, result) COMMENT '日统计索引',
    INDEX idx_operation_logs_hourly_performance (HOUR(created_at), operation_type, duration_ms) COMMENT '小时性能索引',
    
    -- 全文搜索索引
    FULLTEXT idx_fulltext_search (action_detail, error_message, change_summary) COMMENT '操作详情全文搜索',
    
    -- 外键约束
    CONSTRAINT fk_operation_logs_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) 
        ON DELETE SET NULL 
        ON UPDATE CASCADE
        
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='操作日志表 - 记录系统中所有用户操作，用于审计和分析'
  ROW_FORMAT=DYNAMIC
  PARTITION BY RANGE (YEAR(created_at)) (
    PARTITION p2024 VALUES LESS THAN (2025),
    PARTITION p2025 VALUES LESS THAN (2026),
    PARTITION p2026 VALUES LESS THAN (2027),
    PARTITION p_future VALUES LESS THAN MAXVALUE
  );

-- ============================================================================
-- 操作日志表约束和检查
-- ============================================================================

-- IP地址格式约束
ALTER TABLE operation_logs ADD CONSTRAINT chk_operation_logs_ip_format 
CHECK (
    ip_address IS NULL OR
    ip_address REGEXP '^([0-9]{1,3}\\.){3}[0-9]{1,3}$' OR  -- IPv4
    ip_address REGEXP '^[0-9a-fA-F:]+$'                    -- IPv6（简化）
);

-- HTTP状态码约束
ALTER TABLE operation_logs ADD CONSTRAINT chk_response_code_range 
CHECK (
    response_code IS NULL OR 
    (response_code >= 100 AND response_code <= 599)
);

-- 性能指标约束
ALTER TABLE operation_logs ADD CONSTRAINT chk_performance_metrics 
CHECK (
    (duration_ms IS NULL OR duration_ms >= 0) AND
    (cpu_usage IS NULL OR (cpu_usage >= 0 AND cpu_usage <= 100)) AND
    (memory_usage IS NULL OR memory_usage >= 0) AND
    (response_size IS NULL OR response_size >= 0)
);

-- 风险评分约束
ALTER TABLE operation_logs ADD CONSTRAINT chk_risk_score_range 
CHECK (risk_score >= 0 AND risk_score <= 100);

-- 保留期约束
ALTER TABLE operation_logs ADD CONSTRAINT chk_retention_period 
CHECK (retention_period > 0 AND retention_period <= 3650); -- 最长10年

-- 批次信息约束
ALTER TABLE operation_logs ADD CONSTRAINT chk_batch_info_consistency 
CHECK (
    (batch_id IS NULL AND batch_sequence IS NULL AND batch_total IS NULL) OR
    (batch_id IS NOT NULL AND batch_sequence > 0 AND batch_total > 0 AND batch_sequence <= batch_total)
);

-- JSON字段验证
ALTER TABLE operation_logs ADD CONSTRAINT chk_operation_logs_json_valid 
CHECK (
    (action_parameters IS NULL OR JSON_VALID(action_parameters)) AND
    (device_info IS NULL OR JSON_VALID(device_info)) AND
    (location_info IS NULL OR JSON_VALID(location_info)) AND
    (request_headers IS NULL OR JSON_VALID(request_headers)) AND
    (business_context IS NULL OR JSON_VALID(business_context)) AND
    (before_state IS NULL OR JSON_VALID(before_state)) AND
    (after_state IS NULL OR JSON_VALID(after_state)) AND
    (security_flags IS NULL OR JSON_VALID(security_flags)) AND
    (metadata IS NULL OR JSON_VALID(metadata)) AND
    (tags IS NULL OR JSON_VALID(tags)) AND
    (custom_fields IS NULL OR JSON_VALID(custom_fields))
);

-- ============================================================================
-- 操作日志管理触发器
-- ============================================================================

-- 日志风险评估触发器
DELIMITER //
CREATE TRIGGER operation_logs_risk_assessment
BEFORE INSERT ON operation_logs
FOR EACH ROW
BEGIN
    DECLARE base_risk INT DEFAULT 0;
    DECLARE ip_risk INT DEFAULT 0;
    DECLARE time_risk INT DEFAULT 0;
    
    -- 基础风险评分
    CASE NEW.operation_type
        WHEN 'password_change' THEN SET base_risk = 30;
        WHEN 'permission_grant' THEN SET base_risk = 25;
        WHEN 'team_create' THEN SET base_risk = 20;
        WHEN 'file_delete' THEN SET base_risk = 15;
        WHEN 'login' THEN SET base_risk = 10;
        ELSE SET base_risk = 5;
    END CASE;
    
    -- 失败操作增加风险
    IF NEW.result = 'failure' THEN
        SET base_risk = base_risk + 20;
    END IF;
    
    -- 检查IP异常（简化实现）
    IF NEW.ip_address IS NOT NULL THEN
        -- 实际应用中需要检查IP黑名单、地理位置异常等
        -- 这里简化为检查是否为内网IP
        IF NOT (NEW.ip_address REGEXP '^(10\\.|172\\.(1[6-9]|2[0-9]|3[0-1])\\.|192\\.168\\.)') THEN
            SET ip_risk = 10;
        END IF;
    END IF;
    
    -- 非工作时间操作增加风险
    IF HOUR(NEW.created_at) < 8 OR HOUR(NEW.created_at) > 22 THEN
        SET time_risk = 5;
    END IF;
    
    -- 计算总风险评分
    SET NEW.risk_score = base_risk + ip_risk + time_risk;
    SET NEW.risk_score = LEAST(NEW.risk_score, 100);
    
    -- 高风险操作标记为可疑
    IF NEW.risk_score >= 50 THEN
        SET NEW.is_suspicious = TRUE;
    END IF;
    
    -- 生成关联ID（如果未提供）
    IF NEW.correlation_id IS NULL THEN
        SET NEW.correlation_id = UUID();
    END IF;
    
    -- 设置处理时间
    SET NEW.processed_at = CURRENT_TIMESTAMP;
END//
DELIMITER ;

-- 操作统计更新触发器
DELIMITER //
CREATE TRIGGER operation_logs_stats_update
AFTER INSERT ON operation_logs
FOR EACH ROW
BEGIN
    -- 更新用户操作统计（这里简化处理）
    -- 实际应用中可能需要更复杂的统计逻辑
    
    -- 异常操作告警（简化实现）
    IF NEW.is_suspicious = TRUE OR NEW.risk_score >= 70 THEN
        -- 可以在这里触发告警通知
        -- 实际实现中应该通过消息队列或其他异步方式处理
        INSERT INTO notifications (
            user_id, notification_type, title, content, priority,
            related_id, related_type, metadata
        ) 
        SELECT 
            u.id, 'security_alert', '检测到可疑操作', 
            CONCAT('用户 ', u.username, ' 执行了可疑操作：', NEW.action),
            'high', NEW.id, 'operation_log',
            JSON_OBJECT('risk_score', NEW.risk_score, 'operation_type', NEW.operation_type)
        FROM users u 
        WHERE u.role IN ('admin', 'super_admin') 
          AND u.status = 'active'
        LIMIT 1;
    END IF;
END//
DELIMITER ;

-- ============================================================================
-- 操作日志管理存储过程
-- ============================================================================

-- 记录操作日志存储过程
DELIMITER //
CREATE PROCEDURE LogOperation(
    IN user_id_param BIGINT UNSIGNED,
    IN operation_type_param VARCHAR(50),
    IN resource_type_param VARCHAR(20),
    IN resource_id_param BIGINT UNSIGNED,
    IN action_param VARCHAR(100),
    IN result_param VARCHAR(20),
    IN ip_address_param VARCHAR(45),
    IN details_param TEXT,
    OUT log_id BIGINT UNSIGNED
)
BEGIN
    INSERT INTO operation_logs (
        user_id, operation_type, operation_category, resource_type, 
        resource_id, action, result, ip_address, action_detail
    ) VALUES (
        user_id_param, 
        operation_type_param,
        CASE operation_type_param
            WHEN 'login' THEN 'auth'
            WHEN 'logout' THEN 'auth'
            WHEN 'file_upload' THEN 'file'
            WHEN 'file_download' THEN 'file'
            WHEN 'team_create' THEN 'team'
            WHEN 'permission_grant' THEN 'permission'
            ELSE 'system'
        END,
        resource_type_param,
        resource_id_param,
        action_param,
        result_param,
        ip_address_param,
        details_param
    );
    
    SET log_id = LAST_INSERT_ID();
END//
DELIMITER ;

-- 获取用户操作统计存储过程
DELIMITER //
CREATE PROCEDURE GetUserOperationStats(
    IN user_id_param BIGINT UNSIGNED,
    IN days_param INT DEFAULT 30
)
BEGIN
    SELECT 
        operation_category,
        COUNT(*) as total_operations,
        COUNT(CASE WHEN result = 'success' THEN 1 END) as successful_operations,
        COUNT(CASE WHEN result = 'failure' THEN 1 END) as failed_operations,
        COUNT(CASE WHEN is_suspicious = TRUE THEN 1 END) as suspicious_operations,
        AVG(duration_ms) as avg_duration_ms,
        AVG(risk_score) as avg_risk_score,
        MAX(created_at) as last_operation_time
    FROM operation_logs
    WHERE user_id = user_id_param
      AND created_at >= DATE_SUB(CURRENT_DATE, INTERVAL days_param DAY)
    GROUP BY operation_category
    ORDER BY total_operations DESC;
END//
DELIMITER ;

-- 系统操作趋势分析存储过程
DELIMITER //
CREATE PROCEDURE GetOperationTrends(
    IN days_param INT DEFAULT 7
)
BEGIN
    SELECT 
        DATE(created_at) as operation_date,
        operation_category,
        COUNT(*) as total_operations,
        COUNT(CASE WHEN result = 'success' THEN 1 END) as successful_count,
        COUNT(CASE WHEN result = 'failure' THEN 1 END) as failed_count,
        COUNT(CASE WHEN is_suspicious = TRUE THEN 1 END) as suspicious_count,
        AVG(duration_ms) as avg_duration,
        COUNT(DISTINCT user_id) as unique_users,
        COUNT(DISTINCT ip_address) as unique_ips
    FROM operation_logs
    WHERE created_at >= DATE_SUB(CURRENT_DATE, INTERVAL days_param DAY)
    GROUP BY DATE(created_at), operation_category
    ORDER BY operation_date DESC, operation_category;
END//
DELIMITER ;

-- 清理历史日志存储过程
DELIMITER //
CREATE PROCEDURE CleanupOperationLogs()
BEGIN
    DECLARE affected_rows INT DEFAULT 0;
    
    -- 删除超过保留期的日志
    DELETE FROM operation_logs 
    WHERE created_at < DATE_SUB(CURRENT_DATE, INTERVAL retention_period DAY);
    
    SET affected_rows = ROW_COUNT();
    
    -- 优化表空间
    OPTIMIZE TABLE operation_logs;
    
    SELECT CONCAT('已清理 ', affected_rows, ' 条历史操作日志') as result;
END//
DELIMITER ;

-- ============================================================================
-- 操作日志统计视图
-- ============================================================================

-- 操作统计概览视图
CREATE VIEW operation_stats_overview AS
SELECT 
    operation_category,
    operation_type,
    COUNT(*) as total_operations,
    COUNT(CASE WHEN result = 'success' THEN 1 END) as success_count,
    COUNT(CASE WHEN result = 'failure' THEN 1 END) as failure_count,
    ROUND(COUNT(CASE WHEN result = 'success' THEN 1 END) * 100.0 / COUNT(*), 2) as success_rate,
    COUNT(CASE WHEN is_suspicious = TRUE THEN 1 END) as suspicious_count,
    AVG(duration_ms) as avg_duration_ms,
    COUNT(DISTINCT user_id) as unique_users,
    COUNT(DISTINCT ip_address) as unique_ips,
    MAX(created_at) as last_operation_time
FROM operation_logs
WHERE created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
GROUP BY operation_category, operation_type
ORDER BY total_operations DESC;

-- 可疑操作监控视图
CREATE VIEW suspicious_operations_monitor AS
SELECT 
    ol.id,
    ol.user_id,
    u.username,
    ol.operation_type,
    ol.action,
    ol.resource_type,
    ol.resource_name,
    ol.ip_address,
    ol.risk_score,
    ol.result,
    ol.error_message,
    ol.created_at,
    CASE 
        WHEN ol.risk_score >= 80 THEN 'critical'
        WHEN ol.risk_score >= 60 THEN 'high'
        WHEN ol.risk_score >= 40 THEN 'medium'
        ELSE 'low'
    END as risk_level
FROM operation_logs ol
LEFT JOIN users u ON ol.user_id = u.id
WHERE ol.is_suspicious = TRUE
   OR ol.risk_score >= 40
   OR ol.result = 'failure'
ORDER BY ol.risk_score DESC, ol.created_at DESC;

-- 系统性能监控视图
CREATE VIEW system_performance_monitor AS
SELECT 
    operation_type,
    COUNT(*) as operation_count,
    AVG(duration_ms) as avg_duration_ms,
    MIN(duration_ms) as min_duration_ms,
    MAX(duration_ms) as max_duration_ms,
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY duration_ms) as median_duration_ms,
    PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY duration_ms) as p95_duration_ms,
    COUNT(CASE WHEN duration_ms > 5000 THEN 1 END) as slow_operations_count,
    ROUND(COUNT(CASE WHEN duration_ms > 5000 THEN 1 END) * 100.0 / COUNT(*), 2) as slow_operations_rate
FROM operation_logs
WHERE created_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR)
  AND duration_ms IS NOT NULL
GROUP BY operation_type
HAVING operation_count >= 10
ORDER BY avg_duration_ms DESC;