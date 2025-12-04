# Web 智能测试平台 - 用户登录功能

## 项目概述

这是一个现代化的 Web 智能测试平台,采用前后端分离架构。本文档介绍 T01-用户登录功能的实现。

## 技术栈

### 后端
- **语言**: Go 1.21+
- **Web框架**: Gin v1.9.1
- **ORM**: GORM v1.25.5
- **认证**: JWT (golang-jwt/jwt v5.2.0)
- **密码加密**: bcrypt
- **数据库**: PostgreSQL / SQLite / MongoDB

### 前端
- **框架**: React 18+
- **UI库**: Ant Design 5.x
- **状态管理**: Redux Toolkit
- **路由**: React Router v6
- **HTTP客户端**: Axios
- **国际化**: i18next

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 14+
- PostgreSQL 13+ (可选, 默认使用 SQLite)

### 后端启动

```bash
# 1. 进入后端目录
cd backend

# 2. 配置 Go 代理(国内用户)
export GOPROXY=https://goproxy.cn,direct  # Linux/Mac
$env:GOPROXY="https://goproxy.cn,direct"  # Windows PowerShell

# 3. 安装依赖
go mod tidy

# 4. 运行服务
go run cmd/server/main.go

# 服务将在 http://localhost:8080 启动
```

### 前端启动

```bash
# 1. 进入前端目录
cd frontend

# 2. 安装依赖
npm install

# 3. 启动开发服务器
npm start

# 应用将在 http://localhost:3000 启动
```

### 测试登录

1. 访问 http://localhost:3000/login
2. 使用默认管理员账号:
   - 用户名: `admin`
   - 密码: `admin123`
3. 登录成功后自动跳转到首页

## 项目结构

```
webtest/
├── backend/                    # 后端 Go 项目
│   ├── cmd/
│   │   └── server/
│   │       └── main.go         # 程序入口
│   ├── config/
│   │   └── database.go         # 数据库配置
│   ├── internal/
│   │   ├── handlers/
│   │   │   └── auth.go         # 认证处理器
│   │   ├── middleware/
│   │   │   └── auth.go         # 认证中间件
│   │   ├── models/
│   │   │   └── user.go         # User 模型
│   │   ├── repositories/
│   │   │   └── user_repo.go    # 数据访问层
│   │   ├── services/
│   │   │   ├── auth_service.go      # 认证服务
│   │   │   └── auth_service_test.go # 单元测试
│   │   └── utils/
│   │       └── response.go     # 统一响应工具
│   ├── migrations/             # 数据库迁移脚本
│   │   ├── 001_create_users_table.sql
│   │   └── mongodb_schema.js
│   ├── go.mod
│   └── go.sum
│
├── frontend/                   # 前端 React 项目
│   ├── public/
│   ├── src/
│   │   ├── api/
│   │   │   ├── client.js       # Axios 客户端
│   │   │   └── auth.js         # 认证 API
│   │   ├── i18n/
│   │   │   └── index.js        # 国际化配置
│   │   ├── pages/
│   │   │   └── Login/
│   │   │       ├── index.jsx   # 登录页面
│   │   │       ├── index.css
│   │   │       └── Login.test.jsx
│   │   ├── router/
│   │   │   └── index.jsx       # 路由配置
│   │   ├── store/
│   │   │   ├── index.js        # Redux store
│   │   │   └── authSlice.js    # 认证状态
│   │   ├── App.js
│   │   └── index.js
│   ├── .env                    # 环境变量
│   ├── .env.development
│   ├── .env.production
│   ├── jest.config.js
│   └── package.json
│
└── docs/                       # 文档
    ├── API-documentation.md
    ├── backend-deployment.md
    ├── frontend-deployment.md
    ├── manual-test-guide.md
    └── T01-development-summary.md
```

## 核心功能

### 1. 用户登录
- ✅ 用户名/密码表单验证
- ✅ bcrypt 密码加密存储
- ✅ JWT Token 认证
- ✅ 自动跳转到目标页面

### 2. 权限控制
- ✅ JWT 中间件保护路由
- ✅ 前端路由守卫
- ✅ Token 自动刷新重定向

### 3. 需求管理 (T31)
- ✅ 四种需求文档类型支持
  - 整体需求 (Overall Requirements)
  - 整体测试观点 (Overall Test Viewpoint)
  - 变更需求 (Change Requirements)
  - 变更测试观点 (Change Test Viewpoint)
- ✅ Markdown 编辑器集成
  - 编辑/只读模式切换
  - 实时预览渲染
  - 工具栏快捷操作(加粗、斜体、标题、列表等)
- ✅ 版本管理功能
  - 自动生成版本文件名(格式: {项目名}_{文档类型}_{日期}_{时间}.md)
  - 双重保存机制(数据库+文件系统)
  - 版本列表查询(支持按文档类型筛选)
  - 版本下载(Markdown格式)
  - 版本删除(软删除)
- ✅ Markdown 文件导入
  - 文件大小限制(≤5MB)
  - .md 后缀验证
  - UTF-8 编码支持
- ✅ 国际化支持(中文/英文/日文)

**使用指南:**
1. 进入项目详情页 → 需求管理Tab
2. 选择文档类型Tab(整体需求/整体测试观点/变更需求/变更测试观点)
3. 默认只读模式,点击"编辑"按钮进入编辑模式
4. 编辑文档内容:
   - 使用工具栏快捷操作
   - 或直接输入Markdown语法
   - 支持导入现有.md文件
5. 保存操作:
   - "保存"按钮:仅更新数据库,不生成版本
   - "版本保存"按钮:更新数据库+生成版本文件+创建版本记录
6. 切换到版本管理Tab查看历史版本:
   - 整体版本管理:左栏显示整体需求版本,右栏显示整体测试观点版本
   - 变更版本管理:左栏显示变更需求版本,右栏显示变更测试观点版本
7. 版本操作:
   - 下载:下载指定版本的.md文件
   - 删除:软删除版本记录(不影响已生成的文件)

### 4. 自动化测试用例库 (T33)
- ✅ ROLE1-4 四类自动化测试用例管理
  - 多语言支持(中文/日文/英文)
  - 可编辑表格(内联编辑)
  - 拖拽排序
  - 批量删除
- ✅ 版本管理功能
  - 一键保存版本:批量导出ROLE1-4的用例为Excel(19列全语言)
  - 版本列表展示:显示历史版本，包含文件信息和备注
  - 版本下载:一键下载zip压缩包(包含4个Excel文件)
  - 备注编辑:内联编辑版本备注(≤200字符)
  - 版本删除:删除物理文件和数据库记录
- ✅ Excel格式优化
  - 19列:ID, CaseNumber, Screen/Function/Precondition/TestSteps/ExpectedResult (CN/JP/EN), TestResult, Remark
  - 样式美化:表头深蓝背景+粗体白字,自动列宽,文本换行
  - 并发导出:4个ROLE并发处理,提升性能

**使用指南:**
1. 进入项目详情页 → 自动化测试用例Tab
2. 切换到ROLE1/ROLE2/ROLE3/ROLE4 Tab管理对应用例
3. 保存版本:
   - 点击ROLE1 Tab顶部"保存版本"按钮
   - 系统并发导出4个ROLE的用例为Excel
   - 自动切换到"版本管理"Tab
4. 查看版本历史:
   - 版本ID:格式为`{项目名}_{YYYYMMDD_HHMMSS}`
   - 版本文件:显示4个Excel文件及文件大小、用例数量
   - 备注:点击编辑区域可添加/修改备注
5. 版本操作:
   - 下载:下载zip压缩包(包含4个Excel文件)
   - 删除:删除版本记录和物理文件(需二次确认)

### 5. 接口测试用例管理 (T14)
- ✅ ROLE1-4 四类接口测试用例管理
  - 可编辑表格(内联编辑)
  - 插入行(在上方/在下方)
  - 批量删除
  - UUID主键+display_order排序
- ✅ 表格字段
  - No.(显示序号,基于display_order计算)
  - 用例编号、画面、URL、Header、Method、Body、Response、测试结果、备考
  - Method下拉选择:GET/POST/PUT/DELETE/PATCH
  - 测试结果下拉选择:NR/OK/NG
- ✅ 版本管理功能
  - 一键保存版本:导出ROLE1-4的用例为CSV(4个文件)
  - 版本列表展示:显示历史版本,4个文件名分4行显示
  - 版本下载:一键下载ZIP压缩包(包含4个CSV文件)
  - 备注编辑:更新版本备注(≤500字符)
  - 版本删除:删除版本记录和CSV文件

**使用指南:**
1. 进入项目详情页 → 接口测试用例Tab
2. 切换到ROLE1/ROLE2/ROLE3/ROLE4 Tab管理对应用例
3. 表格操作:
   - 内联编辑:单击单元格直接编辑
   - 插入行:点击操作列"在上方插入"或"在下方插入"
   - 批量删除:勾选多行后点击顶部"批量删除"按钮
4. 保存版本:
   - 点击顶部"保存版本"按钮
   - 系统生成4个CSV文件(格式:{项目名}_APITestCase_ROLE1_YYYYMMDD_HHMMSS.csv)
   - 自动切换到"版本管理"Tab
5. 查看版本历史:
   - 版本ID:UUID格式
   - 版本文件:4个文件名分4行显示,每行显示文件名
   - 备注:点击编辑按钮修改备注
6. 版本操作:
   - 下载:下载ZIP压缩包(包含4个CSV文件)
   - 删除:删除版本记录和CSV文件(需二次确认)

**技术特性:**
- 数据库主键使用UUID(VARCHAR 36),确保全局唯一性
- display_order字段控制显示顺序,支持插入行时自动调整
- 插入行后自动重新分配display_order为连续序号(1,2,3...)
- CSV导出时No.列基于display_order计算显示序号

### 6. 国际化
- ✅ 中英文动态切换
- ✅ 所有用户界面文本支持翻译

### 7. 错误处理
- ✅ 统一错误响应格式
- ✅ Axios 拦截器统一处理 401/403/500
- ✅ 用户友好的错误提示

## API 文档

详见 [API-documentation.md](./docs/API-documentation.md)

### 主要接口

#### 登录
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}

响应:
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGc..."
  }
}
```

#### 获取用户信息
```
GET /api/v1/profile
Authorization: Bearer <token>

响应:
{
  "username": "admin",
  "message": "authenticated user profile"
}
```

## 部署

### 开发环境

已配置热重载,修改代码后自动重启。

### 生产环境

详见:
- [后端部署文档](./docs/backend-deployment.md)
- [前端部署文档](./docs/frontend-deployment.md)

**关键步骤**:
1. 修改 `JWT_SECRET` 环境变量
2. 修改管理员默认密码
3. 配置 HTTPS
4. 配置反向代理 (Nginx)
5. 设置数据库备份

## 测试

### 后端测试
```bash
cd backend
go test ./internal/services/... -v
go test ./... -coverprofile=coverage.out
```

### 前端测试
```bash
cd frontend
npm test
npm test -- --coverage
```

### 手动测试
详见 [手动测试指南](./docs/manual-test-guide.md)

## 安全性

### 已实现
- ✅ bcrypt 密码加密 (cost=10)
- ✅ JWT Token 认证 (HS256)
- ✅ 登录失败不泄露用户存在性
- ✅ 参数验证防止注入攻击
- ✅ CORS 支持配置白名单

### 生产环境建议
- ⚠️ 修改默认 JWT_SECRET
- ⚠️ 启用 HTTPS
- ⚠️ 配置 Rate Limiting 防暴力破解
- ⚠️ 添加验证码
- ⚠️ 实现 Refresh Token 机制

## 开发规范

### 后端 (Go)
- 错误透明传递: `fmt.Errorf("msg: %w", err)`
- 函数长度 ≤ 40 行
- 文件长度 ≤ 500 行
- 所有 struct 添加 json tag

### 前端 (React)
- 组件长度 ≤ 200 行
- 使用函数组件 + Hooks
- PascalCase 命名组件
- useEffect 集中管理副作用

## 常见问题

### Q: Go 依赖下载失败?
A: 配置国内代理:
```bash
export GOPROXY=https://goproxy.cn,direct
```

### Q: 前端启动后无法访问后端?
A: 检查 `.env` 文件中 `REACT_APP_API_BASE_URL` 是否正确

### Q: Token 过期后如何处理?
A: 当前需要重新登录,未来版本将实现 Refresh Token

### Q: 如何修改默认密码?
A: 登录后调用修改密码 API (待实现)

## 下一步计划

### 短期 (1-2 周)
- [ ] 实现用户管理功能 (CRUD)
- [ ] 添加角色权限管理
- [ ] 实现 Refresh Token
- [ ] 添加登录验证码

### 中期 (1-2 月)
- [ ] 集成 Swagger API 文档
- [ ] 添加操作日志记录
- [ ] 性能监控 (Prometheus)
- [ ] E2E 测试 (Cypress)

### 长期 (3-6 月)
- [ ] 单点登录 (SSO)
- [ ] OAuth2 第三方登录
- [ ] 多租户支持
- [ ] 微服务架构拆分

## 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

## 许可证

本项目采用 MIT 许可证

## 联系方式

- 项目负责人: [Your Name]
- 邮箱: support@example.com
- 问题反馈: https://github.com/yourorg/webtest/issues

---

**最后更新**: 2025-11-01  
**版本**: v1.0.0
