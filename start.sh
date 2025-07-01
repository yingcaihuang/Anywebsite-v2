#!/bin/bash

# 静态网页托管服务器启动脚本

echo "🚀 启动静态网页托管服务器..."

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

# 创建必要的目录
echo "📁 创建目录..."
mkdir -p static uploads certs

# 构建并启动服务
echo "🔨 构建并启动服务..."
docker-compose up -d --build

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 10

# 检查服务状态
if docker-compose ps | grep -q "Up"; then
    echo "✅ 服务启动成功！"
    echo ""
    echo "📌 访问地址："
    echo "   管理后台: http://localhost:8080/admin"
    echo "   默认账号: admin / admin123"
    echo "   API文档:  http://localhost:8080/api"
    echo "   示例文章: http://localhost:8080/p/welcome"
    echo ""
    echo "🔑 API密钥: demo-api-key-12345"
    echo ""
    echo "📊 查看日志: docker-compose logs -f"
    echo "🛑 停止服务: docker-compose down"
else
    echo "❌ 服务启动失败，请检查日志："
    echo "docker-compose logs"
fi
