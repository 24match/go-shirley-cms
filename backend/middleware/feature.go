package middleware

import (
	"net/http"

	"medical-device-cms/backend/services"
	"github.com/gin-gonic/gin"
)

// FeatureCheckMiddleware 功能模块权限检查中间件
// featureName: 功能模块名称，如 "image_management", "page_config", "multi_language" 等
func FeatureCheckMiddleware(featureName string) gin.HandlerFunc {
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

		// 检查是否为超级管理员（超级管理员不受功能开关限制）
		role, roleExists := c.Get("role")
		if roleExists && role == "superadmin" {
			c.Next()
			return
		}

		// 检查功能开关
		tenantConfigService := services.NewTenantConfigService()
		enabled, err := tenantConfigService.CheckFeature(c, tenantID, featureName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "FEATURE_CHECK_FAILED",
				"message": "功能权限检查失败",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		if !enabled {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    "FEATURE_DISABLED",
				"message": "该功能模块当前不可用，请联系管理员开通",
				"details": gin.H{
					"feature": featureName,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyFeature 要求至少启用一个功能的中间件
// featureNames: 功能模块名称列表
func RequireAnyFeature(featureNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, exists := GetTenantIDFromContext(c)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    "TENANT_REQUIRED",
				"message": "租户上下文不存在",
			})
			c.Abort()
			return
		}

		// 检查是否为超级管理员
		role, roleExists := c.Get("role")
		if roleExists && role == "superadmin" {
			c.Next()
			return
		}

		tenantConfigService := services.NewTenantConfigService()
		
		// 检查是否至少有一个功能启用
		for _, featureName := range featureNames {
			enabled, err := tenantConfigService.CheckFeature(c, tenantID, featureName)
			if err == nil && enabled {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"code":    "FEATURE_DISABLED",
			"message": "所需功能模块当前不可用，请联系管理员开通",
		})
		c.Abort()
	}
}