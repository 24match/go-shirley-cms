package services

import (
	"medical-device-cms/backend/config"
	"medical-device-cms/backend/models"
)

type ContactSubmissionService struct{}

func NewContactSubmissionService() *ContactSubmissionService {
	return &ContactSubmissionService{}
}

// CreateSubmission 创建联系表单提交记录
func (s *ContactSubmissionService) CreateSubmission(name, email, company, inquiry, ipAddress string) (*models.ContactSubmission, error) {
	submission := &models.ContactSubmission{
		Name:      name,
		Email:     email,
		Company:   company,
		Inquiry:   inquiry,
		IPAddress: ipAddress,
		IsRead:    false,
	}

	err := config.DB.Create(submission).Error
	if err != nil {
		return nil, err
	}

	return submission, nil
}

// GetSubmissions 获取联系表单提交列表（支持分页）
func (s *ContactSubmissionService) GetSubmissions(page, pageSize int) ([]models.ContactSubmission, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var submissions []models.ContactSubmission
	var total int64

	// 获取总数
	err := config.DB.Model(&models.ContactSubmission{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询，按创建时间倒序
	offset := (page - 1) * pageSize
	err = config.DB.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&submissions).Error
	if err != nil {
		return nil, 0, err
	}

	return submissions, total, nil
}

// MarkAsRead 标记为已读
func (s *ContactSubmissionService) MarkAsRead(id uint) error {
	return config.DB.Model(&models.ContactSubmission{}).Where("id = ?", id).Update("is_read", true).Error
}

// DeleteSubmission 删除联系表单提交记录
func (s *ContactSubmissionService) DeleteSubmission(id uint) error {
	return config.DB.Delete(&models.ContactSubmission{}, id).Error
}