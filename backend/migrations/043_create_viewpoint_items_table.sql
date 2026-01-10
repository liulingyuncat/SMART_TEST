-- 创建 viewpoint_items 表
-- 用于存储动态创建的AI观点列表
CREATE TABLE IF NOT EXISTS viewpoint_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  project_id INTEGER NOT NULL,
  name VARCHAR(100) NOT NULL,
  content TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME,
  
  -- 外键约束
  FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- 创建复合唯一索引：确保同一项目内观点名称唯一
CREATE UNIQUE INDEX IF NOT EXISTS idx_viewpoint_items_project_name 
ON viewpoint_items(project_id, name) WHERE deleted_at IS NULL;

-- 创建 project_id 索引：优化查询性能
CREATE INDEX IF NOT EXISTS idx_viewpoint_items_project_id 
ON viewpoint_items(project_id);

-- 创建 deleted_at 索引：优化软删除查询
CREATE INDEX IF NOT EXISTS idx_viewpoint_items_deleted_at 
ON viewpoint_items(deleted_at);
