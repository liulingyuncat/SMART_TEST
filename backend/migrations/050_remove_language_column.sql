-- 050_remove_language_column.sql
-- 功能: 移除manual_test_cases表中已废弃的language字段
-- 说明: language字段已不再使用，多语言数据现在通过独立的 *_cn/*_jp/*_en 字段存储
-- 日期: 2025-12-10

-- ========== 删除索引 ==========
DROP INDEX IF EXISTS idx_mtc_language;

-- ========== 删除字段 ==========
-- SQLite 不支持直接删除列，需要重建表
-- 以下SQL适用于PostgreSQL/MySQL，SQLite需要通过重建表实现

-- PostgreSQL/MySQL:
-- ALTER TABLE manual_test_cases DROP COLUMN IF EXISTS language;

-- 对于SQLite，如果需要删除列，请执行以下操作:
-- 1. 创建不包含language列的新表
-- 2. 复制数据
-- 3. 删除旧表
-- 4. 重命名新表

-- 注意: 由于SQLite限制，建议保留language列但不再使用
-- 或在下次大版本升级时通过应用层重建表结构
