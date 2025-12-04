SELECT id, 
       CASE 
         WHEN case_id IS NULL THEN 'NULL' 
         WHEN case_id = '' THEN 'EMPTY' 
         ELSE substr(case_id, 1, 8) 
       END as case_id_status,
       case_type
FROM manual_test_cases 
WHERE case_type='overall' 
ORDER BY id 
LIMIT 10;
