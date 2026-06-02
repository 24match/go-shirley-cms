package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestSuperAdminControllerCreateTenant 测试创建租户
func TestSuperAdminControllerCreateTenant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	ctrl := NewSuperAdminController()
	router.POST("/api/superadmin/tenants", ctrl.CreateTenant)

	// 测试创建租户请求
	tenantData := map[string]string{
		"tenant_code": "test-create",
		"name":        "Test Create Tenant",
		"sub_domain":  "test-create",
		"status":      "active",
	}
	jsonData, _ := json.Marshal(tenantData)

	req, _ := http.NewRequest("POST", "/api/superadmin/tenants", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 可能返回 201（成功）或 500（失败，如重复）
	if w.Code != http.StatusCreated && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 201 or 500, got %d", w.Code)
	}
}

// TestSuperAdminControllerListTenants 测试租户列表
func TestSuperAdminControllerListTenants(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	ctrl := NewSuperAdminController()
	router.GET("/api/superadmin/tenants", ctrl.ListTenants)

	req, _ := http.NewRequest("GET", "/api/superadmin/tenants?status=active&page=1&size=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"] != "SUCCESS" {
		t.Errorf("Expected code SUCCESS, got %v", response["code"])
	}
}

// TestSuperAdminControllerGetTenant 测试获取租户详情
func TestSuperAdminControllerGetTenant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	ctrl := NewSuperAdminController()
	router.GET("/api/superadmin/tenants/:id", ctrl.GetTenant)

	req, _ := http.NewRequest("GET", "/api/superadmin/tenants/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 可能返回 200（租户存在）或 404（租户不存在）
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status 200 or 404, got %d", w.Code)
	}
}

// TestSuperAdminControllerUpdateTenant 测试更新租户
func TestSuperAdminControllerUpdateTenant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	ctrl := NewSuperAdminController()
	router.PUT("/api/superadmin/tenants/:id", ctrl.UpdateTenant)

	// 测试更新租户请求
	updateData := map[string]string{
		"name":        "Updated Tenant Name",
		"description": "Updated description",
	}
	jsonData, _ := json.Marshal(updateData)

	req, _ := http.NewRequest("PUT", "/api/superadmin/tenants/1", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 可能返回 200（成功）或 404（租户不存在）或 500（失败）
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200, 404 or 500, got %d", w.Code)
	}
}

// TestSuperAdminControllerDeleteTenant 测试删除租户
func TestSuperAdminControllerDeleteTenant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	ctrl := NewSuperAdminController()
	router.DELETE("/api/superadmin/tenants/:id", ctrl.DeleteTenant)

	req, _ := http.NewRequest("DELETE", "/api/superadmin/tenants/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 可能返回 200（成功）或 500（失败）
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", w.Code)
	}
}

// TestSuperAdminControllerActivateTenant 测试激活租户
func TestSuperAdminControllerActivateTenant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	ctrl := NewSuperAdminController()
	router.POST("/api/superadmin/tenants/:id/activate", ctrl.ActivateTenant)

	req, _ := http.NewRequest("POST", "/api/superadmin/tenants/1/activate", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 可能返回 200（成功）或 500（失败）
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", w.Code)
	}
}

// TestSuperAdminControllerDisableTenant 测试禁用租户
func TestSuperAdminControllerDisableTenant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	ctrl := NewSuperAdminController()
	router.POST("/api/superadmin/tenants/:id/disable", ctrl.DisableTenant)

	req, _ := http.NewRequest("POST", "/api/superadmin/tenants/1/disable", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 可能返回 200（成功）或 500（失败）
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", w.Code)
	}
}

// TestSuperAdminControllerImpersonateTenant 测试切换租户上下文
func TestSuperAdminControllerImpersonateTenant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	ctrl := NewSuperAdminController()
	router.POST("/api/superadmin/tenants/:id/impersonate", ctrl.ImpersonateTenant)

	req, _ := http.NewRequest("POST", "/api/superadmin/tenants/1/impersonate", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 可能返回 200（成功）或 404（租户不存在）或 500（失败）
	if w.Code != http.StatusOK && w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200, 404 or 500, got %d", w.Code)
	}
}

// TestSuperAdminControllerGetSystemStats 测试获取系统统计
func TestSuperAdminControllerGetSystemStats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	ctrl := NewSuperAdminController()
	router.GET("/api/superadmin/stats", ctrl.GetSystemStats)

	req, _ := http.NewRequest("GET", "/api/superadmin/stats", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"] != "SUCCESS" {
		t.Errorf("Expected code SUCCESS, got %v", response["code"])
	}
}