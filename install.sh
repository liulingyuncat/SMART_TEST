#!/bin/bash
# ===========================================
# SMART TEST 平台部署前安装脚本
# 用于生成必要的运行时文件
# ===========================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$SCRIPT_DIR/backend"
CERTS_DIR="$BACKEND_DIR/certs"
STORAGE_DIR="$BACKEND_DIR/storage"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# ===========================================
# 1. 创建必要的目录结构
# ===========================================
create_directories() {
    log_info "创建必要的目录结构..."
    
    # 存储目录
    mkdir -p "$STORAGE_DIR/static"
    mkdir -p "$STORAGE_DIR/attachments"
    mkdir -p "$STORAGE_DIR/versions"
    mkdir -p "$STORAGE_DIR/exports"
    mkdir -p "$STORAGE_DIR/raw_documents"
    mkdir -p "$STORAGE_DIR/defects"
    mkdir -p "$STORAGE_DIR/tmp"
    
    # 证书目录
    mkdir -p "$CERTS_DIR"
    
    # 后端临时目录
    mkdir -p "$BACKEND_DIR/tmp"
    
    log_info "目录结构创建完成"
}

# ===========================================
# 2. 生成自签名 HTTPS 证书
# ===========================================
generate_certificates() {
    log_info "检查并生成 HTTPS 证书..."
    
    CERT_FILE="$CERTS_DIR/server.crt"
    KEY_FILE="$CERTS_DIR/server.key"
    
    # 检查证书是否已存在
    if [ -f "$CERT_FILE" ] && [ -f "$KEY_FILE" ]; then
        log_warn "证书文件已存在，跳过生成"
        log_info "  - 证书: $CERT_FILE"
        log_info "  - 私钥: $KEY_FILE"
        return 0
    fi
    
    # 检查 openssl 是否可用
    if ! command -v openssl &> /dev/null; then
        log_error "openssl 未安装，无法生成证书"
        log_info "请手动安装 openssl 后重新运行此脚本"
        return 1
    fi
    
    log_info "使用 openssl 生成自签名证书..."
    
    # 证书配置
    CERT_DAYS=365
    CERT_CN="localhost"
    CERT_SUBJ="/C=CN/ST=Beijing/L=Beijing/O=SMART_TEST/OU=Development/CN=$CERT_CN"
    
    # 创建证书扩展配置文件
    CERT_EXT_FILE="$CERTS_DIR/cert_ext.cnf"
    cat > "$CERT_EXT_FILE" << EOF
[req]
default_bits = 2048
prompt = no
default_md = sha256
distinguished_name = dn
x509_extensions = v3_req

[dn]
C = CN
ST = Beijing
L = Beijing
O = SMART_TEST
OU = Development
CN = localhost

[v3_req]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = webtest.local
DNS.3 = *.webtest.local
IP.1 = 127.0.0.1
IP.2 = ::1
EOF

    # 生成私钥和证书
    openssl req -x509 -nodes -days $CERT_DAYS \
        -newkey rsa:2048 \
        -keyout "$KEY_FILE" \
        -out "$CERT_FILE" \
        -config "$CERT_EXT_FILE" \
        -extensions v3_req
    
    # 清理临时配置文件
    rm -f "$CERT_EXT_FILE"
    
    # 设置文件权限
    chmod 600 "$KEY_FILE"
    chmod 644 "$CERT_FILE"
    
    log_info "证书生成完成："
    log_info "  - 证书: $CERT_FILE"
    log_info "  - 私钥: $KEY_FILE"
    log_info "  - 有效期: $CERT_DAYS 天"
    
    # 显示证书信息
    log_info "证书信息："
    openssl x509 -in "$CERT_FILE" -noout -subject -dates
}

# ===========================================
# 3. 创建环境配置文件
# ===========================================
create_env_file() {
    log_info "检查环境配置文件..."
    
    ENV_FILE="$SCRIPT_DIR/.env"
    ENV_EXAMPLE="$SCRIPT_DIR/.env.example"
    
    if [ -f "$ENV_FILE" ]; then
        log_warn ".env 文件已存在，跳过创建"
        return 0
    fi
    
    if [ -f "$ENV_EXAMPLE" ]; then
        log_info "从 .env.example 复制环境配置..."
        cp "$ENV_EXAMPLE" "$ENV_FILE"
        
        # 生成随机密钥
        if command -v openssl &> /dev/null; then
            JWT_SECRET=$(openssl rand -base64 32)
            DB_PASSWORD=$(openssl rand -hex 16)
            MCP_TOKEN=$(openssl rand -base64 24)
            
            # 替换默认值（适用于 Linux/Mac）
            if [[ "$OSTYPE" == "darwin"* ]]; then
                # macOS
                sed -i '' "s/default_jwt_secret_please_change_in_production/$JWT_SECRET/" "$ENV_FILE"
                sed -i '' "s/webtest_default_pass_change_me/$DB_PASSWORD/" "$ENV_FILE"
                sed -i '' "s/default_mcp_token_change_me/$MCP_TOKEN/" "$ENV_FILE"
            else
                # Linux
                sed -i "s/default_jwt_secret_please_change_in_production/$JWT_SECRET/" "$ENV_FILE"
                sed -i "s/webtest_default_pass_change_me/$DB_PASSWORD/" "$ENV_FILE"
                sed -i "s/default_mcp_token_change_me/$MCP_TOKEN/" "$ENV_FILE"
            fi
            
            log_info "已生成随机密钥并写入 .env 文件"
        fi
        
        log_info ".env 文件创建完成"
    else
        log_warn ".env.example 不存在，跳过环境配置"
    fi
}

# ===========================================
# 4. 检查 MCP 配置文件
# ===========================================
check_mcp_config() {
    log_info "检查 MCP 配置文件..."
    
    MCP_CONFIG="$BACKEND_DIR/config/mcp-server.yaml"
    
    if [ -f "$MCP_CONFIG" ]; then
        log_info "MCP 配置文件存在: $MCP_CONFIG"
    else
        log_warn "MCP 配置文件不存在，请确认是否需要手动创建"
    fi
}

# ===========================================
# 5. 设置文件权限
# ===========================================
set_permissions() {
    log_info "设置文件权限..."
    
    # 确保脚本可执行
    chmod +x "$SCRIPT_DIR/docker-entrypoint.sh" 2>/dev/null || true
    chmod +x "$SCRIPT_DIR/install.sh" 2>/dev/null || true
    
    # 存储目录权限
    chmod -R 755 "$STORAGE_DIR" 2>/dev/null || true
    
    log_info "文件权限设置完成"
}

# ===========================================
# 主函数
# ===========================================
main() {
    echo "=========================================="
    echo " SMART TEST 平台部署前安装脚本"
    echo "=========================================="
    echo ""
    
    create_directories
    echo ""
    
    generate_certificates
    echo ""
    
    create_env_file
    echo ""
    
    check_mcp_config
    echo ""
    
    set_permissions
    echo ""
    
    echo "=========================================="
    log_info "安装脚本执行完成！"
    echo "=========================================="
    echo ""
    log_info "下一步操作："
    echo "  1. 检查并修改 .env 文件中的配置"
    echo "  2. 运行 docker-compose up -d 启动服务"
    echo "  3. 访问 https://localhost:8443"
    echo ""
}

# 执行主函数
main "$@"
