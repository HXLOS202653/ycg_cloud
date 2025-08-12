-- +migrate Up
-- 创建迁移: 剩余核心表批量创建
-- 版本: 20250112101800
-- 描述: 批量创建聊天系统和搜索系统的剩余表
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:18:00
-- 依赖: 之前的所有迁移文件
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1

-- ============================================================================
-- 聊天室表 (chat_rooms)
-- ============================================================================

CREATE TABLE chat_rooms (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '聊天室唯一标识',
    name VARCHAR(100) DEFAULT NULL COMMENT '聊天室名称',
    description TEXT DEFAULT NULL COMMENT '聊天室描述',
    room_type ENUM('private', 'group', 'team', 'file_discussion', 'public', 'broadcast') NOT NULL COMMENT '聊天室类型',
    avatar_url VARCHAR(500) DEFAULT NULL COMMENT '聊天室头像URL',
    banner_url VARCHAR(500) DEFAULT NULL COMMENT '聊天室横幅图片',
    
    owner_id BIGINT UNSIGNED DEFAULT NULL COMMENT '聊天室所有者ID',
    team_id BIGINT UNSIGNED DEFAULT NULL COMMENT '关联团队ID',
    file_id BIGINT UNSIGNED DEFAULT NULL COMMENT '关联文件ID（文件讨论）',
    
    max_members INT UNSIGNED DEFAULT 500 COMMENT '最大成员数限制',
    member_count INT UNSIGNED DEFAULT 0 COMMENT '当前成员数',
    active_member_count INT UNSIGNED DEFAULT 0 COMMENT '活跃成员数',
    
    last_message_id BIGINT UNSIGNED DEFAULT NULL COMMENT '最后一条消息ID',
    last_activity_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '最后活动时间',
    
    is_active BOOLEAN DEFAULT TRUE COMMENT '聊天室是否有效',
    is_private BOOLEAN DEFAULT FALSE COMMENT '是否为私有聊天室',
    is_muted BOOLEAN DEFAULT FALSE COMMENT '是否静音',
    
    settings JSON DEFAULT NULL COMMENT '聊天室设置',
    permissions JSON DEFAULT NULL COMMENT '权限配置',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',
    
    INDEX idx_chat_rooms_room_type (room_type),
    INDEX idx_chat_rooms_owner_id (owner_id),
    INDEX idx_chat_rooms_team_id (team_id),
    INDEX idx_chat_rooms_file_id (file_id),
    INDEX idx_chat_rooms_last_activity_at (last_activity_at),
    INDEX idx_chat_rooms_is_active (is_active),
    INDEX idx_chat_rooms_created_at (created_at),
    INDEX idx_chat_rooms_deleted_at (deleted_at),
    
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天室表';

-- ============================================================================
-- 聊天室成员表 (chat_room_members)
-- ============================================================================

CREATE TABLE chat_room_members (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '聊天室成员关系ID',
    room_id BIGINT UNSIGNED NOT NULL COMMENT '聊天室ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    
    role ENUM('owner', 'admin', 'member', 'guest') DEFAULT 'member' COMMENT '成员角色',
    status ENUM('active', 'muted', 'banned', 'left') DEFAULT 'active' COMMENT '成员状态',
    
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
    last_read_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '最后阅读时间',
    last_active_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '最后活跃时间',
    
    unread_count INT UNSIGNED DEFAULT 0 COMMENT '未读消息数',
    notification_enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用通知',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    UNIQUE KEY uk_room_user (room_id, user_id),
    INDEX idx_chat_room_members_room_id (room_id),
    INDEX idx_chat_room_members_user_id (user_id),
    INDEX idx_chat_room_members_role (role),
    INDEX idx_chat_room_members_status (status),
    INDEX idx_chat_room_members_last_active_at (last_active_at),
    
    FOREIGN KEY (room_id) REFERENCES chat_rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天室成员表';

-- ============================================================================
-- 聊天消息表 (chat_messages)
-- ============================================================================

CREATE TABLE chat_messages (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '消息唯一标识',
    room_id BIGINT UNSIGNED NOT NULL COMMENT '聊天室ID',
    sender_id BIGINT UNSIGNED DEFAULT NULL COMMENT '发送者ID',
    
    message_type ENUM('text', 'image', 'file', 'video', 'audio', 'location', 'system', 'quote') DEFAULT 'text' COMMENT '消息类型',
    content TEXT DEFAULT NULL COMMENT '消息内容',
    content_json JSON DEFAULT NULL COMMENT '结构化消息内容',
    
    reply_to_id BIGINT UNSIGNED DEFAULT NULL COMMENT '回复的消息ID',
    forwarded_from_id BIGINT UNSIGNED DEFAULT NULL COMMENT '转发来源消息ID',
    
    file_url VARCHAR(1000) DEFAULT NULL COMMENT '文件URL',
    file_name VARCHAR(255) DEFAULT NULL COMMENT '文件名',
    file_size BIGINT UNSIGNED DEFAULT NULL COMMENT '文件大小',
    file_type VARCHAR(50) DEFAULT NULL COMMENT '文件类型',
    
    is_edited BOOLEAN DEFAULT FALSE COMMENT '是否已编辑',
    is_deleted BOOLEAN DEFAULT FALSE COMMENT '是否已删除',
    is_pinned BOOLEAN DEFAULT FALSE COMMENT '是否置顶',
    
    read_count INT UNSIGNED DEFAULT 0 COMMENT '已读人数',
    reaction_count INT UNSIGNED DEFAULT 0 COMMENT '表情反应数',
    
    edited_at TIMESTAMP NULL DEFAULT NULL COMMENT '编辑时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '删除时间',
    
    INDEX idx_chat_messages_room_id (room_id),
    INDEX idx_chat_messages_sender_id (sender_id),
    INDEX idx_chat_messages_message_type (message_type),
    INDEX idx_chat_messages_reply_to_id (reply_to_id),
    INDEX idx_chat_messages_created_at (created_at),
    INDEX idx_chat_messages_is_deleted (is_deleted),
    
    FOREIGN KEY (room_id) REFERENCES chat_rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (reply_to_id) REFERENCES chat_messages(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='聊天消息表'
PARTITION BY RANGE (YEAR(created_at)) (
    PARTITION p2024 VALUES LESS THAN (2025),
    PARTITION p2025 VALUES LESS THAN (2026),
    PARTITION p2026 VALUES LESS THAN (2027),
    PARTITION p_future VALUES LESS THAN MAXVALUE
);

-- ============================================================================
-- 文件内容索引表 (file_contents)
-- ============================================================================

CREATE TABLE file_contents (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '内容索引ID',
    file_id BIGINT UNSIGNED NOT NULL COMMENT '文件ID',
    
    content_type ENUM('text', 'metadata', 'ocr', 'audio_transcript', 'video_subtitle') DEFAULT 'text' COMMENT '内容类型',
    extracted_text LONGTEXT DEFAULT NULL COMMENT '提取的文本内容',
    language VARCHAR(10) DEFAULT NULL COMMENT '语言代码',
    
    word_count INT UNSIGNED DEFAULT 0 COMMENT '字数统计',
    character_count INT UNSIGNED DEFAULT 0 COMMENT '字符数统计',
    
    extraction_method VARCHAR(50) DEFAULT NULL COMMENT '提取方法',
    extraction_confidence DECIMAL(3,2) DEFAULT NULL COMMENT '提取置信度',
    
    metadata JSON DEFAULT NULL COMMENT '内容元数据',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    
    INDEX idx_file_contents_file_id (file_id),
    INDEX idx_file_contents_content_type (content_type),
    INDEX idx_file_contents_language (language),
    FULLTEXT idx_file_contents_fulltext_content (extracted_text),
    
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件内容索引表';

-- ============================================================================
-- 搜索日志表 (search_logs)
-- ============================================================================

CREATE TABLE search_logs (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '搜索日志ID',
    user_id BIGINT UNSIGNED DEFAULT NULL COMMENT '搜索用户ID',
    session_id VARCHAR(128) DEFAULT NULL COMMENT '会话ID',
    
    search_query VARCHAR(1000) NOT NULL COMMENT '搜索查询内容',
    search_type ENUM('file', 'content', 'user', 'team', 'global') DEFAULT 'file' COMMENT '搜索类型',
    search_scope JSON DEFAULT NULL COMMENT '搜索范围配置',
    
    result_count INT UNSIGNED DEFAULT 0 COMMENT '搜索结果数量',
    response_time_ms INT UNSIGNED DEFAULT NULL COMMENT '搜索响应时间（毫秒）',
    
    filters_applied JSON DEFAULT NULL COMMENT '应用的过滤器',
    sort_criteria VARCHAR(100) DEFAULT NULL COMMENT '排序条件',
    
    clicked_result_id BIGINT UNSIGNED DEFAULT NULL COMMENT '点击的结果ID',
    clicked_position INT UNSIGNED DEFAULT NULL COMMENT '点击结果的位置',
    
    ip_address VARCHAR(45) DEFAULT NULL COMMENT '搜索来源IP',
    user_agent TEXT DEFAULT NULL COMMENT '用户代理',
    
    is_successful BOOLEAN DEFAULT TRUE COMMENT '搜索是否成功',
    error_message TEXT DEFAULT NULL COMMENT '错误信息',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '搜索时间',
    
    INDEX idx_search_logs_user_id (user_id),
    INDEX idx_search_logs_search_type (search_type),
    INDEX idx_search_logs_created_at (created_at),
    INDEX idx_search_logs_is_successful (is_successful),
    FULLTEXT idx_search_logs_fulltext_query (search_query),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='搜索日志表'
PARTITION BY RANGE (YEAR(created_at)) (
    PARTITION p2024 VALUES LESS THAN (2025),
    PARTITION p2025 VALUES LESS THAN (2026),
    PARTITION p2026 VALUES LESS THAN (2027),
    PARTITION p_future VALUES LESS THAN MAXVALUE
);

-- ============================================================================
-- 创建基础视图
-- ============================================================================

-- 聊天室概览视图
CREATE VIEW chat_rooms_overview AS
SELECT 
    cr.id,
    cr.name,
    cr.room_type,
    cr.member_count,
    cr.last_activity_at,
    u.username as owner_name,
    t.name as team_name,
    f.filename as file_name
FROM chat_rooms cr
LEFT JOIN users u ON cr.owner_id = u.id
LEFT JOIN teams t ON cr.team_id = t.id
LEFT JOIN files f ON cr.file_id = f.id
WHERE cr.is_active = TRUE AND cr.deleted_at IS NULL;

-- 搜索统计视图
CREATE VIEW search_statistics AS
SELECT 
    DATE(created_at) as search_date,
    search_type,
    COUNT(*) as total_searches,
    COUNT(CASE WHEN is_successful = TRUE THEN 1 END) as successful_searches,
    AVG(result_count) as avg_results,
    AVG(response_time_ms) as avg_response_time,
    COUNT(DISTINCT user_id) as unique_users
FROM search_logs
WHERE created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
GROUP BY DATE(created_at), search_type
ORDER BY search_date DESC;