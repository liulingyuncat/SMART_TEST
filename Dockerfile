# =============================================================================
# SMART_TEST (webtest) - Multi-stage Dockerfile
# 
# 三阶段构建策略：
#   Stage 1: 前端构建 (Node.js)
#   Stage 2: 后端构建 (Go)
#   Stage 3: 生产镜像 (Alpine)
#
# 镜像: ghcr.io/liulingyuncat/smart_test:latest
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: Frontend Builder
# -----------------------------------------------------------------------------
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# 配置 npm 提高网络稳定性
RUN npm config set fetch-retry-maxtimeout 120000 && \
    npm config set fetch-retry-mintimeout 15000 && \
    npm config set fetch-retries 5 && \
    npm config set fetch-timeout 60000

# 复制依赖文件并安装
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci --production=false --legacy-peer-deps || \
    (echo "First attempt failed, retrying..." && npm ci --production=false --legacy-peer-deps)

# 复制源码并构建
COPY frontend/ ./

# 构建参数：版本号（从 .env 读取，或使用默认值）
ARG VERSION=0.0.1
ENV REACT_APP_VERSION=${VERSION}

# ARM64 构建时降低内存使用，避免 QEMU 崩溃
ARG TARGETARCH
ENV NODE_OPTIONS="--max-old-space-size=3072"
RUN if [ "$TARGETARCH" = "arm64" ]; then \
        # ARM64: 降低并发，减少内存压力
        export NODE_OPTIONS="--max-old-space-size=2048" && \
        npm run build -- --max-workers=2; \
    else \
        # AMD64: 正常构建
        npm run build; \
    fi

# -----------------------------------------------------------------------------
# Stage 2: Backend Builder
# -----------------------------------------------------------------------------
FROM golang:1.24-alpine AS backend-builder

# 安装构建依赖
RUN apk add --no-cache git

WORKDIR /app/backend

# 设置 Go 环境 - 支持多架构
ARG TARGETARCH
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=${TARGETARCH}

# 复制依赖文件并下载
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# 复制源码
COPY backend/ ./

# 构建 Web 服务
RUN go build -ldflags="-s -w" -o server ./cmd/server

# 构建 MCP 服务
RUN go build -ldflags="-s -w" -o mcp-server ./cmd/mcp

# -----------------------------------------------------------------------------
# Stage 3: Production Image
# -----------------------------------------------------------------------------
FROM alpine:3.19

# 安装运行时依赖
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    openssl \
    wget \
    && rm -rf /var/cache/apk/*

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非 root 用户
RUN addgroup -g 1000 webtest && \
    adduser -D -u 1000 -G webtest webtest

WORKDIR /app

# 从构建阶段复制产物
COPY --from=backend-builder /app/backend/server ./server
COPY --from=backend-builder /app/backend/mcp-server ./mcp-server
COPY --from=frontend-builder /app/frontend/build ./frontend/build

# 复制配置文件
COPY backend/config/mcp-server.yaml ./config/mcp-server.yaml

# 复制系统提示词
COPY backend/internal/mcp/prompts ./internal/mcp/prompts

# 创建证书目录（证书将在启动时自动生成）
RUN mkdir -p ./certs

# 复制启动脚本
COPY docker-entrypoint.sh ./docker-entrypoint.sh
RUN chmod +x ./docker-entrypoint.sh

# 创建存储目录并设置权限
RUN mkdir -p ./storage && \
    chown -R webtest:webtest /app

# 切换到非 root 用户
USER webtest

# 暴露端口
EXPOSE 8443 16410

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8443/api/v1/health || exit 1

# 启动命令
CMD ["./docker-entrypoint.sh"]
