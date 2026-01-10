-- Migration 052: Add case_group fields to test_execution_tasks table
-- Description: Add case_group_id and case_group_name fields to store the associated case group for execution tasks
-- Date: 2025-12-25

-- Add case_group_id column
ALTER TABLE test_execution_tasks ADD COLUMN IF NOT EXISTS case_group_id INTEGER DEFAULT 0;

-- Add case_group_name column
ALTER TABLE test_execution_tasks ADD COLUMN IF NOT EXISTS case_group_name VARCHAR(100);

-- Create index for case_group_id
CREATE INDEX IF NOT EXISTS idx_tet_case_group ON test_execution_tasks(case_group_id);

-- Note: For SQLite, use the following syntax instead:
-- ALTER TABLE test_execution_tasks ADD COLUMN case_group_id INTEGER DEFAULT 0;
-- ALTER TABLE test_execution_tasks ADD COLUMN case_group_name VARCHAR(100);
-- CREATE INDEX idx_tet_case_group ON test_execution_tasks(case_group_id);
