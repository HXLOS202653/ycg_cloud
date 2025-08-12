-- +migrate Up
-- 创建迁移: 通知表
-- 版本: 20250112101500
-- 描述: 创建通知表，管理系统通知和用户消息推送
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:15:00
-- 依赖: 20250112100000_create_users_table
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1
-- 说明: 通知表实现统一的消息推送机制，支持多种通知类型和优先级

-- ============================================================================
-- 通知表 (notifications)
-- ============================================================================

CREATE TABLE notifications (
    -- 基础标识字段
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '通知记录唯一标识',
    
    -- 通知接收者和发送者
    user_id BIGINT UNSIGNED NOT NULL COMMENT '通知接收者用户ID',
    sender_id BIGINT UNSIGNED DEFAULT NULL COMMENT '通知发送者用户ID，系统通知为NULL',
    
    -- 通知分类和类型
    notification_type ENUM(
        'system', 'file_share', 'team_invite', 'file_comment', 'file_update', 
        'storage_warning', 'security_alert', 'team_update', 'permission_change',
        'upload_complete', 'download_request', 'collaboration_invite',
        'payment_reminder', 'quota_exceeded', 'backup_complete',
        'virus_detected', 'login_alert', 'password_change'
    ) NOT NULL COMMENT '通知类型分类',
    
    notification_category ENUM('info', 'warning', 'error', 'success', 'promotion') DEFAULT 'info' COMMENT '通知类别',
    
    -- 通知内容
    title VARCHAR(200) NOT NULL COMMENT '通知标题，简短描述',
    content TEXT NOT NULL COMMENT '通知详细内容',
    summary VARCHAR(500) DEFAULT NULL COMMENT '通知摘要，用于预览',
    
    -- 操作和交互
    action_url VARCHAR(500) DEFAULT NULL COMMENT '操作链接URL',
    action_text VARCHAR(50) DEFAULT NULL COMMENT '操作按钮显示文本',
    action_type ENUM('view', 'accept', 'decline', 'download', 'upload', 'configure') DEFAULT NULL COMMENT '操作类型',
    
    -- 关联对象信息
    related_id BIGINT UNSIGNED DEFAULT NULL COMMENT '关联对象ID',
    related_type ENUM('file', 'folder', 'team', 'user', 'share', 'message', 'upload', 'system') DEFAULT NULL COMMENT '关联对象类型',
    related_name VARCHAR(255) DEFAULT NULL COMMENT '关联对象名称，冗余存储便于显示',
    
    -- 通知优先级和重要性
    priority ENUM('low', 'normal', 'high', 'urgent', 'critical') DEFAULT 'normal' COMMENT '通知优先级',
    importance_score INT UNSIGNED DEFAULT 50 COMMENT '重要性评分（0-100）',
    
    -- 通知状态和处理
    status ENUM('pending', 'sent', 'delivered', 'read', 'clicked', 'expired', 'failed') DEFAULT 'pending' COMMENT '通知状态',
    is_read BOOLEAN DEFAULT FALSE COMMENT '是否已读',
    is_archived BOOLEAN DEFAULT FALSE COMMENT '是否已归档',
    is_starred BOOLEAN DEFAULT FALSE COMMENT '是否已标星',
    
    -- 时间管理
    scheduled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '计划发送时间',
    sent_at TIMESTAMP NULL DEFAULT NULL COMMENT '实际发送时间',
    delivered_at TIMESTAMP NULL DEFAULT NULL COMMENT '送达时间',
    read_at TIMESTAMP NULL DEFAULT NULL COMMENT '阅读时间',
    clicked_at TIMESTAMP NULL DEFAULT NULL COMMENT '点击时间',
    expires_at TIMESTAMP NULL DEFAULT NULL COMMENT '过期时间',
    
    -- 推送渠道和方式
    channels JSON DEFAULT NULL COMMENT '推送渠道配置：email, push, sms, in_app',
    delivery_status JSON DEFAULT NULL COMMENT '各渠道的送达状态',
    
    -- 个性化和定制
    template_id VARCHAR(100) DEFAULT NULL COMMENT '通知模板ID',
    template_variables JSON DEFAULT NULL COMMENT '模板变量数据',
    custom_styling JSON DEFAULT NULL COMMENT '自定义样式配置',
    
    -- 批量通知相关
    batch_id VARCHAR(36) DEFAULT NULL COMMENT '批量通知批次ID',
    batch_sequence INT UNSIGNED DEFAULT NULL COMMENT '批次内序号',
    
    -- 用户交互统计
    view_count INT UNSIGNED DEFAULT 0 COMMENT '查看次数',
    click_count INT UNSIGNED DEFAULT 0 COMMENT '点击次数',
    share_count INT UNSIGNED DEFAULT 0 COMMENT '分享次数',
    
    -- 地理位置和设备信息
    user_location VARCHAR(100) DEFAULT NULL COMMENT '用户接收时的地理位置',
    user_device JSON DEFAULT NULL COMMENT '用户设备信息',
    user_timezone VARCHAR(50) DEFAULT NULL COMMENT '用户时区',
    
    -- A/B测试和分析
    experiment_id VARCHAR(100) DEFAULT NULL COMMENT '实验ID，用于A/B测试',
    variant_id VARCHAR(50) DEFAULT NULL COMMENT '变体ID',
    
    -- 扩展属性和元数据
    metadata JSON DEFAULT NULL COMMENT '通知扩展元数据',
    tags JSON DEFAULT NULL COMMENT '通知标签，便于分类',
    custom_fields JSON DEFAULT NULL COMMENT '自定义字段',
    
    -- 时间戳
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '通知创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '通知更新时间',
    deleted_at TIMESTAMP NULL DEFAULT NULL COMMENT '软删除时间',
    
    -- 业务索引设计
    INDEX idx_notifications_user_id (user_id) COMMENT '接收者索引，查询用户通知',
    INDEX idx_notifications_sender_id (sender_id) COMMENT '发送者索引，查询发送记录',
    INDEX idx_notifications_notification_type (notification_type) COMMENT '通知类型索引',
    INDEX idx_notifications_notification_category (notification_category) COMMENT '通知类别索引',
    INDEX idx_notifications_status (status) COMMENT '通知状态索引',
    INDEX idx_notifications_priority (priority) COMMENT '优先级索引',
    INDEX idx_notifications_is_read (is_read) COMMENT '已读状态索引',
    INDEX idx_notifications_is_archived (is_archived) COMMENT '归档状态索引',
    INDEX idx_notifications_is_starred (is_starred) COMMENT '标星状态索引',
    INDEX idx_notifications_related_id_type (related_id, related_type) COMMENT '关联对象复合索引',
    INDEX idx_notifications_scheduled_at (scheduled_at) COMMENT '计划发送时间索引',
    INDEX idx_notifications_expires_at (expires_at) COMMENT '过期时间索引，清理任务',
    INDEX idx_notifications_created_at (created_at) COMMENT '创建时间索引',
    INDEX idx_notifications_batch_id (batch_id) COMMENT '批量通知批次索引',
    INDEX idx_notifications_template_id (template_id) COMMENT '模板ID索引',
    INDEX idx_notifications_experiment_id (experiment_id) COMMENT '实验ID索引',
    INDEX idx_notifications_deleted_at (deleted_at) COMMENT '软删除索引',
    
    -- 复合业务索引
    INDEX idx_notifications_user_status_priority (user_id, status, priority DESC, created_at DESC) COMMENT '用户通知状态优先级复合索引',
    INDEX idx_notifications_user_read_archive (user_id, is_read, is_archived, updated_at DESC) COMMENT '用户阅读归档状态复合索引',
    INDEX idx_notifications_type_priority_time (notification_type, priority, scheduled_at) COMMENT '类型优先级时间复合索引',
    INDEX idx_notifications_sender_type_time (sender_id, notification_type, created_at DESC) COMMENT '发送者类型时间复合索引',
    INDEX idx_notifications_importance_unread (importance_score DESC, is_read, expires_at) COMMENT '重要性未读复合索引',
    
    -- 分区优化索引（基于时间）
    INDEX idx_notifications_daily_stats (DATE(created_at), notification_type, status) COMMENT '日统计索引',
    
    -- 全文搜索索引
    FULLTEXT idx_fulltext_search (title, content, summary) COMMENT '通知内容全文搜索',
    
    -- 外键约束
    CONSTRAINT fk_notifications_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_notifications_sender_id 
        FOREIGN KEY (sender_id) REFERENCES users(id) 
        ON DELETE SET NULL 
        ON UPDATE CASCADE
        
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='系统通知表 - 统一管理所有类型的用户通知和消息推送'
  ROW_FORMAT=DYNAMIC
  PARTITION BY RANGE (YEAR(created_at)) (
    PARTITION p2024 VALUES LESS THAN (2025),
    PARTITION p2025 VALUES LESS THAN (2026),
    PARTITION p2026 VALUES LESS THAN (2027),
    PARTITION p_future VALUES LESS THAN MAXVALUE
  );

-- ============================================================================
-- 通知表约束和检查
-- ============================================================================

-- 通知标题长度约束
ALTER TABLE notifications ADD CONSTRAINT chk_notification_title_length 
CHECK (LENGTH(TRIM(title)) >= 1 AND LENGTH(title) <= 200);

-- 重要性评分约束
ALTER TABLE notifications ADD CONSTRAINT chk_importance_score_range 
CHECK (importance_score >= 0 AND importance_score <= 100);

-- 时间逻辑约束
ALTER TABLE notifications ADD CONSTRAINT chk_notification_time_logic 
CHECK (
    scheduled_at <= COALESCE(expires_at, '2099-12-31 23:59:59') AND
    (sent_at IS NULL OR sent_at >= scheduled_at) AND
    (delivered_at IS NULL OR delivered_at >= sent_at) AND
    (read_at IS NULL OR read_at >= delivered_at) AND
    (clicked_at IS NULL OR clicked_at >= read_at)
);

-- 关联对象约束
ALTER TABLE notifications ADD CONSTRAINT chk_related_object_consistency 
CHECK (
    (related_id IS NULL AND related_type IS NULL) OR
    (related_id IS NOT NULL AND related_type IS NOT NULL)
);

-- 操作约束
ALTER TABLE notifications ADD CONSTRAINT chk_action_consistency 
CHECK (
    (action_url IS NULL AND action_text IS NULL AND action_type IS NULL) OR
    (action_url IS NOT NULL AND action_text IS NOT NULL)
);

-- 统计字段约束
ALTER TABLE notifications ADD CONSTRAINT chk_notification_counters 
CHECK (
    view_count >= 0 AND
    click_count >= 0 AND
    share_count >= 0 AND
    click_count <= view_count
);

-- 批次序号约束
ALTER TABLE notifications ADD CONSTRAINT chk_batch_sequence 
CHECK (
    (batch_id IS NULL AND batch_sequence IS NULL) OR
    (batch_id IS NOT NULL AND batch_sequence > 0)
);

-- JSON字段验证
ALTER TABLE notifications ADD CONSTRAINT chk_notifications_json_valid 
CHECK (
    (channels IS NULL OR JSON_VALID(channels)) AND
    (delivery_status IS NULL OR JSON_VALID(delivery_status)) AND
    (template_variables IS NULL OR JSON_VALID(template_variables)) AND
    (custom_styling IS NULL OR JSON_VALID(custom_styling)) AND
    (user_device IS NULL OR JSON_VALID(user_device)) AND
    (metadata IS NULL OR JSON_VALID(metadata)) AND
    (tags IS NULL OR JSON_VALID(tags)) AND
    (custom_fields IS NULL OR JSON_VALID(custom_fields))
);

-- ============================================================================
-- 通知管理触发器
-- ============================================================================

-- 通知状态自动更新触发器
DELIMITER //
CREATE TRIGGER notifications_status_update
BEFORE UPDATE ON notifications
FOR EACH ROW
BEGIN
    -- 自动设置状态转换的时间戳
    IF OLD.status != NEW.status THEN
        CASE NEW.status
            WHEN 'sent' THEN
                IF NEW.sent_at IS NULL THEN
                    SET NEW.sent_at = CURRENT_TIMESTAMP;
                END IF;
            WHEN 'delivered' THEN
                IF NEW.delivered_at IS NULL THEN
                    SET NEW.delivered_at = CURRENT_TIMESTAMP;
                END IF;
            WHEN 'read' THEN
                IF NEW.read_at IS NULL THEN
                    SET NEW.read_at = CURRENT_TIMESTAMP;
                END IF;
                SET NEW.is_read = TRUE;
            WHEN 'clicked' THEN
                IF NEW.clicked_at IS NULL THEN
                    SET NEW.clicked_at = CURRENT_TIMESTAMP;
                END IF;
        END CASE;
    END IF;
    
    -- 统计更新
    IF NEW.view_count > OLD.view_count AND OLD.status != 'read' THEN
        SET NEW.status = 'read';
        SET NEW.is_read = TRUE;
        IF NEW.read_at IS NULL THEN
            SET NEW.read_at = CURRENT_TIMESTAMP;
        END IF;
    END IF;
    
    -- 过期检查
    IF NEW.expires_at IS NOT NULL AND NEW.expires_at <= CURRENT_TIMESTAMP THEN
        SET NEW.status = 'expired';
    END IF;
END//
DELIMITER ;

-- 通知创建初始化触发器
DELIMITER //
CREATE TRIGGER notifications_initialize
BEFORE INSERT ON notifications
FOR EACH ROW
BEGIN
    -- 生成批次ID（如果是批量通知）
    IF NEW.batch_id IS NULL AND NEW.notification_type = 'system' THEN
        SET NEW.batch_id = UUID();
    END IF;
    
    -- 设置默认推送渠道
    IF NEW.channels IS NULL THEN
        SET NEW.channels = JSON_ARRAY('in_app', 'email');
    END IF;
    
    -- 初始化送达状态
    IF NEW.delivery_status IS NULL THEN
        SET NEW.delivery_status = JSON_OBJECT();
    END IF;
    
    -- 根据用户设置调整推送渠道
    -- 这里简化处理，实际应用中需要查询用户的通知偏好
    
    -- 设置用户时区（如果未指定）
    IF NEW.user_timezone IS NULL THEN
        SET NEW.user_timezone = (
            SELECT timezone 
            FROM users 
            WHERE id = NEW.user_id
        );
    END IF;
END//
DELIMITER ;

-- 通知清理触发器
DELIMITER //
CREATE TRIGGER notifications_cleanup
BEFORE DELETE ON notifications
FOR EACH ROW
BEGIN
    -- 在删除通知前，可以做一些清理工作
    -- 例如清理相关的推送记录、统计数据等
    
    -- 更新用户的未读通知数（这里简化处理）
    -- 实际应用中可能需要更复杂的统计逻辑
    
    -- 记录删除日志（如果需要）
    INSERT INTO operation_logs (
        user_id, action, resource_type, resource_id, 
        details, created_at
    ) VALUES (
        OLD.user_id, 'delete', 'notification', OLD.id,
        JSON_OBJECT('title', OLD.title, 'type', OLD.notification_type),
        CURRENT_TIMESTAMP
    );
END//
DELIMITER ;

-- ============================================================================
-- 通知管理存储过程
-- ============================================================================

-- 创建通知存储过程
DELIMITER //
CREATE PROCEDURE CreateNotification(
    IN user_id_param BIGINT UNSIGNED,
    IN sender_id_param BIGINT UNSIGNED,
    IN type_param VARCHAR(50),
    IN title_param VARCHAR(200),
    IN content_param TEXT,
    IN priority_param VARCHAR(10),
    IN related_id_param BIGINT UNSIGNED,
    IN related_type_param VARCHAR(20),
    OUT notification_id BIGINT UNSIGNED,
    OUT success BOOLEAN,
    OUT error_message VARCHAR(500)
)
BEGIN
    DECLARE user_exists INT DEFAULT 0;
    DECLARE sender_exists INT DEFAULT 0;
    
    SET success = FALSE;
    SET error_message = NULL;
    SET notification_id = NULL;
    
    -- 检查接收用户是否存在
    SELECT COUNT(*) INTO user_exists
    FROM users
    WHERE id = user_id_param AND status = 'active';
    
    IF user_exists = 0 THEN
        SET error_message = '接收用户不存在或已禁用';
        LEAVE CreateNotification;
    END IF;
    
    -- 检查发送用户是否存在（如果指定）
    IF sender_id_param IS NOT NULL THEN
        SELECT COUNT(*) INTO sender_exists
        FROM users
        WHERE id = sender_id_param AND status = 'active';
        
        IF sender_exists = 0 THEN
            SET error_message = '发送用户不存在或已禁用';
            LEAVE CreateNotification;
        END IF;
    END IF;
    
    -- 创建通知
    INSERT INTO notifications (
        user_id, sender_id, notification_type, title, content,
        priority, related_id, related_type
    ) VALUES (
        user_id_param, sender_id_param, type_param, title_param, content_param,
        priority_param, related_id_param, related_type_param
    );
    
    SET notification_id = LAST_INSERT_ID();
    SET success = TRUE;
    
END//
DELIMITER ;

-- 批量标记已读存储过程
DELIMITER //
CREATE PROCEDURE MarkNotificationsAsRead(
    IN user_id_param BIGINT UNSIGNED,
    IN notification_ids JSON,
    OUT affected_count INT
)
BEGIN
    DECLARE i INT DEFAULT 0;
    DECLARE notification_count INT DEFAULT 0;
    DECLARE current_id BIGINT UNSIGNED;
    
    SET affected_count = 0;
    
    IF notification_ids IS NOT NULL THEN
        SET notification_count = JSON_LENGTH(notification_ids);
        
        -- 遍历通知ID列表
        WHILE i < notification_count DO
            SET current_id = JSON_UNQUOTE(JSON_EXTRACT(notification_ids, CONCAT('$[', i, ']')));
            
            UPDATE notifications 
            SET is_read = TRUE,
                read_at = CURRENT_TIMESTAMP,
                status = 'read',
                view_count = view_count + 1
            WHERE id = current_id 
              AND user_id = user_id_param 
              AND is_read = FALSE;
            
            SET affected_count = affected_count + ROW_COUNT();
            SET i = i + 1;
        END WHILE;
    ELSE
        -- 如果没有指定ID，标记所有未读通知为已读
        UPDATE notifications 
        SET is_read = TRUE,
            read_at = CURRENT_TIMESTAMP,
            status = 'read',
            view_count = view_count + 1
        WHERE user_id = user_id_param AND is_read = FALSE;
        
        SET affected_count = ROW_COUNT();
    END IF;
    
END//
DELIMITER ;

-- 清理过期通知存储过程
DELIMITER //
CREATE PROCEDURE CleanExpiredNotifications()
BEGIN
    DECLARE affected_rows INT DEFAULT 0;
    
    -- 标记过期通知
    UPDATE notifications 
    SET status = 'expired'
    WHERE expires_at IS NOT NULL 
      AND expires_at <= CURRENT_TIMESTAMP 
      AND status NOT IN ('expired', 'read');
    
    SET affected_rows = ROW_COUNT();
    
    -- 删除90天前的已过期通知
    DELETE FROM notifications 
    WHERE status = 'expired' 
      AND expires_at <= DATE_SUB(CURRENT_TIMESTAMP, INTERVAL 90 DAY);
    
    SELECT CONCAT('已处理 ', affected_rows, ' 个过期通知，删除 ', ROW_COUNT(), ' 个历史通知') as result;
END//
DELIMITER ;

-- ============================================================================
-- 通知统计视图
-- ============================================================================

-- 用户通知统计视图
CREATE VIEW user_notification_stats AS
SELECT 
    n.user_id,
    u.username,
    COUNT(*) as total_notifications,
    COUNT(CASE WHEN n.is_read = FALSE THEN 1 END) as unread_count,
    COUNT(CASE WHEN n.is_starred = TRUE THEN 1 END) as starred_count,
    COUNT(CASE WHEN n.is_archived = TRUE THEN 1 END) as archived_count,
    COUNT(CASE WHEN n.priority = 'urgent' THEN 1 END) as urgent_count,
    COUNT(CASE WHEN n.priority = 'high' THEN 1 END) as high_priority_count,
    MAX(n.created_at) as last_notification_time,
    AVG(n.importance_score) as avg_importance_score
FROM notifications n
JOIN users u ON n.user_id = u.id
WHERE n.deleted_at IS NULL
GROUP BY n.user_id, u.username;

-- 通知类型统计视图
CREATE VIEW notification_type_stats AS
SELECT 
    notification_type,
    notification_category,
    COUNT(*) as total_count,
    COUNT(CASE WHEN is_read = TRUE THEN 1 END) as read_count,
    COUNT(CASE WHEN status = 'delivered' THEN 1 END) as delivered_count,
    ROUND(COUNT(CASE WHEN is_read = TRUE THEN 1 END) * 100.0 / COUNT(*), 2) as read_rate,
    AVG(importance_score) as avg_importance,
    AVG(click_count) as avg_clicks,
    COUNT(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY) THEN 1 END) as recent_count
FROM notifications
WHERE deleted_at IS NULL
GROUP BY notification_type, notification_category
ORDER BY total_count DESC;

-- 高优先级未读通知视图
CREATE VIEW urgent_unread_notifications AS
SELECT 
    n.id,
    n.user_id,
    u.username,
    u.email,
    n.title,
    n.content,
    n.notification_type,
    n.priority,
    n.importance_score,
    n.created_at,
    n.expires_at,
    TIMESTAMPDIFF(HOUR, n.created_at, NOW()) as hours_since_created
FROM notifications n
JOIN users u ON n.user_id = u.id
WHERE n.is_read = FALSE
  AND n.priority IN ('urgent', 'critical', 'high')
  AND n.status NOT IN ('expired', 'failed')
  AND n.deleted_at IS NULL
  AND (n.expires_at IS NULL OR n.expires_at > NOW())
ORDER BY n.importance_score DESC, n.created_at DESC;