-- 修复Subject字段中的乱码数据
-- 将乱码的"ģ��20" ~ "ģ��34" 修复为 "模块20" ~ "模块34"

-- 查看当前乱码数据
SELECT defect_id, title, subject, hex(subject) as subject_hex 
FROM defects 
WHERE subject LIKE '%ģ%' OR subject LIKE '%��%'
ORDER BY defect_id;

-- 修复DEF-000006 到 DEF-000020的subject
UPDATE defects SET subject = '模块20' WHERE defect_id = 'DEF-000006';
UPDATE defects SET subject = '模块21' WHERE defect_id = 'DEF-000007';
UPDATE defects SET subject = '模块22' WHERE defect_id = 'DEF-000008';
UPDATE defects SET subject = '模块23' WHERE defect_id = 'DEF-000009';
UPDATE defects SET subject = '模块24' WHERE defect_id = 'DEF-000010';
UPDATE defects SET subject = '模块25' WHERE defect_id = 'DEF-000011';
UPDATE defects SET subject = '模块26' WHERE defect_id = 'DEF-000012';
UPDATE defects SET subject = '模块27' WHERE defect_id = 'DEF-000013';
UPDATE defects SET subject = '模块28' WHERE defect_id = 'DEF-000014';
UPDATE defects SET subject = '模块29' WHERE defect_id = 'DEF-000015';
UPDATE defects SET subject = '模块30' WHERE defect_id = 'DEF-000016';
UPDATE defects SET subject = '模块31' WHERE defect_id = 'DEF-000017';
UPDATE defects SET subject = '模块32' WHERE defect_id = 'DEF-000018';
UPDATE defects SET subject = '模块33' WHERE defect_id = 'DEF-000019';
UPDATE defects SET subject = '模块34' WHERE defect_id = 'DEF-000020';

-- 验证修复结果
SELECT defect_id, title, subject, phase
FROM defects 
WHERE defect_id >= 'DEF-000006' AND defect_id <= 'DEF-000020'
ORDER BY defect_id;
