-- +migrate Up
-- 创建迁移: 安全日志表
-- 版本: 20250112101700
-- 描述: 创建安全日志表，记录系统安全相关事件
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:17:00
-- 依赖: 20250112100000_create_users_table
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1
-- 说明: 安全日志表专门记录安全相关事件，用于安全监控和威胁分析

-- ============================================================================
-- 安全日志表 (security_logs)
-- ============================================================================

CREATE TABLE security_logs (
    -- 基础标识字段
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '安全日志唯一标识',
    
    -- 关联用户信息
    user_id BIGINT UNSIGNED DEFAULT NULL COMMENT '关联用户ID，匿名事件为NULL',
    username VARCHAR(100) DEFAULT NULL COMMENT '用户名，冗余存储便于分析',
    session_id VARCHAR(128) DEFAULT NULL COMMENT '会话ID，用于关联会话',
    
    -- 安全事件分类
    event_type ENUM(
        'login_attempt', 'login_success', 'login_failure', 'logout',
        'password_change', 'password_reset', 'email_change', 'phone_change',
        'mfa_enable', 'mfa_disable', 'mfa_challenge', 'mfa_success', 'mfa_failure',
        'account_locked', 'account_unlocked', 'account_suspended', 'account_disabled',
        'permission_escalation', 'unauthorized_access', 'data_breach_attempt',
        'suspicious_activity', 'rate_limit_exceeded', 'ip_blocked', 'api_abuse',
        'file_access_denied', 'admin_access', 'system_intrusion',
        'malware_detected', 'virus_scan_alert', 'data_export_unauthorized',
        'encryption_key_change', 'certificate_expired', 'ssl_error'
    ) NOT NULL COMMENT '安全事件类型',
    
    event_category ENUM('authentication', 'authorization', 'data_access', 'system_security', 'compliance', 'threat_detection') NOT NULL COMMENT '事件类别',
    
    -- 严重程度和风险
    severity ENUM('info', 'warning', 'error', 'critical', 'emergency') DEFAULT 'info' COMMENT '事件严重程度',
    risk_score INT UNSIGNED DEFAULT 0 COMMENT '风险评分（0-100）',
    threat_level ENUM('low', 'medium', 'high', 'critical') DEFAULT 'low' COMMENT '威胁等级',
    
    -- 事件描述和详情
    title VARCHAR(200) NOT NULL COMMENT '事件标题',
    description TEXT NOT NULL COMMENT '事件详细描述',
    remediation_advice TEXT DEFAULT NULL COMMENT '修复建议',
    
    -- 网络和设备信息
    ip_address VARCHAR(45) DEFAULT NULL COMMENT '源IP地址',
    source_port INT UNSIGNED DEFAULT NULL COMMENT '源端口',
    destination_ip VARCHAR(45) DEFAULT NULL COMMENT '目标IP地址',
    destination_port INT UNSIGNED DEFAULT NULL COMMENT '目标端口',
    
    -- 用户代理和设备
    user_agent TEXT DEFAULT NULL COMMENT '用户代理字符串',
    device_fingerprint VARCHAR(64) DEFAULT NULL COMMENT '设备指纹标识',
    device_info JSON DEFAULT NULL COMMENT '设备详细信息',
    
    -- 地理位置信息
    location_info JSON DEFAULT NULL COMMENT '地理位置信息',
    country_code VARCHAR(2) DEFAULT NULL COMMENT '国家代码',
    region VARCHAR(100) DEFAULT NULL COMMENT '地区/省份',
    city VARCHAR(100) DEFAULT NULL COMMENT '城市',
    
    -- 攻击和威胁信息
    attack_type VARCHAR(100) DEFAULT NULL COMMENT '攻击类型',
    attack_vector VARCHAR(100) DEFAULT NULL COMMENT '攻击向量',
    payload TEXT DEFAULT NULL COMMENT '攻击载荷或恶意内容',
    signature VARCHAR(255) DEFAULT NULL COMMENT '威胁特征签名',
    
    -- 处理状态和响应
    status ENUM('detected', 'investigating', 'confirmed', 'mitigated', 'resolved', 'false_positive') DEFAULT 'detected' COMMENT '处理状态',
    is_blocked BOOLEAN DEFAULT FALSE COMMENT '是否已被阻止',
    is_automated_response BOOLEAN DEFAULT FALSE COMMENT '是否为自动响应',
    response_action VARCHAR(200) DEFAULT NULL COMMENT '响应动作',
    
    -- 关联信息
    correlation_id VARCHAR(64) DEFAULT NULL COMMENT '关联ID，用于关联相关事件',
    parent_event_id BIGINT UNSIGNED DEFAULT NULL COMMENT '父事件ID，用于事件链',
    related_events JSON DEFAULT NULL COMMENT '相关事件ID列表',
    
    -- 数据和资源信息
    affected_resource_type ENUM('user', 'file', 'system', 'network', 'application') DEFAULT NULL COMMENT '受影响的资源类型',
    affected_resource_id BIGINT UNSIGNED DEFAULT NULL COMMENT '受影响的资源ID',
    affected_resource_name VARCHAR(255) DEFAULT NULL COMMENT '受影响的资源名称',
    
    -- 合规和法规
    compliance_framework VARCHAR(100) DEFAULT NULL COMMENT '合规框架（如GDPR、SOX等）',
    regulation_impact BOOLEAN DEFAULT FALSE COMMENT '是否涉及法规影响',
    data_classification ENUM('public', 'internal', 'confidential', 'restricted') DEFAULT 'internal' COMMENT '数据分类级别',
    
    -- 检测和分析信息
    detection_method ENUM('rule_based', 'ml_model', 'anomaly_detection', 'signature', 'manual') DEFAULT 'rule_based' COMMENT '检测方法',
    detection_confidence DECIMAL(3,2) DEFAULT 0.50 COMMENT '检测置信度（0.00-1.00）',
    false_positive_probability DECIMAL(3,2) DEFAULT 0.10 COMMENT '误报概率',
    
    -- 时间信息
    event_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '事件实际发生时间',
    detection_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '事件检测时间',
    first_seen_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '首次发现时间',
    last_seen_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '最后发现时间',
    
    -- 统计信息
    occurrence_count INT UNSIGNED DEFAULT 1 COMMENT '事件发生次数',
    impact_score INT UNSIGNED DEFAULT 0 COMMENT '影响评分（0-100）',
    
    -- 外部威胁情报
    threat_intelligence JSON DEFAULT NULL COMMENT '威胁情报信息',
    ioc_indicators JSON DEFAULT NULL COMMENT 'IOC指标（IP、域名、哈希等）',
    attribution JSON DEFAULT NULL COMMENT '攻击归属信息',
    
    -- 扩展属性
    metadata JSON DEFAULT NULL COMMENT '事件扩展元数据',
    tags JSON DEFAULT NULL COMMENT '事件标签',
    custom_fields JSON DEFAULT NULL COMMENT '自定义字段',
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '日志创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '日志更新时间',
    
    -- 业务索引设计
    INDEX idx_security_logs_user_id (user_id) COMMENT '用户ID索引',
    INDEX idx_security_logs_username (username) COMMENT '用户名索引',
    INDEX idx_security_logs_session_id (session_id) COMMENT '会话ID索引',
    INDEX idx_security_logs_event_type (event_type) COMMENT '事件类型索引',
    INDEX idx_security_logs_event_category (event_category) COMMENT '事件类别索引',
    INDEX idx_security_logs_severity (severity) COMMENT '严重程度索引',
    INDEX idx_security_logs_threat_level (threat_level) COMMENT '威胁等级索引',
    INDEX idx_security_logs_risk_score (risk_score) COMMENT '风险评分索引',
    INDEX idx_security_logs_ip_address (ip_address) COMMENT 'IP地址索引',
    INDEX idx_security_logs_status (status) COMMENT '处理状态索引',
    INDEX idx_security_logs_is_blocked (is_blocked) COMMENT '阻止状态索引',
    INDEX idx_security_logs_attack_type (attack_type) COMMENT '攻击类型索引',
    INDEX idx_security_logs_detection_method (detection_method) COMMENT '检测方法索引',
    INDEX idx_security_logs_correlation_id (correlation_id) COMMENT '关联ID索引',
    INDEX idx_security_logs_parent_event_id (parent_event_id) COMMENT '父事件索引',
    INDEX idx_security_logs_affected_resource (affected_resource_type, affected_resource_id) COMMENT '受影响资源索引',
    INDEX idx_security_logs_country_code (country_code) COMMENT '国家代码索引',
    INDEX idx_security_logs_created_at (created_at) COMMENT '创建时间索引',
    INDEX idx_security_logs_event_timestamp (event_timestamp) COMMENT '事件时间索引',
    
    -- 复合业务索引
    INDEX idx_security_logs_user_event_time (user_id, event_type, created_at DESC) COMMENT '用户事件时间复合索引',
    INDEX idx_security_logs_ip_event_severity (ip_address, event_type, severity, created_at) COMMENT 'IP事件严重程度复合索引',
    INDEX idx_security_logs_severity_risk_time (severity, risk_score DESC, created_at DESC) COMMENT '严重程度风险时间复合索引',
    INDEX idx_security_logs_threat_detection (threat_level, detection_method, detection_confidence) COMMENT '威胁检测复合索引',
    INDEX idx_security_logs_security_monitoring (event_category, status, is_blocked, created_at) COMMENT '安全监控复合索引',
    INDEX idx_security_logs_incident_analysis (correlation_id, severity, status, created_at) COMMENT '事件分析复合索引',
    
    -- 分区优化索引（基于时间）
    INDEX idx_security_logs_daily_security_stats (DATE(created_at), event_category, severity) COMMENT '日安全统计索引',
    INDEX idx_security_logs_hourly_threat_stats (HOUR(created_at), threat_level, risk_score) COMMENT '小时威胁统计索引',
    
    -- 全文搜索索引
    FULLTEXT idx_fulltext_search (title, description, remediation_advice) COMMENT '安全事件内容全文搜索',
    
    -- 外键约束
    CONSTRAINT fk_security_logs_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) 
        ON DELETE SET NULL 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_security_logs_parent_event 
        FOREIGN KEY (parent_event_id) REFERENCES security_logs(id) 
        ON DELETE SET NULL 
        ON UPDATE CASCADE
        
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='安全日志表 - 记录系统安全事件，用于威胁检测和安全分析'
  ROW_FORMAT=DYNAMIC
  PARTITION BY RANGE (YEAR(created_at)) (
    PARTITION p2024 VALUES LESS THAN (2025),
    PARTITION p2025 VALUES LESS THAN (2026),
    PARTITION p2026 VALUES LESS THAN (2027),
    PARTITION p_future VALUES LESS THAN MAXVALUE
  );

-- ============================================================================
-- 安全日志表约束和检查
-- ============================================================================

-- IP地址格式约束
ALTER TABLE security_logs ADD CONSTRAINT chk_security_logs_ip_format 
CHECK (
    (ip_address IS NULL OR ip_address REGEXP '^([0-9]{1,3}\\.){3}[0-9]{1,3}$' OR ip_address REGEXP '^[0-9a-fA-F:]+$') AND
    (destination_ip IS NULL OR destination_ip REGEXP '^([0-9]{1,3}\\.){3}[0-9]{1,3}$' OR destination_ip REGEXP '^[0-9a-fA-F:]+$')
);

-- 端口号约束
ALTER TABLE security_logs ADD CONSTRAINT chk_port_ranges 
CHECK (
    (source_port IS NULL OR (source_port >= 1 AND source_port <= 65535)) AND
    (destination_port IS NULL OR (destination_port >= 1 AND destination_port <= 65535))
);

-- 风险和评分约束
ALTER TABLE security_logs ADD CONSTRAINT chk_risk_and_impact_scores 
CHECK (
    risk_score >= 0 AND risk_score <= 100 AND
    impact_score >= 0 AND impact_score <= 100
);

-- 置信度约束
ALTER TABLE security_logs ADD CONSTRAINT chk_confidence_probabilities 
CHECK (
    detection_confidence >= 0.00 AND detection_confidence <= 1.00 AND
    false_positive_probability >= 0.00 AND false_positive_probability <= 1.00
);

-- 发生次数约束
ALTER TABLE security_logs ADD CONSTRAINT chk_occurrence_count 
CHECK (occurrence_count >= 1);

-- 时间逻辑约束
ALTER TABLE security_logs ADD CONSTRAINT chk_security_logs_time_logic 
CHECK (
    detection_timestamp >= event_timestamp AND
    first_seen_at <= last_seen_at AND
    created_at >= event_timestamp
);

-- 国家代码格式约束
ALTER TABLE security_logs ADD CONSTRAINT chk_country_code_format 
CHECK (
    country_code IS NULL OR 
    country_code REGEXP '^[A-Z]{2}$'
);

-- JSON字段验证
ALTER TABLE security_logs ADD CONSTRAINT chk_security_logs_json_valid 
CHECK (
    (device_info IS NULL OR JSON_VALID(device_info)) AND
    (location_info IS NULL OR JSON_VALID(location_info)) AND
    (related_events IS NULL OR JSON_VALID(related_events)) AND
    (threat_intelligence IS NULL OR JSON_VALID(threat_intelligence)) AND
    (ioc_indicators IS NULL OR JSON_VALID(ioc_indicators)) AND
    (attribution IS NULL OR JSON_VALID(attribution)) AND
    (metadata IS NULL OR JSON_VALID(metadata)) AND
    (tags IS NULL OR JSON_VALID(tags)) AND
    (custom_fields IS NULL OR JSON_VALID(custom_fields))
);

-- ============================================================================
-- 安全日志管理触发器
-- ============================================================================

-- 安全事件自动分析触发器
DELIMITER //
CREATE TRIGGER security_logs_auto_analysis
BEFORE INSERT ON security_logs
FOR EACH ROW
BEGIN
    DECLARE similar_events_count INT DEFAULT 0;
    DECLARE recent_failures_count INT DEFAULT 0;
    DECLARE ip_reputation_score INT DEFAULT 50;
    
    -- 生成关联ID（如果未提供）
    IF NEW.correlation_id IS NULL THEN
        SET NEW.correlation_id = UUID();
    END IF;
    
    -- 自动风险评分计算
    SET NEW.risk_score = 0;
    
    -- 基础风险评分
    CASE NEW.event_type
        WHEN 'login_failure' THEN SET NEW.risk_score = NEW.risk_score + 15;
        WHEN 'unauthorized_access' THEN SET NEW.risk_score = NEW.risk_score + 40;
        WHEN 'data_breach_attempt' THEN SET NEW.risk_score = NEW.risk_score + 60;
        WHEN 'system_intrusion' THEN SET NEW.risk_score = NEW.risk_score + 80;
        WHEN 'malware_detected' THEN SET NEW.risk_score = NEW.risk_score + 70;
        ELSE SET NEW.risk_score = NEW.risk_score + 10;
    END CASE;
    
    -- 检查近期同类事件
    SELECT COUNT(*) INTO similar_events_count
    FROM security_logs
    WHERE event_type = NEW.event_type
      AND ip_address = NEW.ip_address
      AND created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR);
    
    -- 频繁事件增加风险
    IF similar_events_count >= 5 THEN
        SET NEW.risk_score = NEW.risk_score + 25;
        SET NEW.threat_level = 'high';
    ELSEIF similar_events_count >= 3 THEN
        SET NEW.risk_score = NEW.risk_score + 15;
        SET NEW.threat_level = 'medium';
    END IF;
    
    -- 检查登录失败次数
    IF NEW.event_type = 'login_failure' THEN
        SELECT COUNT(*) INTO recent_failures_count
        FROM security_logs
        WHERE event_type = 'login_failure'
          AND ip_address = NEW.ip_address
          AND created_at >= DATE_SUB(NOW(), INTERVAL 15 MINUTE);
        
        IF recent_failures_count >= 5 THEN
            SET NEW.risk_score = NEW.risk_score + 30;
            SET NEW.is_blocked = TRUE;
            SET NEW.response_action = 'IP_BLOCKED_AUTO';
        END IF;
    END IF;
    
    -- 确保风险评分不超过100
    SET NEW.risk_score = LEAST(NEW.risk_score, 100);
    
    -- 根据风险评分设置威胁等级
    IF NEW.threat_level IS NULL THEN
        IF NEW.risk_score >= 80 THEN
            SET NEW.threat_level = 'critical';
        ELSEIF NEW.risk_score >= 60 THEN
            SET NEW.threat_level = 'high';
        ELSEIF NEW.risk_score >= 40 THEN
            SET NEW.threat_level = 'medium';
        ELSE
            SET NEW.threat_level = 'low';
        END IF;
    END IF;
    
    -- 设置严重程度
    IF NEW.severity = 'info' THEN
        CASE NEW.threat_level
            WHEN 'critical' THEN SET NEW.severity = 'critical';
            WHEN 'high' THEN SET NEW.severity = 'error';
            WHEN 'medium' THEN SET NEW.severity = 'warning';
            ELSE SET NEW.severity = 'info';
        END CASE;
    END IF;
    
    -- 自动设置检测时间戳
    SET NEW.detection_timestamp = CURRENT_TIMESTAMP;
    
    -- 更新最后发现时间
    SET NEW.last_seen_at = CURRENT_TIMESTAMP;
END//
DELIMITER ;

-- 安全告警触发器
DELIMITER //
CREATE TRIGGER security_logs_alert_trigger
AFTER INSERT ON security_logs
FOR EACH ROW
BEGIN
    -- 高风险事件自动创建通知
    IF NEW.severity IN ('critical', 'error') OR NEW.risk_score >= 70 THEN
        INSERT INTO notifications (
            user_id, notification_type, title, content, priority,
            related_id, related_type, metadata
        )
        SELECT 
            u.id, 'security_alert', 
            CONCAT('安全警报: ', NEW.title),
            CONCAT('检测到高风险安全事件: ', NEW.description),
            CASE 
                WHEN NEW.severity = 'critical' THEN 'urgent'
                WHEN NEW.severity = 'error' THEN 'high'
                ELSE 'normal'
            END,
            NEW.id, 'security_log',
            JSON_OBJECT(
                'risk_score', NEW.risk_score,
                'threat_level', NEW.threat_level,
                'event_type', NEW.event_type,
                'ip_address', NEW.ip_address
            )
        FROM users u
        WHERE u.role IN ('admin', 'super_admin')
          AND u.status = 'active'
        LIMIT 5; -- 限制通知数量
    END IF;
    
    -- 更新安全统计（这里简化处理）
    -- 实际应用中可能需要更复杂的统计逻辑
END//
DELIMITER ;

-- 安全事件聚合触发器
DELIMITER //
CREATE TRIGGER security_logs_aggregation
AFTER INSERT ON security_logs
FOR EACH ROW
BEGIN
    DECLARE existing_event_id BIGINT UNSIGNED DEFAULT NULL;
    
    -- 检查是否存在相似的未解决事件（用于事件聚合）
    SELECT id INTO existing_event_id
    FROM security_logs
    WHERE event_type = NEW.event_type
      AND ip_address = NEW.ip_address
      AND user_id = NEW.user_id
      AND status NOT IN ('resolved', 'false_positive')
      AND created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR)
      AND id != NEW.id
    ORDER BY created_at DESC
    LIMIT 1;
    
    -- 如果找到相似事件，更新发生次数
    IF existing_event_id IS NOT NULL THEN
        UPDATE security_logs 
        SET occurrence_count = occurrence_count + 1,
            last_seen_at = NEW.created_at,
            risk_score = LEAST(risk_score + 5, 100),
            updated_at = CURRENT_TIMESTAMP
        WHERE id = existing_event_id;
        
        -- 设置当前事件的父事件
        UPDATE security_logs 
        SET parent_event_id = existing_event_id
        WHERE id = NEW.id;
    END IF;
END//
DELIMITER ;

-- ============================================================================
-- 安全日志管理存储过程
-- ============================================================================

-- 记录安全事件存储过程
DELIMITER //
CREATE PROCEDURE LogSecurityEvent(
    IN user_id_param BIGINT UNSIGNED,
    IN event_type_param VARCHAR(50),
    IN severity_param VARCHAR(20),
    IN title_param VARCHAR(200),
    IN description_param TEXT,
    IN ip_address_param VARCHAR(45),
    IN user_agent_param TEXT,
    OUT log_id BIGINT UNSIGNED
)
BEGIN
    DECLARE event_category_val VARCHAR(50);
    
    -- 根据事件类型确定类别
    SET event_category_val = CASE event_type_param
        WHEN 'login_attempt' THEN 'authentication'
        WHEN 'login_success' THEN 'authentication'
        WHEN 'login_failure' THEN 'authentication'
        WHEN 'unauthorized_access' THEN 'authorization'
        WHEN 'data_breach_attempt' THEN 'data_access'
        WHEN 'malware_detected' THEN 'threat_detection'
        ELSE 'system_security'
    END;
    
    INSERT INTO security_logs (
        user_id, event_type, event_category, severity, 
        title, description, ip_address, user_agent
    ) VALUES (
        user_id_param, event_type_param, event_category_val, severity_param,
        title_param, description_param, ip_address_param, user_agent_param
    );
    
    SET log_id = LAST_INSERT_ID();
END//
DELIMITER ;

-- 获取安全威胁统计存储过程
DELIMITER //
CREATE PROCEDURE GetSecurityThreatStats(
    IN days_param INT DEFAULT 7
)
BEGIN
    SELECT 
        event_category,
        threat_level,
        COUNT(*) as total_events,
        COUNT(CASE WHEN status = 'resolved' THEN 1 END) as resolved_events,
        COUNT(CASE WHEN is_blocked = TRUE THEN 1 END) as blocked_events,
        AVG(risk_score) as avg_risk_score,
        MAX(risk_score) as max_risk_score,
        COUNT(DISTINCT ip_address) as unique_ips,
        COUNT(DISTINCT user_id) as affected_users
    FROM security_logs
    WHERE created_at >= DATE_SUB(CURRENT_DATE, INTERVAL days_param DAY)
    GROUP BY event_category, threat_level
    ORDER BY threat_level DESC, total_events DESC;
END//
DELIMITER ;

-- 检测异常活动存储过程
DELIMITER //
CREATE PROCEDURE DetectAnomalousActivity(
    IN hours_param INT DEFAULT 24
)
BEGIN
    -- 检测异常IP活动
    SELECT 
        ip_address,
        COUNT(*) as event_count,
        COUNT(DISTINCT event_type) as event_types,
        COUNT(DISTINCT user_id) as affected_users,
        AVG(risk_score) as avg_risk_score,
        MAX(severity) as max_severity,
        MIN(created_at) as first_event,
        MAX(created_at) as last_event
    FROM security_logs
    WHERE created_at >= DATE_SUB(NOW(), INTERVAL hours_param HOUR)
      AND (risk_score >= 50 OR severity IN ('error', 'critical'))
    GROUP BY ip_address
    HAVING event_count >= 10 OR affected_users >= 3
    ORDER BY avg_risk_score DESC, event_count DESC;
END//
DELIMITER ;

-- 清理历史安全日志存储过程
DELIMITER //
CREATE PROCEDURE CleanupSecurityLogs()
BEGIN
    DECLARE affected_rows INT DEFAULT 0;
    
    -- 删除超过2年的低风险日志
    DELETE FROM security_logs 
    WHERE created_at < DATE_SUB(CURRENT_DATE, INTERVAL 730 DAY)
      AND severity = 'info'
      AND risk_score < 30;
    
    SET affected_rows = ROW_COUNT();
    
    -- 优化表空间
    OPTIMIZE TABLE security_logs;
    
    SELECT CONCAT('已清理 ', affected_rows, ' 条历史安全日志') as result;
END//
DELIMITER ;

-- ============================================================================
-- 安全日志统计视图
-- ============================================================================

-- 安全威胁概览视图
CREATE VIEW security_threat_overview AS
SELECT 
    event_category,
    threat_level,
    COUNT(*) as total_events,
    COUNT(CASE WHEN severity = 'critical' THEN 1 END) as critical_events,
    COUNT(CASE WHEN severity = 'error' THEN 1 END) as error_events,
    COUNT(CASE WHEN is_blocked = TRUE THEN 1 END) as blocked_events,
    COUNT(CASE WHEN status = 'resolved' THEN 1 END) as resolved_events,
    AVG(risk_score) as avg_risk_score,
    COUNT(DISTINCT ip_address) as unique_ips,
    COUNT(DISTINCT user_id) as affected_users,
    MAX(created_at) as last_event_time
FROM security_logs
WHERE created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
GROUP BY event_category, threat_level
ORDER BY threat_level DESC, total_events DESC;

-- 高风险IP监控视图
CREATE VIEW high_risk_ip_monitor AS
SELECT 
    ip_address,
    country_code,
    city,
    COUNT(*) as total_events,
    COUNT(DISTINCT event_type) as event_types,
    COUNT(DISTINCT user_id) as targeted_users,
    AVG(risk_score) as avg_risk_score,
    MAX(risk_score) as max_risk_score,
    COUNT(CASE WHEN is_blocked = TRUE THEN 1 END) as blocked_attempts,
    MIN(created_at) as first_seen,
    MAX(created_at) as last_seen,
    CASE 
        WHEN AVG(risk_score) >= 80 THEN 'critical'
        WHEN AVG(risk_score) >= 60 THEN 'high'
        WHEN AVG(risk_score) >= 40 THEN 'medium'
        ELSE 'low'
    END as risk_level
FROM security_logs
WHERE created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)
  AND risk_score >= 30
GROUP BY ip_address, country_code, city
HAVING total_events >= 5
ORDER BY avg_risk_score DESC, total_events DESC;

-- 未解决安全事件视图
CREATE VIEW unresolved_security_incidents AS
SELECT 
    sl.id,
    sl.user_id,
    u.username,
    sl.event_type,
    sl.severity,
    sl.threat_level,
    sl.title,
    sl.ip_address,
    sl.risk_score,
    sl.status,
    sl.created_at,
    TIMESTAMPDIFF(HOUR, sl.created_at, NOW()) as hours_since_detection,
    sl.occurrence_count
FROM security_logs sl
LEFT JOIN users u ON sl.user_id = u.id
WHERE sl.status NOT IN ('resolved', 'false_positive')
  AND sl.severity IN ('error', 'critical', 'warning')
ORDER BY sl.severity DESC, sl.risk_score DESC, sl.created_at DESC;