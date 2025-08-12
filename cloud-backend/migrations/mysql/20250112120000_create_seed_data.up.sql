-- ============================================================================
-- 种子数据脚本 - 初始化系统基础数据
-- 文件: 20250112120000_create_seed_data.up.sql
-- 描述: 创建管理员用户、默认角色、系统配置等种子数据
-- ============================================================================

-- 设置字符集和排序规则
SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci;

-- ============================================================================
-- 1. 管理员用户和角色数据
-- ============================================================================

-- 插入超级管理员账户
INSERT IGNORE INTO users (
    id,
    username, 
    email, 
    email_verified, 
    password_hash, 
    real_name, 
    role, 
    status,
    storage_quota,
    upload_bandwidth_limit,
    download_bandwidth_limit,
    max_file_size,
    timezone,
    language,
    created_at,
    updated_at
) VALUES (
    1,
    'superadmin',
    'superadmin@ycgcloud.com',
    TRUE,
    '$2a$12$rQ8Kd5K5K5K5K5K5K5K5Ku',  -- bcrypt(admin123) - 生产环境需修改
    '超级管理员',
    'super_admin',
    'active',
    1099511627776,  -- 1TB配额
    1073741824,     -- 1GB/s上传
    1073741824,     -- 1GB/s下载
    107374182400,   -- 100GB单文件限制
    'Asia/Shanghai',
    'zh-CN',
    '2025-01-01 00:00:00',
    CURRENT_TIMESTAMP
);

-- 插入系统管理员账户
INSERT IGNORE INTO users (
    id,
    username, 
    email, 
    email_verified, 
    password_hash, 
    real_name, 
    role, 
    status,
    storage_quota,
    upload_bandwidth_limit,
    download_bandwidth_limit,
    max_file_size,
    timezone,
    language,
    created_at,
    updated_at
) VALUES (
    2,
    'admin',
    'admin@ycgcloud.com',
    TRUE,
    '$2a$12$aQ8Kd5K5K5K5K5K5K5K5Ku',  -- bcrypt(admin123) - 生产环境需修改
    '系统管理员',
    'admin',
    'active',
    536870912000,   -- 500GB配额
    536870912,      -- 512MB/s上传
    536870912,      -- 512MB/s下载
    53687091200,    -- 50GB单文件限制
    'Asia/Shanghai',
    'zh-CN',
    '2025-01-01 00:00:00',
    CURRENT_TIMESTAMP
);

-- 插入运维管理员账户
INSERT IGNORE INTO users (
    id,
    username, 
    email, 
    email_verified, 
    password_hash, 
    real_name, 
    role, 
    status,
    storage_quota,
    upload_bandwidth_limit,
    download_bandwidth_limit,
    max_file_size,
    timezone,
    language,
    created_at,
    updated_at
) VALUES (
    3,
    'operator',
    'operator@ycgcloud.com',
    TRUE,
    '$2a$12$bQ8Kd5K5K5K5K5K5K5K5Ku',  -- bcrypt(operator123) - 生产环境需修改
    '运维管理员',
    'admin',
    'active',
    107374182400,   -- 100GB配额
    104857600,      -- 100MB/s上传
    104857600,      -- 100MB/s下载
    10737418240,    -- 10GB单文件限制
    'Asia/Shanghai',
    'zh-CN',
    '2025-01-01 00:00:00',
    CURRENT_TIMESTAMP
);

-- 插入系统用户（用于自动化操作）
INSERT IGNORE INTO users (
    id,
    username, 
    email, 
    email_verified, 
    password_hash, 
    real_name, 
    role, 
    status,
    storage_quota,
    upload_bandwidth_limit,
    download_bandwidth_limit,
    max_file_size,
    timezone,
    language,
    created_at,
    updated_at
) VALUES (
    4,
    'system',
    'system@ycgcloud.com',
    TRUE,
    '$2a$12$cQ8Kd5K5K5K5K5K5K5K5Ku',  -- bcrypt(system123) - 生产环境需修改
    '系统用户',
    'system',
    'active',
    0,              -- 无存储配额
    1073741824,     -- 1GB/s上传
    1073741824,     -- 1GB/s下载
    107374182400,   -- 100GB单文件限制
    'Asia/Shanghai',
    'zh-CN',
    '2025-01-01 00:00:00',
    CURRENT_TIMESTAMP
);

-- 插入示例普通用户
INSERT IGNORE INTO users (
    id,
    username, 
    email, 
    email_verified, 
    password_hash, 
    real_name, 
    role, 
    status,
    storage_quota,
    upload_bandwidth_limit,
    download_bandwidth_limit,
    max_file_size,
    timezone,
    language,
    created_at,
    updated_at
) VALUES (
    5,
    'demo_user',
    'demo@ycgcloud.com',
    TRUE,
    '$2a$12$dQ8Kd5K5K5K5K5K5K5K5Ku',  -- bcrypt(demo123) - 生产环境需修改
    '演示用户',
    'user',
    'active',
    10737418240,    -- 10GB配额
    10485760,       -- 10MB/s上传
    10485760,       -- 10MB/s下载
    1073741824,     -- 1GB单文件限制
    'Asia/Shanghai',
    'zh-CN',
    '2025-01-01 00:00:00',
    CURRENT_TIMESTAMP
);

-- ============================================================================
-- 2. 用户设置默认数据
-- ============================================================================

-- 为管理员用户创建默认设置
INSERT IGNORE INTO user_settings (
    user_id,
    theme,
    language,
    timezone,
    items_per_page,
    sync_enabled,
    sync_bandwidth_limit,
    auto_preview,
    trash_retention_days,
    share_default_expire,
    enable_notifications,
    notification_email,
    notification_browser,
    notification_mobile,
    auto_logout_minutes,
    enable_two_factor,
    cache_size_mb,
    concurrent_uploads,
    concurrent_downloads,
    created_at,
    updated_at
) VALUES 
-- 超级管理员设置
(1, 'dark', 'zh-CN', 'Asia/Shanghai', 50, TRUE, 104857600, TRUE, 30, 7, TRUE, TRUE, TRUE, TRUE, 60, FALSE, 1024, 5, 10, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
-- 系统管理员设置
(2, 'light', 'zh-CN', 'Asia/Shanghai', 50, TRUE, 52428800, TRUE, 30, 7, TRUE, TRUE, TRUE, TRUE, 30, FALSE, 512, 3, 8, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
-- 运维管理员设置
(3, 'light', 'zh-CN', 'Asia/Shanghai', 30, TRUE, 10485760, TRUE, 7, 3, TRUE, TRUE, FALSE, FALSE, 120, FALSE, 256, 2, 5, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
-- 系统用户设置
(4, 'light', 'zh-CN', 'Asia/Shanghai', 100, FALSE, 0, FALSE, 0, 0, FALSE, FALSE, FALSE, FALSE, 0, FALSE, 128, 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
-- 演示用户设置
(5, 'light', 'zh-CN', 'Asia/Shanghai', 20, TRUE, 10485760, TRUE, 7, 1, TRUE, TRUE, TRUE, FALSE, 30, FALSE, 128, 2, 3, '2025-01-01 00:00:00', CURRENT_TIMESTAMP);

-- ============================================================================
-- 3. 系统配置数据
-- ============================================================================

-- 插入系统基础配置
INSERT IGNORE INTO system_configs (
    id,
    config_key,
    config_value,
    value_type,
    category,
    subcategory,
    config_group,
    display_name,
    description,
    is_public,
    is_editable,
    is_system,
    is_sensitive,
    access_level,
    is_active,
    environment,
    feature_flag,
    config_source,
    created_by,
    updated_by,
    created_at,
    updated_at
) VALUES 
-- 系统基础配置
(1, 'system.app_name', 'YCG云盘系统', 'string', 'system', 'basic', 'application', '应用名称', '系统应用名称', TRUE, TRUE, FALSE, FALSE, 'public', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(2, 'system.app_version', '1.0.0', 'string', 'system', 'basic', 'application', '应用版本', '当前系统版本号', TRUE, FALSE, TRUE, FALSE, 'public', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(3, 'system.maintenance_mode', 'false', 'boolean', 'system', 'basic', 'application', '维护模式', '系统维护模式开关', FALSE, TRUE, TRUE, FALSE, 'admin', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(4, 'system.max_file_size', '107374182400', 'integer', 'storage', 'limits', 'file_management', '最大文件大小', '单个文件最大大小(字节)', FALSE, TRUE, FALSE, FALSE, 'admin', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(5, 'system.default_storage_quota', '10737418240', 'integer', 'storage', 'limits', 'user_management', '默认存储配额', '新用户默认存储配额(字节)', FALSE, TRUE, FALSE, FALSE, 'admin', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),

-- 安全配置
(6, 'security.password_min_length', '8', 'integer', 'security', 'password', 'authentication', '密码最小长度', '用户密码最小长度要求', FALSE, TRUE, FALSE, FALSE, 'admin', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(7, 'security.session_timeout', '7200', 'integer', 'security', 'session', 'authentication', '会话超时时间', '用户会话超时时间(秒)', FALSE, TRUE, FALSE, FALSE, 'admin', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(8, 'security.max_login_attempts', '5', 'integer', 'security', 'login', 'authentication', '最大登录尝试次数', '用户登录失败最大次数', FALSE, TRUE, FALSE, FALSE, 'admin', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(9, 'security.lockout_duration', '900', 'integer', 'security', 'login', 'authentication', '账户锁定时长', '登录失败后账户锁定时长(秒)', FALSE, TRUE, FALSE, FALSE, 'admin', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),

-- 文件管理配置
(10, 'file.allowed_extensions', '["jpg","jpeg","png","gif","bmp","webp","pdf","doc","docx","xls","xlsx","ppt","pptx","txt","md","zip","rar","7z","mp4","avi","mkv","mp3","wav","flac"]', 'json', 'file', 'upload', 'file_management', '允许的文件扩展名', '系统允许上传的文件扩展名列表', FALSE, TRUE, FALSE, FALSE, 'admin', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(11, 'file.chunk_size', '5242880', 'integer', 'file', 'upload', 'file_management', '分片上传大小', '文件分片上传时每片大小(字节)', FALSE, TRUE, FALSE, FALSE, 'admin', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(12, 'file.preview_enabled', 'true', 'boolean', 'file', 'preview', 'file_management', '文件预览功能', '是否启用文件在线预览功能', TRUE, TRUE, FALSE, FALSE, 'public', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),

-- 分享配置
(13, 'share.default_expire_days', '7', 'integer', 'share', 'policy', 'sharing', '默认分享过期天数', '文件分享默认过期天数', FALSE, TRUE, FALSE, FALSE, 'admin', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(14, 'share.max_expire_days', '365', 'integer', 'share', 'policy', 'sharing', '最大分享过期天数', '文件分享最大过期天数', FALSE, TRUE, FALSE, FALSE, 'admin', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(15, 'share.require_password', 'false', 'boolean', 'share', 'policy', 'sharing', '强制密码保护', '是否强制要求分享设置密码', FALSE, TRUE, FALSE, FALSE, 'admin', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),

-- 通知配置
(16, 'notification.email_enabled', 'true', 'boolean', 'notification', 'email', 'messaging', '邮件通知开关', '是否启用邮件通知功能', FALSE, TRUE, FALSE, FALSE, 'admin', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(17, 'notification.browser_enabled', 'true', 'boolean', 'notification', 'browser', 'messaging', '浏览器通知开关', '是否启用浏览器推送通知', TRUE, TRUE, FALSE, FALSE, 'public', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),

-- OSS配置 (敏感信息)
(18, 'oss.access_key_id', 'your_access_key_id', 'string', 'storage', 'oss', 'external_service', 'OSS访问密钥ID', '阿里云OSS访问密钥ID', FALSE, TRUE, FALSE, TRUE, 'super_admin', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(19, 'oss.access_key_secret', 'your_access_key_secret', 'string', 'storage', 'oss', 'external_service', 'OSS访问密钥Secret', '阿里云OSS访问密钥Secret', FALSE, TRUE, FALSE, TRUE, 'super_admin', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(20, 'oss.bucket_name', 'ycg-cloud-storage', 'string', 'storage', 'oss', 'external_service', 'OSS存储桶名称', '阿里云OSS存储桶名称', FALSE, TRUE, FALSE, FALSE, 'admin', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(21, 'oss.endpoint', 'oss-cn-hangzhou.aliyuncs.com', 'string', 'storage', 'oss', 'external_service', 'OSS端点地址', '阿里云OSS服务端点', FALSE, TRUE, FALSE, FALSE, 'admin', TRUE, 'production', NULL, 'manual', 1, 1, '2025-01-01 00:00:00', CURRENT_TIMESTAMP);

-- ============================================================================
-- 4. 示例团队数据
-- ============================================================================

-- 插入示例团队
INSERT IGNORE INTO teams (
    id,
    name,
    slug,
    description,
    owner_id,
    team_type,
    status,
    is_active,
    plan_type,
    verification_level,
    member_count,
    max_members,
    storage_quota,
    storage_used,
    storage_warning_threshold,
    max_projects,
    max_shared_links,
    max_file_size,
    total_files_created,
    total_shares_created,
    country_code,
    created_at,
    updated_at
) VALUES (
    1,
    '系统管理团队',
    'system-admin-team',
    '系统管理员专用团队，负责系统维护和管理工作',
    1,  -- superadmin
    'organization',
    'active',
    TRUE,
    'enterprise',
    'verified',
    3,
    100,
    5497558138880,  -- 5TB
    0,
    0.8,
    1000,
    10000,
    107374182400,
    0,
    0,
    'CN',
    '2025-01-01 00:00:00',
    CURRENT_TIMESTAMP
);

-- 插入团队成员关系
INSERT IGNORE INTO team_members (
    team_id,
    user_id,
    role,
    status,
    invited_by,
    invited_at,
    joined_at,
    access_level,
    display_name,
    bio,
    job_title,
    department,
    storage_quota,
    storage_used,
    activity_score,
    created_at,
    updated_at
) VALUES 
-- 超级管理员
(1, 1, 'owner', 'active', NULL, '2025-01-01 00:00:00', '2025-01-01 00:00:00', 'full', '超级管理员', '系统超级管理员', 'CTO', '技术部', 1099511627776, 0, 100.0, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
-- 系统管理员
(1, 2, 'admin', 'active', 1, '2025-01-01 00:00:00', '2025-01-01 00:00:00', 'admin', '系统管理员', '系统日常管理', '系统管理员', '技术部', 536870912000, 0, 90.0, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
-- 运维管理员
(1, 3, 'moderator', 'active', 1, '2025-01-01 00:00:00', '2025-01-01 00:00:00', 'limited', '运维管理员', '系统运维监控', '运维工程师', '技术部', 107374182400, 0, 85.0, '2025-01-01 00:00:00', CURRENT_TIMESTAMP);

-- ============================================================================
-- 5. 默认文件夹数据
-- ============================================================================

-- 为每个用户创建根目录和默认文件夹
INSERT IGNORE INTO files (
    id,
    filename,
    file_path,
    file_type,
    is_folder,
    user_id,
    parent_id,
    file_size,
    storage_location,
    created_at,
    updated_at
) VALUES 
-- 超级管理员根目录
(1, '根目录', '/', 'folder', TRUE, 1, NULL, 0, 'local', '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(2, '文档', '/文档', 'folder', TRUE, 1, 1, 0, 'local', '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(3, '图片', '/图片', 'folder', TRUE, 1, 1, 0, 'local', '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(4, '视频', '/视频', 'folder', TRUE, 1, 1, 0, 'local', '2025-01-01 00:00:00', CURRENT_TIMESTAMP),

-- 系统管理员根目录
(5, '根目录', '/', 'folder', TRUE, 2, NULL, 0, 'local', '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(6, '系统文档', '/系统文档', 'folder', TRUE, 2, 5, 0, 'local', '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(7, '备份文件', '/备份文件', 'folder', TRUE, 2, 5, 0, 'local', '2025-01-01 00:00:00', CURRENT_TIMESTAMP),

-- 演示用户根目录
(8, '根目录', '/', 'folder', TRUE, 5, NULL, 0, 'local', '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(9, '我的文档', '/我的文档', 'folder', TRUE, 5, 8, 0, 'local', '2025-01-01 00:00:00', CURRENT_TIMESTAMP),
(10, '照片', '/照片', 'folder', TRUE, 5, 8, 0, 'local', '2025-01-01 00:00:00', CURRENT_TIMESTAMP);

-- 插入文件夹扩展属性
INSERT IGNORE INTO folders (
    id,
    folder_type,
    folder_category,
    total_size,
    file_count,
    folder_count,
    total_file_count,
    total_folder_count,
    default_permission,
    max_file_size,
    max_folder_size,
    max_file_versions,
    is_template,
    collaboration_enabled,
    discussion_enabled,
    task_management,
    search_weight,
    cache_ttl,
    stats_updated_at
) VALUES 
-- 根目录配置
(1, 'root', 'personal', 0, 0, 3, 0, 3, 'private', 107374182400, NULL, 50, FALSE, TRUE, TRUE, TRUE, 1.0, 3600, CURRENT_TIMESTAMP),
(5, 'root', 'personal', 0, 0, 2, 0, 2, 'private', 53687091200, NULL, 30, FALSE, TRUE, TRUE, TRUE, 1.0, 3600, CURRENT_TIMESTAMP),
(8, 'root', 'personal', 0, 0, 2, 0, 2, 'private', 1073741824, NULL, 10, FALSE, TRUE, FALSE, FALSE, 1.0, 3600, CURRENT_TIMESTAMP),

-- 普通文件夹配置
(2, 'standard', 'documents', 0, 0, 0, 0, 0, 'private', 107374182400, NULL, 50, FALSE, TRUE, TRUE, FALSE, 2.0, 1800, CURRENT_TIMESTAMP),
(3, 'standard', 'media', 0, 0, 0, 0, 0, 'private', 107374182400, NULL, 20, FALSE, FALSE, FALSE, FALSE, 1.5, 1800, CURRENT_TIMESTAMP),
(4, 'standard', 'media', 0, 0, 0, 0, 0, 'private', 107374182400, NULL, 10, FALSE, FALSE, FALSE, FALSE, 1.2, 1800, CURRENT_TIMESTAMP),
(6, 'standard', 'documents', 0, 0, 0, 0, 0, 'private', 53687091200, NULL, 100, FALSE, TRUE, TRUE, TRUE, 3.0, 1800, CURRENT_TIMESTAMP),
(7, 'standard', 'backup', 0, 0, 0, 0, 0, 'private', 53687091200, NULL, 5, FALSE, FALSE, FALSE, FALSE, 1.0, 7200, CURRENT_TIMESTAMP),
(9, 'standard', 'documents', 0, 0, 0, 0, 0, 'private', 1073741824, NULL, 10, FALSE, TRUE, FALSE, FALSE, 2.0, 1800, CURRENT_TIMESTAMP),
(10, 'standard', 'media', 0, 0, 0, 0, 0, 'private', 1073741824, NULL, 5, FALSE, FALSE, FALSE, FALSE, 1.5, 1800, CURRENT_TIMESTAMP);

-- ============================================================================
-- 6. 示例通知数据
-- ============================================================================

-- 插入欢迎通知
INSERT IGNORE INTO notifications (
    id,
    user_id,
    sender_id,
    notification_type,
    notification_category,
    title,
    content,
    status,
    priority,
    is_read,
    related_type,
    related_id,
    scheduled_at,
    expires_at,
    importance_score,
    created_at,
    updated_at
) VALUES 
-- 超级管理员欢迎通知
(1, 1, 4, 'system', 'welcome', '欢迎使用YCG云盘系统', '您好，超级管理员！欢迎使用YCG云盘系统。您拥有系统最高权限，请合理使用管理功能。', 'delivered', 'high', FALSE, 'user', 1, '2025-01-01 00:00:00', '2025-12-31 23:59:59', 90, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),

-- 系统管理员欢迎通知
(2, 2, 4, 'system', 'welcome', '欢迎加入管理团队', '您好，系统管理员！欢迎加入YCG云盘管理团队。请查看管理手册了解您的职责范围。', 'delivered', 'medium', FALSE, 'user', 2, '2025-01-01 00:00:00', '2025-12-31 23:59:59', 80, '2025-01-01 00:00:00', CURRENT_TIMESTAMP),

-- 演示用户欢迎通知
(3, 5, 4, 'system', 'welcome', '欢迎使用云盘服务', '您好！欢迎使用YCG云盘服务。您可以开始上传、管理和分享您的文件。如有问题请联系客服。', 'delivered', 'medium', FALSE, 'user', 5, '2025-01-01 00:00:00', '2025-12-31 23:59:59', 70, '2025-01-01 00:00:00', CURRENT_TIMESTAMP);

-- ============================================================================
-- 总结信息
-- ============================================================================

-- 插入种子数据创建日志
INSERT IGNORE INTO operation_logs (
    user_id,
    session_id,
    operation_type,
    operation_category,
    resource_type,
    resource_id,
    action,
    result,
    ip_address,
    user_agent,
    description,
    correlation_id,
    created_at
) VALUES (
    1,
    'seed-data-init',
    'data_management',
    'initialization',
    'database',
    'seed_data',
    'create',
    'success',
    '127.0.0.1',
    'MySQL/Migration',
    '初始化系统种子数据：创建管理员用户、默认配置、示例团队和文件夹',
    'seed-init-' || UNIX_TIMESTAMP(),
    CURRENT_TIMESTAMP
);

-- 显示种子数据创建完成信息
SELECT 
    '种子数据创建完成' as message,
    COUNT(DISTINCT u.id) as users_created,
    COUNT(DISTINCT t.id) as teams_created,
    COUNT(DISTINCT sc.id) as configs_created,
    COUNT(DISTINCT f.id) as folders_created,
    COUNT(DISTINCT n.id) as notifications_created
FROM users u
CROSS JOIN teams t  
CROSS JOIN system_configs sc
CROSS JOIN files f
CROSS JOIN notifications n
WHERE u.id IN (1,2,3,4,5)
  AND t.id = 1
  AND sc.id BETWEEN 1 AND 21
  AND f.id BETWEEN 1 AND 10
  AND n.id BETWEEN 1 AND 3;