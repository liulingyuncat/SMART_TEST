#!/bin/bash
# =============================================================================
# SMART_TEST (webtest) - Docker Entrypoint Script
#
# 在单个容器内同时运行 Web 服务和 MCP 服务
# 支持 Docker 信号传递实现优雅停机
# =============================================================================

set -e

# 定义变量存储子进程 PID
SERVER_PID=""
MCP_PID=""

# 自动生成 TLS 证书（如果不存在）
generate_certs() {
    if [ ! -f "/app/certs/server.crt" ] || [ ! -f "/app/certs/server.key" ]; then
        echo "[entrypoint] Generating self-signed TLS certificates..."
        openssl req -x509 -newkey rsa:4096 -nodes \
            -keyout /app/certs/server.key \
            -out /app/certs/server.crt \
            -days 7300 \
            -subj "/C=CN/ST=Beijing/L=Beijing/O=SmartTest/CN=localhost"
        echo "[entrypoint] TLS certificates generated successfully"
    else
        echo "[entrypoint] Using existing TLS certificates"
    fi
}

# 检查数据库连接和必需的表
check_database() {
    local max_attempts=30
    local attempt=0
    
    echo "[entrypoint] Checking database connection..."
    
    # 等待数据库就绪
    while [ $attempt -lt $max_attempts ]; do
        if [ "${DB_TYPE:-sqlite}" = "postgres" ]; then
            # PostgreSQL 连接检查
            if nc -z "${DB_HOST:-postgres}" "${DB_PORT:-5432}" 2>/dev/null; then
                echo "[entrypoint] Database connection established"
                break
            fi
        else
            # SQLite 不需要连接检查
            echo "[entrypoint] Using SQLite database"
            break
        fi
        
        attempt=$((attempt + 1))
        echo "[entrypoint] Waiting for database... (${attempt}/${max_attempts})"
        sleep 2
    done
    
    if [ $attempt -eq $max_attempts ]; then
        echo "[entrypoint] ERROR: Failed to connect to database after ${max_attempts} attempts"
        exit 1
    fi
    
    echo "[entrypoint] Database connectivity verified"
}

# 信号处理函数 - 优雅关闭所有子进程
cleanup() {
    echo "[entrypoint] Received shutdown signal, stopping services..."
    
    if [ -n "$SERVER_PID" ] && kill -0 "$SERVER_PID" 2>/dev/null; then
        echo "[entrypoint] Stopping Web service (PID: $SERVER_PID)..."
        kill -TERM "$SERVER_PID" 2>/dev/null || true
    fi
    
    if [ -n "$MCP_PID" ] && kill -0 "$MCP_PID" 2>/dev/null; then
        echo "[entrypoint] Stopping MCP service (PID: $MCP_PID)..."
        kill -TERM "$MCP_PID" 2>/dev/null || true
    fi
    
    # 等待子进程退出
    wait
    
    echo "[entrypoint] All services stopped."
    exit 0
}

# 注册信号处理
trap cleanup SIGTERM SIGINT SIGQUIT

echo "[entrypoint] Starting SMART_TEST services..."

# 生成证书
generate_certs

# 检查数据库连接
check_database

echo "[entrypoint] Database type: ${DB_TYPE:-sqlite}"
echo "[entrypoint] Note: Table schema will be auto-created by GORM on first startup"

# 安装 Playwright 驱动（如果未安装）
# playwright-go 在首次 Run() 时会自动下载驱动，这里不需要手动安装
# 驱动会被下载到 /root/.cache/ms-playwright-go/
echo "[entrypoint] Playwright driver will be auto-downloaded on first use"

# 启动 Web 服务（后台运行）
echo "[entrypoint] Starting Web service on port 8443..."
./server &
SERVER_PID=$!
echo "[entrypoint] Web service started (PID: $SERVER_PID)"

# 等待 Web 服务启动
sleep 2

# 启动 MCP 服务（后台运行）
echo "[entrypoint] Starting MCP service on port 16410..."
./mcp-server -config ./config/mcp-server.yaml &
MCP_PID=$!
echo "[entrypoint] MCP service started (PID: $MCP_PID)"

echo "[entrypoint] All services started successfully."
echo "[entrypoint] Web service: https://localhost:8443"
echo "[entrypoint] MCP service: http://localhost:16410"

# 等待任一进程退出
wait -n

# 如果任一进程退出，记录并清理
echo "[entrypoint] A service has stopped unexpectedly, shutting down..."
cleanup
