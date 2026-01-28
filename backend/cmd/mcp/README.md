#### PEVVD Intelligent Test Platform MCP Server

MCP (Model Context Protocol) 服务器，为AI助手提供访问PEVVD Intelligent Test Platform数据的能力。

## 快速开始

### 构建

```powershell
# 在 backend 目录下
go build -o build/mcp-server.exe ./cmd/mcp/...

# 或使用构建脚本
.\scripts\build-mcp.ps1 build
```

### 运行

```powershell
# 使用默认配置
.\build\mcp-server.exe

# 指定配置文件
.\build\mcp-server.exe --config ./config/mcp-server.yaml

# 查看版本
.\build\mcp-server.exe --version
```

## 配置

### 配置文件

默认配置文件路径: `./config/mcp-server.yaml`

```yaml
server:
  mode: stdio              # 传输模式: stdio | sse
  log_level: info          # 日志级别: debug | info | warn | error
  log_format: json         # 日志格式: json | text

backend:
  base_url: https://localhost:8443
  timeout: 30s
  retry_count: 3
  retry_delay: 1s

auth:
  token_env: MCP_AUTH_TOKEN    # 从环境变量读取Token
  token_file: ""               # 或从文件读取Token
  validate_on_start: true      # 启动时验证Token
```

### 环境变量

| 变量名              | 说明         | 默认值                |
| ------------------- | ------------ | --------------------- |
| `MCP_AUTH_TOKEN`  | JWT认证Token | (必需)                |
| `MCP_BACKEND_URL` | 后端API地址  | https://localhost:8443 |
| `MCP_SERVER_MODE` | 服务器模式   | stdio                 |
| `MCP_LOG_LEVEL`   | 日志级别     | info                  |

## 协议

### MCP版本

- 协议版本: `2025-06-18`
- 传输协议: JSON-RPC 2.0 over stdio

### 支持的方法

| 方法            | 说明                  |
| --------------- | --------------------- |
| `initialize`  | 初始化连接            |
| `initialized` | 确认初始化完成 (通知) |
| `tools/list`  | 列出可用工具          |
| `tools/call`  | 调用工具              |

## 可用工具 (30个)

### 原始文档 (2个)

- `list_raw_documents` - 获取项目原始文档列表
- `get_raw_document` - 获取原始文档详情

### 需求管理 (4个)

- `list_requirement_items` - 获取需求条目列表
- `get_requirement_item` - 获取需求条目详情
- `create_requirement_item` - 创建需求条目
- `update_requirement_item` - 更新需求条目

### 测试观点 (4个)

- `list_viewpoint_items` - 获取测试观点列表
- `get_viewpoint_item` - 获取测试观点详情
- `create_viewpoint_item` - 创建测试观点
- `update_viewpoint_item` - 更新测试观点

### 手工用例 (5个)

- `list_manual_groups` - 获取手工测试用例集列表
- `list_manual_cases` - 获取手工用例列表
- `create_case_group` - 创建用例集
- `create_manual_case` - 创建手工用例
- `update_manual_case` - 更新手工用例

### Web自动化用例 (3个)

- `list_web_cases` - 获取Web测试用例列表
- `create_web_case` - 创建Web测试用例
- `update_web_case` - 更新Web测试用例

### API自动化用例 (3个)

- `list_api_cases` - 获取API测试用例列表
- `create_api_case` - 创建API测试用例
- `update_api_case` - 更新API测试用例

### 执行任务 (3个)

- `list_execution_tasks` - 获取执行任务列表
- `get_execution_task_cases` - 获取任务关联用例
- `update_execution_case_result` - 更新执行结果

### 缺陷管理 (2个)

- `list_defects` - 获取缺陷列表
- `update_defect` - 更新缺陷

### AI报告 (2个)

- `create_ai_report` - 创建AI报告
- `update_ai_report` - 更新AI报告

## 开发

### 运行测试

```powershell
# 运行所有MCP测试
go test -v ./internal/mcp/...

# 使用构建脚本
.\scripts\build-mcp.ps1 test

# 带覆盖率
.\scripts\build-mcp.ps1 test-coverage
```

### 项目结构

```
backend/
├── cmd/mcp/
│   └── main.go              # 入口程序
├── config/
│   └── mcp-server.yaml      # 配置文件模板
├── internal/mcp/
│   ├── config/              # 配置管理
│   │   └── config.go
│   ├── transport/           # 传输层
│   │   ├── transport.go     # 接口定义
│   │   └── stdio.go         # stdio实现
│   ├── protocol/            # 协议层
│   │   ├── jsonrpc.go       # JSON-RPC解析
│   │   └── router.go        # 消息路由
│   ├── client/              # 后端客户端
│   │   ├── auth.go          # 认证管理
│   │   └── client.go        # HTTP客户端
│   ├── tools/               # 工具层
│   │   ├── handler.go       # 处理器接口
│   │   ├── registry.go      # 注册表
│   │   ├── validator.go     # 参数验证
│   │   └── handlers/        # 工具处理器实现
│   │       ├── register.go
│   │       ├── raw_document.go
│   │       ├── requirement.go
│   │       ├── viewpoint.go
│   │       ├── manual_case.go
│   │       ├── web_case.go
│   │       ├── api_case.go
│   │       ├── review.go
│   │       ├── execution.go
│   │       ├── defect.go
│   │       └── report.go
│   └── server.go            # 服务器主体
└── scripts/
    └── build-mcp.ps1        # 构建脚本
```

## 在VS Code中配置

在 `.vscode/mcp.json` 中添加:

```json
{
  "servers": {
    "webtest": {
      "type": "stdio",
      "command": "path/to/mcp-server.exe",
      "args": ["--config", "path/to/mcp-server.yaml"],
      "env": {
        "MCP_AUTH_TOKEN": "your-jwt-token"
      }
    }
  }
}
```

## 故障排除

### Token无效

确保设置了有效的 `MCP_AUTH_TOKEN` 环境变量。可以通过以下命令获取Token:

```powershell
# 登录获取Token
curl -X POST https://localhost:8443/api/v1/auth/login `
  -H "Content-Type: application/json" `
  -d '{"username":"admin","password":"password"}' `
  --insecure
```

### 连接后端失败

检查:

1. 后端服务是否运行 (`https://localhost:8443`)
2. `MCP_BACKEND_URL` 环境变量是否正确设置
3. 网络是否可达
4. HTTPS证书配置（开发环境使用自签名证书）

### 日志查看

MCP服务器将日志输出到 stderr，不会干扰 stdio 通信:

```powershell
# 将stderr重定向到文件
.\build\mcp-server.exe 2> mcp.log
```
