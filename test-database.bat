@echo off
chcp 65001 >nul
echo ======================================
echo   数据库连接测试工具
echo ======================================
echo.

echo [1/4] 检查 MySQL 容器状态...
docker ps | findstr static-hosting-mysql >nul 2>&1
if errorlevel 1 (
    echo [错误] MySQL 容器未运行
    echo 请先运行: start-local.bat
    pause
    exit /b 1
) else (
    echo [成功] MySQL 容器正在运行
)

echo.
echo [2/4] 检查端口连接...
netstat -an | findstr ":3306" >nul 2>&1
if errorlevel 1 (
    echo [警告] 端口 3306 未监听
) else (
    echo [成功] 端口 3306 正在监听
)

echo.
echo [3/4] 测试数据库连接...
docker exec static-hosting-mysql mysqladmin ping -h localhost -u app -ppassword >nul 2>&1
if errorlevel 1 (
    echo [错误] 数据库连接失败
    echo 检查数据库日志:
    docker logs static-hosting-mysql --tail 10
    echo.
) else (
    echo [成功] 数据库连接正常
)

echo.
echo [4/4] 测试数据库访问...
docker exec static-hosting-mysql mysql -u app -ppassword -e "SHOW DATABASES;" >nul 2>&1
if errorlevel 1 (
    echo [错误] 数据库访问失败
) else (
    echo [成功] 数据库访问正常
    echo.
    echo 数据库列表:
    docker exec static-hosting-mysql mysql -u app -ppassword -e "SHOW DATABASES;"
)

echo.
echo ======================================
echo   连接信息
echo ======================================
echo 主机: localhost
echo 端口: 3306
echo 用户: app
echo 密码: password
echo 数据库: static_hosting
echo.

echo ======================================
echo   故障排除建议
echo ======================================
echo 如果连接失败，请尝试:
echo 1. 重启 MySQL 容器: docker restart static-hosting-mysql
echo 2. 删除并重新创建: docker rm -f static-hosting-mysql
echo 3. 检查端口占用: netstat -an ^| findstr :3306
echo 4. 查看详细日志: docker logs static-hosting-mysql
echo.

pause
