-- 010_down_api_test_cases_table.sql
-- 功能: 回滚010_create_api_test_cases_table.sql迁移
-- 用途: 删除api_test_cases表和api_test_case_versions表

BEGIN;

-- 删除版本管理表(先删除,因为有外键依赖)
DROP TABLE IF EXISTS api_test_case_versions;

-- 删除接口测试用例表
DROP TABLE IF EXISTS api_test_cases;

COMMIT;
