.PHONY: dev server mcp executor clean build test help install-executor

# 默认目标
.DEFAULT_GOAL := help

# 变量定义
BACKEND_DIR := backend
FRONTEND_DIR := frontend
SERVER_CMD := $(BACKEND_DIR)/cmd/server
MCP_CMD := $(BACKEND_DIR)/cmd/mcp
EXECUTOR_DIR := playwright-executor
BUILD_DIR := bin
SERVER_BIN := $(BUILD_DIR)/server
MCP_BIN := $(BUILD_DIR)/mcp

# 颜色输出
COLOR_RESET := \033[0m
COLOR_BLUE := \033[34m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m

##@ 开发命令

dev: dev-certs ## 启动开发环境（仅后端，前端需单独启动）
	@echo "$(COLOR_BLUE)启动开发环境...$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)注意: 本命令仅启动后端服务，Docker 和前端需单独启动$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)  - Docker: make dev-deps$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)  - 前端: make frontend-dev$(COLOR_RESET)"
	@echo ""
	@if [ -f .env ]; then \
		echo "$(COLOR_GREEN)✓ 从 .env 加载环境变量$(COLOR_RESET)"; \
		export $$(grep -v '^#' .env | xargs); \
	else \
		echo "$(COLOR_YELLOW)⚠ .env 文件不存在，使用默认配置$(COLOR_RESET)"; \
	fi; \
	trap 'echo "$(COLOR_YELLOW)正在关闭服务...$(COLOR_RESET)"; pkill -P $$$$; exit' INT; \
	echo "$(COLOR_GREEN)启动 Server...$(COLOR_RESET)"; \
	cd $(BACKEND_DIR) && go run ./cmd/server/main.go & \
	SERVER_PID=$$!; \
	sleep 2; \
	echo "$(COLOR_GREEN)启动 MCP Server...$(COLOR_RESET)"; \
	cd $(BACKEND_DIR) && go run ./cmd/mcp/main.go & \
	MCP_PID=$$!; \
	echo "$(COLOR_GREEN)✓ 后端服务已启动$(COLOR_RESET)"; \
	echo "  - Server PID: $$SERVER_PID"; \
	echo "  - MCP PID: $$MCP_PID"; \
	echo "  - Server URL: https://localhost:8443"; \
	echo "  - MCP URL: http://localhost:16410"; \
	echo ""; \
	echo "$(COLOR_YELLOW)前端开发服务器：在另一个终端运行 'make frontend-dev'$(COLOR_RESET)"; \
	echo "$(COLOR_YELLOW)或使用 'make dev-all' 一键启动全部服务$(COLOR_RESET)"; \
	echo ""; \
	echo "$(COLOR_YELLOW)按 Ctrl+C 退出...$(COLOR_RESET)"; \
	wait

dev-all: dev-deps dev-certs frontend-build ## 一键启动全部服务（Docker + 构建前端 + 启动后端）
	@echo "$(COLOR_BLUE)启动完整开发环境...$(COLOR_RESET)"
	@if [ -f .env ]; then \
		echo "$(COLOR_GREEN)✓ 从 .env 加载环境变量$(COLOR_RESET)"; \
		export $$(grep -v '^#' .env | xargs); \
	else \
		echo "$(COLOR_YELLOW)⚠ .env 文件不存在，使用默认配置$(COLOR_RESET)"; \
	fi; \
	trap 'echo "$(COLOR_YELLOW)正在关闭服务...$(COLOR_RESET)"; pkill -P $$$$; exit' INT; \
	echo "$(COLOR_GREEN)启动 Server...$(COLOR_RESET)"; \
	cd $(BACKEND_DIR) && go run ./cmd/server/main.go & \
	SERVER_PID=$$!; \
	sleep 2; \
	echo "$(COLOR_GREEN)启动 MCP Server...$(COLOR_RESET)"; \
	cd $(BACKEND_DIR) && go run ./cmd/mcp/main.go & \
	MCP_PID=$$!; \
	echo "$(COLOR_GREEN)✓ 全部服务已启动$(COLOR_RESET)"; \
	echo "  - Server PID: $$SERVER_PID"; \
	echo "  - MCP PID: $$MCP_PID"; \
	echo "  - Frontend: https://localhost:8443"; \
	echo "  - MCP URL: http://localhost:16410"; \
	echo "$(COLOR_YELLOW)按 Ctrl+C 退出...$(COLOR_RESET)"; \
	wait

dev-restart: dev-stop dev ## 重启开发环境

dev-deps: ## 启动开发依赖（PostgreSQL + Playwright）
	@echo "$(COLOR_BLUE)启动开发依赖服务...$(COLOR_RESET)"
	docker compose -f docker-compose.dev.yml up -d
	@echo "$(COLOR_GREEN)✓ 开发依赖已启动$(COLOR_RESET)"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Playwright Runner: ws://localhost:53729/"
	@echo "  - Playwright Executor: http://localhost:53730/"

dev-stop: ## 停止开发环境（关闭进程和容器）
	@echo "$(COLOR_BLUE)停止开发环境...$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)关闭 Go 进程...$(COLOR_RESET)"
	@pkill -9 -f 'go run.*cmd/server' 2>/dev/null || true
	@pkill -9 -f 'go run.*cmd/mcp' 2>/dev/null || true
	@pkill -9 -f 'go-build.*server.*main' 2>/dev/null || true
	@pkill -9 -f 'go-build.*mcp.*main' 2>/dev/null || true
	@sleep 1
	@echo "$(COLOR_YELLOW)关闭 Docker 容器...$(COLOR_RESET)"
	@docker compose -f docker-compose.dev.yml down 2>/dev/null || true
	@echo "$(COLOR_GREEN)✓ 开发环境已停止$(COLOR_RESET)"

dev-clean: dev-stop ## 停止并清理开发环境（包括数据）
	@echo "$(COLOR_BLUE)清理开发环境...$(COLOR_RESET)"
	@docker compose -f docker-compose.dev.yml down -v
	@rm -f backend/webtest.db backend/webtest.db-shm backend/webtest.db-wal
	@echo "$(COLOR_GREEN)✓ 开发环境已清理$(COLOR_RESET)"

dev-deps-stop: ## 停止开发依赖
	@echo "$(COLOR_BLUE)停止开发依赖服务...$(COLOR_RESET)"
	docker compose -f docker-compose.dev.yml down
	@echo "$(COLOR_GREEN)✓ 开发依赖已停止$(COLOR_RESET)"

dev-deps-logs: ## 查看开发依赖日志
	docker compose -f docker-compose.dev.yml logs -f

dev-certs: ## 生成开发用 HTTPS 证书
	@echo "$(COLOR_BLUE)检查证书...$(COLOR_RESET)"
	@if [ ! -f certs/server.crt ]; then \
		echo "$(COLOR_YELLOW)证书不存在，正在生成...$(COLOR_RESET)"; \
		mkdir -p certs; \
		cd $(BACKEND_DIR)/cmd/gencert && go build -o ../../../$(BUILD_DIR)/gencert main.go; \
		cd ../../.. && ./$(BUILD_DIR)/gencert; \
		echo "$(COLOR_GREEN)✓ 证书生成完成$(COLOR_RESET)"; \
	else \
		echo "$(COLOR_GREEN)✓ 证书已存在$(COLOR_RESET)"; \
	fi

server: ## 仅启动 Server
	@echo "$(COLOR_BLUE)启动 Server...$(COLOR_RESET)"
	@if [ -f .env ]; then \
		export $$(grep -v '^#' .env | xargs); \
	fi; \
	cd $(SERVER_CMD) && go run main.go

mcp: ## 仅启动 MCP Server
	@echo "$(COLOR_BLUE)启动 MCP Server...$(COLOR_RESET)"
	cd $(MCP_CMD) && go run main.go

##@ 前端命令

frontend-dev: ## 启动前端开发服务器 (http://localhost:3000)
	@echo "$(COLOR_BLUE)启动前端开发服务器...$(COLOR_RESET)"
	@if [ ! -d "$(FRONTEND_DIR)/node_modules" ]; then \
		echo "$(COLOR_YELLOW)node_modules 不存在，正在安装依赖...$(COLOR_RESET)"; \
		cd $(FRONTEND_DIR) && npm install --legacy-peer-deps; \
	fi
	@echo "$(COLOR_GREEN)前端开发服务器将在 http://localhost:3000 启动$(COLOR_RESET)"
	@echo "$(COLOR_YELLOW)API 代理到: https://localhost:8443$(COLOR_RESET)"
	cd $(FRONTEND_DIR) && npm start

frontend-build: ## 构建前端生产文件
	@echo "$(COLOR_BLUE)构建前端...$(COLOR_RESET)"
	@if [ ! -d "$(FRONTEND_DIR)/node_modules" ]; then \
		echo "$(COLOR_YELLOW)node_modules 不存在，正在安装依赖...$(COLOR_RESET)"; \
		cd $(FRONTEND_DIR) && npm install --legacy-peer-deps; \
	fi
	cd $(FRONTEND_DIR) && npm run build
	@echo "$(COLOR_GREEN)✓ 前端构建完成: $(FRONTEND_DIR)/build$(COLOR_RESET)"

frontend-install: ## 安装前端依赖
	@echo "$(COLOR_BLUE)安装前端依赖...$(COLOR_RESET)"
	cd $(FRONTEND_DIR) && npm install --legacy-peer-deps
	@echo "$(COLOR_GREEN)✓ 前端依赖安装完成$(COLOR_RESET)"

frontend-clean: ## 清理前端构建文件
	@echo "$(COLOR_BLUE)清理前端构建文件...$(COLOR_RESET)"
	rm -rf $(FRONTEND_DIR)/build
	@echo "$(COLOR_GREEN)✓ 前端清理完成$(COLOR_RESET)"

##@ 构建命令

build: build-server build-mcp ## 构建所有二进制文件

build-server: ## 构建 Server 二进制文件
	@echo "$(COLOR_BLUE)构建 Server...$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	cd $(SERVER_CMD) && go build -o ../../$(SERVER_BIN) main.go
	@echo "$(COLOR_GREEN)✓ Server 构建完成: $(SERVER_BIN)$(COLOR_RESET)"

build-mcp: ## 构建 MCP Server 二进制文件
	@echo "$(COLOR_BLUE)构建 MCP Server...$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	cd $(MCP_CMD) && go build -o ../../$(MCP_BIN) main.go
	@echo "$(COLOR_GREEN)✓ MCP Server 构建完成: $(MCP_BIN)$(COLOR_RESET)"

##@ 测试命令

test: ## 运行测试
	@echo "$(COLOR_BLUE)运行测试...$(COLOR_RESET)"
	cd $(BACKEND_DIR) && go test -v ./...

test-cover: ## 运行测试并生成覆盖率报告
	@echo "$(COLOR_BLUE)运行测试（覆盖率）...$(COLOR_RESET)"
	cd $(BACKEND_DIR) && go test -v -coverprofile=coverage.out ./...
	cd $(BACKEND_DIR) && go tool cover -html=coverage.out -o coverage.html
	@echo "$(COLOR_GREEN)✓ 覆盖率报告: $(BACKEND_DIR)/coverage.html$(COLOR_RESET)"

##@ 清理命令

clean: ## 清理构建文件
	@echo "$(COLOR_BLUE)清理构建文件...$(COLOR_RESET)"
	rm -rf $(BUILD_DIR)
	cd $(BACKEND_DIR) && rm -f coverage.out coverage.html
	@echo "$(COLOR_GREEN)✓ 清理完成$(COLOR_RESET)"

##@ 辅助命令

install-executor: ## 安装 Playwright Executor 依赖
	@echo "$(COLOR_BLUE)安装 Executor 依赖...$(COLOR_RESET)"
	cd $(EXECUTOR_DIR) && npm install
	@echo "$(COLOR_GREEN)✓ Executor 依赖安装完成$(COLOR_RESET)"

deps: ## 安装依赖
	@echo "$(COLOR_BLUE)安装 Go 依赖...$(COLOR_RESET)"
	cd $(BACKEND_DIR) && go mod download
	@$(MAKE) install-executor
	@echo "$(COLOR_GREEN)✓ 所有依赖安装完成$(COLOR_RESET)"

tidy: ## 整理 Go 模块
	@echo "$(COLOR_BLUE)整理 Go 模块...$(COLOR_RESET)"
	cd $(BACKEND_DIR) && go mod tidy
	@echo "$(COLOR_GREEN)✓ 模块整理完成$(COLOR_RESET)"

fmt: ## 格式化代码
	@echo "$(COLOR_BLUE)格式化代码...$(COLOR_RESET)"
	cd $(BACKEND_DIR) && go fmt ./...
	@echo "$(COLOR_GREEN)✓ 代码格式化完成$(COLOR_RESET)"

lint: ## 运行 linter
	@echo "$(COLOR_BLUE)运行 linter...$(COLOR_RESET)"
	cd $(BACKEND_DIR) && golangci-lint run ./...

help: ## 显示帮助信息
	@awk 'BEGIN {FS = ":.*##"; printf "\n使用方法:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
