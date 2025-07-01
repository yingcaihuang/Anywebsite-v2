package web

import (
	"net/http"
	"static-hosting-server/internal/auth"
	"static-hosting-server/internal/config"
	"static-hosting-server/internal/models"
	"static-hosting-server/internal/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebHandler struct {
	db             *gorm.DB
	cfg            *config.Config
	authService    *auth.AuthService
	articleService *services.ArticleService
}

func NewWebHandler(db *gorm.DB, cfg *config.Config) *WebHandler {
	authService := auth.NewAuthService(db, cfg)
	articleService := services.NewArticleService(db, cfg)

	return &WebHandler{
		db:             db,
		cfg:            cfg,
		authService:    authService,
		articleService: articleService,
	}
}

func SetupRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	handler := NewWebHandler(db, cfg)

	// 管理后台路由
	admin := router.Group("/admin")
	{
		// 登录页面（暂时简化）
		admin.GET("/login", handler.LoginPage)
		admin.POST("/login", handler.Login)
		admin.POST("/logout", handler.Logout) // 添加退出登录

		// 需要认证的路由
		authenticated := admin.Group("")
		authenticated.Use(handler.authService.AdminAuthMiddleware())
		{
			authenticated.GET("", handler.Dashboard)
			authenticated.GET("/dashboard", handler.Dashboard)

			// 文章管理
			authenticated.GET("/articles", handler.ArticlesList)
			authenticated.GET("/articles/new", handler.NewArticlePage)
			authenticated.POST("/articles", handler.CreateArticleWeb)
			authenticated.GET("/articles/:id/edit", handler.EditArticlePage)
			authenticated.POST("/articles/:id", handler.UpdateArticleWeb)
			authenticated.POST("/articles/:id/delete", handler.DeleteArticleWeb)
		}
	}
}

// 登录页面
func (h *WebHandler) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "管理员登录",
	})
}

// 登录处理
func (h *WebHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 验证管理员凭据
	if h.authService.ValidateAdminCredentials(username, password) {
		// 生成会话token并设置cookie
		sessionToken := h.authService.GenerateSessionToken()
		c.SetCookie("admin_session", sessionToken, 3600*24*7, "/", "", false, true) // 7天有效期
		c.Redirect(http.StatusFound, "/admin/dashboard")
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "管理员登录",
		"error": "用户名或密码错误",
	})
}

// 仪表板
func (h *WebHandler) Dashboard(c *gin.Context) {
	// 获取统计信息
	stats, err := h.getStatistics()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title": "管理后台",
		"stats": stats,
	})
}

// 文章列表
func (h *WebHandler) ArticlesList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit := 20
	status := c.Query("status")

	articles, total, err := h.articleService.ListArticles(page, limit, status)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "articles_list.html", gin.H{
		"title":    "文章管理",
		"articles": articles,
		"total":    total,
		"page":     page,
		"limit":    limit,
		"status":   status,
	})
}

// 新建文章页面
func (h *WebHandler) NewArticlePage(c *gin.Context) {
	c.HTML(http.StatusOK, "article_form.html", gin.H{
		"title":  "新建文章",
		"action": "/admin/articles",
		"method": "POST",
	})
}

// 创建文章（Web表单）
func (h *WebHandler) CreateArticleWeb(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	slug := c.PostForm("slug")
	status := c.PostForm("status")
	expiresAtStr := c.PostForm("expires_at")

	var expiresAt *time.Time
	if expiresAtStr != "" {
		if parsed, err := time.Parse("2006-01-02T15:04", expiresAtStr); err == nil {
			expiresAt = &parsed
		}
	}

	_, err := h.articleService.CreateArticle(title, content, slug, status, expiresAt)
	if err != nil {
		c.HTML(http.StatusBadRequest, "article_form.html", gin.H{
			"title":  "新建文章",
			"action": "/admin/articles",
			"method": "POST",
			"error":  err.Error(),
			"form_data": gin.H{
				"title":      title,
				"content":    content,
				"slug":       slug,
				"status":     status,
				"expires_at": expiresAtStr,
			},
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/articles")
}

// 编辑文章页面
func (h *WebHandler) EditArticlePage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Invalid article ID",
		})
		return
	}

	article, err := h.articleService.GetArticleByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Article not found",
		})
		return
	}

	c.HTML(http.StatusOK, "article_form.html", gin.H{
		"title":   "编辑文章",
		"action":  "/admin/articles/" + id,
		"method":  "POST",
		"article": article,
	})
}

// 更新文章（Web表单）
func (h *WebHandler) UpdateArticleWeb(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Invalid article ID",
		})
		return
	}

	title := c.PostForm("title")
	content := c.PostForm("content")
	status := c.PostForm("status")
	expiresAtStr := c.PostForm("expires_at")

	var expiresAt *time.Time
	if expiresAtStr != "" {
		if parsed, err := time.Parse("2006-01-02T15:04", expiresAtStr); err == nil {
			expiresAt = &parsed
		}
	}

	_, err := h.articleService.UpdateArticle(id, title, content, status, expiresAt)
	if err != nil {
		article, _ := h.articleService.GetArticleByID(id)
		c.HTML(http.StatusBadRequest, "article_form.html", gin.H{
			"title":   "编辑文章",
			"action":  "/admin/articles/" + id,
			"method":  "POST",
			"article": article,
			"error":   err.Error(),
		})
		return
	}

	c.Redirect(http.StatusFound, "/admin/articles")
}

// 删除文章（Web表单）
func (h *WebHandler) DeleteArticleWeb(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	if err := h.articleService.DeleteArticle(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "/admin/articles")
}

// 退出登录
func (h *WebHandler) Logout(c *gin.Context) {
	// 清除会话cookie
	c.SetCookie("admin_session", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/admin/login")
}

// 获取统计信息
func (h *WebHandler) getStatistics() (map[string]interface{}, error) {
	var totalArticles int64
	var publishedArticles int64
	var draftArticles int64

	// 修正：使用正确的模型引用
	if err := h.db.Model(&models.Article{}).Count(&totalArticles).Error; err != nil {
		return nil, err
	}

	if err := h.db.Model(&models.Article{}).Where("status = ?", "published").Count(&publishedArticles).Error; err != nil {
		return nil, err
	}

	if err := h.db.Model(&models.Article{}).Where("status = ?", "draft").Count(&draftArticles).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_articles":     totalArticles,
		"published_articles": publishedArticles,
		"draft_articles":     draftArticles,
	}, nil
}
