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
	)
	if err != nil {
		return err
	}

	initDefaultUser()
	initDefaultPageConfig()
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