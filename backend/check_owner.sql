-- 检查项目的owner_id和对应的用户信息
SELECT 
    p.id as project_id,
    p.name as project_name,
    p.owner_id,
    u.id as user_id,
    u.username,
    u.nickname,
    CASE 
        WHEN u.id IS NULL THEN '用户不存在'
        WHEN u.nickname IS NULL OR u.nickname = '' THEN '昵称为空'
        ELSE 'OK'
    END as status
FROM projects p
LEFT JOIN users u ON u.id = p.owner_id
WHERE p.deleted_at IS NULL;

-- 查看所有用户的信息
SELECT id, username, nickname, role FROM users WHERE deleted_at IS NULL;
