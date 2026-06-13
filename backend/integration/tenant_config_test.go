package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"medical-device-cms/backend/config"
	"medical-device-cms/backend/models"
	"medical-device-cms/backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 测试辅助：创建测试数据库
func setupIntegrationTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移所有表结构
	db.AutoMigrate(
		&models.User{},
		&models.Tenant{},
		&models.TenantConfig{},
		&models.Image{},
		&models.ContentItem{},
		&models.ContactSubmission{},
		&models.ModuleConfig{},
		&models.PageConfig{},
		&models.LanguageText{},
		&models.SiteSetting{},
		&models.AuditLog{},
	)

	return db
}

// 测试辅助：生成 JWT Token
func generateTestToken(userID uint, role string, tenantID uint) string {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"role":      role,
		"tenant_id": tenantID,
		"exp":       time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-secret-key"))
	return tokenString
}

// 测试辅助：创建测试用户
func createTestUser(db *gorm.DB, username string, role string, tenantID uint) *models.User {
	user := models.User{
		TenantID: tenantID,
		Username: username,
		Password: "hashed_password",
		Role:     role,
	}
	db.Create(&user)
	return &user
}

// 测试辅助：创建测试租户
func createTestTenant(db *gorm.DB, code string) *models.Tenant {
	tenant := models.Tenant{
		TenantCode: code,
		Name:       "Test Tenant " + code,
		Status:     models.TenantStatusActive,
	}
	db.Create(&tenant)
	return &tenant
}

// TestTenantConfigAutoInit 测试租户创建时配置自动初始化
func TestTenantConfigAutoInit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	config.DB = setupIntegrationTestDB()

	// 创建租户（应该自动创建配置）
	tenant := createTestTenant(config.DB, "auto-init-tenant")

	// 验证配置已自动创建
	var tenantConfig models.TenantConfig
	err := config.DB.Where("tenant_id = ?", tenant.ID).First(&tenantConfig).Error

	if err != nil {
		t.Errorf("Expected TenantConfig to be auto-created, got error: %v", err)
	}

	if tenantConfig.TenantID != tenant.ID {
		t.Errorf("Expected TenantID %d, got %d", tenant.ID, tenantConfig.TenantID)
	}

	if tenantConfig.SubscriptionPlan != "free" {
		t.Errorf("Expected default SubscriptionPlan 'free', got '%s'", tenantConfig.SubscriptionPlan)
	}

	// 验证默认功能标志
	var featureFlags map[string]bool
	json.Unmarshal([]byte(tenantConfig.FeatureFlags), &featureFlags)

	if !featureFlags["image_management"] {
		t.Error("Expected image_management to be enabled by default")
	}

	if !featureFlags["page_config"] {
		t.Error("Expected page_config to be enabled by default")
	}

	// 验证默认配额
	var resourceQuota map[string]int
	json.Unmarshal([]byte(tenantConfig.ResourceQuota), &resourceQuota)

	if resourceQuota["max_images"] != 50 {
		t.Errorf("Expected max_images 50 for free plan, got %d", resourceQuota["max_images"])
	}

	if resourceQuota["max_storage_mb"] != 512 {
		t.Errorf("Expected max_storage_mb 512 for free plan, got %d", resourceQuota["max_storage_mb"])
	}
}

// TestQuotaLimitImageUpload 测试配额限制功能（上传图片超限场景）
func TestQuotaLimitImageUpload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	config.DB = setupIntegrationTestDB()

	// 创建租户
	tenant := createTestTenant(config.DB, "quota-limit-tenant")

	// 创建租户管理员用户
	tenantAdmin := createTestUser(config.DB, "tenant@quota", "tenant_admin", tenant.ID)

	// 创建配置，设置图片配额为 5
	featureFlags := map[string]bool{
		"image_management": true,
		"page_config":      true,
	}
	resourceQuota := map[string]int{
		"max_images":     5,
		"max_storage_mb": 1024,
	}
	resourceUsage := map[string]int{
		"used_images":     4,
		"used_storage_mb": 100,
	}

	flagsJSON, _ := json.Marshal(featureFlags)
	quotaJSON, _ := json.Marshal(resourceQuota)
	usageJSON, _ := json.Marshal(resourceUsage)

	tenantConfig := models.TenantConfig{
		TenantID:         tenant.ID,
		FeatureFlags:     string(flagsJSON),
		ResourceQuota:    string(quotaJSON),
		ResourceUsage:    string(usageJSON),
		SubscriptionPlan: "free",
	}
	config.DB.Create(&tenantConfig)

	// 创建路由
	r := gin.Default()
	routes.RegisterRoutes(r)

	// 更新使用量到上限
	config.DB.Model(&models.TenantConfig{}).Where("tenant_id = ?", tenant.ID).Update("resource_usage", `{"used_images":5,"used_storage_mb":100}`)

	// 测试：上传第 6 张图片（应该失败）
	t.Run("UploadImage_ExceedQuota", func(t *testing.T) {
		token := generateTestToken(tenantAdmin.ID, "tenant_admin", tenant.ID)

		body := bytes.NewBufferString("fake-image-data")
		req, _ := http.NewRequest("POST", "/api/admin/images", body)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "multipart/form-data")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// 应该返回 403 配额超限
		if w.Code != 403 {
			t.Errorf("Expected status 403 for quota exceeded, got %d", w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "QUOTA_EXCEEDED" {
			t.Errorf("Expected code QUOTA_EXCEEDED, got %v", response["code"])
		}
	})
}

// TestFeatureSwitchEffect 测试功能开关效果
func TestFeatureSwitchEffect(t *testing.T) {
	gin.SetMode(gin.TestMode)

	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	config.DB = setupIntegrationTestDB()

	// 创建租户
	tenant := createTestTenant(config.DB, "feature-switch-tenant")

	// 创建租户管理员用户
	tenantAdmin := createTestUser(config.DB, "tenant@feature", "tenant_admin", tenant.ID)

	// 创建配置，禁用图片管理功能
	featureFlags := map[string]bool{
		"image_management":   false,
		"page_config":        true,
		"multi_language":     true,
		"contact_form":       true,
		"content_management": true,
	}
	resourceQuota := map[string]int{
		"max_images":     100,
		"max_storage_mb": 1024,
	}
	resourceUsage := map[string]int{
		"used_images":     0,
		"used_storage_mb": 0,
	}

	flagsJSON, _ := json.Marshal(featureFlags)
	quotaJSON, _ := json.Marshal(resourceQuota)
	usageJSON, _ := json.Marshal(resourceUsage)

	tenantConfig := models.TenantConfig{
		TenantID:         tenant.ID,
		FeatureFlags:     string(flagsJSON),
		ResourceQuota:    string(quotaJSON),
		ResourceUsage:    string(usageJSON),
		SubscriptionPlan: "free",
	}
	config.DB.Create(&tenantConfig)

	// 创建路由
	r := gin.Default()
	routes.RegisterRoutes(r)

	token := generateTestToken(tenantAdmin.ID, "tenant_admin", tenant.ID)

	// 测试：访问图片管理 API（应该被禁用）
	t.Run("AccessImageManagement_Disabled", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/admin/images", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// 应该返回 403 功能禁用
		if w.Code != 403 {
			t.Errorf("Expected status 403 for disabled feature, got %d", w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "FEATURE_DISABLED" {
			t.Errorf("Expected code FEATURE_DISABLED, got %v", response["code"])
		}
	})

	// 启用图片管理功能
	featureFlags["image_management"] = true
	flagsJSON, _ = json.Marshal(featureFlags)
	config.DB.Model(&models.TenantConfig{}).Where("tenant_id = ?", tenant.ID).Update("feature_flags", string(flagsJSON))

	// 测试：访问图片管理 API（现在应该允许）
	t.Run("AccessImageManagement_Enabled", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/admin/images", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// 注意：由于没有实际数据，可能返回 200 空列表或其他状态
		// 但不应返回 FEATURE_DISABLED
		if w.Code == 403 {
			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			if response["code"] == "FEATURE_DISABLED" {
				t.Error("Expected feature to be enabled, but got FEATURE_DISABLED")
			}
		}
	})
}

// TestSuperAdminConfigTenant 测试超级管理员配置租户完整流程
func TestSuperAdminConfigTenant(t *testing.T) {
	gin.SetMode(gin.TestMode)

	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	config.DB = setupIntegrationTestDB()

	// 创建超级管理员
	superAdmin := createTestUser(config.DB, "superadmin", "superadmin", 0)

	// 创建租户
	tenant := createTestTenant(config.DB, "config-test-tenant")

	// 创建默认配置
	tenantConfig := models.TenantConfig{
		TenantID:         tenant.ID,
		FeatureFlags:     `{"image_management":true,"page_config":true}`,
		ResourceQuota:    `{"max_images":50,"max_storage_mb":512}`,
		ResourceUsage:    `{"used_images":0,"used_storage_mb":0}`,
		SubscriptionPlan: "free",
	}
	config.DB.Create(&tenantConfig)

	// 创建路由
	r := gin.Default()
	routes.RegisterRoutes(r)

	token := generateTestToken(superAdmin.ID, "superadmin", 0)

	// 测试：获取租户配置
	t.Run("GetTenantConfig", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/superadmin/tenants/%d/config", tenant.ID), nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["data"] == nil {
			t.Error("Expected data in response")
		}
	})

	// 测试：更新租户配置
	t.Run("UpdateTenantConfig", func(t *testing.T) {
		updateReq := map[string]interface{}{
			"feature_flags": map[string]bool{
				"image_management":   true,
				"page_config":        true,
				"multi_language":     true,
				"contact_form":       true,
				"content_management": true,
			},
			"resource_quota": map[string]int{
				"max_images":        500,
				"max_storage_mb":    5120,
				"max_content_items": 100,
				"max_users":         10,
			},
			"subscription_plan": "pro",
		}

		body, _ := json.Marshal(updateReq)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/superadmin/tenants/%d/config", tenant.ID), bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// 验证配置已更新
		var updatedConfig models.TenantConfig
		config.DB.Where("tenant_id = ?", tenant.ID).First(&updatedConfig)

		if updatedConfig.SubscriptionPlan != "pro" {
			t.Errorf("Expected SubscriptionPlan 'pro', got '%s'", updatedConfig.SubscriptionPlan)
		}

		var updatedFlags map[string]bool
		json.Unmarshal([]byte(updatedConfig.FeatureFlags), &updatedFlags)

		if !updatedFlags["multi_language"] {
			t.Error("Expected multi_language to be enabled after update")
		}
	})

	// 测试：重置配额
	t.Run("ResetQuota", func(t *testing.T) {
		// 先设置一些使用量
		config.DB.Model(&models.TenantConfig{}).Where("tenant_id = ?", tenant.ID).Update("resource_usage", `{"used_images":25,"used_storage_mb":256}`)

		req, _ := http.NewRequest("POST", fmt.Sprintf("/api/superadmin/tenants/%d/quota/reset", tenant.ID), nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// 验证配额已重置
		var resetConfig models.TenantConfig
		config.DB.Where("tenant_id = ?", tenant.ID).First(&resetConfig)

		var usage map[string]int
		json.Unmarshal([]byte(resetConfig.ResourceUsage), &usage)

		if usage["used_images"] != 0 {
			t.Errorf("Expected used_images 0 after reset, got %d", usage["used_images"])
		}

		if usage["used_storage_mb"] != 0 {
			t.Errorf("Expected used_storage_mb 0 after reset, got %d", usage["used_storage_mb"])
		}
	})
}

// TestTenantQuotaAPI 测试租户配额 API
func TestTenantQuotaAPI(t *testing.T) {
	gin.SetMode(gin.TestMode)

	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	config.DB = setupIntegrationTestDB()

	// 创建租户
	tenant := createTestTenant(config.DB, "quota-api-tenant")

	// 创建租户管理员
	tenantAdmin := createTestUser(config.DB, "tenant@quotaapi", "tenant_admin", tenant.ID)

	// 创建配置
	tenantConfig := models.TenantConfig{
		TenantID:         tenant.ID,
		FeatureFlags:     `{"image_management":true,"page_config":true}`,
		ResourceQuota:    `{"max_images":100,"max_storage_mb":1024}`,
		ResourceUsage:    `{"used_images":25,"used_storage_mb":256}`,
		SubscriptionPlan: "free",
	}
	config.DB.Create(&tenantConfig)

	// 创建路由
	r := gin.Default()
	routes.RegisterRoutes(r)

	token := generateTestToken(tenantAdmin.ID, "tenant_admin", tenant.ID)

	// 测试：获取配额使用情况
	t.Run("GetQuota", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/tenant/quota", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data, ok := response["data"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected data object in response")
		}

		quota, ok := data["quota"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected quota object")
		}

		if int(quota["max_images"].(float64)) != 100 {
			t.Errorf("Expected max_images 100, got %v", quota["max_images"])
		}

		usage, ok := data["usage"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected usage object")
		}

		if int(usage["used_images"].(float64)) != 25 {
			t.Errorf("Expected used_images 25, got %v", usage["used_images"])
		}
	})

	// 测试：获取功能模块状态
	t.Run("GetFeatures", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/tenant/features", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data, ok := response["data"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected data object in response")
		}

		features, ok := data["features"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected features object")
		}

		if !features["image_management"].(bool) {
			t.Error("Expected image_management to be true")
		}
	})
}

// TestTenantIsolation 测试租户隔离（租户只能访问自己的配置）
func TestTenantIsolation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	config.DB = setupIntegrationTestDB()

	// 创建两个租户
	tenant1 := createTestTenant(config.DB, "isolation-tenant-1")
	tenant2 := createTestTenant(config.DB, "isolation-tenant-2")

	// 创建两个租户的管理员
	tenant1Admin := createTestUser(config.DB, "tenant1@isolation", "tenant_admin", tenant1.ID)
	createTestUser(config.DB, "tenant2@isolation", "tenant_admin", tenant2.ID)

	// 为两个租户创建不同的配置
	config.DB.Create(&models.TenantConfig{
		TenantID:         tenant1.ID,
		FeatureFlags:     `{"image_management":true}`,
		ResourceQuota:    `{"max_images":100}`,
		ResourceUsage:    `{"used_images":10}`,
		SubscriptionPlan: "free",
	})

	config.DB.Create(&models.TenantConfig{
		TenantID:         tenant2.ID,
		FeatureFlags:     `{"image_management":false}`,
		ResourceQuota:    `{"max_images":50}`,
		ResourceUsage:    `{"used_images":5}`,
		SubscriptionPlan: "pro",
	})

	// 创建路由
	r := gin.Default()
	routes.RegisterRoutes(r)

	// 测试：租户 1 管理员尝试访问租户 2 的配置（应该失败）
	t.Run("Tenant1_AccessTenant2Config", func(t *testing.T) {
		token := generateTestToken(tenant1Admin.ID, "tenant_admin", tenant1.ID)

		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/superadmin/tenants/%d/config", tenant2.ID), nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// 应该返回 403 权限不足
		if w.Code != 403 {
			t.Errorf("Expected status 403 for cross-tenant access, got %d", w.Code)
		}
	})

	// 测试：租户 1 管理员只能访问自己的配置
	t.Run("Tenant1_AccessOwnConfig", func(t *testing.T) {
		token := generateTestToken(tenant1Admin.ID, "tenant_admin", tenant1.ID)

		req, _ := http.NewRequest("GET", "/api/tenant/config", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200 for own config, got %d", w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data, ok := response["data"].(map[string]interface{})
		if !ok {
			t.Fatal("Expected data object in response")
		}

		// 验证返回的是租户 1 的配置
		if data["subscription_plan"] != "free" {
			t.Errorf("Expected subscription_plan 'free' for tenant1, got %v", data["subscription_plan"])
		}
	})
}