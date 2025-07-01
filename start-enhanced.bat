@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

REM ====================================
REM 静态网页托管服务器启动脚本 (Windows)
REM ====================================

echo.
echo ========================================
echo   静态网页托管服务器
echo ========================================
echo.

REM 检查 Docker 是否安装
echo [1/5] 检查 Docker 环境...
docker --version >nul 2>&1
if errorlevel 1 (
    echo [错误] Docker 未安装或未启动
    echo 请先安装并启动 Docker Desktop
    echo 下载地址: https://www.docker.com/products/docker-desktop
    echo.
    pause
    exit /b 1
) else (
    echo [成功] Docker 已安装
)

docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo [错误] Docker Compose 未安装
    echo 请确保 Docker Desktop 包含 Docker Compose
    echo.
    pause
    exit /b 1
) else (
    echo [成功] Docker Compose 已安装
)

REM 创建必要的目录
echo.
echo [2/5] 创建项目目录...
if not exist "static" (
    mkdir static
    echo [创建] static 目录
)
if not exist "uploads" (
    mkdir uploads
    echo [创建] uploads 目录
)
if not exist "certs" (
    mkdir certs
    echo [创建] certs 目录
)

REM 检查配置文件
echo.
echo [3/5] 检查配置文件...
if not exist "docker-compose.yml" (
    echo [错误] 未找到 docker-compose.yml 文件
    echo 请确保在项目根目录运行此脚本
    pause
    exit /b 1
) else (
    echo [成功] 配置文件检查完成
)

REM 停止现有服务（如果有）
echo.
echo [4/5] 清理现有服务...
docker-compose down >nul 2>&1

REM 构建并启动服务
echo.
echo [5/5] 构建并启动服务...
echo 正在构建镜像，首次运行可能需要几分钟...
docker-compose up -d --build

if errorlevel 1 (
    echo.
    echo [错误] 服务启动失败！
    echo 请检查错误信息并尝试以下解决方案：
    echo 1. 确保 Docker Desktop 正在运行
    echo 2. 检查端口 8080 和 3306 是否被占用
    echo 3. 查看详细日志: docker-compose logs
    echo.
    pause
    exit /b 1
)

REM 等待服务完全启动
echo.
echo 等待服务启动...
timeout /t 15 /nobreak >nul

REM 检查服务状态
echo 检查服务状态...
docker-compose ps | findstr "Up" >nul
if errorlevel 1 (
    echo.
    echo [警告] 服务可能未完全启动
    echo 请查看日志: docker-compose logs
    echo.
) else (
    echo.
    echo ========================================
    echo   服务启动成功！
    echo ========================================
    echo.
    echo [访问地址]
    echo   管理后台: http://localhost:8080/admin
    echo   示例文章: http://localhost:8080/p/welcome
    echo.
    echo [登录信息]
    echo   用户名: admin
    echo   密码:   admin123
    echo.
    echo [API 信息]
    echo   API 地址: http://localhost:8080/api/articles
    echo   API 密钥: demo-api-key-12345
    echo.
    echo [常用命令]
    echo   查看日志: docker-compose logs -f
    echo   停止服务: docker-compose down
    echo   重启服务: docker-compose restart
    echo.
    echo 提示: 首次启动会自动创建数据库和示例数据
    echo.
)

echo 按任意键退出...
pause >nul
