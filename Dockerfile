# =============================================================================
# SMART_TEST (webtest) - Multi-stage Dockerfile
# 
# 三阶段构建策略：
#   Stage 1: 前端构建 (Node.js Alpine)
#   Stage 2: 后端构建 (Go Debian) - 包含 playwright-go 驱动安装
#   Stage 3: 生产镜像 (Debian Slim) - 支持 Playwright WebSocket 连接
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
# 使用 Debian 版本的 Go 镜像，以便正确安装 playwright-go 驱动
# -----------------------------------------------------------------------------
FROM golang:1.24-bookworm AS backend-builder

# 安装构建依赖
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    && rm -rf /var/lib/apt/lists/*

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

# 安装 playwright-go 驱动（下载到 /root/.cache/ms-playwright-go/）
RUN go install github.com/playwright-community/playwright-go/cmd/playwright@v0.5200.1 && \
    /go/bin/playwright install

# -----------------------------------------------------------------------------
# Stage 3: Production Image
# 使用 Debian 而非 Alpine，因为 Playwright 不支持 musl libc
# -----------------------------------------------------------------------------
FROM debian:bookworm-slim

# 安装运行时依赖
# 注意：不再需要 docker-cli，使用 WebSocket 连接 playwright-runner
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    tzdata \
    openssl \
    wget \
    curl \
    bash \
    netcat-openbsd \
    && rm -rf /var/lib/apt/lists/*

# 设置时区
ENV TZ=Asia/Shanghai

WORKDIR /app

# 从构建阶段复制产物
COPY --from=backend-builder /app/backend/server ./server
COPY --from=backend-builder /app/backend/mcp-server ./mcp-server
COPY --from=frontend-builder /app/frontend/build ./frontend/build

# 复制 playwright-go 驱动文件（以 root 运行，直接使用 /root）
COPY --from=backend-builder /root/.cache/ms-playwright-go /root/.cache/ms-playwright-go

# 复制配置文件
COPY backend/config/mcp-server.yaml ./config/mcp-server.yaml

# 复制系统提示词
COPY backend/internal/mcp/prompts ./internal/mcp/prompts

# 创建证书目录（证书将在启动时自动生成）
RUN mkdir -p ./certs

# 复制启动脚本
COPY docker-entrypoint.sh ./docker-entrypoint.sh
RUN chmod +x ./docker-entrypoint.sh

# 创建存储目录
RUN mkdir -p ./storage

# 注意: 以 root 用户运行以避免挂载卷的权限问题
# 在生产环境中，通过 docker-compose 的 volume 挂载可能导致权限冲突
# 使用 root 可以确保应用始终能访问挂载的目录

# 暴露端口
EXPOSE 8443 16410

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8443/api/v1/health || exit 1

# 启动命令
CMD ["./docker-entrypoint.sh"]
