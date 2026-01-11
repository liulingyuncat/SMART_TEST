#!/bin/bash
# =============================================================================
# SMART TEST 平台 - 生产环境一键部署脚本
#
# 用法: curl -sSL https://raw.githubusercontent.com/liulingyuncat/SMART_TEST/main/install.sh | bash
# 或者: ./install.sh
#
# 本脚本会生成:
#   - docker-compose.yml
#   - .env (包含随机生成的安全密钥)
#   - storage/ 目录 (数据持久化)
# =============================================================================

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   SMART TEST 平台 - 生产环境部署${NC}"
echo -e "${BLUE}========================================${NC}"
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
    
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        log_error "Docker Compose 未安装"
        exit 1
    fi
    
    log_info "Docker 环境检查通过 ✓"
}

# ===========================================
# 2. 生成随机密钥
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
# 3. 创建 docker-compose.yml
# ===========================================
create_docker_compose() {
    log_info "生成 docker-compose.yml..."
    
    if [ -f "docker-compose.yml" ]; then
        log_warn "docker-compose.yml 已存在，跳过生成"
        return 0
    fi
    
    cat > docker-compose.yml << 'COMPOSE_EOF'
# =============================================================================
# SMART TEST - Production Docker Compose
# 生成时间: $(date '+%Y-%m-%d %H:%M:%S')
# =============================================================================

version: '3.8'

services:
  webtest:
    image: ghcr.io/liulingyuncat/smart_test:latest
    container_name: webtest
    restart: unless-stopped
    ports:
      - "8443:8443"   # Web HTTPS
      - "16410:16410" # MCP HTTP
    environment:
      # 数据库配置
      - DB_TYPE=postgres
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=webtest
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=webtest
      # 应用配置
      - JWT_SECRET=${JWT_SECRET}
      - TZ=Asia/Shanghai
      # MCP 配置 (可选)
      - MCP_AUTH_TOKEN=${MCP_AUTH_TOKEN}
      - MCP_BACKEND_URL=https://localhost:8443
    volumes:
      - ./storage:/app/storage
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - webtest-network
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8443/api/v1/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 15s

  postgres:
    image: postgres:16-alpine
    container_name: webtest-postgres
    restart: unless-stopped
    environment:
      - POSTGRES_DB=webtest
      - POSTGRES_USER=webtest
      - POSTGRES_PASSWORD=${DB_PASSWORD}
      - TZ=Asia/Shanghai
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - webtest-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U webtest -d webtest"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s

networks:
  webtest-network:
    driver: bridge

volumes:
  postgres_data:
COMPOSE_EOF
    
    log_info "docker-compose.yml 创建完成 ✓"
}

# ===========================================
# 4. 创建 .env 文件
# ===========================================
create_env_file() {
    log_info "生成 .env 配置文件..."
    
    if [ -f ".env" ]; then
        log_warn ".env 文件已存在，跳过生成"
        log_warn "如需重新生成，请先删除现有 .env 文件"
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

# 数据库密码
DB_PASSWORD=${db_password}

# JWT 签名密钥 (至少32字符)
JWT_SECRET=${jwt_secret}

# MCP 服务认证令牌 (可选，用于 AI 工具集成)
MCP_AUTH_TOKEN=${mcp_token}
ENV_EOF
    
    chmod 600 .env
    log_info ".env 文件创建完成 ✓"
}

# ===========================================
# 5. 创建数据目录
# ===========================================
create_directories() {
    log_info "创建数据目录..."
    mkdir -p storage
    log_info "storage/ 目录创建完成 ✓"
}

# ===========================================
# 6. 拉取镜像
# ===========================================
pull_image() {
    log_info "拉取最新 Docker 镜像..."
    docker pull ghcr.io/liulingyuncat/smart_test:latest
    log_info "镜像拉取完成 ✓"
}

# ===========================================
# 主函数
# ===========================================
main() {
    check_docker
    echo ""
    
    create_docker_compose
    create_env_file
    create_directories
    echo ""
    
    # 询问是否拉取镜像
    read -p "是否立即拉取 Docker 镜像? [Y/n] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]] || [[ -z $REPLY ]]; then
        pull_image
    fi
    
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}   部署准备完成！${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo "生成的文件:"
    echo "  - docker-compose.yml  (容器编排)"
    echo "  - .env                (环境配置)"
    echo "  - storage/            (数据目录)"
    echo ""
    echo "启动服务:"
    echo -e "  ${BLUE}docker compose up -d${NC}"
    echo ""
    echo "访问地址:"
    echo -e "  Web 界面: ${BLUE}https://localhost:8443${NC}"
    echo -e "  MCP 服务: ${BLUE}http://localhost:16410${NC}"
    echo ""
    echo "默认账号:"
    echo -e "  用户名: ${BLUE}admin${NC}"
    echo -e "  密码:   ${BLUE}admin123${NC}"
    echo ""
    echo -e "${YELLOW}⚠️  首次登录后请立即修改默认密码！${NC}"
    echo ""
}

# 执行
main "$@"
