-- 创建手工测试用例表
-- 支持多数据库: PostgreSQL/SQLite/MongoDB
-- 创建日期: 2025-11-06
-- 任务: T10-手工测试用例-元数据编辑

CREATE TABLE IF NOT EXISTS manual_test_cases (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    project_id INTEGER NOT NULL,
    language VARCHAR(20) DEFAULT '中文',
    test_version VARCHAR(50),
    test_env VARCHAR(100),
    test_date VARCHAR(20),
    executor VARCHAR(50),
    case_number VARCHAR(50),
    major_function VARCHAR(100),
    middle_function VARCHAR(100),
    minor_function VARCHAR(100),
    precondition TEXT,
    test_steps TEXT,
    expected_result TEXT,
    test_result VARCHAR(10) DEFAULT 'NR',
    remark TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_mtc_project ON manual_test_cases(project_id);
CREATE INDEX IF NOT EXISTS idx_mtc_language ON manual_test_cases(language);
CREATE INDEX IF NOT EXISTS idx_mtc_major_func ON manual_test_cases(major_function);
CREATE INDEX IF NOT EXISTS idx_mtc_result ON manual_test_cases(test_result);
CREATE INDEX IF NOT EXISTS idx_mtc_deleted_at ON manual_test_cases(deleted_at);

-- SQLite 注释(使用虚拟注释列)
-- 表说明: 手工测试用例表,支持元数据管理和多语言筛选
-- language字段: 用例语言版本(中文/English/日本語)
-- test_result字段: 测试结果(OK/NG/Block/NR)
-- major_function/middle_function/minor_function: 三级功能分类
