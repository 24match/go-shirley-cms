package services

import (
	"os"
	"strconv"
	"strings"

	"medical-device-cms/backend/config"
	"medical-device-cms/backend/models"
)

type ModuleService struct{}

func NewModuleService() *ModuleService {
	return &ModuleService{}
}

func (s *ModuleService) GetModuleConfigs(moduleName string) (interface{}, error) {
	if moduleName != "" {
		var moduleConfig models.ModuleConfig
		err := config.DB.Where("module_name = ?", moduleName).First(&moduleConfig).Error
		return moduleConfig, err
	}

	var moduleConfigs []models.ModuleConfig
	err := config.DB.Order("sort_order ASC").Find(&moduleConfigs).Error
	return moduleConfigs, err
}

func (s *ModuleService) GetModuleConfig(name string) (*models.ModuleConfig, error) {
	var moduleConfig models.ModuleConfig
	err := config.DB.Where("module_name = ?", name).First(&moduleConfig).Error
	return &moduleConfig, err
}

func (s *ModuleService) SaveModuleConfig(moduleName string, updates map[string]interface{}) (*models.ModuleConfig, error) {
	var moduleConfig models.ModuleConfig

	contentType := updates["contentType"].(string)
	if strings.HasPrefix(contentType, "multipart/form-data") {
		enabled := updates["enabled"].(string) == "true"
		title := updates["title"].(string)
		subtitle := updates["subtitle"].(string)
		content := updates["content"].(string)
		sortOrder, _ := strconv.Atoi(updates["sortOrder"].(string))
		description := updates["description"].(string)
		extraData := updates["extraData"].(string)

		config.DB.Where("module_name = ?", moduleName).First(&moduleConfig)
		if moduleConfig.ID == 0 {
			moduleConfig = models.ModuleConfig{ModuleName: moduleName}
		}

		moduleConfig.Enabled = enabled
		moduleConfig.Title = title
		moduleConfig.Subtitle = subtitle
		moduleConfig.Content = content
		moduleConfig.SortOrder = sortOrder
		moduleConfig.Description = description
		moduleConfig.ExtraData = extraData

		if newFile, ok := updates["file"].(interface{}); ok && newFile != nil {
			if moduleConfig.ImagePath != "" {
				os.Remove("./uploads/" + moduleConfig.ImagePath)
			}
			file := newFile.(map[string]interface{})
			filename := generateFilename(file["filename"].(string))
			fileData := file["data"].([]byte)
			os.WriteFile("./uploads/"+filename, fileData, 0644)
			moduleConfig.ImagePath = filename
		}

		err := config.DB.Save(&moduleConfig).Error
		return &moduleConfig, err
	} else {
		var req models.ModuleConfig
		if data, ok := updates["data"].(models.ModuleConfig); ok {
			req = data
		}

		if req.ModuleName == "" {
			req.ModuleName = moduleName
		}

		config.DB.Where("module_name = ?", req.ModuleName).First(&models.ModuleConfig{})
		err := config.DB.Save(&req).Error
		return &req, err
	}
}

func (s *ModuleService) DeleteModuleConfig(name string) error {
	var moduleConfig models.ModuleConfig
	if err := config.DB.Where("module_name = ?", name).First(&moduleConfig).Error; err != nil {
		return err
	}

	if moduleConfig.ImagePath != "" {
		os.Remove("./uploads/" + moduleConfig.ImagePath)
	}

	return config.DB.Delete(&moduleConfig).Error
}

func (s *ModuleService) GetPublicModuleConfigs() ([]models.ModuleConfig, error) {
	var moduleConfigs []models.ModuleConfig
	err := config.DB.Where("enabled = ?", true).Order("sort_order ASC").Find(&moduleConfigs).Error
	return moduleConfigs, err
}