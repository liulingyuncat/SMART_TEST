-- Migration 022: Add case snapshot fields to execution_case_results
-- Description: Store case content snapshot for display without re-fetching from original table
-- Date: 2025-01-28

-- Add display_id field
ALTER TABLE execution_case_results ADD COLUMN display_id INTEGER DEFAULT 0;

-- Add case content snapshot fields - Chinese
ALTER TABLE execution_case_results ADD COLUMN screen_cn VARCHAR(500);
ALTER TABLE execution_case_results ADD COLUMN function_cn VARCHAR(500);
ALTER TABLE execution_case_results ADD COLUMN precondition_cn TEXT;
ALTER TABLE execution_case_results ADD COLUMN test_steps_cn TEXT;
ALTER TABLE execution_case_results ADD COLUMN expected_result_cn TEXT;

-- Add case content snapshot fields - Japanese
ALTER TABLE execution_case_results ADD COLUMN screen_jp VARCHAR(500);
ALTER TABLE execution_case_results ADD COLUMN function_jp VARCHAR(500);
ALTER TABLE execution_case_results ADD COLUMN precondition_jp TEXT;
ALTER TABLE execution_case_results ADD COLUMN test_steps_jp TEXT;
ALTER TABLE execution_case_results ADD COLUMN expected_result_jp TEXT;

-- Add case content snapshot fields - English
ALTER TABLE execution_case_results ADD COLUMN screen_en VARCHAR(500);
ALTER TABLE execution_case_results ADD COLUMN function_en VARCHAR(500);
ALTER TABLE execution_case_results ADD COLUMN precondition_en TEXT;
ALTER TABLE execution_case_results ADD COLUMN test_steps_en TEXT;
ALTER TABLE execution_case_results ADD COLUMN expected_result_en TEXT;

-- Create index on display_id for sorting
CREATE INDEX idx_ecr_display_id ON execution_case_results(display_id);
