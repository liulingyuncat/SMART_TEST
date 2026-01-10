-- 051_add_case_group_and_web_versions.sql
-- 功能: AIWeb用例库改造-1 数据库变更
-- 说明: 
--   1. 为auto_test_cases表新增case_group字段，支持Web用例集管理
--   2. 创建web_case_versions表，支持Web用例多语言版本管理
-- 创建日期: 2025-12-11

-- ========== 扩展auto_test_cases表 ==========
-- 新增case_group字段用于关联用例集
ALTER TABLE auto_test_cases ADD COLUMN case_group VARCHAR(100) DEFAULT '';

-- 创建case_group字段索引，优化按用例集查询性能
CREATE INDEX IF NOT EXISTS idx_atc_case_group ON auto_test_cases(case_group);

-- ========== 创建web_case_versions表 ==========
-- Web用例版本管理表，存储多语言版本保存记录
CREATE TABLE IF NOT EXISTS web_case_versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    version_id VARCHAR(100) NOT NULL,
    project_id INTEGER NOT NULL,
    project_name VARCHAR(100) NOT NULL,
    zip_filename VARCHAR(255) NOT NULL,
    zip_path VARCHAR(500) NOT NULL,
    file_size BIGINT,
    case_count INTEGER DEFAULT 0,
    remark VARCHAR(200) DEFAULT '',
    created_by INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_wcv_version_id ON web_case_versions(version_id);
CREATE INDEX IF NOT EXISTS idx_wcv_project_id ON web_case_versions(project_id);
CREATE INDEX IF NOT EXISTS idx_wcv_created_by ON web_case_versions(created_by);
CREATE INDEX IF NOT EXISTS idx_wcv_created_at ON web_case_versions(created_at DESC);

-- 表说明:
-- web_case_versions: Web用例版本表
-- version_id: 版本标识，格式为 项目名_时间戳
-- zip_filename: zip文件名，包含CN/JP/EN/All四个Excel文件
-- zip_path: zip文件存储路径 storage/versions/web-cases/{projectId}/
-- file_size: zip文件大小（字节）
-- case_count: 版本包含的用例总数
-- remark: 版本备注信息
