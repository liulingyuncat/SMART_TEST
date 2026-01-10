-- 迁移缺陷ID格式从 DEF-XXXXXX 到 XXXXXX
-- 此迁移将现有的缺陷ID格式从 DEF-XXXXXX 改为纯数字格式 XXXXXX

-- 首先备份现有数据（可选）
-- CREATE TABLE defects_backup AS SELECT * FROM defects;

-- 更新缺陷ID格式
UPDATE defects
SET defect_id = SUBSTR(defect_id, 5)
WHERE defect_id LIKE 'DEF-%';

-- 验证更新结果
SELECT COUNT(*) as total_defects, COUNT(CASE WHEN defect_id NOT LIKE 'DEF-%' THEN 1 END) as updated_defects
FROM defects;