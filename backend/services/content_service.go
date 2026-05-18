package services

import (
	"os"
	"strconv"

	"medical-device-cms/backend/config"
	"medical-device-cms/backend/models"
)

type ContentService struct{}

func NewContentService() *ContentService {
	return &ContentService{}
}

func (s *ContentService) GetContentItems(section string) ([]models.ContentItem, error) {
	var items []models.ContentItem
	query := config.DB.Order("sort_order ASC")
	if section != "" {
		query = query.Where("section = ?", section)
	}
	err := query.Find(&items).Error
	return items, err
}

func (s *ContentService) CreateContentItem(section, title, description, imagePath string, sortOrder int, icon string) (*models.ContentItem, error) {
	item := models.ContentItem{
		Section:     section,
		Title:       title,
		Description: description,
		ImagePath:   imagePath,
		SortOrder:   sortOrder,
		Icon:        icon,
	}
	err := config.DB.Create(&item).Error
	return &item, err
}

func (s *ContentService) UpdateContentItem(id uint, updates map[string]interface{}) (*models.ContentItem, error) {
	var item models.ContentItem
	if err := config.DB.First(&item, id).Error; err != nil {
		return nil, err
	}

	if section, ok := updates["section"].(string); ok && section != "" {
		item.Section = section
	}
	if title, ok := updates["title"].(string); ok && title != "" {
		item.Title = title
	}
	if description, ok := updates["description"].(string); ok && description != "" {
		item.Description = description
	}
	if sortOrderStr, ok := updates["sort_order"].(string); ok && sortOrderStr != "" {
		sortOrder, _ := strconv.Atoi(sortOrderStr)
		item.SortOrder = sortOrder
	}

	if newFile, ok := updates["file"].(interface{}); ok && newFile != nil {
		if item.ImagePath != "" {
			os.Remove("./uploads/" + item.ImagePath)
		}
		file := newFile.(map[string]interface{})
		filename := generateFilename(file["filename"].(string))
		fileData := file["data"].([]byte)
		os.WriteFile("./uploads/"+filename, fileData, 0644)
		item.ImagePath = filename
	} else {
		if imagePath, ok := updates["image_path"].(string); ok && imagePath != "" {
			item.ImagePath = imagePath
		}
	}

	err := config.DB.Save(&item).Error
	return &item, err
}

func (s *ContentService) DeleteContentItem(id uint) error {
	var item models.ContentItem
	if err := config.DB.First(&item, id).Error; err != nil {
		return err
	}
	return config.DB.Delete(&item).Error
}