package services

import (
	"testing"

	"medical-device-cms/backend/config"
	"medical-device-cms/backend/models"

	"github.com/stretchr/testify/assert"
)

func setupTestDB() {
	// 使用内存数据库或测试数据库
	config.InitDB()
}

func TestMergeExtraData_ContactFields(t *testing.T) {
	setupTestDB()

	mc := &models.ModuleConfig{
		ModuleName: "contact",
		Enabled:    true,
		ExtraData:  `{"email":"test@example.com","phone":"+86 123 4567 8900","whatsapp":"+1 234 567 8900","address":"No. 123, Industrial Park"}`,
	}

	result := mergeExtraData(mc)

	// 验证基础字段
	assert.Equal(t, "contact", result["moduleName"])
	assert.Equal(t, true, result["enabled"])

	// 验证 extraData 中的字段被正确展开到顶层
	assert.Equal(t, "test@example.com", result["email"])
	assert.Equal(t, "+86 123 4567 8900", result["phone"])
	assert.Equal(t, "+1 234 567 8900", result["whatsapp"])
	assert.Equal(t, "No. 123, Industrial Park", result["address"])
}

func TestMergeExtraData_EmptyExtraData(t *testing.T) {
	setupTestDB()

	mc := &models.ModuleConfig{
		ModuleName: "contact",
		Enabled:    true,
		ExtraData:  "",
	}

	result := mergeExtraData(mc)

	// 验证基础字段存在
	assert.Equal(t, "contact", result["moduleName"])
	assert.Equal(t, true, result["enabled"])

	// 验证没有额外的字段
	assert.Nil(t, result["email"])
	assert.Nil(t, result["phone"])
}

func TestMergeExtraData_InvalidJSON(t *testing.T) {
	setupTestDB()

	mc := &models.ModuleConfig{
		ModuleName: "contact",
		Enabled:    true,
		ExtraData:  `invalid json`,
	}

	result := mergeExtraData(mc)

	// 验证基础字段存在
	assert.Equal(t, "contact", result["moduleName"])

	// 验证解析失败时不会崩溃，没有额外字段
	assert.Nil(t, result["email"])
}

func TestGetPublicModuleConfigs_FilteringByEnabled(t *testing.T) {
	setupTestDB()

	// 清理测试数据
	defer func() {
		config.DB.Where("module_name IN ?", []string{"contact_test", "banner_test"}).Delete(&models.ModuleConfig{})
	}()

	// 创建测试数据
	config.DB.Create(&models.ModuleConfig{
		ModuleName: "contact_test",
		Enabled:    true,
		ExtraData:  `{"email":"enabled@example.com"}`,
	})
	config.DB.Create(&models.ModuleConfig{
		ModuleName: "banner_test",
		Enabled:    false,
		ExtraData:  `{"email":"disabled@example.com"}`,
	})

	service := NewModuleService()
	result, err := service.GetPublicModuleConfigs()

	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(result), 1)

	// 验证只返回 enabled=true 的模块
	foundContact := false
	foundBanner := false
	for _, r := range result {
		if r["moduleName"] == "contact_test" {
			foundContact = true
			assert.Equal(t, true, r["enabled"])
			assert.Equal(t, "enabled@example.com", r["email"])
		}
		if r["moduleName"] == "banner_test" {
			foundBanner = true
		}
	}

	assert.True(t, foundContact, "Should find enabled contact_test module")
	assert.False(t, foundBanner, "Should NOT find disabled banner_test module")
}