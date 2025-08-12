-- +migrate Up
-- 创建迁移: 分享访问日志表
-- 版本: 20250112101000
-- 描述: 创建分享访问日志表，记录所有分享的访问、下载等操作日志
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:10:00
-- 依赖: 20250112100900_create_file_shares_table, 20250112100300_create_files_table
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1
-- 说明: 分享访问日志表用于记录和分析分享的使用情况，支持安全审计和数据统计

-- ============================================================================
-- 分享访问日志表 (share_access_logs)
-- ============================================================================

CREATE TABLE share_access_logs (
    -- 基础标识字段
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '访问日志唯一标识',
    
    -- 关联的分享信息
    share_id BIGINT UNSIGNED NOT NULL COMMENT '关联的分享记录ID',
    share_code VARCHAR(8) NOT NULL COMMENT '分享码，冗余存储便于查询',
    
    -- 访问者信息
    visitor_ip VARCHAR(45) NOT NULL COMMENT '访问者IP地址，支持IPv4和IPv6',
    visitor_user_agent TEXT DEFAULT NULL COMMENT '访问者User-Agent字符串，包含浏览器和设备信息',
    visitor_fingerprint VARCHAR(64) DEFAULT NULL COMMENT '访问者设备指纹，用于唯一访客识别',
    visitor_session_id VARCHAR(128) DEFAULT NULL COMMENT '访问者会话ID，跟踪同一会话的多次操作',
    
    -- 地理位置信息
    visitor_location VARCHAR(100) DEFAULT NULL COMMENT '访问者地理位置，格式：城市,省份,国家',
    visitor_country_code VARCHAR(2) DEFAULT NULL COMMENT '访问者国家代码，ISO 3166-1 alpha-2',
    visitor_region VARCHAR(50) DEFAULT NULL COMMENT '访问者地区/省份',
    visitor_city VARCHAR(50) DEFAULT NULL COMMENT '访问者城市',
    visitor_timezone VARCHAR(50) DEFAULT NULL COMMENT '访问者时区',
    visitor_coordinates VARCHAR(50) DEFAULT NULL COMMENT '访问者经纬度坐标（如果可获取）',
    
    -- 访问类型和操作
    access_type ENUM('view', 'download', 'upload', 'comment', 'preview', 'stream', 'search', 'list') NOT NULL COMMENT '访问类型',
    operation_detail VARCHAR(100) DEFAULT NULL COMMENT '操作详情，如下载的具体文件名',
    
    -- 访问的文件信息
    file_id BIGINT UNSIGNED DEFAULT NULL COMMENT '访问的具体文件ID（如果是单文件操作）',
    file_name VARCHAR(255) DEFAULT NULL COMMENT '访问的文件名，冗余存储',
    file_size BIGINT UNSIGNED DEFAULT NULL COMMENT '访问的文件大小（下载时记录）',
    
    -- 访问结果和状态
    success BOOLEAN DEFAULT TRUE COMMENT '操作是否成功完成',
    response_code INT UNSIGNED DEFAULT 200 COMMENT 'HTTP响应状态码',
    error_message TEXT DEFAULT NULL COMMENT '错误详细信息（如果操作失败）',
    error_code VARCHAR(50) DEFAULT NULL COMMENT '错误代码分类',
    
    -- 性能和网络信息
    response_time INT UNSIGNED DEFAULT NULL COMMENT '响应时间（毫秒）',
    download_speed BIGINT UNSIGNED DEFAULT NULL COMMENT '下载速度（字节/秒，仅下载操作）',
    bytes_transferred BIGINT UNSIGNED DEFAULT NULL COMMENT '传输的字节数',
    
    -- 访问来源和渠道
    referrer_url VARCHAR(1000) DEFAULT NULL COMMENT '来源页面URL',
    referrer_domain VARCHAR(100) DEFAULT NULL COMMENT '来源域名',
    access_channel ENUM('web', 'mobile', 'api', 'embed', 'direct') DEFAULT 'web' COMMENT '访问渠道',
    
    -- 设备和浏览器信息
    device_type ENUM('desktop', 'mobile', 'tablet', 'bot', 'unknown') DEFAULT 'unknown' COMMENT '设备类型',
    browser_name VARCHAR(50) DEFAULT NULL COMMENT '浏览器名称',
    browser_version VARCHAR(20) DEFAULT NULL COMMENT '浏览器版本',
    os_name VARCHAR(50) DEFAULT NULL COMMENT '操作系统名称',
    os_version VARCHAR(20) DEFAULT NULL COMMENT '操作系统版本',
    
    -- 安全和验证信息
    authentication_method ENUM('none', 'password', 'token', 'whitelist', 'oauth') DEFAULT 'none' COMMENT '认证方式',
    password_attempts INT UNSIGNED DEFAULT 0 COMMENT '密码尝试次数（密码保护的分享）',
    is_suspicious BOOLEAN DEFAULT FALSE COMMENT '是否被标记为可疑访问',
    risk_score INT UNSIGNED DEFAULT 0 COMMENT '风险评分（0-100）',
    blocked_reason VARCHAR(200) DEFAULT NULL COMMENT '被阻止的原因（如果访问被拒绝）',
    
    -- 会话和行为分析
    session_duration INT UNSIGNED DEFAULT NULL COMMENT '会话持续时间（秒）',
    pages_viewed INT UNSIGNED DEFAULT 1 COMMENT '本次会话查看的页面数',
    files_downloaded INT UNSIGNED DEFAULT 0 COMMENT '本次会话下载的文件数',
    
    -- API访问相关
    api_key_id VARCHAR(64) DEFAULT NULL COMMENT 'API密钥标识（API访问时）',
    api_version VARCHAR(10) DEFAULT NULL COMMENT 'API版本',
    
    -- 扩展信息
    custom_headers JSON DEFAULT NULL COMMENT '自定义HTTP头信息',
    metadata JSON DEFAULT NULL COMMENT '扩展元数据，如特殊标记、实验参数等',
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '访问发生时间',
    processed_at TIMESTAMP NULL DEFAULT NULL COMMENT '日志处理时间（用于批量分析）',
    
    -- 业务索引设计
    INDEX idx_share_access_logs_share_id (share_id) COMMENT '分享ID索引，查询分享的访问记录',
    INDEX idx_share_access_logs_share_code (share_code) COMMENT '分享码索引',
    INDEX idx_share_access_logs_visitor_ip (visitor_ip) COMMENT '访问者IP索引，安全分析',
    INDEX idx_share_access_logs_access_type (access_type) COMMENT '访问类型索引，操作统计',
    INDEX idx_share_access_logs_file_id (file_id) COMMENT '文件ID索引，文件访问统计',
    INDEX idx_share_access_logs_created_at (created_at) COMMENT '访问时间索引，时间序列查询',
    INDEX idx_share_access_logs_success (success) COMMENT '成功状态索引，错误分析',
    INDEX idx_share_access_logs_device_type (device_type) COMMENT '设备类型索引，设备统计',
    INDEX idx_share_access_logs_country_code (visitor_country_code) COMMENT '国家代码索引，地理分布',
    INDEX idx_share_access_logs_is_suspicious (is_suspicious) COMMENT '可疑访问索引，安全监控',
    INDEX idx_share_access_logs_authentication_method (authentication_method) COMMENT '认证方式索引',
    INDEX idx_share_access_logs_visitor_fingerprint (visitor_fingerprint) COMMENT '设备指纹索引，用户识别',
    
    -- 复合业务索引
    INDEX idx_share_access_logs_share_time (share_id, created_at DESC) COMMENT '分享时间复合索引，访问历史查询',
    INDEX idx_share_access_logs_ip_time (visitor_ip, created_at DESC) COMMENT 'IP时间复合索引，IP行为分析',
    INDEX idx_share_access_logs_type_success_time (access_type, success, created_at DESC) COMMENT '类型成功时间复合索引',
    INDEX idx_share_access_logs_file_access_stats (file_id, access_type, success, created_at) COMMENT '文件访问统计优化索引',
    INDEX idx_share_access_logs_security_analysis (visitor_ip, is_suspicious, risk_score, created_at) COMMENT '安全分析复合索引',
    INDEX idx_share_access_logs_performance_analysis (access_type, response_time, bytes_transferred) COMMENT '性能分析索引',
    
    -- 分区优化索引（基于时间）
    INDEX idx_daily_stats (DATE(created_at), access_type, success) COMMENT '日统计索引',
    
    -- 外键约束
    CONSTRAINT fk_share_access_logs_share_id 
        FOREIGN KEY (share_id) REFERENCES file_shares(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_share_access_logs_file_id 
        FOREIGN KEY (file_id) REFERENCES files(id) 
        ON DELETE SET NULL 
        ON UPDATE CASCADE
        
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='分享访问日志表 - 记录所有分享的访问行为，支持安全审计和数据分析'
  ROW_FORMAT=DYNAMIC
  PARTITION BY RANGE (YEAR(created_at)) (
    PARTITION p2024 VALUES LESS THAN (2025),
    PARTITION p2025 VALUES LESS THAN (2026),
    PARTITION p2026 VALUES LESS THAN (2027),
    PARTITION p_future VALUES LESS THAN MAXVALUE
  );

-- ============================================================================
-- 分享访问日志约束和检查
-- ============================================================================

-- IP地址格式约束
ALTER TABLE share_access_logs ADD CONSTRAINT chk_visitor_ip_format 
CHECK (
    visitor_ip REGEXP '^([0-9]{1,3}\\.){3}[0-9]{1,3}$' OR  -- IPv4
    visitor_ip REGEXP '^[0-9a-fA-F:]+$'                    -- IPv6（简化检查）
);

-- 响应码范围约束
ALTER TABLE share_access_logs ADD CONSTRAINT chk_response_code_range 
CHECK (response_code >= 100 AND response_code <= 599);

-- 性能指标约束
ALTER TABLE share_access_logs ADD CONSTRAINT chk_performance_metrics 
CHECK (
    (response_time IS NULL OR response_time >= 0) AND
    (download_speed IS NULL OR download_speed >= 0) AND
    (bytes_transferred IS NULL OR bytes_transferred >= 0)
);

-- 风险评分约束
ALTER TABLE share_access_logs ADD CONSTRAINT chk_risk_score_range 
CHECK (risk_score >= 0 AND risk_score <= 100);

-- 密码尝试次数约束
ALTER TABLE share_access_logs ADD CONSTRAINT chk_password_attempts 
CHECK (password_attempts >= 0 AND password_attempts <= 10);

-- 会话统计约束
ALTER TABLE share_access_logs ADD CONSTRAINT chk_session_stats 
CHECK (
    (session_duration IS NULL OR session_duration >= 0) AND
    pages_viewed >= 0 AND
    files_downloaded >= 0
);

-- 国家代码格式约束
ALTER TABLE share_access_logs ADD CONSTRAINT chk_country_code_format 
CHECK (
    visitor_country_code IS NULL OR 
    visitor_country_code REGEXP '^[A-Z]{2}$'
);

-- JSON字段验证
ALTER TABLE share_access_logs ADD CONSTRAINT chk_share_logs_json_valid 
CHECK (
    (custom_headers IS NULL OR JSON_VALID(custom_headers)) AND
    (metadata IS NULL OR JSON_VALID(metadata))
);

-- ============================================================================
-- 访问日志分析触发器
-- ============================================================================

-- 访问统计更新触发器
DELIMITER //
CREATE TRIGGER share_access_logs_update_stats
AFTER INSERT ON share_access_logs
FOR EACH ROW
BEGIN
    -- 更新分享表的统计信息
    IF NEW.access_type = 'view' AND NEW.success = TRUE THEN
        UPDATE file_shares 
        SET view_count = view_count + 1,
            last_accessed_at = NEW.created_at
        WHERE id = NEW.share_id;
    END IF;
    
    IF NEW.access_type = 'download' AND NEW.success = TRUE THEN
        UPDATE file_shares 
        SET download_count = download_count + 1,
            last_downloaded_at = NEW.created_at
        WHERE id = NEW.share_id;
    END IF;
    
    -- 更新文件表的访问统计
    IF NEW.file_id IS NOT NULL THEN
        IF NEW.access_type = 'view' AND NEW.success = TRUE THEN
            UPDATE files 
            SET view_count = view_count + 1,
                last_accessed_at = NEW.created_at
            WHERE id = NEW.file_id;
        END IF;
        
        IF NEW.access_type = 'download' AND NEW.success = TRUE THEN
            UPDATE files 
            SET download_count = download_count + 1
            WHERE id = NEW.file_id;
        END IF;
    END IF;
END//
DELIMITER ;

-- 可疑访问检测触发器
DELIMITER //
CREATE TRIGGER share_access_logs_security_check
BEFORE INSERT ON share_access_logs
FOR EACH ROW
BEGIN
    DECLARE recent_failures INT DEFAULT 0;
    DECLARE ip_access_count INT DEFAULT 0;
    
    -- 检查最近5分钟内同一IP的失败访问次数
    SELECT COUNT(*) INTO recent_failures
    FROM share_access_logs
    WHERE visitor_ip = NEW.visitor_ip
      AND success = FALSE
      AND created_at >= DATE_SUB(NOW(), INTERVAL 5 MINUTE);
    
    -- 检查同一IP在1小时内的访问次数
    SELECT COUNT(*) INTO ip_access_count
    FROM share_access_logs
    WHERE visitor_ip = NEW.visitor_ip
      AND created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR);
    
    -- 计算风险评分
    SET NEW.risk_score = 0;
    
    -- 频繁失败访问
    IF recent_failures >= 5 THEN
        SET NEW.risk_score = NEW.risk_score + 30;
        SET NEW.is_suspicious = TRUE;
    END IF;
    
    -- 高频访问
    IF ip_access_count >= 100 THEN
        SET NEW.risk_score = NEW.risk_score + 25;
        SET NEW.is_suspicious = TRUE;
    END IF;
    
    -- 密码尝试过多
    IF NEW.password_attempts >= 3 THEN
        SET NEW.risk_score = NEW.risk_score + 20;
        SET NEW.is_suspicious = TRUE;
    END IF;
    
    -- 异常User-Agent
    IF NEW.visitor_user_agent IS NULL OR LENGTH(NEW.visitor_user_agent) < 10 THEN
        SET NEW.risk_score = NEW.risk_score + 15;
    END IF;
    
    -- 确保风险评分不超过100
    IF NEW.risk_score > 100 THEN
        SET NEW.risk_score = 100;
    END IF;
END//
DELIMITER ;

-- ============================================================================
-- 访问日志分析存储过程
-- ============================================================================

-- 获取分享访问统计的存储过程
DELIMITER //
CREATE PROCEDURE GetShareAccessStatistics(
    IN share_id_param BIGINT UNSIGNED,
    IN days_param INT DEFAULT 7
)
BEGIN
    SELECT 
        DATE(created_at) as access_date,
        COUNT(*) as total_accesses,
        COUNT(DISTINCT visitor_ip) as unique_visitors,
        COUNT(CASE WHEN access_type = 'view' THEN 1 END) as views,
        COUNT(CASE WHEN access_type = 'download' THEN 1 END) as downloads,
        COUNT(CASE WHEN success = FALSE THEN 1 END) as failed_accesses,
        AVG(response_time) as avg_response_time,
        SUM(bytes_transferred) as total_bytes_transferred
    FROM share_access_logs
    WHERE share_id = share_id_param
      AND created_at >= DATE_SUB(CURRENT_DATE, INTERVAL days_param DAY)
    GROUP BY DATE(created_at)
    ORDER BY access_date DESC;
END//
DELIMITER ;

-- 获取访问地理分布统计
DELIMITER //
CREATE PROCEDURE GetAccessGeographyStats(
    IN share_id_param BIGINT UNSIGNED,
    IN days_param INT DEFAULT 30
)
BEGIN
    SELECT 
        visitor_country_code,
        visitor_region,
        COUNT(*) as access_count,
        COUNT(DISTINCT visitor_ip) as unique_visitors,
        COUNT(CASE WHEN access_type = 'download' THEN 1 END) as downloads,
        AVG(response_time) as avg_response_time
    FROM share_access_logs
    WHERE share_id = share_id_param
      AND created_at >= DATE_SUB(CURRENT_DATE, INTERVAL days_param DAY)
      AND visitor_country_code IS NOT NULL
    GROUP BY visitor_country_code, visitor_region
    ORDER BY access_count DESC;
END//
DELIMITER ;

-- 安全威胁检测存储过程
DELIMITER //
CREATE PROCEDURE DetectSecurityThreats(
    IN hours_param INT DEFAULT 24
)
BEGIN
    -- 检测高风险IP
    SELECT 
        visitor_ip,
        COUNT(*) as access_attempts,
        COUNT(CASE WHEN success = FALSE THEN 1 END) as failed_attempts,
        COUNT(CASE WHEN is_suspicious = TRUE THEN 1 END) as suspicious_attempts,
        AVG(risk_score) as avg_risk_score,
        MAX(risk_score) as max_risk_score,
        COUNT(DISTINCT share_id) as shares_accessed
    FROM share_access_logs
    WHERE created_at >= DATE_SUB(NOW(), INTERVAL hours_param HOUR)
      AND (is_suspicious = TRUE OR risk_score >= 50)
    GROUP BY visitor_ip
    HAVING failed_attempts >= 5 OR suspicious_attempts >= 3
    ORDER BY avg_risk_score DESC, failed_attempts DESC;
END//
DELIMITER ;

-- 清理历史日志的存储过程
DELIMITER //
CREATE PROCEDURE CleanOldAccessLogs(
    IN retention_days INT DEFAULT 90
)
BEGIN
    DECLARE affected_rows INT DEFAULT 0;
    
    -- 删除超过保留期的日志
    DELETE FROM share_access_logs 
    WHERE created_at < DATE_SUB(CURRENT_DATE, INTERVAL retention_days DAY);
    
    SET affected_rows = ROW_COUNT();
    
    -- 优化表空间
    OPTIMIZE TABLE share_access_logs;
    
    SELECT CONCAT('已清理 ', affected_rows, ' 条历史访问日志') as result;
END//
DELIMITER ;

-- ============================================================================
-- 访问日志分析视图
-- ============================================================================

-- 热门分享访问统计视图
CREATE VIEW popular_share_access_stats AS
SELECT 
    sal.share_id,
    fs.share_code,
    fs.share_name,
    COUNT(*) as total_accesses,
    COUNT(DISTINCT sal.visitor_ip) as unique_visitors,
    COUNT(CASE WHEN sal.access_type = 'view' THEN 1 END) as total_views,
    COUNT(CASE WHEN sal.access_type = 'download' THEN 1 END) as total_downloads,
    COUNT(CASE WHEN sal.success = FALSE THEN 1 END) as failed_accesses,
    ROUND(AVG(sal.response_time), 2) as avg_response_time,
    MAX(sal.created_at) as last_access_time,
    ROUND(COUNT(CASE WHEN sal.success = FALSE THEN 1 END) * 100.0 / COUNT(*), 2) as failure_rate
FROM share_access_logs sal
JOIN file_shares fs ON sal.share_id = fs.id
WHERE sal.created_at >= DATE_SUB(CURRENT_DATE, INTERVAL 30 DAY)
GROUP BY sal.share_id, fs.share_code, fs.share_name
HAVING total_accesses >= 10
ORDER BY total_accesses DESC, unique_visitors DESC;

-- 访问设备统计视图
CREATE VIEW access_device_stats AS
SELECT 
    device_type,
    browser_name,
    os_name,
    COUNT(*) as access_count,
    COUNT(DISTINCT visitor_ip) as unique_users,
    COUNT(CASE WHEN access_type = 'download' THEN 1 END) as downloads,
    ROUND(AVG(response_time), 2) as avg_response_time,
    ROUND(COUNT(CASE WHEN success = FALSE THEN 1 END) * 100.0 / COUNT(*), 2) as failure_rate
FROM share_access_logs
WHERE created_at >= DATE_SUB(CURRENT_DATE, INTERVAL 7 DAY)
  AND device_type IS NOT NULL
GROUP BY device_type, browser_name, os_name
ORDER BY access_count DESC;

-- 可疑访问监控视图
CREATE VIEW suspicious_access_monitor AS
SELECT 
    visitor_ip,
    COUNT(*) as total_attempts,
    COUNT(CASE WHEN success = FALSE THEN 1 END) as failed_attempts,
    COUNT(CASE WHEN is_suspicious = TRUE THEN 1 END) as suspicious_attempts,
    COUNT(DISTINCT share_id) as shares_accessed,
    AVG(risk_score) as avg_risk_score,
    MAX(risk_score) as max_risk_score,
    MIN(created_at) as first_attempt,
    MAX(created_at) as last_attempt,
    GROUP_CONCAT(DISTINCT visitor_country_code) as countries
FROM share_access_logs
WHERE created_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR)
  AND (is_suspicious = TRUE OR risk_score >= 30)
GROUP BY visitor_ip
ORDER BY avg_risk_score DESC, failed_attempts DESC;