-- 回滚prompts表
DROP INDEX IF EXISTS idx_prompts_name;
DROP INDEX IF EXISTS idx_prompts_user;
DROP INDEX IF EXISTS idx_prompts_project_scope;
DROP TABLE IF EXISTS prompts;
