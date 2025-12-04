-- 为projects表的name字段添加唯一索引
-- 确保项目名称全局唯一,防止并发创建重名项目
CREATE UNIQUE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
