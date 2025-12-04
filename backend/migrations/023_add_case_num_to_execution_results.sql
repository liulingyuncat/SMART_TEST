-- Migration 023: Add case_num field to execution_case_results
-- Description: Store user-defined CaseID for display consistency with AIWeb用例库
-- Date: 2025-01-28

-- Add case_num field (user-defined CaseID like "Login223", "forget113")
ALTER TABLE execution_case_results ADD COLUMN case_num VARCHAR(100);
