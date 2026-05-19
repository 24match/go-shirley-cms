package common

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAPIResponseStructure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		response       APIResponse
		expectedCode   int
		expectedMsg    string
		expectData     bool
	}{
		{
			name:           "success response with data",
			response:       Success(map[string]string{"key": "value"}),
			expectedCode:   SuccessCode,
			expectedMsg:    "success",
			expectData:     true,
		},
		{
			name:           "success response with custom message",
			response:       SuccessWithMessage(nil, "Created successfully"),
			expectedCode:   SuccessCode,
			expectedMsg:    "Created successfully",
			expectData:     false,
		},
		{
			name:           "bad request error",
			response:       BadRequest("Invalid input"),
			expectedCode:   ErrBadRequest,
			expectedMsg:    "Invalid input",
			expectData:     false,
		},
		{
			name:           "unauthorized error",
			response:       Unauthorized("Login required"),
			expectedCode:   ErrUnauthorized,
			expectedMsg:    "Login required",
			expectData:     false,
		},
		{
			name:           "not found error",
			response:       NotFound("Resource not found"),
			expectedCode:   ErrNotFound,
			expectedMsg:    "Resource not found",
			expectData:     false,
		},
		{
			name:           "internal server error",
			response:       InternalServerError("Server error"),
			expectedCode:   ErrInternalServerError,
			expectedMsg:    "Server error",
			expectData:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedCode, tc.response.Code)
			assert.Equal(t, tc.expectedMsg, tc.response.Message)
			assert.NotZero(t, tc.response.Timestamp)
			assert.True(t, tc.response.Timestamp <= time.Now().Unix())

			if tc.expectData {
				assert.NotNil(t, tc.response.Data)
			} else {
				assert.Nil(t, tc.response.Data)
			}
		})
	}
}

func TestJSONResponseFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		handlerFunc    gin.HandlerFunc
		expectedCode   int
		expectedHTTP   int
	}{
		{
			name: "JSONSuccess",
			handlerFunc: func(c *gin.Context) {
				JSONSuccess(c, map[string]string{"key": "value"})
			},
			expectedCode: SuccessCode,
			expectedHTTP: http.StatusOK,
		},
		{
			name: "JSONBadRequest",
			handlerFunc: func(c *gin.Context) {
				JSONBadRequest(c, "Bad request")
			},
			expectedCode: ErrBadRequest,
			expectedHTTP: http.StatusBadRequest,
		},
		{
			name: "JSONUnauthorized",
			handlerFunc: func(c *gin.Context) {
				JSONUnauthorized(c, "Unauthorized")
			},
			expectedCode: ErrUnauthorized,
			expectedHTTP: http.StatusUnauthorized,
		},
		{
			name: "JSONNotFound",
			handlerFunc: func(c *gin.Context) {
				JSONNotFound(c, "Not found")
			},
			expectedCode: ErrNotFound,
			expectedHTTP: http.StatusNotFound,
		},
		{
			name: "JSONInternalServerError",
			handlerFunc: func(c *gin.Context) {
				JSONInternalServerError(c, "Internal error")
			},
			expectedCode: ErrInternalServerError,
			expectedHTTP: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			r := gin.New()
			r.GET("/test", tc.handlerFunc)
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedHTTP, w.Code)

			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCode, response.Code)
			assert.NotZero(t, response.Timestamp)
		})
	}
}

func TestErrorStatusCodeMapping(t *testing.T) {
	tests := []struct {
		name     string
		errCode  int
		httpCode int
	}{
		{"success", SuccessCode, http.StatusOK},
		{"bad request", ErrBadRequest, http.StatusBadRequest},
		{"unauthorized", ErrUnauthorized, http.StatusUnauthorized},
		{"forbidden", ErrForbidden, http.StatusForbidden},
		{"not found", ErrNotFound, http.StatusNotFound},
		{"conflict", ErrConflict, http.StatusConflict},
		{"internal error", ErrInternalServerError, http.StatusInternalServerError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			r := gin.New()
			r.GET("/test", func(c *gin.Context) {
				JSON(c, tc.httpCode, Error(tc.errCode, "test"))
			})
			r.ServeHTTP(w, req)

			assert.Equal(t, tc.httpCode, w.Code)

			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.errCode, response.Code)
		})
	}
}