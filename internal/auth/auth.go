package auth

import (
	"net/http"
	"static-hosting-server/internal/config"
	"static-hosting-server/internal/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:  db,
		cfg: cfg,
	}
}

// API Key 认证中间件
func (a *AuthService) APIKeyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			// 尝试从查询参数获取
			apiKey = c.Query("api_key")
		}

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "API key required",
			})
			c.Abort()
			return
		}

		// 验证API密钥
		var dbAPIKey models.APIKey
		if err := a.db.Where("key = ? AND is_active = ?", apiKey, true).First(&dbAPIKey).Error; err != nil {
			// 检查配置中的静态API密钥
			if !a.isStaticAPIKey(apiKey) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Invalid or inactive API key",
				})
				c.Abort()
				return
			}
		} else {
			// 检查API密钥是否过期
			if dbAPIKey.ExpiresAt != nil && dbAPIKey.ExpiresAt.Before(time.Now()) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "API key has expired",
				})
				c.Abort()
				return
			}

			// 更新最后使用时间
			now := time.Now()
			a.db.Model(&dbAPIKey).Update("last_used_at", &now)

			// 将API密钥信息存储在上下文中
			c.Set("api_key", &dbAPIKey)
		}

		c.Next()
	}
}

// 检查是否为配置中的静态API密钥
func (a *AuthService) isStaticAPIKey(apiKey string) bool {
	for _, key := range a.cfg.Security.APIKeys {
		if key == apiKey {
			return true
		}
	}
	return false
}

// 管理员认证中间件（用于后台管理）
func (a *AuthService) AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查会话Cookie
		sessionToken, err := c.Cookie("admin_session")
		if err != nil || sessionToken == "" {
			// 如果是AJAX请求，返回JSON
			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
				c.Abort()
				return
			}
			// 重定向到登录页面
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		// 验证会话token（这里可以实现更复杂的会话验证）
		if !a.validateSessionToken(sessionToken) {
			// 清除无效cookie
			c.SetCookie("admin_session", "", -1, "/", "", false, true)

			if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "会话已过期"})
				c.Abort()
				return
			}
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		c.Next()
	}
}

// 验证会话token
func (a *AuthService) validateSessionToken(token string) bool {
	// 简单的token验证，实际项目中可以使用更安全的方式
	// 这里可以检查数据库中的会话记录、过期时间等
	return token == "valid_admin_session_2025"
}

// 生成会话token
func (a *AuthService) GenerateSessionToken() string {
	// 实际项目中应该生成随机token并存储到数据库
	return "valid_admin_session_2025"
}

// 验证管理员登录凭据
func (a *AuthService) ValidateAdminCredentials(username, password string) bool {
	// 简化的验证逻辑，实际项目中应该查询数据库
	// 并使用bcrypt等方式验证密码
	return username == "admin" && password == "admin123"
}

// 生成新的API密钥
func (a *AuthService) GenerateAPIKey(name string, permissions string, expiresAt *time.Time) (*models.APIKey, error) {
	// 生成随机密钥
	key := generateRandomKey(32)

	apiKey := &models.APIKey{
		Name:        name,
		Key:         key,
		IsActive:    true,
		Permissions: permissions,
		ExpiresAt:   expiresAt,
	}

	if err := a.db.Create(apiKey).Error; err != nil {
		return nil, err
	}

	return apiKey, nil
}

// 生成随机密钥
func generateRandomKey(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	// 使用更好的随机种子
	timeNano := time.Now().UnixNano()
	for i := range result {
		// 使用不同的种子避免相同的字符
		timeNano = timeNano*1103515245 + 12345 // 线性同余生成器
		result[i] = charset[timeNano%int64(len(charset))]
	}
	return string(result)
}
