@echo off
chcp 65001 >nul
REM 静态网页托管服务器启动脚本 (Windows)

echo 启动静态网页托管服务器...

REM 检查 Docker 是否安装
docker --version >nul 2>&1
if errorlevel 1 (
    echo [错误] Docker 未安装，请先安装 Docker Desktop
    echo 下载地址: https://www.docker.com/products/docker-desktop
    pause
    exit /b 1
)

docker-compose --version >nul 2>&1
if errorlevel 1 (
    echo [错误] Docker Compose 未安装，请先安装 Docker Compose
    pause
    exit /b 1
)

REM 创建必要的目录
echo [信息] 创建目录...
if not exist "static" mkdir static
if not exist "uploads" mkdir uploads
if not exist "certs" mkdir certs

REM 构建并启动服务
echo [信息] 构建并启动服务...
docker-compose up -d --build

REM 等待服务启动
echo [信息] 等待服务启动...
timeout /t 10 /nobreak >nul

REM 检查服务状态
docker-compose ps | findstr "Up" >nul
if errorlevel 1 (
    echo [错误] 服务启动失败，请检查日志：
    echo   命令: docker-compose logs
    pause
    exit /b 1
) else (
    echo [成功] 服务启动成功！
    echo.
    echo ===== 访问地址 =====
    echo 管理后台: http://localhost:8080/admin
    echo 默认账号: admin / admin123
    echo API 接口: http://localhost:8080/api
    echo 示例文章: http://localhost:8080/p/welcome
    echo.
    echo ===== API 信息 =====
    echo API 密钥: demo-api-key-12345
    echo.
    echo ===== 常用命令 =====
    echo 查看日志: docker-compose logs -f
    echo 停止服务: docker-compose down
    echo.
    pause
)
