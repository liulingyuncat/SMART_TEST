-- 010_create_api_test_cases_table.sql
-- Create API test cases table and version management table

BEGIN;

-- Create api_test_cases table
CREATE TABLE api_test_cases (
    id VARCHAR(36) PRIMARY KEY,
    project_id INTEGER NOT NULL,
    case_type VARCHAR(20) NOT NULL DEFAULT 'role1' CHECK (case_type IN ('role1','role2','role3','role4')),
    case_number VARCHAR(50),
    screen VARCHAR(100),
    url TEXT,
    header TEXT,
    method VARCHAR(10) NOT NULL DEFAULT 'GET' CHECK (method IN ('GET','POST','PUT','DELETE','PATCH')),
    body TEXT,
    response TEXT,
    test_result VARCHAR(10) DEFAULT 'NR' CHECK (test_result IN ('OK','NG','NR')),
    remark TEXT,
    display_order INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- Indexes
CREATE INDEX idx_api_cases_project_type ON api_test_cases(project_id, case_type);
CREATE INDEX idx_api_cases_display_order ON api_test_cases(project_id, case_type, display_order);

-- Create api_test_case_versions table
CREATE TABLE api_test_case_versions (
    id VARCHAR(36) PRIMARY KEY,
    project_id INTEGER NOT NULL,
    filename_role1 VARCHAR(255) NOT NULL,
    filename_role2 VARCHAR(255) NOT NULL,
    filename_role3 VARCHAR(255) NOT NULL,
    filename_role4 VARCHAR(255) NOT NULL,
    remark TEXT,
    created_by INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- Indexes
CREATE INDEX idx_api_versions_project ON api_test_case_versions(project_id);
CREATE INDEX idx_api_versions_created_at ON api_test_case_versions(created_at DESC);

COMMIT;
