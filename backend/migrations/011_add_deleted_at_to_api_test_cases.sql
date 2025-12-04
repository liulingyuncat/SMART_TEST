-- 011_add_deleted_at_to_api_test_cases.sql
-- 功能: 为api_test_cases表添加deleted_at字段(GORM软删除支持)

BEGIN;

-- 添加deleted_at列到api_test_cases表
ALTER TABLE api_test_cases ADD COLUMN deleted_at TIMESTAMP;

-- 添加索引(可选,优化软删除查询性能)
CREATE INDEX idx_api_cases_deleted_at ON api_test_cases(deleted_at);

COMMIT;

-- ========== 回滚脚本 ==========
-- 执行回滚: 
-- BEGIN;
-- DROP INDEX IF EXISTS idx_api_cases_deleted_at;
-- 注意: SQLite不支持DROP COLUMN,需要重建表或接受保留该列
-- COMMIT;
