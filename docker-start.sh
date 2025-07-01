#!/bin/bash

# Docker Compose 启动脚本
# 包含完整的健康检查和初始化验证

set -e

echo "🐳 启动 Docker Compose 服务..."
echo "========================================"

# 检查 Docker 和 Docker Compose 是否可用
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装或未启动"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装"
    exit 1
fi

# 清理可能存在的旧容器
echo "🧹 清理旧容器..."
docker-compose down --volumes --remove-orphans

# 构建并启动服务
echo "🔨 构建并启动服务..."
docker-compose up --build -d

# 等待 MySQL 健康检查通过
echo "⏳ 等待 MySQL 数据库启动..."
timeout=300  # 5分钟超时
elapsed=0
interval=10

while [ $elapsed -lt $timeout ]; do
    if docker-compose exec mysql mysqladmin ping -h localhost -u app -ppassword --silent; then
        echo "✅ MySQL 数据库已启动"
        break
    fi
    echo "⏳ MySQL 仍在启动中... (${elapsed}s/${timeout}s)"
    sleep $interval
    elapsed=$((elapsed + interval))
done

if [ $elapsed -ge $timeout ]; then
    echo "❌ MySQL 启动超时"
    docker-compose logs mysql
    exit 1
fi

# 验证数据库初始化
echo "🔍 验证数据库初始化..."
docker-compose exec mysql mysql -u app -ppassword static_hosting -e "
SELECT 
    'Database verification' as check_type,
    (SELECT COUNT(*) FROM articles) as articles_count,
    (SELECT COUNT(*) FROM users) as users_count,
    (SELECT COUNT(*) FROM api_keys) as api_keys_count;
"

# 等待 Web 服务健康检查通过
echo "⏳ 等待 Web 服务启动..."
timeout=120  # 2分钟超时
elapsed=0

while [ $elapsed -lt $timeout ]; do
    if curl -s http://localhost:8080/ > /dev/null; then
        echo "✅ Web 服务已启动"
        break
    fi
    echo "⏳ Web 服务仍在启动中... (${elapsed}s/${timeout}s)"
    sleep $interval
    elapsed=$((elapsed + interval))
done

if [ $elapsed -ge $timeout ]; then
    echo "❌ Web 服务启动超时"
    docker-compose logs web
    exit 1
fi

# 测试 API 端点
echo "🧪 测试 API 端点..."
if curl -s -H "X-API-Key: demo-api-key-12345" http://localhost:8080/api/articles/welcome > /dev/null; then
    echo "✅ API 认证工作正常"
else
    echo "⚠️ API 认证可能有问题，请检查日志"
fi

# 显示服务状态
echo ""
echo "🎉 所有服务启动完成！"
echo "========================================"
echo "📱 Web 界面: http://localhost:8080"
echo "🔧 管理后台: http://localhost:8080/admin"
echo "🔑 默认账号: admin / password"
echo "🚀 API 密钥: demo-api-key-12345"
echo "🗄️ 数据库端口: localhost:3306"
echo ""
echo "📊 服务状态:"
docker-compose ps

echo ""
echo "📝 查看日志: docker-compose logs -f"
echo "🛑 停止服务: docker-compose down"
echo "🔄 重启服务: docker-compose restart"
