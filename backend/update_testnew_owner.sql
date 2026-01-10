-- 更新TestNew项目的创建人为dayun（用户ID=3）
UPDATE projects SET owner_id = 3 WHERE id = 1 AND name = 'TestNew';

-- 验证更新结果
SELECT 
    p.id as project_id,
    p.name as project_name,
    p.owner_id,
    u.username,
    u.nickname as owner_name
FROM projects p
LEFT JOIN users u ON u.id = p.owner_id
WHERE p.id = 1;
