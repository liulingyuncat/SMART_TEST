-- Migration: Create case_review_items table
-- Purpose: 支持多文档审阅管理,替代原test_case_reviews单文档模式
-- Date: 2025-12-09

CREATE TABLE IF NOT EXISTS case_review_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    content TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    -- 外键约束:级联删除
    CONSTRAINT fk_review_items_project FOREIGN KEY (project_id) 
        REFERENCES projects(id) ON DELETE CASCADE,
    
    -- 唯一约束:同项目内审阅名称唯一
    CONSTRAINT uq_review_items_name UNIQUE (project_id, name)
);

-- 创建索引:加速按项目查询
CREATE INDEX idx_review_items_project ON case_review_items(project_id);
