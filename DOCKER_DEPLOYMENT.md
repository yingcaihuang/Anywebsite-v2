# Docker 部署指南

本指南提供了完整的 Docker 部署流程，包括数据库初始化、数据迁移和系统验证。

## 🚀 快速开始

### 1. 使用自动化脚本启动（推荐）

**Windows:**
```bash
./docker-start.bat
```

**Linux/macOS:**
```bash
chmod +x docker-start.sh
./docker-start.sh
```

### 2. 手动启动

```bash
# 清理旧容器
docker-compose down --volumes --remove-orphans

# 构建并启动
docker-compose up --build -d

# 查看日志
docker-compose logs -f
```

## 📊 数据库初始化

### 自动初始化

Docker Compose 会自动执行以下初始化操作：

1. **创建数据库和用户**
   - 数据库：`static_hosting`
   - 用户：`app` / `password`
   - 字符集：`utf8mb4`

2. **创建表结构**
   - `articles` - 文章表（UUID主键）
   - `users` - 用户表
   - `api_keys` - API密钥表
   - `sessions` - 会话表

3. **插入初始数据**
   - 默认管理员：`admin` / `password`
   - API密钥：`demo-api-key-12345`, `n8n-integration-key`
   - 示例文章：欢迎页面和API文档

### 初始化脚本位置

- 主初始化脚本：`./scripts/docker-init.sql`
- MySQL配置：`./docker/mysql/conf.d/docker.cnf`

## 🔍 验证部署

### 自动验证

```bash
# Linux/macOS
chmod +x docker/verify.sh
./docker/verify.sh

# 或手动执行验证
docker-compose exec mysql mysql -u app -ppassword static_hosting -e "
SELECT 
    'Verification' as type,
    (SELECT COUNT(*) FROM articles) as articles,
    (SELECT COUNT(*) FROM users) as users,
    (SELECT COUNT(*) FROM api_keys) as api_keys;
"
```

### 手动验证检查点

1. **服务状态检查**
   ```bash
   docker-compose ps
   ```

2. **Web服务测试**
   ```bash
   curl http://localhost:8080/
   ```

3. **API认证测试**
   ```bash
   curl -H "X-API-Key: demo-api-key-12345" http://localhost:8080/
   ```

4. **管理后台测试**
   - 访问：http://localhost:8080/admin
   - 账号：admin / password

## 🗃️ 数据库管理

### 备份数据库

```bash
# 使用备份脚本
chmod +x docker/backup.sh
./docker/backup.sh backup

# 手动备份
docker-compose exec mysql mysqldump -u app -ppassword static_hosting > backup.sql
```

### 恢复数据库

```bash
# 使用备份脚本
./docker/backup.sh restore ./backups/backup_20250701_120000.sql

# 手动恢复
docker-compose exec -T mysql mysql -u app -ppassword static_hosting < backup.sql
```

### 查看备份

```bash
./docker/backup.sh list
```

## 🔧 配置说明

### 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `SHS_DATABASE_HOST` | `mysql` | 数据库主机 |
| `SHS_DATABASE_USER` | `app` | 数据库用户 |
| `SHS_DATABASE_PASSWORD` | `password` | 数据库密码 |
| `SHS_DATABASE_DBNAME` | `static_hosting` | 数据库名称 |
| `SHS_SERVER_DOMAIN` | `localhost:8080` | 服务器域名 |

### 端口映射

| 服务 | 容器端口 | 主机端口 | 说明 |
|------|----------|----------|------|
| Web | 8080 | 8080 | HTTP服务 |
| Web | 8443 | 8443 | HTTPS服务 |
| MySQL | 3306 | 3306 | 数据库服务 |

### 存储卷

| 本地路径 | 容器路径 | 说明 |
|----------|----------|------|
| `./static` | `/root/static` | 静态文件 |
| `./certs` | `/root/certs` | SSL证书 |
| `./uploads` | `/root/uploads` | 上传文件 |
| `./configs` | `/root/configs` | 配置文件 |
| `mysql_data` | `/var/lib/mysql` | MySQL数据 |

## 🚨 故障排除

### 常见问题

1. **MySQL启动失败**
   ```bash
   # 查看MySQL日志
   docker-compose logs mysql
   
   # 重新初始化数据库
   docker-compose down --volumes
   docker-compose up -d mysql
   ```

2. **Web服务无法访问**
   ```bash
   # 查看Web服务日志
   docker-compose logs web
   
   # 检查端口占用
   netstat -tulpn | grep :8080
   ```

3. **数据库连接失败**
   ```bash
   # 测试数据库连接
   docker-compose exec mysql mysql -u app -ppassword static_hosting -e "SELECT 1;"
   ```

4. **UUID迁移问题**
   ```bash
   # 重新执行初始化
   docker-compose exec mysql mysql -u app -ppassword static_hosting < scripts/docker-init.sql
   ```

### 日志查看

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f web
docker-compose logs -f mysql

# 查看最近日志
docker-compose logs --tail=100 web
```

### 性能监控

```bash
# 查看容器资源使用
docker stats

# 查看数据库性能
docker-compose exec mysql mysql -u app -ppassword -e "SHOW PROCESSLIST;"
```

## 📈 性能优化

### MySQL 优化配置

配置文件：`./docker/mysql/conf.d/docker.cnf`

主要优化项：
- InnoDB缓冲池：256MB
- 连接数限制：100
- 查询缓存：32MB
- 日志优化：减少磁盘I/O

### 容器资源限制

在 `docker-compose.yml` 中添加资源限制：

```yaml
services:
  mysql:
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '1.0'
        reservations:
          memory: 256M
          cpus: '0.5'
```

## 🔄 更新和维护

### 更新应用

```bash
# 拉取最新代码
git pull

# 重新构建并启动
docker-compose up --build -d

# 验证更新
./docker/verify.sh
```

### 清理和重置

```bash
# 完全清理（删除所有数据）
docker-compose down --volumes --rmi all
docker system prune -a

# 重新部署
./docker-start.sh
```

## 📞 支持

如果遇到问题：

1. 查看日志：`docker-compose logs`
2. 验证系统：`./docker/verify.sh`
3. 检查网络：`docker network ls`
4. 重启服务：`docker-compose restart`

## 🎯 生产环境建议

1. **安全性**
   - 修改默认密码
   - 使用环境变量文件
   - 启用SSL证书
   - 配置防火墙

2. **备份策略**
   - 定期自动备份
   - 异地备份存储
   - 备份恢复测试

3. **监控**
   - 容器健康检查
   - 日志聚合
   - 性能监控
   - 告警设置
