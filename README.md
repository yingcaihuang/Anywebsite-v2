# 静态网页托管服务器

一个使用 Golang 实现的功能完整的静态网页托管服务器，支持文章发布、过期管理、后台管理、API鉴权和自动SSL证书管理。

## 功能特点

- ✅ **API文章发布**: 通过RESTful API发布文章，设置过期时间
- ✅ **后台管理**: 基于Web的管理界面，支持文章增删改查
- ✅ **API鉴权**: 支持API密钥认证，确保接口安全
- ✅ **n8n兼容**: API响应格式符合n8n集成规范
- ✅ **静态文件生成**: 自动将文章转换为静态HTML页面
- ✅ **过期管理**: 定时清理过期文章和静态文件
- ⏳ **ACME证书**: 自动申请和续期SSL证书
- ✅ **Docker部署**: 支持Docker Compose一键部署

## 快速开始

### 使用 Docker Compose (推荐)

1. 克隆项目
```bash
git clone <repository-url>
cd static-hosting-server
```

2. 启动服务
```bash
docker-compose up -d
```

3. 访问应用
- 管理后台: http://localhost:8080/admin (admin/admin123)
- API文档: http://localhost:8080/api
- 示例文章: http://localhost:8080/p/welcome

### 本地开发

1. 安装依赖
```bash
go mod tidy
```

2. 配置数据库
```bash
# 启动MySQL数据库
docker run -d --name mysql \
  -e MYSQL_ROOT_PASSWORD=rootpassword \
  -e MYSQL_DATABASE=static_hosting \
  -e MYSQL_USER=app \
  -e MYSQL_PASSWORD=password \
  -p 3306:3306 \
  mysql:8.0
```

3. 运行应用
```bash
go run cmd/server/main.go
```

## API 使用说明

### 认证

所有API请求需要在Header中包含API密钥：
```
X-API-Key: demo-api-key-12345
```

### 创建文章

```bash
curl -X POST http://localhost:8080/api/articles \
  -H "Content-Type: application/json" \
  -H "X-API-Key: demo-api-key-12345" \
  -d '{
    "title": "我的文章",
    "content": "<h1>Hello World</h1><p>这是文章内容</p>",
    "slug": "my-article",
    "status": "published",
    "expires_at": "2024-12-31T23:59:59Z"
  }'
```

### 响应格式 (n8n兼容)

```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "我的文章",
    "slug": "my-article",
    "status": "published",
    "created_at": "2024-01-01T12:00:00Z"
  },
  "url": "http://localhost:8080/p/my-article"
}
```

### 更新文章

```bash
curl -X PUT http://localhost:8080/api/articles/1 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: demo-api-key-12345" \
  -d '{
    "title": "更新的标题",
    "content": "<h1>更新的内容</h1>",
    "status": "published"
  }'
```

### 删除文章

```bash
curl -X DELETE http://localhost:8080/api/articles/1 \
  -H "X-API-Key: demo-api-key-12345"
```

### 获取文章列表

```bash
curl -X GET "http://localhost:8080/api/articles?page=1&limit=10&status=published" \
  -H "X-API-Key: demo-api-key-12345"
```

## 配置说明

配置文件位于 `configs/config.yml`：

```yaml
server:
  port: "8080"
  mode: "debug"  # debug, release
  domain: "localhost:8080"

database:
  host: "localhost"
  port: 3306
  user: "app"
  password: "password"
  dbname: "static_hosting"
  charset: "utf8mb4"

security:
  api_keys:
    - "demo-api-key-12345"
    - "n8n-integration-key"

storage:
  static_path: "./static"
  uploads_path: "./uploads"
  certs_path: "./certs"
```

## 环境变量

支持通过环境变量覆盖配置，环境变量前缀为 `SHS_`：

- `SHS_DATABASE_HOST`: 数据库主机
- `SHS_DATABASE_USER`: 数据库用户名
- `SHS_DATABASE_PASSWORD`: 数据库密码
- `SHS_SERVER_DOMAIN`: 服务器域名

## 目录结构

```
.
├── cmd/
│   └── server/           # 主程序入口
├── internal/
│   ├── api/             # API路由和处理器
│   ├── auth/            # 认证中间件
│   ├── config/          # 配置管理
│   ├── database/        # 数据库连接
│   ├── models/          # 数据模型
│   ├── scheduler/       # 定时任务
│   ├── services/        # 业务逻辑
│   └── web/             # Web管理界面
├── templates/           # HTML模板
├── static/              # 静态文件目录
├── configs/             # 配置文件
├── scripts/             # 脚本文件
├── docker-compose.yml   # Docker编排
└── Dockerfile          # Docker镜像
```

## 后台管理

访问 http://localhost:8080/admin 进入管理后台

默认账号：
- 用户名: admin
- 密码: admin123

功能包括：
- 文章列表和搜索
- 创建和编辑文章
- 文章状态管理
- 过期时间设置

## n8n 集成

本服务器的API完全兼容n8n工作流，可以直接在n8n中使用：

1. 在n8n中添加HTTP请求节点
2. 设置API密钥认证
3. 使用提供的API端点进行集成

API响应格式符合n8n标准，包含：
- `success`: 操作是否成功
- `data`: 返回数据
- `error`: 错误信息（如有）
- `url`: 生成的文章URL（如适用）

## 开发计划

- [ ] ACME证书自动管理
- [ ] 更多认证方式支持
- [ ] 文章分类和标签
- [ ] 文件上传功能
- [ ] 访问统计
- [ ] 更多主题模板

## 贡献

欢迎提交Issue和Pull Request来改进这个项目。

## 许可证

MIT License
