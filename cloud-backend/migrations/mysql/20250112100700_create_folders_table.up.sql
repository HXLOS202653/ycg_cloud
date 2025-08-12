-- +migrate Up
-- 创建迁移: 文件夹扩展属性表
-- 版本: 20250112100700
-- 描述: 创建文件夹扩展属性表，管理文件夹的高级功能和特殊属性
-- 作者: 系统自动生成
-- 创建时间: 2025-01-12 10:07:00
-- 依赖: 20250112100300_create_files_table
-- 数据库版本要求: MySQL 8.0.31+
-- GORM版本: 1.30.1
-- 说明: 此表是files表的扩展，专门存储文件夹类型记录的额外属性和功能

-- ============================================================================
-- 文件夹扩展属性表 (folders)
-- ============================================================================

CREATE TABLE folders (
    -- 基础标识字段（与files表关联）
    id BIGINT UNSIGNED PRIMARY KEY COMMENT '文件夹ID，直接关联files表的主键',
    
    -- 文件夹类型和分类
    folder_type ENUM('normal', 'shared', 'team', 'system', 'template', 'backup', 'recycle') DEFAULT 'normal' COMMENT '文件夹类型',
    folder_category VARCHAR(50) DEFAULT NULL COMMENT '文件夹分类：work、personal、project、archive等',
    
    -- 文件夹容量和统计信息
    file_count INT UNSIGNED DEFAULT 0 COMMENT '直接子文件数量（不包含子文件夹）',
    folder_count INT UNSIGNED DEFAULT 0 COMMENT '直接子文件夹数量',
    total_file_count INT UNSIGNED DEFAULT 0 COMMENT '总文件数量（包含所有子目录）',
    total_folder_count INT UNSIGNED DEFAULT 0 COMMENT '总文件夹数量（包含所有子目录）',
    total_size BIGINT UNSIGNED DEFAULT 0 COMMENT '文件夹总大小（包含所有子目录，字节）',
    
    -- 文件夹权限和安全设置
    default_permission ENUM('private', 'shared', 'public', 'team') DEFAULT 'private' COMMENT '默认权限级别',
    inherit_permissions BOOLEAN DEFAULT TRUE COMMENT '是否继承父文件夹权限',
    auto_share_subfolder BOOLEAN DEFAULT FALSE COMMENT '是否自动分享子文件夹',
    
    -- 文件夹显示和布局设置
    view_mode ENUM('grid', 'list', 'detail', 'timeline') DEFAULT 'grid' COMMENT '默认显示模式',
    sort_order ENUM('name_asc', 'name_desc', 'size_asc', 'size_desc', 'date_asc', 'date_desc', 'type_asc', 'type_desc') DEFAULT 'name_asc' COMMENT '默认排序方式',
    cover_image_url VARCHAR(1000) DEFAULT NULL COMMENT '文件夹封面图片URL',
    color_theme VARCHAR(20) DEFAULT NULL COMMENT '文件夹主题颜色：blue、green、red、purple等',
    icon_type VARCHAR(50) DEFAULT 'folder' COMMENT '文件夹图标类型',
    
    -- 文件夹功能设置
    enable_auto_backup BOOLEAN DEFAULT FALSE COMMENT '是否启用自动备份',
    backup_schedule VARCHAR(100) DEFAULT NULL COMMENT '备份计划表达式（cron格式）',
    enable_version_control BOOLEAN DEFAULT TRUE COMMENT '是否启用版本控制',
    max_file_versions INT UNSIGNED DEFAULT 10 COMMENT '最大保留版本数',
    
    -- 文件夹同步设置
    sync_enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用同步',
    sync_priority ENUM('low', 'normal', 'high', 'urgent') DEFAULT 'normal' COMMENT '同步优先级',
    exclude_patterns JSON DEFAULT NULL COMMENT '同步排除模式，如文件名模式、扩展名等',
    
    -- 文件夹监控和通知
    monitor_changes BOOLEAN DEFAULT FALSE COMMENT '是否监控文件变更',
    notification_settings JSON DEFAULT NULL COMMENT '通知设置：新文件、修改、删除等',
    
    -- 文件夹限制和配额
    max_file_size BIGINT UNSIGNED DEFAULT NULL COMMENT '单文件最大大小限制（字节），NULL表示继承用户设置',
    max_file_count INT UNSIGNED DEFAULT NULL COMMENT '最大文件数量限制，NULL表示无限制',
    max_folder_size BIGINT UNSIGNED DEFAULT NULL COMMENT '文件夹最大总大小（字节），NULL表示无限制',
    allowed_file_types JSON DEFAULT NULL COMMENT '允许的文件类型，NULL表示无限制',
    forbidden_file_types JSON DEFAULT NULL COMMENT '禁止的文件类型',
    
    -- 文件夹模板和规则
    is_template BOOLEAN DEFAULT FALSE COMMENT '是否为模板文件夹',
    template_name VARCHAR(100) DEFAULT NULL COMMENT '模板名称',
    auto_organize_rules JSON DEFAULT NULL COMMENT '自动整理规则：按日期、类型、大小等',
    naming_rules JSON DEFAULT NULL COMMENT '文件命名规则和约定',
    
    -- 文件夹协作功能
    collaboration_enabled BOOLEAN DEFAULT FALSE COMMENT '是否启用协作功能',
    discussion_enabled BOOLEAN DEFAULT FALSE COMMENT '是否启用讨论功能',
    task_management BOOLEAN DEFAULT FALSE COMMENT '是否启用任务管理',
    
    -- 文件夹生命周期
    auto_cleanup_enabled BOOLEAN DEFAULT FALSE COMMENT '是否启用自动清理',
    cleanup_rules JSON DEFAULT NULL COMMENT '清理规则：文件年龄、大小、访问时间等',
    archive_threshold_days INT UNSIGNED DEFAULT NULL COMMENT '归档阈值天数',
    
    -- 文件夹搜索和索引
    enable_fulltext_search BOOLEAN DEFAULT TRUE COMMENT '是否启用全文搜索',
    search_weight DECIMAL(3,2) DEFAULT 1.00 COMMENT '搜索权重（0.01-9.99）',
    indexed_at TIMESTAMP NULL DEFAULT NULL COMMENT '最后索引时间',
    
    -- 文件夹缓存设置
    cache_enabled BOOLEAN DEFAULT TRUE COMMENT '是否启用缓存',
    cache_ttl INT UNSIGNED DEFAULT 3600 COMMENT '缓存生存时间（秒）',
    
    -- 扩展属性和元数据
    custom_fields JSON DEFAULT NULL COMMENT '自定义字段和属性',
    metadata JSON DEFAULT NULL COMMENT '文件夹元数据',
    workflow_settings JSON DEFAULT NULL COMMENT '工作流设置',
    
    -- 时间戳
    stats_updated_at TIMESTAMP NULL DEFAULT NULL COMMENT '统计信息最后更新时间',
    settings_updated_at TIMESTAMP NULL DEFAULT NULL COMMENT '设置最后更新时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
    
    -- 业务索引设计
    INDEX idx_folders_folder_type (folder_type) COMMENT '文件夹类型索引',
    INDEX idx_folders_folder_category (folder_category) COMMENT '文件夹分类索引',
    INDEX idx_folders_total_size (total_size) COMMENT '总大小索引，用于空间统计',
    INDEX idx_folders_file_count (file_count) COMMENT '文件数量索引',
    INDEX idx_folders_default_permission (default_permission) COMMENT '默认权限索引',
    INDEX idx_folders_is_template (is_template) COMMENT '模板文件夹索引',
    INDEX idx_folders_collaboration_enabled (collaboration_enabled) COMMENT '协作功能索引',
    INDEX idx_folders_stats_updated_at (stats_updated_at) COMMENT '统计更新时间索引',
    INDEX idx_folders_created_at (created_at) COMMENT '创建时间索引',
    
    -- 复合业务索引
    INDEX idx_folders_type_category (folder_type, folder_category) COMMENT '类型分类复合索引',
    INDEX idx_folders_size_count (total_size, total_file_count) COMMENT '大小文件数复合索引',
    INDEX idx_folders_collaboration_features (collaboration_enabled, discussion_enabled, task_management) COMMENT '协作功能复合索引',
    
    -- 外键约束：确保该ID在files表中存在且为文件夹
    CONSTRAINT fk_folders_id 
        FOREIGN KEY (id) REFERENCES files(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE
        
) ENGINE=InnoDB 
  DEFAULT CHARSET=utf8mb4 
  COLLATE=utf8mb4_unicode_ci 
  COMMENT='文件夹扩展属性表 - 存储文件夹的高级功能、权限设置和统计信息'
  ROW_FORMAT=DYNAMIC;

-- ============================================================================
-- 文件夹表约束和触发器
-- ============================================================================

-- 确保只有files表中的文件夹记录才能在此表中有对应记录
ALTER TABLE folders ADD CONSTRAINT chk_folders_is_folder 
CHECK (
    id IN (SELECT id FROM files WHERE is_folder = TRUE)
);

-- 文件和文件夹数量约束
ALTER TABLE folders ADD CONSTRAINT chk_folder_counts 
CHECK (
    file_count >= 0 AND 
    folder_count >= 0 AND 
    total_file_count >= file_count AND 
    total_folder_count >= folder_count
);

-- 大小约束
ALTER TABLE folders ADD CONSTRAINT chk_folder_sizes 
CHECK (
    total_size >= 0 AND
    (max_file_size IS NULL OR max_file_size > 0) AND
    (max_folder_size IS NULL OR max_folder_size > 0)
);

-- 版本数量约束
ALTER TABLE folders ADD CONSTRAINT chk_max_versions 
CHECK (max_file_versions > 0 AND max_file_versions <= 1000);

-- 搜索权重约束
ALTER TABLE folders ADD CONSTRAINT chk_search_weight 
CHECK (search_weight >= 0.01 AND search_weight <= 9.99);

-- 缓存TTL约束
ALTER TABLE folders ADD CONSTRAINT chk_cache_ttl 
CHECK (cache_ttl > 0 AND cache_ttl <= 86400); -- 最长1天

-- JSON字段验证
ALTER TABLE folders ADD CONSTRAINT chk_folders_json_valid 
CHECK (
    (exclude_patterns IS NULL OR JSON_VALID(exclude_patterns)) AND
    (notification_settings IS NULL OR JSON_VALID(notification_settings)) AND
    (allowed_file_types IS NULL OR JSON_VALID(allowed_file_types)) AND
    (forbidden_file_types IS NULL OR JSON_VALID(forbidden_file_types)) AND
    (auto_organize_rules IS NULL OR JSON_VALID(auto_organize_rules)) AND
    (naming_rules IS NULL OR JSON_VALID(naming_rules)) AND
    (cleanup_rules IS NULL OR JSON_VALID(cleanup_rules)) AND
    (custom_fields IS NULL OR JSON_VALID(custom_fields)) AND
    (metadata IS NULL OR JSON_VALID(metadata)) AND
    (workflow_settings IS NULL OR JSON_VALID(workflow_settings))
);

-- 创建文件夹统计自动更新触发器
DELIMITER //
CREATE TRIGGER folders_stats_update
AFTER INSERT ON files
FOR EACH ROW
BEGIN
    -- 当在文件夹中添加文件或子文件夹时，更新父文件夹统计
    IF NEW.parent_id IS NOT NULL THEN
        IF NEW.is_folder = TRUE THEN
            -- 添加了子文件夹
            UPDATE folders 
            SET folder_count = folder_count + 1,
                total_folder_count = total_folder_count + 1,
                stats_updated_at = CURRENT_TIMESTAMP
            WHERE id = NEW.parent_id;
        ELSE
            -- 添加了文件
            UPDATE folders 
            SET file_count = file_count + 1,
                total_file_count = total_file_count + 1,
                total_size = total_size + NEW.file_size,
                stats_updated_at = CURRENT_TIMESTAMP
            WHERE id = NEW.parent_id;
        END IF;
        
        -- 递归更新所有父级文件夹的总计数
        CALL UpdateParentFolderStats(NEW.parent_id, NEW.is_folder, NEW.file_size, 1);
    END IF;
END//
DELIMITER ;

-- 创建文件夹统计删除更新触发器
DELIMITER //
CREATE TRIGGER folders_stats_delete
AFTER UPDATE ON files
FOR EACH ROW
BEGIN
    -- 处理软删除对文件夹统计的影响
    IF OLD.is_deleted = FALSE AND NEW.is_deleted = TRUE AND NEW.parent_id IS NOT NULL THEN
        IF NEW.is_folder = TRUE THEN
            -- 删除了子文件夹
            UPDATE folders 
            SET folder_count = GREATEST(folder_count - 1, 0),
                total_folder_count = GREATEST(total_folder_count - 1, 0),
                stats_updated_at = CURRENT_TIMESTAMP
            WHERE id = NEW.parent_id;
        ELSE
            -- 删除了文件
            UPDATE folders 
            SET file_count = GREATEST(file_count - 1, 0),
                total_file_count = GREATEST(total_file_count - 1, 0),
                total_size = GREATEST(total_size - NEW.file_size, 0),
                stats_updated_at = CURRENT_TIMESTAMP
            WHERE id = NEW.parent_id;
        END IF;
        
        -- 递归更新所有父级文件夹的总计数
        CALL UpdateParentFolderStats(NEW.parent_id, NEW.is_folder, NEW.file_size, -1);
        
    ELSEIF OLD.is_deleted = TRUE AND NEW.is_deleted = FALSE AND NEW.parent_id IS NOT NULL THEN
        -- 恢复文件/文件夹
        IF NEW.is_folder = TRUE THEN
            UPDATE folders 
            SET folder_count = folder_count + 1,
                total_folder_count = total_folder_count + 1,
                stats_updated_at = CURRENT_TIMESTAMP
            WHERE id = NEW.parent_id;
        ELSE
            UPDATE folders 
            SET file_count = file_count + 1,
                total_file_count = total_file_count + 1,
                total_size = total_size + NEW.file_size,
                stats_updated_at = CURRENT_TIMESTAMP
            WHERE id = NEW.parent_id;
        END IF;
        
        -- 递归更新所有父级文件夹的总计数
        CALL UpdateParentFolderStats(NEW.parent_id, NEW.is_folder, NEW.file_size, 1);
    END IF;
END//
DELIMITER ;

-- 创建文件夹权限继承触发器
DELIMITER //
CREATE TRIGGER folders_permission_inherit
AFTER INSERT ON folders
FOR EACH ROW
BEGIN
    -- 如果启用权限继承，从父文件夹继承权限
    IF NEW.inherit_permissions = TRUE THEN
        SET @parent_folder_id = (SELECT parent_id FROM files WHERE id = NEW.id);
        
        IF @parent_folder_id IS NOT NULL THEN
            -- 复制父文件夹的权限设置
            UPDATE folders 
            SET default_permission = (SELECT default_permission FROM folders WHERE id = @parent_folder_id),
                auto_share_subfolder = (SELECT auto_share_subfolder FROM folders WHERE id = @parent_folder_id)
            WHERE id = NEW.id;
        END IF;
    END IF;
END//
DELIMITER ;

-- ============================================================================
-- 文件夹管理存储过程
-- ============================================================================

-- 创建递归更新父文件夹统计的存储过程
DELIMITER //
CREATE PROCEDURE UpdateParentFolderStats(
    IN folder_id BIGINT UNSIGNED,
    IN is_folder_item BOOLEAN,
    IN file_size BIGINT UNSIGNED,
    IN count_delta INT
)
BEGIN
    DECLARE parent_folder_id BIGINT UNSIGNED;
    DECLARE done INT DEFAULT FALSE;
    
    -- 获取当前文件夹的父文件夹ID
    SELECT parent_id INTO parent_folder_id 
    FROM files 
    WHERE id = folder_id AND is_folder = TRUE;
    
    -- 递归更新父级文件夹统计
    WHILE parent_folder_id IS NOT NULL DO
        IF is_folder_item THEN
            UPDATE folders 
            SET total_folder_count = GREATEST(total_folder_count + count_delta, 0),
                stats_updated_at = CURRENT_TIMESTAMP
            WHERE id = parent_folder_id;
        ELSE
            UPDATE folders 
            SET total_file_count = GREATEST(total_file_count + count_delta, 0),
                total_size = GREATEST(total_size + (file_size * count_delta), 0),
                stats_updated_at = CURRENT_TIMESTAMP
            WHERE id = parent_folder_id;
        END IF;
        
        -- 获取下一个父级文件夹
        SELECT parent_id INTO parent_folder_id 
        FROM files 
        WHERE id = parent_folder_id AND is_folder = TRUE;
    END WHILE;
END//
DELIMITER ;

-- 创建重新计算文件夹统计的存储过程
DELIMITER //
CREATE PROCEDURE RecalculateFolderStats(
    IN target_folder_id BIGINT UNSIGNED
)
BEGIN
    DECLARE direct_file_count INT DEFAULT 0;
    DECLARE direct_folder_count INT DEFAULT 0;
    DECLARE total_size_calc BIGINT DEFAULT 0;
    
    -- 计算直接子文件统计
    SELECT 
        COUNT(CASE WHEN is_folder = FALSE THEN 1 END),
        COUNT(CASE WHEN is_folder = TRUE THEN 1 END),
        COALESCE(SUM(CASE WHEN is_folder = FALSE THEN file_size END), 0)
    INTO direct_file_count, direct_folder_count, total_size_calc
    FROM files 
    WHERE parent_id = target_folder_id AND is_deleted = FALSE;
    
    -- 更新直接统计
    UPDATE folders 
    SET file_count = direct_file_count,
        folder_count = direct_folder_count,
        stats_updated_at = CURRENT_TIMESTAMP
    WHERE id = target_folder_id;
    
    -- 递归计算总计统计（这里简化处理，实际可能需要更复杂的递归逻辑）
    UPDATE folders 
    SET total_file_count = (
            SELECT COUNT(*) 
            FROM files 
            WHERE file_path LIKE CONCAT((SELECT file_path FROM files WHERE id = target_folder_id), '%')
              AND is_folder = FALSE 
              AND is_deleted = FALSE
        ),
        total_folder_count = (
            SELECT COUNT(*) 
            FROM files 
            WHERE file_path LIKE CONCAT((SELECT file_path FROM files WHERE id = target_folder_id), '%')
              AND is_folder = TRUE 
              AND is_deleted = FALSE
              AND id != target_folder_id
        ),
        total_size = (
            SELECT COALESCE(SUM(file_size), 0)
            FROM files 
            WHERE file_path LIKE CONCAT((SELECT file_path FROM files WHERE id = target_folder_id), '%')
              AND is_folder = FALSE 
              AND is_deleted = FALSE
        )
    WHERE id = target_folder_id;
    
    SELECT CONCAT('文件夹 ', target_folder_id, ' 统计信息已重新计算') as result;
END//
DELIMITER ;

-- 创建批量初始化文件夹记录的存储过程
DELIMITER //
CREATE PROCEDURE InitializeFolderRecords()
BEGIN
    DECLARE affected_rows INT DEFAULT 0;
    
    -- 为所有现有的文件夹创建folders表记录
    INSERT INTO folders (id, folder_type, created_at)
    SELECT id, 
           CASE 
               WHEN user_id = 1 THEN 'system'  -- 假设用户ID 1是系统用户
               WHEN is_public = TRUE THEN 'shared'
               ELSE 'normal'
           END as folder_type,
           created_at
    FROM files 
    WHERE is_folder = TRUE 
      AND is_deleted = FALSE
      AND id NOT IN (SELECT id FROM folders);
    
    SET affected_rows = ROW_COUNT();
    
    SELECT CONCAT('已初始化 ', affected_rows, ' 个文件夹记录') as result;
END//
DELIMITER ;

-- ============================================================================
-- 文件夹统计和管理视图
-- ============================================================================

-- 创建文件夹详细统计视图
CREATE VIEW folder_detailed_stats AS
SELECT 
    f.id,
    f.filename as folder_name,
    f.file_path,
    f.user_id,
    u.username,
    fd.folder_type,
    fd.folder_category,
    fd.file_count,
    fd.folder_count,
    fd.total_file_count,
    fd.total_folder_count,
    fd.total_size,
    ROUND(fd.total_size / 1024 / 1024, 2) as total_size_mb,
    fd.default_permission,
    fd.collaboration_enabled,
    fd.stats_updated_at,
    f.created_at,
    f.updated_at
FROM files f
JOIN folders fd ON f.id = fd.id
LEFT JOIN users u ON f.user_id = u.id
WHERE f.is_folder = TRUE AND f.is_deleted = FALSE;

-- 创建文件夹容量排行视图
CREATE VIEW folder_size_ranking AS
SELECT 
    f.id,
    f.filename as folder_name,
    f.file_path,
    u.username,
    fd.total_size,
    ROUND(fd.total_size / 1024 / 1024, 2) as size_mb,
    fd.total_file_count,
    fd.total_folder_count,
    RANK() OVER (ORDER BY fd.total_size DESC) as size_rank
FROM files f
JOIN folders fd ON f.id = fd.id
LEFT JOIN users u ON f.user_id = u.id
WHERE f.is_folder = TRUE 
  AND f.is_deleted = FALSE 
  AND fd.total_size > 0
ORDER BY fd.total_size DESC;

-- 创建协作文件夹活跃度视图
CREATE VIEW collaborative_folders AS
SELECT 
    f.id,
    f.filename as folder_name,
    f.user_id,
    u.username,
    fd.collaboration_enabled,
    fd.discussion_enabled,
    fd.task_management,
    fd.total_file_count,
    fd.stats_updated_at,
    f.updated_at as last_activity
FROM files f
JOIN folders fd ON f.id = fd.id
LEFT JOIN users u ON f.user_id = u.id
WHERE f.is_folder = TRUE 
  AND f.is_deleted = FALSE 
  AND fd.collaboration_enabled = TRUE
ORDER BY f.updated_at DESC;

-- ============================================================================
-- 初始化现有文件夹
-- ============================================================================

-- 为所有现有的文件夹创建folders表记录
CALL InitializeFolderRecords();