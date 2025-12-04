-- 010_create_api_test_cases_table.sql
-- 功能: 创建接口测试用例表(api_test_cases)和版本管理表(api_test_case_versions)
-- 主键设计: UUID(id VARCHAR(36)) + 显示序号(display_order INTEGER)
-- 支持: ROLE1-ROLE4分类管理,版本管理(CSV导出+ZIP下载)

BEGIN;

-- ========== 创建api_test_cases表 ==========

CREATE TABLE api_test_cases (
    -- 主键(UUID)
    id VARCHAR(36) PRIMARY KEY,       -- UUID主键,全局唯一标识符
    
    -- 基本信息
    project_id INTEGER NOT NULL,
    case_type VARCHAR(20) NOT NULL DEFAULT 'role1' CHECK (case_type IN ('role1','role2','role3','role4')),
    
    -- 用例基本字段
    case_number VARCHAR(50),          -- 用例编号(用户自定义)
    screen VARCHAR(100),              -- 画面/接口所属模块
    
    -- API专属字段
    url TEXT,                         -- 接口地址
    header TEXT,                      -- 请求头(支持多行)
    method VARCHAR(10) NOT NULL DEFAULT 'GET' CHECK (method IN ('GET','POST','PUT','DELETE','PATCH')),  -- HTTP方法
    body TEXT,                        -- 请求体(支持多行)
    response TEXT,                    -- 预期响应(支持多行)
    
    -- 测试结果
    test_result VARCHAR(10) DEFAULT 'NR' CHECK (test_result IN ('OK','NG','NR')),  -- 测试结果:OK通过/NG失败/NR未执行
    remark TEXT,                      -- 备注
    
    -- 显示顺序(用于计算前端No.列序号)
    display_order INTEGER,            -- 显示顺序,前端根据此字段计算No.序号(1,2,3...)
    
    -- 审计字段
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- 外键约束
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

-- 索引
CREATE INDEX idx_api_cases_project_type ON api_test_cases(project_id, case_type);
CREATE INDEX idx_api_cases_display_order ON api_test_cases(project_id, case_type, display_order);

-- ========== 创建api_test_case_versions表 ==========

CREATE TABLE api_test_case_versions (
    -- 主键(UUID)
    id VARCHAR(36) PRIMARY KEY,       -- 版本UUID
    
    -- 关联信息
    project_id INTEGER NOT NULL,
    
    -- 四个CSV文件名(每个ROLE一个)
    filename_role1 VARCHAR(255) NOT NULL,  -- ROLE1文件名
    filename_role2 VARCHAR(255) NOT NULL,  -- ROLE2文件名
    filename_role3 VARCHAR(255) NOT NULL,  -- ROLE3文件名
    filename_role4 VARCHAR(255) NOT NULL,  -- ROLE4文件名
    
    -- 版本信息
    remark TEXT,                      -- 版本备注(限制500字符,前端校验)
    created_by INTEGER NOT NULL,      -- 创建人ID
    
    -- 审计字段
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- 外键约束
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- 索引
CREATE INDEX idx_api_versions_project ON api_test_case_versions(project_id);
CREATE INDEX idx_api_versions_created_at ON api_test_case_versions(created_at DESC);

COMMIT;

-- ========== 回滚脚本 ==========
-- 执行回滚: sqlite3 webtest.db < 010_down_api_test_cases_table.sql

-- BEGIN;
-- DROP TABLE IF EXISTS api_test_case_versions;
-- DROP TABLE IF EXISTS api_test_cases;
-- COMMIT;
