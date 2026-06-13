package middleware

import (
	"net/http"

	"medical-device-cms/backend/services"
	"github.com/gin-gonic/gin"
)

// QuotaCheckMiddleware 配额检查中间件
// resourceType: 资源类型，如 "images", "content_items", "users", "storage_mb"
func QuotaCheckMiddleware(resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取租户 ID
		tenantID, exists := GetTenantIDFromContext(c)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    "TENANT_REQUIRED",
				"message": "租户上下文不存在",
			})
			c.Abort()
			return
		}

		// 检查是否为超级管理员（超级管理员不受配额限制）
		role, roleExists := c.Get("role")
		if roleExists && role == "superadmin" {
			c.Next()
			return
		}

		// 检查配额
		quotaService := services.NewQuotaService()
		isSufficient, _, err := quotaService.CheckQuota(c, tenantID, resourceType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "QUOTA_CHECK_FAILED",
				"message": "配额检查失败",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		if !isSufficient {
			// 获取配额信息用于返回详细错误
			_, info, _ := quotaService.CheckQuota(c, tenantID, resourceType)
			c.JSON(http.StatusForbidden, gin.H{
				"code":    "QUOTA_EXCEEDED",
				"message": "资源配额已超限",
				"details": gin.H{
					"resource_type": resourceType,
					"used":          info.Used,
					"limit":         info.Limit,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// QuotaIncrementMiddleware 配额增加中间件
// 在请求成功后自动增加资源使用量
// resourceType: 资源类型
// amount: 增加数量，默认为 1
func QuotaIncrementMiddleware(resourceType string, amount ...int) gin.HandlerFunc {
	incrementAmount := 1
	if len(amount) > 0 {
		incrementAmount = amount[0]
	}

	return func(c *gin.Context) {
		c.Next()

		// 只在请求成功后增加配额
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			tenantID, exists := GetTenantIDFromContext(c)
			if !exists {
				return
			}

			// 跳过超级管理员
			role, roleExists := c.Get("role")
			if roleExists && role == "superadmin" {
				return
			}

			// 增加使用量
			tenantConfigService := services.NewTenantConfigService()
			tenantConfigService.IncrementUsage(c, tenantID, resourceType, incrementAmount)
		}
	}
}