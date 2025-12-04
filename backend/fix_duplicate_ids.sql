-- 修复overall用例的ID，使其连续递增
-- 使用ROW_NUMBER()为每条记录分配唯一ID

BEGIN TRANSACTION;

-- 创建临时表存储新的ID映射
CREATE TEMPORARY TABLE temp_id_mapping AS
SELECT 
    case_id,
    ROW_NUMBER() OVER (ORDER BY id, created_at) as new_id
FROM manual_test_cases
WHERE case_type = 'overall';

-- 更新ID
UPDATE manual_test_cases
SET id = (
    SELECT new_id 
    FROM temp_id_mapping 
    WHERE temp_id_mapping.case_id = manual_test_cases.case_id
)
WHERE case_type = 'overall';

-- 同样处理change用例
DROP TABLE IF EXISTS temp_id_mapping;

CREATE TEMPORARY TABLE temp_id_mapping AS
SELECT 
    case_id,
    ROW_NUMBER() OVER (ORDER BY id, created_at) as new_id
FROM manual_test_cases
WHERE case_type = 'change';

UPDATE manual_test_cases
SET id = (
    SELECT new_id 
    FROM temp_id_mapping 
    WHERE temp_id_mapping.case_id = manual_test_cases.case_id
)
WHERE case_type = 'change';

-- AI用例
DROP TABLE IF EXISTS temp_id_mapping;

CREATE TEMPORARY TABLE temp_id_mapping AS
SELECT 
    case_id,
    ROW_NUMBER() OVER (ORDER BY id, created_at) as new_id
FROM manual_test_cases
WHERE case_type = 'ai';

UPDATE manual_test_cases
SET id = (
    SELECT new_id 
    FROM temp_id_mapping 
    WHERE temp_id_mapping.case_id = manual_test_cases.case_id
)
WHERE case_type = 'ai';

DROP TABLE IF EXISTS temp_id_mapping;

COMMIT;

-- 验证结果
SELECT case_type, id, COUNT(*) as count 
FROM manual_test_cases 
GROUP BY case_type, id 
HAVING COUNT(*) > 1;
