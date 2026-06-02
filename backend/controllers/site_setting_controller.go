package controllers

import (
	"medical-device-cms/backend/common"
	"medical-device-cms/backend/services"

	"github.com/gin-gonic/gin"
)

type SiteSettingController struct {
	siteSettingService *services.SiteSettingService
}

func NewSiteSettingController() *SiteSettingController {
	return &SiteSettingController{
		siteSettingService: services.NewSiteSettingService(),
	}
}

// GetSettings 获取站点设置
// @Summary 获取站点设置
// @Description 获取所有站点配置项
// @Tags 站点设置
// @Accept json
// @Produce json
// @Success 200 {object} common.APIResponse{data=[]models.SiteSetting} "站点设置列表"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Security APIKey
// @Router /api/admin/site-settings [get]
func (c *SiteSettingController) GetSettings(ctx *gin.Context) {
	settings, err := c.siteSettingService.GetSettings()
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, settings)
}

// SaveSettings 保存站点设置
// @Summary 保存站点设置
// @Description 保存站点标题、Logo 名称等配置项
// @Tags 站点设置
// @Accept multipart/form-data
// @Produce json
// @Param zhSiteTitle formData string false "中文站点标题"
// @Param enSiteTitle formData string false "英文站点标题"
// @Param zhSiteLogo formData string false "中文站点 Logo 名称"
// @Param enSiteLogo formData string false "英文站点 Logo 名称"
// @Param siteLogoColor formData string false "Logo 后缀颜色"
// @Success 200 {object} common.APIResponse{data=[]models.SiteSetting} "保存成功"
// @Failure 400 {object} common.APIResponse "请求参数错误"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Security APIKey
// @Router /api/admin/site-settings [post]
func (c *SiteSettingController) SaveSettings(ctx *gin.Context) {
	zhTitle := ctx.PostForm("zhSiteTitle")
	enTitle := ctx.PostForm("enSiteTitle")
	zhLogo := ctx.PostForm("zhSiteLogo")
	enLogo := ctx.PostForm("enSiteLogo")
	logoColor := ctx.PostForm("siteLogoColor")

	if zhTitle == "" && enTitle == "" && zhLogo == "" && enLogo == "" {
		common.JSONBadRequest(ctx, "At least one field is required")
		return
	}

	err := c.siteSettingService.SaveSettings(zhTitle, enTitle, zhLogo, enLogo, logoColor)
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}

	settings, err := c.siteSettingService.GetSettings()
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, settings)
}

// GetPublicSettings 获取公开站点设置
// @Summary 获取公开站点设置
// @Description 获取站点配置的公开信息（无需认证）
// @Tags 公开接口
// @Accept json
// @Produce json
// @Success 200 {object} common.APIResponse{data=[]models.SiteSetting} "站点设置列表"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Router /api/public/site-settings [get]
func (c *SiteSettingController) GetPublicSettings(ctx *gin.Context) {
	settings, err := c.siteSettingService.GetSettings()
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, settings)
}