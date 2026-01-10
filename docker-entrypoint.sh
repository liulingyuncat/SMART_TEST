#!/bin/sh
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
echo "[entrypoint] Database type: ${DB_TYPE:-sqlite}"

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
