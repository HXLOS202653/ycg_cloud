-- +migrate Up
-- 创建迁移: 文件版本表
-- 版本: 20250112100400
-- 描述: 创建文件版本历史管理表，支持版本控制、差异对比和里程碑管理
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:04:00
-- 依赖: 20250112100300_create_files_table
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1

-- ============================================================================
-- 文件版本历史管理表 (file_versions)
-- ============================================================================

CREATE TABLE file_versions (
    -- 基础标识字段
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '版本记录ID',
    file_id BIGINT UNSIGNED NOT NULL COMMENT '关联的文件ID，引用files表',
    version_number INT UNSIGNED NOT NULL COMMENT '版本号，从1开始递增',
    
    -- 版本文件基础信息
    filename VARCHAR(255) NOT NULL COMMENT '该版本的文件名，可能与当前文件名不同',
    file_size BIGINT UNSIGNED NOT NULL COMMENT '该版本文件大小（字节）',
    file_extension VARCHAR(20) DEFAULT NULL COMMENT '该版本文件扩展名',
    mime_type VARCHAR(100) DEFAULT NULL COMMENT '该版本文件MIME类型',
    
    -- 文件完整性校验
    md5_hash VARCHAR(32) NOT NULL COMMENT '该版本文件MD5哈希值，唯一标识',
    sha256_hash VARCHAR(64) DEFAULT NULL COMMENT '该版本文件SHA256哈希值',
    
    -- 云存储信息
    oss_bucket VARCHAR(100) DEFAULT NULL COMMENT 'OSS存储桶名称',
    oss_key VARCHAR(500) NOT NULL COMMENT 'OSS对象键，该版本的存储路径',
    oss_url VARCHAR(1000) NOT NULL COMMENT 'OSS访问URL，该版本的下载地址',
    oss_etag VARCHAR(100) DEFAULT NULL COMMENT 'OSS ETag，用于验证上传完整性',
    
    -- 版本控制属性
    is_milestone BOOLEAN DEFAULT FALSE COMMENT '是否为里程碑版本，重要版本标记',
    is_active BOOLEAN DEFAULT TRUE COMMENT '版本是否有效，无效版本不可下载',
    is_draft BOOLEAN DEFAULT FALSE COMMENT '是否为草稿版本，未正式发布',
    
    -- 版本变更信息
    change_description TEXT DEFAULT NULL COMMENT '版本变更说明，描述此版本的修改内容',
    change_type ENUM('create', 'update', 'rename', 'move', 'restore') DEFAULT 'update' COMMENT '变更类型',
    
    -- 版本创建者
    created_by BIGINT UNSIGNED NOT NULL COMMENT '版本创建者用户ID',
    
    -- 差异和对比信息
    file_diff JSON DEFAULT NULL COMMENT '与上一版本的差异信息，包含变更统计',
    previous_version_id BIGINT UNSIGNED DEFAULT NULL COMMENT '上一个版本的ID，用于版本链',
    
    -- 版本元数据
    version_metadata JSON DEFAULT NULL COMMENT '版本相关元数据，如编辑器信息、设备信息等',
    compression_ratio DECIMAL(5,4) DEFAULT NULL COMMENT '压缩比率，原始大小/存储大小',
    
    -- 版本生命周期
    expires_at TIMESTAMP NULL DEFAULT NULL COMMENT '版本过期时间，过期后可被清理',
    archived_at TIMESTAMP NULL DEFAULT NULL COMMENT '版本归档时间',
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '版本创建时间',
    
    -- 唯一约束：同一文件的版本号必须唯一
    UNIQUE KEY uk_file_version (file_id, version_number),
    
    -- 业务索引
    INDEX idx_file_versions_file_id (file_id) COMMENT '文件ID索引，查询文件的所有版本',
    INDEX idx_file_versions_version_number (version_number) COMMENT '版本号索引，版本排序',
    INDEX idx_file_versions_created_by (created_by) COMMENT '创建者索引，查询用户创建的版本',
    INDEX idx_file_versions_is_milestone (is_milestone) COMMENT '里程碑版本索引',
    INDEX idx_file_versions_is_active (is_active) COMMENT '有效版本索引',
    INDEX idx_file_versions_created_at (created_at) COMMENT '创建时间索引，时间排序',
    INDEX idx_file_versions_md5_hash (md5_hash) COMMENT 'MD5哈希索引，去重和查找',
    INDEX idx_file_versions_change_type (change_type) COMMENT '变更类型索引',
    INDEX idx_file_versions_expires_at (expires_at) COMMENT '过期时间索引，清理任务使用',
    
    -- 复合索引
    INDEX idx_file_versions_file_active_version (file_id, is_active, version_number DESC) COMMENT '文件有效版本查询优化',
    INDEX idx_file_versions_file_milestone (file_id, is_milestone, created_at DESC) COMMENT '文件里程碑版本查询',
    INDEX idx_file_versions_user_versions (created_by, created_at DESC) COMMENT '用户版本历史查询',
    
    -- 外键约束
    CONSTRAINT fk_file_versions_file_id 
        FOREIGN KEY (file_id) REFERENCES files(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_file_versions_created_by 
        FOREIGN KEY (created_by) REFERENCES users(id) 
        ON DELETE RESTRICT  -- 限制删除：不允许删除有版本记录的用户
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_file_versions_previous 
        FOREIGN KEY (previous_version_id) REFERENCES file_versions(id) 
        ON DELETE SET NULL  -- 如果上一版本被删除，设为NULL
        ON UPDATE CASCADE
        
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='文件版本历史管理表 - 支持完整的版本控制、差异对比和里程碑管理'
  ROW_FORMAT=DYNAMIC;

-- ============================================================================
-- 文件版本表约束和触发器
-- ============================================================================

-- 文件大小合理性约束
ALTER TABLE file_versions ADD CONSTRAINT chk_version_file_size 
CHECK (file_size > 0 AND file_size <= 107374182400); -- 不超过100GB

-- 版本号必须大于0
ALTER TABLE file_versions ADD CONSTRAINT chk_version_number_positive 
CHECK (version_number > 0);

-- 哈希值格式约束
ALTER TABLE file_versions ADD CONSTRAINT chk_version_md5_format 
CHECK (md5_hash REGEXP '^[a-f0-9]{32}$');

ALTER TABLE file_versions ADD CONSTRAINT chk_version_sha256_format 
CHECK (sha256_hash IS NULL OR sha256_hash REGEXP '^[a-f0-9]{64}$');

-- 压缩比率合理性约束
ALTER TABLE file_versions ADD CONSTRAINT chk_compression_ratio 
CHECK (compression_ratio IS NULL OR (compression_ratio > 0 AND compression_ratio <= 10));

-- JSON字段验证
ALTER TABLE file_versions ADD CONSTRAINT chk_version_json_valid 
CHECK (
    (file_diff IS NULL OR JSON_VALID(file_diff)) AND
    (version_metadata IS NULL OR JSON_VALID(version_metadata))
);

-- OSS键和URL不能为空约束
ALTER TABLE file_versions ADD CONSTRAINT chk_oss_info_not_empty 
CHECK (
    LENGTH(TRIM(oss_key)) > 0 AND 
    LENGTH(TRIM(oss_url)) > 0
);

-- 创建版本自动编号触发器
DELIMITER //
CREATE TRIGGER file_versions_auto_version 
BEFORE INSERT ON file_versions
FOR EACH ROW
BEGIN
    -- 如果版本号为0或NULL，自动生成下一个版本号
    IF NEW.version_number IS NULL OR NEW.version_number = 0 THEN
        SELECT COALESCE(MAX(version_number), 0) + 1 
        INTO NEW.version_number 
        FROM file_versions 
        WHERE file_id = NEW.file_id;
    END IF;
    
    -- 设置上一版本ID
    IF NEW.previous_version_id IS NULL AND NEW.version_number > 1 THEN
        SELECT id INTO NEW.previous_version_id 
        FROM file_versions 
        WHERE file_id = NEW.file_id 
          AND version_number = NEW.version_number - 1 
          AND is_active = TRUE
        LIMIT 1;
    END IF;
    
    -- 如果是第一个版本，自动标记为里程碑
    IF NEW.version_number = 1 THEN
        SET NEW.is_milestone = TRUE;
        SET NEW.change_type = 'create';
        SET NEW.change_description = COALESCE(NEW.change_description, '文件初始版本');
    END IF;
END//
DELIMITER ;

-- 创建版本更新文件表触发器
DELIMITER //
CREATE TRIGGER file_versions_update_files 
AFTER INSERT ON file_versions
FOR EACH ROW
BEGIN
    -- 更新files表中的版本号
    UPDATE files 
    SET version = NEW.version_number,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.file_id;
    
    -- 如果这是一个有效的新版本，停用之前的草稿版本
    IF NEW.is_active = TRUE AND NEW.is_draft = FALSE THEN
        UPDATE file_versions 
        SET is_active = FALSE 
        WHERE file_id = NEW.file_id 
          AND id != NEW.id 
          AND is_draft = TRUE;
    END IF;
END//
DELIMITER ;

-- 创建版本清理触发器
DELIMITER //
CREATE TRIGGER file_versions_cleanup 
AFTER UPDATE ON file_versions
FOR EACH ROW
BEGIN
    -- 如果版本被标记为过期且非里程碑版本，自动归档
    IF OLD.expires_at IS NULL AND NEW.expires_at IS NOT NULL 
       AND NEW.expires_at <= CURRENT_TIMESTAMP 
       AND NEW.is_milestone = FALSE THEN
        UPDATE file_versions 
        SET archived_at = CURRENT_TIMESTAMP,
            is_active = FALSE
        WHERE id = NEW.id;
    END IF;
END//
DELIMITER ;

-- ============================================================================
-- 版本管理存储过程
-- ============================================================================

-- 创建版本清理存储过程
DELIMITER //
CREATE PROCEDURE CleanupExpiredVersions(
    IN retention_days INT DEFAULT 90,
    IN keep_milestone BOOLEAN DEFAULT TRUE
)
BEGIN
    DECLARE affected_rows INT DEFAULT 0;
    DECLARE cleanup_date TIMESTAMP DEFAULT DATE_SUB(CURRENT_TIMESTAMP, INTERVAL retention_days DAY);
    
    -- 标记过期的非里程碑版本
    UPDATE file_versions 
    SET expires_at = CURRENT_TIMESTAMP,
        archived_at = CURRENT_TIMESTAMP,
        is_active = FALSE
    WHERE created_at < cleanup_date
      AND is_active = TRUE
      AND (keep_milestone = FALSE OR is_milestone = FALSE)
      AND expires_at IS NULL;
    
    SET affected_rows = ROW_COUNT();
    
    SELECT CONCAT('已标记 ', affected_rows, ' 个版本为过期') as result;
END//
DELIMITER ;

-- 创建获取文件版本历史的存储过程
DELIMITER //
CREATE PROCEDURE GetFileVersionHistory(
    IN target_file_id BIGINT UNSIGNED,
    IN include_inactive BOOLEAN DEFAULT FALSE
)
BEGIN
    SELECT 
        fv.*,
        u.username as created_by_name,
        LAG(fv.file_size) OVER (ORDER BY fv.version_number) as previous_size,
        fv.file_size - LAG(fv.file_size) OVER (ORDER BY fv.version_number) as size_change
    FROM file_versions fv
    LEFT JOIN users u ON fv.created_by = u.id
    WHERE fv.file_id = target_file_id
      AND (include_inactive = TRUE OR fv.is_active = TRUE)
    ORDER BY fv.version_number DESC;
END//
DELIMITER ;

-- ============================================================================
-- 版本统计视图
-- ============================================================================

-- 创建文件版本统计视图
CREATE VIEW file_version_stats AS
SELECT 
    f.id as file_id,
    f.filename,
    f.user_id,
    u.username,
    COUNT(fv.id) as total_versions,
    COUNT(CASE WHEN fv.is_milestone = TRUE THEN 1 END) as milestone_count,
    COUNT(CASE WHEN fv.is_active = TRUE THEN 1 END) as active_versions,
    MAX(fv.version_number) as latest_version,
    SUM(fv.file_size) as total_version_size,
    MIN(fv.created_at) as first_version_time,
    MAX(fv.created_at) as latest_version_time
FROM files f
LEFT JOIN file_versions fv ON f.id = fv.file_id
LEFT JOIN users u ON f.user_id = u.id
WHERE f.is_folder = FALSE AND f.is_deleted = FALSE
GROUP BY f.id, f.filename, f.user_id, u.username;

-- 创建版本变更活跃度视图
CREATE VIEW version_activity AS
SELECT 
    DATE(fv.created_at) as activity_date,
    COUNT(*) as versions_created,
    COUNT(DISTINCT fv.file_id) as files_updated,
    COUNT(DISTINCT fv.created_by) as active_users,
    SUM(fv.file_size) as total_size_versioned,
    COUNT(CASE WHEN fv.is_milestone = TRUE THEN 1 END) as milestones_created
FROM file_versions fv
WHERE fv.created_at >= DATE_SUB(CURRENT_DATE, INTERVAL 30 DAY)
GROUP BY DATE(fv.created_at)
ORDER BY activity_date DESC;

-- ============================================================================
-- 补充约束检查
-- ============================================================================

-- 版本号范围约束
ALTER TABLE file_versions ADD CONSTRAINT chk_version_number_range 
CHECK (version_number > 0 AND version_number <= 10000);

-- 哈希值格式增强约束
ALTER TABLE file_versions ADD CONSTRAINT chk_version_hash_enhanced 
CHECK (
    md5_hash REGEXP '^[a-f0-9]{32}$' AND
    (sha256_hash IS NULL OR sha256_hash REGEXP '^[a-f0-9]{64}$')
);

-- 变更类型与摘要一致性约束
ALTER TABLE file_versions ADD CONSTRAINT chk_change_type_summary_consistency 
CHECK (
    (change_type IN ('create', 'update', 'rename', 'move') AND change_summary IS NOT NULL AND LENGTH(TRIM(change_summary)) > 0) OR
    (change_type = 'delete')
);