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

# 检查和创建存储目录
check_storage_directories() {
    local storage_base="${STORAGE_BASE_PATH:-/app/storage}"
    
    echo "[entrypoint] Checking storage directories..."
    echo "[entrypoint] Storage base path: ${storage_base}"
    echo "[entrypoint] Running as user: $(whoami), UID: $(id -u), GID: $(id -g)"
    
    # 创建基础目录（如果不存在）
    if [ ! -d "${storage_base}" ]; then
        echo "[entrypoint] Creating storage directory..."
        mkdir -p "${storage_base}" || {
            echo "[entrypoint] ERROR: Cannot create storage directory: ${storage_base}"
            exit 1
        }
    fi
    
    # 确保目录可写
    if [ ! -w "${storage_base}" ]; then
        echo "[entrypoint] Storage directory is not writable, fixing permissions..."
        chmod 755 "${storage_base}" || {
            echo "[entrypoint] WARNING: Failed to fix permissions on ${storage_base}"
        }
    fi
    
    # 创建子目录结构（按功能分组）
    local subdirs=(
        "raw_documents"         # 原始需求文档上传
        "test_files"            # 测试文件
        "screenshots"           # 执行截图
        "defects"               # 缺陷附件上传
        "versions"              # 版本文件
        "versions/auto-cases"   # Web自动化用例版本
        "versions/api-cases"    # API用例版本
        "versions/web-cases"    # Web用例版本打包
        "temp"                  # 临时文件
    )
    
    for subdir in "${subdirs[@]}"; do
        local dir_path="${storage_base}/${subdir}"
        if [ ! -d "${dir_path}" ]; then
            echo "[entrypoint] Creating directory: ${dir_path}"
            mkdir -p "${dir_path}" || {
                echo "[entrypoint] WARNING: Failed to create ${dir_path}"
            }
        fi
    done
    
    # 显示最终权限状态
    echo "[entrypoint] Storage directory permissions:"
    ls -ld "${storage_base}" || true
    echo "[entrypoint] Storage directories verified"
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

# 检查存储目录
check_storage_directories

# 检查数据库连接
check_database

echo "[entrypoint] Database type: ${DB_TYPE:-sqlite}"
echo "[entrypoint] Note: Table schema will be auto-created by GORM on first startup"

# 安装 Playwright 驱动（如果未安装）
# playwright-go 在首次 Run() 时会自动下载驱动，这里不需要手动安装
# 驱动会被下载到 /root/.cache/ms-playwright-go/
echo "[entrypoint] Playwright driver will be auto-downloaded on first use"

# 检查并记录 PROMPTS_DIR 配置
if [ -n "$PROMPTS_DIR" ]; then
    echo "[entrypoint] PROMPTS_DIR is set to: $PROMPTS_DIR"
    if [ -d "$PROMPTS_DIR" ]; then
        echo "[entrypoint] Prompts directory exists"
        PROMPT_COUNT=$(find "$PROMPTS_DIR" -name "*.prompt.md" | wc -l)
        echo "[entrypoint] Found $PROMPT_COUNT prompt files in $PROMPTS_DIR"
        if [ "$PROMPT_COUNT" -gt 0 ]; then
            echo "[entrypoint] Listing prompt files:"
            ls -la "$PROMPTS_DIR"/*.prompt.md 2>/dev/null || echo "[entrypoint] No .prompt.md files found"
        else
            echo "[entrypoint] WARNING: No .prompt.md files found in $PROMPTS_DIR"
        fi
    else
        echo "[entrypoint] ERROR: PROMPTS_DIR directory does not exist: $PROMPTS_DIR"
        echo "[entrypoint] Checking /app/internal/mcp/prompts..."
        if [ -d "/app/internal/mcp/prompts" ]; then
            echo "[entrypoint] Found prompts at /app/internal/mcp/prompts, using this path"
            export PROMPTS_DIR="/app/internal/mcp/prompts"
        fi
    fi
else
    echo "[entrypoint] WARNING: PROMPTS_DIR environment variable is not set"
    echo "[entrypoint] Checking default location /app/internal/mcp/prompts..."
    if [ -d "/app/internal/mcp/prompts" ]; then
        PROMPT_COUNT=$(find "/app/internal/mcp/prompts" -name "*.prompt.md" | wc -l)
        echo "[entrypoint] Found $PROMPT_COUNT prompt files in default location"
        export PROMPTS_DIR="/app/internal/mcp/prompts"
        echo "[entrypoint] Set PROMPTS_DIR to: $PROMPTS_DIR"
    else
        echo "[entrypoint] ERROR: Prompts directory not found in default location"
    fi
fi

# 启动 Web 服务（后台运行）
echo "[entrypoint] Starting Web service on port 8443..."
echo "[entrypoint] Environment: DB_TYPE=${DB_TYPE}, PROMPTS_DIR=$PROMPTS_DIR"
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
