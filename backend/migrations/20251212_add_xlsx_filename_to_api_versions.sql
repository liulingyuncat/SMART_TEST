-- 为api_test_case_versions表添加xlsx_filename字段
-- 迁移日期: 2025-12-12

-- 添加新字段（SQLite只支持ADD COLUMN）
ALTER TABLE api_test_case_versions ADD COLUMN xlsx_filename VARCHAR(255) DEFAULT '';

-- 注意：SQLite不支持MODIFY COLUMN，旧字段保持不变
-- filename_role1/2/3/4字段保留用于向后兼容
