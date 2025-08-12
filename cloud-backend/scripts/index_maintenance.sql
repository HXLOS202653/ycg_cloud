-- ============================================================================
-- 数据库索引维护脚本
-- 用途: 定期维护和优化数据库索引
-- 执行频率: 建议每周执行一次
-- ============================================================================

-- 设置会话变量
SET SESSION sql_mode = 'STRICT_TRANS_TABLES,NO_ZERO_DATE,NO_ZERO_IN_DATE,ERROR_FOR_DIVISION_BY_ZERO';

-- ============================================================================
-- 1. 索引使用情况分析
-- ============================================================================

SELECT '=== 索引使用情况分析 ===' as 'Analysis Report';

-- 查看索引使用统计
SELECT 
    '索引使用统计' as 'Report Type',
    object_name as 'Table Name',
    index_name as 'Index Name',
    count_read as 'Read Count',
    count_write as 'Write Count',
    count_read + count_write as 'Total Usage',
    CASE 
        WHEN count_read + count_write = 0 THEN '未使用'
        WHEN count_read + count_write < 100 THEN '低使用'
        WHEN count_read + count_write < 1000 THEN '中等使用'
        ELSE '高使用'
    END as 'Usage Level'
FROM performance_schema.table_io_waits_summary_by_index_usage
WHERE object_schema = DATABASE()
    AND index_name IS NOT NULL
    AND index_name != 'PRIMARY'
ORDER BY count_read + count_write DESC;

-- 查找未使用的索引
SELECT 
    '未使用索引检查' as 'Report Type',
    t.table_name as 'Table Name',
    t.index_name as 'Index Name',
    GROUP_CONCAT(t.column_name ORDER BY t.seq_in_index) as 'Columns',
    'DROP INDEX' as 'Recommended Action'
FROM information_schema.statistics t
LEFT JOIN performance_schema.table_io_waits_summary_by_index_usage p 
    ON t.table_schema = p.object_schema 
    AND t.table_name = p.object_name 
    AND t.index_name = p.index_name
WHERE t.table_schema = DATABASE()
    AND t.index_name != 'PRIMARY'
    AND (p.index_name IS NULL OR (p.count_read + p.count_write = 0))
GROUP BY t.table_name, t.index_name
ORDER BY t.table_name, t.index_name;

-- ============================================================================
-- 2. 慢查询分析
-- ============================================================================

SELECT '=== 慢查询分析 ===' as 'Analysis Report';

-- 分析慢查询模式
SELECT 
    '慢查询统计' as 'Report Type',
    LEFT(digest_text, 100) as 'Query Pattern',
    count_star as 'Execution Count',
    ROUND(avg_timer_wait/1000000000, 3) as 'Avg Time (sec)',
    ROUND(max_timer_wait/1000000000, 3) as 'Max Time (sec)',
    ROUND(sum_timer_wait/1000000000, 3) as 'Total Time (sec)',
    ROUND(sum_rows_examined/count_star, 0) as 'Avg Rows Examined',
    ROUND(sum_rows_sent/count_star, 0) as 'Avg Rows Sent',
    CASE 
        WHEN sum_rows_examined/count_star > sum_rows_sent/count_star * 100 THEN '需要索引优化'
        WHEN avg_timer_wait/1000000000 > 2 THEN '查询较慢'
        ELSE '正常'
    END as 'Analysis'
FROM performance_schema.events_statements_summary_by_digest
WHERE digest_text IS NOT NULL
    AND digest_text NOT LIKE '%performance_schema%'
    AND digest_text NOT LIKE '%information_schema%'
    AND avg_timer_wait > 1000000000  -- 大于1秒的查询
ORDER BY avg_timer_wait DESC
LIMIT 20;

-- ============================================================================
-- 3. 表空间和索引大小分析
-- ============================================================================

SELECT '=== 表空间和索引大小分析 ===' as 'Analysis Report';

-- 分析表和索引大小
SELECT 
    '表大小统计' as 'Report Type',
    table_name as 'Table Name',
    ROUND(((data_length + index_length) / 1024 / 1024), 2) as 'Total Size (MB)',
    ROUND((data_length / 1024 / 1024), 2) as 'Data Size (MB)',
    ROUND((index_length / 1024 / 1024), 2) as 'Index Size (MB)',
    ROUND((index_length / (data_length + index_length)) * 100, 2) as 'Index Ratio (%)',
    table_rows as 'Row Count',
    CASE 
        WHEN (index_length / (data_length + index_length)) > 0.5 THEN '索引过多'
        WHEN (index_length / (data_length + index_length)) < 0.1 THEN '索引不足'
        ELSE '正常'
    END as 'Index Health'
FROM information_schema.tables
WHERE table_schema = DATABASE()
    AND table_type = 'BASE TABLE'
ORDER BY (data_length + index_length) DESC;

-- ============================================================================
-- 4. 索引基数分析
-- ============================================================================

SELECT '=== 索引基数分析 ===' as 'Analysis Report';

-- 分析索引选择性
SELECT 
    '索引选择性分析' as 'Report Type',
    table_name as 'Table Name',
    index_name as 'Index Name',
    column_name as 'Column Name',
    cardinality as 'Cardinality',
    CASE 
        WHEN cardinality IS NULL THEN '需要更新统计信息'
        WHEN cardinality = 0 THEN '无区分度'
        WHEN cardinality = 1 THEN '区分度极低'
        WHEN cardinality < 10 THEN '区分度低'
        WHEN cardinality < 100 THEN '区分度中等'
        ELSE '区分度高'
    END as 'Selectivity'
FROM information_schema.statistics
WHERE table_schema = DATABASE()
    AND index_name != 'PRIMARY'
ORDER BY table_name, index_name, seq_in_index;

-- ============================================================================
-- 5. 全文索引状态检查
-- ============================================================================

SELECT '=== 全文索引状态检查 ===' as 'Analysis Report';

-- 检查全文索引配置
SELECT 
    '全文索引配置' as 'Report Type',
    'ft_min_word_len' as 'Parameter',
    @@ft_min_word_len as 'Value',
    '最小词长度' as 'Description';

SELECT 
    '全文索引配置' as 'Report Type',
    'ngram_token_size' as 'Parameter',
    @@ngram_token_size as 'Value',
    'N-gram分词大小' as 'Description';

-- 检查全文索引表
SELECT 
    '全文索引列表' as 'Report Type',
    table_name as 'Table Name',
    index_name as 'Index Name',
    GROUP_CONCAT(column_name ORDER BY seq_in_index) as 'Indexed Columns'
FROM information_schema.statistics
WHERE table_schema = DATABASE()
    AND index_type = 'FULLTEXT'
GROUP BY table_name, index_name
ORDER BY table_name, index_name;

-- ============================================================================
-- 6. 执行索引维护操作
-- ============================================================================

SELECT '=== 开始执行索引维护操作 ===' as 'Maintenance Operations';

-- 更新表统计信息
ANALYZE TABLE users;
SELECT 'ANALYZE TABLE users' as 'Operation', 'Completed' as 'Status';

ANALYZE TABLE files;
SELECT 'ANALYZE TABLE files' as 'Operation', 'Completed' as 'Status';

ANALYZE TABLE file_versions;
SELECT 'ANALYZE TABLE file_versions' as 'Operation', 'Completed' as 'Status';

ANALYZE TABLE upload_tasks;
SELECT 'ANALYZE TABLE upload_tasks' as 'Operation', 'Completed' as 'Status';

ANALYZE TABLE file_chunks;
SELECT 'ANALYZE TABLE file_chunks' as 'Operation', 'Completed' as 'Status';

ANALYZE TABLE file_shares;
SELECT 'ANALYZE TABLE file_shares' as 'Operation', 'Completed' as 'Status';

ANALYZE TABLE teams;
SELECT 'ANALYZE TABLE teams' as 'Operation', 'Completed' as 'Status';

ANALYZE TABLE team_members;
SELECT 'ANALYZE TABLE team_members' as 'Operation', 'Completed' as 'Status';

ANALYZE TABLE notifications;
SELECT 'ANALYZE TABLE notifications' as 'Operation', 'Completed' as 'Status';

-- 优化包含全文索引的表（月度执行）
-- 注意：OPTIMIZE TABLE 会锁表，建议在低峰期执行
-- OPTIMIZE TABLE files;
-- OPTIMIZE TABLE file_shares;
-- OPTIMIZE TABLE teams;

SELECT 'OPTIMIZE TABLE operations' as 'Operation', 'Skipped (uncomment to enable)' as 'Status';

-- ============================================================================
-- 7. 性能配置检查
-- ============================================================================

SELECT '=== 性能配置检查 ===' as 'Configuration Check';

-- 检查重要的MySQL配置参数
SELECT 
    '缓冲池配置' as 'Category',
    'innodb_buffer_pool_size' as 'Parameter',
    @@innodb_buffer_pool_size as 'Current Value',
    'bytes' as 'Unit',
    CASE 
        WHEN @@innodb_buffer_pool_size < 1073741824 THEN '建议增加到至少1GB'
        ELSE '配置合理'
    END as 'Recommendation';

SELECT 
    '索引统计配置' as 'Category',
    'innodb_stats_persistent' as 'Parameter',
    @@innodb_stats_persistent as 'Current Value',
    'boolean' as 'Unit',
    CASE 
        WHEN @@innodb_stats_persistent = 1 THEN '配置正确'
        ELSE '建议开启持久化统计'
    END as 'Recommendation';

SELECT 
    '慢查询配置' as 'Category',
    'slow_query_log' as 'Parameter',
    @@slow_query_log as 'Current Value',
    'boolean' as 'Unit',
    CASE 
        WHEN @@slow_query_log = 1 THEN '配置正确'
        ELSE '建议开启慢查询日志'
    END as 'Recommendation';

SELECT 
    '慢查询配置' as 'Category',
    'long_query_time' as 'Parameter',
    @@long_query_time as 'Current Value',
    'seconds' as 'Unit',
    CASE 
        WHEN @@long_query_time <= 2 THEN '配置合理'
        ELSE '建议设置为2秒以下'
    END as 'Recommendation';

-- ============================================================================
-- 8. 索引碎片检查
-- ============================================================================

SELECT '=== 索引碎片检查 ===' as 'Fragmentation Check';

-- 检查表碎片情况
SELECT 
    '表碎片分析' as 'Report Type',
    table_name as 'Table Name',
    ROUND(data_length/1024/1024, 2) as 'Data Size (MB)',
    ROUND(data_free/1024/1024, 2) as 'Free Space (MB)',
    ROUND((data_free/(data_length+data_free))*100, 2) as 'Fragmentation (%)',
    CASE 
        WHEN data_free = 0 THEN '无碎片'
        WHEN (data_free/(data_length+data_free))*100 < 10 THEN '碎片较少'
        WHEN (data_free/(data_length+data_free))*100 < 25 THEN '中等碎片'
        ELSE '严重碎片，建议优化'
    END as 'Fragmentation Level'
FROM information_schema.tables
WHERE table_schema = DATABASE()
    AND table_type = 'BASE TABLE'
    AND data_length > 0
ORDER BY (data_free/(data_length+data_free)) DESC;

-- ============================================================================
-- 9. 分区表维护（如果使用了分区）
-- ============================================================================

SELECT '=== 分区表维护检查 ===' as 'Partition Maintenance';

-- 检查分区表信息
SELECT 
    '分区表列表' as 'Report Type',
    table_name as 'Table Name',
    partition_name as 'Partition Name',
    partition_method as 'Partition Method',
    partition_expression as 'Partition Expression',
    table_rows as 'Rows',
    ROUND((data_length + index_length)/1024/1024, 2) as 'Size (MB)'
FROM information_schema.partitions
WHERE table_schema = DATABASE()
    AND partition_name IS NOT NULL
ORDER BY table_name, partition_ordinal_position;

-- 检查是否需要添加新分区（针对时间分区表）
SELECT 
    '分区维护建议' as 'Report Type',
    'operation_logs' as 'Table Name',
    CONCAT('如果当前年份是 ', YEAR(NOW()), '，建议检查是否需要添加 p', YEAR(NOW())+1, ' 分区') as 'Recommendation';

SELECT 
    '分区维护建议' as 'Report Type',
    'chat_messages' as 'Table Name',
    CONCAT('如果当前年份是 ', YEAR(NOW()), '，建议检查是否需要添加 p', YEAR(NOW())+1, ' 分区') as 'Recommendation';

-- ============================================================================
-- 10. 维护总结报告
-- ============================================================================

SELECT '=== 维护总结报告 ===' as 'Summary Report';

SELECT 
    '维护完成时间' as 'Item',
    NOW() as 'Value',
    '索引维护操作已完成' as 'Status';

SELECT 
    '下次维护建议' as 'Item',
    DATE_ADD(NOW(), INTERVAL 7 DAY) as 'Next Maintenance Date',
    '建议每周执行一次维护' as 'Frequency';

-- 生成维护建议
SELECT 
    '维护建议' as 'Category',
    CASE 
        WHEN COUNT(*) > 0 THEN CONCAT('发现 ', COUNT(*), ' 个未使用的索引，建议评估后删除')
        ELSE '所有索引都在正常使用'
    END as 'Recommendation'
FROM information_schema.statistics t
LEFT JOIN performance_schema.table_io_waits_summary_by_index_usage p 
    ON t.table_schema = p.object_schema 
    AND t.table_name = p.object_name 
    AND t.index_name = p.index_name
WHERE t.table_schema = DATABASE()
    AND t.index_name != 'PRIMARY'
    AND (p.index_name IS NULL OR (p.count_read + p.count_write = 0));

-- ============================================================================
-- 可选：清理操作（请根据实际情况启用）
-- ============================================================================

-- 清理过期的性能统计数据（可选）
-- CALL sys.ps_truncate_all_tables(FALSE);

-- 重置慢查询统计（可选）
-- SET GLOBAL slow_query_log = 'OFF';
-- SET GLOBAL slow_query_log = 'ON';

SELECT '维护脚本执行完成' as 'Final Status', NOW() as 'Completion Time';

-- ============================================================================
-- 脚本使用说明
-- ============================================================================

/*
使用说明：

1. 执行频率：
   - 索引统计分析：每周执行
   - ANALYZE TABLE：每周执行
   - OPTIMIZE TABLE：每月执行（在低峰期）

2. 注意事项：
   - OPTIMIZE TABLE 会锁表，建议在维护窗口执行
   - 某些配置修改需要 SUPER 权限
   - 生产环境执行前请先在测试环境验证

3. 监控指标：
   - 未使用索引数量
   - 慢查询数量和执行时间
   - 表碎片率
   - 索引选择性

4. 优化建议：
   - 删除未使用的索引
   - 为慢查询创建合适的索引
   - 定期清理碎片
   - 更新统计信息

5. 自动化：
   - 可以将此脚本添加到 cron 作业中定期执行
   - 建议配合监控系统使用
   - 重要操作建议人工确认后执行
*/