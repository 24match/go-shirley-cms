package services

import (
	"context"
	"encoding/json"
	"testing"

	"medical-device-cms/backend/config"
	"medical-device-cms/backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestGetTenantConfig 测试获取租户配置
func TestGetTenantConfig(t *testing.T) {
	// 保存原有 DB 引用
	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	// 设置测试数据库
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	config.DB = db

	// 迁移表结构
	config.DB.AutoMigrate(&models.Tenant{}, &models.TenantConfig{})

	// 创建测试租户
	tenant := models.Tenant{
		TenantCode: "test-tenant",
		Name:       "Test Tenant",
		Status:     models.TenantStatusActive,
	}
	config.DB.Create(&tenant)

	// 创建测试配置
	featureFlags := map[string]bool{
		"image_management":   true,
		"page_config":        true,
		"multi_language":     false,
		"contact_form":       true,
		"content_management": true,
	}
	resourceQuota := map[string]int{
		"max_images":        100,
		"max_storage_mb":    1024,
		"max_content_items": 50,
		"max_users":         5,
	}
	resourceUsage := map[string]int{
		"used_images":        25,
		"used_storage_mb":    256,
		"used_content_items": 10,
		"used_users":         2,
	}

	flagsJSON, _ := json.Marshal(featureFlags)
	quotaJSON, _ := json.Marshal(resourceQuota)
	usageJSON, _ := json.Marshal(resourceUsage)

	tenantConfig := models.TenantConfig{
		TenantID:              tenant.ID,
		FeatureFlags:          string(flagsJSON),
		ResourceQuota:         string(quotaJSON),
		ResourceUsage:         string(usageJSON),
		SubscriptionPlan:      "pro",
		SubscriptionExpiresAt: nil,
	}
	config.DB.Create(&tenantConfig)

	// 测试：正常获取配置
	t.Run("GetTenantConfig_Success", func(t *testing.T) {
		service := NewTenantConfigService()
		dto, err := service.GetTenantConfig(context.Background(), tenant.ID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if dto == nil {
			t.Fatal("Expected config DTO, got nil")
		}

		if dto.TenantID != tenant.ID {
			t.Errorf("Expected TenantID %d, got %d", tenant.ID, dto.TenantID)
		}

		if dto.SubscriptionPlan != "pro" {
			t.Errorf("Expected SubscriptionPlan 'pro', got '%s'", dto.SubscriptionPlan)
		}

		// 验证功能标志
		if !dto.FeatureFlags["image_management"] {
			t.Error("Expected image_management to be true")
		}

		if dto.FeatureFlags["multi_language"] {
			t.Error("Expected multi_language to be false")
		}

		// 验证配额
		if dto.ResourceQuota["max_images"] != 100 {
			t.Errorf("Expected max_images 100, got %d", dto.ResourceQuota["max_images"])
		}

		// 验证使用量
		if dto.ResourceUsage["used_images"] != 25 {
			t.Errorf("Expected used_images 25, got %d", dto.ResourceUsage["used_images"])
		}
	})

	// 测试：租户配置不存在
	t.Run("GetTenantConfig_NotFound", func(t *testing.T) {
		service := NewTenantConfigService()
		_, err := service.GetTenantConfig(context.Background(), 99999)

		if err == nil {
			t.Error("Expected error for non-existent config, got nil")
		}
	})
}

// TestUpdateTenantConfig 测试更新租户配置
func TestUpdateTenantConfig(t *testing.T) {
	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	config.DB = db
	config.DB.AutoMigrate(&models.Tenant{}, &models.TenantConfig{})

	// 创建测试租户
	tenant := models.Tenant{
		TenantCode: "test-tenant-update",
		Name:       "Test Tenant Update",
		Status:     models.TenantStatusActive,
	}
	config.DB.Create(&tenant)

	// 创建初始配置
	initialFlags := map[string]bool{
		"image_management": true,
		"page_config":      true,
	}
	initialQuota := map[string]int{
		"max_images":     50,
		"max_storage_mb": 512,
	}

	flagsJSON, _ := json.Marshal(initialFlags)
	quotaJSON, _ := json.Marshal(initialQuota)

	tenantConfig := models.TenantConfig{
		TenantID:         tenant.ID,
		FeatureFlags:     string(flagsJSON),
		ResourceQuota:    string(quotaJSON),
		ResourceUsage:    `{"used_images":0,"used_storage_mb":0}`,
		SubscriptionPlan: "free",
	}
	config.DB.Create(&tenantConfig)

	t.Run("UpdateTenantConfig_Success", func(t *testing.T) {
		service := NewTenantConfigService()

		newFlags := map[string]bool{
			"image_management":   true,
			"page_config":        true,
			"multi_language":     true,
			"contact_form":       true,
			"content_management": true,
		}

		newQuota := map[string]int{
			"max_images":        500,
			"max_storage_mb":    5120,
			"max_content_items": 100,
			"max_users":         10,
		}

		req := &UpdateTenantConfigRequest{
			FeatureFlags:     newFlags,
			ResourceQuota:    newQuota,
			SubscriptionPlan: "pro",
		}

		dto, err := service.UpdateTenantConfig(context.Background(), tenant.ID, req)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if dto == nil {
			t.Fatal("Expected config DTO, got nil")
		}

		if dto.SubscriptionPlan != "pro" {
			t.Errorf("Expected SubscriptionPlan 'pro', got '%s'", dto.SubscriptionPlan)
		}

		if !dto.FeatureFlags["multi_language"] {
			t.Error("Expected multi_language to be true after update")
		}

		if dto.ResourceQuota["max_images"] != 500 {
			t.Errorf("Expected max_images 500, got %d", dto.ResourceQuota["max_images"])
		}
	})

	t.Run("UpdateTenantConfig_NotFound", func(t *testing.T) {
		service := NewTenantConfigService()

		req := &UpdateTenantConfigRequest{
			SubscriptionPlan: "enterprise",
		}

		_, err := service.UpdateTenantConfig(context.Background(), 99999, req)

		if err == nil {
			t.Error("Expected error for non-existent config, got nil")
		}
	})
}

// TestCheckFeature 测试功能开关检查
func TestCheckFeature(t *testing.T) {
	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	config.DB = db
	config.DB.AutoMigrate(&models.Tenant{}, &models.TenantConfig{})

	tenant := models.Tenant{
		TenantCode: "test-tenant-feature",
		Name:       "Test Tenant Feature",
		Status:     models.TenantStatusActive,
	}
	config.DB.Create(&tenant)

	featureFlags := map[string]bool{
		"image_management":   true,
		"page_config":        true,
		"multi_language":     false,
		"contact_form":       true,
		"content_management": true,
	}

	flagsJSON, _ := json.Marshal(featureFlags)

	tenantConfig := models.TenantConfig{
		TenantID:     tenant.ID,
		FeatureFlags: string(flagsJSON),
		ResourceQuota: `{"max_images":50,"max_storage_mb":512}`,
		ResourceUsage: `{"used_images":0,"used_storage_mb":0}`,
		SubscriptionPlan: "free",
	}
	config.DB.Create(&tenantConfig)

	service := NewTenantConfigService()

	t.Run("CheckFeature_Enabled", func(t *testing.T) {
		enabled, err := service.CheckFeature(context.Background(), tenant.ID, "image_management")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if !enabled {
			t.Error("Expected image_management to be enabled")
		}
	})

	t.Run("CheckFeature_Disabled", func(t *testing.T) {
		enabled, err := service.CheckFeature(context.Background(), tenant.ID, "multi_language")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if enabled {
			t.Error("Expected multi_language to be disabled")
		}
	})

	t.Run("CheckFeature_NotExists", func(t *testing.T) {
		enabled, err := service.CheckFeature(context.Background(), tenant.ID, "unknown_feature")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if enabled {
			t.Error("Expected unknown_feature to return false")
		}
	})

	t.Run("CheckFeature_TenantNotFound", func(t *testing.T) {
		_, err := service.CheckFeature(context.Background(), 99999, "image_management")

		if err == nil {
			t.Error("Expected error for non-existent tenant, got nil")
		}
	})
}

// TestCheckQuota 测试配额检查
func TestCheckQuota(t *testing.T) {
	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	config.DB = db
	config.DB.AutoMigrate(&models.Tenant{}, &models.TenantConfig{})

	tenant := models.Tenant{
		TenantCode: "test-tenant-quota",
		Name:       "Test Tenant Quota",
		Status:     models.TenantStatusActive,
	}
	config.DB.Create(&tenant)

	resourceQuota := map[string]int{
		"max_images":        100,
		"max_storage_mb":    1024,
		"max_content_items": 50,
		"max_users":         5,
	}

	resourceUsage := map[string]int{
		"used_images":        80,
		"used_storage_mb":    512,
		"used_content_items": 25,
		"used_users":         2,
	}

	quotaJSON, _ := json.Marshal(resourceQuota)
	usageJSON, _ := json.Marshal(resourceUsage)

	tenantConfig := models.TenantConfig{
		TenantID:         tenant.ID,
		FeatureFlags:     `{"image_management":true}`,
		ResourceQuota:    string(quotaJSON),
		ResourceUsage:    string(usageJSON),
		SubscriptionPlan: "pro",
	}
	config.DB.Create(&tenantConfig)

	service := NewTenantConfigService()

	t.Run("CheckQuota_Sufficient", func(t *testing.T) {
		sufficient, used, limit, err := service.CheckQuota(context.Background(), tenant.ID, "images")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if !sufficient {
			t.Error("Expected quota to be sufficient")
		}

		if used != 80 {
			t.Errorf("Expected used 80, got %d", used)
		}

		if limit != 100 {
			t.Errorf("Expected limit 100, got %d", limit)
		}
	})

	t.Run("CheckQuota_Exceeded", func(t *testing.T) {
		// 更新使用量到上限
		newUsage := map[string]int{
			"used_images": 100,
		}
		usageJSON, _ := json.Marshal(newUsage)
		config.DB.Model(&models.TenantConfig{}).Where("tenant_id = ?", tenant.ID).Update("resource_usage", string(usageJSON))

		sufficient, used, limit, err := service.CheckQuota(context.Background(), tenant.ID, "images")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if sufficient {
			t.Error("Expected quota to be exceeded")
		}

		if used != 100 {
			t.Errorf("Expected used 100, got %d", used)
		}

		if limit != 100 {
			t.Errorf("Expected limit 100, got %d", limit)
		}
	})

	t.Run("CheckQuota_Unlimited", func(t *testing.T) {
		// 设置无限制配额
		newQuota := map[string]int{
			"max_images": -1,
		}
		quotaJSON, _ := json.Marshal(newQuota)
		config.DB.Model(&models.TenantConfig{}).Where("tenant_id = ?", tenant.ID).Update("resource_quota", string(quotaJSON))

		sufficient, _, limit, err := service.CheckQuota(context.Background(), tenant.ID, "images")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if !sufficient {
			t.Error("Expected unlimited quota to be sufficient")
		}

		if limit != -1 {
			t.Errorf("Expected limit -1, got %d", limit)
		}
	})

	t.Run("CheckQuota_NoRestriction", func(t *testing.T) {
		// 测试不存在的配额类型
		sufficient, _, _, err := service.CheckQuota(context.Background(), tenant.ID, "unknown_resource")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if !sufficient {
			t.Error("Expected no restriction to be sufficient")
		}
	})
}

// TestIncrementUsage 测试增加使用量
func TestIncrementUsage(t *testing.T) {
	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	config.DB = db
	config.DB.AutoMigrate(&models.Tenant{}, &models.TenantConfig{})

	tenant := models.Tenant{
		TenantCode: "test-tenant-increment",
		Name:       "Test Tenant Increment",
		Status:     models.TenantStatusActive,
	}
	config.DB.Create(&tenant)

	tenantConfig := models.TenantConfig{
		TenantID:         tenant.ID,
		FeatureFlags:     `{"image_management":true}`,
		ResourceQuota:    `{"max_images":100,"max_storage_mb":1024}`,
		ResourceUsage:    `{"used_images":10,"used_storage_mb":100}`,
		SubscriptionPlan: "free",
	}
	config.DB.Create(&tenantConfig)

	service := NewTenantConfigService()

	t.Run("IncrementUsage_Success", func(t *testing.T) {
		err := service.IncrementUsage(context.Background(), tenant.ID, "images", 5)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		// 验证使用量已增加
		var tc models.TenantConfig
		config.DB.Where("tenant_id = ?", tenant.ID).First(&tc)

		var usage map[string]int
		json.Unmarshal([]byte(tc.ResourceUsage), &usage)

		if usage["used_images"] != 15 {
			t.Errorf("Expected used_images 15, got %d", usage["used_images"])
		}
	})

	t.Run("IncrementUsage_NewResource", func(t *testing.T) {
		err := service.IncrementUsage(context.Background(), tenant.ID, "content_items", 3)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		var tc models.TenantConfig
		config.DB.Where("tenant_id = ?", tenant.ID).First(&tc)

		var usage map[string]int
		json.Unmarshal([]byte(tc.ResourceUsage), &usage)

		if usage["used_content_items"] != 3 {
			t.Errorf("Expected used_content_items 3, got %d", usage["used_content_items"])
		}
	})
}

// TestResetQuotaUsage 测试重置配额使用统计
func TestResetQuotaUsage(t *testing.T) {
	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	config.DB = db
	config.DB.AutoMigrate(&models.Tenant{}, &models.TenantConfig{})

	tenant := models.Tenant{
		TenantCode: "test-tenant-reset",
		Name:       "Test Tenant Reset",
		Status:     models.TenantStatusActive,
	}
	config.DB.Create(&tenant)

	tenantConfig := models.TenantConfig{
		TenantID:         tenant.ID,
		FeatureFlags:     `{"image_management":true}`,
		ResourceQuota:    `{"max_images":100,"max_storage_mb":1024}`,
		ResourceUsage:    `{"used_images":50,"used_storage_mb":512}`,
		SubscriptionPlan: "free",
	}
	config.DB.Create(&tenantConfig)

	service := NewTenantConfigService()

	t.Run("ResetQuotaUsage_SingleResource", func(t *testing.T) {
		err := service.ResetQuotaUsage(context.Background(), tenant.ID, "images")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		var tc models.TenantConfig
		config.DB.Where("tenant_id = ?", tenant.ID).First(&tc)

		var usage map[string]int
		json.Unmarshal([]byte(tc.ResourceUsage), &usage)

		if usage["used_images"] != 0 {
			t.Errorf("Expected used_images 0, got %d", usage["used_images"])
		}

		if usage["used_storage_mb"] != 512 {
			t.Errorf("Expected used_storage_mb to remain 512, got %d", usage["used_storage_mb"])
		}
	})

	t.Run("ResetQuotaUsage_All", func(t *testing.T) {
		// 先恢复一些使用量
		config.DB.Model(&models.TenantConfig{}).Where("tenant_id = ?", tenant.ID).Update("resource_usage", `{"used_images":50,"used_storage_mb":512}`)

		err := service.ResetQuotaUsage(context.Background(), tenant.ID, "")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		var tc models.TenantConfig
		config.DB.Where("tenant_id = ?", tenant.ID).First(&tc)

		var usage map[string]int
		json.Unmarshal([]byte(tc.ResourceUsage), &usage)

		if usage["used_images"] != 0 {
			t.Errorf("Expected used_images 0, got %d", usage["used_images"])
		}

		if usage["used_storage_mb"] != 0 {
			t.Errorf("Expected used_storage_mb 0, got %d", usage["used_storage_mb"])
		}
	})
}

// TestGetFeatures 测试获取功能模块状态
func TestGetFeatures(t *testing.T) {
	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	config.DB = db
	config.DB.AutoMigrate(&models.Tenant{}, &models.TenantConfig{})

	tenant := models.Tenant{
		TenantCode: "test-tenant-features",
		Name:       "Test Tenant Features",
		Status:     models.TenantStatusActive,
	}
	config.DB.Create(&tenant)

	featureFlags := map[string]bool{
		"image_management":   true,
		"page_config":        false,
		"multi_language":     true,
		"contact_form":       true,
		"content_management": false,
	}

	flagsJSON, _ := json.Marshal(featureFlags)

	tenantConfig := models.TenantConfig{
		TenantID:         tenant.ID,
		FeatureFlags:     string(flagsJSON),
		ResourceQuota:    `{"max_images":50}`,
		ResourceUsage:    `{"used_images":0}`,
		SubscriptionPlan: "free",
	}
	config.DB.Create(&tenantConfig)

	service := NewTenantConfigService()

	t.Run("GetFeatures_Success", func(t *testing.T) {
		features, err := service.GetFeatures(context.Background(), tenant.ID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(features) != 5 {
			t.Errorf("Expected 5 features, got %d", len(features))
		}

		if !features["image_management"] {
			t.Error("Expected image_management to be true")
		}

		if features["page_config"] {
			t.Error("Expected page_config to be false")
		}
	})

	t.Run("GetFeatures_EmptyFlags", func(t *testing.T) {
		// 创建空配置的租户
		tenant2 := models.Tenant{
			TenantCode: "test-tenant-empty",
			Name:       "Test Tenant Empty",
			Status:     models.TenantStatusActive,
		}
		config.DB.Create(&tenant2)

		config.DB.Create(&models.TenantConfig{
			TenantID:         tenant2.ID,
			FeatureFlags:     "",
			ResourceQuota:    `{"max_images":50}`,
			ResourceUsage:    `{"used_images":0}`,
			SubscriptionPlan: "free",
		})

		features, err := service.GetFeatures(context.Background(), tenant2.ID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if features == nil {
			t.Error("Expected empty map, got nil")
		}

		if len(features) != 0 {
			t.Errorf("Expected 0 features, got %d", len(features))
		}
	})
}

// TestCreateDefaultConfig 测试创建默认配置
func TestCreateDefaultConfig(t *testing.T) {
	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	config.DB = db
	config.DB.AutoMigrate(&models.Tenant{}, &models.TenantConfig{})

	tenant := models.Tenant{
		TenantCode: "test-tenant-default",
		Name:       "Test Tenant Default",
		Status:     models.TenantStatusActive,
	}
	config.DB.Create(&tenant)

	service := NewTenantConfigService()

	t.Run("CreateDefaultConfig_Free", func(t *testing.T) {
		err := service.CreateDefaultConfig(context.Background(), tenant.ID, "free")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		var tc models.TenantConfig
		config.DB.Where("tenant_id = ?", tenant.ID).First(&tc)

		if tc.ID == 0 {
			t.Fatal("Expected config to be created")
		}

		if tc.SubscriptionPlan != "free" {
			t.Errorf("Expected SubscriptionPlan 'free', got '%s'", tc.SubscriptionPlan)
		}

		var flags map[string]bool
		json.Unmarshal([]byte(tc.FeatureFlags), &flags)

		if !flags["image_management"] {
			t.Error("Expected image_management to be true by default")
		}

		var quota map[string]int
		json.Unmarshal([]byte(tc.ResourceQuota), &quota)

		if quota["max_images"] != 50 {
			t.Errorf("Expected max_images 50 for free plan, got %d", quota["max_images"])
		}
	})

	t.Run("CreateDefaultConfig_Pro", func(t *testing.T) {
		tenant2 := models.Tenant{
			TenantCode: "test-tenant-pro",
			Name:       "Test Tenant Pro",
			Status:     models.TenantStatusActive,
		}
		config.DB.Create(&tenant2)

		err := service.CreateDefaultConfig(context.Background(), tenant2.ID, "pro")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		var tc models.TenantConfig
		config.DB.Where("tenant_id = ?", tenant2.ID).First(&tc)

		var quota map[string]int
		json.Unmarshal([]byte(tc.ResourceQuota), &quota)

		if quota["max_images"] != 500 {
			t.Errorf("Expected max_images 500 for pro plan, got %d", quota["max_images"])
		}

		if quota["max_storage_mb"] != 5120 {
			t.Errorf("Expected max_storage_mb 5120 for pro plan, got %d", quota["max_storage_mb"])
		}
	})

	t.Run("CreateDefaultConfig_Enterprise", func(t *testing.T) {
		tenant3 := models.Tenant{
			TenantCode: "test-tenant-enterprise",
			Name:       "Test Tenant Enterprise",
			Status:     models.TenantStatusActive,
		}
		config.DB.Create(&tenant3)

		err := service.CreateDefaultConfig(context.Background(), tenant3.ID, "enterprise")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		var tc models.TenantConfig
		config.DB.Where("tenant_id = ?", tenant3.ID).First(&tc)

		var quota map[string]int
		json.Unmarshal([]byte(tc.ResourceQuota), &quota)

		if quota["max_images"] != -1 {
			t.Errorf("Expected max_images -1 (unlimited) for enterprise plan, got %d", quota["max_images"])
		}
	})

	t.Run("CreateDefaultConfig_AlreadyExists", func(t *testing.T) {
		// 再次为已存在配置的租户创建配置
		err := service.CreateDefaultConfig(context.Background(), tenant.ID, "free")

		if err != nil {
			t.Errorf("Expected no error for existing config, got %v", err)
		}

		var count int64
		config.DB.Model(&models.TenantConfig{}).Where("tenant_id = ?", tenant.ID).Count(&count)

		if count != 1 {
			t.Errorf("Expected 1 config record, got %d", count)
		}
	})
}