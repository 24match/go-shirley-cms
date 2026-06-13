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
	DB, err = gorm.Open(sqlite.Open("backend/medical.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// 首先迁移租户和审计日志表
	err = DB.AutoMigrate(
		&models.Tenant{},
		&models.AuditLog{},
	)
	if err != nil {
		return err
	}

	// 然后迁移其他表（现在包含 tenant_id 字段）
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
		&models.TenantConfig{},
	)
	if err != nil {
		return err
	}

	// 执行 SaaS 迁移（如果需要）
	MigrateToSaaS()
	
	// 初始化租户配置
	InitTenantConfigs()

	// 初始化默认数据
	initDefaultUser()
	initDefaultPageConfig()
	initDefaultSiteSettings()
	return nil
}

// MigrateToSaaS 执行 SaaS 架构迁移
// 创建默认租户，将现有数据关联到默认租户
func MigrateToSaaS() {
	// 检查是否已存在租户
	var count int64
	DB.Model(&models.Tenant{}).Count(&count)
	if count > 0 {
		return // 已有租户，无需迁移
	}

	// 创建默认租户
	defaultTenant := models.Tenant{
		TenantCode:         "default",
		Name:               "Default Tenant",
		Status:             models.TenantStatusActive,
		SubDomain:          "default",
		SubscriptionPlan:   models.SubscriptionPlanFree,
		SubscriptionExpiresAt: nil,
	}
	DB.Create(&defaultTenant)

	// 将所有现有用户关联到默认租户
	DB.Model(&models.User{}).Where("tenant_id = 0").Update("tenant_id", defaultTenant.ID)

	// 将所有现有图片关联到默认租户
	DB.Model(&models.Image{}).Where("tenant_id = 0").Update("tenant_id", defaultTenant.ID)

	// 将所有现有页面配置关联到默认租户
	DB.Model(&models.PageConfig{}).Where("tenant_id = 0").Update("tenant_id", defaultTenant.ID)

	// 将所有现有模块配置关联到默认租户
	DB.Model(&models.ModuleConfig{}).Where("tenant_id = 0").Update("tenant_id", defaultTenant.ID)

	// 将所有现有内容项关联到默认租户
	DB.Model(&models.ContentItem{}).Where("tenant_id = 0").Update("tenant_id", defaultTenant.ID)

	// 将所有现有语言文本关联到默认租户
	DB.Model(&models.LanguageText{}).Where("tenant_id = 0").Update("tenant_id", defaultTenant.ID)

	// 将所有现有语言版本文本关联到默认租户
	DB.Model(&models.LanguageTextVersion{}).Where("tenant_id = 0").Update("tenant_id", defaultTenant.ID)

	// 将所有现有站点设置关联到默认租户
	DB.Model(&models.SiteSetting{}).Where("tenant_id = 0").Update("tenant_id", defaultTenant.ID)

	// 将所有现有联系表单提交关联到默认租户
	DB.Model(&models.ContactSubmission{}).Where("tenant_id = 0").Update("tenant_id", defaultTenant.ID)
}

func initDefaultUser() {
	// 检查是否已存在超级管理员
	var superAdminCount int64
	DB.Model(&models.User{}).Where("role = ?", "superadmin").Count(&superAdminCount)
	if superAdminCount == 0 {
		// 创建默认超级管理员账户
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("superadmin123"), bcrypt.DefaultCost)
		DB.Create(&models.User{
			Username: "superadmin",
			Password: string(hashedPassword),
			Role:     "superadmin",
			TenantID: 0, // 超级管理员不属于任何租户
		})
	}

	// 检查是否已存在租户管理员
	var tenantAdminCount int64
	DB.Model(&models.User{}).Where("role = ?", "tenant_admin").Count(&tenantAdminCount)
	if tenantAdminCount == 0 {
		// 获取默认租户
		var defaultTenant models.Tenant
		DB.Where("tenant_code = ?", "default").First(&defaultTenant)
		
		if defaultTenant.ID > 0 {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
			DB.Create(&models.User{
				TenantID: defaultTenant.ID,
				Username: "admin",
				Password: string(hashedPassword),
				Role:     "tenant_admin",
			})
		}
	}
}

func initDefaultPageConfig() {
	// 获取默认租户
	var defaultTenant models.Tenant
	DB.Where("tenant_code = ?", "default").First(&defaultTenant)
	
	defaultTenantID := uint(0)
	if defaultTenant.ID > 0 {
		defaultTenantID = defaultTenant.ID
	}

	defaultConfigs := map[string]string{
		"banner":   `{"enabled": true, "title": "Professional Medical Device Manufacturer & Supplier", "subtitle": "CE, FDA, ISO Certified Medical Devices. We provide high-quality hospital equipment, surgical supplies and disposable medical products with OEM/ODM global wholesale service."}`,
		"about":    `{"enabled": true, "title": "About Our Company", "content": "Founded in 2010, our company is a professional manufacturer and exporter engaged in the research, development, production and sales of high-standard medical devices."}`,
		"products": `{"enabled": true, "title": "Our Main Products", "show_count": 6}`,
	}

	for pageName, configData := range defaultConfigs {
		var config models.PageConfig
		DB.Where("page_name = ? AND tenant_id = ?", pageName, defaultTenantID).First(&config)
		if config.ID == 0 {
			DB.Create(&models.PageConfig{
				TenantID:   defaultTenantID,
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

// InitTenantConfigs 为所有租户初始化配置记录
func InitTenantConfigs() {
	var tenants []models.Tenant
	DB.Find(&tenants)
	
	for _, tenant := range tenants {
		var config models.TenantConfig
		DB.Where("tenant_id = ?", tenant.ID).First(&config)
		if config.ID == 0 {
			// 创建默认配置
			defaultConfig := &models.TenantConfig{
				TenantID: tenant.ID,
				FeatureFlags: `{"image_management":true,"page_config":true,"multi_language":true,"contact_form":true,"content_management":true}`,
				ResourceQuota: `{"max_images":50,"max_storage_mb":512,"max_content_items":20,"max_users":3}`,
				ResourceUsage: `{"used_images":0,"used_storage_mb":0,"used_content_items":0,"used_users":0}`,
				SubscriptionPlan: tenant.SubscriptionPlan,
				SubscriptionExpiresAt: tenant.SubscriptionExpiresAt,
			}
			DB.Create(defaultConfig)
		}
	}
}

func initDefaultSiteSettings() {
	// 获取默认租户
	var defaultTenant models.Tenant
	DB.Where("tenant_code = ?", "default").First(&defaultTenant)
	
	defaultTenantID := uint(0)
	if defaultTenant.ID > 0 {
		defaultTenantID = defaultTenant.ID
	}

	var count int64
	DB.Model(&models.SiteSetting{}).Where("key = ? AND tenant_id = ?", "zh_site_title", defaultTenantID).Count(&count)
	if count == 0 {
		DB.Create(&models.SiteSetting{
			TenantID: defaultTenantID,
			Key:      "zh_site_title",
			Value:    "专业医疗器械制造商与出口商",
		})
	}

	DB.Model(&models.SiteSetting{}).Where("key = ? AND tenant_id = ?", "en_site_title", defaultTenantID).Count(&count)
	if count == 0 {
		DB.Create(&models.SiteSetting{
			TenantID: defaultTenantID,
			Key:      "en_site_title",
			Value:    "Professional Medical Device Manufacturer & Exporter",
		})
	}

	// 初始化默认 Logo 配置
	DB.Model(&models.SiteSetting{}).Where("key = ? AND tenant_id = ?", "zh_site_logo", defaultTenantID).Count(&count)
	if count == 0 {
		DB.Create(&models.SiteSetting{
			TenantID: defaultTenantID,
			Key:      "zh_site_logo",
			Value:    "医疗",
		})
	}

	DB.Model(&models.SiteSetting{}).Where("key = ? AND tenant_id = ?", "en_site_logo", defaultTenantID).Count(&count)
	if count == 0 {
		DB.Create(&models.SiteSetting{
			TenantID: defaultTenantID,
			Key:      "en_site_logo",
			Value:    "MEDICAL",
		})
	}

	DB.Model(&models.SiteSetting{}).Where("key = ? AND tenant_id = ?", "site_logo_color", defaultTenantID).Count(&count)
	if count == 0 {
		DB.Create(&models.SiteSetting{
			TenantID: defaultTenantID,
			Key:      "site_logo_color",
			Value:    "#06a499",
		})
	}
}