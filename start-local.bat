@echo off
chcp 65001 >nul
echo ======================================
echo   本地构建启动脚本 (网络问题备用方案)
echo ======================================
echo.

echo [1/4] 检查 Go 环境...
go version >nul 2>&1
if errorlevel 1 (
    echo [错误] Go 未安装，请先安装 Go 1.21+
    echo 下载地址: https://golang.org/dl/
    pause
    exit /b 1
) else (
    echo [成功] Go 环境检查通过
)

echo.
echo [2/4] 下载依赖...
echo 设置 Go 代理为国内镜像...
set GOPROXY=https://goproxy.cn,direct
set GO111MODULE=on

echo 正在下载项目依赖...
go mod tidy
if errorlevel 1 (
    echo [错误] 依赖下载失败
    pause
    exit /b 1
) else (
    echo [成功] 依赖下载完成
)

echo.
echo [3/4] 编译应用...
echo 正在编译静态网页托管服务器...
go build -o server.exe ./cmd/server
if errorlevel 1 (
    echo [错误] 编译失败
    pause
    exit /b 1
) else (
    echo [成功] 编译完成
)

echo.
echo [4/4] 启动 MySQL 数据库...
echo 检查是否已有 MySQL 容器运行...
docker ps | findstr static-hosting-mysql >nul 2>&1
if not errorlevel 1 (
    echo [信息] MySQL 容器已在运行
) else (
    echo 启动 MySQL 容器...
    docker run -d --name static-hosting-mysql ^
      -e MYSQL_ROOT_PASSWORD=rootpassword ^
      -e MYSQL_DATABASE=static_hosting ^
      -e MYSQL_USER=app ^
      -e MYSQL_PASSWORD=password ^
      -p 3306:3306 ^
      --health-cmd="mysqladmin ping -h localhost -u app -ppassword" ^
      --health-interval=10s ^
      --health-timeout=5s ^
      --health-retries=5 ^
      mysql:8.0

    if errorlevel 1 (
        echo [警告] MySQL 容器启动失败，可能已存在
        echo 尝试删除旧容器并重新创建...
        docker rm -f static-hosting-mysql >nul 2>&1
        docker run -d --name static-hosting-mysql ^
          -e MYSQL_ROOT_PASSWORD=rootpassword ^
          -e MYSQL_DATABASE=static_hosting ^
          -e MYSQL_USER=app ^
          -e MYSQL_PASSWORD=password ^
          -p 3306:3306 ^
          --health-cmd="mysqladmin ping -h localhost -u app -ppassword" ^
          --health-interval=10s ^
          --health-timeout=5s ^
          --health-retries=5 ^
          mysql:8.0
    )
)

echo 等待 MySQL 完全启动...
echo 这可能需要 30-60 秒，请耐心等待...

REM 循环检查 MySQL 是否准备就绪
set /a count=0
:check_mysql
set /a count+=1
if %count% gtr 30 (
    echo [错误] MySQL 启动超时
    echo 请检查 Docker 状态: docker logs static-hosting-mysql
    pause
    exit /b 1
)

echo 检查 MySQL 连接状态 [%count%/30]...
docker exec static-hosting-mysql mysqladmin ping -h localhost -u app -ppassword >nul 2>&1
if errorlevel 1 (
    timeout /t 2 /nobreak >nul
    goto check_mysql
) else (
    echo [成功] MySQL 已准备就绪！
)

echo.
echo [完成] 准备工作完成！
echo.
echo ===== 启动应用 =====
echo 运行以下命令启动应用：
echo.
echo   .\server.exe
echo.
echo ===== 访问信息 =====
echo 管理后台: http://localhost:8080/admin
echo 默认账号: admin / admin123
echo API 接口: http://localhost:8080/api
echo.
echo ===== 停止服务 =====
echo 停止应用: Ctrl+C
echo 停止数据库: docker stop static-hosting-mysql
echo.

choice /c yn /m "现在启动应用吗? (Y/N)"
if errorlevel 2 goto :end

echo.
echo 正在启动应用...
.\server.exe

:end
echo.
echo 脚本执行完成。
pause
