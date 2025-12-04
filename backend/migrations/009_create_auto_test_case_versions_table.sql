-- Migration 009: Create auto_test_case_versions table
-- Description: Add version management for automated test cases
-- Date: 2025-11-21

CREATE TABLE IF NOT EXISTS auto_test_case_versions (
    id SERIAL PRIMARY KEY,
    version_id INTEGER NOT NULL,                    -- 版本ID(同批次4个ROLE共享)
    project_id INTEGER NOT NULL,                    -- 项目ID
    project_name VARCHAR(100) NOT NULL,             -- 项目名称(冗余存储,用于文件命名)
    role_type VARCHAR(10) NOT NULL CHECK (role_type IN ('role1', 'role2', 'role3', 'role4')),
    filename VARCHAR(255) NOT NULL,                 -- Excel文件名
    file_path VARCHAR(500) NOT NULL,                -- 文件存储路径
    file_size BIGINT,                               -- 文件大小(字节)
    case_count INTEGER DEFAULT 0,                   -- 该ROLE的用例数量
    remark VARCHAR(200),                            -- 备注
    created_by INTEGER,                             -- 创建用户ID
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_auto_versions_project FOREIGN KEY (project_id) 
        REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_auto_versions_user FOREIGN KEY (created_by) 
        REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_auto_versions_version_id ON auto_test_case_versions(version_id);
CREATE INDEX idx_auto_versions_project ON auto_test_case_versions(project_id);
CREATE INDEX idx_auto_versions_created ON auto_test_case_versions(created_at DESC);

COMMENT ON TABLE auto_test_case_versions IS '自动化测试用例版本记录表';
COMMENT ON COLUMN auto_test_case_versions.version_id IS '版本ID(同批次4个ROLE共享此ID)';
COMMENT ON COLUMN auto_test_case_versions.role_type IS 'ROLE类型: role1/role2/role3/role4';
COMMENT ON COLUMN auto_test_case_versions.filename IS 'Excel文件名(格式:{项目名}_Autotestcase_ROLE{X}_{timestamp}.xlsx)';
