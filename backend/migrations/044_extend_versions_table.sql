-- 扩展 versions 表
-- 添加 item_type 字段用于区分版本类型（requirement-batch, viewpoint-batch等）
-- 添加 file_list 字段用于存储ZIP内包含的MD文件列表（JSON数组格式）

ALTER TABLE versions ADD COLUMN item_type VARCHAR(50);
ALTER TABLE versions ADD COLUMN file_list TEXT;

-- 创建 item_type 索引：优化按类型查询版本列表
CREATE INDEX IF NOT EXISTS idx_versions_item_type 
ON versions(item_type);
