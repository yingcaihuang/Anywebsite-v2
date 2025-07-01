-- 数据库初始化脚本

-- 设置字符集
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建默认管理员用户（如果需要）
-- INSERT INTO users (username, email, password, role, created_at, updated_at) 
-- VALUES ('admin', 'admin@example.com', '$2a$10$...', 'admin', NOW(), NOW());

-- 创建默认API密钥
-- INSERT INTO api_keys (name, `key`, is_active, permissions, created_at, updated_at)
-- VALUES ('Default API Key', 'demo-api-key-12345', true, '{"read": true, "write": true}', NOW(), NOW());

-- 插入示例文章
INSERT INTO articles (title, content, slug, status, created_at, updated_at) VALUES 
('欢迎使用静态网页托管服务器', 
'<h2>欢迎！</h2>
<p>这是一个功能强大的静态网页托管服务器，支持：</p>
<ul>
<li>通过 API 发布文章</li>
<li>设置文章过期时间</li>
<li>后台管理界面</li>
<li>API 鉴权验证</li>
<li>自动 SSL 证书管理</li>
</ul>
<h3>快速开始</h3>
<p>访问 <code>/admin</code> 进入管理后台，默认账号：admin / admin123</p>
<p>使用 API 密钥：<code>demo-api-key-12345</code> 来调用 API 接口。</p>', 
'welcome', 
'published', 
NOW(), 
NOW());

SET FOREIGN_KEY_CHECKS = 1;
