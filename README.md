# SMART TEST 智能测试平台

[![CI/CD](https://github.com/liulingyuncat/SMART_TEST/actions/workflows/ci.yml/badge.svg)](https://github.com/liulingyuncat/SMART_TEST/actions/workflows/ci.yml)
[![Docker Build](https://github.com/liulingyuncat/SMART_TEST/actions/workflows/docker-build.yml/badge.svg)](https://github.com/liulingyuncat/SMART_TEST/actions/workflows/docker-build.yml)
[![GitHub Container Registry](https://img.shields.io/badge/ghcr.io-smart__test-blue?logo=docker)](https://github.com/liulingyuncat/SMART_TEST/pkgs/container/smart_test)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![Node Version](https://img.shields.io/badge/Node-20-339933?logo=node.js)](https://nodejs.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 简介

现代化的智能测试管理平台，采用 Go + React 前后端分离架构，支持需求管理、测试用例管理、缺陷跟踪、AI 辅助测试和质量报告生成。

## 🚀 快速开始

### 使用 Docker 部署（推荐）

```bash
# 1. 克隆项目
git clone https://github.com/liulingyuncat/SMART_TEST.git
cd SMART_TEST

# 2. 运行安装脚本（生成证书和配置）
./install.sh

# 3. 启动服务
docker-compose up -d

# 4. 访问应用
# 前端: https://localhost:8443
# 后端: https://localhost:8443/api
```

**默认管理员账号**: `admin` / `admin123`

### 生产环境部署

```bash
# 1. 运行安装脚本（自动生成证书和随机密钥）
./install.sh

# 2. 或手动编辑 .env 文件
DB_PASSWORD=your_secure_password
JWT_SECRET=your_jwt_secret_min_32_chars
MCP_AUTH_TOKEN=your_mcp_token

# 3. 使用 GitHub Container Registry（支持 AMD64/ARM64）
docker pull ghcr.io/liulingyuncat/smart_test:latest
docker-compose up -d
```

> 📌 **架构支持**: 支持 AMD64 (Intel/AMD x86_64)、ARM64 (Apple Silicon、ARM 服务器)

> 📖 详细文档: [Docker 使用指南](./DOCKER.md) | [GitHub Packages 部署](./GITHUB_PACKAGES_DEPLOYMENT.md)

## 技术栈

### 后端

- **语言**: Go 1.21+
- **Web框架**: Gin v1.9.1
- **ORM**: GORM v1.25.5
**后端**: Go 1.21+ · Gin · GORM · JWT · SQLite  
**前端**: React 18 · Ant Design 5 · Redux Toolkit · Axios  
**部署**: Docker · Docker Compose · HTTPS  
**AI**: MCP (Model Context Protocol) Server
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

# 3本地开发

### 前置要求
- Go 1.21+
- Node.js 20+
- OpenSSL (生成证书)

### 后端开发
```bash
cd backend
export GOPROXY=https://goproxy.cn,direct  # 国内用户
go mod tidy
go run cmd/server/main.go  # 启动后端服务
```

### 前端开发
```bash
cd frontend
npm install
npm start  # 启动开发服务器
``` user.go         # User 模型
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
│   ├── .env                    # Go 后端
│   ├── cmd/                    # 主程序入口
│   │   ├── server/             # Web 服务器
│   │   ├── mcp/                # MCP 服务器
│   │   └── gencert/            # 证书生成工具
│   ├── internal/               # 内部包
│   │   ├── handlers/           # HTTP 处理器
│   │   ├── middleware/         # 中间件
│   │   ├── models/             # 数据模型
│   │   ├── services/           # 业务逻辑
│   │   └── mcp/                # MCP 协议实现
│   ├── migrations/             # 数据库迁移
│   └── config/                 # 配置文件
│
├── frontend/                   # React 前端
│   ├── src/
│   │   ├── api/                # API 客户端
│   │   ├── pages/              # 页面组件
│   │   ├── components/         # 公共组件
│   │   ├── store/              # Redux 状态
│   │   └── locales/            # 国际化资源
│   └── public/                 # 静态资源
│
├── install.sh                  # 部署前安装脚本
├── docker-compose.yml          # Docker 编排
├── Dockerfile                  # Docker 镜像
└── .env.example                # 环境变量模板
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
用户认证
- JWT Token 认证
- bcrypt 密码加密
- 角色权限管理
- 前后端路由守卫

### 需求管理
   - 下载:下载zip压缩包(包含4个Excel文件)
   - 删除:删除版本记录和物理文件(需二次确认)

### 5. 接口测试用例管理 (T14)

- ✅ ROLE1-4 四类接口测试用例管理
  - 可编辑表格(内联编辑)
  - 插入行(在上方/在下方)
  四种需求文档类型（整体需求/测试观点、变更需求/测试观点）
- Markdown 编辑器（实时预览、工具栏）
- 版本管理（自动命名、文件存储、下载）
- Markdown 文件导入（≤5MB）
- 三语言支持（中/英/日）

### 测试用例管理
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
#### 手工测试用例
- ROLE1-4 分类管理
- 多语言用例（中/英/日）
- 内联编辑、拖拽排序、批量操作
- Excel 导出（19列全语言）
- 版本管理（ZIP 打包、备注编辑）

#### Web 自动化用例
- Playwright 用例管理
- 版本保存和回退
- 执行结果记录

#### 接口测试用例
3. 配置 HTTPS
4. 配置反向代理 (Nginx)
5. 设置数据库备份

## 测试

### 后端测试

``ROLE1-4 分类管理
- 表格编辑（内联编辑、插入行、批量删除）
- 字段：用例编号、URL、Method、Header、Body、Response 等
- CSV 导出和版本管理
- UUID 主键 + display_order 排序

### 缺陷管理
- 缺陷生命周期管理
- 附件上传和评论
- 缺陷统计和趋势图表
- 导出 Excel 报告

### 测试执行
- 任务创建和分配
- 用例筛选和批量执行
- 实时进度跟踪
- 燃尽图和进度统计

### AI 质量报告
- Markdown 报告编辑
- SVG/Recharts 图表集成
- PDF/HTML 导出
- 模板管理

### MCP 协议支持
- AI 辅助用例生成
- 智能需求分析
- 自动化脚本生成
- 15+ 工具集成
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
- [ ] 实示例

### 用户登录
```bash
curl -X POST https://localhost:8443/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 获取项目列表
```bash
curl https://localhost:8443/api/v1/projects \
  -H "Authorization: Bearer YOUR_TOKEN"*最后更新**: 2025-11-01
**版本**: v1.0.0
说明

### install.sh 脚本
部署前运行 `./install.sh` 自动完成：
- 创建必要目录（storage、certs、tmp）
- 生成自签名 HTTPS 证书
- 创建 .env 文件并生成随机密钥
- 设置文件权限

### Docker Compose
```yaml
services:
  webtest:
    image: ghcr.io/liulingyuncat/smart_test:latest
    ports:
      - "8443:8443"
    volumes:
      - ./backend/storage:/app/storage
      - ./backend/certs:/app/certs
    environment:
      - JWT_SECRET=${JWT_SECRET}
      - DB_PASSWORD=${DB_PASSWORD}
```

### 生产环境检查清单
- [ ] 运行 `./install.sh` 生成证书和配置
- [ ] 修改 `.env` 中的默认密码
- [ ] 配置反向代理（Nginx）
- [ ] 设置数据库备份计划
- [ ] 修改管理员默认密码```bash
**后端**: 错误透明传递 · 函数 ≤40 行 · 文件 ≤500 行  
**前端**: 函数组件 + Hooks · 组件 ≤200 行 · PascalCase 命名
- bcrypt 密码加密（cost=10）
- JWT Token 认证（HS256）
- HTTPS 自签名证书
- CORS 白名单配置
- SQL 注入防护
- XSS 过滤**Q: Go 依赖下载失败？**  
A: 配置代理 `export GOPROXY=https://goproxy.cn,direct`

**Q: HTTPS 证书错误？**  
A: 运行 `./install.sh` 重新生成证书

**Q: Docker 容器无法启动？**  
A: 检查 `.env` 文件和端口占用（8443）

**Q: 如何修改默认密码？**  
A: 登录后在个人中心修改开发路线

- [ ] Refresh Token 机制
- [ ] 登录验证码
- [ ] Swagger API 文档
- [ ] 操作日志审计
- [ ] 性能监控（Prometheus）
- [ ] SSO 单点登录
- [ ] 多租户支持

欢迎提交 Issue 和 Pull Request！

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/Feature`)
3. 提交更改 (`git commit -m 'Add Feature'`)
4. 推送分支 (`git push origin feature/Feature`)
5. 提交 Pull Request

## 许可证

MIT License

---

**最后更新**: 2026-01-11 | 