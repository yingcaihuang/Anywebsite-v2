-- 数据库迁移脚本：将Article表的ID从INT改为VARCHAR(36) UUID
-- 注意：这个脚本会清空现有数据，因为改变主键类型比较复杂
-- 在生产环境中使用前请备份数据！

USE static_hosting;

-- 1. 备份现有数据（可选）
-- CREATE TABLE articles_backup AS SELECT * FROM articles;

-- 2. 删除外键约束（如果有的话）
-- 由于我们的设计中Article是独立表，通常没有外键指向它

-- 3. 重建articles表
DROP TABLE IF EXISTS articles;

CREATE TABLE articles (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content LONGTEXT,
    slug VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'draft',
    expires_at DATETIME(3),
    created_at DATETIME(3),
    updated_at DATETIME(3),
    deleted_at DATETIME(3),
    INDEX idx_articles_deleted_at (deleted_at),
    UNIQUE INDEX uni_articles_slug (slug)
);

-- 4. 如果需要恢复数据，可以使用以下语句（需要为每行生成UUID）
-- INSERT INTO articles (id, title, content, slug, status, expires_at, created_at, updated_at, deleted_at)
-- SELECT UUID(), title, content, slug, status, expires_at, created_at, updated_at, deleted_at
-- FROM articles_backup;

-- 完成迁移
SELECT 'Article table migration completed' as status;
