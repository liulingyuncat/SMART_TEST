#!/bin/bash
# =============================================================================
# SMART TEST 平台 - 生产环境部署准备脚本
#
# 功能: 在执行 docker compose 之前，准备必要的环境配置
#
# 使用方法:
#   1. 下载 docker-compose.yml 和 install.sh
#   2. 运行 ./install.sh 生成配置
#   3. 运行 docker compose pull && docker compose up -d
#
# 本脚本会生成:
#   - .env (包含随机生成的安全密钥)
#   - storage/ 目录 (数据持久化)
# =============================================================================

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }
log_step() { echo -e "${CYAN}[STEP]${NC} $1"; }

echo ""
echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}   SMART TEST 平台 - 生产环境部署准备${NC}"
echo -e "${BLUE}================================================${NC}"
echo ""

# ===========================================
# 1. 检查 Docker 环境
# ===========================================
check_docker() {
    log_info "检查 Docker 环境..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker"
        echo "  安装指南: https://docs.docker.com/get-docker/"
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        log_error "Docker 服务未运行，请启动 Docker"
        exit 1
    fi
    
    # 检查 docker compose 命令
    if docker compose version &> /dev/null; then
        COMPOSE_CMD="docker compose"
    elif command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose"
    else
        log_error "Docker Compose 未安装"
        exit 1
    fi
    
    log_info "Docker 环境检查通过 ✓"
    log_info "Compose 命令: $COMPOSE_CMD"
}

# ===========================================
# 2. 检查 docker-compose.yml
# ===========================================
check_compose_file() {
    log_info "检查 docker-compose.yml..."
    
    if [ ! -f "docker-compose.yml" ]; then
        log_error "docker-compose.yml 文件不存在！"
        echo ""
        echo "请先下载 docker-compose.yml 文件:"
        echo -e "  ${BLUE}curl -O https://raw.githubusercontent.com/liulingyuncat/SMART_TEST/main/docker-compose.yml${NC}"
        echo ""
        exit 1
    fi
    
    log_info "docker-compose.yml 检查通过 ✓"
}

# ===========================================
# 3. 生成随机密钥
# ===========================================
generate_secret() {
    local length=${1:-32}
    if command -v openssl &> /dev/null; then
        openssl rand -base64 $length | tr -d '=/+' | head -c $length
    else
        head -c $length /dev/urandom | base64 | tr -d '=/+' | head -c $length
    fi
}

# ===========================================
# 4. 创建 .env 文件
# ===========================================
create_env_file() {
    log_info "生成 .env 配置文件..."
    
    if [ -f ".env" ]; then
        log_warn ".env 文件已存在，跳过生成"
        log_warn "如需重新生成密钥，请先删除现有 .env 文件"
        return 0
    fi
    
    local db_password=$(generate_secret 24)
    local jwt_secret=$(generate_secret 48)
    local mcp_token=$(generate_secret 32)
    
    cat > .env << ENV_EOF
# =============================================================================
# SMART TEST - 环境变量配置
# 生成时间: $(date '+%Y-%m-%d %H:%M:%S')
# 警告: 此文件包含敏感信息，请勿提交到版本控制！
# =============================================================================

# 数据库密码 (PostgreSQL)
DB_PASSWORD=${db_password}

# JWT 签名密钥 (用于用户认证，至少32字符)
JWT_SECRET=${jwt_secret}

# MCP 服务认证令牌 (可选，用于 AI 工具集成)
MCP_AUTH_TOKEN=${mcp_token}
ENV_EOF
    
    chmod 600 .env
    log_info ".env 文件创建完成 ✓"
    log_info "文件权限已设为 600 (仅所有者可读写)"
}

# ===========================================
# 5. 创建数据目录
# ===========================================
create_directories() {
    log_info "创建数据持久化目录..."
    
    # 应用数据目录 (附件、导出文件等)
    mkdir -p storage
    log_info "storage/ 目录创建完成 ✓"
    
    # PostgreSQL 数据目录
    mkdir -p data/postgres
    log_info "data/postgres/ 目录创建完成 ✓"
}

# ===========================================
# 6. 显示后续步骤
# ===========================================
show_next_steps() {
    echo ""
    echo -e "${GREEN}================================================${NC}"
    echo -e "${GREEN}   环境准备完成！${NC}"
    echo -e "${GREEN}================================================${NC}"
    echo ""
    echo -e "${CYAN}已生成的文件和目录:${NC}"
    echo "  ├── .env            # 环境变量配置 (含随机生成的安全密钥)"
    echo "  ├── storage/        # 应用数据目录 (附件、导出文件等)"
    echo "  └── data/postgres/  # PostgreSQL 数据库目录"
    echo ""
    echo -e "${CYAN}========================================${NC}"
    echo -e "${CYAN}   后续步骤说明${NC}"
    echo -e "${CYAN}========================================${NC}"
    echo ""
    
    # Step 1: 拉取镜像
    log_step "步骤 1: 拉取 Docker 镜像"
    echo ""
    echo "  执行命令:"
    echo -e "    ${BLUE}$COMPOSE_CMD pull${NC}"
    echo ""
    echo "  说明: 从 GitHub Container Registry (ghcr.io) 拉取最新的应用镜像。"
    echo "        首次拉取可能需要几分钟，取决于网络速度。"
    echo "        镜像包含: 前端 React 应用 + 后端 Go 服务 + MCP 服务器"
    echo ""
    
    # Step 2: 启动服务
    log_step "步骤 2: 启动服务"
    echo ""
    echo "  执行命令:"
    echo -e "    ${BLUE}$COMPOSE_CMD up -d${NC}"
    echo ""
    echo "  说明: 以后台模式启动所有服务容器。"
    echo "        -d 参数表示 detached mode (后台运行)"
    echo "        将启动以下服务:"
    echo "          • webtest      - 主应用 (Web + API + MCP)"
    echo "          • postgres     - PostgreSQL 数据库"
    echo ""
    
    # Step 3: 查看状态
    log_step "步骤 3: 查看服务状态"
    echo ""
    echo "  执行命令:"
    echo -e "    ${BLUE}$COMPOSE_CMD ps${NC}"
    echo ""
    echo "  说明: 查看所有容器的运行状态。"
    echo "        正常情况下，所有服务状态应为 'Up' 或 'running'。"
    echo ""
    
    # Step 4: 查看日志
    log_step "步骤 4: 查看服务日志 (可选)"
    echo ""
    echo "  执行命令:"
    echo -e "    ${BLUE}$COMPOSE_CMD logs -f webtest${NC}    # 实时查看应用日志"
    echo -e "    ${BLUE}$COMPOSE_CMD logs -f postgres${NC}   # 实时查看数据库日志"
    echo -e "    ${BLUE}$COMPOSE_CMD logs --tail=100${NC}    # 查看最近100行日志"
    echo ""
    echo "  说明: -f 参数表示 follow mode，实时输出新日志。"
    echo "        按 Ctrl+C 退出日志查看。"
    echo ""
    
    # 快捷命令
    echo -e "${CYAN}========================================${NC}"
    echo -e "${CYAN}   快捷命令 (一键启动)${NC}"
    echo -e "${CYAN}========================================${NC}"
    echo ""
    echo "  如果想一次性完成拉取和启动，执行:"
    echo -e "    ${BLUE}$COMPOSE_CMD pull && $COMPOSE_CMD up -d${NC}"
    echo ""
    
    # 访问信息
    echo -e "${CYAN}========================================${NC}"
    echo -e "${CYAN}   访问信息${NC}"
    echo -e "${CYAN}========================================${NC}"
    echo ""
    echo "  Web 管理界面:"
    echo -e "    ${BLUE}https://localhost:8443${NC}"
    echo ""
    echo "  MCP 服务地址 (用于 AI 工具集成):"
    echo -e "    ${BLUE}http://localhost:16410${NC}"
    echo ""
    echo "  默认管理员账号:"
    echo -e "    用户名: ${BLUE}admin${NC}"
    echo -e "    密码:   ${BLUE}admin123${NC}"
    echo ""
    echo -e "  ${YELLOW}⚠️  首次登录后请立即修改默认密码！${NC}"
    echo ""
    
    # 常用运维命令
    echo -e "${CYAN}========================================${NC}"
    echo -e "${CYAN}   常用运维命令${NC}"
    echo -e "${CYAN}========================================${NC}"
    echo ""
    echo "  停止服务:"
    echo -e "    ${BLUE}$COMPOSE_CMD down${NC}"
    echo ""
    echo "  重启服务:"
    echo -e "    ${BLUE}$COMPOSE_CMD restart${NC}"
    echo ""
    echo "  更新到最新版本:"
    echo -e "    ${BLUE}$COMPOSE_CMD pull && $COMPOSE_CMD up -d${NC}"
    echo ""
    echo "  完全清理 (包括数据卷):"
    echo -e "    ${BLUE}$COMPOSE_CMD down -v${NC}"
    echo -e "    ${YELLOW}⚠️  警告: 此命令会删除所有数据！${NC}"
    echo ""
}

# ===========================================
# 主函数
# ===========================================
main() {
    check_docker
    echo ""
    
    check_compose_file
    echo ""
    
    create_env_file
    create_directories
    
    show_next_steps
}

# 执行
main "$@"
