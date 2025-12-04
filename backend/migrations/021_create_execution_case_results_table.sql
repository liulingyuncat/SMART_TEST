-- Migration 021: Create execution_case_results table
-- Description: Store test execution results for each case in a task
-- Date: 2025-11-28

-- Create execution_case_results table
CREATE TABLE IF NOT EXISTS execution_case_results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_uuid VARCHAR(36) NOT NULL,
    case_id VARCHAR(36) NOT NULL,
    case_type VARCHAR(20) NOT NULL CHECK (case_type IN ('overall', 'acceptance', 'change', 'ai', 'role1', 'role2', 'role3', 'role4', 'api')),
    test_result VARCHAR(10) NOT NULL DEFAULT 'NR' CHECK (test_result IN ('NR', 'OK', 'NG', 'Block')),
    bug_id VARCHAR(50),
    remark TEXT,
    updated_by INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_ecr_task FOREIGN KEY (task_uuid) 
        REFERENCES test_execution_tasks(task_uuid) ON DELETE CASCADE,
    CONSTRAINT fk_ecr_updater FOREIGN KEY (updated_by) 
        REFERENCES users(id)
);

-- Create unique index on task_uuid + case_id (one result per case per task)
CREATE UNIQUE INDEX idx_task_case ON execution_case_results(task_uuid, case_id);

-- Create index on task_uuid for efficient querying by task
CREATE INDEX idx_ecr_task_uuid ON execution_case_results(task_uuid);

-- Create index on case_id for lookup by case
CREATE INDEX idx_ecr_case_id ON execution_case_results(case_id);

-- Create index on test_result for statistics queries
CREATE INDEX idx_ecr_test_result ON execution_case_results(test_result);
