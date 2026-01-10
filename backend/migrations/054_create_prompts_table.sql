-- 创建prompts表
CREATE TABLE IF NOT EXISTS prompts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    version TEXT NOT NULL DEFAULT '1.0',
    content TEXT NOT NULL,
    arguments TEXT,
    scope TEXT NOT NULL,
    user_id INTEGER,
    created_by INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- 创建唯一索引
CREATE UNIQUE INDEX IF NOT EXISTS uq_prompt_name ON prompts(project_id, name, scope, COALESCE(user_id, 0)) WHERE deleted_at IS NULL;

-- 创建普通索引用于查询优化  
CREATE INDEX IF NOT EXISTS idx_prompts_project_scope ON prompts(project_id, scope) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_prompts_user ON prompts(user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_prompts_name ON prompts(name) WHERE deleted_at IS NULL;
