-- Migration: 003_add_project_edit_delete_indexes.sql
-- Purpose: Create indexes for project-related tables to optimize cascade delete performance
-- Date: 2025-11-04

-- Index for project_members table
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_project_members_project_id 
ON project_members(project_id);

-- Index for manual_test_cases table
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_manual_test_cases_project_id 
ON manual_test_cases(project_id);

-- Index for auto_test_cases table
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_auto_test_cases_project_id 
ON auto_test_cases(project_id);

-- Index for api_test_cases table
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_api_test_cases_project_id 
ON api_test_cases(project_id);
