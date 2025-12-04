-- Migration 020: Create test_execution_tasks table
-- Description: Add support for test execution task management
-- Date: 2025-11-27

-- Create test_execution_tasks table
CREATE TABLE IF NOT EXISTS test_execution_tasks (
    task_uuid VARCHAR(36) PRIMARY KEY,
    project_id INTEGER NOT NULL,
    task_name VARCHAR(50) NOT NULL,
    execution_type VARCHAR(20) NOT NULL CHECK (execution_type IN ('manual', 'automation', 'api')),
    task_status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (task_status IN ('pending', 'in_progress', 'completed')),
    start_date DATE,
    end_date DATE,
    test_version VARCHAR(50),
    test_env VARCHAR(100),
    test_date DATE,
    executor VARCHAR(50),
    task_description TEXT,
    created_by INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    CONSTRAINT fk_tet_project FOREIGN KEY (project_id) 
        REFERENCES projects(id) ON DELETE CASCADE,
    CONSTRAINT fk_tet_creator FOREIGN KEY (created_by) 
        REFERENCES users(id)
);

-- Create indexes
CREATE INDEX idx_tet_project ON test_execution_tasks(project_id);
CREATE INDEX idx_tet_status ON test_execution_tasks(task_status);
CREATE INDEX idx_tet_creator ON test_execution_tasks(created_by);

-- Create unique index for case-insensitive task name within project (excluding soft-deleted)
-- Note: SQLite partial unique index syntax is compatible with PostgreSQL
-- This partial unique index only applies when deleted_at IS NULL
CREATE UNIQUE INDEX idx_tet_project_name ON test_execution_tasks(project_id, LOWER(task_name)) WHERE deleted_at IS NULL;
