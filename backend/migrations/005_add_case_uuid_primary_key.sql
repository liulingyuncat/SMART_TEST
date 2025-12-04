-- 迁移说明: 将ID从主键改为普通序号字段，添加case_id(UUID)作为新主键
-- 创建日期: 2025-11-12
-- 任务: T11-手工测试用例-表格CRUD
-- 重要: 执行前请备份数据库！

-- 步骤1: 创建新表结构（带UUID主键和普通ID字段）
CREATE TABLE IF NOT EXISTS manual_test_cases_new (
    case_id VARCHAR(36) PRIMARY KEY,  -- UUID主键
    id INTEGER NOT NULL,               -- 显示序号（普通字段，可修改，用于重排序）
    project_id INTEGER NOT NULL,
    case_type VARCHAR(20) DEFAULT 'overall',
    test_version VARCHAR(50),
    test_env VARCHAR(100),
    test_date VARCHAR(20),
    executor VARCHAR(50),
    case_number VARCHAR(50),
    
    -- 单语言字段(AI用例使用)
    major_function VARCHAR(100),
    middle_function VARCHAR(100),
    minor_function VARCHAR(100),
    precondition TEXT,
    test_steps TEXT,
    expected_result TEXT,
    
    -- 多语言字段(整体/变更用例使用)
    major_function_cn VARCHAR(100),
    major_function_jp VARCHAR(100),
    major_function_en VARCHAR(100),
    middle_function_cn VARCHAR(100),
    middle_function_jp VARCHAR(100),
    middle_function_en VARCHAR(100),
    minor_function_cn VARCHAR(100),
    minor_function_jp VARCHAR(100),
    minor_function_en VARCHAR(100),
    precondition_cn TEXT,
    precondition_jp TEXT,
    precondition_en TEXT,
    test_steps_cn TEXT,
    test_steps_jp TEXT,
    test_steps_en TEXT,
    expected_result_cn TEXT,
    expected_result_jp TEXT,
    expected_result_en TEXT,
    
    -- 共用字段
    test_result VARCHAR(10) DEFAULT 'NR',
    remark TEXT,
    language VARCHAR(20) DEFAULT '中文',
    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- 步骤2: 迁移旧数据（为每条记录生成UUID）
INSERT INTO manual_test_cases_new (
    case_id, id, project_id, case_type, test_version, test_env, test_date, executor, case_number,
    major_function, middle_function, minor_function, precondition, test_steps, expected_result,
    major_function_cn, major_function_jp, major_function_en,
    middle_function_cn, middle_function_jp, middle_function_en,
    minor_function_cn, minor_function_jp, minor_function_en,
    precondition_cn, precondition_jp, precondition_en,
    test_steps_cn, test_steps_jp, test_steps_en,
    expected_result_cn, expected_result_jp, expected_result_en,
    test_result, remark, language, created_at, updated_at, deleted_at
)
SELECT 
    lower(hex(randomblob(16))) as case_id,  -- 生成UUID
    id, project_id, case_type, test_version, test_env, test_date, executor, case_number,
    major_function, middle_function, minor_function, precondition, test_steps, expected_result,
    major_function_cn, major_function_jp, major_function_en,
    middle_function_cn, middle_function_jp, middle_function_en,
    minor_function_cn, minor_function_jp, minor_function_en,
    precondition_cn, precondition_jp, precondition_en,
    test_steps_cn, test_steps_jp, test_steps_en,
    expected_result_cn, expected_result_jp, expected_result_en,
    test_result, remark, language, created_at, updated_at, deleted_at
FROM manual_test_cases;

-- 步骤3: 删除旧表
DROP TABLE manual_test_cases;

-- 步骤4: 重命名新表
ALTER TABLE manual_test_cases_new RENAME TO manual_test_cases;

-- 步骤5: 重建索引
CREATE INDEX IF NOT EXISTS idx_mtc_project ON manual_test_cases(project_id);
CREATE INDEX IF NOT EXISTS idx_mtc_type ON manual_test_cases(case_type);
CREATE INDEX IF NOT EXISTS idx_mtc_display_id ON manual_test_cases(id);
CREATE INDEX IF NOT EXISTS idx_mtc_major_func ON manual_test_cases(major_function);
CREATE INDEX IF NOT EXISTS idx_mtc_major_func_cn ON manual_test_cases(major_function_cn);
CREATE INDEX IF NOT EXISTS idx_mtc_major_func_jp ON manual_test_cases(major_function_jp);
CREATE INDEX IF NOT EXISTS idx_mtc_major_func_en ON manual_test_cases(major_function_en);
CREATE INDEX IF NOT EXISTS idx_mtc_language ON manual_test_cases(language);
CREATE INDEX IF NOT EXISTS idx_mtc_deleted_at ON manual_test_cases(deleted_at);

-- SQLite 注释
-- 迁移说明: case_id现在是UUID主键，id是普通显示序号字段
-- id字段可以自由修改，用于重新排序功能
