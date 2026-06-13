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

// TestQuotaService_CheckQuota 测试配额检查
func TestQuotaService_CheckQuota(t *testing.T) {
	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	config.DB = db
	config.DB.AutoMigrate(&models.Tenant{}, &models.TenantConfig{})

	tenant := models.Tenant{
		TenantCode: "test-quota-tenant",
		Name:       "Test Quota Tenant",
		Status:     models.TenantStatusActive,
	}
	config.DB.Create(&tenant)

	// 创建测试配置
	resourceQuota := map[string]int{
		"max_images":        100,
		"max_storage_mb":    1024,
		"max_content_items": 50,
		"max_users":         5,
	}

	resourceUsage := map[string]int{
		"used_images":        50,
		"used_storage_mb":    256,
		"used_content_items": 20,
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

	service := NewQuotaService()

	t.Run("CheckQuota_Sufficient", func(t *testing.T) {
		sufficient, info, err := service.CheckQuota(context.Background(), tenant.ID, "images")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if !sufficient {
			t.Error("Expected quota to be sufficient")
		}

		if info == nil {
			t.Fatal("Expected quota info, got nil")
		}

		if info.Used != 50 {
			t.Errorf("Expected Used 50, got %d", info.Used)
		}

		if info.Limit != 100 {
			t.Errorf("Expected Limit 100, got %d", info.Limit)
		}

		if info.Available != 50 {
			t.Errorf("Expected Available 50, got %d", info.Available)
		}

		if info.IsUnlimited {
			t.Error("Expected IsUnlimited to be false")
		}
	})

	t.Run("CheckQuota_Exceeded", func(t *testing.T) {
		// 更新使用量到上限
		newUsage := map[string]int{
			"used_images": 100,
		}
		usageJSON, _ := json.Marshal(newUsage)
		config.DB.Model(&models.TenantConfig{}).Where("tenant_id = ?", tenant.ID).Update("resource_usage", string(usageJSON))

		sufficient, info, err := service.CheckQuota(context.Background(), tenant.ID, "images")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if sufficient {
			t.Error("Expected quota to be exceeded")
		}

		if info.Used != 100 {
			t.Errorf("Expected Used 100, got %d", info.Used)
		}

		if info.Available != 0 {
			t.Errorf("Expected Available 0, got %d", info.Available)
		}
	})

	t.Run("CheckQuota_Unlimited", func(t *testing.T) {
		// 设置无限制配额
		newQuota := map[string]int{
			"max_images": -1,
		}
		quotaJSON, _ := json.Marshal(newQuota)
		config.DB.Model(&models.TenantConfig{}).Where("tenant_id = ?", tenant.ID).Update("resource_quota", string(quotaJSON))

		sufficient, info, err := service.CheckQuota(context.Background(), tenant.ID, "images")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if !sufficient {
			t.Error("Expected unlimited quota to be sufficient")
		}

		if !info.IsUnlimited {
			t.Error("Expected IsUnlimited to be true")
		}

		if info.Limit != -1 {
			t.Errorf("Expected Limit -1, got %d", info.Limit)
		}
	})

	t.Run("CheckQuota_NoRestriction", func(t *testing.T) {
		sufficient, info, err := service.CheckQuota(context.Background(), tenant.ID, "unknown_resource")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if !sufficient {
			t.Error("Expected no restriction to be sufficient")
		}

		if !info.IsUnlimited {
			t.Error("Expected IsUnlimited to be true for non-existent quota")
		}
	})

	t.Run("CheckQuota_TenantNotFound", func(t *testing.T) {
		_, _, err := service.CheckQuota(context.Background(), 99999, "images")

		if err == nil {
			t.Error("Expected error for non-existent tenant, got nil")
		}
	})
}

// TestQuotaService_CheckAndIncrement 测试检查并增加使用量
func TestQuotaService_CheckAndIncrement(t *testing.T) {
	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	config.DB = db
	config.DB.AutoMigrate(&models.Tenant{}, &models.TenantConfig{})

	tenant := models.Tenant{
		TenantCode: "test-increment-tenant",
		Name:       "Test Increment Tenant",
		Status:     models.TenantStatusActive,
	}
	config.DB.Create(&tenant)

	tenantConfig := models.TenantConfig{
		TenantID:         tenant.ID,
		FeatureFlags:     `{"image_management":true}`,
		ResourceQuota:    `{"max_images":100,"max_storage_mb":1024}`,
		ResourceUsage:    `{"used_images":50,"used_storage_mb":100}`,
		SubscriptionPlan: "free",
	}
	config.DB.Create(&tenantConfig)

	service := NewQuotaService()

	t.Run("CheckAndIncrement_Success", func(t *testing.T) {
		success, err := service.CheckAndIncrement(context.Background(), tenant.ID, "images", 10)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if !success {
			t.Error("Expected increment to succeed")
		}

		// 验证使用量已增加
		var tc models.TenantConfig
		config.DB.Where("tenant_id = ?", tenant.ID).First(&tc)

		var usage map[string]int
		json.Unmarshal([]byte(tc.ResourceUsage), &usage)

		if usage["used_images"] != 60 {
			t.Errorf("Expected used_images 60, got %d", usage["used_images"])
		}
	})

	t.Run("CheckAndIncrement_Exceeded", func(t *testing.T) {
		// 设置使用量接近上限
		config.DB.Model(&models.TenantConfig{}).Where("tenant_id = ?", tenant.ID).Update("resource_usage", `{"used_images":95,"used_storage_mb":100}`)

		success, err := service.CheckAndIncrement(context.Background(), tenant.ID, "images", 10)

		if err == nil {
			t.Error("Expected error for exceeded quota, got nil")
		}

		if success {
			t.Error("Expected increment to fail when quota exceeded")
		}
	})

	t.Run("CheckAndIncrement_Unlimited", func(t *testing.T) {
		// 设置无限制配额
		config.DB.Model(&models.TenantConfig{}).Where("tenant_id = ?", tenant.ID).Update("resource_quota", `{"max_images":-1,"max_storage_mb":1024}`)
		config.DB.Model(&models.TenantConfig{}).Where("tenant_id = ?", tenant.ID).Update("resource_usage", `{"used_images":1000,"used_storage_mb":100}`)

		success, err := service.CheckAndIncrement(context.Background(), tenant.ID, "images", 100)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if !success {
			t.Error("Expected increment to succeed for unlimited quota")
		}
	})

	t.Run("CheckAndIncrement_NewResource", func(t *testing.T) {
		success, err := service.CheckAndIncrement(context.Background(), tenant.ID, "content_items", 5)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if !success {
			t.Error("Expected increment to succeed for new resource")
		}

		var tc models.TenantConfig
		config.DB.Where("tenant_id = ?", tenant.ID).First(&tc)

		var usage map[string]int
		json.Unmarshal([]byte(tc.ResourceUsage), &usage)

		if usage["used_content_items"] != 5 {
			t.Errorf("Expected used_content_items 5, got %d", usage["used_content_items"])
		}
	})
}

// TestQuotaService_GetUsage 测试获取使用量
func TestQuotaService_GetUsage(t *testing.T) {
	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	config.DB = db
	config.DB.AutoMigrate(&models.Tenant{}, &models.TenantConfig{})

	tenant := models.Tenant{
		TenantCode: "test-usage-tenant",
		Name:       "Test Usage Tenant",
		Status:     models.TenantStatusActive,
	}
	config.DB.Create(&tenant)

	tenantConfig := models.TenantConfig{
		TenantID:         tenant.ID,
		FeatureFlags:     `{"image_management":true}`,
		ResourceQuota:    `{"max_images":100,"max_storage_mb":1024}`,
		ResourceUsage:    `{"used_images":50,"used_storage_mb":256}`,
		SubscriptionPlan: "free",
	}
	config.DB.Create(&tenantConfig)

	service := NewQuotaService()

	t.Run("GetUsage_Success", func(t *testing.T) {
		usage, err := service.GetUsage(context.Background(), tenant.ID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if usage == nil {
			t.Fatal("Expected usage map, got nil")
		}

		if usage["used_images"] != 50 {
			t.Errorf("Expected used_images 50, got %d", usage["used_images"])
		}

		if usage["used_storage_mb"] != 256 {
			t.Errorf("Expected used_storage_mb 256, got %d", usage["used_storage_mb"])
		}
	})

	t.Run("GetUsage_TenantNotFound", func(t *testing.T) {
		_, err := service.GetUsage(context.Background(), 99999)

		if err == nil {
			t.Error("Expected error for non-existent tenant, got nil")
		}
	})
}

// TestQuotaService_GetQuota 测试获取配额
func TestQuotaService_GetQuota(t *testing.T) {
	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	config.DB = db
	config.DB.AutoMigrate(&models.Tenant{}, &models.TenantConfig{})

	tenant := models.Tenant{
		TenantCode: "test-quota-get-tenant",
		Name:       "Test Quota Get Tenant",
		Status:     models.TenantStatusActive,
	}
	config.DB.Create(&tenant)

	tenantConfig := models.TenantConfig{
		TenantID:         tenant.ID,
		FeatureFlags:     `{"image_management":true}`,
		ResourceQuota:    `{"max_images":100,"max_storage_mb":1024}`,
		ResourceUsage:    `{"used_images":50,"used_storage_mb":256}`,
		SubscriptionPlan: "free",
	}
	config.DB.Create(&tenantConfig)

	service := NewQuotaService()

	t.Run("GetQuota_Success", func(t *testing.T) {
		quota, err := service.GetQuota(context.Background(), tenant.ID)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if quota == nil {
			t.Fatal("Expected quota map, got nil")
		}

		if quota["max_images"] != 100 {
			t.Errorf("Expected max_images 100, got %d", quota["max_images"])
		}

		if quota["max_storage_mb"] != 1024 {
			t.Errorf("Expected max_storage_mb 1024, got %d", quota["max_storage_mb"])
		}
	})

	t.Run("GetQuota_TenantNotFound", func(t *testing.T) {
		_, err := service.GetQuota(context.Background(), 99999)

		if err == nil {
			t.Error("Expected error for non-existent tenant, got nil")
		}
	})
}

// TestQuotaService_ResetQuotaUsage 测试重置配额使用量
func TestQuotaService_ResetQuotaUsage(t *testing.T) {
	originalDB := config.DB
	defer func() { config.DB = originalDB }()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	config.DB = db
	config.DB.AutoMigrate(&models.Tenant{}, &models.TenantConfig{})

	tenant := models.Tenant{
		TenantCode: "test-reset-tenant",
		Name:       "Test Reset Tenant",
		Status:     models.TenantStatusActive,
	}
	config.DB.Create(&tenant)

	tenantConfig := models.TenantConfig{
		TenantID:         tenant.ID,
		FeatureFlags:     `{"image_management":true}`,
		ResourceQuota:    `{"max_images":100,"max_storage_mb":1024}`,
		ResourceUsage:    `{"used_images":50,"used_storage_mb":256}`,
		SubscriptionPlan: "free",
	}
	config.DB.Create(&tenantConfig)

	service := NewQuotaService()

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

		if usage["used_storage_mb"] != 256 {
			t.Errorf("Expected used_storage_mb to remain 256, got %d", usage["used_storage_mb"])
		}
	})

	t.Run("ResetQuotaUsage_All", func(t *testing.T) {
		// 恢复一些使用量
		config.DB.Model(&models.TenantConfig{}).Where("tenant_id = ?", tenant.ID).Update("resource_usage", `{"used_images":50,"used_storage_mb":256}`)

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

	t.Run("ResetQuotaUsage_TenantNotFound", func(t *testing.T) {
		err := service.ResetQuotaUsage(context.Background(), 99999, "images")

		if err == nil {
			t.Error("Expected error for non-existent tenant, got nil")
		}
	})
}