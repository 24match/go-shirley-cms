package services

import (
	"medical-device-cms/backend/config"
	"medical-device-cms/backend/models"
)

type LanguageService struct{}

func NewLanguageService() *LanguageService {
	return &LanguageService{}
}

func (s *LanguageService) GetPublicLanguageTexts() (map[string]map[string]string, error) {
	var texts []models.LanguageText
	if err := config.DB.Find(&texts).Error; err != nil {
		return nil, err
	}

	result := make(map[string]map[string]string)
	result["en"] = make(map[string]string)
	result["zh"] = make(map[string]string)

	for _, text := range texts {
		result["en"][text.Key] = text.EnText
		result["zh"][text.Key] = text.ZhText
	}

	return result, nil
}

func (s *LanguageService) GetLanguageTexts(module string) ([]models.LanguageText, error) {
	var texts []models.LanguageText
	query := config.DB
	if module != "" {
		query = query.Where("module = ?", module)
	}
	err := query.Find(&texts).Error
	return texts, err
}

func (s *LanguageService) GetLanguageText(id uint) (*models.LanguageText, error) {
	var text models.LanguageText
	err := config.DB.First(&text, id).Error
	return &text, err
}

func (s *LanguageService) CreateLanguageText(key, module, enText, zhText, description string) (*models.LanguageText, error) {
	var existing models.LanguageText
	if err := config.DB.Where("key = ?", key).First(&existing).Error; err == nil {
		return nil, nil
	}

	text := models.LanguageText{
		Key:         key,
		Module:      module,
		EnText:      enText,
		ZhText:      zhText,
		Description: description,
		Version:     1,
	}

	err := config.DB.Create(&text).Error
	return &text, err
}

func (s *LanguageService) UpdateLanguageText(id uint, enText, zhText, description string) (*models.LanguageText, error) {
	var text models.LanguageText
	if err := config.DB.First(&text, id).Error; err != nil {
		return nil, err
	}

	config.DB.Create(&models.LanguageTextVersion{
		LanguageTextID: text.ID,
		Key:            text.Key,
		Module:         text.Module,
		EnText:         text.EnText,
		ZhText:         text.ZhText,
		Description:    text.Description,
		Version:        text.Version,
		UpdatedAt:      text.UpdatedAt,
	})

	text.EnText = enText
	text.ZhText = zhText
	text.Description = description
	text.Version++

	err := config.DB.Save(&text).Error
	return &text, err
}

func (s *LanguageService) DeleteLanguageText(id uint) error {
	var text models.LanguageText
	if err := config.DB.First(&text, id).Error; err != nil {
		return err
	}

	config.DB.Delete(&models.LanguageTextVersion{}, "language_text_id = ?", id)
	return config.DB.Delete(&text).Error
}

func (s *LanguageService) GetLanguageTextVersions(id uint) ([]models.LanguageTextVersion, error) {
	var versions []models.LanguageTextVersion
	err := config.DB.Where("language_text_id = ?", id).Order("version DESC").Find(&versions).Error
	return versions, err
}

func (s *LanguageService) RestoreLanguageTextVersion(id, version uint) (*models.LanguageText, error) {
	var text models.LanguageText
	if err := config.DB.First(&text, id).Error; err != nil {
		return nil, err
	}

	var oldVersion models.LanguageTextVersion
	if err := config.DB.Where("language_text_id = ? AND version = ?", id, version).First(&oldVersion).Error; err != nil {
		return nil, err
	}

	config.DB.Create(&models.LanguageTextVersion{
		LanguageTextID: text.ID,
		Key:            text.Key,
		Module:         text.Module,
		EnText:         text.EnText,
		ZhText:         text.ZhText,
		Description:    text.Description,
		Version:        text.Version,
		UpdatedAt:      text.UpdatedAt,
	})

	text.EnText = oldVersion.EnText
	text.ZhText = oldVersion.ZhText
	text.Description = oldVersion.Description
	text.Version++

	err := config.DB.Save(&text).Error
	return &text, err
}