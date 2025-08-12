-- +migrate Up
-- 创建迁移: 文件分片表
-- 版本: 20250112100600
-- 描述: 创建文件分片信息管理表，支持分片上传、断点续传和分片状态跟踪
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:06:00
-- 依赖: 20250112100500_create_upload_tasks_table
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1

-- ============================================================================
-- 文件分片信息管理表 (file_chunks)
-- ============================================================================

CREATE TABLE file_chunks (
    -- 基础标识字段
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '分片记录ID',
    upload_id VARCHAR(64) NOT NULL COMMENT '关联的上传任务ID，引用upload_tasks表',
    chunk_number INT UNSIGNED NOT NULL COMMENT '分片序号，从1开始递增',
    
    -- 分片基础信息
    chunk_size INT UNSIGNED NOT NULL COMMENT '分片实际大小（字节）',
    chunk_offset BIGINT UNSIGNED NOT NULL COMMENT '分片在原文件中的偏移量（字节）',
    
    -- 分片完整性校验
    chunk_md5 VARCHAR(32) NOT NULL COMMENT '分片MD5哈希值，用于完整性验证',
    chunk_sha1 VARCHAR(40) DEFAULT NULL COMMENT '分片SHA1哈希值，额外验证',
    chunk_crc32 VARCHAR(8) DEFAULT NULL COMMENT '分片CRC32校验值，快速验证',
    
    -- 分片存储信息
    chunk_path VARCHAR(500) DEFAULT NULL COMMENT '分片临时存储路径',
    oss_bucket VARCHAR(100) DEFAULT NULL COMMENT '分片存储的OSS桶名',
    oss_key VARCHAR(500) DEFAULT NULL COMMENT '分片在OSS中的对象键',
    oss_etag VARCHAR(100) DEFAULT NULL COMMENT 'OSS ETag，上传成功的标识',
    
    -- 分片状态管理
    status ENUM('pending', 'preparing', 'uploading', 'completed', 'failed', 'cancelled', 'expired') DEFAULT 'pending' COMMENT '分片上传状态',
    
    -- 上传控制信息
    upload_url VARCHAR(1000) DEFAULT NULL COMMENT '分片上传的预签名URL',
    upload_method ENUM('PUT', 'POST') DEFAULT 'PUT' COMMENT '上传HTTP方法',
    upload_headers JSON DEFAULT NULL COMMENT '上传需要的HTTP头信息',
    
    -- 重试和错误处理
    retry_count INT UNSIGNED DEFAULT 0 COMMENT '当前重试次数',
    max_retries INT UNSIGNED DEFAULT 3 COMMENT '最大重试次数',
    last_error TEXT DEFAULT NULL COMMENT '最后一次错误信息',
    error_code VARCHAR(50) DEFAULT NULL COMMENT '错误代码',
    
    -- 性能统计
    upload_speed BIGINT UNSIGNED DEFAULT 0 COMMENT '分片上传速度（字节/秒）',
    upload_duration INT UNSIGNED DEFAULT 0 COMMENT '分片上传耗时（毫秒）',
    
    -- 并发控制
    worker_id VARCHAR(32) DEFAULT NULL COMMENT '处理该分片的工作进程ID',
    locked_at TIMESTAMP NULL DEFAULT NULL COMMENT '分片锁定时间，防止并发处理',
    lock_expires_at TIMESTAMP NULL DEFAULT NULL COMMENT '锁定过期时间',
    
    -- URL过期管理
    url_expires_at TIMESTAMP NULL DEFAULT NULL COMMENT '上传URL过期时间',
    url_generated_at TIMESTAMP NULL DEFAULT NULL COMMENT 'URL生成时间',
    
    -- 时间戳记录
    upload_started_at TIMESTAMP NULL DEFAULT NULL COMMENT '分片开始上传时间',
    uploaded_at TIMESTAMP NULL DEFAULT NULL COMMENT '分片上传完成时间',
    verified_at TIMESTAMP NULL DEFAULT NULL COMMENT '分片验证完成时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '分片记录创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '分片记录更新时间',
    
    -- 业务唯一约束：同一上传任务的分片序号必须唯一
    UNIQUE KEY uk_upload_chunk (upload_id, chunk_number),
    
    -- 基础业务索引
    INDEX idx_file_chunks_upload_id (upload_id) COMMENT '上传任务ID索引，查询任务的所有分片',
    INDEX idx_file_chunks_chunk_number (chunk_number) COMMENT '分片序号索引，排序和查找',
    INDEX idx_file_chunks_status (status) COMMENT '状态索引，按状态筛选分片',
    INDEX idx_file_chunks_chunk_md5 (chunk_md5) COMMENT 'MD5哈希索引，重复分片检测',
    INDEX idx_file_chunks_worker_id (worker_id) COMMENT '工作进程ID索引，并发控制',
    INDEX idx_file_chunks_created_at (created_at) COMMENT '创建时间索引，时间排序',
    INDEX idx_file_chunks_uploaded_at (uploaded_at) COMMENT '上传时间索引，完成状态查询',
    
    -- 时间相关索引
    INDEX idx_file_chunks_url_expires (url_expires_at) COMMENT 'URL过期时间索引，清理过期URL',
    INDEX idx_file_chunks_lock_expires (lock_expires_at) COMMENT '锁过期时间索引，清理过期锁',
    INDEX idx_file_chunks_locked_at (locked_at) COMMENT '锁定时间索引，并发控制',
    
    -- 复合业务索引
    INDEX idx_file_chunks_upload_status_chunk (upload_id, status, chunk_number) COMMENT '上传任务状态分片查询优化',
    INDEX idx_file_chunks_status_retry (status, retry_count, updated_at) COMMENT '重试任务查询优化',
    INDEX idx_file_chunks_upload_completed (upload_id, status, uploaded_at) COMMENT '完成分片统计优化',
    INDEX idx_file_chunks_worker_status (worker_id, status, locked_at) COMMENT '工作进程任务查询',
    
    -- 外键约束
    CONSTRAINT fk_file_chunks_upload_id 
        FOREIGN KEY (upload_id) REFERENCES upload_tasks(upload_id) 
        ON DELETE CASCADE  -- 上传任务删除时，级联删除所有分片
        ON UPDATE CASCADE
        
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='文件分片信息管理表 - 支持分片上传、断点续传和并发控制'
  ROW_FORMAT=DYNAMIC;

-- ============================================================================
-- 文件分片表约束和触发器
-- ============================================================================

-- 分片大小合理性约束
ALTER TABLE file_chunks ADD CONSTRAINT chk_chunk_size_positive 
CHECK (chunk_size > 0 AND chunk_size <= 104857600); -- 不超过100MB

-- 分片序号必须大于0
ALTER TABLE file_chunks ADD CONSTRAINT chk_chunk_number_positive 
CHECK (chunk_number > 0);

-- 偏移量不能为负
ALTER TABLE file_chunks ADD CONSTRAINT chk_chunk_offset_positive 
CHECK (chunk_offset >= 0);

-- 重试次数约束
ALTER TABLE file_chunks ADD CONSTRAINT chk_chunk_retry_count 
CHECK (retry_count >= 0 AND retry_count <= max_retries);

-- 最大重试次数约束
ALTER TABLE file_chunks ADD CONSTRAINT chk_chunk_max_retries 
CHECK (max_retries >= 0 AND max_retries <= 20);

-- 哈希值格式约束
ALTER TABLE file_chunks ADD CONSTRAINT chk_chunk_md5_format 
CHECK (chunk_md5 REGEXP '^[a-f0-9]{32}$');

ALTER TABLE file_chunks ADD CONSTRAINT chk_chunk_sha1_format 
CHECK (chunk_sha1 IS NULL OR chunk_sha1 REGEXP '^[a-f0-9]{40}$');

ALTER TABLE file_chunks ADD CONSTRAINT chk_chunk_crc32_format 
CHECK (chunk_crc32 IS NULL OR chunk_crc32 REGEXP '^[a-fA-F0-9]{8}$');

-- JSON字段验证
ALTER TABLE file_chunks ADD CONSTRAINT chk_chunk_json_valid 
CHECK (upload_headers IS NULL OR JSON_VALID(upload_headers));

-- 时间逻辑约束
ALTER TABLE file_chunks ADD CONSTRAINT chk_chunk_time_logic 
CHECK (
    (uploaded_at IS NULL OR uploaded_at >= created_at) AND
    (upload_started_at IS NULL OR upload_started_at >= created_at) AND
    (verified_at IS NULL OR verified_at >= created_at) AND
    (lock_expires_at IS NULL OR lock_expires_at >= locked_at)
);

-- 创建分片状态更新触发器
DELIMITER //
CREATE TRIGGER file_chunks_status_update 
BEFORE UPDATE ON file_chunks
FOR EACH ROW
BEGIN
    -- 记录状态变化时间
    IF OLD.status != NEW.status THEN
        CASE NEW.status
            WHEN 'uploading' THEN
                IF NEW.upload_started_at IS NULL THEN
                    SET NEW.upload_started_at = CURRENT_TIMESTAMP;
                END IF;
            WHEN 'completed' THEN
                IF NEW.uploaded_at IS NULL THEN
                    SET NEW.uploaded_at = CURRENT_TIMESTAMP;
                END IF;
                IF NEW.verified_at IS NULL THEN
                    SET NEW.verified_at = CURRENT_TIMESTAMP;
                END IF;
            WHEN 'failed' THEN
                SET NEW.retry_count = NEW.retry_count + 1;
        END CASE;
    END IF;
    
    -- 计算上传耗时
    IF OLD.status != 'completed' AND NEW.status = 'completed' 
       AND NEW.upload_started_at IS NOT NULL THEN
        SET NEW.upload_duration = TIMESTAMPDIFF(MICROSECOND, NEW.upload_started_at, CURRENT_TIMESTAMP) / 1000;
        
        -- 计算上传速度
        IF NEW.upload_duration > 0 THEN
            SET NEW.upload_speed = ROUND(NEW.chunk_size * 1000 / NEW.upload_duration);
        END IF;
    END IF;
    
    -- 清理锁定状态
    IF NEW.status IN ('completed', 'failed', 'cancelled') THEN
        SET NEW.locked_at = NULL;
        SET NEW.lock_expires_at = NULL;
        SET NEW.worker_id = NULL;
    END IF;
END//
DELIMITER ;

-- 创建分片锁定管理触发器
DELIMITER //
CREATE TRIGGER file_chunks_lock_update 
BEFORE UPDATE ON file_chunks
FOR EACH ROW
BEGIN
    -- 设置锁定时自动计算过期时间
    IF OLD.locked_at IS NULL AND NEW.locked_at IS NOT NULL THEN
        IF NEW.lock_expires_at IS NULL THEN
            SET NEW.lock_expires_at = DATE_ADD(NEW.locked_at, INTERVAL 30 MINUTE); -- 默认30分钟锁定
        END IF;
    END IF;
    
    -- 检查锁定是否过期
    IF NEW.locked_at IS NOT NULL AND NEW.lock_expires_at IS NOT NULL 
       AND NEW.lock_expires_at <= CURRENT_TIMESTAMP THEN
        -- 锁定已过期，清理锁定状态
        SET NEW.locked_at = NULL;
        SET NEW.lock_expires_at = NULL;
        SET NEW.worker_id = NULL;
        
        -- 如果状态是uploading但锁定过期，重置为pending
        IF NEW.status = 'uploading' THEN
            SET NEW.status = 'pending';
        END IF;
    END IF;
END//
DELIMITER ;

-- 创建分片完成后更新上传任务触发器
DELIMITER //
CREATE TRIGGER file_chunks_update_task 
AFTER UPDATE ON file_chunks
FOR EACH ROW
BEGIN
    -- 当分片状态变为completed时，更新上传任务的进度
    IF OLD.status != 'completed' AND NEW.status = 'completed' THEN
        UPDATE upload_tasks 
        SET uploaded_chunks = uploaded_chunks + 1,
            uploaded_size = uploaded_size + NEW.chunk_size,
            last_chunk_at = CURRENT_TIMESTAMP
        WHERE upload_id = NEW.upload_id;
        
        -- 如果是第一个分片，记录首片时间
        UPDATE upload_tasks 
        SET first_chunk_at = CURRENT_TIMESTAMP
        WHERE upload_id = NEW.upload_id 
          AND first_chunk_at IS NULL;
    END IF;
    
    -- 当分片状态变为failed时，更新失败计数
    IF OLD.status != 'failed' AND NEW.status = 'failed' THEN
        UPDATE upload_tasks 
        SET failed_chunks = failed_chunks + 1
        WHERE upload_id = NEW.upload_id;
    END IF;
END//
DELIMITER ;

-- ============================================================================
-- 分片管理存储过程
-- ============================================================================

-- 创建清理过期分片锁的存储过程
DELIMITER //
CREATE PROCEDURE CleanupExpiredChunkLocks()
BEGIN
    DECLARE affected_rows INT DEFAULT 0;
    
    -- 清理过期的分片锁定
    UPDATE file_chunks 
    SET locked_at = NULL,
        lock_expires_at = NULL,
        worker_id = NULL,
        status = CASE 
            WHEN status = 'uploading' THEN 'pending'
            ELSE status 
        END
    WHERE lock_expires_at IS NOT NULL 
      AND lock_expires_at <= CURRENT_TIMESTAMP;
    
    SET affected_rows = ROW_COUNT();
    
    SELECT CONCAT('已清理 ', affected_rows, ' 个过期的分片锁定') as result;
END//
DELIMITER ;

-- 创建获取上传任务分片状态的存储过程
DELIMITER //
CREATE PROCEDURE GetUploadChunkStatus(
    IN target_upload_id VARCHAR(64)
)
BEGIN
    SELECT 
        COUNT(*) as total_chunks,
        COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed_chunks,
        COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_chunks,
        COUNT(CASE WHEN status = 'uploading' THEN 1 END) as uploading_chunks,
        COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_chunks,
        COALESCE(SUM(CASE WHEN status = 'completed' THEN chunk_size END), 0) as completed_size,
        COALESCE(AVG(CASE WHEN status = 'completed' THEN upload_speed END), 0) as avg_upload_speed,
        MIN(chunk_number) as first_chunk,
        MAX(chunk_number) as last_chunk,
        MAX(uploaded_at) as last_upload_time
    FROM file_chunks 
    WHERE upload_id = target_upload_id;
END//
DELIMITER ;

-- 创建批量重置失败分片的存储过程
DELIMITER //
CREATE PROCEDURE ResetFailedChunks(
    IN target_upload_id VARCHAR(64),
    IN max_retry_exceeded BOOLEAN DEFAULT FALSE
)
BEGIN
    DECLARE affected_rows INT DEFAULT 0;
    
    IF max_retry_exceeded THEN
        -- 重置所有失败的分片，包括超过最大重试次数的
        UPDATE file_chunks 
        SET status = 'pending',
            retry_count = 0,
            last_error = NULL,
            error_code = NULL,
            locked_at = NULL,
            lock_expires_at = NULL,
            worker_id = NULL
        WHERE upload_id = target_upload_id 
          AND status = 'failed';
    ELSE
        -- 只重置未超过最大重试次数的失败分片
        UPDATE file_chunks 
        SET status = 'pending',
            last_error = NULL,
            error_code = NULL,
            locked_at = NULL,
            lock_expires_at = NULL,
            worker_id = NULL
        WHERE upload_id = target_upload_id 
          AND status = 'failed'
          AND retry_count < max_retries;
    END IF;
    
    SET affected_rows = ROW_COUNT();
    
    SELECT CONCAT('已重置 ', affected_rows, ' 个失败分片') as result;
END//
DELIMITER ;

-- ============================================================================
-- 分片统计视图
-- ============================================================================

-- 创建分片状态统计视图
CREATE VIEW chunk_status_stats AS
SELECT 
    ut.upload_id,
    ut.user_id,
    ut.filename,
    ut.status as task_status,
    COUNT(fc.id) as total_chunks,
    COUNT(CASE WHEN fc.status = 'completed' THEN 1 END) as completed_chunks,
    COUNT(CASE WHEN fc.status = 'failed' THEN 1 END) as failed_chunks,
    COUNT(CASE WHEN fc.status = 'uploading' THEN 1 END) as uploading_chunks,
    COUNT(CASE WHEN fc.status = 'pending' THEN 1 END) as pending_chunks,
    ROUND(COUNT(CASE WHEN fc.status = 'completed' THEN 1 END) * 100.0 / COUNT(fc.id), 2) as completion_percentage,
    COALESCE(SUM(CASE WHEN fc.status = 'completed' THEN fc.chunk_size END), 0) as uploaded_bytes,
    COALESCE(AVG(CASE WHEN fc.status = 'completed' THEN fc.upload_speed END), 0) as avg_chunk_speed
FROM upload_tasks ut
LEFT JOIN file_chunks fc ON ut.upload_id = fc.upload_id
GROUP BY ut.upload_id, ut.user_id, ut.filename, ut.status;

-- 创建分片性能分析视图
CREATE VIEW chunk_performance AS
SELECT 
    DATE(fc.created_at) as chunk_date,
    COUNT(*) as total_chunks_created,
    COUNT(CASE WHEN fc.status = 'completed' THEN 1 END) as completed_chunks,
    COUNT(CASE WHEN fc.status = 'failed' THEN 1 END) as failed_chunks,
    ROUND(COUNT(CASE WHEN fc.status = 'completed' THEN 1 END) * 100.0 / COUNT(*), 2) as success_rate,
    COALESCE(AVG(CASE WHEN fc.status = 'completed' THEN fc.upload_speed END), 0) as avg_speed,
    COALESCE(AVG(CASE WHEN fc.status = 'completed' THEN fc.upload_duration END), 0) as avg_duration,
    SUM(CASE WHEN fc.status = 'completed' THEN fc.chunk_size END) as total_bytes_uploaded
FROM file_chunks fc
WHERE fc.created_at >= DATE_SUB(CURRENT_DATE, INTERVAL 7 DAY)
GROUP BY DATE(fc.created_at)
ORDER BY chunk_date DESC;