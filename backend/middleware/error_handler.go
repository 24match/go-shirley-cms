package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"medical-device-cms/backend/common"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				common.JSONInternalServerError(c, "Internal server error")
			}
		}()

		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			handleError(c, err)
		}
	}
}

func handleError(c *gin.Context, err *gin.Error) {
	statusCode := http.StatusInternalServerError
	errorCode := common.ErrInternalServerError
	message := "Internal server error"

	switch err.Type {
	case gin.ErrorTypeBind:
		statusCode = http.StatusBadRequest
		errorCode = common.ErrBadRequest
		message = "Invalid request parameters"
	case gin.ErrorTypeRender:
		statusCode = http.StatusInternalServerError
		errorCode = common.ErrInternalServerError
		message = "Failed to render response"
	case gin.ErrorTypePrivate:
		statusCode = http.StatusInternalServerError
		errorCode = common.ErrInternalServerError
		message = err.Error()
	case gin.ErrorTypePublic:
		statusCode = http.StatusBadRequest
		errorCode = common.ErrBadRequest
		message = err.Error()
	default:
		statusCode = http.StatusInternalServerError
		errorCode = common.ErrInternalServerError
		message = "Unexpected error"
	}

	c.JSON(statusCode, common.Error(errorCode, message))
}