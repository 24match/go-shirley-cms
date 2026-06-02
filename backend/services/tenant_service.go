package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"medical-device-cms/backend/config"
	"medical-device-cms/backend/models"
)

// TenantService 租户服务
type TenantService struct{}

// NewTenantService 创建租户服务实例
func NewTenantService() *TenantService {
	return &TenantService{}
}

// CreateTenantRequest 创建租户请求
type CreateTenantRequest struct {
	Name              string     `json:"name" binding:"required"`
	Status            string     `json:"status"`
	SubDomain         string     `json:"sub_domain"`
	CustomDomain      string     `json:"custom_domain"`
	SubscriptionPlan  string     `json:"subscription_plan"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at"`
	AdminUsername     string     `json:"admin_username" binding:"required"`
	AdminPassword     string     `json:"admin_password" binding:"required"`
}

// UpdateTenantRequest 更新租户请求
type UpdateTenantRequest struct {
	Name              string     `json:"name"`
	Status            string     `json:"status"`
	SubDomain         string     `json:"sub_domain"`
	CustomDomain      string     `json:"custom_domain"`
	SubscriptionPlan  string     `json:"subscription_plan"`
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at"`
}

// TenantFilter 租户查询过滤器
type TenantFilter struct {
	Status string `form:"status"`
	Page   int    `form:"page" binding:"min=1"`
	Size   int    `form:"size" binding:"min=1,max=100"`
}

// CreateTenant 创建租户
func (s *TenantService) CreateTenant(ctx context.Context, req *CreateTenantRequest) (*models.Tenant, error) {
	// 生成唯一 tenant_code
	tenantCode := s.generateTenantCode(req.Name)

	// 检查 tenant_code 是否已存在
	var existing models.Tenant
	if err := config.DB.Where("tenant_code = ?", tenantCode).First(&existing).Error; err == nil {
		// 如果已存在，添加时间戳
		tenantCode = fmt.Sprintf("%s-%d", tenantCode, time.Now().UnixNano()%10000)
	}

	// 设置默认值
	if req.Status == "" {
		req.Status = models.TenantStatusActive
	}
	if req.SubscriptionPlan == "" {
		req.SubscriptionPlan = models.SubscriptionPlanFree
	}

	// 创建租户
	tenant := models.Tenant{
		TenantCode:          tenantCode,
		Name:                req.Name,
		Status:              req.Status,
		SubDomain:           req.SubDomain,
		CustomDomain:        req.CustomDomain,
		SubscriptionPlan:    req.SubscriptionPlan,
		SubscriptionExpiresAt: req.SubscriptionExpiresAt,
	}

	if err := config.DB.Create(&tenant).Error; err != nil {
		return nil, err
	}

	// 创建租户管理员账户
	if req.AdminUsername != "" && req.AdminPassword != "" {
		adminUser := models.User{
			TenantID: tenant.ID,
			Username: req.AdminUsername,
			Role:     "tenant_admin",
		}
		// 密码加密在 user_service 中处理
		if err := config.DB.Create(&adminUser).Error; err != nil {
			// 回滚租户创建
			config.DB.Delete(&tenant)
			return nil, fmt.Errorf("failed to create admin user: %w", err)
		}
	}

	return &tenant, nil
}

// GetTenant 获取租户信息
func (s *TenantService) GetTenant(ctx context.Context, id uint) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := config.DB.First(&tenant, id).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

// GetTenantByCode 根据 tenant_code 获取租户
func (s *TenantService) GetTenantByCode(ctx context.Context, tenantCode string) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := config.DB.Where("tenant_code = ?", tenantCode).First(&tenant).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

// ListTenants 获取租户列表
func (s *TenantService) ListTenants(ctx context.Context, filter *TenantFilter) ([]*models.Tenant, int64, error) {
	var tenants []*models.Tenant
	var total int64

	query := config.DB.Model(&models.Tenant{})

	// 应用过滤条件
	if filter != nil && filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if filter != nil && filter.Page > 0 && filter.Size > 0 {
		offset := (filter.Page - 1) * filter.Size
		query = query.Offset(offset).Limit(filter.Size)
	}

	// 按创建时间倒序
	query = query.Order("created_at DESC")

	if err := query.Find(&tenants).Error; err != nil {
		return nil, 0, err
	}

	return tenants, total, nil
}

// UpdateTenant 更新租户信息
func (s *TenantService) UpdateTenant(ctx context.Context, id uint, req *UpdateTenantRequest) (*models.Tenant, error) {
	var tenant models.Tenant
	if err := config.DB.First(&tenant, id).Error; err != nil {
		return nil, err
	}

	// 更新字段
	if req.Name != "" {
		tenant.Name = req.Name
	}
	if req.Status != "" {
		tenant.Status = req.Status
	}
	if req.SubDomain != "" {
		tenant.SubDomain = req.SubDomain
	}
	if req.CustomDomain != "" {
		tenant.CustomDomain = req.CustomDomain
	}
	if req.SubscriptionPlan != "" {
		tenant.SubscriptionPlan = req.SubscriptionPlan
	}
	if req.SubscriptionExpiresAt != nil {
		tenant.SubscriptionExpiresAt = req.SubscriptionExpiresAt
	}

	if err := config.DB.Save(&tenant).Error; err != nil {
		return nil, err
	}

	return &tenant, nil
}

// DeleteTenant 删除租户（软删除）
func (s *TenantService) DeleteTenant(ctx context.Context, id uint) error {
	var tenant models.Tenant
	if err := config.DB.First(&tenant, id).Error; err != nil {
		return err
	}

	// 使用 GORM 的软删除
	if err := config.DB.Delete(&tenant).Error; err != nil {
		return err
	}

	return nil
}

// ActivateTenant 激活租户
func (s *TenantService) ActivateTenant(ctx context.Context, id uint) error {
	return config.DB.Model(&models.Tenant{}).Where("id = ?", id).Update("status", models.TenantStatusActive).Error
}

// DisableTenant 禁用租户
func (s *TenantService) DisableTenant(ctx context.Context, id uint) error {
	return config.DB.Model(&models.Tenant{}).Where("id = ?", id).Update("status", models.TenantStatusDisabled).Error
}

// SuspendTenant 暂停租户
func (s *TenantService) SuspendTenant(ctx context.Context, id uint) error {
	return config.DB.Model(&models.Tenant{}).Where("id = ?", id).Update("status", models.TenantStatusSuspended).Error
}

// generateTenantCode 生成唯一 tenant_code
func (s *TenantService) generateTenantCode(name string) string {
	// 将名称转换为小写
	code := strings.ToLower(name)
	// 替换空格和特殊字符为连字符
	code = strings.ReplaceAll(code, " ", "-")
	code = strings.ReplaceAll(code, "_", "-")
	// 移除非法字符
	code = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, code)
	// 移除连续的连字符
	for strings.Contains(code, "--") {
		code = strings.ReplaceAll(code, "--", "-")
	}
	// 移除开头和结尾的连字符
	code = strings.Trim(code, "-")
	// 如果为空，使用默认值
	if code == "" {
		code = "tenant"
	}
	return code
}