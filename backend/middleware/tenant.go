package middleware

import (
	"net/http"
	"strings"

	"medical-device-cms/backend/config"
	"medical-device-cms/backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	// TenantContextKey 租户上下文键
	TenantContextKey = "tenant"
	// TenantIDContextKey 租户 ID 上下文键
	TenantIDContextKey = "tenant_id"
	// TenantCodeContextKey 租户代码上下文键
	TenantCodeContextKey = "tenant_code"
)

// TenantMiddleware 租户识别中间件
// 从域名或路径中识别租户，并将租户信息注入到上下文中
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要租户识别的路径
		if shouldSkipTenantCheck(c.Request.URL.Path) {
			c.Next()
			return
		}

		tenantCode := extractTenantCode(c)

		if tenantCode == "" {
			// 如果没有找到租户代码，尝试从 Header 获取
			tenantCode = c.GetHeader("X-Tenant-Code")
		}

		if tenantCode == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    "TENANT_CODE_REQUIRED",
				"message": "租户代码是必需的",
			})
			return
		}

		// 查询租户信息
		var tenant models.Tenant
		if err := config.DB.Where("tenant_code = ?", tenantCode).First(&tenant).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"code":    "TENANT_NOT_FOUND",
				"message": "租户不存在",
			})
			return
		}

		// 检查租户状态
		if tenant.Status == models.TenantStatusDisabled {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    "TENANT_DISABLED",
				"message": "租户已被禁用",
			})
			return
		}

		if tenant.Status == models.TenantStatusSuspended {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    "TENANT_SUSPENDED",
				"message": "租户已被暂停",
			})
			return
		}

		if tenant.Status == models.TenantStatusPending {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    "TENANT_PENDING",
				"message": "租户尚未激活",
			})
			return
		}

		// 将租户信息注入到上下文中
		c.Set(TenantContextKey, &tenant)
		c.Set(TenantIDContextKey, tenant.ID)
		c.Set(TenantCodeContextKey, tenant.TenantCode)

		c.Next()
	}
}

// shouldSkipTenantCheck 检查是否应该跳过租户检查
func shouldSkipTenantCheck(path string) bool {
	// 超级管理员 API 不需要租户识别
	if strings.HasPrefix(path, "/api/superadmin/") {
		return true
	}

	// 认证端点不需要租户识别
	if path == "/api/login" || path == "/api/register" {
		return true
	}

	// Swagger 文档不需要租户识别
	if strings.HasPrefix(path, "/swagger/") {
		return true
	}

	return false
}

// extractTenantCode 从请求中提取租户代码
func extractTenantCode(c *gin.Context) string {
	// 1. 首先尝试从子域名提取
	host := c.Request.Host
	parts := strings.Split(host, ".")
	if len(parts) >= 3 {
		// 假设域名格式为：tenant.example.com
		// 排除 www 和 api 等常见前缀
		subdomain := parts[0]
		if subdomain != "www" && subdomain != "api" && subdomain != "localhost" {
			return subdomain
		}
	}

	// 2. 尝试从路径前缀提取 (如 /tenant-code/...)
	path := strings.TrimPrefix(c.Request.URL.Path, "/")
	parts = strings.Split(path, "/")
	if len(parts) > 0 && parts[0] != "" && parts[0] != "api" && parts[0] != "admin" && parts[0] != "frontend" {
		// 验证是否是有效的租户代码格式
		if isValidTenantCode(parts[0]) {
			return parts[0]
		}
	}

	return ""
}

// isValidTenantCode 验证是否是有效的租户代码格式
func isValidTenantCode(code string) bool {
	if len(code) == 0 || len(code) > 50 {
		return false
	}
	// 租户代码只能包含小写字母、数字和连字符
	for _, r := range code {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-') {
			return false
		}
	}
	return true
}

// GetTenantFromContext 从上下文获取租户信息
func GetTenantFromContext(c *gin.Context) (*models.Tenant, bool) {
	tenant, exists := c.Get(TenantContextKey)
	if !exists {
		return nil, false
	}
	t, ok := tenant.(*models.Tenant)
	return t, ok
}

// GetTenantIDFromContext 从上下文获取租户 ID
func GetTenantIDFromContext(c *gin.Context) (uint, bool) {
	tenantID, exists := c.Get(TenantIDContextKey)
	if !exists {
		return 0, false
	}
	id, ok := tenantID.(uint)
	return id, ok
}

// GetTenantCodeFromContext 从上下文获取租户代码
func GetTenantCodeFromContext(c *gin.Context) (string, bool) {
	tenantCode, exists := c.Get(TenantCodeContextKey)
	if !exists {
		return "", false
	}
	code, ok := tenantCode.(string)
	return code, ok
}

// SetTenantContext 设置租户上下文（用于超级管理员切换租户）
func SetTenantContext(c *gin.Context, tenant *models.Tenant) {
	c.Set(TenantContextKey, tenant)
	c.Set(TenantIDContextKey, tenant.ID)
	c.Set(TenantCodeContextKey, tenant.TenantCode)
}

// WithTenantScope GORM Scope 用于自动添加租户过滤
func WithTenantScope(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		tenantID, exists := GetTenantIDFromContext(c)
		if !exists {
			// 如果没有租户上下文，返回空结果
			return db.Where("1 = 0")
		}
		return db.Where("tenant_id = ?", tenantID)
	}
}