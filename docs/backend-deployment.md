# Web 智能测试平台 - 后端部署文档

## 环境要求

- Go 1.21 或更高版本
- PostgreSQL 13+ / SQLite 3+ / MongoDB 4.4+

## 安装依赖

```bash
cd backend
go mod tidy
```

## 环境变量配置

创建 `.env` 文件:
```bash
# 数据库配置
DB_TYPE=sqlite          # postgres, sqlite, mongodb
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=webtest

# JWT 密钥(生产环境必须修改)
JWT_SECRET=your-secret-key-change-in-production

# 服务器配置
SERVER_PORT=8080
```

## 数据库初始化

### PostgreSQL
```bash
psql -U postgres -c "CREATE DATABASE webtest;"
psql -U postgres -d webtest -f migrations/001_create_users_table.sql
```

### MongoDB
```bash
mongosh webtest < migrations/mongodb_schema.js
```

### SQLite
应用启动时自动创建数据库文件 `webtest.db`

## 运行应用

### 开发环境
```bash
go run cmd/server/main.go
```

### 生产环境编译
```bash
go build -o webtest-server cmd/server/main.go
./webtest-server
```

## 测试

运行单元测试:
```bash
go test ./internal/services/... -v
```

运行所有测试:
```bash
go test ./... -v
```

生成覆盖率报告:
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Docker 部署

### Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

### 构建和运行
```bash
docker build -t webtest-backend .
docker run -p 8080:8080 \
  -e DB_TYPE=postgres \
  -e DB_HOST=host.docker.internal \
  -e JWT_SECRET=your-secret \
  webtest-backend
```

## Systemd 服务配置

创建 `/etc/systemd/system/webtest.service`:
```ini
[Unit]
Description=Web Test Platform Backend
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/webtest
ExecStart=/opt/webtest/webtest-server
Restart=on-failure
Environment="JWT_SECRET=your-secret"
Environment="DB_TYPE=postgres"

[Install]
WantedBy=multi-user.target
```

启动服务:
```bash
sudo systemctl daemon-reload
sudo systemctl enable webtest
sudo systemctl start webtest
```

## 默认管理员账号

系统启动时自动创建:
- 用户名: `admin`, 密码: `admin123`
- 用户名: `root`, 密码: `root123`

**生产环境部署后请立即修改默认密码!**

## API 测试

```bash
# 登录测试
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 认证测试
curl http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 性能优化

1. 数据库连接池: MaxOpenConns=25, MaxIdleConns=10
2. GORM 日志级别根据环境调整
3. 使用索引优化查询性能
4. 考虑使用 Redis 缓存 Token

## 监控和日志

建议集成:
- Prometheus + Grafana 监控
- ELK Stack 日志收集
- Sentry 错误追踪
