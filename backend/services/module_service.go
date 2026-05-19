package services

import (
	"encoding/json"
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
		if err != nil {
			return moduleConfig, err
		}
		// Merge extra data into module config for frontend
		return mergeExtraData(&moduleConfig), nil
	}

	var moduleConfigs []models.ModuleConfig
	err := config.DB.Order("sort_order ASC").Find(&moduleConfigs).Error
	if err != nil {
		return moduleConfigs, err
	}
	// Merge extra data for each module
	var results []map[string]interface{}
	for _, mc := range moduleConfigs {
		results = append(results, mergeExtraData(&mc))
	}
	return results, nil
}

// mergeExtraData 将 ExtraData JSON 中的字段合并到模块配置中，以便前端可以访问
func mergeExtraData(mc *models.ModuleConfig) map[string]interface{} {
	result := make(map[string]interface{})
	// 先复制模型中的基础字段
	result["id"] = mc.ID
	result["createdAt"] = mc.CreatedAt
	result["updatedAt"] = mc.UpdatedAt
	result["moduleName"] = mc.ModuleName
	result["enabled"] = mc.Enabled
	result["zhTitle"] = mc.ZhTitle
	result["enTitle"] = mc.EnTitle
	result["zhSubtitle"] = mc.ZhSubtitle
	result["enSubtitle"] = mc.EnSubtitle
	result["zhContent"] = mc.ZhContent
	result["enContent"] = mc.EnContent
	result["title"] = mc.Title
	result["subtitle"] = mc.Subtitle
	result["content"] = mc.Content
	result["imagePath"] = mc.ImagePath
	result["sortOrder"] = mc.SortOrder
	result["extraData"] = mc.ExtraData
	result["zhDescription"] = mc.ZhDescription
	result["enDescription"] = mc.EnDescription
	result["description"] = mc.Description

	// 解析并合并 ExtraData 中的字段
	if mc.ExtraData != "" {
		var extra map[string]interface{}
		if err := json.Unmarshal([]byte(mc.ExtraData), &extra); err == nil {
			for k, v := range extra {
				result[k] = v
			}
		}
	}

	return result
}

func (s *ModuleService) GetModuleConfig(name string) (*models.ModuleConfig, error) {
	var moduleConfig models.ModuleConfig
	err := config.DB.Where("module_name = ?", name).First(&moduleConfig).Error
	return &moduleConfig, err
}

func (s *ModuleService) SaveModuleConfig(moduleName string, updates map[string]interface{}) (*models.ModuleConfig, error) {
	var moduleConfig models.ModuleConfig

	contentType := getStringValue(updates, "contentType", "")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		enabled := getStringValue(updates, "enabled", "") == "true"
		title := getStringValue(updates, "title", "")
		subtitle := getStringValue(updates, "subtitle", "")
		content := getStringValue(updates, "content", "")
		sortOrderStr := getStringValue(updates, "sortOrder", "0")
		sortOrder, _ := strconv.Atoi(sortOrderStr)
		description := getStringValue(updates, "description", "")
		extraData := getStringValue(updates, "extraData", "")
		// 新增多语言字段
		zhTitle := getStringValue(updates, "zhTitle", "")
		enTitle := getStringValue(updates, "enTitle", "")
		zhSubtitle := getStringValue(updates, "zhSubtitle", "")
		enSubtitle := getStringValue(updates, "enSubtitle", "")
		zhContent := getStringValue(updates, "zhContent", "")
		enContent := getStringValue(updates, "enContent", "")
		zhDescription := getStringValue(updates, "zhDescription", "")
		enDescription := getStringValue(updates, "enDescription", "")

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
		// 设置多语言字段
		moduleConfig.ZhTitle = zhTitle
		moduleConfig.EnTitle = enTitle
		moduleConfig.ZhSubtitle = zhSubtitle
		moduleConfig.EnSubtitle = enSubtitle
		moduleConfig.ZhContent = zhContent
		moduleConfig.EnContent = enContent
		moduleConfig.ZhDescription = zhDescription
		moduleConfig.EnDescription = enDescription

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

func getStringValue(updates map[string]interface{}, key, defaultValue string) string {
	if val, ok := updates[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return defaultValue
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