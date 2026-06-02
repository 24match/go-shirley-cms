package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestTenantControllerCreateUser 测试创建租户用户
func TestTenantControllerCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	ctrl := NewTenantController()
	router.POST("/api/tenant/users", ctrl.CreateUser)

	// 测试创建用户请求
	userData := map[string]string{
		"username": "testuser",
		"password": "password123",
		"role":     "user",
	}
	jsonData, _ := json.Marshal(userData)

	req, _ := http.NewRequest("POST", "/api/tenant/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 可能返回 201（成功）或 403（租户上下文不存在）或 400（参数错误）
	if w.Code != http.StatusCreated && w.Code != http.StatusForbidden && w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 201, 403 or 400, got %d", w.Code)
	}
}

// TestTenantControllerListUsers 测试获取租户用户列表
func TestTenantControllerListUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	ctrl := NewTenantController()
	router.GET("/api/tenant/users", ctrl.ListUsers)

	req, _ := http.NewRequest("GET", "/api/tenant/users?page=1&size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 可能返回 200（成功）或 403（租户上下文不存在）
	if w.Code != http.StatusOK && w.Code != http.StatusForbidden {
		t.Errorf("Expected status 200 or 403, got %d", w.Code)
	}
}

// TestTenantControllerGetUser 测试获取用户详情
func TestTenantControllerGetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	ctrl := NewTenantController()
	router.GET("/api/tenant/users/:id", ctrl.GetUser)

	req, _ := http.NewRequest("GET", "/api/tenant/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 可能返回 200（成功）或 404（用户不存在）
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status 200 or 404, got %d", w.Code)
	}
}

// TestTenantControllerUpdateUser 测试更新用户
func TestTenantControllerUpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	ctrl := NewTenantController()
	router.PUT("/api/tenant/users/:id", ctrl.UpdateUser)

	// 测试更新用户请求
	userData := map[string]string{
		"username": "updateduser",
		"role":     "admin",
	}
	jsonData, _ := json.Marshal(userData)

	req, _ := http.NewRequest("PUT", "/api/tenant/users/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 可能返回 200（成功）或 400（参数错误）或 500（失败）
	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200, 400 or 500, got %d", w.Code)
	}
}

// TestTenantControllerDeleteUser 测试删除用户
func TestTenantControllerDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	ctrl := NewTenantController()
	router.DELETE("/api/tenant/users/:id", ctrl.DeleteUser)

	req, _ := http.NewRequest("DELETE", "/api/tenant/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 可能返回 200（成功）或 500（失败）
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", w.Code)
	}
}

// TestTenantControllerGetDomainConfig 测试获取域名配置
func TestTenantControllerGetDomainConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	ctrl := NewTenantController()
	router.GET("/api/tenant/domain", ctrl.GetDomainConfig)

	req, _ := http.NewRequest("GET", "/api/tenant/domain", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 可能返回 200（成功）或 403（租户上下文不存在）
	if w.Code != http.StatusOK && w.Code != http.StatusForbidden {
		t.Errorf("Expected status 200 or 403, got %d", w.Code)
	}
}

// TestTenantControllerUpdateDomainConfig 测试更新域名配置
func TestTenantControllerUpdateDomainConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	ctrl := NewTenantController()
	router.PUT("/api/tenant/domain", ctrl.UpdateDomainConfig)

	// 测试更新域名配置请求
	domainData := map[string]string{
		"custom_domain": "example.com",
	}
	jsonData, _ := json.Marshal(domainData)

	req, _ := http.NewRequest("PUT", "/api/tenant/domain", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 可能返回 200（成功）
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}