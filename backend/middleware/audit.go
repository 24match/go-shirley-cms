package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"medical-device-cms/backend/services"

	"github.com/gin-gonic/gin"
)

// responseWriter 自定义响应写入器，用于捕获响应内容
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// AuditMiddleware 审计日志中间件
// 记录所有 API 请求的操作日志
func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要审计的路径
		if shouldSkipAudit(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 记录请求开始时间
		startTime := time.Now()

		// 创建自定义响应写入器
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = writer

		// 获取操作人 ID（从认证中间件注入的 claims）
		operatorID := getOperatorID(c)

		// 获取租户 ID
		tenantID, _ := GetTenantIDFromContext(c)

		// 处理请求
		c.Next()

		// 计算处理时间
		duration := time.Since(startTime)

		// 获取请求体
		requestBody := getRequestBody(c)

		// 获取响应体
		responseBody := writer.body.String()

		// 确定操作类型
		action := determineAction(c.Request.Method, c.Writer.Status())

		// 确定资源类型
		resourceType := determineResourceType(c.Request.URL.Path)

		// 异步记录审计日志（不阻塞请求）
		go func() {
			auditService := services.NewAuditService()
			
			// 记录审计日志
			_ = auditService.LogAction(
				operatorID,
				tenantID,
				action,
				resourceType,
				0, // 资源 ID（需要从路径中提取，这里简化处理）
				requestBody,
				responseBody,
			)
		}()

		// 记录请求日志
		log.Printf("[%s] %s %s %d %v %s",
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
			c.Request.UserAgent(),
		)
	}
}

// shouldSkipAudit 检查是否应该跳过审计
func shouldSkipAudit(path string) bool {
	// 跳过静态文件
	if strings.HasPrefix(path, "/frontend/") ||
		strings.HasPrefix(path, "/uploads/") ||
		strings.HasPrefix(path, "/admin/") {
		return true
	}

	// 跳过 Swagger 文档
	if strings.HasPrefix(path, "/swagger/") {
		return true
	}

	// 跳过健康检查
	if path == "/health" || path == "/ready" {
		return true
	}

	return false
}

// getOperatorID 从上下文获取操作人 ID
func getOperatorID(c *gin.Context) uint {
	// 尝试从上下文获取用户 ID
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}

	// 尝试从 claims 获取
	if claims, exists := c.Get("claims"); exists {
		if claimMap, ok := claims.(map[string]interface{}); ok {
			if id, ok := claimMap["user_id"].(float64); ok {
				return uint(id)
			}
		}
	}

	return 0 // 匿名用户
}

// getRequestBody 获取请求体
func getRequestBody(c *gin.Context) string {
	// 只记录 JSON 请求体
	contentType := c.GetHeader("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return ""
	}

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}

	// 恢复请求体以便后续处理
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	return string(body)
}

// determineAction 根据 HTTP 方法确定操作类型
func determineAction(method string, statusCode int) string {
	switch method {
	case http.MethodPost:
		if statusCode == http.StatusCreated || statusCode == http.StatusOK {
			return "CREATE"
		}
		return "CREATE_FAILED"
	case http.MethodPut, http.MethodPatch:
		if statusCode == http.StatusOK {
			return "UPDATE"
		}
		return "UPDATE_FAILED"
	case http.MethodDelete:
		if statusCode == http.StatusOK || statusCode == http.StatusNoContent {
			return "DELETE"
		}
		return "DELETE_FAILED"
	case http.MethodGet:
		return "READ"
	default:
		return "UNKNOWN"
	}
}

// determineResourceType 根据路径确定资源类型
func determineResourceType(path string) string {
	// 提取 API 路径中的资源类型
	// 例如：/api/admin/images -> image
	// /api/admin/content -> content
	// /api/admin/modules -> module
	// /api/admin/lang -> language
	// /api/admin/users -> user

	parts := strings.Split(strings.TrimPrefix(path, "/api/"), "/")
	if len(parts) >= 2 {
		resource := parts[1]
		// 标准化资源名称
		switch {
		case strings.HasPrefix(resource, "image"):
			return "image"
		case strings.HasPrefix(resource, "content"):
			return "content"
		case strings.HasPrefix(resource, "module"):
			return "module"
		case strings.HasPrefix(resource, "lang"):
			return "language"
		case strings.HasPrefix(resource, "user"):
			return "user"
		case strings.HasPrefix(resource, "tenant"):
			return "tenant"
		case strings.HasPrefix(resource, "contact"):
			return "contact"
		case strings.HasPrefix(resource, "site-setting"):
			return "site_setting"
		default:
			return resource
		}
	}

	return "unknown"
}

// LogManualAction 手动记录审计日志的辅助函数
// 用于在业务逻辑中需要记录特殊操作时使用
func LogManualAction(c *gin.Context, action string, resourceType string, resourceID uint, beforeValue, afterValue interface{}) {
	operatorID := getOperatorID(c)
	tenantID, _ := GetTenantIDFromContext(c)

	auditService := services.NewAuditService()
	_ = auditService.LogAction(
		operatorID,
		tenantID,
		action,
		resourceType,
		resourceID,
		beforeValue,
		afterValue,
	)
}