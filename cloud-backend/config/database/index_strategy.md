# 数据库索引策略配置

## 📋 索引设计原则

基于数据库设计文档（07数据库设计.md），本系统的索引策略遵循以下原则：

### 1. 核心设计原则
- **性能优化**: 合理的索引设计和分区策略
- **查询效率**: 基于业务查询模式设计复合索引
- **写入性能**: 避免过多索引影响写入性能
- **存储优化**: 控制索引大小，避免占用过多磁盘空间
- **命名规范**: 严格遵循 `idx_tablename_fieldname` 规范

### 2. 索引类型策略

#### 2.1 主键索引 (PRIMARY KEY)
- 所有表使用 BIGINT UNSIGNED AUTO_INCREMENT 主键
- 主键自动创建聚簇索引，优化范围查询

#### 2.2 唯一索引 (UNIQUE KEY)
- 业务唯一性约束：用户名、邮箱、分享码等
- 防止重复数据，提供快速查找
- 命名规范：`uk_tablename_fieldname`

#### 2.3 普通索引 (INDEX)
- 高频查询字段：用户ID、状态、时间等
- 外键字段：关联表查询优化
- 命名规范：`idx_tablename_fieldname`

#### 2.4 复合索引 (COMPOSITE INDEX)
- 多字段联合查询优化
- 遵循最左前缀原则
- 命名规范：`idx_tablename_descriptive_name`

#### 2.5 全文索引 (FULLTEXT INDEX)
- 文件名、内容、描述等文本搜索
- 使用 ngram 解析器支持中文
- 命名规范：`idx_tablename_fulltext_description`

## 🗄️ 核心表索引策略

### 1. 用户管理模块

#### users 表索引策略
```sql
-- 登录验证索引
INDEX idx_users_username (username) COMMENT '用户名登录查询',
INDEX idx_users_email (email) COMMENT '邮箱登录查询',
INDEX idx_users_phone (phone) COMMENT '手机号登录查询',

-- 状态筛选索引
INDEX idx_users_status (status) COMMENT '用户状态筛选',
INDEX idx_users_role (role) COMMENT '角色权限查询',

-- 时间序列索引
INDEX idx_users_created_at (created_at) COMMENT '注册时间排序',
INDEX idx_users_last_login_at (last_login_at) COMMENT '活跃度统计',
INDEX idx_users_deleted_at (deleted_at) COMMENT '软删除查询',

-- 存储管理索引
INDEX idx_users_storage_used (storage_used) COMMENT '存储配额管理',
INDEX idx_users_locked_until (locked_until) COMMENT '账户解锁查询'
```

#### user_sessions 表索引策略
```sql
-- 会话管理索引
INDEX idx_user_sessions_user_id (user_id) COMMENT '用户会话查询',
INDEX idx_user_sessions_session_token (session_token) COMMENT 'Token验证',
INDEX idx_user_sessions_expires_at (expires_at) COMMENT '过期会话清理',
INDEX idx_user_sessions_is_active (is_active) COMMENT '活跃会话筛选',
INDEX idx_user_sessions_last_activity_at (last_activity_at) COMMENT '空闲检测',
INDEX idx_user_sessions_ip_address (ip_address) COMMENT '安全分析',

-- 复合业务索引
INDEX idx_user_sessions_user_active (user_id, is_active, expires_at) COMMENT '用户有效会话查询'
```

### 2. 文件管理模块

#### files 表索引策略
```sql
-- 核心业务索引
INDEX idx_files_user_parent_deleted (user_id, parent_id, is_deleted) COMMENT '文件树查询核心索引',
INDEX idx_files_filename (filename) COMMENT '文件名搜索',
INDEX idx_files_file_path (file_path(255)) COMMENT '路径查询',

-- 文件分类索引
INDEX idx_files_file_type (file_type) COMMENT '文件类型筛选',
INDEX idx_files_mime_type (mime_type) COMMENT 'MIME类型查询',
INDEX idx_files_file_extension (file_extension) COMMENT '扩展名索引',

-- 去重和完整性索引
INDEX idx_files_md5_hash (md5_hash) COMMENT 'MD5去重查询',
INDEX idx_files_sha256_hash (sha256_hash) COMMENT 'SHA256校验',

-- 状态筛选索引
INDEX idx_files_is_folder (is_folder) COMMENT '文件夹筛选',
INDEX idx_files_is_deleted (is_deleted) COMMENT '软删除筛选',
INDEX idx_files_is_favorite (is_favorite) COMMENT '收藏文件',
INDEX idx_files_is_public (is_public) COMMENT '公开文件',

-- 时间相关索引
INDEX idx_files_deleted_at (deleted_at) COMMENT '回收站查询',
INDEX idx_files_created_at (created_at) COMMENT '创建时间排序',
INDEX idx_files_updated_at (updated_at) COMMENT '修改时间',
INDEX idx_files_last_accessed_at (last_accessed_at) COMMENT '访问热度统计',

-- 统计和管理索引
INDEX idx_files_file_size (file_size) COMMENT '文件大小排序',
INDEX idx_files_download_count (download_count) COMMENT '下载热度',
INDEX idx_files_view_count (view_count) COMMENT '查看统计',

-- 安全相关索引
INDEX idx_files_virus_scan_status (virus_scan_status) COMMENT '病毒扫描状态',
INDEX idx_files_access_level (access_level) COMMENT '访问级别',

-- 复合业务索引
INDEX idx_files_user_type_deleted (user_id, file_type, is_deleted) COMMENT '用户文件类型筛选',
INDEX idx_files_parent_folder_name (parent_id, is_folder, filename) COMMENT '文件夹内容浏览',
INDEX idx_files_user_favorite (user_id, is_favorite, updated_at) COMMENT '用户收藏文件',

-- 全文搜索索引
FULLTEXT idx_files_fulltext_search (filename, description) WITH PARSER ngram COMMENT '文件名描述全文搜索'
```

### 3. 文件版本管理

#### file_versions 表索引策略
```sql
-- 版本管理索引
INDEX idx_file_versions_file_id (file_id) COMMENT '文件版本查询',
INDEX idx_file_versions_version_number (version_number) COMMENT '版本号排序',
INDEX idx_file_versions_created_by (created_by) COMMENT '创建者查询',
INDEX idx_file_versions_is_milestone (is_milestone) COMMENT '里程碑版本',
INDEX idx_file_versions_is_active (is_active) COMMENT '当前版本',
INDEX idx_file_versions_created_at (created_at) COMMENT '版本时间',

-- 完整性索引
INDEX idx_file_versions_md5_hash (md5_hash) COMMENT '版本去重',
INDEX idx_file_versions_change_type (change_type) COMMENT '变更类型',
INDEX idx_file_versions_expires_at (expires_at) COMMENT '版本过期清理',

-- 复合索引
INDEX idx_file_versions_file_active_version (file_id, is_active, version_number DESC) COMMENT '文件当前版本查询',
INDEX idx_file_versions_file_milestone (file_id, is_milestone, created_at DESC) COMMENT '里程碑版本查询',
INDEX idx_file_versions_user_versions (created_by, created_at DESC) COMMENT '用户版本历史'
```

### 4. 上传任务管理

#### upload_tasks 表索引策略
```sql
-- 任务管理索引
INDEX idx_upload_tasks_upload_id (upload_id) COMMENT '上传任务查询',
INDEX idx_upload_tasks_user_id (user_id) COMMENT '用户上传任务',
INDEX idx_upload_tasks_status (status) COMMENT '任务状态筛选',
INDEX idx_upload_tasks_md5_hash (md5_hash) COMMENT '秒传支持',
INDEX idx_upload_tasks_expires_at (expires_at) COMMENT '任务过期清理',
INDEX idx_upload_tasks_created_at (created_at) COMMENT '任务时间排序',

-- 会话和安全索引
INDEX idx_upload_tasks_upload_token (upload_token) COMMENT '上传令牌验证',
INDEX idx_upload_tasks_session_id (session_id) COMMENT '会话管理',
INDEX idx_upload_tasks_client_ip (client_ip) COMMENT '安全分析',

-- 复合业务索引
INDEX idx_upload_tasks_user_status_created (user_id, status, created_at DESC) COMMENT '用户任务状态查询',
INDEX idx_upload_tasks_status_expires (status, expires_at) COMMENT '状态过期查询',
INDEX idx_upload_tasks_user_progress (user_id, upload_percentage, updated_at DESC) COMMENT '上传进度查询',
INDEX idx_upload_tasks_hash_size (md5_hash, file_size) COMMENT '文件去重优化'
```

#### file_chunks 表索引策略
```sql
-- 分片管理索引
INDEX idx_file_chunks_upload_id (upload_id) COMMENT '上传任务分片查询',
INDEX idx_file_chunks_chunk_number (chunk_number) COMMENT '分片序号排序',
INDEX idx_file_chunks_status (status) COMMENT '分片状态筛选',
INDEX idx_file_chunks_chunk_md5 (chunk_md5) COMMENT '分片去重检测',
INDEX idx_file_chunks_worker_id (worker_id) COMMENT '并发控制',

-- 时间管理索引
INDEX idx_file_chunks_created_at (created_at) COMMENT '分片创建时间',
INDEX idx_file_chunks_uploaded_at (uploaded_at) COMMENT '上传完成时间',
INDEX idx_file_chunks_url_expires (url_expires_at) COMMENT 'URL过期清理',
INDEX idx_file_chunks_lock_expires (lock_expires_at) COMMENT '锁过期清理',
INDEX idx_file_chunks_locked_at (locked_at) COMMENT '并发锁控制',

-- 复合业务索引
INDEX idx_file_chunks_upload_status_chunk (upload_id, status, chunk_number) COMMENT '任务分片状态查询',
INDEX idx_file_chunks_status_retry (status, retry_count, updated_at) COMMENT '重试任务查询',
INDEX idx_file_chunks_upload_completed (upload_id, status, uploaded_at) COMMENT '完成分片统计',
INDEX idx_file_chunks_worker_status (worker_id, status, locked_at) COMMENT '工作进程任务查询'
```

### 5. 分享管理模块

#### file_shares 表索引策略
```sql
-- 分享核心索引
INDEX idx_file_shares_share_code (share_code) COMMENT '分享码访问查询',
INDEX idx_file_shares_file_id (file_id) COMMENT '文件分享查询',
INDEX idx_file_shares_folder_id (folder_id) COMMENT '文件夹分享查询',
INDEX idx_file_shares_user_id (user_id) COMMENT '用户分享管理',

-- 分享类型和状态索引
INDEX idx_file_shares_share_type (share_type) COMMENT '分享类型筛选',
INDEX idx_file_shares_access_type (access_type) COMMENT '访问类型控制',
INDEX idx_file_shares_status (status) COMMENT '分享状态管理',
INDEX idx_file_shares_is_active (is_active) COMMENT '有效分享筛选',

-- 时间管理索引
INDEX idx_file_shares_expires_at (expires_at) COMMENT '过期分享清理',
INDEX idx_file_shares_created_at (created_at) COMMENT '分享时间排序',
INDEX idx_file_shares_last_accessed_at (last_accessed_at) COMMENT '访问活跃度',

-- 统计索引
INDEX idx_file_shares_download_count (download_count) COMMENT '下载热度统计',
INDEX idx_file_shares_view_count (view_count) COMMENT '查看统计',
INDEX idx_file_shares_deleted_at (deleted_at) COMMENT '软删除管理',

-- 复合业务索引
INDEX idx_file_shares_user_status (user_id, status, created_at DESC) COMMENT '用户分享状态查询',
INDEX idx_file_shares_type_status (share_type, status, expires_at) COMMENT '类型状态过期查询',
INDEX idx_file_shares_active_expires (is_active, expires_at, last_accessed_at) COMMENT '活跃分享过期查询',

-- 全文搜索索引
FULLTEXT idx_file_shares_fulltext_search (share_name, share_description) COMMENT '分享内容全文搜索'
```

#### share_access_logs 表索引策略
```sql
-- 访问记录索引
INDEX idx_share_access_logs_share_id (share_id) COMMENT '分享访问记录查询',
INDEX idx_share_access_logs_share_code (share_code) COMMENT '分享码访问统计',
INDEX idx_share_access_logs_visitor_ip (visitor_ip) COMMENT 'IP安全分析',
INDEX idx_share_access_logs_access_type (access_type) COMMENT '操作统计',
INDEX idx_share_access_logs_file_id (file_id) COMMENT '文件访问统计',

-- 时间和状态索引
INDEX idx_share_access_logs_created_at (created_at) COMMENT '访问时间查询',
INDEX idx_share_access_logs_success (success) COMMENT '成功率分析',
INDEX idx_share_access_logs_device_type (device_type) COMMENT '设备统计',

-- 地理和安全索引
INDEX idx_share_access_logs_country_code (visitor_country_code) COMMENT '地理分布统计',
INDEX idx_share_access_logs_is_suspicious (is_suspicious) COMMENT '可疑访问监控',
INDEX idx_share_access_logs_authentication_method (authentication_method) COMMENT '认证方式统计',
INDEX idx_share_access_logs_visitor_fingerprint (visitor_fingerprint) COMMENT '设备指纹识别',

-- 复合分析索引
INDEX idx_share_access_logs_share_time (share_id, created_at DESC) COMMENT '分享访问历史',
INDEX idx_share_access_logs_ip_time (visitor_ip, created_at DESC) COMMENT 'IP行为分析',
INDEX idx_share_access_logs_type_success_time (access_type, success, created_at DESC) COMMENT '操作成功率分析',
INDEX idx_share_access_logs_file_access_stats (file_id, access_type, success, created_at) COMMENT '文件访问统计',
INDEX idx_share_access_logs_security_analysis (visitor_ip, is_suspicious, risk_score, created_at) COMMENT '安全分析',
INDEX idx_share_access_logs_performance_analysis (access_type, response_time, bytes_transferred) COMMENT '性能分析'
```

### 6. 团队协作模块

#### teams 表索引策略
```sql
-- 团队基础索引
INDEX idx_teams_name (name) COMMENT '团队名称搜索',
INDEX idx_teams_slug (slug) COMMENT '团队标识符URL访问',
INDEX idx_teams_owner_id (owner_id) COMMENT '团队所有者查询',
INDEX idx_teams_team_type (team_type) COMMENT '团队类型筛选',
INDEX idx_teams_status (status) COMMENT '团队状态管理',
INDEX idx_teams_is_active (is_active) COMMENT '活跃团队筛选',

-- 业务管理索引
INDEX idx_teams_plan_type (plan_type) COMMENT '订阅计划管理',
INDEX idx_teams_verification_level (verification_level) COMMENT '认证级别',
INDEX idx_teams_created_at (created_at) COMMENT '创建时间排序',
INDEX idx_teams_last_activity_at (last_activity_at) COMMENT '活跃度排序',

-- 统计索引
INDEX idx_teams_storage_used (storage_used) COMMENT '存储使用统计',
INDEX idx_teams_member_count (member_count) COMMENT '成员数量统计',
INDEX idx_teams_deleted_at (deleted_at) COMMENT '软删除管理',
INDEX idx_teams_country_code (country_code) COMMENT '地理分布统计',

-- 复合业务索引
INDEX idx_teams_owner_status (owner_id, status, is_active) COMMENT '所有者状态查询',
INDEX idx_teams_type_plan (team_type, plan_type, is_active) COMMENT '类型计划查询',
INDEX idx_teams_activity_stats (last_activity_at DESC, member_count, storage_used) COMMENT '活跃度统计',
INDEX idx_teams_storage_quota_usage (storage_used, storage_quota, storage_warning_threshold) COMMENT '存储使用情况',

-- 全文搜索索引
FULLTEXT idx_teams_fulltext_search (name, description) COMMENT '团队名称描述全文搜索'
```

#### team_members 表索引策略
```sql
-- 成员关系索引
INDEX idx_team_members_team_id (team_id) COMMENT '团队成员查询',
INDEX idx_team_members_user_id (user_id) COMMENT '用户团队查询',
INDEX idx_team_members_role (role) COMMENT '角色筛选',
INDEX idx_team_members_status (status) COMMENT '成员状态筛选',
INDEX idx_team_members_invited_by (invited_by) COMMENT '邀请记录查询',

-- 时间索引
INDEX idx_team_members_joined_at (joined_at) COMMENT '加入时间排序',
INDEX idx_team_members_last_active_at (last_active_at) COMMENT '活跃度统计',
INDEX idx_team_members_deleted_at (deleted_at) COMMENT '软删除管理',

-- 组织架构索引
INDEX idx_team_members_access_level (access_level) COMMENT '访问级别管理',
INDEX idx_team_members_department (department) COMMENT '部门查询',
INDEX idx_team_members_external_id (external_id) COMMENT '外部系统集成',

-- 复合业务索引
INDEX idx_team_members_team_status_role (team_id, status, role) COMMENT '团队状态角色查询',
INDEX idx_team_members_team_active_members (team_id, status, last_active_at DESC) COMMENT '团队活跃成员',
INDEX idx_team_members_user_teams_active (user_id, status, joined_at DESC) COMMENT '用户活跃团队',
INDEX idx_team_members_invitation_tracking (invited_by, status, invited_at DESC) COMMENT '邀请跟踪',
INDEX idx_team_members_activity_analysis (team_id, activity_score DESC, last_active_at) COMMENT '活跃度分析',
INDEX idx_team_members_storage_usage (team_id, storage_used DESC, storage_quota) COMMENT '存储使用分析',

-- 全文搜索索引
FULLTEXT idx_team_members_fulltext_search (display_name, bio, job_title) COMMENT '成员信息全文搜索'
```

## 🔧 索引优化配置

### 1. MySQL 配置优化

#### 1.1 InnoDB 索引优化
```sql
-- InnoDB 缓冲池大小（推荐系统内存的70-80%）
SET GLOBAL innodb_buffer_pool_size = 2147483648; -- 2GB

-- 索引统计信息更新
SET GLOBAL innodb_stats_persistent = ON;
SET GLOBAL innodb_stats_auto_recalc = ON;

-- 并发读写优化
SET GLOBAL innodb_thread_concurrency = 0;
SET GLOBAL innodb_read_io_threads = 8;
SET GLOBAL innodb_write_io_threads = 8;
```

#### 1.2 全文索引优化
```sql
-- 中文分词优化
SET GLOBAL ft_min_word_len = 1;
SET GLOBAL ngram_token_size = 2;

-- 全文索引缓存大小
SET GLOBAL ft_query_expansion_limit = 20;
SET GLOBAL ft_boolean_syntax = '+ -><()~*:""&|';
```

### 2. 索引维护策略

#### 2.1 定期维护任务
```sql
-- 重建索引统计信息（每周执行）
ANALYZE TABLE users, files, file_versions, upload_tasks;

-- 优化全文索引（每月执行）
OPTIMIZE TABLE files, file_shares, teams;

-- 检查索引碎片（每月执行）
SELECT 
    table_name,
    index_name,
    cardinality,
    pages,
    pages/cardinality as pages_per_key
FROM information_schema.statistics 
WHERE table_schema = 'cloud_storage'
AND cardinality > 0
ORDER BY pages_per_key DESC;
```

#### 2.2 索引监控指标
```sql
-- 查询索引使用情况
SELECT 
    object_schema,
    object_name,
    index_name,
    count_star,
    count_read,
    count_write,
    count_read/count_star as read_ratio
FROM performance_schema.table_io_waits_summary_by_index_usage
WHERE object_schema = 'cloud_storage'
ORDER BY count_star DESC;

-- 查询慢查询中的索引使用
SELECT 
    digest_text,
    count_star,
    avg_timer_wait/1000000000 as avg_time_sec,
    rows_examined_avg,
    rows_sent_avg
FROM performance_schema.events_statements_summary_by_digest
WHERE digest_text LIKE '%SELECT%'
ORDER BY avg_timer_wait DESC
LIMIT 10;
```

### 3. 分区策略

#### 3.1 时间分区表
```sql
-- 日志表按年分区
CREATE TABLE operation_logs (
    -- ... 字段定义 ...
) PARTITION BY RANGE (YEAR(created_at)) (
    PARTITION p2024 VALUES LESS THAN (2025),
    PARTITION p2025 VALUES LESS THAN (2026),
    PARTITION p2026 VALUES LESS THAN (2027),
    PARTITION p_future VALUES LESS THAN MAXVALUE
);

-- 聊天消息表按年分区
CREATE TABLE chat_messages (
    -- ... 字段定义 ...
) PARTITION BY RANGE (YEAR(created_at)) (
    PARTITION p2024 VALUES LESS THAN (2025),
    PARTITION p2025 VALUES LESS THAN (2026),
    PARTITION p2026 VALUES LESS THAN (2027),
    PARTITION p_future VALUES LESS THAN MAXVALUE
);
```

#### 3.2 分区维护
```sql
-- 添加新年度分区
ALTER TABLE operation_logs ADD PARTITION (
    PARTITION p2027 VALUES LESS THAN (2028)
);

-- 删除旧分区数据
ALTER TABLE operation_logs DROP PARTITION p2022;
```

## 📊 性能监控与调优

### 1. 索引性能监控

#### 1.1 关键性能指标
- **索引命中率**: 查询是否使用了预期索引
- **索引选择性**: 索引的区分度和过滤效果
- **查询响应时间**: 各类查询的平均响应时间
- **索引维护成本**: 写入操作对索引的影响

#### 1.2 监控查询语句
```sql
-- 监控索引使用情况
EXPLAIN SELECT * FROM files 
WHERE user_id = 1 AND is_deleted = 0 
ORDER BY created_at DESC;

-- 分析查询成本
EXPLAIN FORMAT=JSON SELECT f.*, fv.version_number
FROM files f
LEFT JOIN file_versions fv ON f.id = fv.file_id AND fv.is_active = 1
WHERE f.user_id = 1 AND f.is_deleted = 0;
```

### 2. 索引优化建议

#### 2.1 查询优化原则
1. **最左前缀原则**: 复合索引按查询频率排列字段顺序
2. **索引覆盖**: 尽量使用覆盖索引减少回表查询
3. **范围查询优化**: 范围字段放在复合索引最后
4. **NULL值处理**: 考虑NULL值对索引效果的影响

#### 2.2 常见优化场景
```sql
-- 优化文件树查询
-- 原查询: WHERE user_id = ? AND parent_id = ? AND is_deleted = 0
-- 优化索引: idx_files_user_parent_deleted (user_id, parent_id, is_deleted)

-- 优化分页查询
-- 原查询: WHERE user_id = ? ORDER BY created_at DESC LIMIT ?, ?
-- 优化索引: idx_files_user_created (user_id, created_at DESC)

-- 优化统计查询
-- 原查询: WHERE status = 'active' AND created_at >= ?
-- 优化索引: idx_table_status_created (status, created_at)
```

### 3. 存储优化策略

#### 3.1 索引大小控制
- 限制VARCHAR字段索引长度
- 使用前缀索引优化长文本字段
- 定期清理无用索引

#### 3.2 数据归档策略
- 历史数据定期归档
- 冷数据分离存储
- 日志数据定期清理

## 🚀 实施建议

### 1. 部署阶段
1. **开发环境**: 完整索引配置，性能测试
2. **测试环境**: 数据量压测，索引效果验证
3. **生产环境**: 渐进式索引创建，监控性能影响

### 2. 监控告警
1. **慢查询监控**: 超过1秒的查询自动告警
2. **索引使用率**: 定期检查未使用的索引
3. **存储空间**: 索引占用空间监控

### 3. 持续优化
1. **定期审查**: 月度索引使用情况分析
2. **业务调整**: 根据业务变化调整索引策略
3. **版本升级**: 跟进MySQL新特性优化索引

---

**配置完成后，请执行以下验证步骤：**
1. 运行 EXPLAIN 分析关键查询
2. 监控索引创建对写入性能的影响
3. 验证全文搜索功能正常工作
4. 检查索引命名规范符合项目标准