-- Migration 024: Add deleted_at column to execution_case_results
-- Description: Support soft delete for execution case results
-- Date: 2025-01-28

ALTER TABLE execution_case_results ADD COLUMN deleted_at DATETIME;
CREATE INDEX idx_ecr_deleted_at ON execution_case_results(deleted_at);
