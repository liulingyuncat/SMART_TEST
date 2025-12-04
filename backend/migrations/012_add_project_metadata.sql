-- Migration: Add project metadata fields (status, owner_id)
-- Created: 2025-11-25
-- Description: Add status and owner_id fields to projects table for project management

-- Step 1: Add status field with default value 'pending'
-- Project status: pending(待开始), in-progress(进行中), completed(已完结)
ALTER TABLE projects 
ADD COLUMN status VARCHAR(20) DEFAULT 'pending' NOT NULL;

-- Step 2: Add owner_id field (nullable foreign key to users table)
-- Project owner user ID, references users.id
ALTER TABLE projects 
ADD COLUMN owner_id INT NULL;

-- Step 3: Create indexes for better query performance
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_owner_id ON projects(owner_id);

-- Step 4: Add foreign key constraint with ON DELETE SET NULL
ALTER TABLE projects 
ADD CONSTRAINT fk_projects_owner 
FOREIGN KEY (owner_id) REFERENCES users(id) 
ON DELETE SET NULL 
ON UPDATE CASCADE;

-- Step 5: Update existing data - set status to 'in-progress' for historical projects
UPDATE projects 
SET status = 'in-progress' 
WHERE status IS NULL OR status = '';

-- Rollback script (for reference, execute manually if needed):
-- ALTER TABLE projects DROP FOREIGN KEY fk_projects_owner;
-- DROP INDEX idx_projects_status ON projects;
-- DROP INDEX idx_projects_owner_id ON projects;
-- ALTER TABLE projects DROP COLUMN owner_id;
-- ALTER TABLE projects DROP COLUMN status;
