-- 检查缺陷中Subject和Phase为空的记录
SELECT 
    defect_id,
    title,
    subject,
    phase,
    created_at
FROM defects
WHERE subject = '' OR subject IS NULL OR phase = '' OR phase IS NULL
ORDER BY created_at DESC;

-- 统计空Subject/Phase的数量
SELECT 
    COUNT(*) as total_defects,
    SUM(CASE WHEN subject = '' OR subject IS NULL THEN 1 ELSE 0 END) as empty_subject,
    SUM(CASE WHEN phase = '' OR phase IS NULL THEN 1 ELSE 0 END) as empty_phase
FROM defects;

-- 查看最近创建的缺陷（用于验证修复是否生效）
SELECT 
    defect_id,
    title,
    subject,
    phase,
    created_at
FROM defects
ORDER BY created_at DESC
LIMIT 10;
