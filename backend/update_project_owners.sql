-- 为所有owner_id为空的项目设置owner_id为admin账号(id=1)
-- 假设admin账号的昵称为'admin'或'Admin'

-- 更新所有owner_id为NULL的项目,设置为admin账号
UPDATE projects 
SET owner_id = 1 
WHERE owner_id IS NULL;

-- 验证更新结果
SELECT id, name, owner_id, created_at 
FROM projects 
ORDER BY created_at DESC;
