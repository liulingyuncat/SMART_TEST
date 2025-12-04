-- 005_add_multilingual_fields.sql
-- 功能: 为manual_test_cases表添加18个多语言字段,支持整体/变更用例的多语言管理
-- 注意: 保留现有的6个单语言字段(major_function等)用于AI用例

BEGIN;

-- ========== 新增18个多语言字段(供整体/变更用例使用) ==========

-- 大功能分类 - 三语言
ALTER TABLE manual_test_cases ADD COLUMN major_function_cn VARCHAR(100);
ALTER TABLE manual_test_cases ADD COLUMN major_function_jp VARCHAR(100);
ALTER TABLE manual_test_cases ADD COLUMN major_function_en VARCHAR(100);

-- 中功能分类 - 三语言
ALTER TABLE manual_test_cases ADD COLUMN middle_function_cn VARCHAR(100);
ALTER TABLE manual_test_cases ADD COLUMN middle_function_jp VARCHAR(100);
ALTER TABLE manual_test_cases ADD COLUMN middle_function_en VARCHAR(100);

-- 小功能分类 - 三语言
ALTER TABLE manual_test_cases ADD COLUMN minor_function_cn VARCHAR(100);
ALTER TABLE manual_test_cases ADD COLUMN minor_function_jp VARCHAR(100);
ALTER TABLE manual_test_cases ADD COLUMN minor_function_en VARCHAR(100);

-- 前置条件 - 三语言
ALTER TABLE manual_test_cases ADD COLUMN precondition_cn TEXT;
ALTER TABLE manual_test_cases ADD COLUMN precondition_jp TEXT;
ALTER TABLE manual_test_cases ADD COLUMN precondition_en TEXT;

-- 测试步骤 - 三语言
ALTER TABLE manual_test_cases ADD COLUMN test_steps_cn TEXT;
ALTER TABLE manual_test_cases ADD COLUMN test_steps_jp TEXT;
ALTER TABLE manual_test_cases ADD COLUMN test_steps_en TEXT;

-- 期待值 - 三语言
ALTER TABLE manual_test_cases ADD COLUMN expected_result_cn TEXT;
ALTER TABLE manual_test_cases ADD COLUMN expected_result_jp TEXT;
ALTER TABLE manual_test_cases ADD COLUMN expected_result_en TEXT;

-- ========== 创建索引(用于筛选和排序) ==========

CREATE INDEX idx_mtc_major_func_cn ON manual_test_cases(major_function_cn);
CREATE INDEX idx_mtc_major_func_jp ON manual_test_cases(major_function_jp);
CREATE INDEX idx_mtc_major_func_en ON manual_test_cases(major_function_en);

-- ========== 数据迁移:将现有整体/变更用例的单语言字段复制到CN字段 ==========

UPDATE manual_test_cases SET 
    major_function_cn = major_function,
    middle_function_cn = middle_function,
    minor_function_cn = minor_function,
    precondition_cn = precondition,
    test_steps_cn = test_steps,
    expected_result_cn = expected_result
WHERE case_type IN ('overall', 'change') AND major_function IS NOT NULL;

-- ========== 重要说明 ==========
-- ❌ 不删除单语言字段(major_function, middle_function, minor_function, precondition, test_steps, expected_result)
--    原因: AI用例需要使用这些字段
--
-- ❌ 不删除Language字段
--    原因: 保留用于向后兼容,避免数据迁移风险
--
-- ✅ 字段使用场景:
--    - AI用例(ai): 使用单语言字段(major_function等)
--    - 整体用例(overall): 使用多语言字段(major_function_cn/jp/en等)
--    - 变更用例(change): 使用多语言字段(major_function_cn/jp/en等)

COMMIT;
