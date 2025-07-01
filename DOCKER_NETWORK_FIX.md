# Docker 网络问题解决方案

## 问题描述
Docker 构建时出现网络连接错误，无法连接到 Docker Hub 或下载镜像。

## 解决方案

### 方案一：使用镜像加速器 (推荐)

1. **配置 Docker Desktop 镜像加速**
   - 打开 Docker Desktop
   - 进入 Settings (设置)
   - 选择 Docker Engine
   - 在配置文件中添加镜像加速器：

```json
{
  "registry-mirrors": [
    "https://dockerproxy.com",
    "https://hub-mirror.c.163.com",
    "https://mirror.baidubce.com"
  ]
}
```

2. **应用配置并重启 Docker**

### 方案二：修改网络设置

1. **切换 DNS 服务器**
   - 将 DNS 设置为：`8.8.8.8` 或 `114.114.114.114`
   - Windows: 网络设置 → 更改适配器选项 → 属性 → IPv4

2. **检查防火墙设置**
   - 确保 Docker Desktop 被允许通过防火墙
   - 暂时关闭杀毒软件的网络防护

### 方案三：使用代理 (如果有)

如果您使用代理，需要配置 Docker 代理：

在 Docker Desktop 设置中添加：
```json
{
  "proxies": {
    "default": {
      "httpProxy": "http://proxy.example.com:8080",
      "httpsProxy": "http://proxy.example.com:8080"
    }
  }
}
```

### 方案四：离线构建 (备用方案)

如果网络问题持续，可以：

1. **使用本地 Go 构建**
```bash
# 直接在本地编译
go build -o server.exe ./cmd/server

# 然后使用更简单的 Dockerfile
```

2. **使用预构建镜像**
我们已经修改了 Dockerfile 使用国内镜像源。

## 快速修复步骤

1. **重新配置 Docker Desktop**
   - 停止 Docker Desktop
   - 应用镜像加速器配置
   - 重启 Docker Desktop

2. **清理 Docker 缓存**
```bash
docker system prune -a
```

3. **重新构建**
```bash
docker-compose build --no-cache
```

## 验证修复

运行以下命令验证网络连接：
```bash
# 测试 Docker Hub 连接
docker pull hello-world

# 测试镜像加速器
docker pull alpine:3.18
```

## 常见错误代码

- `connectex: A connection attempt failed` - 网络连接超时
- `failed to fetch oauth token` - 认证服务器连接失败
- `dial tcp: i/o timeout` - DNS 解析超时

## 预防措施

1. 使用固定版本的镜像标签而不是 `latest`
2. 配置镜像加速器
3. 定期清理 Docker 缓存
4. 使用多阶段构建减少网络依赖

---

**提示**: 修改后的 Dockerfile 已经包含了镜像加速优化，应该能解决大部分网络问题。
