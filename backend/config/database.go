package config

import (
	"encoding/json"

	"medical-device-cms/backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("medical.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Image{},
		&models.PageConfig{},
		&models.ModuleConfig{},
		&models.ContentItem{},
		&models.LanguageText{},
		&models.LanguageTextVersion{},
		&models.SiteSetting{},
		&models.ContactSubmission{},
	)
	if err != nil {
		return err
	}

	initDefaultUser()
	initDefaultPageConfig()
	initDefaultSiteSettings()
	return nil
}

func initDefaultUser() {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		DB.Create(&models.User{
			Username: "admin",
			Password: string(hashedPassword),
			Role:     "admin",
		})
	}
}

func initDefaultPageConfig() {
	defaultConfigs := map[string]string{
		"banner":   `{"enabled": true, "title": "Professional Medical Device Manufacturer & Supplier", "subtitle": "CE, FDA, ISO Certified Medical Devices. We provide high-quality hospital equipment, surgical supplies and disposable medical products with OEM/ODM global wholesale service."}`,
		"about":    `{"enabled": true, "title": "About Our Company", "content": "Founded in 2010, our company is a professional manufacturer and exporter engaged in the research, development, production and sales of high-standard medical devices."}`,
		"products": `{"enabled": true, "title": "Our Main Products", "show_count": 6}`,
	}

	for pageName, configData := range defaultConfigs {
		var config models.PageConfig
		DB.Where("page_name = ?", pageName).First(&config)
		if config.ID == 0 {
			DB.Create(&models.PageConfig{
				PageName:   pageName,
				ConfigData: configData,
			})
		}
	}
}

func GetPageConfig(pageName string) (map[string]interface{}, error) {
	var config models.PageConfig
	if err := DB.Where("page_name = ?", pageName).First(&config).Error; err != nil {
		return nil, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(config.ConfigData), &data); err != nil {
		return nil, err
	}
	return data, nil
}

func GetAllPageConfigs() ([]models.PageConfig, error) {
	var configs []models.PageConfig
	if err := DB.Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

func initDefaultSiteSettings() {
	var count int64
	DB.Model(&models.SiteSetting{}).Where("key = ?", "zh_site_title").Count(&count)
	if count == 0 {
		DB.Create(&models.SiteSetting{
			Key:   "zh_site_title",
			Value: "专业医疗器械制造商与出口商",
		})
	}

	DB.Model(&models.SiteSetting{}).Where("key = ?", "en_site_title").Count(&count)
	if count == 0 {
		DB.Create(&models.SiteSetting{
			Key:   "en_site_title",
			Value: "Professional Medical Device Manufacturer & Exporter",
		})
	}

	// 初始化默认 Logo 配置
	DB.Model(&models.SiteSetting{}).Where("key = ?", "zh_site_logo").Count(&count)
	if count == 0 {
		DB.Create(&models.SiteSetting{
			Key:   "zh_site_logo",
			Value: "医疗",
		})
	}

	DB.Model(&models.SiteSetting{}).Where("key = ?", "en_site_logo").Count(&count)
	if count == 0 {
		DB.Create(&models.SiteSetting{
			Key:   "en_site_logo",
			Value: "MEDICAL",
		})
	}

	DB.Model(&models.SiteSetting{}).Where("key = ?", "site_logo_color").Count(&count)
	if count == 0 {
		DB.Create(&models.SiteSetting{
			Key:   "site_logo_color",
			Value: "#06a499",
		})
	}
}
