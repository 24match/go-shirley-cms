package services

import (
	"testing"
	"time"

	"medical-device-cms/backend/models"
)

// TestCreateTenant 测试租户创建功能
func TestCreateTenant(t *testing.T) {
	service := NewTenantService()

	// 测试创建新租户
	tenant, err := service.CreateTenant("test-tenant", "Test Tenant", "test", "")
	if err != nil {
		t.Fatalf("Failed to create tenant: %v", err)
	}

	if tenant.TenantCode != "test-tenant" {
		t.Errorf("Expected TenantCode to be 'test-tenant', got '%s'", tenant.TenantCode)
	}

	if tenant.Name != "Test Tenant" {
		t.Errorf("Expected Name to be 'Test Tenant', got '%s'", tenant.Name)
	}

	if tenant.Status != models.TenantStatusActive {
		t.Errorf("Expected Status to be Active, got %d", tenant.Status)
	}
}

// TestCreateTenantDuplicateCode 测试租户代码唯一性
func TestCreateTenantDuplicateCode(t *testing.T) {
	service := NewTenantService()

	// 首先创建一个租户
	_, err := service.CreateTenant("duplicate-test", "First Tenant", "test", "")
	if err != nil {
		t.Fatalf("Failed to create first tenant: %v", err)
	}

	// 尝试创建相同代码的租户
	_, err = service.CreateTenant("duplicate-test", "Second Tenant", "test", "")
	if err == nil {
		t.Error("Expected error when creating tenant with duplicate code, got nil")
	}
}

// TestGetTenant 测试租户查询功能
func TestGetTenant(t *testing.T) {
	service := NewTenantService()

	// 创建一个测试租户
	tenant, err := service.CreateTenant("get-test", "Get Test Tenant", "test", "")
	if err != nil {
		t.Fatalf("Failed to create tenant: %v", err)
	}

	// 通过 ID 查询
	foundTenant, err := service.GetTenantByID(tenant.ID)
	if err != nil {
		t.Fatalf("Failed to get tenant by ID: %v", err)
	}

	if foundTenant.ID != tenant.ID {
		t.Errorf("Expected ID %d, got %d", tenant.ID, foundTenant.ID)
	}

	// 通过 TenantCode 查询
	foundTenant2, err := service.GetTenantByCode("get-test")
	if err != nil {
		t.Fatalf("Failed to get tenant by code: %v", err)
	}

	if foundTenant2.ID != tenant.ID {
		t.Errorf("Expected ID %d, got %d", tenant.ID, foundTenant2.ID)
	}
}

// TestListTenants 测试租户列表功能
func TestListTenants(t *testing.T) {
	service := NewTenantService()

	// 创建多个测试租户
	service.CreateTenant("list-test-1", "List Test 1", "test", "")
	service.CreateTenant("list-test-2", "List Test 2", "test", "")

	// 获取租户列表
	tenants, total, err := service.ListTenants("", 0, 10)
	if err != nil {
		t.Fatalf("Failed to list tenants: %v", err)
	}

	if total < 2 {
		t.Errorf("Expected at least 2 tenants, got %d", total)
	}

	// 测试状态过滤
	activeTenants, activeTotal, err := service.ListTenants("active", 0, 10)
	if err != nil {
		t.Fatalf("Failed to list active tenants: %v", err)
	}

	if activeTotal < 2 {
		t.Errorf("Expected at least 2 active tenants, got %d", activeTotal)
	}

	_ = activeTenants
}

// TestUpdateTenant 测试租户更新功能
func TestUpdateTenant(t *testing.T) {
	service := NewTenantService()

	// 创建测试租户
	tenant, err := service.CreateTenant("update-test", "Update Test", "test", "")
	if err != nil {
		t.Fatalf("Failed to create tenant: %v", err)
	}

	// 更新租户信息
	updates := map[string]interface{}{
		"name":        "Updated Test Tenant",
		"sub_domain":  "updated",
		"description": "Test description",
	}

	updatedTenant, err := service.UpdateTenant(tenant.ID, updates)
	if err != nil {
		t.Fatalf("Failed to update tenant: %v", err)
	}

	if updatedTenant.Name != "Updated Test Tenant" {
		t.Errorf("Expected Name to be 'Updated Test Tenant', got '%s'", updatedTenant.Name)
	}

	if updatedTenant.SubDomain != "updated" {
		t.Errorf("Expected SubDomain to be 'updated', got '%s'", updatedTenant.SubDomain)
	}
}

// TestActivateTenant 测试租户激活功能
func TestActivateTenant(t *testing.T) {
	service := NewTenantService()

	// 创建测试租户
	tenant, err := service.CreateTenant("activate-test", "Activate Test", "test", "")
	if err != nil {
		t.Fatalf("Failed to create tenant: %v", err)
	}

	// 先禁用租户
	_, err = service.DisableTenant(tenant.ID)
	if err != nil {
		t.Fatalf("Failed to disable tenant: %v", err)
	}

	// 验证租户已禁用
	disabledTenant, _ := service.GetTenantByID(tenant.ID)
	if disabledTenant.Status != models.TenantStatusDisabled {
		t.Errorf("Expected Status to be Disabled, got %d", disabledTenant.Status)
	}

	// 激活租户
	activatedTenant, err := service.ActivateTenant(tenant.ID)
	if err != nil {
		t.Fatalf("Failed to activate tenant: %v", err)
	}

	if activatedTenant.Status != models.TenantStatusActive {
		t.Errorf("Expected Status to be Active, got %d", activatedTenant.Status)
	}
}

// TestDisableTenant 测试租户禁用功能
func TestDisableTenant(t *testing.T) {
	service := NewTenantService()

	// 创建测试租户
	tenant, err := service.CreateTenant("disable-test", "Disable Test", "test", "")
	if err != nil {
		t.Fatalf("Failed to create tenant: %v", err)
	}

	// 禁用租户
	disabledTenant, err := service.DisableTenant(tenant.ID)
	if err != nil {
		t.Fatalf("Failed to disable tenant: %v", err)
	}

	if disabledTenant.Status != models.TenantStatusDisabled {
		t.Errorf("Expected Status to be Disabled, got %d", disabledTenant.Status)
	}
}

// TestDeleteTenant 测试租户删除功能
func TestDeleteTenant(t *testing.T) {
	service := NewTenantService()

	// 创建测试租户
	tenant, err := service.CreateTenant("delete-test", "Delete Test", "test", "")
	if err != nil {
		t.Fatalf("Failed to create tenant: %v", err)
	}

	// 删除租户
	err = service.DeleteTenant(tenant.ID)
	if err != nil {
		t.Fatalf("Failed to delete tenant: %v", err)
	}

	// 验证租户已删除（软删除，状态变为 deleted）
	deletedTenant, err := service.GetTenantByID(tenant.ID)
	if err != nil {
		t.Fatalf("Failed to get deleted tenant: %v", err)
	}

	if deletedTenant.Status != models.TenantStatusDeleted {
		t.Errorf("Expected Status to be Deleted, got %d", deletedTenant.Status)
	}
}

// TestGetSystemStats 测试系统统计功能（超级管理员）
func TestGetSystemStats(t *testing.T) {
	service := NewTenantService()

	// 获取系统统计
	stats, err := service.GetSystemStats()
	if err != nil {
		t.Fatalf("Failed to get system stats: %v", err)
	}

	// 验证统计数据结构
	if stats["total_tenants"] == nil {
		t.Error("Expected total_tenants in stats")
	}

	if stats["active_tenants"] == nil {
		t.Error("Expected active_tenants in stats")
	}

	if stats["total_users"] == nil {
		t.Error("Expected total_users in stats")
	}
}

// TestGenerateTenantCode 测试租户代码生成
func TestGenerateTenantCode(t *testing.T) {
	service := NewTenantService()

	// 测试生成租户代码
	code := service.GenerateTenantCode("Test Company Name")
	
	if code == "" {
		t.Error("Expected non-empty tenant code")
	}

	// 验证代码格式（应该只包含小写字母和数字）
	for _, r := range code {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-') {
			t.Errorf("Invalid character '%c' in tenant code '%s'", r, code)
		}
	}
}

// TestImpersonateTenant 测试租户切换功能
func TestImpersonateTenant(t *testing.T) {
	service := NewTenantService()

	// 创建测试租户
	tenant, err := service.CreateTenant("impersonate-test", "Impersonate Test", "test", "")
	if err != nil {
		t.Fatalf("Failed to create tenant: %v", err)
	}

	// 测试切换租户（这里主要验证租户存在性检查）
	foundTenant, err := service.GetTenantByID(tenant.ID)
	if err != nil {
		t.Fatalf("Failed to get tenant for impersonation: %v", err)
	}

	if foundTenant.ID != tenant.ID {
		t.Errorf("Expected tenant ID %d, got %d", tenant.ID, foundTenant.ID)
	}
}

// TestTenantSubscription 测试租户订阅功能
func TestTenantSubscription(t *testing.T) {
	service := NewTenantService()

	// 创建带订阅信息的租户
	expiresAt := time.Now().AddDate(0, 1, 0) // 1 个月后过期
	tenant, err := service.CreateTenantWithSubscription(
		"subscription-test",
		"Subscription Test",
		"test",
		"",
		"premium",
		&expiresAt,
	)
	if err != nil {
		t.Fatalf("Failed to create tenant with subscription: %v", err)
	}

	if tenant.SubscriptionPlan != "premium" {
		t.Errorf("Expected SubscriptionPlan to be 'premium', got '%s'", tenant.SubscriptionPlan)
	}

	// 验证过期时间
	if tenant.SubscriptionExpiresAt == nil {
		t.Error("Expected SubscriptionExpiresAt to be set")
	}
}