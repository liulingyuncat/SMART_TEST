# SMART-TEST 特性文档

## 1. 文档概述



### 1.1 文档目的

本文档旨在描述 SMART-TEST 平台的功能特性和需求规格，作为产品开发、测试和验收的依据。文档从代码实现反向整理，确保与实际系统保持一致。

### 1.2 适用范围

本文档适用于：
- 产品经理：了解系统功能范围
- 开发团队：明确功能实现要求
- 测试团队：制定测试策略和用例
- 运维团队：了解系统部署和运维需求
- AI 助手：通过 MCP 协议与系统交互

### 1.3 术语定义

| 术语 | 定义 |
|------|------|
| MCP | Model Context Protocol，模型上下文协议，用于 AI 助手与平台交互 |
| 用例集 | 测试用例的逻辑分组，支持手工、Web自动化、API 三种类型 |
| 测试观点 | 从需求衍生的测试设计视角，指导用例编写 |
| 执行任务 | 基于用例集创建的测试执行计划 |
| 缺陷 | 测试过程中发现的问题记录 |
| 版本 | 用例或需求的历史快照，支持版本管理 |

## 2. 产品概述



### 2.1 产品定位

SMART-TEST 是一个 AI 辅助的智能软件测试平台，定位为：

- **人机协作测试平台**：通过 Web 界面和 MCP 协议，实现人类测试人员与 AI 助手的无缝协作
- **测试全生命周期管理**：覆盖需求分析、测试设计、用例管理、执行跟踪、缺陷管理全流程
- **知识中心**：构建统一的测试知识库，支持多版本、多语言管理

### 2.2 目标用户

| 用户类型 | 使用场景 |
|----------|----------|
| 测试工程师 | 编写/执行测试用例、提交缺陷、查看测试报告 |
| 测试经理 | 管理项目、分配任务、监控测试进度 |
| 开发人员 | 查看缺陷详情、验证修复结果 |
| AI 助手 | 通过 MCP 协议自动化测试设计、生成报告 |
| 系统管理员 | 用户管理、系统配置、权限分配 |

### 2.3 核心价值

1. **AI 能力增强**：利用 LLM 模型辅助测试设计，提升效率和质量
2. **单一事实来源**：统一管理需求和用例，消除信息孤岛
3. **多语言支持**：支持中文、英文、日文等多语言用例管理
4. **自动化集成**：内置 Playwright Web 自动化执行能力
5. **灵活扩展**：MCP 协议支持多种 AI 助手集成

## 3. 功能特性



### 3.1 用户与权限管理



#### 3.1.1 用户认证

**功能描述**：
- 用户名/密码登录认证
- JWT Token 机制，支持 Token 刷新
- 会话管理与自动登出

**实现依据**：
- `backend/internal/handlers/auth.go`：登录/登出接口
- `backend/internal/services/auth_service.go`：认证逻辑
- `frontend/src/api/auth.js`：前端认证 API

#### 3.1.2 角色权限

**角色类型**：
| 角色 | 权限说明 |
|------|----------|
| system_admin | 系统管理员，拥有所有权限 |
| project_admin | 项目管理员，管理项目成员和配置 |
| member | 普通成员，执行测试任务 |

**实现依据**：
- `backend/internal/models/user.go`：用户模型定义
- `backend/internal/middleware/`：权限中间件

#### 3.1.3 项目成员管理

**功能描述**：
- 添加/移除项目成员
- 分配成员角色
- 成员权限查询

**实现依据**：
- `backend/internal/models/project_member.go`：成员模型
- `frontend/src/components/MemberTransfer.jsx`：成员管理组件

### 3.2 项目管理



#### 3.2.1 项目创建与配置

**功能描述**：
- 创建新项目
- 配置项目基本信息（名称、描述）
- 配置缺陷分类（主题、阶段）

**实现依据**：
- `backend/internal/handlers/project_handler.go`：项目接口
- `backend/internal/services/project_service.go`：项目服务
- `frontend/src/pages/ProjectManagement/`：项目管理页面

#### 3.2.2 项目信息管理

**功能描述**：
- 查看项目详情
- 编辑项目信息
- 项目列表展示与筛选

**实现依据**：
- `frontend/src/pages/ProjectList/`：项目列表
- `frontend/src/components/ProjectInfoTab.jsx`：项目信息标签页

### 3.3 需求管理



#### 3.3.1 原始文档导入

**功能描述**：
- 支持多种格式文档导入
- 文档解析与内容提取
- 原始文档存储与版本管理

**实现依据**：
- `backend/internal/handlers/raw_document_handler.go`
- `backend/internal/services/raw_document_service.go`
- `frontend/src/pages/ProjectDetail/RawDocumentTab/`

#### 3.3.2 需求条目管理

**功能描述**：
- 创建、编辑、删除需求条目
- 需求层级结构支持
- 需求与测试观点关联

**实现依据**：
- `backend/internal/handlers/requirement_item_handler.go`
- `backend/internal/services/requirement_item_service.go`
- `backend/internal/models/requirement_item.go`
- `frontend/src/pages/ProjectDetail/RequirementManagement/`

#### 3.3.3 需求版本控制

**功能描述**：
- 需求变更历史记录
- 版本对比与回滚
- 版本备注

**实现依据**：
- `backend/internal/handlers/version_handler.go`
- `backend/internal/services/version_service.go`
- `backend/internal/models/version.go`

### 3.4 测试观点管理



#### 3.4.1 观点创建与编辑

**功能描述**：
- 创建测试观点（Markdown 格式）
- 编辑和删除观点
- 观点分类管理

**实现依据**：
- `backend/internal/handlers/viewpoint_item_handler.go`
- `backend/internal/services/viewpoint_item_service.go`
- `backend/internal/models/viewpoint_item.go`

#### 3.4.2 观点与需求关联

**功能描述**：
- 建立观点与需求的关联
- 追溯测试覆盖度

**实现依据**：
- MCP viewpoint 工具支持关联查询

### 3.5 测试用例管理



#### 3.5.1 用例集管理

**功能描述**：
- 创建用例集（Case Group）
- 支持三种类型：手工测试、Web 自动化、API 接口
- 用例集描述和元数据管理

**实现依据**：
- `backend/internal/handlers/case_group_handler.go`
- `backend/internal/models/case_group.go`

#### 3.5.2 手工测试用例

**功能描述**：
- 创建、编辑、删除手工测试用例
- 支持批量操作
- 用例字段：编号、标题、前置条件、步骤、预期结果

**实现依据**：
- `backend/internal/handlers/manual_cases_handler.go`
- `backend/internal/services/manual_test_case_service.go`
- `backend/internal/models/manual_test_case.go`
- `frontend/src/pages/ProjectDetail/ManualTestTabs/`

#### 3.5.3 Web自动化用例

**功能描述**：
- 基于 Playwright 的 Web 自动化用例
- 用例字段包括页面、操作、定位器、断言
- 支持自动化执行

**实现依据**：
- backend/internal/handlers/auto_test_case.go
- backend/internal/services/auto_test_case_service.go
- frontend/src/pages/ProjectDetail/AutoTestTabs/

#### 3.5.4 API接口用例

**功能描述**：
- API 接口测试用例管理
- 用例字段包括URL、方法、请求参数、期望响应
- 支持批量创建和更新

**实现依据**：
- backend/internal/handlers/api_test_case.go
- backend/internal/services/api_test_case_service.go
- frontend/src/pages/ProjectDetail/ApiTestTabs/

#### 3.5.5 用例版本管理

**功能描述**：
- 用例变更历史记录
- 版本对比与回滚

**实现依据**：
- case_version.go 和 web_case_version.go 模型
- web_version_service.go 服务层

#### 3.5.6 多语言支持

**功能描述**：
- 支持中文(CN)、英文(EN)、日文(JP)
- 多语言用例字段存储
- 语言切换和翻译辅助

**实现依据**：
- 数据模型中的 `_cn`、`_en`、`_jp` 后缀字段
- `frontend/src/i18n.js`：国际化配置
- `frontend/src/locales/`：语言包

### 3.6 用例评审



#### 3.6.1 评审流程

**功能描述**：
- 创建用例评审任务
- 评审意见记录
- 评审状态跟踪

**实现依据**：
- backend/internal/handlers/review_handler.go
- backend/internal/services/review_service.go
- backend/internal/models/case_review.go

#### 3.6.2 评审条目管理

**功能描述**：
- 创建评审条目
- 评审意见分类
- AI 辅助评审建议

**实现依据**：
- backend/internal/handlers/review_item_handler.go
- backend/internal/services/review_item_service.go
- backend/internal/models/case_review_item.go

### 3.7 测试执行



#### 3.7.1 执行任务管理

**功能描述**：
- 创建执行任务
- 关联用例集
- 任务状态管理

**实现依据**：
- backend/internal/handlers/execution_task_handler.go
- backend/internal/services/execution_task_service.go
- backend/internal/models/execution_task.go
- frontend/src/pages/ProjectDetail/TestExecution/

#### 3.7.2 执行结果记录

**功能描述**：
- 记录用例执行结果（OK/NG/BLOCK/NR）
- 执行备注和截图
- 结果统计分析

**实现依据**：
- backend/internal/handlers/execution_case_result_handler.go
- backend/internal/services/execution_case_result_service.go
- backend/internal/models/execution_case_result.go

#### 3.7.3 Playwright自动化执行

**功能描述**：
- 基于 Playwright 的 Web 自动化执行
- 执行结果自动记录
- 失败截图和日志

**实现依据**：
- Web 自动化用例集成 Playwright 框架

### 3.8 缺陷管理



#### 3.8.1 缺陷创建与跟踪

**功能描述**：
- 创建缺陷记录
- 缺陷状态跟踪（Open/Fixed/Closed）
- 缺陷与用例关联

**实现依据**：
- backend/internal/handlers/defect_handler.go
- backend/internal/services/defect_service.go
- backend/internal/models/defect.go
- frontend/src/pages/ProjectDetail/DefectManagement/

#### 3.8.2 缺陷分类配置

**功能描述**：
- 配置缺陷主题分类
- 配置缺陷阶段分类
- 自定义分类管理

**实现依据**：
- backend/internal/handlers/defect_config_handler.go
- backend/internal/services/defect_config_service.go
- backend/internal/models/defect_subject.go
- backend/internal/models/defect_phase.go

#### 3.8.3 缺陷评论与附件

**功能描述**：
- 缺陷评论功能
- 附件上传和管理
- 评论历史记录

**实现依据**：
- backend/internal/handlers/defect_comment_handler.go
- backend/internal/handlers/defect_attachment_handler.go
- backend/internal/models/defect_comment.go
- backend/internal/models/defect_attachment.go

### 3.9 AI能力集成



#### 3.9.1 MCP协议支持

**功能描述**：
- 实现 MCP 协议
- 支持 stdio 和 HTTP 传输
- 提供 39 个 MCP 工具

**实现依据**：
- backend/cmd/mcp/main.go
- backend/internal/mcp/server.go

#### 3.9.2 AI辅助测试设计

**功能描述**：
- AI 辅助需求分析
- AI 辅助观点生成
- AI 辅助用例设计

**实现依据**：
- MCP 工具集提供需求、观点、用例的 CRUD 接口
- 外部 AI 助手通过 MCP 协议调用

#### 3.9.3 AI质量报告

**功能描述**：
- AI 生成测试质量报告
- 报告内容支持 Markdown 格式
- 报告版本管理

**实现依据**：
- backend/internal/handlers/ai_report_handler.go
- backend/internal/services/ai_report_service.go
- backend/internal/models/ai_report.go
- frontend/src/pages/ProjectDetail/AIReportManagement/

## 4. MCP工具集



### 4.1 用户与项目工具

- get_current_user_info：获取当前用户信息
- get_current_project_name：获取当前项目信息

### 4.2 需求与观点工具

| 工具名称 | 功能说明 |
|----------|----------|
| list_requirement_items | 获取需求列表 |
| get_requirement_item | 获取需求详情 |
| create_requirement_item | 创建需求 |
| update_requirement_item | 更新需求 |
| list_viewpoint_items | 获取观点列表 |
| get_viewpoint_item | 获取观点详情 |
| create_viewpoint_item | 创建观点 |
| update_viewpoint_item | 更新观点 |

### 4.3 用例管理工具

**手工用例工具**：
- list_manual_groups / list_manual_cases
- create_case_group / create_manual_cases
- update_manual_case / update_manual_cases

**Web自动化工具**：
- list_web_groups / list_web_cases
- get_web_group_metadata
- create_web_group / create_web_cases
- update_web_cases

**API用例工具**：
- list_api_groups / list_api_cases
- get_api_group_metadata
- create_api_group / create_api_case
- update_api_case

### 4.4 执行与缺陷工具

| 工具名称 | 功能说明 |
|----------|----------|
| list_execution_tasks | 获取执行任务列表 |
| get_execution_task_cases | 获取任务用例 |
| update_execution_case_result | 更新执行结果 |
| list_defects | 获取缺陷列表 |
| update_defect | 更新缺陷 |

### 4.5 AI报告工具

- create_ai_report：创建 AI 测试报告
- update_ai_report：更新 AI 测试报告

## 5. 非功能性需求



### 5.1 性能需求

- 支持 100 并发用户
- API 响应时间小于 500ms
- 数据库查询优化

### 5.2 安全需求

- HTTPS 加密传输
- JWT Token 认证
- 密码加密存储

### 5.3 可用性需求

- 响应式 Web 界面
- 友好的错误提示
- 操作确认机制

### 5.4 国际化需求

- 界面支持中/英/日三语切换
- 用例内容支持多语言存储
- Ant Design 国际化组件