-- +migrate Up
-- 创建迁移: 文件上传任务表
-- 版本: 20250112100500
-- 描述: 创建分片上传任务管理表，支持断点续传、并发上传和上传状态跟踪
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:05:00
-- 依赖: 20250112100300_create_files_table
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1

-- ============================================================================
-- 文件上传任务管理表 (upload_tasks)
-- ============================================================================

CREATE TABLE upload_tasks (
    -- 基础标识字段
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '上传任务自增ID',
    upload_id VARCHAR(64) NOT NULL UNIQUE COMMENT '上传任务唯一标识，UUID格式',
    
    -- 用户和文件基础信息
    user_id BIGINT UNSIGNED NOT NULL COMMENT '上传用户ID',
    filename VARCHAR(255) NOT NULL COMMENT '原始文件名，保持用户上传时的命名',
    parent_id BIGINT UNSIGNED DEFAULT NULL COMMENT '目标父文件夹ID，NULL表示根目录',
    
    -- 文件属性信息
    file_size BIGINT UNSIGNED NOT NULL COMMENT '文件总大小（字节）',
    file_type VARCHAR(50) DEFAULT NULL COMMENT '文件类型分类：document、image、video等',
    mime_type VARCHAR(100) DEFAULT NULL COMMENT '文件MIME类型',
    file_extension VARCHAR(20) DEFAULT NULL COMMENT '文件扩展名',
    content_type VARCHAR(100) DEFAULT NULL COMMENT 'HTTP Content-Type头',
    
    -- 文件完整性信息
    md5_hash VARCHAR(32) DEFAULT NULL COMMENT '文件完整MD5哈希值，用于验证',
    sha256_hash VARCHAR(64) DEFAULT NULL COMMENT '文件SHA256哈希值，增强验证',
    expected_checksum VARCHAR(128) DEFAULT NULL COMMENT '客户端提供的预期校验和',
    
    -- 分片上传配置
    chunk_size INT UNSIGNED DEFAULT 2097152 COMMENT '分片大小（字节），默认2MB',
    total_chunks INT UNSIGNED NOT NULL COMMENT '总分片数量',
    uploaded_chunks INT UNSIGNED DEFAULT 0 COMMENT '已成功上传的分片数量',
    failed_chunks INT UNSIGNED DEFAULT 0 COMMENT '上传失败的分片数量',
    
    -- 上传进度统计
    uploaded_size BIGINT UNSIGNED DEFAULT 0 COMMENT '已上传大小（字节）',
    upload_percentage DECIMAL(5,2) DEFAULT 0.00 COMMENT '上传进度百分比（0.00-100.00）',
    
    -- 任务状态管理
    status ENUM('pending', 'preparing', 'uploading', 'merging', 'verifying', 'completed', 'failed', 'cancelled', 'expired') DEFAULT 'pending' COMMENT '上传任务状态',
    upload_mode ENUM('single', 'chunk', 'resumable') DEFAULT 'chunk' COMMENT '上传模式',
    
    -- 安全和权限
    upload_token VARCHAR(128) NOT NULL COMMENT '上传令牌，用于验证上传权限',
    session_id VARCHAR(64) DEFAULT NULL COMMENT '关联的用户会话ID',
    client_ip VARCHAR(45) NOT NULL COMMENT '客户端IP地址',
    user_agent TEXT DEFAULT NULL COMMENT '客户端User-Agent信息',
    
    -- 存储路径信息
    temp_dir VARCHAR(500) DEFAULT NULL COMMENT '临时文件存储目录',
    final_oss_key VARCHAR(500) DEFAULT NULL COMMENT '最终OSS对象键',
    upload_url VARCHAR(1000) DEFAULT NULL COMMENT '上传目标URL',
    
    -- 错误处理
    error_message TEXT DEFAULT NULL COMMENT '详细错误信息',
    error_code VARCHAR(50) DEFAULT NULL COMMENT '错误代码',
    last_error_at TIMESTAMP NULL DEFAULT NULL COMMENT '最后错误时间',
    
    -- 性能统计
    upload_speed BIGINT UNSIGNED DEFAULT 0 COMMENT '平均上传速度（字节/秒）',
    peak_speed BIGINT UNSIGNED DEFAULT 0 COMMENT '峰值上传速度（字节/秒）',
    estimated_time INT UNSIGNED DEFAULT 0 COMMENT '预估剩余时间（秒）',
    actual_duration INT UNSIGNED DEFAULT 0 COMMENT '实际上传耗时（秒）',
    
    -- 重试机制
    retry_count INT UNSIGNED DEFAULT 0 COMMENT '当前重试次数',
    max_retries INT UNSIGNED DEFAULT 3 COMMENT '最大重试次数',
    retry_delay INT UNSIGNED DEFAULT 5 COMMENT '重试延迟（秒）',
    
    -- 并发控制
    concurrent_chunks INT UNSIGNED DEFAULT 3 COMMENT '并发上传分片数量',
    max_concurrent INT UNSIGNED DEFAULT 5 COMMENT '最大并发数限制',
    
    -- 生命周期管理
    expires_at TIMESTAMP NOT NULL COMMENT '任务过期时间，过期后自动清理',
    cleanup_at TIMESTAMP NULL DEFAULT NULL COMMENT '清理时间，何时删除临时文件',
    
    -- 客户端信息
    client_metadata JSON DEFAULT NULL COMMENT '客户端元数据，设备信息、应用版本等',
    upload_options JSON DEFAULT NULL COMMENT '上传选项配置，压缩、加密等',
    
    -- 时间戳记录
    started_at TIMESTAMP NULL DEFAULT NULL COMMENT '开始上传时间',
    first_chunk_at TIMESTAMP NULL DEFAULT NULL COMMENT '第一个分片上传时间',
    last_chunk_at TIMESTAMP NULL DEFAULT NULL COMMENT '最后一个分片上传时间',
    completed_at TIMESTAMP NULL DEFAULT NULL COMMENT '上传完成时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '任务创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '任务更新时间',
    
    -- 业务索引设计
    INDEX idx_upload_tasks_upload_id (upload_id) COMMENT '上传ID索引，快速查找任务',
    INDEX idx_upload_tasks_user_id (user_id) COMMENT '用户ID索引，查询用户的上传任务',
    INDEX idx_upload_tasks_status (status) COMMENT '状态索引，按状态筛选任务',
    INDEX idx_upload_tasks_md5_hash (md5_hash) COMMENT 'MD5哈希索引，秒传功能支持',
    INDEX idx_upload_tasks_expires_at (expires_at) COMMENT '过期时间索引，清理任务使用',
    INDEX idx_upload_tasks_created_at (created_at) COMMENT '创建时间索引，时间排序',
    INDEX idx_upload_tasks_upload_token (upload_token) COMMENT '上传令牌索引，权限验证',
    INDEX idx_upload_tasks_session_id (session_id) COMMENT '会话ID索引，会话管理',
    INDEX idx_upload_tasks_client_ip (client_ip) COMMENT '客户端IP索引，安全分析',
    
    -- 复合业务索引
    INDEX idx_upload_tasks_user_status_created (user_id, status, created_at DESC) COMMENT '用户任务状态查询优化',
    INDEX idx_upload_tasks_status_expires (status, expires_at) COMMENT '状态过期时间查询优化',
    INDEX idx_upload_tasks_user_progress (user_id, upload_percentage, updated_at DESC) COMMENT '用户上传进度查询',
    INDEX idx_upload_tasks_hash_size (md5_hash, file_size) COMMENT '文件去重查询优化',
    
    -- 外键约束
    CONSTRAINT fk_upload_tasks_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_upload_tasks_parent_id 
        FOREIGN KEY (parent_id) REFERENCES files(id) 
        ON DELETE SET NULL  -- 如果父文件夹被删除，设为NULL（根目录）
        ON UPDATE CASCADE
        
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='文件上传任务管理表 - 支持分片上传、断点续传和并发控制'
  ROW_FORMAT=DYNAMIC;

-- ============================================================================
-- 上传任务表约束和触发器
-- ============================================================================

-- 文件大小合理性约束
ALTER TABLE upload_tasks ADD CONSTRAINT chk_upload_file_size 
CHECK (file_size > 0 AND file_size <= 107374182400); -- 不超过100GB

-- 分片大小合理性约束
ALTER TABLE upload_tasks ADD CONSTRAINT chk_chunk_size_range 
CHECK (chunk_size >= 1048576 AND chunk_size <= 104857600); -- 1MB到100MB

-- 分片数量合理性约束
ALTER TABLE upload_tasks ADD CONSTRAINT chk_total_chunks_positive 
CHECK (total_chunks > 0 AND total_chunks <= 10000); -- 最多10000个分片

-- 进度百分比约束
ALTER TABLE upload_tasks ADD CONSTRAINT chk_upload_percentage 
CHECK (upload_percentage >= 0.00 AND upload_percentage <= 100.00);

-- 重试次数约束
ALTER TABLE upload_tasks ADD CONSTRAINT chk_retry_counts 
CHECK (retry_count >= 0 AND max_retries >= 0 AND retry_count <= max_retries);

-- 并发数量约束
ALTER TABLE upload_tasks ADD CONSTRAINT chk_concurrent_limits 
CHECK (
    concurrent_chunks > 0 AND concurrent_chunks <= 20 AND
    max_concurrent > 0 AND max_concurrent <= 50 AND
    concurrent_chunks <= max_concurrent
);

-- 哈希格式约束
ALTER TABLE upload_tasks ADD CONSTRAINT chk_upload_md5_format 
CHECK (md5_hash IS NULL OR md5_hash REGEXP '^[a-f0-9]{32}$');

ALTER TABLE upload_tasks ADD CONSTRAINT chk_upload_sha256_format 
CHECK (sha256_hash IS NULL OR sha256_hash REGEXP '^[a-f0-9]{64}$');

-- 上传ID格式约束（UUID格式）
ALTER TABLE upload_tasks ADD CONSTRAINT chk_upload_id_format 
CHECK (upload_id REGEXP '^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$');

-- JSON字段验证
ALTER TABLE upload_tasks ADD CONSTRAINT chk_upload_json_valid 
CHECK (
    (client_metadata IS NULL OR JSON_VALID(client_metadata)) AND
    (upload_options IS NULL OR JSON_VALID(upload_options))
);

-- 时间逻辑约束
ALTER TABLE upload_tasks ADD CONSTRAINT chk_upload_time_logic 
CHECK (
    expires_at > created_at AND
    (started_at IS NULL OR started_at >= created_at) AND
    (completed_at IS NULL OR completed_at >= created_at)
);

-- 创建上传进度自动更新触发器
DELIMITER //
CREATE TRIGGER upload_tasks_progress_update 
BEFORE UPDATE ON upload_tasks
FOR EACH ROW
BEGIN
    -- 自动计算上传进度百分比
    IF NEW.total_chunks > 0 THEN
        SET NEW.upload_percentage = ROUND((NEW.uploaded_chunks / NEW.total_chunks) * 100, 2);
    END IF;
    
    -- 更新上传状态逻辑
    IF OLD.status != 'completed' AND NEW.uploaded_chunks = NEW.total_chunks AND NEW.total_chunks > 0 THEN
        SET NEW.status = 'merging';
    END IF;
    
    -- 记录开始时间
    IF OLD.status IN ('pending', 'preparing') AND NEW.status = 'uploading' AND NEW.started_at IS NULL THEN
        SET NEW.started_at = CURRENT_TIMESTAMP;
    END IF;
    
    -- 记录完成时间
    IF OLD.status != 'completed' AND NEW.status = 'completed' AND NEW.completed_at IS NULL THEN
        SET NEW.completed_at = CURRENT_TIMESTAMP;
        -- 计算实际耗时
        IF NEW.started_at IS NOT NULL THEN
            SET NEW.actual_duration = TIMESTAMPDIFF(SECOND, NEW.started_at, CURRENT_TIMESTAMP);
        END IF;
    END IF;
    
    -- 记录错误时间
    IF NEW.status IN ('failed', 'cancelled') AND OLD.status NOT IN ('failed', 'cancelled') THEN
        SET NEW.last_error_at = CURRENT_TIMESTAMP;
    END IF;
END//
DELIMITER ;

-- 创建上传任务过期检查触发器
DELIMITER //
CREATE TRIGGER upload_tasks_expiry_check 
BEFORE INSERT ON upload_tasks
FOR EACH ROW
BEGIN
    -- 确保过期时间大于创建时间
    IF NEW.expires_at <= CURRENT_TIMESTAMP THEN
        SIGNAL SQLSTATE '45000' 
        SET MESSAGE_TEXT = '上传任务过期时间必须大于当前时间';
    END IF;
    
    -- 设置默认的清理时间（过期后1小时）
    IF NEW.cleanup_at IS NULL THEN
        SET NEW.cleanup_at = DATE_ADD(NEW.expires_at, INTERVAL 1 HOUR);
    END IF;
    
    -- 生成上传令牌（如果未提供）
    IF NEW.upload_token IS NULL OR LENGTH(NEW.upload_token) = 0 THEN
        SET NEW.upload_token = CONCAT(
            'upload_',
            UNIX_TIMESTAMP(),
            '_',
            LEFT(MD5(CONCAT(NEW.upload_id, NEW.user_id, RAND())), 16)
        );
    END IF;
END//
DELIMITER ;

-- ============================================================================
-- 上传任务管理存储过程
-- ============================================================================

-- 创建清理过期上传任务的存储过程
DELIMITER //
CREATE PROCEDURE CleanupExpiredUploads()
BEGIN
    DECLARE affected_rows INT DEFAULT 0;
    
    -- 删除已过期且超过清理时间的任务
    DELETE FROM upload_tasks 
    WHERE status IN ('expired', 'failed', 'cancelled') 
      AND cleanup_at < CURRENT_TIMESTAMP;
    
    SET affected_rows = ROW_COUNT();
    
    -- 标记过期但未清理的任务
    UPDATE upload_tasks 
    SET status = 'expired',
        error_message = '上传任务已过期',
        last_error_at = CURRENT_TIMESTAMP
    WHERE expires_at < CURRENT_TIMESTAMP 
      AND status NOT IN ('completed', 'expired', 'failed', 'cancelled');
    
    SET affected_rows = affected_rows + ROW_COUNT();
    
    SELECT CONCAT('已处理 ', affected_rows, ' 个过期上传任务') as result;
END//
DELIMITER ;

-- 创建获取用户上传统计的存储过程
DELIMITER //
CREATE PROCEDURE GetUserUploadStats(
    IN target_user_id BIGINT UNSIGNED,
    IN days_back INT DEFAULT 30
)
BEGIN
    DECLARE start_date TIMESTAMP DEFAULT DATE_SUB(CURRENT_TIMESTAMP, INTERVAL days_back DAY);
    
    SELECT 
        COUNT(*) as total_uploads,
        COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_uploads,
        COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_uploads,
        COUNT(CASE WHEN status IN ('uploading', 'merging', 'verifying') THEN 1 END) as active_uploads,
        COALESCE(SUM(CASE WHEN status = 'completed' THEN file_size END), 0) as total_uploaded_size,
        COALESCE(AVG(CASE WHEN status = 'completed' THEN actual_duration END), 0) as avg_upload_time,
        COALESCE(AVG(CASE WHEN status = 'completed' THEN upload_speed END), 0) as avg_upload_speed,
        COALESCE(MAX(upload_speed), 0) as peak_upload_speed
    FROM upload_tasks 
    WHERE user_id = target_user_id 
      AND created_at >= start_date;
END//
DELIMITER ;

-- ============================================================================
-- 上传任务统计视图
-- ============================================================================

-- 创建上传任务统计视图
CREATE VIEW upload_task_stats AS
SELECT 
    u.id as user_id,
    u.username,
    COUNT(ut.id) as total_tasks,
    COUNT(CASE WHEN ut.status = 'completed' THEN 1 END) as completed_tasks,
    COUNT(CASE WHEN ut.status = 'failed' THEN 1 END) as failed_tasks,
    COUNT(CASE WHEN ut.status IN ('uploading', 'merging') THEN 1 END) as active_tasks,
    COALESCE(SUM(CASE WHEN ut.status = 'completed' THEN ut.file_size END), 0) as total_uploaded_bytes,
    COALESCE(AVG(CASE WHEN ut.status = 'completed' THEN ut.upload_speed END), 0) as avg_upload_speed,
    MAX(ut.created_at) as last_upload_time
FROM users u
LEFT JOIN upload_tasks ut ON u.id = ut.user_id
WHERE u.deleted_at IS NULL
GROUP BY u.id, u.username;

-- 创建上传性能监控视图
CREATE VIEW upload_performance AS
SELECT 
    DATE(ut.created_at) as upload_date,
    COUNT(*) as daily_uploads,
    COUNT(CASE WHEN ut.status = 'completed' THEN 1 END) as completed_uploads,
    COUNT(CASE WHEN ut.status = 'failed' THEN 1 END) as failed_uploads,
    ROUND(COUNT(CASE WHEN ut.status = 'completed' THEN 1 END) * 100.0 / COUNT(*), 2) as success_rate,
    SUM(CASE WHEN ut.status = 'completed' THEN ut.file_size END) as total_bytes_uploaded,
    AVG(CASE WHEN ut.status = 'completed' THEN ut.upload_speed END) as avg_speed,
    AVG(CASE WHEN ut.status = 'completed' THEN ut.actual_duration END) as avg_duration
FROM upload_tasks ut
WHERE ut.created_at >= DATE_SUB(CURRENT_DATE, INTERVAL 30 DAY)
GROUP BY DATE(ut.created_at)
ORDER BY upload_date DESC;