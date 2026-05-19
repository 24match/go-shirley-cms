package common

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

const (
	SuccessCode = 200

	ErrBadRequest      = 400
	ErrUnauthorized    = 401
	ErrForbidden       = 403
	ErrNotFound        = 404
	ErrMethodNotAllowed = 405
	ErrConflict        = 409

	ErrInternalServerError = 500
	ErrDatabaseError       = 501
	ErrValidationError     = 502
	ErrServiceUnavailable  = 503
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