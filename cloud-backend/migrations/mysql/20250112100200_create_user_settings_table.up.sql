-- +migrate Up
-- 创建迁移: 用户设置表
-- 版本: 20250112100200
-- 描述: 创建用户个性化设置表，管理同步、预览、通知等偏好设置
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:02:00
-- 依赖: 20250112100000_create_users_table
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1

-- ============================================================================
-- 用户个性化设置表 (user_settings)
-- ============================================================================

CREATE TABLE user_settings (
    -- 基础标识字段
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '设置ID',
    user_id BIGINT UNSIGNED NOT NULL UNIQUE COMMENT '用户ID，一对一关联users表',
    
    -- 同步设置
    sync_enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用文件同步功能',
    sync_bandwidth_limit INT UNSIGNED DEFAULT 5242880 COMMENT '同步带宽限制，默认5MB/s（5*1024*1024字节）',
    sync_only_wifi BOOLEAN DEFAULT FALSE COMMENT '仅在WiFi环境下同步（移动端）',
    sync_auto_pause BOOLEAN DEFAULT TRUE COMMENT '检测到大文件时是否自动暂停同步',
    
    -- 文件预览设置
    auto_preview BOOLEAN DEFAULT TRUE COMMENT '是否自动预览文件',
    preview_image_quality ENUM('low', 'medium', 'high', 'original') DEFAULT 'medium' COMMENT '图片预览质量',
    thumbnail_quality ENUM('low', 'medium', 'high') DEFAULT 'medium' COMMENT '缩略图质量',
    video_quality ENUM('360p', '720p', '1080p', '4k', 'auto') DEFAULT 'auto' COMMENT '视频播放质量',
    auto_play_video BOOLEAN DEFAULT FALSE COMMENT '是否自动播放视频',
    preload_thumbnails BOOLEAN DEFAULT TRUE COMMENT '是否预加载缩略图',
    
    -- 文件管理设置
    auto_cleanup_trash BOOLEAN DEFAULT TRUE COMMENT '是否自动清理回收站',
    trash_retention_days INT UNSIGNED DEFAULT 30 COMMENT '回收站文件保留天数，0表示永久保留',
    show_hidden_files BOOLEAN DEFAULT FALSE COMMENT '是否显示隐藏文件（以.开头的文件）',
    default_sort_order ENUM('name_asc', 'name_desc', 'size_asc', 'size_desc', 'date_asc', 'date_desc', 'type_asc', 'type_desc') DEFAULT 'name_asc' COMMENT '默认文件排序方式',
    list_view_type ENUM('grid', 'list', 'detail') DEFAULT 'grid' COMMENT '文件列表显示类型',
    items_per_page INT UNSIGNED DEFAULT 50 COMMENT '每页显示文件数量，范围10-200',
    
    -- 分享设置
    share_default_expire INT UNSIGNED DEFAULT 7 COMMENT '默认分享过期天数，0表示永久',
    share_require_password BOOLEAN DEFAULT FALSE COMMENT '分享是否默认需要密码',
    share_default_permission ENUM('view', 'download', 'edit') DEFAULT 'view' COMMENT '分享默认权限',
    share_notify_updates BOOLEAN DEFAULT TRUE COMMENT '分享文件更新时是否通知',
    allow_public_share BOOLEAN DEFAULT TRUE COMMENT '是否允许创建公开分享链接',
    
    -- 通知设置
    notification_file_shared BOOLEAN DEFAULT TRUE COMMENT '文件被分享时是否通知',
    notification_comment_added BOOLEAN DEFAULT TRUE COMMENT '收到评论时是否通知',
    notification_mention BOOLEAN DEFAULT TRUE COMMENT '被@提及时是否通知',
    notification_storage_warning BOOLEAN DEFAULT TRUE COMMENT '存储空间警告通知',
    notification_security_alert BOOLEAN DEFAULT TRUE COMMENT '安全警报通知（异常登录等）',
    notification_system_update BOOLEAN DEFAULT TRUE COMMENT '系统更新通知',
    
    -- 免打扰时间设置
    quiet_hours_enabled BOOLEAN DEFAULT FALSE COMMENT '是否启用免打扰时间',
    quiet_hours_start TIME DEFAULT NULL COMMENT '免打扰开始时间（24小时格式）',
    quiet_hours_end TIME DEFAULT NULL COMMENT '免打扰结束时间（24小时格式）',
    quiet_hours_timezone VARCHAR(50) DEFAULT 'Asia/Shanghai' COMMENT '免打扰时间的时区',
    
    -- 安全设置
    login_notification BOOLEAN DEFAULT TRUE COMMENT '新设备登录时是否通知',
    download_notification BOOLEAN DEFAULT FALSE COMMENT '文件下载时是否通知',
    require_password_for_sensitive BOOLEAN DEFAULT TRUE COMMENT '敏感操作是否需要密码确认',
    auto_logout_minutes INT UNSIGNED DEFAULT 480 COMMENT '自动注销时间（分钟），0表示不自动注销，默认8小时',
    
    -- 界面偏好设置
    sidebar_collapsed BOOLEAN DEFAULT FALSE COMMENT '侧边栏是否默认收起',
    show_file_extensions BOOLEAN DEFAULT TRUE COMMENT '是否显示文件扩展名',
    compact_mode BOOLEAN DEFAULT FALSE COMMENT '是否启用紧凑模式',
    enable_animations BOOLEAN DEFAULT TRUE COMMENT '是否启用界面动画',
    
    -- 高级设置
    enable_debug_mode BOOLEAN DEFAULT FALSE COMMENT '是否启用调试模式（显示更多技术信息）',
    cache_size_mb INT UNSIGNED DEFAULT 100 COMMENT '本地缓存大小限制（MB），仅客户端使用',
    concurrent_uploads INT UNSIGNED DEFAULT 3 COMMENT '并发上传数量限制，范围1-10',
    concurrent_downloads INT UNSIGNED DEFAULT 5 COMMENT '并发下载数量限制，范围1-20',
    
    -- 时间戳字段
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '设置创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '设置更新时间',
    
    -- 外键约束
    CONSTRAINT fk_user_settings_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
    
    -- 索引设计
    INDEX idx_user_settings_user_id (user_id) COMMENT '用户ID索引，用于快速查找用户设置',
    INDEX idx_user_settings_updated_at (updated_at) COMMENT '更新时间索引，用于同步检查'
    
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='用户个性化设置表 - 存储用户的各项偏好设置和配置选项'
  ROW_FORMAT=DYNAMIC;

-- ============================================================================
-- 用户设置约束和触发器
-- ============================================================================

-- 添加数值范围约束
ALTER TABLE user_settings ADD CONSTRAINT chk_sync_bandwidth_limit 
CHECK (sync_bandwidth_limit >= 1048576 AND sync_bandwidth_limit <= 104857600); -- 1MB/s 到 100MB/s

ALTER TABLE user_settings ADD CONSTRAINT chk_trash_retention_days 
CHECK (trash_retention_days >= 0 AND trash_retention_days <= 365);

ALTER TABLE user_settings ADD CONSTRAINT chk_items_per_page 
CHECK (items_per_page >= 10 AND items_per_page <= 200);

ALTER TABLE user_settings ADD CONSTRAINT chk_share_default_expire 
CHECK (share_default_expire >= 0 AND share_default_expire <= 365);

ALTER TABLE user_settings ADD CONSTRAINT chk_auto_logout_minutes 
CHECK (auto_logout_minutes = 0 OR (auto_logout_minutes >= 5 AND auto_logout_minutes <= 10080)); -- 5分钟到7天

ALTER TABLE user_settings ADD CONSTRAINT chk_cache_size_mb 
CHECK (cache_size_mb >= 10 AND cache_size_mb <= 10240); -- 10MB到10GB

ALTER TABLE user_settings ADD CONSTRAINT chk_concurrent_uploads 
CHECK (concurrent_uploads >= 1 AND concurrent_uploads <= 10);

ALTER TABLE user_settings ADD CONSTRAINT chk_concurrent_downloads 
CHECK (concurrent_downloads >= 1 AND concurrent_downloads <= 20);

-- 免打扰时间逻辑约束
ALTER TABLE user_settings ADD CONSTRAINT chk_quiet_hours_logic 
CHECK (
    (quiet_hours_enabled = FALSE) OR 
    (quiet_hours_enabled = TRUE AND quiet_hours_start IS NOT NULL AND quiet_hours_end IS NOT NULL)
);

-- 创建用户设置初始化触发器
DELIMITER //
CREATE TRIGGER user_settings_init_defaults 
BEFORE INSERT ON user_settings
FOR EACH ROW
BEGIN
    -- 确保免打扰时区与用户时区一致
    IF NEW.quiet_hours_timezone IS NULL THEN
        SELECT timezone INTO NEW.quiet_hours_timezone 
        FROM users 
        WHERE id = NEW.user_id;
    END IF;
    
    -- 根据用户角色设置不同的默认值
    SET @user_role = (SELECT role FROM users WHERE id = NEW.user_id);
    
    IF @user_role IN ('admin', 'super_admin') THEN
        -- 管理员用户的特殊设置
        SET NEW.enable_debug_mode = COALESCE(NEW.enable_debug_mode, TRUE);
        SET NEW.concurrent_uploads = COALESCE(NEW.concurrent_uploads, 5);
        SET NEW.concurrent_downloads = COALESCE(NEW.concurrent_downloads, 10);
    END IF;
END//
DELIMITER ;

-- ============================================================================
-- 初始化默认用户设置
-- ============================================================================

-- 为已存在的用户创建默认设置
INSERT INTO user_settings (user_id, created_at)
SELECT id, CURRENT_TIMESTAMP
FROM users 
WHERE deleted_at IS NULL
  AND id NOT IN (SELECT user_id FROM user_settings);

-- ============================================================================
-- 用户设置统计视图
-- ============================================================================

CREATE VIEW user_settings_stats AS
SELECT 
    COUNT(*) as total_users,
    COUNT(CASE WHEN sync_enabled = TRUE THEN 1 END) as sync_enabled_count,
    COUNT(CASE WHEN auto_preview = TRUE THEN 1 END) as auto_preview_count,
    COUNT(CASE WHEN quiet_hours_enabled = TRUE THEN 1 END) as quiet_hours_enabled_count,
    COUNT(CASE WHEN enable_debug_mode = TRUE THEN 1 END) as debug_mode_count,
    AVG(sync_bandwidth_limit) as avg_sync_bandwidth,
    AVG(trash_retention_days) as avg_trash_retention,
    AVG(items_per_page) as avg_items_per_page
FROM user_settings;