package common

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// APIResponse 统一 API 响应结构
// @Description 所有 API 响应的统一包装格式，包含状态码、消息、数据和 timestamp
type APIResponse struct {
	// 状态码，200 表示成功，其他表示错误
	Code int `json:"code" example:"200"`
	// 消息提示，成功时为"success"，错误时为错误描述
	Message string `json:"message" example:"success"`
	// 响应数据，成功时返回业务数据，错误时为 null
	Data interface{} `json:"data,omitempty"`
	// Unix 时间戳（秒）
	Timestamp int64 `json:"timestamp" example:"1717200000"`
}

// 错误码常量定义
const (
	// SuccessCode 成功状态码
	SuccessCode = 200

	// ErrBadRequest 请求参数错误
	ErrBadRequest = 400
	// ErrUnauthorized 未授权访问
	ErrUnauthorized = 401
	// ErrForbidden 禁止访问
	ErrForbidden = 403
	// ErrNotFound 资源不存在
	ErrNotFound = 404
	// ErrMethodNotAllowed 请求方法不允许
	ErrMethodNotAllowed = 405
	// ErrConflict 资源冲突
	ErrConflict = 409

	// ErrInternalServerError 服务器内部错误
	ErrInternalServerError = 500
	// ErrDatabaseError 数据库错误
	ErrDatabaseError = 501
	// ErrValidationError 数据验证错误
	ErrValidationError = 502
	// ErrServiceUnavailable 服务不可用
	ErrServiceUnavailable = 503

	// 租户相关错误码
	// ErrTenantNotFound 租户不存在
	ErrTenantNotFound = 1001
	// ErrTenantDisabled 租户已被禁用
	ErrTenantDisabled = 1002
	// ErrTenantSuspended 租户已被暂停
	ErrTenantSuspended = 1003
	// ErrTenantPending 租户尚未激活
	ErrTenantPending = 1004
	// ErrTenantCodeRequired 租户代码是必需的
	ErrTenantCodeRequired = 1005
	// ErrCrossTenantAccess 跨租户访问被拒绝
	ErrCrossTenantAccess = 1006

	// 超级管理员相关错误码
	// ErrSuperAdminRequired 需要超级管理员权限
	ErrSuperAdminRequired = 2001
)

func Success(data interface{}) APIResponse {
	return APIResponse{
		Code:      SuccessCode,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

func SuccessWithMessage(data interface{}, message string) APIResponse {
	return APIResponse{
		Code:      SuccessCode,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

func Error(code int, message string) APIResponse {
	return APIResponse{
		Code:      code,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().Unix(),
	}
}

func BadRequest(message string) APIResponse {
	return Error(ErrBadRequest, message)
}

func Unauthorized(message string) APIResponse {
	return Error(ErrUnauthorized, message)
}

func Forbidden(message string) APIResponse {
	return Error(ErrForbidden, message)
}

func NotFound(message string) APIResponse {
	return Error(ErrNotFound, message)
}

func InternalServerError(message string) APIResponse {
	return Error(ErrInternalServerError, message)
}

func DatabaseError(message string) APIResponse {
	return Error(ErrDatabaseError, message)
}

func ValidationError(message string) APIResponse {
	return Error(ErrValidationError, message)
}

func JSON(c *gin.Context, httpStatus int, response APIResponse) {
	c.JSON(httpStatus, response)
}

func JSONSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Success(data))
}

func JSONSuccessWithMessage(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, SuccessWithMessage(data, message))
}

func JSONError(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Error(code, message))
}

func JSONBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, BadRequest(message))
}

func JSONUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Unauthorized(message))
}

func JSONForbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Forbidden(message))
}

func JSONNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, NotFound(message))
}

func JSONInternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, InternalServerError(message))
}

// 租户相关错误响应
func JSONTenantNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Error(ErrTenantNotFound, message))
}

func JSONTenantDisabled(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Error(ErrTenantDisabled, message))
}

func JSONTenantSuspended(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Error(ErrTenantSuspended, message))
}

func JSONTenantPending(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Error(ErrTenantPending, message))
}

func JSONTenantCodeRequired(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Error(ErrTenantCodeRequired, message))
}

func JSONCrossTenantAccess(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Error(ErrCrossTenantAccess, message))
}

func JSONSuperAdminRequired(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Error(ErrSuperAdminRequired, message))
}