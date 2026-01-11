#!/bin/bash
# =============================================================================
# SMART TEST 平台 - 数据备份与迁移脚本
#
# 功能:
#   backup  - 备份数据库和应用数据到压缩包
#   restore - 从备份压缩包恢复数据
#
# 用法:
#   ./migrate.sh backup              # 备份到默认目录
#   ./migrate.sh backup /path/to    # 备份到指定目录
#   ./migrate.sh restore backup.tar.gz  # 从备份恢复
#
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

# 配置
BACKUP_DIR="./backups"
TIMESTAMP=$(date '+%Y%m%d_%H%M%S')
BACKUP_NAME="smarttest_backup_${TIMESTAMP}"

# ===========================================
# 显示帮助信息
# ===========================================
show_help() {
    echo ""
    echo -e "${BLUE}SMART TEST 平台 - 数据备份与迁移工具${NC}"
    echo ""
    echo "用法:"
    echo "  $0 backup [目标目录]      备份数据"
    echo "  $0 restore <备份文件>     恢复数据"
    echo "  $0 help                   显示帮助"
    echo ""
    echo "示例:"
    echo "  $0 backup                 # 备份到 ./backups/ 目录"
    echo "  $0 backup /mnt/backup     # 备份到指定目录"
    echo "  $0 restore smarttest_backup_20260111_120000.tar.gz"
    echo ""
    echo "备份内容:"
    echo "  - data/postgres/          PostgreSQL 数据库文件"
    echo "  - storage/                应用数据 (附件、导出文件等)"
    echo "  - .env                    环境配置文件"
    echo ""
}

# ===========================================
# 检查 Docker 环境
# ===========================================
check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装"
        exit 1
    fi
    
    if docker compose version &> /dev/null; then
        COMPOSE_CMD="docker compose"
    elif command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose"
    else
        log_error "Docker Compose 未安装"
        exit 1
    fi
}

# ===========================================
# 检查服务状态
# ===========================================
check_services_running() {
    if docker ps --format '{{.Names}}' | grep -q "webtest-postgres"; then
        return 0
    else
        return 1
    fi
}

# ===========================================
# 备份数据
# ===========================================
do_backup() {
    local target_dir="${1:-$BACKUP_DIR}"
    
    echo ""
    echo -e "${BLUE}================================================${NC}"
    echo -e "${BLUE}   SMART TEST - 数据备份${NC}"
    echo -e "${BLUE}================================================${NC}"
    echo ""
    
    # 检查数据目录是否存在
    if [ ! -d "data/postgres" ] && [ ! -d "storage" ]; then
        log_error "未找到数据目录 (data/postgres 或 storage)"
        log_error "请确认在正确的部署目录中执行此脚本"
        exit 1
    fi
    
    # 创建备份目录
    mkdir -p "$target_dir"
    local backup_file="${target_dir}/${BACKUP_NAME}.tar.gz"
    
    log_info "开始备份..."
    log_info "备份时间: $(date '+%Y-%m-%d %H:%M:%S')"
    echo ""
    
    # 检查服务状态并建议停止
    if check_services_running; then
        log_warn "检测到服务正在运行"
        log_warn "建议在备份前停止服务以确保数据一致性"
        echo ""
        read -p "是否停止服务后继续备份? [Y/n] " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]] || [[ -z $REPLY ]]; then
            log_info "停止服务..."
            $COMPOSE_CMD down
            SERVICES_STOPPED=true
        else
            log_warn "在服务运行时备份，数据可能不一致"
            SERVICES_STOPPED=false
        fi
    else
        SERVICES_STOPPED=false
    fi
    
    # 收集要备份的目录
    local backup_items=""
    
    if [ -d "data/postgres" ]; then
        backup_items="$backup_items data/postgres"
        log_step "包含: data/postgres/ (PostgreSQL 数据库)"
    fi
    
    if [ -d "storage" ]; then
        backup_items="$backup_items storage"
        log_step "包含: storage/ (应用数据)"
    fi
    
    if [ -f ".env" ]; then
        backup_items="$backup_items .env"
        log_step "包含: .env (环境配置)"
    fi
    
    echo ""
    log_info "正在创建压缩包..."
    
    # 创建备份
    tar -czvf "$backup_file" $backup_items 2>/dev/null
    
    # 计算文件大小
    local file_size=$(du -h "$backup_file" | cut -f1)
    
    # 重启服务
    if [ "$SERVICES_STOPPED" = true ]; then
        echo ""
        log_info "重新启动服务..."
        $COMPOSE_CMD up -d
    fi
    
    echo ""
    echo -e "${GREEN}================================================${NC}"
    echo -e "${GREEN}   备份完成！${NC}"
    echo -e "${GREEN}================================================${NC}"
    echo ""
    echo -e "备份文件: ${BLUE}${backup_file}${NC}"
    echo -e "文件大小: ${BLUE}${file_size}${NC}"
    echo ""
    
    # 显示迁移步骤
    show_migration_steps "$backup_file"
}

# ===========================================
# 显示迁移步骤说明
# ===========================================
show_migration_steps() {
    local backup_file="$1"
    local backup_filename=$(basename "$backup_file")
    
    echo -e "${CYAN}================================================${NC}"
    echo -e "${CYAN}   如何将数据迁移到新服务器${NC}"
    echo -e "${CYAN}================================================${NC}"
    echo ""
    
    log_step "步骤 1: 将备份文件传输到新服务器"
    echo ""
    echo "  使用 scp 复制文件:"
    echo -e "    ${BLUE}scp ${backup_file} user@new-server:/path/to/smarttest/${NC}"
    echo ""
    echo "  或使用 rsync:"
    echo -e "    ${BLUE}rsync -avz ${backup_file} user@new-server:/path/to/smarttest/${NC}"
    echo ""
    
    log_step "步骤 2: 在新服务器上准备环境"
    echo ""
    echo "  创建部署目录并进入:"
    echo -e "    ${BLUE}mkdir -p /path/to/smarttest && cd /path/to/smarttest${NC}"
    echo ""
    echo "  下载必要文件:"
    echo -e "    ${BLUE}curl -O https://raw.githubusercontent.com/liulingyuncat/SMART_TEST/main/docker-compose.yml${NC}"
    echo -e "    ${BLUE}curl -O https://raw.githubusercontent.com/liulingyuncat/SMART_TEST/main/migrate.sh${NC}"
    echo -e "    ${BLUE}chmod +x migrate.sh${NC}"
    echo ""
    
    log_step "步骤 3: 恢复数据"
    echo ""
    echo "  执行恢复命令:"
    echo -e "    ${BLUE}./migrate.sh restore ${backup_filename}${NC}"
    echo ""
    
    log_step "步骤 4: 启动服务"
    echo ""
    echo "  拉取镜像并启动:"
    echo -e "    ${BLUE}docker compose pull${NC}"
    echo -e "    ${BLUE}docker compose up -d${NC}"
    echo ""
    
    log_step "步骤 5: 验证迁移"
    echo ""
    echo "  检查服务状态:"
    echo -e "    ${BLUE}docker compose ps${NC}"
    echo ""
    echo "  查看日志确认正常启动:"
    echo -e "    ${BLUE}docker compose logs -f${NC}"
    echo ""
    echo "  访问 Web 界面验证数据:"
    echo -e "    ${BLUE}https://new-server:8443${NC}"
    echo ""
    
    echo -e "${CYAN}================================================${NC}"
    echo -e "${CYAN}   重要提示${NC}"
    echo -e "${CYAN}================================================${NC}"
    echo ""
    echo "  • 确保新服务器已安装 Docker 和 Docker Compose"
    echo "  • 迁移后请验证所有数据是否完整"
    echo "  • 建议保留备份文件以便需要时回滚"
    echo "  • 如果修改了服务器地址，可能需要更新 .env 中的配置"
    echo ""
}

# ===========================================
# 恢复数据
# ===========================================
do_restore() {
    local backup_file="$1"
    
    echo ""
    echo -e "${BLUE}================================================${NC}"
    echo -e "${BLUE}   SMART TEST - 数据恢复${NC}"
    echo -e "${BLUE}================================================${NC}"
    echo ""
    
    # 检查备份文件
    if [ -z "$backup_file" ]; then
        log_error "请指定备份文件"
        echo "用法: $0 restore <备份文件>"
        exit 1
    fi
    
    if [ ! -f "$backup_file" ]; then
        log_error "备份文件不存在: $backup_file"
        exit 1
    fi
    
    log_info "备份文件: $backup_file"
    log_info "恢复时间: $(date '+%Y-%m-%d %H:%M:%S')"
    echo ""
    
    # 检查服务状态
    if check_services_running; then
        log_warn "检测到服务正在运行，恢复前需要停止服务"
        read -p "是否停止服务并继续恢复? [y/N] " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "已取消恢复操作"
            exit 0
        fi
        log_info "停止服务..."
        $COMPOSE_CMD down
    fi
    
    # 检查现有数据
    if [ -d "data/postgres" ] && [ "$(ls -A data/postgres 2>/dev/null)" ]; then
        echo ""
        log_warn "检测到现有数据目录 data/postgres/"
        log_warn "恢复操作将覆盖现有数据！"
        echo ""
        read -p "是否继续? 输入 'YES' 确认: " confirm
        if [ "$confirm" != "YES" ]; then
            log_info "已取消恢复操作"
            exit 0
        fi
        
        # 备份现有数据
        local old_backup="data_before_restore_${TIMESTAMP}"
        log_info "备份现有数据到 ${old_backup}/"
        mkdir -p "$old_backup"
        [ -d "data" ] && mv data "$old_backup/"
        [ -d "storage" ] && mv storage "$old_backup/"
        [ -f ".env" ] && cp .env "$old_backup/"
    fi
    
    # 解压恢复
    log_info "正在解压备份文件..."
    tar -xzvf "$backup_file"
    
    echo ""
    echo -e "${GREEN}================================================${NC}"
    echo -e "${GREEN}   数据恢复完成！${NC}"
    echo -e "${GREEN}================================================${NC}"
    echo ""
    
    # 显示恢复的内容
    echo "已恢复的内容:"
    [ -d "data/postgres" ] && echo "  ✓ data/postgres/ (PostgreSQL 数据库)"
    [ -d "storage" ] && echo "  ✓ storage/ (应用数据)"
    [ -f ".env" ] && echo "  ✓ .env (环境配置)"
    echo ""
    
    echo -e "${CYAN}后续步骤:${NC}"
    echo ""
    echo "  1. 检查 .env 配置文件是否需要修改"
    echo -e "     ${BLUE}cat .env${NC}"
    echo ""
    echo "  2. 拉取最新镜像"
    echo -e "     ${BLUE}docker compose pull${NC}"
    echo ""
    echo "  3. 启动服务"
    echo -e "     ${BLUE}docker compose up -d${NC}"
    echo ""
    echo "  4. 检查服务状态"
    echo -e "     ${BLUE}docker compose ps${NC}"
    echo ""
}

# ===========================================
# 主函数
# ===========================================
main() {
    check_docker
    
    case "${1:-}" in
        backup)
            do_backup "$2"
            ;;
        restore)
            do_restore "$2"
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            show_help
            exit 1
            ;;
    esac
}

# 执行
main "$@"
