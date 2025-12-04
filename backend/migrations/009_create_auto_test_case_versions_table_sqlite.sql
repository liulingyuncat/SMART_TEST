-- Migration 009: Create auto_test_case_versions table for SQLite
-- Description: Add version management for automated test cases
-- Date: 2025-11-21

CREATE TABLE IF NOT EXISTS auto_test_case_versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    version_id TEXT NOT NULL,                       -- 版本ID(格式:{项目名}_{YYYYMMDD_HHMMSS})
    project_id INTEGER NOT NULL,                    -- 项目ID
    project_name TEXT NOT NULL,                     -- 项目名称(冗余存储,用于文件命名)
    role_type TEXT NOT NULL CHECK (role_type IN ('role1', 'role2', 'role3', 'role4')),
    filename TEXT NOT NULL,                         -- Excel文件名
    file_path TEXT NOT NULL,                        -- 文件存储路径
    file_size INTEGER,                              -- 文件大小(字节)
    case_count INTEGER DEFAULT 0,                   -- 该ROLE的用例数量
    remark TEXT DEFAULT '',                         -- 备注(最大200字符)
    created_by INTEGER,                             -- 创建用户ID
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

-- 创建索引优化查询性能
CREATE INDEX IF NOT EXISTS idx_auto_versions_project_version ON auto_test_case_versions(project_id, version_id);
CREATE INDEX IF NOT EXISTS idx_auto_versions_created ON auto_test_case_versions(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_auto_versions_role ON auto_test_case_versions(role_type);
