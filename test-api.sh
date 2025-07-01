#!/bin/bash

# API 测试脚本

API_BASE="http://localhost:8080/api"
API_KEY="demo-api-key-12345"

echo "🧪 静态网页托管服务器 API 测试"
echo "=================================="

# 测试创建文章
echo "📝 测试创建文章..."
RESPONSE=$(curl -s -X POST "$API_BASE/articles" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $API_KEY" \
  -d '{
    "title": "API测试文章",
    "content": "<h1>这是通过API创建的文章</h1><p>内容包含HTML格式。</p><ul><li>支持列表</li><li>支持链接</li></ul>",
    "slug": "api-test-article",
    "status": "published"
  }')

echo "响应: $RESPONSE"
echo ""

# 解析文章ID
ARTICLE_ID=$(echo $RESPONSE | grep -o '"id":[0-9]*' | grep -o '[0-9]*')
echo "创建的文章ID: $ARTICLE_ID"
echo ""

# 测试获取文章
if [ ! -z "$ARTICLE_ID" ]; then
    echo "📖 测试获取文章..."
    curl -s -X GET "$API_BASE/articles/$ARTICLE_ID" \
      -H "X-API-Key: $API_KEY" | jq '.'
    echo ""
fi

# 测试更新文章
if [ ! -z "$ARTICLE_ID" ]; then
    echo "✏️ 测试更新文章..."
    curl -s -X PUT "$API_BASE/articles/$ARTICLE_ID" \
      -H "Content-Type: application/json" \
      -H "X-API-Key: $API_KEY" \
      -d '{
        "title": "更新后的API测试文章",
        "content": "<h1>文章已更新</h1><p>这是更新后的内容。</p>",
        "status": "published"
      }' | jq '.'
    echo ""
fi

# 测试获取文章列表
echo "📋 测试获取文章列表..."
curl -s -X GET "$API_BASE/articles?limit=5" \
  -H "X-API-Key: $API_KEY" | jq '.'
echo ""

# 测试删除文章
if [ ! -z "$ARTICLE_ID" ]; then
    read -p "是否删除测试文章? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "🗑️ 测试删除文章..."
        curl -s -X DELETE "$API_BASE/articles/$ARTICLE_ID" \
          -H "X-API-Key: $API_KEY" | jq '.'
        echo ""
    fi
fi

echo "✅ API 测试完成！"
