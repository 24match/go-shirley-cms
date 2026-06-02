package services

import (
	"encoding/json"
	"medical-device-cms/backend/config"
	"medical-device-cms/backend/models"
)

type SiteSettingService struct{}

func NewSiteSettingService() *SiteSettingService {
	return &SiteSettingService{}
}

type SiteSettings struct {
	ZhSiteTitle     string `json:"zhSiteTitle"`
	EnSiteTitle     string `json:"enSiteTitle"`
	ZhSiteLogo      string `json:"zhSiteLogo"`
	EnSiteLogo      string `json:"enSiteLogo"`
	SiteLogoColor   string `json:"siteLogoColor"`
}

func (s *SiteSettingService) GetSettings() (*SiteSettings, error) {
	settings := &SiteSettings{
		ZhSiteTitle:   "专业医疗器械制造商与出口商",
		EnSiteTitle:   "Professional Medical Device Manufacturer & Exporter",
		ZhSiteLogo:    "医疗",
		EnSiteLogo:    "MEDICAL",
		SiteLogoColor: "#06a499",
	}

	var zhSetting models.SiteSetting
	if err := config.DB.Where("key = ?", "zh_site_title").First(&zhSetting).Error; err == nil {
		settings.ZhSiteTitle = zhSetting.Value
	}

	var enSetting models.SiteSetting
	if err := config.DB.Where("key = ?", "en_site_title").First(&enSetting).Error; err == nil {
		settings.EnSiteTitle = enSetting.Value
	}

	var zhLogoSetting models.SiteSetting
	if err := config.DB.Where("key = ?", "zh_site_logo").First(&zhLogoSetting).Error; err == nil {
		settings.ZhSiteLogo = zhLogoSetting.Value
	}

	var enLogoSetting models.SiteSetting
	if err := config.DB.Where("key = ?", "en_site_logo").First(&enLogoSetting).Error; err == nil {
		settings.EnSiteLogo = enLogoSetting.Value
	}

	var logoColorSetting models.SiteSetting
	if err := config.DB.Where("key = ?", "site_logo_color").First(&logoColorSetting).Error; err == nil {
		settings.SiteLogoColor = logoColorSetting.Value
	}

	return settings, nil
}

func (s *SiteSettingService) SaveSettings(zhTitle, enTitle, zhLogo, enLogo, logoColor string) error {
	if err := s.saveSetting("zh_site_title", zhTitle); err != nil {
		return err
	}
	if err := s.saveSetting("en_site_title", enTitle); err != nil {
		return err
	}
	if err := s.saveSetting("zh_site_logo", zhLogo); err != nil {
		return err
	}
	if err := s.saveSetting("en_site_logo", enLogo); err != nil {
		return err
	}
	if err := s.saveSetting("site_logo_color", logoColor); err != nil {
		return err
	}
	return nil
}

func (s *SiteSettingService) saveSetting(key, value string) error {
	var setting models.SiteSetting
	if err := config.DB.Where("key = ?", key).First(&setting).Error; err != nil {
		setting = models.SiteSetting{
			Key:   key,
			Value: value,
		}
		return config.DB.Create(&setting).Error
	}

	setting.Value = value
	return config.DB.Save(&setting).Error
}

func (s *SiteSettingService) GetSettingsAsJSON() (string, error) {
	settings, err := s.GetSettings()
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(settings)
	if err != nil {
		return "", err
	}
	return string(data), nil
}