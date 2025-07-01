package services

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"static-hosting-server/internal/config"
	"static-hosting-server/internal/models"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ArticleService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewArticleService(db *gorm.DB, cfg *config.Config) *ArticleService {
	return &ArticleService{
		db:  db,
		cfg: cfg,
	}
}

// 创建文章
func (s *ArticleService) CreateArticle(title, content, slug, status string, expiresAt *time.Time) (*models.Article, error) {
	// 如果没有提供slug，从标题生成
	if slug == "" {
		slug = s.generateSlugFromTitle(title)
	}

	// 检查slug是否已存在
	var existingArticle models.Article
	if err := s.db.Where("slug = ?", slug).First(&existingArticle).Error; err == nil {
		return nil, fmt.Errorf("article with slug '%s' already exists", slug)
	}

	// 设置默认状态
	if status == "" {
		status = "draft"
	}

	article := &models.Article{
		Title:     title,
		Content:   content,
		Slug:      slug,
		Status:    status,
		ExpiresAt: expiresAt,
	}

	if err := s.db.Create(article).Error; err != nil {
		return nil, err
	}

	// 如果状态为已发布，生成静态文件
	if status == "published" {
		if err := s.generateStaticFiles(article); err != nil {
			// 记录错误但不回滚创建操作
			fmt.Printf("Failed to generate static files for article %d: %v\n", article.ID, err)
		}
	}

	return article, nil
}

// 根据ID获取文章
func (s *ArticleService) GetArticleByID(id string) (*models.Article, error) {
	var article models.Article
	if err := s.db.Where("id = ?", id).First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// 根据slug获取已发布的文章
func (s *ArticleService) GetPublishedArticleBySlug(slug string) (*models.Article, error) {
	var article models.Article
	if err := s.db.Where("slug = ? AND status = ?", slug, "published").First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// 更新文章
func (s *ArticleService) UpdateArticle(id string, title, content, status string, expiresAt *time.Time) (*models.Article, error) {
	var article models.Article
	if err := s.db.Where("id = ?", id).First(&article).Error; err != nil {
		return nil, err
	}

	oldStatus := article.Status

	// 更新字段
	updates := make(map[string]interface{})
	if title != "" {
		updates["title"] = title
	}
	if content != "" {
		updates["content"] = content
	}
	if status != "" {
		updates["status"] = status
	}
	if expiresAt != nil {
		updates["expires_at"] = expiresAt
	}

	if err := s.db.Model(&article).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 重新获取更新后的文章
	if err := s.db.Where("id = ?", id).First(&article).Error; err != nil {
		return nil, err
	}

	// 处理状态变化
	if oldStatus != article.Status {
		if article.Status == "published" {
			// 生成静态文件
			if err := s.generateStaticFiles(&article); err != nil {
				fmt.Printf("Failed to generate static files for article %d: %v\n", article.ID, err)
			}
		} else if oldStatus == "published" {
			// 删除静态文件
			if err := s.removeStaticFiles(article.Slug); err != nil {
				fmt.Printf("Failed to remove static files for article %d: %v\n", article.ID, err)
			}
		}
	} else if article.Status == "published" {
		// 更新静态文件
		if err := s.generateStaticFiles(&article); err != nil {
			fmt.Printf("Failed to update static files for article %d: %v\n", article.ID, err)
		}
	}

	return &article, nil
}

// 删除文章
func (s *ArticleService) DeleteArticle(id string) error {
	var article models.Article
	if err := s.db.Where("id = ?", id).First(&article).Error; err != nil {
		return err
	}

	// 删除静态文件
	if article.Status == "published" {
		if err := s.removeStaticFiles(article.Slug); err != nil {
			fmt.Printf("Failed to remove static files for article %s: %v\n", article.ID, err)
		}
	}

	return s.db.Delete(&article).Error
}

// 获取文章列表
func (s *ArticleService) ListArticles(page, limit int, status string) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	query := s.db.Model(&models.Article{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// 生成静态HTML文件
func (s *ArticleService) generateStaticFiles(article *models.Article) error {
	// 创建文章目录
	articleDir := filepath.Join(s.cfg.Storage.StaticPath, "articles", article.Slug)
	if err := os.MkdirAll(articleDir, 0755); err != nil {
		return fmt.Errorf("failed to create article directory: %w", err)
	}

	// 读取模板并添加自定义函数
	tmplPath := "templates/article.html"
	tmpl, err := template.New("article.html").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}).ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// 创建HTML文件
	htmlPath := filepath.Join(articleDir, "index.html")
	file, err := os.Create(htmlPath)
	if err != nil {
		return fmt.Errorf("failed to create HTML file: %w", err)
	}
	defer file.Close()

	// 渲染模板
	data := struct {
		Article *models.Article
		Domain  string
	}{
		Article: article,
		Domain:  s.cfg.Server.Domain,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// 删除静态文件
func (s *ArticleService) removeStaticFiles(slug string) error {
	articleDir := filepath.Join(s.cfg.Storage.StaticPath, "articles", slug)
	return os.RemoveAll(articleDir)
}

// 从标题生成slug
func (s *ArticleService) generateSlugFromTitle(title string) string {
	// 简化的slug生成逻辑
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// 移除特殊字符（这里简化处理）
	var result strings.Builder
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}

	slug = result.String()

	// 添加时间戳以确保唯一性
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%d", slug, timestamp)
}

// 清理过期文章
func (s *ArticleService) CleanupExpiredArticles() error {
	var expiredArticles []models.Article
	if err := s.db.Where("expires_at IS NOT NULL AND expires_at < ? AND status = ?",
		time.Now(), "published").Find(&expiredArticles).Error; err != nil {
		return err
	}

	for _, article := range expiredArticles {
		// 删除静态文件
		if err := s.removeStaticFiles(article.Slug); err != nil {
			fmt.Printf("Failed to remove static files for expired article %d: %v\n", article.ID, err)
		}

		// 更新状态为过期
		if err := s.db.Model(&article).Update("status", "expired").Error; err != nil {
			fmt.Printf("Failed to update status for expired article %d: %v\n", article.ID, err)
		}
	}

	return nil
}
