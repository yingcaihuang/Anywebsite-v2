-- 完整的数据库初始化脚本
-- 适用于Docker容器首次启动时的数据库初始化

-- 设置字符集和基本配置
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO';
SET @OLD_TIME_ZONE=@@TIME_ZONE, TIME_ZONE='+08:00';

-- 使用数据库
USE static_hosting;

-- 清理现有表（如果存在）
DROP TABLE IF EXISTS articles;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS sessions;

-- ============================================================================
-- 创建文章表（使用UUID主键）
-- ============================================================================
CREATE TABLE articles (
    id VARCHAR(36) PRIMARY KEY COMMENT 'UUID主键',
    title VARCHAR(255) NOT NULL COMMENT '文章标题',
    content LONGTEXT COMMENT '文章内容（支持HTML）',
    slug VARCHAR(255) NOT NULL COMMENT 'URL友好的标识符',
    status VARCHAR(20) DEFAULT 'draft' COMMENT '文章状态：draft, published, archived',
    expires_at DATETIME(3) NULL COMMENT '过期时间',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    deleted_at DATETIME(3) NULL COMMENT '删除时间（软删除）',
    
    UNIQUE INDEX uni_articles_slug (slug),
    INDEX idx_articles_status (status),
    INDEX idx_articles_created_at (created_at),
    INDEX idx_articles_expires_at (expires_at),
    INDEX idx_articles_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章表';

-- ============================================================================
-- 创建用户表（用于后台管理）
-- ============================================================================
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    email VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱',
    password VARCHAR(255) NOT NULL COMMENT '加密后的密码',
    role VARCHAR(20) DEFAULT 'admin' COMMENT '用户角色：admin, editor',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否激活',
    last_login_at DATETIME(3) NULL COMMENT '最后登录时间',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    
    INDEX idx_users_username (username),
    INDEX idx_users_email (email),
    INDEX idx_users_role (role),
    INDEX idx_users_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================================================
-- 创建API密钥表
-- ============================================================================
CREATE TABLE api_keys (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL COMMENT 'API密钥名称',
    api_key VARCHAR(255) NOT NULL UNIQUE COMMENT 'API密钥',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否激活',
    permissions JSON COMMENT '权限配置',
    last_used_at DATETIME(3) NULL COMMENT '最后使用时间',
    expires_at DATETIME(3) NULL COMMENT '过期时间',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    updated_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
    
    UNIQUE INDEX uni_api_keys_key (api_key),
    INDEX idx_api_keys_is_active (is_active),
    INDEX idx_api_keys_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='API密钥表';

-- ============================================================================
-- 创建会话表（用于后台登录状态管理）
-- ============================================================================
CREATE TABLE sessions (
    id VARCHAR(128) PRIMARY KEY COMMENT '会话ID',
    user_id INT NOT NULL COMMENT '用户ID',
    ip_address VARCHAR(45) COMMENT 'IP地址',
    user_agent TEXT COMMENT '用户代理',
    payload LONGTEXT COMMENT '会话数据',
    last_activity DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '最后活动时间',
    created_at DATETIME(3) DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    
    INDEX idx_sessions_user_id (user_id),
    INDEX idx_sessions_last_activity (last_activity),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会话表';

-- ============================================================================
-- 插入初始数据
-- ============================================================================

-- 插入默认管理员用户
-- 密码为 "password" 的bcrypt哈希值
INSERT INTO users (username, email, password, role, is_active, created_at, updated_at) VALUES 
('admin', 'admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', TRUE, NOW(3), NOW(3));

-- 插入默认API密钥
INSERT INTO api_keys (name, api_key, is_active, permissions, created_at, updated_at) VALUES 
('Demo API Key', 'demo-api-key-12345', TRUE, '{"read": true, "write": true, "delete": true}', NOW(3), NOW(3)),
('N8N Integration Key', 'n8n-integration-key', TRUE, '{"read": true, "write": true, "delete": false}', NOW(3), NOW(3));

-- 插入示例文章（使用UUID函数生成ID）
INSERT INTO articles (id, title, content, slug, status, created_at, updated_at) VALUES 
(UUID(), '欢迎使用静态网页托管服务器', 
'<div style="max-width: 800px; margin: 0 auto; padding: 20px; font-family: Arial, sans-serif;">
    <h1 style="color: #2c3e50; text-align: center;">🎉 欢迎使用静态网页托管服务器</h1>
    
    <div style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 20px; border-radius: 10px; margin: 20px 0;">
        <h2 style="margin-top: 0;">✨ 主要功能</h2>
        <ul style="list-style-type: none; padding-left: 0;">
            <li style="margin: 10px 0;">🚀 <strong>API文章发布</strong> - 通过RESTful API轻松发布和管理文章</li>
            <li style="margin: 10px 0;">⏰ <strong>定时过期</strong> - 设置文章自动过期时间</li>
            <li style="margin: 10px 0;">🔐 <strong>后台管理</strong> - 安全的Session认证系统</li>
            <li style="margin: 10px 0;">🔑 <strong>API认证</strong> - X-API-Key密钥验证</li>
            <li style="margin: 10px 0;">🆔 <strong>UUID主键</strong> - 使用UUID作为文章唯一标识</li>
            <li style="margin: 10px 0;">🎨 <strong>HTML渲染</strong> - 支持富文本和HTML内容</li>
            <li style="margin: 10px 0;">🔒 <strong>SSL证书</strong> - 自动HTTPS证书管理</li>
        </ul>
    </div>

    <div style="background: #f8f9fa; padding: 20px; border-radius: 10px; border-left: 4px solid #28a745;">
        <h3 style="color: #28a745; margin-top: 0;">🚀 快速开始</h3>
        <p><strong>后台管理：</strong> 访问 <code style="background: #e9ecef; padding: 2px 6px; border-radius: 4px;">/admin</code> 进入管理后台</p>
        <p><strong>默认账号：</strong> <code style="background: #e9ecef; padding: 2px 6px; border-radius: 4px;">admin</code> / <code style="background: #e9ecef; padding: 2px 6px; border-radius: 4px;">password</code></p>
        <p><strong>API密钥：</strong> <code style="background: #e9ecef; padding: 2px 6px; border-radius: 4px;">demo-api-key-12345</code></p>
    </div>

    <div style="background: #fff3cd; padding: 20px; border-radius: 10px; border-left: 4px solid #ffc107; margin: 20px 0;">
        <h3 style="color: #856404; margin-top: 0;">📚 API 使用示例</h3>
        <h4>创建文章：</h4>
        <pre style="background: #f8f9fa; padding: 10px; border-radius: 4px; overflow-x: auto;"><code>curl -X POST http://localhost:8080/api/articles \\
  -H "Content-Type: application/json" \\
  -H "X-API-Key: demo-api-key-12345" \\
  -d "{
    \"title\": \"我的文章\",
    \"content\": \"<h1>这是文章内容</h1>\",
    \"slug\": \"my-article\",
    \"status\": \"published\"
  }"</code></pre>
        
        <h4>获取文章：</h4>
        <pre style="background: #f8f9fa; padding: 10px; border-radius: 4px; overflow-x: auto;"><code>curl -H "X-API-Key: demo-api-key-12345" \\
  http://localhost:8080/api/articles/{article_id}</code></pre>
    </div>

    <div style="text-align: center; margin-top: 30px; padding-top: 20px; border-top: 2px solid #e9ecef;">
        <p style="color: #6c757d;">
            <strong>系统启动时间：</strong> ' || NOW() || '<br>
            <strong>数据库版本：</strong> MySQL 8.0 with UUID Support<br>
            <strong>字符集：</strong> UTF8MB4
        </p>
    </div>
</div>', 
'welcome', 
'published', 
NOW(3), 
NOW(3)),

(UUID(), 'API 测试文档', 
'<div style="max-width: 800px; margin: 0 auto; padding: 20px; font-family: ''Segoe UI'', Tahoma, Geneva, Verdana, sans-serif;">
    <h1 style="color: #2c3e50;">📖 API 测试文档</h1>
    
    <h2 style="color: #3498db;">🔗 API 端点</h2>
    <table style="width: 100%; border-collapse: collapse; margin: 20px 0;">
        <thead>
            <tr style="background: #f8f9fa;">
                <th style="border: 1px solid #dee2e6; padding: 12px; text-align: left;">方法</th>
                <th style="border: 1px solid #dee2e6; padding: 12px; text-align: left;">端点</th>
                <th style="border: 1px solid #dee2e6; padding: 12px; text-align: left;">描述</th>
            </tr>
        </thead>
        <tbody>
            <tr>
                <td style="border: 1px solid #dee2e6; padding: 12px;"><code>POST</code></td>
                <td style="border: 1px solid #dee2e6; padding: 12px;"><code>/api/articles</code></td>
                <td style="border: 1px solid #dee2e6; padding: 12px;">创建新文章</td>
            </tr>
            <tr style="background: #f8f9fa;">
                <td style="border: 1px solid #dee2e6; padding: 12px;"><code>GET</code></td>
                <td style="border: 1px solid #dee2e6; padding: 12px;"><code>/api/articles/{id}</code></td>
                <td style="border: 1px solid #dee2e6; padding: 12px;">获取指定文章</td>
            </tr>
            <tr>
                <td style="border: 1px solid #dee2e6; padding: 12px;"><code>PUT</code></td>
                <td style="border: 1px solid #dee2e6; padding: 12px;"><code>/api/articles/{id}</code></td>
                <td style="border: 1px solid #dee2e6; padding: 12px;">更新文章</td>
            </tr>
            <tr style="background: #f8f9fa;">
                <td style="border: 1px solid #dee2e6; padding: 12px;"><code>DELETE</code></td>
                <td style="border: 1px solid #dee2e6; padding: 12px;"><code>/api/articles/{id}</code></td>
                <td style="border: 1px solid #dee2e6; padding: 12px;">删除文章</td>
            </tr>
        </tbody>
    </table>

    <h2 style="color: #e74c3c;">🔒 认证要求</h2>
    <div style="background: #ffe6e6; padding: 15px; border-radius: 8px; border-left: 4px solid #e74c3c;">
        <p><strong>所有API请求都需要在请求头中包含：</strong></p>
        <code style="background: #f8f9fa; padding: 8px; display: block; border-radius: 4px;">X-API-Key: demo-api-key-12345</code>
    </div>

    <h2 style="color: #27ae60;">✅ 测试状态</h2>
    <ul style="list-style-type: none; padding-left: 0;">
        <li style="margin: 8px 0;">✅ Session认证系统</li>
        <li style="margin: 8px 0;">✅ API Key认证</li>
        <li style="margin: 8px 0;">✅ UUID主键系统</li>
        <li style="margin: 8px 0;">✅ HTML内容渲染</li>
        <li style="margin: 8px 0;">✅ CRUD操作完整性</li>
        <li style="margin: 8px 0;">✅ 数据库连接稳定</li>
    </ul>
</div>', 
'api-docs', 
'published', 
NOW(3), 
NOW(3));

-- ============================================================================
-- 验证初始化结果
-- ============================================================================

-- 显示表结构信息
SELECT 
    'Database initialization completed successfully' as status,
    (SELECT COUNT(*) FROM articles) as articles_count,
    (SELECT COUNT(*) FROM users) as users_count,
    (SELECT COUNT(*) FROM api_keys) as api_keys_count;

-- 显示创建的表
SHOW TABLES;

-- 恢复设置
SET FOREIGN_KEY_CHECKS = 1;
SET SQL_MODE=@OLD_SQL_MODE;
SET TIME_ZONE=@OLD_TIME_ZONE;

-- 输出完成信息
SELECT 
    'Docker database initialization completed' as message,
    NOW(3) as completed_at,
    'Ready for application startup' as status;
