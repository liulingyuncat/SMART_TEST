# SMART_TEST 部署指南

本文档介绍如何使用 Docker 部署 SMART_TEST (webtest) 智能测试平台。

## 目录

- [系统要求](#系统要求)
- [快速开始](#快速开始)
- [配置说明](#配置说明)
- [常见操作](#常见操作)
- [故障排除](#故障排除)
- [开发环境](#开发环境)

---

## 系统要求

| 组件 | 最低版本 | 推荐版本 |
|------|----------|----------|
| Docker | 20.10+ | 24.0+ |
| Docker Compose | 2.0+ | 2.20+ |
| 内存 | 2GB | 4GB+ |
| 磁盘空间 | 10GB | 20GB+ |

---

## 快速开始

### 1. 克隆仓库

```bash
git clone https://github.com/liulingyuncat/SMART_TEST.git
cd SMART_TEST
```

### 2. 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑配置（必须修改密码和密钥）
nano .env  # 或使用其他编辑器
```

**必须修改的配置项：**
- `DB_PASSWORD` - 数据库密码
- `JWT_SECRET` - JWT 签名密钥
- `MCP_AUTH_TOKEN` - MCP 认证令牌
- `POSTGRES_PASSWORD` - PostgreSQL 密码（需与 DB_PASSWORD 一致）

### 3. 启动服务

```bash
# 启动所有服务（后台运行）
docker-compose up -d

# 查看服务状态
docker-compose ps
```

### 4. 验证部署

- **Web 界面**: https://localhost:8443
- **MCP 服务**: http://localhost:16410

首次访问会看到自签名证书警告，可以选择继续访问。

---

## 配置说明

### 环境变量列表

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `DB_TYPE` | 数据库类型 | `postgres` |
| `DB_HOST` | 数据库主机 | `postgres` |
| `DB_PORT` | 数据库端口 | `5432` |
| `DB_USER` | 数据库用户 | `webtest` |
| `DB_PASSWORD` | 数据库密码 | (必填) |
| `DB_NAME` | 数据库名 | `webtest` |
| `DB_MAX_OPEN_CONNS` | 最大连接数 | `25` |
| `DB_MAX_IDLE_CONNS` | 最大空闲连接 | `10` |
| `JWT_SECRET` | JWT 签名密钥 | (必填) |
| `MCP_AUTH_TOKEN` | MCP 认证令牌 | (必填) |
| `TZ` | 时区 | `Asia/Shanghai` |

### 端口说明

| 端口 | 服务 | 协议 |
|------|------|------|
| 8443 | Web 服务 | HTTPS |
| 16410 | MCP 服务 | HTTP |
| 5432 | PostgreSQL | TCP (内部) |

---

## 常见操作

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f webtest
docker-compose logs -f postgres
```

### 更新镜像

```bash
# 拉取最新镜像
docker-compose pull

# 重启服务
docker-compose up -d
```

### 备份数据库

```bash
# 备份 PostgreSQL 数据
docker exec webtest-postgres pg_dump -U webtest webtest > backup_$(date +%Y%m%d).sql

# 恢复数据
docker exec -i webtest-postgres psql -U webtest webtest < backup_20260110.sql
```

### 停止服务

```bash
# 停止服务（保留数据）
docker-compose down

# 停止并删除数据卷（⚠️ 会丢失所有数据）
docker-compose down -v
```

---

## 故障排除

### 服务无法启动

1. 检查端口占用：
   ```bash
   netstat -tlnp | grep -E '8443|16410|5432'
   ```

2. 检查日志：
   ```bash
   docker-compose logs webtest
   ```

3. 检查数据库连接：
   ```bash
   docker exec webtest-postgres pg_isready -U webtest
   ```

### 数据库连接失败

1. 确认 PostgreSQL 已启动：
   ```bash
   docker-compose ps postgres
   ```

2. 检查密码配置：
   - 确保 `.env` 中 `DB_PASSWORD` 和 `POSTGRES_PASSWORD` 一致

### 健康检查失败

检查健康检查端点：
```bash
curl -k https://localhost:8443/api/v1/health
```

预期响应：
```json
{"status": "ok", "db": "connected"}
```

---

## 开发环境

### 使用开发数据库

如果只需要 PostgreSQL 用于本地开发：

```bash
# 启动开发数据库
docker-compose -f docker-compose.dev.yml up -d

# 连接信息
# Host: localhost
# Port: 5432
# User: webtest
# Password: webtest_dev_password (或 .env 中配置)
# Database: webtest
```

### 本地开发配置

在 `.env` 中设置：

```bash
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=webtest
DB_PASSWORD=webtest_dev_password
DB_NAME=webtest
```

然后启动本地服务：

```bash
cd backend
go run ./cmd/server
```

---

## 架构说明

```
┌─────────────────────────────────────────┐
│           Docker Compose                │
│  ┌─────────────────────────────────┐   │
│  │       webtest Container         │   │
│  │  ┌───────────┐ ┌─────────────┐  │   │
│  │  │Web Service│ │ MCP Service │  │   │
│  │  │   :8443   │ │   :16410    │  │   │
│  │  └───────────┘ └─────────────┘  │   │
│  └─────────────────────────────────┘   │
│                  │                      │
│  ┌───────────────▼─────────────────┐   │
│  │     PostgreSQL Container        │   │
│  │           :5432                 │   │
│  └─────────────────────────────────┘   │
└─────────────────────────────────────────┘
```

---

## 更多信息

- **GitHub 仓库**: https://github.com/liulingyuncat/SMART_TEST
- **镜像仓库**: ghcr.io/liulingyuncat/smart_test
