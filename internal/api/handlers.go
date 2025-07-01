package api

import (
	"net/http"
	"static-hosting-server/internal/auth"
	"static-hosting-server/internal/config"
	"static-hosting-server/internal/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db             *gorm.DB
	cfg            *config.Config
	authService    *auth.AuthService
	articleService *services.ArticleService
}

func NewHandler(db *gorm.DB, cfg *config.Config) *Handler {
	authService := auth.NewAuthService(db, cfg)
	articleService := services.NewArticleService(db, cfg)

	return &Handler{
		db:             db,
		cfg:            cfg,
		authService:    authService,
		articleService: articleService,
	}
}

func SetupRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	handler := NewHandler(db, cfg)

	// API 路由组
	api := router.Group("/api")
	api.Use(handler.authService.APIKeyAuthMiddleware())
	{
		// 文章相关API
		articles := api.Group("/articles")
		{
			articles.POST("", handler.CreateArticle)
			articles.GET("/:id", handler.GetArticle)
			articles.PUT("/:id", handler.UpdateArticle)
			articles.DELETE("/:id", handler.DeleteArticle)
			articles.GET("", handler.ListArticles)
		}

		// API密钥管理
		apiKeys := api.Group("/keys")
		{
			apiKeys.POST("", handler.CreateAPIKey)
			apiKeys.GET("", handler.ListAPIKeys)
			apiKeys.DELETE("/:id", handler.DeleteAPIKey)
		}
	}

	// 公开的文章访问API
	router.GET("/p/:slug", handler.GetPublishedArticle)
}

// n8n 兼容的响应格式
type N8nResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	URL     string      `json:"url,omitempty"`
}

// 创建文章
func (h *Handler) CreateArticle(c *gin.Context) {
	var req struct {
		Title     string     `json:"title" binding:"required"`
		Content   string     `json:"content" binding:"required"`
		Slug      string     `json:"slug"`
		Status    string     `json:"status"`
		ExpiresAt *time.Time `json:"expires_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, N8nResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	article, err := h.articleService.CreateArticle(req.Title, req.Content, req.Slug, req.Status, req.ExpiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, N8nResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// 生成发布URL
	publishURL := ""
	if article.Status == "published" {
		publishURL = h.cfg.Server.Domain + "/p/" + article.Slug
	}

	c.JSON(http.StatusCreated, N8nResponse{
		Success: true,
		Data:    article,
		URL:     publishURL,
	})
}

// 获取文章
func (h *Handler) GetArticle(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, N8nResponse{
			Success: false,
			Error:   "Invalid article ID",
		})
		return
	}

	article, err := h.articleService.GetArticleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, N8nResponse{
			Success: false,
			Error:   "Article not found",
		})
		return
	}

	c.JSON(http.StatusOK, N8nResponse{
		Success: true,
		Data:    article,
	})
}

// 更新文章
func (h *Handler) UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, N8nResponse{
			Success: false,
			Error:   "Invalid article ID",
		})
		return
	}

	var req struct {
		Title     string     `json:"title"`
		Content   string     `json:"content"`
		Status    string     `json:"status"`
		ExpiresAt *time.Time `json:"expires_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, N8nResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	article, err := h.articleService.UpdateArticle(id, req.Title, req.Content, req.Status, req.ExpiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, N8nResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// 生成发布URL
	publishURL := ""
	if article.Status == "published" {
		publishURL = h.cfg.Server.Domain + "/p/" + article.Slug
	}

	c.JSON(http.StatusOK, N8nResponse{
		Success: true,
		Data:    article,
		URL:     publishURL,
	})
}

// 删除文章
func (h *Handler) DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, N8nResponse{
			Success: false,
			Error:   "Invalid article ID",
		})
		return
	}

	if err := h.articleService.DeleteArticle(id); err != nil {
		c.JSON(http.StatusInternalServerError, N8nResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, N8nResponse{
		Success: true,
	})
}

// 获取文章列表
func (h *Handler) ListArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	articles, total, err := h.articleService.ListArticles(page, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, N8nResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, N8nResponse{
		Success: true,
		Data: gin.H{
			"articles": articles,
			"total":    total,
			"page":     page,
			"limit":    limit,
		},
	})
}

// 获取已发布的文章（公开访问）
func (h *Handler) GetPublishedArticle(c *gin.Context) {
	slug := c.Param("slug")

	article, err := h.articleService.GetPublishedArticleBySlug(slug)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"message": "Article not found",
		})
		return
	}

	// 检查是否过期
	if article.ExpiresAt != nil && article.ExpiresAt.Before(time.Now()) {
		c.HTML(http.StatusGone, "expired.html", gin.H{
			"message": "This article has expired",
		})
		return
	}

	c.HTML(http.StatusOK, "article.html", gin.H{
		"article": article,
	})
}

// API密钥管理
func (h *Handler) CreateAPIKey(c *gin.Context) {
	var req struct {
		Name        string     `json:"name" binding:"required"`
		Permissions string     `json:"permissions"`
		ExpiresAt   *time.Time `json:"expires_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, N8nResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	apiKey, err := h.authService.GenerateAPIKey(req.Name, req.Permissions, req.ExpiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, N8nResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, N8nResponse{
		Success: true,
		Data:    apiKey,
	})
}

func (h *Handler) ListAPIKeys(c *gin.Context) {
	// 实现API密钥列表逻辑
	c.JSON(http.StatusOK, N8nResponse{
		Success: true,
		Data:    []interface{}{},
	})
}

func (h *Handler) DeleteAPIKey(c *gin.Context) {
	// 实现删除API密钥逻辑
	c.JSON(http.StatusOK, N8nResponse{
		Success: true,
	})
}
