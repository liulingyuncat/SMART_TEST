#!/bin/bash
# =============================================================================
# SMART_TEST - AutoMigrate 模型完整性检查脚本
#
# 功能：
#   1. 扫描 backend/internal/models/ 目录下所有定义的 GORM 模型
#   2. 检查 backend/cmd/server/main.go 中的 AutoMigrate 调用
#   3. 验证所有模型都已正确注册到 AutoMigrate
#   4. 报告缺失的模型
#
# 用法：
#   ./check_automigrate.sh           # 在项目根目录运行
#   cd backend && ../check_automigrate.sh  # 或在 backend 目录运行
#
# 退出码：
#   0 - 所有模型都已注册
#   1 - 发现缺失的模型或检查失败
# =============================================================================

set -e

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[✓]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[!]${NC} $1"; }
log_error() { echo -e "${RED}[✗]${NC} $1"; }
log_step() { echo -e "${BLUE}[→]${NC} $1"; }

# 确定项目根目录
if [ -f "go.mod" ]; then
    # 在 backend 目录
    PROJECT_ROOT=".."
    BACKEND_DIR="."
elif [ -f "backend/go.mod" ]; then
    # 在项目根目录
    PROJECT_ROOT="."
    BACKEND_DIR="backend"
else
    log_error "无法找到项目根目录，请在项目根目录或 backend 目录下运行此脚本"
    exit 1
fi

MODELS_DIR="${BACKEND_DIR}/internal/models"
MAIN_FILE="${BACKEND_DIR}/cmd/server/main.go"

echo ""
echo "========================================"
echo "  SMART_TEST AutoMigrate 完整性检查"
echo "========================================"
echo ""

# 检查必需文件是否存在
log_step "检查必需文件..."
if [ ! -d "$MODELS_DIR" ]; then
    log_error "模型目录不存在: $MODELS_DIR"
    exit 1
fi

if [ ! -f "$MAIN_FILE" ]; then
    log_error "主文件不存在: $MAIN_FILE"
    exit 1
fi

log_info "找到模型目录: $MODELS_DIR"
log_info "找到主文件: $MAIN_FILE"
echo ""

# 第一步：提取所有 GORM 模型（包含 gorm.Model 或 GORM 标签的结构体）
log_step "扫描所有 GORM 模型定义..."

# 提取真正的数据库模型（包含 gorm.Model 或 gorm 标签的结构体）
DEFINED_MODELS=()
for model_file in "$MODELS_DIR"/*.go; do
    if [ -f "$model_file" ]; then
        # 读取所有结构体定义
        while IFS= read -r line; do
            model_name=$(echo "$line" | awk '{print $2}')
            if [ -n "$model_name" ]; then
                # 提取该结构体的完整定义（从 type 到 }）
                struct_content=$(awk "/^type $model_name struct/,/^}/" "$model_file")
                
                # 检查是否包含 gorm.Model 或 gorm: 标签（真正的数据库模型）
                if echo "$struct_content" | grep -qE "(gorm\.Model|gorm:)"; then
                    DEFINED_MODELS+=("$model_name")
                fi
            fi
        done < <(grep "^type .* struct" "$model_file" 2>/dev/null || true)
    fi
done

# 去重并排序
DEFINED_MODELS=($(printf '%s\n' "${DEFINED_MODELS[@]}" | sort -u))

log_info "发现 ${#DEFINED_MODELS[@]} 个 GORM 模型定义"
for model in "${DEFINED_MODELS[@]}"; do
    echo "   - $model"
done
echo ""

# 第二步：提取 AutoMigrate 中注册的模型
log_step "检查 AutoMigrate 注册列表..."

# 提取 AutoMigrate 调用中的所有模型
# 匹配模式：&models.ModelName{}
REGISTERED_MODELS=()
if grep -q "db.AutoMigrate" "$MAIN_FILE"; then
    # 提取 AutoMigrate 区块内的所有 &models.XXX{}
    while IFS= read -r line; do
        model_name=$(echo "$line" | sed -n 's/.*&models\.\([A-Za-z0-9_]*\){}.*/\1/p')
        if [ -n "$model_name" ]; then
            REGISTERED_MODELS+=("$model_name")
        fi
    done < <(sed -n '/db.AutoMigrate(/,/); err != nil/p' "$MAIN_FILE")
else
    log_error "在 $MAIN_FILE 中未找到 AutoMigrate 调用"
    exit 1
fi

# 去重并排序
REGISTERED_MODELS=($(printf '%s\n' "${REGISTERED_MODELS[@]}" | sort -u))

log_info "发现 ${#REGISTERED_MODELS[@]} 个已注册模型"
for model in "${REGISTERED_MODELS[@]}"; do
    echo "   - $model"
done
echo ""

# 第三步：对比差异
log_step "对比差异..."

MISSING_MODELS=()
for model in "${DEFINED_MODELS[@]}"; do
    found=false
    for registered in "${REGISTERED_MODELS[@]}"; do
        if [ "$model" = "$registered" ]; then
            found=true
            break
        fi
    done
    
    if [ "$found" = false ]; then
        MISSING_MODELS+=("$model")
    fi
done

# 第四步：输出结果
echo ""
echo "========================================"
echo "  检查结果"
echo "========================================"
echo ""

if [ ${#MISSING_MODELS[@]} -eq 0 ]; then
    log_info "✅ 所有模型都已正确注册到 AutoMigrate"
    echo ""
    exit 0
else
    log_error "❌ 发现 ${#MISSING_MODELS[@]} 个模型未注册到 AutoMigrate："
    echo ""
    for model in "${MISSING_MODELS[@]}"; do
        echo -e "${RED}   ✗ $model${NC}"
    done
    echo ""
    log_warn "请在 $MAIN_FILE 的 AutoMigrate 调用中添加以下行："
    echo ""
    for model in "${MISSING_MODELS[@]}"; do
        echo -e "${YELLOW}   &models.${model}{},"
    done
    echo -e "${NC}"
    exit 1
fi
