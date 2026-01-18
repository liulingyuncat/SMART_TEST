# SMART-TEST 架构设计文档

## 1. 文档概述



### 1.1 文档目的

本文档描述 SMART-TEST 平台的系统架构设计，包括后端服务、前端应用、部署架构等技术细节，作为开发和运维的技术参考。

### 1.2 适用范围

本文档适用于：
- 后端开发人员
- 前端开发人员
- DevOps 工程师
- 系统架构师

### 1.3 术语定义

| 术语 | 定义 |
|------|------|
| Handler | 处理 HTTP 请求的控制器层 |
| Service | 业务逻辑层 |
| Repository | 数据访问层 |
| GORM | Go 语言 ORM 框架 |

## 2. 系统架构概览



### 2.1 架构原则

1. **分层架构**：清晰的 Handler-Service-Repository 三层结构
2. **关注点分离**：Web 服务与 MCP 服务独立部署
3. **可扩展性**：支持多数据库、可插拔的 MCP 工具
4. **容器化优先**：Docker 容器化部署

### 2.2 整体架构图

```
客户端层
  Web 浏览器 / AI 助手 / Playwright Runner
         |
服务层
  Web 服务 (Gin, :8443)    MCP 服务 (:16410)
         |
数据层
  PostgreSQL / SQLite
```

### 2.3 技术栈选型

| 层级 | 技术选型 | 说明 |
|------|----------|------|
| 后端语言 | Go 1.24 | 高性能、并发支持 |
| Web 框架 | Gin | 轻量级 HTTP 框架 |
| ORM | GORM | Go ORM 框架 |
| 数据库 | PostgreSQL/SQLite | 多数据库支持 |
| 前端框架 | React 18 | 组件化开发 |
| UI 组件库 | Ant Design | 企业级 UI |
| 状态管理 | Redux Toolkit | 状态管理 |
| 容器化 | Docker | 容器部署 |
| CI/CD | GitHub Actions | 持续集成 |

## 3. 后端架构



### 3.1 Web服务架构



#### 3.1.1 入口与启动流程

**入口文件**：backend/cmd/server/main.go

**启动流程**：
1. 加载数据库配置
2. 初始化数据库连接
3. 执行自动迁移
4. 初始化 Repository 层
5. 初始化 Service 层
6. 注册路由和中间件
7. 启动 HTTPS 服务

#### 3.1.2 路由与中间件

**Web 框架**：Gin

**中间件**：
- 认证中间件：JWT Token 验证
- CORS 中间件：跨域支持
- 日志中间件：请求日志记录

**路由结构**：
- /api/v1/auth：认证相关
- /api/v1/users：用户管理
- /api/v1/projects：项目管理
- /api/v1/cases：用例管理

#### 3.1.3 分层架构

**三层架构**：
```
Handlers (internal/handlers/)
    ↓
Services (internal/services/)
    ↓
Repositories (internal/repositories/)
    ↓
Models (internal/models/)
```

**职责划分**：
- Handlers：请求解析、响应封装
- Services：业务逻辑、事务管理
- Repositories：数据访问、CRUD 操作

### 3.2 MCP服务架构



#### 3.2.1 MCP协议实现

**协议版本**：2025-06-18

**实现文件**：backend/internal/mcp/

**核心组件**：
- server.go：MCP 服务器主体
- protocol/：消息路由和协议处理
- tools/：工具注册和调用
- transport/：传输层（stdio/HTTP）
- prompts/：提示词管理

#### 3.2.2 工具注册机制

**工具目录**：backend/internal/mcp/tools/handlers/

**注册机制**：
1. 定义工具结构和参数 Schema
2. 实现 Handler 函数
3. 在 register.go 中注册到 ToolRegistry
4. 服务启动时自动加载

**工具分类**：
- user_project.go：用户项目工具
- requirement.go：需求工具
- viewpoint.go：观点工具
- manual_case.go：手工用例工具
- web_case.go：Web 用例工具
- api_case.go：API 用例工具
- execution.go：执行任务工具
- defect.go：缺陷工具
- report.go：报告工具

#### 3.2.3 传输层设计

**支持的传输方式**：
1. **stdio**：标准输入输出，用于 CLI 集成
2. **HTTP SSE**：HTTP Server-Sent Events

**实现文件**：backend/internal/mcp/transport/

**配置方式**：通过 mcp-server.yaml 配置文件

### 3.3 数据层设计



#### 3.3.1 数据模型

**核心模型**：
| 模型 | 说明 |
|------|------|
| User | 用户信息 |
| Project | 项目信息 |
| ProjectMember | 项目成员关系 |
| RequirementItem | 需求条目 |
| ViewpointItem | 测试观点 |
| ManualTestCase | 手工测试用例 |
| AutoTestCase | Web 自动化用例 |
| ApiTestCase | API 接口用例 |
| ExecutionTask | 执行任务 |
| Defect | 缺陷 |
| AIReport | AI 报告 |

#### 3.3.2 仓储模式

**设计模式**：Repository Pattern

**目录**：backend/internal/repositories/

**职责**：封装数据库操作，提供统一 CRUD 接口

#### 3.3.3 数据库支持

**支持的数据库**：
1. **PostgreSQL**：生产环境推荐
2. **SQLite**：开发和轻量部署

**配置方式**：通过环境变量 DB_TYPE 切换

**ORM**：GORM，支持自动迁移

## 4. 前端架构



### 4.1 技术框架

**核心技术**：
- React 18
- Ant Design 5.x
- Redux Toolkit
- React Router
- i18next（国际化）
- Axios（HTTP 客户端）

### 4.2 目录结构

```
frontend/src/
├── api/          # API 接口封装
├── components/   # 公共组件
├── pages/        # 页面组件
├── store/        # Redux 状态管理
├── router/       # 路由配置
├── locales/      # 语言包
├── utils/        # 工具函数
└── App.js        # 应用入口
```

### 4.3 状态管理

**技术**：Redux Toolkit

**Store 结构**：
- authSlice：认证状态
- projectSlice：项目状态

**目录**：frontend/src/store/

### 4.4 路由设计

**技术**：React Router

**主要路由**：
- /login：登录页
- /projects：项目列表
- /project/:id：项目详情
- /users：用户管理
- /profile：个人中心

### 4.5 组件设计

**组件分类**：
- 布局组件：MainLayout, Header, Sidebar
- 业务组件：ProjectInfoTab, MemberTransfer
- 通用组件：RoleGuard, LanguageSwitch

**UI 框架**：Ant Design 5.x

### 4.6 国际化实现

**技术**：i18next

**支持语言**：
- zh：中文
- en：英文
- ja：日文

**语言包目录**：frontend/src/locales/

## 5. 部署架构



### 5.1 容器化设计

**Dockerfile**：多阶段构建

构建阶段：
1. 后端构建：Go 编译
2. 前端构建：npm build
3. 运行阶段：Alpine 最小镜像

**镜像仓库**：ghcr.io

### 5.2 Docker Compose编排

**服务定义**：
- webtest：主服务（Web + MCP）
- postgres：PostgreSQL 数据库

**端口映射**：
- 8443：Web 服务（HTTPS）
- 16410：MCP 服务

**持久化**：
- ./storage：应用存储
- ./data/postgres：数据库文件

### 5.3 CI/CD流水线

**平台**：GitHub Actions

**工作流**：
1. ci.yml：持续集成
   - 后端测试（go vet, go test）
   - 前端测试（npm lint）
   - Docker 镜像构建和推送

2. docker-build.yml：手动构建
   - 支持自定义标签
   - 多平台支持（amd64, arm64）

### 5.4 环境配置

**环境变量**：
| 变量 | 说明 |
|------|------|
| DB_TYPE | 数据库类型 |
| DB_HOST | 数据库主机 |
| DB_PORT | 数据库端口 |
| DB_USER | 数据库用户 |
| DB_PASSWORD | 数据库密码 |
| JWT_SECRET | JWT 密钥 |
| MCP_AUTH_TOKEN | MCP 认证令牌 |

## 6. 安全架构



### 6.1 认证机制

**Web 服务认证**：
- JWT Token 机制
- Token 存储于 localStorage

**MCP 服务认证**：
- Bearer Token 认证
- 支持动态 Token 模式

### 6.2 授权控制

**角色权限**：
- system_admin：系统管理员
- project_admin：项目管理员
- member：普通成员

**实现方式**：中间件拦截检查

### 6.3 数据安全

**安全措施**：
- HTTPS 加密传输
- 密码 bcrypt 加密存储
- 敏感数据脱敏

## 7. 扩展性设计



### 7.1 MCP工具扩展

**扩展步骤**：
1. 创建新工具文件
2. 定义工具 Schema 和 Handler
3. 在 register.go 中注册
4. 重启 MCP 服务

### 7.2 多数据库支持

**支持的数据库**：
- PostgreSQL 生产推荐
- SQLite 开发轻量

**切换方式**：设置 DB_TYPE 环境变量

### 7.3 插件化架构

**扩展点**：
- MCP 工具插件化注册
- 数据库驱动可替换
- 前端组件模块化
- 多语言动态加载

**未来规划**：
- 更多 AI 模型接口支持
- 测试报告模板定制
- 第三方系统集成 API