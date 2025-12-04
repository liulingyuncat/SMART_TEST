-- 检查test_case_reviews表结构
SELECT 
    column_name, 
    data_type, 
    character_maximum_length, 
    is_nullable,
    column_default
FROM information_schema.columns 
WHERE table_name = 'test_case_reviews'
ORDER BY ordinal_position;

-- 检查约束
SELECT 
    tc.constraint_name, 
    tc.constraint_type,
    kcu.column_name
FROM information_schema.table_constraints tc
JOIN information_schema.key_column_usage kcu 
    ON tc.constraint_name = kcu.constraint_name
WHERE tc.table_name = 'test_case_reviews'
ORDER BY tc.constraint_type, tc.constraint_name;

-- 检查索引
SELECT 
    indexname,
    indexdef
FROM pg_indexes
WHERE tablename = 'test_case_reviews';

-- 查看当前数据
SELECT id, project_id, case_type, LENGTH(content) as content_length, created_at, updated_at
FROM test_case_reviews
LIMIT 10;
