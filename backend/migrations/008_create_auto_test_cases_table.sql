-- 008_create_auto_test_cases_table.sql
-- 功能: 创建自动化测试用例表(auto_test_cases),支持ROLE1-ROLE4分类管理
-- 多语言字段: 15个(画面/功能/前置条件/测试步骤/期待值 × CN/JP/EN)
-- 主键设计: UUID(case_id) + 显示序号(id)

BEGIN;

-- ========== 创建auto_test_cases表 ==========

CREATE TABLE auto_test_cases (
    -- 主键(UUID)
    case_id VARCHAR(36) PRIMARY KEY,  -- UUID主键,全局唯一标识
    
    -- 显示序号和基本信息
    id INTEGER NOT NULL,              -- 显示序号(可重排,用于前端展示)
    project_id INTEGER NOT NULL,
    case_type VARCHAR(20) NOT NULL DEFAULT 'role1' CHECK (case_type IN ('role1','role2','role3','role4')),
    
    -- 元数据(冗余存储便于导出)
    test_version VARCHAR(50),         -- 测试版本
    test_date VARCHAR(20),            -- 测试日期(YYYY-MM-DD格式)
    
    -- 公共字段(不区分语言)
    case_number VARCHAR(50),          -- 用例编号
    test_result VARCHAR(10) DEFAULT 'NR' CHECK (test_result IN ('OK','NG','NR')),  -- 测试结果:OK通过/NG失败/NR未执行
    remark TEXT,                      -- 备考
    
    -- 多语言字段 - 画面(三语言)
    screen_cn VARCHAR(100),
    screen_jp VARCHAR(100),
    screen_en VARCHAR(100),
    
    -- 多语言字段 - 功能(三语言,简化为单一功能字段)
    function_cn VARCHAR(200),
    function_jp VARCHAR(200),
    function_en VARCHAR(200),
    
    -- 多语言字段 - 前置条件(三语言)
    precondition_cn TEXT,
    precondition_jp TEXT,
    precondition_en TEXT,
    
    -- 多语言字段 - 测试步骤(三语言)
    test_steps_cn TEXT,
    test_steps_jp TEXT,
    test_steps_en TEXT,
    
    -- 多语言字段 - 期待值(三语言)
    expected_result_cn TEXT,
    expected_result_jp TEXT,
    expected_result_en TEXT,
    
    -- 审计字段
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP              -- 软删除标记
);

-- ========== 索引设计 ==========

-- 项目和类型组合索引(最常用查询)
CREATE INDEX idx_atc_project ON auto_test_cases(project_id);
CREATE INDEX idx_atc_type ON auto_test_cases(case_type);

-- 显示序号索引(用于排序)
CREATE INDEX idx_atc_display_id ON auto_test_cases(id);

-- 测试结果索引(用于统计)
CREATE INDEX idx_atc_result ON auto_test_cases(test_result);

-- 功能字段索引(用于搜索筛选)
CREATE INDEX idx_atc_function_cn ON auto_test_cases(function_cn);
CREATE INDEX idx_atc_function_jp ON auto_test_cases(function_jp);
CREATE INDEX idx_atc_function_en ON auto_test_cases(function_en);

-- 软删除索引
CREATE INDEX idx_atc_deleted_at ON auto_test_cases(deleted_at);

-- ========== 外键约束 ==========

-- 引用projects表(级联删除)
-- 注意: SQLite不支持ALTER TABLE ADD CONSTRAINT语法,需在建表时定义或使用PRAGMA

-- 若数据库为PostgreSQL/MySQL,使用以下语句:
-- ALTER TABLE auto_test_cases ADD CONSTRAINT fk_atc_project 
--     FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE;

-- SQLite处理方式: 在应用层保证引用完整性,或重建表时添加外键

COMMIT;

-- ========== 验证脚本(可选执行) ==========
-- SELECT name FROM sqlite_master WHERE type='table' AND name='auto_test_cases';
-- PRAGMA table_info(auto_test_cases);
