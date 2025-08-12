-- +migrate Up
-- 创建迁移: 文件表
-- 版本: 20250112100300
-- 描述: 创建文件和文件夹统一管理表，支持层级结构、版本控制、多媒体属性和云存储集成
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:03:00
-- 依赖: 20250112100000_create_users_table
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1

-- ============================================================================
-- 文件和文件夹统一管理表 (files)
-- ============================================================================

CREATE TABLE files (
    -- 基础标识字段
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '文件/文件夹唯一标识',
    filename VARCHAR(255) NOT NULL COMMENT '文件名或文件夹名，支持Unicode字符',
    file_path VARCHAR(1000) NOT NULL COMMENT '文件完整路径，从根目录开始的绝对路径',
    
    -- 基础属性
    file_size BIGINT UNSIGNED DEFAULT 0 COMMENT '文件大小（字节），文件夹为0',
    file_type VARCHAR(50) DEFAULT NULL COMMENT '文件类型分类：document、image、video、audio、archive、code等',
    mime_type VARCHAR(100) DEFAULT NULL COMMENT '标准MIME类型，如image/jpeg、application/pdf',
    file_extension VARCHAR(20) DEFAULT NULL COMMENT '文件扩展名，统一小写存储',
    
    -- 文件完整性校验
    md5_hash VARCHAR(32) DEFAULT NULL COMMENT '文件MD5哈希值，用于去重和完整性验证',
    sha256_hash VARCHAR(64) DEFAULT NULL COMMENT '文件SHA256哈希值，增强安全性验证',
    
    -- 所有权和层级关系
    user_id BIGINT UNSIGNED NOT NULL COMMENT '文件所有者ID，关联users表',
    parent_id BIGINT UNSIGNED DEFAULT NULL COMMENT '父文件夹ID，NULL表示根目录，自引用files表',
    
    -- 文件状态标识
    is_folder BOOLEAN DEFAULT FALSE COMMENT '是否为文件夹，true=文件夹，false=文件',
    is_deleted BOOLEAN DEFAULT FALSE COMMENT '软删除标记，支持回收站功能',
    is_favorite BOOLEAN DEFAULT FALSE COMMENT '是否收藏，用户快捷访问',
    is_locked BOOLEAN DEFAULT FALSE COMMENT '是否锁定，禁止修改和删除',
    is_encrypted BOOLEAN DEFAULT FALSE COMMENT '是否启用客户端加密存储',
    is_public BOOLEAN DEFAULT FALSE COMMENT '是否为公开文件，无需认证即可访问',
    
    -- 加密相关
    encryption_key_id VARCHAR(64) DEFAULT NULL COMMENT '加密密钥标识符，关联密钥管理系统',
    encryption_algorithm VARCHAR(20) DEFAULT NULL COMMENT '加密算法：AES-256-GCM、ChaCha20-Poly1305',
    
    -- 版本控制
    version INT UNSIGNED DEFAULT 1 COMMENT '当前版本号，从1开始递增',
    
    -- 统计计数器
    download_count INT UNSIGNED DEFAULT 0 COMMENT '下载次数统计',
    view_count INT UNSIGNED DEFAULT 0 COMMENT '查看/预览次数统计',
    share_count INT UNSIGNED DEFAULT 0 COMMENT '分享次数统计',
    comment_count INT UNSIGNED DEFAULT 0 COMMENT '评论数量统计',
    like_count INT UNSIGNED DEFAULT 0 COMMENT '点赞数量统计',
    
    -- 软删除相关
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间，软删除标记',
    original_parent_id BIGINT UNSIGNED DEFAULT NULL COMMENT '删除前的原始父文件夹ID，用于恢复',
    
    -- 云存储集成 (阿里云OSS)
    oss_bucket VARCHAR(100) DEFAULT NULL COMMENT 'OSS存储桶名称',
    oss_key VARCHAR(500) DEFAULT NULL COMMENT 'OSS对象键，唯一标识存储对象',
    oss_url VARCHAR(1000) DEFAULT NULL COMMENT 'OSS完整访问URL',
    oss_region VARCHAR(50) DEFAULT NULL COMMENT 'OSS存储区域',
    
    -- 多媒体预览URL
    thumbnail_url VARCHAR(1000) DEFAULT NULL COMMENT '缩略图URL，支持图片和视频封面',
    preview_url VARCHAR(1000) DEFAULT NULL COMMENT '预览文件URL，如PDF转图片',
    video_cover_url VARCHAR(1000) DEFAULT NULL COMMENT '视频封面图URL',
    
    -- 多媒体属性
    duration INT UNSIGNED DEFAULT NULL COMMENT '音视频时长（秒）',
    width INT UNSIGNED DEFAULT NULL COMMENT '图片/视频宽度（像素）',
    height INT UNSIGNED DEFAULT NULL COMMENT '图片/视频高度（像素）',
    bit_rate INT UNSIGNED DEFAULT NULL COMMENT '音视频比特率（bps）',
    frame_rate DECIMAL(8,3) DEFAULT NULL COMMENT '视频帧率（fps）',
    audio_channels TINYINT UNSIGNED DEFAULT NULL COMMENT '音频声道数',
    sample_rate INT UNSIGNED DEFAULT NULL COMMENT '音频采样率（Hz）',
    
    -- JSON扩展属性
    exif_data JSON DEFAULT NULL COMMENT '图片EXIF元数据，包含拍摄信息、GPS坐标等',
    metadata JSON DEFAULT NULL COMMENT '文件元数据，包含自定义属性和扩展信息',
    tags JSON DEFAULT NULL COMMENT '文件标签数组，支持多标签分类',
    
    -- 文件描述和备注
    description TEXT DEFAULT NULL COMMENT '文件描述和备注信息',
    
    -- 安全扫描
    virus_scan_status ENUM('pending', 'scanning', 'clean', 'infected', 'failed', 'skipped') DEFAULT 'pending' COMMENT '病毒扫描状态',
    virus_scan_result JSON DEFAULT NULL COMMENT '病毒扫描结果详情，包含威胁类型和处理建议',
    virus_scanned_at TIMESTAMP NULL DEFAULT NULL COMMENT '最后一次病毒扫描时间',
    
    -- 访问控制
    access_level ENUM('private', 'shared', 'public', 'restricted') DEFAULT 'private' COMMENT '访问级别控制',
    download_permission ENUM('owner', 'shared_users', 'anyone') DEFAULT 'owner' COMMENT '下载权限控制',
    
    -- 时间戳
    last_accessed_at TIMESTAMP NULL DEFAULT NULL COMMENT '最后访问时间，用于热度统计',
    content_updated_at TIMESTAMP NULL DEFAULT NULL COMMENT '内容最后更新时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    
    -- 索引策略
    -- 主要查询索引
    INDEX idx_files_user_parent_deleted (user_id, parent_id, is_deleted) COMMENT '用户文件树查询的核心索引',
    INDEX idx_files_filename (filename) COMMENT '文件名搜索索引',
    INDEX idx_files_file_path (file_path(255)) COMMENT '路径查询索引，限制255字符',
    
    -- 分类和类型索引
    INDEX idx_files_file_type (file_type) COMMENT '文件类型筛选索引',
    INDEX idx_files_mime_type (mime_type) COMMENT 'MIME类型查询索引',
    INDEX idx_files_file_extension (file_extension) COMMENT '文件扩展名索引',
    
    -- 完整性和去重索引
    INDEX idx_files_md5_hash (md5_hash) COMMENT 'MD5哈希查询索引，用于去重',
    INDEX idx_files_sha256_hash (sha256_hash) COMMENT 'SHA256哈希查询索引',
    
    -- 状态筛选索引
    INDEX idx_files_is_folder (is_folder) COMMENT '文件夹筛选索引',
    INDEX idx_files_is_deleted (is_deleted) COMMENT '软删除筛选索引',
    INDEX idx_files_is_favorite (is_favorite) COMMENT '收藏文件索引',
    INDEX idx_files_is_public (is_public) COMMENT '公开文件索引',
    
    -- 时间相关索引
    INDEX idx_files_deleted_at (deleted_at) COMMENT '删除时间索引，回收站查询',
    INDEX idx_files_created_at (created_at) COMMENT '创建时间索引，时间排序',
    INDEX idx_files_updated_at (updated_at) COMMENT '更新时间索引',
    INDEX idx_files_last_accessed_at (last_accessed_at) COMMENT '访问时间索引，热度统计',
    
    -- 大小和统计索引
    INDEX idx_files_file_size (file_size) COMMENT '文件大小索引，大小排序和统计',
    INDEX idx_files_download_count (download_count) COMMENT '下载次数索引，热门文件',
    INDEX idx_files_view_count (view_count) COMMENT '查看次数索引',
    
    -- 安全相关索引
    INDEX idx_files_virus_scan_status (virus_scan_status) COMMENT '病毒扫描状态索引',
    INDEX idx_files_access_level (access_level) COMMENT '访问级别索引',
    
    -- 复合业务索引
    INDEX idx_files_user_type_deleted (user_id, file_type, is_deleted) COMMENT '用户文件类型筛选索引',
    INDEX idx_files_parent_folder_name (parent_id, is_folder, filename) COMMENT '文件夹内容浏览优化索引',
    INDEX idx_files_user_favorite (user_id, is_favorite, updated_at) COMMENT '用户收藏文件索引',
    
    -- 全文搜索索引 (MySQL 8.0 ngram支持中文)
    FULLTEXT idx_fulltext_search (filename, description) WITH PARSER ngram COMMENT '文件名和描述全文搜索索引',
    
    -- 外键约束
    CONSTRAINT fk_files_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_files_parent_id 
        FOREIGN KEY (parent_id) REFERENCES files(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE
        
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='文件和文件夹统一管理表 - 支持层级结构、版本控制、多媒体属性和云存储集成'
  ROW_FORMAT=DYNAMIC;

-- ============================================================================
-- 文件表约束和触发器
-- ============================================================================

-- 文件路径格式约束
ALTER TABLE files ADD CONSTRAINT chk_file_path_format 
CHECK (
    file_path REGEXP '^/([^/\0]+/)*[^/\0]*$' OR file_path = '/'
);

-- 文件大小合理性约束
ALTER TABLE files ADD CONSTRAINT chk_file_size_range 
CHECK (
    file_size >= 0 AND 
    (is_folder = TRUE OR file_size <= 107374182400) -- 文件夹或文件不超过100GB
);

-- 哈希值格式约束
ALTER TABLE files ADD CONSTRAINT chk_md5_format 
CHECK (md5_hash IS NULL OR md5_hash REGEXP '^[a-f0-9]{32}$');

ALTER TABLE files ADD CONSTRAINT chk_sha256_format 
CHECK (sha256_hash IS NULL OR sha256_hash REGEXP '^[a-f0-9]{64}$');

-- 文件扩展名格式约束
ALTER TABLE files ADD CONSTRAINT chk_extension_format 
CHECK (file_extension IS NULL OR file_extension REGEXP '^[a-z0-9]{1,20}$');

-- 多媒体属性合理性约束
ALTER TABLE files ADD CONSTRAINT chk_media_dimensions 
CHECK (
    (width IS NULL AND height IS NULL) OR 
    (width > 0 AND width <= 65535 AND height > 0 AND height <= 65535)
);

ALTER TABLE files ADD CONSTRAINT chk_duration_range 
CHECK (duration IS NULL OR (duration >= 0 AND duration <= 86400)); -- 最长24小时

ALTER TABLE files ADD CONSTRAINT chk_frame_rate_range 
CHECK (frame_rate IS NULL OR (frame_rate > 0 AND frame_rate <= 240));

-- JSON字段格式验证
ALTER TABLE files ADD CONSTRAINT chk_json_fields_valid 
CHECK (
    (exif_data IS NULL OR JSON_VALID(exif_data)) AND
    (metadata IS NULL OR JSON_VALID(metadata)) AND
    (tags IS NULL OR JSON_VALID(tags)) AND
    (virus_scan_result IS NULL OR JSON_VALID(virus_scan_result))
);

-- 文件夹逻辑约束
ALTER TABLE files ADD CONSTRAINT chk_folder_logic 
CHECK (
    (is_folder = FALSE) OR 
    (is_folder = TRUE AND file_size = 0 AND mime_type IS NULL AND file_extension IS NULL)
);

-- 创建文件统计更新触发器
DELIMITER //
CREATE TRIGGER files_stats_update 
AFTER INSERT ON files
FOR EACH ROW
BEGIN
    -- 更新用户存储使用量
    IF NEW.is_folder = FALSE AND NEW.is_deleted = FALSE THEN
        UPDATE users 
        SET storage_used = storage_used + NEW.file_size 
        WHERE id = NEW.user_id;
    END IF;
    
    -- 如果是新文件夹，确保路径正确
    IF NEW.is_folder = TRUE AND NEW.parent_id IS NOT NULL THEN
        UPDATE files 
        SET file_path = CONCAT(
            (SELECT file_path FROM files WHERE id = NEW.parent_id),
            CASE WHEN RIGHT((SELECT file_path FROM files WHERE id = NEW.parent_id), 1) = '/' 
                 THEN '' ELSE '/' END,
            NEW.filename,
            '/'
        )
        WHERE id = NEW.id;
    END IF;
END//
DELIMITER ;

-- 创建文件删除统计触发器
DELIMITER //
CREATE TRIGGER files_stats_delete 
AFTER UPDATE ON files
FOR EACH ROW
BEGIN
    -- 处理软删除对存储使用量的影响
    IF OLD.is_deleted = FALSE AND NEW.is_deleted = TRUE AND NEW.is_folder = FALSE THEN
        -- 文件被删除，减少存储使用量
        UPDATE users 
        SET storage_used = storage_used - NEW.file_size 
        WHERE id = NEW.user_id;
        
        -- 记录删除前的父文件夹
        IF NEW.original_parent_id IS NULL THEN
            UPDATE files 
            SET original_parent_id = OLD.parent_id 
            WHERE id = NEW.id;
        END IF;
        
    ELSEIF OLD.is_deleted = TRUE AND NEW.is_deleted = FALSE AND NEW.is_folder = FALSE THEN
        -- 文件被恢复，增加存储使用量
        UPDATE users 
        SET storage_used = storage_used + NEW.file_size 
        WHERE id = NEW.user_id;
    END IF;
END//
DELIMITER ;

-- ============================================================================
-- 初始化根文件夹
-- ============================================================================

-- 为每个用户创建根文件夹
INSERT INTO files (
    filename, 
    file_path, 
    user_id, 
    parent_id, 
    is_folder, 
    file_type,
    description,
    created_at
)
SELECT 
    'ROOT' as filename,
    '/' as file_path,
    id as user_id,
    NULL as parent_id,
    TRUE as is_folder,
    'folder' as file_type,
    '用户根目录' as description,
    CURRENT_TIMESTAMP as created_at
FROM users 
WHERE deleted_at IS NULL;

-- ============================================================================
-- 文件管理视图
-- ============================================================================

-- 创建文件统计视图
CREATE VIEW file_stats AS
SELECT 
    u.id as user_id,
    u.username,
    COUNT(f.id) as total_files,
    COUNT(CASE WHEN f.is_folder = TRUE THEN 1 END) as folder_count,
    COUNT(CASE WHEN f.is_folder = FALSE THEN 1 END) as file_count,
    COUNT(CASE WHEN f.is_deleted = TRUE THEN 1 END) as deleted_count,
    COUNT(CASE WHEN f.is_favorite = TRUE THEN 1 END) as favorite_count,
    COALESCE(SUM(CASE WHEN f.is_folder = FALSE AND f.is_deleted = FALSE THEN f.file_size END), 0) as total_size,
    MAX(f.created_at) as last_upload_time
FROM users u
LEFT JOIN files f ON u.id = f.user_id
WHERE u.deleted_at IS NULL
GROUP BY u.id, u.username;

-- 创建热门文件视图
CREATE VIEW popular_files AS
SELECT 
    f.*,
    u.username as owner_name,
    (f.download_count * 3 + f.view_count * 2 + f.share_count * 5 + f.like_count * 1) as popularity_score
FROM files f
JOIN users u ON f.user_id = u.id
WHERE f.is_folder = FALSE 
  AND f.is_deleted = FALSE 
  AND f.access_level IN ('public', 'shared')
ORDER BY popularity_score DESC;