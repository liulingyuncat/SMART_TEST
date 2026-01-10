-- 创建用例集表
-- 用于独立维护用例集信息，即使用例集内没有用例也保留用例集
-- 创建日期: 2025-12-10

CREATE TABLE IF NOT EXISTS case_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    case_type VARCHAR(20) NOT NULL DEFAULT 'overall',
    group_name VARCHAR(100) NOT NULL,
    description TEXT,
    display_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    UNIQUE(project_id, case_type, group_name)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_cg_project_type ON case_groups(project_id, case_type);
CREATE INDEX IF NOT EXISTS idx_cg_deleted_at ON case_groups(deleted_at);
CREATE INDEX IF NOT EXISTS idx_cg_display_order ON case_groups(display_order);

-- 表说明:
-- case_groups: 用例集表，独立维护用例集信息
-- group_name: 用例集名称，在同一项目和类型下唯一
-- case_type: 用例类型(overall/change/acceptance)
-- display_order: 显示顺序
