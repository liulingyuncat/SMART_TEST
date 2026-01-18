@echo off
REM =============================================================================
REM SMART_TEST - Development Environment Script (Windows)
REM 
REM 用法:
REM   dev.bat start        - 启动开发环境（仅后端）
REM   dev.bat start-full   - 启动完整环境（构建前端+启动后端）
REM   dev.bat frontend     - 启动前端开发服务器
REM   dev.bat build        - 构建前端
REM   dev.bat deps         - 启动开发依赖（PostgreSQL + Playwright）
REM   dev.bat stop         - 停止开发环境
REM   dev.bat restart      - 重启开发环境
REM   dev.bat clean        - 清理开发环境（包括数据）
REM   dev.bat logs         - 查看 Docker 日志
REM =============================================================================

setlocal enabledelayedexpansion

set BACKEND_DIR=backend
set FRONTEND_DIR=frontend
set CERT_DIR=certs
set BUILD_DIR=bin

if "%1"=="" (
    echo 用法: dev.bat [start^|start-full^|frontend^|build^|deps^|stop^|restart^|clean^|logs]
    exit /b 1
)

if "%1"=="start" goto start
if "%1"=="start-full" goto start-full
if "%1"=="frontend" goto frontend
if "%1"=="build" goto build-frontend
if "%1"=="deps" goto deps
if "%1"=="stop" goto stop
if "%1"=="restart" goto restart
if "%1"=="clean" goto clean
if "%1"=="logs" goto logs
echo 错误: 未知命令 "%1"
exit /b 1

:start
echo [启动开发环境 - 仅后端]
echo.
echo 注意: 本命令仅启动后端服务，Docker 和前端需单独启动
echo   - Docker: dev.bat deps
echo   - 前端: dev.bat frontend
echo   或使用 'dev.bat start-full' 一键启动全部服务
echo.

REM 检查证书
echo [1/3] 检查 HTTPS 证书...
if not exist "%CERT_DIR%\server.crt" (
    echo 证书不存在，正在生成...
    if not exist "%CERT_DIR%" mkdir "%CERT_DIR%"
    if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"
    
    cd %BACKEND_DIR%\cmd\gencert
    go build -o ..\..\..\%BUILD_DIR%\gencert.exe main.go
    cd ..\..\..
    %BUILD_DIR%\gencert.exe
    
    if errorlevel 1 (
        echo 错误: 证书生成失败
        exit /b 1
    )
    echo 证书生成完成
) else (
    echo 证书已存在
)

REM 启动 Server
echo [2/3] 启动 Web Server...
cd %BACKEND_DIR%

REM 加载 .env 环境变量
set ENV_VARS=
if exist "..\.env" (
    echo ✓ 从 .env 加载环境变量
    for /f "usebackq tokens=* delims=" %%a in ("..\.env") do (
        set line=%%a
        REM 跳过注释和空行
        if not "!line:~0,1!"=="#" if not "!line!"=="" (
            set ENV_VARS=!ENV_VARS! ^& set %%a
        )
    )
) else (
    echo ⚠ .env 文件不存在，使用默认配置
)

start "SMART_TEST Server" /min cmd /c "!ENV_VARS! ^& go run .\cmd\server\main.go"
cd ..

timeout /t 2 /nobreak >nul

REM 启动 MCP Server
echo [3/3] 启动 MCP Server...
cd %BACKEND_DIR%
start "SMART_TEST MCP" /min cmd /c "!ENV_VARS! ^& go run .\cmd\mcp\main.go"
cd ..

echo.
echo ✓ 后端服务已启动
echo   - Web Server:  https://localhost:8443
echo   - MCP Server:  http://localhost:16410
echo   - Playwright:  http://localhost:53730
echo.
echo 前端开发服务器: 在另一个终端运行 'dev.bat frontend'
echo 或使用 'dev.bat start-full' 一键启动全部服务
echo.
echo 使用 'dev.bat stop' 停止服务
goto end

:deps
echo [启动开发依赖服务]
echo.
docker compose -f docker-compose.dev.yml up -d
if errorlevel 1 (
    echo 错误: Docker 容器启动失败
    exit /b 1
)
echo.
echo ✓ 开发依赖已启动
echo   - PostgreSQL:           localhost:5432
echo   - Playwright Runner:    ws://localhost:53729/
echo   - Playwright Executor:  http://localhost:53730/
echo.
goto end

:start-full
echo [启动完整开发环境]
echo.

REM 启动 Docker 依赖
echo [1/5] 启动 Docker 容器...
docker compose -f docker-compose.dev.yml up -d
if errorlevel 1 (
    echo 错误: Docker 容器启动失败
    exit /b 1
)

REM 检查证书
echo [2/5] 检查 HTTPS 证书...
if not exist "%CERT_DIR%\server.crt" (
    echo 证书不存在，正在生成...
    if not exist "%CERT_DIR%" mkdir "%CERT_DIR%"
    if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"
    
    cd %BACKEND_DIR%\cmd\gencert
    go build -o ..\..\..\%BUILD_DIR%\gencert.exe main.go
    cd ..\..\..
    %BUILD_DIR%\gencert.exe
    
    if errorlevel 1 (
        echo 错误: 证书生成失败
        exit /b 1
    )
    echo 证书生成完成
) else (
    echo 证书已存在
)

REM 构建前端
echo [3/5] 构建前端...
cd %FRONTEND_DIR%
if not exist "node_modules" (
    echo 安装前端依赖...
    call npm install --legacy-peer-deps
)
call npm run build
if errorlevel 1 (
    echo 错误: 前端构建失败
    exit /b 1
)
cd ..

REM 启动 Server
echo [4/5] 启动 Web Server...
cd %BACKEND_DIR%

REM 加载 .env 环境变量
set ENV_VARS=
if exist "..\.env" (
    echo ✓ 从 .env 加载环境变量
    for /f "usebackq tokens=* delims=" %%a in ("..\.env") do (
        set line=%%a
        REM 跳过注释和空行
        if not "!line:~0,1!"=="#" if not "!line!"=="" (
            set ENV_VARS=!ENV_VARS! ^& set %%a
        )
    )
) else (
    echo ⚠ .env 文件不存在，使用默认配置
)

start "SMART_TEST Server" /min cmd /c "!ENV_VARS! ^& go run .\cmd\server\main.go"
cd ..

timeout /t 2 /nobreak >nul

REM 启动 MCP Server
echo [5/5] 启动 MCP Server...
cd %BACKEND_DIR%
start "SMART_TEST MCP" /min cmd /c "!ENV_VARS! ^& go run .\cmd\mcp\main.go"
cd ..

echo.
echo ✓ 完整开发环境已启动
echo   - Frontend:    https://localhost:8443
echo   - MCP Server:  http://localhost:16410
echo   - Playwright:  http://localhost:53730
echo.
echo 使用 'dev.bat stop' 停止服务
goto end

:frontend
echo [启动前端开发服务器]
echo.
cd %FRONTEND_DIR%
if not exist "node_modules" (
    echo 安装前端依赖...
    call npm install --legacy-peer-deps
)
echo.
echo 前端开发服务器将在 http://localhost:3000 启动
echo API 代理到: https://localhost:8443
echo.
call npm start
goto end

:build-frontend
echo [构建前端]
echo.
cd %FRONTEND_DIR%
if not exist "node_modules" (
    echo 安装前端依赖...
    call npm install --legacy-peer-deps
)
call npm run build
if errorlevel 1 (
    echo 错误: 前端构建失败
    exit /b 1
)
cd ..
echo.
echo ✓ 前端构建完成: %FRONTEND_DIR%\build
goto end

:stop
echo [停止开发环境]
echo.

echo 停止 Go 进程...
taskkill /FI "WindowTitle eq SMART_TEST Server*" /F >nul 2>&1
taskkill /FI "WindowTitle eq SMART_TEST MCP*" /F >nul 2>&1
timeout /t 1 /nobreak >nul

echo 停止 Docker 容器...
docker compose -f docker-compose.dev.yml down

echo.
echo ✓ 开发环境已停止
goto end

:restart
echo [重启开发环境]
call %0 stop
timeout /t 2 /nobreak >nul
call %0 start
goto end

:clean
echo [清理开发环境]
echo.

call %0 stop

echo 清理 Docker 数据...
docker compose -f docker-compose.dev.yml down -v

echo 清理 SQLite 数据库...
if exist "%BACKEND_DIR%\webtest.db" del /q "%BACKEND_DIR%\webtest.db"
if exist "%BACKEND_DIR%\webtest.db-shm" del /q "%BACKEND_DIR%\webtest.db-shm"
if exist "%BACKEND_DIR%\webtest.db-wal" del /q "%BACKEND_DIR%\webtest.db-wal"

echo.
echo ✓ 开发环境已清理
goto end

:logs
echo [查看 Docker 日志]
docker compose -f docker-compose.dev.yml logs -f
goto end

:end
endlocal
