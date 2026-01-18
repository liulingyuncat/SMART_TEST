# Playwright Executor Service

Playwright 脚本执行服务，用于执行用户提交的 Playwright 测试脚本。

## 架构说明

**为什么需要这个服务？**

根据 [Playwright 官方文档](https://playwright.dev/docs/docker)，`playwright run-server` 仅提供远程浏览器实例，**不支持直接接收 JS 脚本执行**。客户端必须通过 `browserType.connect()` 连接后，在客户端执行 Playwright API 调用。

因此，如果要执行 JavaScript Playwright 脚本，必须有一个 Node.js 运行时来：
1. 解析并执行用户的 JavaScript 脚本
2. 调用 Playwright API 控制浏览器
3. 捕获执行结果或错误

**架构图：**

```
┌─────────────────┐     HTTP      ┌─────────────────┐    WebSocket    ┌─────────────────┐
│  Go Backend     │ ───────────► │  Playwright     │ ──────────────► │  Playwright     │
│  (webtest)      │    :53730    │  Executor       │     :53729      │  Runner         │
│                 │              │  (Node.js)      │                 │  (run-server)   │
└─────────────────┘              └─────────────────┘                 └─────────────────┘
```

## 端口分配

- **53729**: Playwright Runner (run-server WebSocket)
- **53730**: Playwright Executor (HTTP API)

## 本地开发

### 1. 启动开发依赖（一次性）

```bash
make dev-deps
```

这会启动 Docker 容器：
- PostgreSQL (localhost:5432)
- Playwright Runner (ws://localhost:53729/)
- Playwright Executor (http://localhost:53730/)

### 2. 启动后端服务

```bash
make dev
```

这会启动：
- **Backend Server** - Go 主服务
- **MCP Server** - Model Context Protocol 服务

按 `Ctrl+C` 退出。

### 3. 停止开发依赖

```bash
make dev-deps-stop
```

### 4. 查看依赖服务日志

```bash
make dev-deps-logs
```

## 环境变量

### Playwright Executor

- `PLAYWRIGHT_WS`: Playwright Server WebSocket 地址（默认: `ws://playwright-runner:53729/`）
- `PORT`: 服务端口（默认: `53730`）

### Backend Server

- `PLAYWRIGHT_EXECUTOR_URL`: Executor 服务地址（默认: `http://playwright-executor:53730`）

## 测试脚本格式

脚本格式为带有 `page` 参数的异步函数：

```javascript
async (page) => {
  await page.goto('http://localhost:3002/mail/view');
  await page.getByRole('textbox', { name: '用户名' }).fill('neteye');
  await page.getByRole('textbox', { name: '密码' }).fill('neteye@123');
  await page.getByRole('button', { name: '登 录' }).click();
  await page.waitForURL('**/mail/view**');
  await expect(page.getByText('neteye')).toBeVisible();
}
```

## API

### POST /execute

执行 Playwright 脚本。

**请求：**

```json
{
  "scriptCode": "async (page) => { await page.goto('https://example.com'); }",
  "timeout": 60000
}
```

**响应（成功）：**

```json
{
  "success": true,
  "output": "Script executed successfully",
  "responseTime": 1234
}
```

**响应（失败）：**

```json
{
  "success": false,
  "error": "page.goto: net::ERR_CONNECTION_REFUSED at http://localhost:3002/mail/view",
  "stack": "Error: page.goto: net::ERR_CONNECTION_REFUSED...",
  "responseTime": 5678
}
```

## 问题排查

### 1. "Cannot connect to Playwright Server"

确保 Playwright Runner 正在运行：

```bash
docker ps | grep playwright-runner
```

### 2. "Script execution failed"

查看 Executor 日志：

```bash
# 如果使用 Docker
docker logs playwright-executor

# 如果使用 make dev
# 日志会直接输出到终端
```

### 3. 端口冲突

修改端口：

```bash
# Executor
PORT=3002 make executor

# Server  
# 修改 backend/cmd/server/main.go 中的端口配置
```

## Docker 部署

参见 [docker-compose.yml](../docker-compose.yml)。

## 更多命令

查看所有可用命令：

```bash
make help
```
