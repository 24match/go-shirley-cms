package services

import (
	"context"
	"encoding/json"
	"fmt"

	"medical-device-cms/backend/config"
	"medical-device-cms/backend/models"
	"gorm.io/gorm"
)

// QuotaService 配额管理服务
type QuotaService struct {
	tenantConfigService *TenantConfigService
}

// NewQuotaService 创建配额管理服务实例
func NewQuotaService() *QuotaService {
	return &QuotaService{
		tenantConfigService: NewTenantConfigService(),
	}
}

// QuotaInfo 配额信息
type QuotaInfo struct {
	ResourceType string `json:"resource_type"`
	Used         int    `json:"used"`
	Limit        int    `json:"limit"`
	Available    int    `json:"available"`
	IsUnlimited  bool   `json:"is_unlimited"`
}

// CheckQuota 检查配额是否充足
// 返回：是否充足、配额信息、错误
func (s *QuotaService) CheckQuota(ctx context.Context, tenantID uint, resourceType string) (bool, *QuotaInfo, error) {
	tc, err := s.tenantConfigService.GetTenantConfigByTenantID(ctx, tenantID)
	if err != nil {
		return false, nil, err
	}

	var quota map[string]int
	var usage map[string]int

	if err := json.Unmarshal([]byte(tc.ResourceQuota), &quota); err != nil {
		return false, nil, err
	}

	if err := json.Unmarshal([]byte(tc.ResourceUsage), &usage); err != nil {
		return false, nil, err
	}

	limitKey := "max_" + resourceType
	usageKey := "used_" + resourceType

	limit, exists := quota[limitKey]
	if !exists {
		// 无此配额限制，返回充足
		return true, &QuotaInfo{
			ResourceType: resourceType,
			Used:         0,
			Limit:        -1,
			Available:    -1,
			IsUnlimited:  true,
		}, nil
	}

	// -1 表示无限制
	if limit == -1 {
		used := usage[usageKey]
		return true, &QuotaInfo{
			ResourceType: resourceType,
			Used:         used,
			Limit:        -1,
			Available:    -1,
			IsUnlimited:  true,
		}, nil
	}

	used := usage[usageKey]
	available := limit - used

	return used < limit, &QuotaInfo{
		ResourceType: resourceType,
		Used:         used,
		Limit:        limit,
		Available:    available,
		IsUnlimited:  false,
	}, nil
}

// CheckAndIncrement 检查并增加使用量（原子操作）
// 返回：是否成功、错误
func (s *QuotaService) CheckAndIncrement(ctx context.Context, tenantID uint, resourceType string, amount int) (bool, error) {
	var success bool
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
			// 无此配额限制，直接增加
			return s.incrementUsageInTx(tx, &tc, resourceType, amount, usage)
		}

		// -1 表示无限制
		if limit == -1 {
			return s.incrementUsageInTx(tx, &tc, resourceType, amount, usage)
		}

		usageKey := "used_" + resourceType
		current, exists := usage[usageKey]
		if !exists {
			current = 0
		}

		if current+amount > limit {
			return fmt.Errorf("quota exceeded: %s (used: %d, limit: %d)", resourceType, current, limit)
		}

		return s.incrementUsageInTx(tx, &tc, resourceType, amount, usage)
	})

	if err != nil {
		return false, err
	}
	return success, nil
}

// incrementUsageInTx 在事务中增加使用量
func (s *QuotaService) incrementUsageInTx(tx *gorm.DB, tc *models.TenantConfig, resourceType string, amount int, usage map[string]int) error {
	usageKey := "used_" + resourceType
	current, exists := usage[usageKey]
	if !exists {
		current = 0
	}
	usage[usageKey] = current + amount

	usageJSON, err := json.Marshal(usage)
	if err != nil {
		return err
	}
	tc.ResourceUsage = string(usageJSON)

	return tx.Save(tc).Error
}

// GetUsage 获取租户资源使用情况
func (s *QuotaService) GetUsage(ctx context.Context, tenantID uint) (map[string]int, error) {
	tc, err := s.tenantConfigService.GetTenantConfigByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	var usage map[string]int
	if err := json.Unmarshal([]byte(tc.ResourceUsage), &usage); err != nil {
		return nil, err
	}

	return usage, nil
}

// GetQuota 获取租户配额配置
func (s *QuotaService) GetQuota(ctx context.Context, tenantID uint) (map[string]int, error) {
	tc, err := s.tenantConfigService.GetTenantConfigByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	var quota map[string]int
	if err := json.Unmarshal([]byte(tc.ResourceQuota), &quota); err != nil {
		return nil, err
	}

	return quota, nil
}

// GetAllQuotaInfo 获取所有配额信息
func (s *QuotaService) GetAllQuotaInfo(ctx context.Context, tenantID uint) ([]*QuotaInfo, error) {
	quota, err := s.GetQuota(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	usage, err := s.GetUsage(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	var result []*QuotaInfo

	for limitKey, limit := range quota {
		// 从 "max_images" 提取 "images"
		resourceType := limitKey[4:] // 去掉 "max_" 前缀
		usageKey := "used_" + resourceType

		used := usage[usageKey]
		isUnlimited := limit == -1

		var available int
		if isUnlimited {
			available = -1
		} else {
			available = limit - used
		}

		result = append(result, &QuotaInfo{
			ResourceType: resourceType,
			Used:         used,
			Limit:        limit,
			Available:    available,
			IsUnlimited:  isUnlimited,
		})
	}

	return result, nil
}

// ResetQuotaUsage 重置配额使用统计
func (s *QuotaService) ResetQuotaUsage(ctx context.Context, tenantID uint, resourceType string) error {
	tc, err := s.tenantConfigService.GetTenantConfigByTenantID(ctx, tenantID)
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

// UpdateQuota 更新配额配置
func (s *QuotaService) UpdateQuota(ctx context.Context, tenantID uint, newQuota map[string]int) error {
	tc, err := s.tenantConfigService.GetTenantConfigByTenantID(ctx, tenantID)
	if err != nil {
		return err
	}

	quotaJSON, err := json.Marshal(newQuota)
	if err != nil {
		return err
	}
	tc.ResourceQuota = string(quotaJSON)

	return config.DB.Save(tc).Error
}

// IsQuotaExceeded 检查配额是否已超限
func (s *QuotaService) IsQuotaExceeded(ctx context.Context, tenantID uint, resourceType string) (bool, *QuotaInfo, error) {
	isSufficient, info, err := s.CheckQuota(ctx, tenantID, resourceType)
	if err != nil {
		return false, nil, err
	}
	return !isSufficient, info, nil
}