@echo off
chcp 65001 > nul
setlocal enabledelayedexpansion

:: Docker Compose Windows 启动脚本
:: 包含完整的健康检查和初始化验证

echo 🐳 启动 Docker Compose 服务...
echo ========================================

:: 检查 Docker 和 Docker Compose 是否可用
docker --version > nul 2>&1
if !errorlevel! neq 0 (
    echo ❌ Docker 未安装或未启动
    pause
    exit /b 1
)

docker-compose --version > nul 2>&1
if !errorlevel! neq 0 (
    echo ❌ Docker Compose 未安装
    pause
    exit /b 1
)

:: 清理可能存在的旧容器
echo 🧹 清理旧容器...
docker-compose down --volumes --remove-orphans

:: 构建并启动服务
echo 🔨 构建并启动服务...
docker-compose up --build -d

:: 等待 MySQL 健康检查通过
echo ⏳ 等待 MySQL 数据库启动...
set /a timeout=300
set /a elapsed=0
set /a interval=10

:mysql_wait_loop
if !elapsed! geq !timeout! (
    echo ❌ MySQL 启动超时
    docker-compose logs mysql
    pause
    exit /b 1
)

docker-compose exec mysql mysqladmin ping -h localhost -u app -ppassword --silent > nul 2>&1
if !errorlevel! equ 0 (
    echo ✅ MySQL 数据库已启动
    goto mysql_ready
)

echo ⏳ MySQL 仍在启动中... (!elapsed!s/!timeout!s)
timeout /t !interval! /nobreak > nul
set /a elapsed=!elapsed!+!interval!
goto mysql_wait_loop

:mysql_ready

:: 验证数据库初始化
echo 🔍 验证数据库初始化...
docker-compose exec mysql mysql -u app -ppassword static_hosting -e "SELECT 'Database verification' as check_type, (SELECT COUNT(*) FROM articles) as articles_count, (SELECT COUNT(*) FROM users) as users_count, (SELECT COUNT(*) FROM api_keys) as api_keys_count;"

:: 等待 Web 服务健康检查通过
echo ⏳ 等待 Web 服务启动...
set /a timeout=120
set /a elapsed=0

:web_wait_loop
if !elapsed! geq !timeout! (
    echo ❌ Web 服务启动超时
    docker-compose logs web
    pause
    exit /b 1
)

curl -s http://localhost:8080/ > nul 2>&1
if !errorlevel! equ 0 (
    echo ✅ Web 服务已启动
    goto web_ready
)

echo ⏳ Web 服务仍在启动中... (!elapsed!s/!timeout!s)
timeout /t !interval! /nobreak > nul
set /a elapsed=!elapsed!+!interval!
goto web_wait_loop

:web_ready

:: 测试 API 端点
echo 🧪 测试 API 端点...
curl -s -H "X-API-Key: demo-api-key-12345" http://localhost:8080/ > nul 2>&1
if !errorlevel! equ 0 (
    echo ✅ API 服务工作正常
) else (
    echo ⚠️ API 可能有问题，请检查日志
)

:: 显示服务状态
echo.
echo 🎉 所有服务启动完成！
echo ========================================
echo 📱 Web 界面: http://localhost:8080
echo 🔧 管理后台: http://localhost:8080/admin
echo 🔑 默认账号: admin / password
echo 🚀 API 密钥: demo-api-key-12345
echo 🗄️ 数据库端口: localhost:3306
echo.
echo 📊 服务状态:
docker-compose ps

echo.
echo 📝 查看日志: docker-compose logs -f
echo 🛑 停止服务: docker-compose down
echo 🔄 重启服务: docker-compose restart
echo.
echo 按任意键退出...
pause > nul
