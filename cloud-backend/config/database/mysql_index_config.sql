-- ============================================================================
-- MySQL 数据库索引策略配置文件
-- 基于: 07数据库设计.md 索引设计
-- 创建时间: 2025-01-12
-- ============================================================================

-- 设置MySQL会话变量
SET SESSION sql_mode = 'STRICT_TRANS_TABLES,NO_ZERO_DATE,NO_ZERO_IN_DATE,ERROR_FOR_DIVISION_BY_ZERO';
SET SESSION innodb_strict_mode = ON;

-- ============================================================================
-- 1. MySQL 全文索引配置优化
-- ============================================================================

-- 全文索引基础配置
SET GLOBAL ft_min_word_len = 1;                    -- 最小词长度，支持中文单字
SET GLOBAL ft_query_expansion_limit = 20;          -- 查询扩展限制
SET GLOBAL ft_boolean_syntax = '+ -><()~*:""&|';   -- 布尔搜索语法

-- ngram 分词器配置（中文支持）
SET GLOBAL ngram_token_size = 2;                   -- 中文分词粒度

-- ============================================================================
-- 2. InnoDB 索引性能优化配置
-- ============================================================================

-- InnoDB 缓冲池配置（根据服务器内存调整）
-- 推荐设置为系统内存的 70-80%
-- SET GLOBAL innodb_buffer_pool_size = 2147483648;   -- 2GB (需要重启MySQL)

-- 索引统计信息配置
SET GLOBAL innodb_stats_persistent = ON;           -- 持久化统计信息
SET GLOBAL innodb_stats_auto_recalc = ON;          -- 自动重计算统计信息
SET GLOBAL innodb_stats_persistent_sample_pages = 20; -- 统计信息采样页数

-- 并发配置
SET GLOBAL innodb_thread_concurrency = 0;          -- 无限制并发（让系统自动调节）
SET GLOBAL innodb_read_io_threads = 8;             -- 读IO线程数
SET GLOBAL innodb_write_io_threads = 8;            -- 写IO线程数

-- 锁等待配置
SET GLOBAL innodb_lock_wait_timeout = 50;          -- 锁等待超时（秒）
SET GLOBAL innodb_deadlock_detect = ON;            -- 死锁检测

-- ============================================================================
-- 3. 查询缓存配置（MySQL 5.7及以下版本）
-- ============================================================================

-- 注意：MySQL 8.0 已移除查询缓存，以下配置仅适用于旧版本
-- SET GLOBAL query_cache_type = ON;
-- SET GLOBAL query_cache_size = 268435456;          -- 256MB
-- SET GLOBAL query_cache_limit = 1048576;           -- 1MB

-- ============================================================================
-- 4. 性能监控配置
-- ============================================================================

-- 启用性能模式
SET GLOBAL performance_schema = ON;

-- 慢查询日志配置
SET GLOBAL slow_query_log = ON;
SET GLOBAL slow_query_log_file = '/var/log/mysql/slow-query.log';
SET GLOBAL long_query_time = 1.0;                  -- 慢查询阈值：1秒
SET GLOBAL log_queries_not_using_indexes = ON;     -- 记录未使用索引的查询

-- ============================================================================
-- 5. 索引优化相关配置
-- ============================================================================

-- 优化器配置
SET GLOBAL optimizer_switch = 'index_merge=on,index_merge_union=on,index_merge_sort_union=on,index_merge_intersection=on';

-- 排序缓冲区大小
SET GLOBAL sort_buffer_size = 2097152;             -- 2MB
SET GLOBAL read_buffer_size = 131072;              -- 128KB
SET GLOBAL read_rnd_buffer_size = 262144;          -- 256KB

-- 连接相关配置
SET GLOBAL max_connections = 1000;                 -- 最大连接数
SET GLOBAL max_connect_errors = 100000;            -- 最大连接错误数

-- ============================================================================
-- 6. 临时表和内存表配置
-- ============================================================================

SET GLOBAL tmp_table_size = 134217728;             -- 128MB
SET GLOBAL max_heap_table_size = 134217728;        -- 128MB

-- ============================================================================
-- 7. 二进制日志配置（主从复制）
-- ============================================================================

-- 如果需要主从复制，配置以下参数
-- SET GLOBAL log_bin = ON;
-- SET GLOBAL binlog_format = 'ROW';
-- SET GLOBAL sync_binlog = 1;
-- SET GLOBAL innodb_flush_log_at_trx_commit = 1;

-- ============================================================================
-- 8. 表维护相关配置
-- ============================================================================

-- 自动统计信息更新
SET GLOBAL innodb_stats_on_metadata = OFF;         -- 关闭元数据统计信息更新

-- 自适应哈希索引
SET GLOBAL innodb_adaptive_hash_index = ON;        -- 启用自适应哈希索引

-- ============================================================================
-- 9. 数据库连接超时配置
-- ============================================================================

SET GLOBAL wait_timeout = 28800;                   -- 8小时
SET GLOBAL interactive_timeout = 28800;            -- 8小时
SET GLOBAL connect_timeout = 10;                   -- 10秒

-- ============================================================================
-- 10. 字符集和排序规则配置
-- ============================================================================

-- 确保使用 UTF8MB4 字符集支持完整的 Unicode
SET GLOBAL character_set_server = 'utf8mb4';
SET GLOBAL collation_server = 'utf8mb4_unicode_ci';

-- ============================================================================
-- 验证配置
-- ============================================================================

-- 显示当前重要配置
SELECT 
    '=== 重要配置检查 ===' as 'Configuration Check';

SELECT 
    'InnoDB Buffer Pool Size' as 'Parameter',
    @@innodb_buffer_pool_size as 'Value',
    'bytes' as 'Unit';

SELECT 
    'Full-Text Min Word Length' as 'Parameter',
    @@ft_min_word_len as 'Value',
    'characters' as 'Unit';

SELECT 
    'NGram Token Size' as 'Parameter',
    @@ngram_token_size as 'Value',
    'characters' as 'Unit';

SELECT 
    'Slow Query Time' as 'Parameter',
    @@long_query_time as 'Value',
    'seconds' as 'Unit';

SELECT 
    'Max Connections' as 'Parameter',
    @@max_connections as 'Value',
    'connections' as 'Unit';

-- ============================================================================
-- 常用监控查询
-- ============================================================================

-- 查看索引使用情况的查询模板
/*
-- 索引使用统计
SELECT 
    object_schema as 'Database',
    object_name as 'Table',
    index_name as 'Index',
    count_read as 'Read Count',
    count_write as 'Write Count',
    count_read + count_write as 'Total Usage'
FROM performance_schema.table_io_waits_summary_by_index_usage
WHERE object_schema = 'your_database_name'
    AND count_star > 0
ORDER BY count_star DESC;

-- 慢查询统计
SELECT 
    digest_text as 'Query Pattern',
    count_star as 'Execution Count',
    ROUND(avg_timer_wait/1000000000, 3) as 'Avg Time (sec)',
    ROUND(max_timer_wait/1000000000, 3) as 'Max Time (sec)',
    ROUND(sum_timer_wait/1000000000, 3) as 'Total Time (sec)'
FROM performance_schema.events_statements_summary_by_digest
WHERE digest_text IS NOT NULL
    AND avg_timer_wait > 1000000000  -- 大于1秒的查询
ORDER BY avg_timer_wait DESC
LIMIT 10;

-- 未使用的索引检查
SELECT 
    t.table_schema as 'Database',
    t.table_name as 'Table',
    t.index_name as 'Unused Index',
    t.column_name as 'Column'
FROM information_schema.statistics t
LEFT JOIN performance_schema.table_io_waits_summary_by_index_usage p 
    ON t.table_schema = p.object_schema 
    AND t.table_name = p.object_name 
    AND t.index_name = p.index_name
WHERE t.table_schema = 'your_database_name'
    AND t.index_name != 'PRIMARY'
    AND p.index_name IS NULL
ORDER BY t.table_name, t.index_name;
*/

-- ============================================================================
-- 配置完成提示
-- ============================================================================

SELECT 
    '数据库索引策略配置完成' as 'Status',
    '请检查上述配置参数是否符合预期' as 'Next Step',
    '建议执行索引使用情况监控查询进行验证' as 'Recommendation';

-- 提示：某些配置需要重启MySQL服务才能生效
-- 包括：innodb_buffer_pool_size 等
SELECT 
    '注意：部分配置需要重启MySQL服务' as 'Important Notice',
    'innodb_buffer_pool_size 等参数需要重启生效' as 'Details';