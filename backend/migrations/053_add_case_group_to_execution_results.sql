-- Migration 053: Add case_group fields to execution_case_results table
-- Description: Add case_group_id and case_group_name fields to execution_case_results
--              to store which case group each execution result belongs to
-- Date: 2025-12-26

-- Add case_group_id column (matches model: `gorm:"type:int;default:0;index:idx_ecr_group_id"`)
ALTER TABLE execution_case_results ADD COLUMN case_group_id INTEGER DEFAULT 0;

-- Add case_group_name column (matches model: `gorm:"type:varchar(100)"`)
ALTER TABLE execution_case_results ADD COLUMN case_group_name VARCHAR(100);

-- Create index for case_group_id for efficient filtering
CREATE INDEX idx_ecr_group_id ON execution_case_results(case_group_id);

-- Verification query (optional, run manually to verify):
-- SELECT sql FROM sqlite_master WHERE type='table' AND name='execution_case_results';
-- PRAGMA table_info(execution_case_results);
