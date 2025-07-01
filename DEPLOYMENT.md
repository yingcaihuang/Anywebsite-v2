# 🚀 静态网页托管服务器 - 部署指南

## 快速启动

### 方法一：使用 Docker Compose (推荐)

1. **确保已安装 Docker**
   - Windows: 安装 Docker Desktop
   - Linux: 安装 Docker 和 Docker Compose

2. **启动服务**
   ```bash
   # Linux/Mac
   chmod +x start.sh
   ./start.sh
   
   # Windows
   start.bat
   ```

3. **访问应用**
   - 管理后台: http://localhost:8080/admin
   - 默认账号: `admin` / `admin123`
   - 示例文章: http://localhost:8080/p/welcome

### 方法二：本地开发运行

1. **安装依赖**
   ```bash
   go mod tidy
   ```

2. **启动 MySQL 数据库**
   ```bash
   docker run -d --name mysql \
     -e MYSQL_ROOT_PASSWORD=rootpassword \
     -e MYSQL_DATABASE=static_hosting \
     -e MYSQL_USER=app \
     -e MYSQL_PASSWORD=password \
     -p 3306:3306 \
     mysql:8.0
   ```

3. **运行应用**
   ```bash
   go run cmd/server/main.go
   ```

## 🔧 配置说明

### 环境变量

可以通过环境变量覆盖配置，前缀为 `SHS_`：

```bash
export SHS_DATABASE_HOST=localhost
export SHS_DATABASE_USER=app
export SHS_DATABASE_PASSWORD=password
export SHS_DATABASE_DBNAME=static_hosting
export SHS_SERVER_DOMAIN=yourdomain.com
```

### 配置文件

编辑 `configs/config.yml`：

```yaml
server:
  port: "8080"
  mode: "release"  # 生产环境使用 release
  domain: "yourdomain.com"

database:
  host: "mysql"
  port: 3306
  user: "app"
  password: "your_secure_password"
  dbname: "static_hosting"

security:
  api_keys:
    - "your-secure-api-key"
    - "n8n-integration-key"
```

## 📚 API 使用指南

### 认证

所有API请求需要在Header中包含API密钥：
```
X-API-Key: your-api-key
```

### 创建文章

```bash
curl -X POST http://localhost:8080/api/articles \
  -H "Content-Type: application/json" \
  -H "X-API-Key: demo-api-key-12345" \
  -d '{
    "title": "我的文章",
    "content": "<h1>文章标题</h1><p>文章内容</p>",
    "slug": "my-article",
    "status": "published",
    "expires_at": "2024-12-31T23:59:59Z"
  }'
```

### n8n 集成示例

在 n8n 中创建 HTTP 请求节点：

1. **URL**: `http://your-server:8080/api/articles`
2. **Method**: POST
3. **Headers**: 
   - `Content-Type`: application/json
   - `X-API-Key`: your-api-key
4. **Body**: JSON格式的文章数据

响应格式：
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "文章标题",
    "slug": "article-slug",
    "status": "published"
  },
  "url": "http://your-server:8080/p/article-slug"
}
```

## 🛠️ 生产环境部署

### 1. 域名和SSL

修改 `docker-compose.yml` 添加 SSL 支持：

```yaml
services:
  web:
    ports:
      - "80:8080"
      - "443:8443"
    environment:
      - SHS_SERVER_DOMAIN=yourdomain.com
      - SHS_ACME_EMAIL=admin@yourdomain.com
```

### 2. 数据备份

```bash
# 备份数据库
docker exec mysql mysqldump -u app -p static_hosting > backup.sql

# 备份静态文件
tar -czf static_backup.tar.gz static/
```

### 3. 反向代理 (Nginx 示例)

```nginx
server {
    listen 80;
    server_name yourdomain.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 🧪 测试

运行API测试脚本：
```bash
chmod +x test-api.sh
./test-api.sh
```

## 🔧 故障排除

### 常见问题

1. **数据库连接失败**
   ```bash
   # 检查数据库状态
   docker-compose logs mysql
   
   # 重新启动数据库
   docker-compose restart mysql
   ```

2. **静态文件生成失败**
   ```bash
   # 检查目录权限
   chmod -R 755 static/
   
   # 检查应用日志
   docker-compose logs web
   ```

3. **API密钥无效**
   - 确保在配置文件中正确设置API密钥
   - 检查请求Header中的密钥格式

### 日志查看

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f web
docker-compose logs -f mysql
```

## 📊 监控

可以集成以下监控工具：
- Prometheus + Grafana
- ELK Stack
- Application Performance Monitoring (APM)

## 🔄 更新和维护

### 更新应用

```bash
# 停止服务
docker-compose down

# 拉取新代码
git pull

# 重新构建和启动
docker-compose up -d --build
```

### 数据库迁移

应用启动时会自动执行数据库迁移，无需手动操作。

---

## 🎉 完成！

您的静态网页托管服务器现在已经准备就绪。这是一个功能完整的解决方案，包含：

- ✅ RESTful API (n8n兼容)
- ✅ 现代化管理后台
- ✅ 自动过期管理
- ✅ 静态文件生成
- ✅ Docker 容器化部署
- ✅ API 认证和安全
- ✅ 完整的文档和测试脚本

**下一步计划：**
- [ ] ACME证书自动管理 (高级功能)
- [ ] 更多主题模板
- [ ] 文件上传功能
- [ ] 访问统计和分析

如有问题，请查看日志或提交Issue。祝使用愉快！ 🚀
