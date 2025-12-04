-- Migration 006: Create test_case_reviews and test_case_versions tables
-- Description: Add support for case review and version management features
-- Date: 2025-11-13

-- Create test_case_reviews table
CREATE TABLE IF NOT EXISTS test_case_reviews (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL,
    case_type VARCHAR(20) NOT NULL CHECK (case_type IN ('ai', 'overall', 'change')),
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_reviews_project FOREIGN KEY (project_id) 
        REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT uq_reviews_project_type UNIQUE (project_id, case_type)
);

CREATE INDEX idx_reviews_project ON test_case_reviews(project_id);

COMMENT ON TABLE test_case_reviews IS '测试用例评审记录表';
COMMENT ON COLUMN test_case_reviews.project_id IS '项目ID';
COMMENT ON COLUMN test_case_reviews.case_type IS '用例类型: ai/overall/change';
COMMENT ON COLUMN test_case_reviews.content IS 'Markdown格式评审内容';

-- Create test_case_versions table
CREATE TABLE IF NOT EXISTS test_case_versions (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL,
    case_type VARCHAR(20) NOT NULL DEFAULT 'overall' CHECK (case_type IN ('ai', 'overall', 'change')),
    filename VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT,
    created_by INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_versions_project FOREIGN KEY (project_id) 
        REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_versions_user FOREIGN KEY (created_by) 
        REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_versions_project ON test_case_versions(project_id);
CREATE INDEX idx_versions_created ON test_case_versions(created_at DESC);

COMMENT ON TABLE test_case_versions IS '测试用例版本记录表';
COMMENT ON COLUMN test_case_versions.project_id IS '项目ID';
COMMENT ON COLUMN test_case_versions.case_type IS '用例类型: ai/overall/change (当前仅支持overall)';
COMMENT ON COLUMN test_case_versions.filename IS '版本文件名';
COMMENT ON COLUMN test_case_versions.file_path IS '文件存储路径(相对路径)';
COMMENT ON COLUMN test_case_versions.file_size IS '文件大小(字节)';
COMMENT ON COLUMN test_case_versions.created_by IS '创建用户ID';
