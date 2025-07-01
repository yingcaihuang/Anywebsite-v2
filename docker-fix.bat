@echo off
chcp 65001 >nul
echo ======================================
echo   Docker 网络连接诊断工具
echo ======================================
echo.

echo [1/6] 检查 Docker 状态...
docker --version >nul 2>&1
if errorlevel 1 (
    echo [错误] Docker 未安装或未启动
    goto :end
) else (
    echo [成功] Docker 已安装并运行
)

echo.
echo [2/6] 检查网络连接...
ping -n 1 8.8.8.8 >nul 2>&1
if errorlevel 1 (
    echo [警告] 网络连接可能有问题
) else (
    echo [成功] 网络连接正常
)

echo.
echo [3/6] 测试 DNS 解析...
nslookup hub.docker.com >nul 2>&1
if errorlevel 1 (
    echo [警告] DNS 解析可能有问题
    echo 建议更换 DNS 为 8.8.8.8 或 114.114.114.114
) else (
    echo [成功] DNS 解析正常
)

echo.
echo [4/6] 清理 Docker 缓存...
echo 正在清理构建缓存...
docker builder prune -f >nul 2>&1
echo [完成] Docker 缓存已清理

echo.
echo [5/6] 测试镜像拉取...
echo 正在测试镜像下载...
docker pull hello-world >nul 2>&1
if errorlevel 1 (
    echo [错误] 镜像下载失败，可能需要配置镜像加速器
    echo 请查看 DOCKER_NETWORK_FIX.md 文档
) else (
    echo [成功] 镜像下载正常
    docker rmi hello-world >nul 2>&1
)

echo.
echo [6/6] 尝试重新构建...
echo 使用优化配置重新构建项目...
docker-compose build --no-cache

if errorlevel 1 (
    echo.
    echo [失败] 构建仍然失败
    echo 解决建议：
    echo 1. 检查网络连接和防火墙设置
    echo 2. 配置 Docker 镜像加速器
    echo 3. 查看详细文档: DOCKER_NETWORK_FIX.md
    echo 4. 如果问题持续，考虑使用本地构建方式
) else (
    echo.
    echo [成功] 构建完成！
    echo 现在可以运行: docker-compose up -d
)

:end
echo.
echo 按任意键退出...
pause >nul
