package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"medical-device-cms/backend/config"
	"medical-device-cms/backend/models"
	"gorm.io/gorm"
)

// 审计操作类型常量
const (
	AuditActionTenantConfigUpdate = "TENANT_CONFIG_UPDATE"
)

// TenantConfigService 租户配置服务
type TenantConfigService struct{}

// NewTenantConfigService 创建租户配置服务实例
func NewTenantConfigService() *TenantConfigService {
	return &TenantConfigService{}
}

// TenantConfigDTO 租户配置数据传输对象
type TenantConfigDTO struct {
	ID                    uint              `json:"id"`
	TenantID              uint              `json:"tenant_id"`
	FeatureFlags          map[string]bool   `json:"feature_flags"`
	ResourceQuota         map[string]int    `json:"resource_quota"`
	ResourceUsage         map[string]int    `json:"resource_usage"`
	SubscriptionPlan      string            `json:"subscription_plan"`
	SubscriptionExpiresAt *string           `json:"subscription_expires_at,omitempty"`
}

// UpdateTenantConfigRequest 更新租户配置请求
type UpdateTenantConfigRequest struct {
	FeatureFlags          map[string]bool `json:"feature_flags,omitempty"`
	ResourceQuota         map[string]int  `json:"resource_quota,omitempty"`
	SubscriptionPlan      string          `json:"subscription_plan,omitempty"`
	SubscriptionExpiresAt *string         `json:"subscription_expires_at,omitempty"`
}

// GetTenantConfig 获取租户配置
func (s *TenantConfigService) GetTenantConfig(ctx context.Context, tenantID uint) (*TenantConfigDTO, error) {
	var tc models.TenantConfig
	if err := config.DB.Where("tenant_id = ?", tenantID).First(&tc).Error; err != nil {
		return nil, err
	}
	return s.toDTO(&tc), nil
}

// GetTenantConfigByTenantID 根据租户 ID 获取配置（内部使用）
func (s *TenantConfigService) GetTenantConfigByTenantID(ctx context.Context, tenantID uint) (*models.TenantConfig, error) {
	var tc models.TenantConfig
	if err := config.DB.Where("tenant_id = ?", tenantID).First(&tc).Error; err != nil {
		return nil, err
	}
	return &tc, nil
}

// UpdateTenantConfig 更新租户配置
func (s *TenantConfigService) UpdateTenantConfig(ctx context.Context, tenantID uint, req *UpdateTenantConfigRequest) (*TenantConfigDTO, error) {
	var tc models.TenantConfig
	if err := config.DB.Where("tenant_id = ?", tenantID).First(&tc).Error; err != nil {
		return nil, fmt.Errorf("tenant config not found: %w", err)
	}

	// 记录原始配置用于审计
	originalConfig := s.toDTO(&tc)

	// 更新功能模块配置
	if req.FeatureFlags != nil {
		flagsJSON, err := json.Marshal(req.FeatureFlags)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal feature flags: %w", err)
		}
		tc.FeatureFlags = string(flagsJSON)
	}

	// 更新资源配额配置
	if req.ResourceQuota != nil {
		quotaJSON, err := json.Marshal(req.ResourceQuota)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal resource quota: %w", err)
		}
		tc.ResourceQuota = string(quotaJSON)
	}

	// 更新订阅计划
	if req.SubscriptionPlan != "" {
		tc.SubscriptionPlan = req.SubscriptionPlan
	}

	// 更新订阅过期时间
	if req.SubscriptionExpiresAt != nil {
		if *req.SubscriptionExpiresAt == "" {
			tc.SubscriptionExpiresAt = nil
		} else {
			// 解析时间字符串
			expiresAt, err := time.Parse(time.RFC3339, *req.SubscriptionExpiresAt)
			if err != nil {
				// 尝试其他格式
				expiresAt, err = time.Parse("2006-01-02", *req.SubscriptionExpiresAt)
				if err != nil {
					return nil, fmt.Errorf("invalid date format: %w", err)
				}
			}
			tc.SubscriptionExpiresAt = &expiresAt
		}
	}

	if err := config.DB.Save(&tc).Error; err != nil {
		return nil, err
	}

	// 记录审计日志
	auditService := NewAuditService()
	operatorID := uint(0) // 从上下文中获取操作人 ID（如果有）
	
	// 记录配置变更审计日志
	_ = auditService.LogAction(
		operatorID,
		tenantID,
		AuditActionTenantConfigUpdate,
		"TenantConfig",
		tc.ID,
		originalConfig,
		s.toDTO(&tc),
	)

	return s.toDTO(&tc), nil
}

// CheckFeature 检查租户是否启用了某个功能
func (s *TenantConfigService) CheckFeature(ctx context.Context, tenantID uint, featureName string) (bool, error) {
	tc, err := s.GetTenantConfigByTenantID(ctx, tenantID)
	if err != nil {
		return false, err
	}

	var flags map[string]bool
	if tc.FeatureFlags == "" {
		return false, nil
	}

	if err := json.Unmarshal([]byte(tc.FeatureFlags), &flags); err != nil {
		return false, err
	}

	enabled, exists := flags[featureName]
	if !exists {
		return false, nil
	}
	return enabled, nil
}

// CheckQuota 检查租户资源配额
// 返回：是否充足、已使用量、配额限制、错误
func (s *TenantConfigService) CheckQuota(ctx context.Context, tenantID uint, resourceType string) (bool, int, int, error) {
	tc, err := s.GetTenantConfigByTenantID(ctx, tenantID)
	if err != nil {
		return false, 0, 0, err
	}

	var quota map[string]int
	var usage map[string]int

	if err := json.Unmarshal([]byte(tc.ResourceQuota), &quota); err != nil {
		return false, 0, 0, err
	}

	if err := json.Unmarshal([]byte(tc.ResourceUsage), &usage); err != nil {
		return false, 0, 0, err
	}

	limit, exists := quota["max_"+resourceType]
	if !exists {
		return true, 0, 0, nil // 无此配额限制
	}

	// -1 表示无限制
	if limit == -1 {
		return true, 0, -1, nil
	}

	used := usage["used_"+resourceType]
	return used < limit, used, limit, nil
}

// IncrementUsage 增加资源使用量
func (s *TenantConfigService) IncrementUsage(ctx context.Context, tenantID uint, resourceType string, amount int) error {
	tc, err := s.GetTenantConfigByTenantID(ctx, tenantID)
	if err != nil {
		return err
	}

	var usage map[string]int
	if err := json.Unmarshal([]byte(tc.ResourceUsage), &usage); err != nil {
		return err
	}

	key := "used_" + resourceType
	current, exists := usage[key]
	if !exists {
		current = 0
	}
	usage[key] = current + amount

	usageJSON, err := json.Marshal(usage)
	if err != nil {
		return err
	}
	tc.ResourceUsage = string(usageJSON)

	return config.DB.Save(tc).Error
}

// CheckAndIncrementUsage 检查并增加资源使用量（原子操作）
// 返回：是否成功（配额充足）、错误
func (s *TenantConfigService) CheckAndIncrementUsage(ctx context.Context, tenantID uint, resourceType string, amount int) (bool, error) {
	// 开启事务
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		var tc models.TenantConfig
		if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("tenant_id = ?", tenantID).First(&tc).Error; err != nil {
			return err
		}

		var quota map[string]int
		var usage map[string]int

		if err := json.Unmarshal([]byte(tc.ResourceQuota), &quota); err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(tc.ResourceUsage), &usage); err != nil {
			return err
		}

		limitKey := "max_" + resourceType
		limit, exists := quota[limitKey]
		if !exists {
			return nil // 无此配额限制
		}

		if limit == -1 {
			return nil // 无限制
		}

		usageKey := "used_" + resourceType
		current, exists := usage[usageKey]
		if !exists {
			current = 0
		}

		if current+amount > limit {
			return fmt.Errorf("quota exceeded: %s", resourceType)
		}

		usage[usageKey] = current + amount
		usageJSON, err := json.Marshal(usage)
		if err != nil {
			return err
		}
		tc.ResourceUsage = string(usageJSON)

		return tx.Save(&tc).Error
	})

	if err != nil {
		return false, err
	}
	return true, nil
}

// ResetQuotaUsage 重置租户配额使用统计
func (s *TenantConfigService) ResetQuotaUsage(ctx context.Context, tenantID uint, resourceType string) error {
	tc, err := s.GetTenantConfigByTenantID(ctx, tenantID)
	if err != nil {
		return err
	}

	var usage map[string]int
	if err := json.Unmarshal([]byte(tc.ResourceUsage), &usage); err != nil {
		return err
	}

	if resourceType == "" {
		// 重置所有配额
		usage = make(map[string]int)
	} else {
		key := "used_" + resourceType
		usage[key] = 0
	}

	usageJSON, err := json.Marshal(usage)
	if err != nil {
		return err
	}
	tc.ResourceUsage = string(usageJSON)

	return config.DB.Save(tc).Error
}

// GetQuotaUsage 获取租户配额使用情况
func (s *TenantConfigService) GetQuotaUsage(ctx context.Context, tenantID uint) (map[string]int, map[string]int, error) {
	tc, err := s.GetTenantConfigByTenantID(ctx, tenantID)
	if err != nil {
		return nil, nil, err
	}

	var quota map[string]int
	var usage map[string]int

	if err := json.Unmarshal([]byte(tc.ResourceQuota), &quota); err != nil {
		return nil, nil, err
	}

	if err := json.Unmarshal([]byte(tc.ResourceUsage), &usage); err != nil {
		return nil, nil, err
	}

	return quota, usage, nil
}

// GetFeatures 获取租户功能模块状态
func (s *TenantConfigService) GetFeatures(ctx context.Context, tenantID uint) (map[string]bool, error) {
	tc, err := s.GetTenantConfigByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	var flags map[string]bool
	if tc.FeatureFlags == "" {
		return make(map[string]bool), nil
	}

	if err := json.Unmarshal([]byte(tc.FeatureFlags), &flags); err != nil {
		return nil, err
	}

	return flags, nil
}

// toDTO 将模型转换为 DTO
func (s *TenantConfigService) toDTO(tc *models.TenantConfig) *TenantConfigDTO {
	dto := &TenantConfigDTO{
		ID:               tc.ID,
		TenantID:         tc.TenantID,
		SubscriptionPlan: tc.SubscriptionPlan,
	}

	// 解析功能模块
	if tc.FeatureFlags != "" {
		json.Unmarshal([]byte(tc.FeatureFlags), &dto.FeatureFlags)
	} else {
		dto.FeatureFlags = make(map[string]bool)
	}

	// 解析资源配额
	if tc.ResourceQuota != "" {
		json.Unmarshal([]byte(tc.ResourceQuota), &dto.ResourceQuota)
	} else {
		dto.ResourceQuota = make(map[string]int)
	}

	// 解析资源使用
	if tc.ResourceUsage != "" {
		json.Unmarshal([]byte(tc.ResourceUsage), &dto.ResourceUsage)
	} else {
		dto.ResourceUsage = make(map[string]int)
	}

	// 处理订阅过期时间
	if tc.SubscriptionExpiresAt != nil {
		expiresStr := tc.SubscriptionExpiresAt.Format(time.RFC3339)
		dto.SubscriptionExpiresAt = &expiresStr
	}

	return dto
}

// CreateDefaultConfig 为租户创建默认配置
func (s *TenantConfigService) CreateDefaultConfig(ctx context.Context, tenantID uint, subscriptionPlan string) error {
	// 检查配置是否已存在
	var existing models.TenantConfig
	if err := config.DB.Where("tenant_id = ?", tenantID).First(&existing).Error; err == nil {
		return nil // 已存在，无需创建
	}

	defaultFlags := map[string]bool{
		"image_management":   true,
		"page_config":        true,
		"multi_language":     true,
		"contact_form":       true,
		"content_management": true,
	}

	defaultQuota := map[string]int{
		"max_images":        50,
		"max_storage_mb":    512,
		"max_content_items": 20,
		"max_users":         3,
	}

	defaultUsage := map[string]int{
		"used_images":        0,
		"used_storage_mb":    0,
		"used_content_items": 0,
		"used_users":         0,
	}

	// 根据订阅计划调整默认值
	switch subscriptionPlan {
	case "pro":
		defaultQuota["max_images"] = 500
		defaultQuota["max_storage_mb"] = 5120
		defaultQuota["max_content_items"] = 100
		defaultQuota["max_users"] = 10
		defaultFlags["multi_language"] = true
	case "enterprise":
		defaultQuota["max_images"] = -1
		defaultQuota["max_storage_mb"] = -1
		defaultQuota["max_content_items"] = -1
		defaultQuota["max_users"] = -1
	}

	flagsJSON, _ := json.Marshal(defaultFlags)
	quotaJSON, _ := json.Marshal(defaultQuota)
	usageJSON, _ := json.Marshal(defaultUsage)

	configModel := &models.TenantConfig{
		TenantID:              tenantID,
		FeatureFlags:          string(flagsJSON),
		ResourceQuota:         string(quotaJSON),
		ResourceUsage:         string(usageJSON),
		SubscriptionPlan:      subscriptionPlan,
		SubscriptionExpiresAt: nil,
	}

	return config.DB.Create(configModel).Error
}