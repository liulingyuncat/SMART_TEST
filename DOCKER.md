# 🐳 Docker 使用指南

## 📦 获取镜像

### 从 GitHub Container Registry 拉取

```bash
# 拉取最新版本
docker pull ghcr.io/liulingyuncat/smart_test:latest

# 拉取特定版本
docker pull ghcr.io/liulingyuncat/smart_test:v1.0.0

# 拉取特定提交
docker pull ghcr.io/liulingyuncat/smart_test:sha-abc1234
```

### 认证（私有仓库）

如果仓库是私有的，需要先登录：

```bash
# 使用 GitHub Personal Access Token (PAT)
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# 或使用交互式登录
docker login ghcr.io
```

## 🚀 快速启动

### 使用 Docker Run

```bash
docker run -d \
  --name smart-test \
  -p 8443:8443 \
  -p 16410:16410 \
  -v smart-test-data:/app/storage \
  ghcr.io/liulingyuncat/smart_test:latest
```

### 使用 Docker Compose

```yaml
version: '3.8'

services:
  smart-test:
    image: ghcr.io/liulingyuncat/smart_test:latest
    container_name: smart-test
    ports:
      - "8443:8443"
      - "16410:16410"
    volumes:
      - smart-test-data:/app/storage
      - ./certs:/app/certs
    environment:
      - TZ=Asia/Shanghai
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8443/api/v1/health"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  smart-test-data:
```

启动服务：

```bash
docker-compose up -d
```

## 🔧 环境变量配置

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `TZ` | 时区 | `Asia/Shanghai` |
| `PORT` | HTTP 端口 | `8443` |
| `MCP_PORT` | MCP 服务端口 | `16410` |

## 📂 数据持久化

### 重要目录

- `/app/storage` - 存储上传的文件和数据库
- `/app/certs` - SSL 证书（首次启动自动生成）
- `/app/config` - 配置文件

### 挂载示例

```bash
docker run -d \
  --name smart-test \
  -p 8443:8443 \
  -p 16410:16410 \
  -v $(pwd)/storage:/app/storage \
  -v $(pwd)/certs:/app/certs \
  -v $(pwd)/config:/app/config \
  ghcr.io/liulingyuncat/smart_test:latest
```

## 🏗️ 本地构建镜像

```bash
# 构建镜像
docker build -t smart-test:local .

# 多架构构建
docker buildx build --platform linux/amd64,linux/arm64 -t smart-test:local .
```

## 📊 CI/CD 工作流

### 自动构建触发条件

1. **推送到 main 分支**
   - 自动运行测试
   - 测试通过后构建并推送镜像
   - 标签: `latest`, `sha-<commit>`

2. **创建版本标签** (如 `v1.0.0`)
   - 自动构建并推送版本镜像
   - 标签: `v1.0.0`, `1.0.0`, `1.0`, `1`

3. **手动触发**
   - 通过 GitHub Actions 手动触发
   - 可自定义标签

### 查看构建状态

访问项目的 GitHub Actions 页面：
```
https://github.com/liulingyuncat/SMART_TEST/actions
```

### 查看已发布的镜像

访问项目的 Packages 页面：
```
https://github.com/liulingyuncat/SMART_TEST/pkgs/container/smart_test
```

## 🔍 健康检查

```bash
# 检查容器健康状态
docker ps

# 查看健康检查日志
docker inspect --format='{{json .State.Health}}' smart-test

# 手动健康检查
curl -f http://localhost:8443/api/v1/health || echo "Service unhealthy"
```

## 📝 日志查看

```bash
# 查看实时日志
docker logs -f smart-test

# 查看最近 100 行日志
docker logs --tail 100 smart-test

# 查看带时间戳的日志
docker logs -t smart-test
```

## 🛠️ 故障排查

### 容器无法启动

```bash
# 查看详细错误信息
docker logs smart-test

# 进入容器调试
docker exec -it smart-test sh
```

### 端口冲突

```bash
# 检查端口占用
netstat -tuln | grep 8443

# 使用其他端口
docker run -d -p 9443:8443 ghcr.io/liulingyuncat/smart_test:latest
```

### 权限问题

容器使用非 root 用户 `webtest` (UID=1000, GID=1000)。确保挂载的目录有正确权限：

```bash
# 设置目录权限
chown -R 1000:1000 ./storage ./certs ./config
```

## 🔐 安全建议

1. **使用特定版本标签**，避免使用 `latest`
2. **定期更新镜像**，获取安全补丁
3. **使用 Docker secrets** 管理敏感信息
4. **限制容器资源**：

```bash
docker run -d \
  --name smart-test \
  --memory="512m" \
  --cpus="1.0" \
  -p 8443:8443 \
  ghcr.io/liulingyuncat/smart_test:latest
```

## 📚 更多资源

- [Dockerfile 源码](./Dockerfile)
- [Docker Compose 配置](./docker-compose.yml)
- [CI/CD 工作流](./.github/workflows/ci.yml)
- [项目文档](./README.md)

## 🆘 获取帮助

遇到问题？欢迎提交 Issue：
```
https://github.com/liulingyuncat/SMART_TEST/issues
```
