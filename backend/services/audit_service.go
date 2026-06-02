package services

import (
	"encoding/json"
	"medical-device-cms/backend/config"
	"medical-device-cms/backend/models"
)

type AuditService struct{}

func NewAuditService() *AuditService {
	return &AuditService{}
}

// LogAction 记录操作日志
// operatorID: 操作人 ID
// tenantID: 租户 ID
// action: 操作类型（CREATE, UPDATE, DELETE, LOGIN, LOGOUT 等）
// resourceType: 资源类型（user, image, content, module 等）
// resourceID: 资源 ID
// beforeValue: 操作前的值（JSON 格式）
// afterValue: 操作后的值（JSON 格式）
func (s *AuditService) LogAction(
	operatorID uint,
	tenantID uint,
	action string,
	resourceType string,
	resourceID uint,
	beforeValue interface{},
	afterValue interface{},
) error {
	var beforeJSON, afterJSON string

	if beforeValue != nil {
		data, err := json.Marshal(beforeValue)
		if err == nil {
			beforeJSON = string(data)
		}
	}

	if afterValue != nil {
		data, err := json.Marshal(afterValue)
		if err == nil {
			afterJSON = string(data)
		}
	}

	auditLog := models.AuditLog{
		OperatorID:   operatorID,
		TenantID:     tenantID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		BeforeValue:  beforeJSON,
		AfterValue:   afterJSON,
	}

	return config.DB.Create(&auditLog).Error
}

// ListLogs 查询审计日志列表
// tenantID: 租户 ID（0 表示查询所有租户，仅超级管理员可用）
// operatorID: 操作人 ID（0 表示不限制）
// action: 操作类型（空表示不限制）
// resourceType: 资源类型（空表示不限制）
// limit: 分页限制
// offset: 分页偏移
func (s *AuditService) ListLogs(
	tenantID uint,
	operatorID uint,
	action string,
	resourceType string,
	limit int,
	offset int,
) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	query := config.DB.Model(&models.AuditLog{})

	// 租户过滤（tenantID 为 0 时不过滤，用于超级管理员）
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}

	// 操作人过滤
	if operatorID > 0 {
		query = query.Where("operator_id = ?", operatorID)
	}

	// 操作类型过滤
	if action != "" {
		query = query.Where("action = ?", action)
	}

	// 资源类型过滤
	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error

	return logs, total, err
}

// GetLogsByTenant 按租户查询审计日志
func (s *AuditService) GetLogsByTenant(tenantID uint, limit int, offset int) ([]models.AuditLog, int64, error) {
	return s.ListLogs(tenantID, 0, "", "", limit, offset)
}

// GetLogsByOperator 按操作人查询审计日志
func (s *AuditService) GetLogsByOperator(tenantID uint, operatorID uint, limit int, offset int) ([]models.AuditLog, int64, error) {
	return s.ListLogs(tenantID, operatorID, "", "", limit, offset)
}

// GetLogsByAction 按操作类型查询审计日志
func (s *AuditService) GetLogsByAction(tenantID uint, action string, limit int, offset int) ([]models.AuditLog, int64, error) {
	return s.ListLogs(tenantID, 0, action, "", limit, offset)
}

// GetLogsByResource 按资源类型查询审计日志
func (s *AuditService) GetLogsByResource(tenantID uint, resourceType string, limit int, offset int) ([]models.AuditLog, int64, error) {
	return s.ListLogs(tenantID, 0, "", resourceType, limit, offset)
}

// GetAuditStats 获取审计统计信息
func (s *AuditService) GetAuditStats(tenantID uint) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总操作数
	var totalActions int64
	query := config.DB.Model(&models.AuditLog{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	query.Count(&totalActions)
	stats["total_actions"] = totalActions

	// 按操作类型统计
	type ActionCount struct {
		Action string
		Count  int64
	}
	var actionCounts []ActionCount
	query.Select("action, COUNT(*) as count").
		Group("action").
		Scan(&actionCounts)
	stats["actions_by_type"] = actionCounts

	// 按资源类型统计
	type ResourceTypeCount struct {
		ResourceType string
		Count        int64
	}
	var resourceTypeCounts []ResourceTypeCount
	query.Select("resource_type, COUNT(*) as count").
		Group("resource_type").
		Scan(&resourceTypeCounts)
	stats["actions_by_resource"] = resourceTypeCounts

	return stats, nil
}