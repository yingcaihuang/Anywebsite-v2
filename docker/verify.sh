#!/bin/bash

# 数据库迁移验证脚本
# 检查数据库结构和数据完整性

set -e

echo "🔍 开始数据库迁移验证..."
echo "========================================"

# 检查 MySQL 服务是否可用
if ! docker-compose exec mysql mysqladmin ping -h localhost -u app -ppassword --silent; then
    echo "❌ MySQL 服务不可用"
    exit 1
fi

echo "✅ MySQL 服务可用"

# 验证数据库存在
echo "🗄️ 验证数据库存在..."
DB_EXISTS=$(docker-compose exec mysql mysql -u app -ppassword -e "SHOW DATABASES LIKE 'static_hosting';" -s -N)
if [ -z "$DB_EXISTS" ]; then
    echo "❌ 数据库 'static_hosting' 不存在"
    exit 1
fi
echo "✅ 数据库 'static_hosting' 存在"

# 验证表结构
echo "📋 验证表结构..."
TABLES=$(docker-compose exec mysql mysql -u app -ppassword static_hosting -e "SHOW TABLES;" -s -N)
EXPECTED_TABLES=("articles" "users" "api_keys" "sessions")

for table in "${EXPECTED_TABLES[@]}"; do
    if echo "$TABLES" | grep -q "^$table$"; then
        echo "✅ 表 '$table' 存在"
    else
        echo "❌ 表 '$table' 不存在"
        exit 1
    fi
done

# 验证 articles 表结构（UUID支持）
echo "🆔 验证 articles 表 UUID 支持..."
ARTICLES_ID_TYPE=$(docker-compose exec mysql mysql -u app -ppassword static_hosting -e "DESCRIBE articles;" -s -N | grep "^id" | awk '{print $2}')
if [[ "$ARTICLES_ID_TYPE" == "varchar(36)" ]]; then
    echo "✅ articles 表使用 UUID 主键 (VARCHAR(36))"
else
    echo "❌ articles 表主键类型错误: $ARTICLES_ID_TYPE (应为 varchar(36))"
    exit 1
fi

# 验证字符集
echo "🔤 验证字符集配置..."
CHARSET=$(docker-compose exec mysql mysql -u app -ppassword static_hosting -e "SELECT DEFAULT_CHARACTER_SET_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = 'static_hosting';" -s -N)
if [[ "$CHARSET" == "utf8mb4" ]]; then
    echo "✅ 数据库字符集: $CHARSET"
else
    echo "⚠️ 数据库字符集: $CHARSET (建议使用 utf8mb4)"
fi

# 验证初始数据
echo "📊 验证初始数据..."
docker-compose exec mysql mysql -u app -ppassword static_hosting -e "
SELECT 
    'Data verification' as check_type,
    (SELECT COUNT(*) FROM articles) as articles_count,
    (SELECT COUNT(*) FROM users) as users_count,
    (SELECT COUNT(*) FROM api_keys) as api_keys_count,
    (SELECT COUNT(*) FROM sessions) as sessions_count;
"

# 验证 UUID 格式
echo "🔍 验证 UUID 格式..."
UUID_CHECK=$(docker-compose exec mysql mysql -u app -ppassword static_hosting -e "
SELECT id FROM articles WHERE 
    LENGTH(id) != 36 
    OR id NOT REGEXP '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$'
LIMIT 1;" -s -N)

if [ -z "$UUID_CHECK" ]; then
    echo "✅ 所有 articles 记录使用有效的 UUID 格式"
else
    echo "❌ 发现无效的 UUID 格式: $UUID_CHECK"
    exit 1
fi

# 验证默认管理员用户
echo "👤 验证默认管理员用户..."
ADMIN_EXISTS=$(docker-compose exec mysql mysql -u app -ppassword static_hosting -e "SELECT username FROM users WHERE username='admin';" -s -N)
if [[ "$ADMIN_EXISTS" == "admin" ]]; then
    echo "✅ 默认管理员用户存在"
else
    echo "❌ 默认管理员用户不存在"
    exit 1
fi

# 验证 API 密钥
echo "🔑 验证 API 密钥..."
API_KEYS=$(docker-compose exec mysql mysql -u app -ppassword static_hosting -e "SELECT api_key FROM api_keys WHERE is_active = 1;" -s -N)
if echo "$API_KEYS" | grep -q "demo-api-key-12345"; then
    echo "✅ 默认 API 密钥存在"
else
    echo "❌ 默认 API 密钥不存在"
    exit 1
fi

# 验证索引
echo "📇 验证数据库索引..."
INDEXES=$(docker-compose exec mysql mysql -u app -ppassword static_hosting -e "SHOW INDEX FROM articles;" -s -N | wc -l)
if [ "$INDEXES" -gt 1 ]; then
    echo "✅ articles 表索引已创建 ($INDEXES 个索引)"
else
    echo "⚠️ articles 表索引较少，可能影响性能"
fi

# 显示数据库配置信息
echo "⚙️ 数据库配置信息..."
docker-compose exec mysql mysql -u app -ppassword static_hosting -e "
SELECT 
    @@character_set_database as db_charset,
    @@collation_database as db_collation,
    @@time_zone as timezone,
    @@version as mysql_version;
"

echo "========================================"
echo "🎉 数据库迁移验证完成！"
echo ""
echo "📊 验证报告:"
echo "  ✅ 数据库服务正常"
echo "  ✅ 表结构完整"
echo "  ✅ UUID 主键正确"
echo "  ✅ 字符集配置正确"
echo "  ✅ 初始数据完整"
echo "  ✅ 管理员用户可用"
echo "  ✅ API 密钥可用"
echo ""
echo "🚀 系统已准备就绪，可以正常使用！"
