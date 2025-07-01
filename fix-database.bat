@echo off
chcp 65001 >nul
echo ======================================
echo   数据库连接问题快速修复
echo ======================================
echo.

echo [1/5] 停止现有 MySQL 容器...
docker stop static-hosting-mysql >nul 2>&1
docker rm static-hosting-mysql >nul 2>&1

echo [2/5] 清理网络和卷...
docker network prune -f >nul 2>&1

echo [3/5] 重新创建 MySQL 容器...
docker run -d --name static-hosting-mysql ^
  -e MYSQL_ROOT_PASSWORD=rootpassword ^
  -e MYSQL_DATABASE=static_hosting ^
  -e MYSQL_USER=app ^
  -e MYSQL_PASSWORD=password ^
  -p 127.0.0.1:3306:3306 ^
  --health-cmd="mysqladmin ping -h localhost -u app -ppassword" ^
  --health-interval=10s ^
  --health-timeout=5s ^
  --health-retries=5 ^
  mysql:8.0 --default-authentication-plugin=mysql_native_password

if errorlevel 1 (
    echo [错误] MySQL 容器创建失败
    pause
    exit /b 1
)

echo [4/5] 等待 MySQL 完全启动...
echo 等待时间可能需要 1-2 分钟，请耐心等待...

REM 等待健康检查通过
set /a count=0
:wait_healthy
set /a count+=1
if %count% gtr 60 (
    echo [错误] MySQL 启动超时
    echo 查看日志: docker logs static-hosting-mysql
    pause
    exit /b 1
)

docker inspect static-hosting-mysql --format="{{.State.Health.Status}}" 2>nul | findstr "healthy" >nul
if errorlevel 1 (
    echo 等待中... [%count%/60]
    timeout /t 2 /nobreak >nul
    goto wait_healthy
)

echo [成功] MySQL 健康检查通过！

echo [5/5] 验证数据库连接...
docker exec static-hosting-mysql mysql -u app -ppassword -e "SELECT 1;" >nul 2>&1
if errorlevel 1 (
    echo [错误] 数据库验证失败
    docker logs static-hosting-mysql --tail 10
    pause
    exit /b 1
) else (
    echo [成功] 数据库连接验证通过！
)

echo.
echo ======================================
echo   修复完成！
echo ======================================
echo 现在可以重新运行应用:
echo   .\server.exe
echo.
echo 或者运行完整启动脚本:
echo   start-local.bat
echo.

pause
