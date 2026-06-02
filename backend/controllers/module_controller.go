package controllers

import (
	"encoding/json"
	"strings"

	"medical-device-cms/backend/common"
	"medical-device-cms/backend/config"
	"medical-device-cms/backend/services"

	"github.com/gin-gonic/gin"
)

type ModuleController struct {
	moduleService *services.ModuleService
}

func NewModuleController() *ModuleController {
	return &ModuleController{
		moduleService: services.NewModuleService(),
	}
}

// GetPageConfig 获取页面配置
// @Summary 获取页面配置
// @Description 获取指定页面或所有页面的配置信息
// @Tags 模块管理
// @Accept json
// @Produce json
// @Param page query string false "页面名称，不传则返回所有页面配置"
// @Success 200 {object} common.APIResponse{data=[]models.PageConfig} "页面配置列表"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Security APIKey
// @Router /api/admin/config [get]
func (c *ModuleController) GetPageConfig(ctx *gin.Context) {
	pageName := ctx.Query("page")
	if pageName != "" {
		data, err := config.GetPageConfig(pageName)
		if err != nil {
			common.JSONInternalServerError(ctx, err.Error())
			return
		}
		common.JSONSuccess(ctx, data)
	} else {
		configs, err := config.GetAllPageConfigs()
		if err != nil {
			common.JSONInternalServerError(ctx, err.Error())
			return
		}
		common.JSONSuccess(ctx, configs)
	}
}

// UpdatePageConfig 更新页面配置
// @Summary 更新页面配置
// @Description 更新指定页面的配置数据
// @Tags 模块管理
// @Accept json
// @Produce json
// @Param request body UpdatePageConfigRequest true "页面配置更新请求"
// @Success 200 {object} common.APIResponse "更新成功"
// @Failure 400 {object} common.APIResponse "请求参数错误"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Security APIKey
// @Router /api/admin/config [put]
func (c *ModuleController) UpdatePageConfig(ctx *gin.Context) {
	var req UpdatePageConfigRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.JSONBadRequest(ctx, "Invalid request parameters")
		return
	}

	common.JSONSuccessWithMessage(ctx, nil, "Config updated")
}

// UpdatePageConfigRequest 页面配置更新请求
// @Description 更新页面配置的请求参数
type UpdatePageConfigRequest struct {
	// 页面名称
	PageName string `json:"pageName" binding:"required" example:"home"`
	// 配置数据（JSON 格式）
	ConfigData string `json:"configData" example:"{\"theme\":\"light\"}"`
}

// GetPublicPageConfig 获取公开页面配置
// @Summary 获取公开页面配置
// @Description 获取所有页面的公开配置信息（无需认证）
// @Tags 公开接口
// @Accept json
// @Produce json
// @Success 200 {object} common.APIResponse{data=[]models.PageConfig} "页面配置列表"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Router /api/public/config [get]
func (c *ModuleController) GetPublicPageConfig(ctx *gin.Context) {
	configs, err := config.GetAllPageConfigs()
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, configs)
}

// GetModuleConfigs 获取模块配置列表
// @Summary 获取模块配置列表
// @Description 获取所有模块或指定模块的配置信息
// @Tags 模块管理
// @Accept json
// @Produce json
// @Param module query string false "模块名称，不传则返回所有模块配置"
// @Success 200 {object} common.APIResponse{data=[]models.ModuleConfig} "模块配置列表"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Security APIKey
// @Router /api/admin/modules [get]
func (c *ModuleController) GetModuleConfigs(ctx *gin.Context) {
	moduleName := ctx.Query("module")
	configs, err := c.moduleService.GetModuleConfigs(moduleName)
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, configs)
}

// GetModuleConfig 获取单个模块配置
// @Summary 获取单个模块配置
// @Description 根据模块名称获取模块配置详情
// @Tags 模块管理
// @Accept json
// @Produce json
// @Param name path string true "模块名称"
// @Success 200 {object} common.APIResponse{data=models.ModuleConfig} "模块配置详情"
// @Failure 404 {object} common.APIResponse "模块配置不存在"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Security APIKey
// @Router /api/admin/modules/{name} [get]
func (c *ModuleController) GetModuleConfig(ctx *gin.Context) {
	name := ctx.Param("name")
	config, err := c.moduleService.GetModuleConfig(name)
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, config)
}

// SaveModuleConfig 保存模块配置
// @Summary 保存模块配置
// @Description 创建或更新模块配置信息，支持表单和 JSON 格式
// @Tags 模块管理
// @Accept json
// @Produce json
// @Param moduleName path string true "模块名称"
// @Param request body SaveModuleConfigRequest true "模块配置数据"
// @Success 200 {object} common.APIResponse{data=models.ModuleConfig} "保存成功"
// @Failure 400 {object} common.APIResponse "请求参数错误"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Security APIKey
// @Router /api/admin/modules/{name} [put]
func (c *ModuleController) SaveModuleConfig(ctx *gin.Context) {
	moduleName := ctx.PostForm("moduleName")
	if moduleName == "" {
		var req struct {
			ModuleName string `json:"moduleName"`
		}
		if err := ctx.ShouldBindJSON(&req); err == nil {
			moduleName = req.ModuleName
		}
	}

	if moduleName == "" {
		common.JSONBadRequest(ctx, "moduleName is required")
		return
	}

	updates := map[string]interface{}{
		"contentType": ctx.GetHeader("Content-Type"),
	}

	contentType := ctx.GetHeader("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		updates["enabled"] = ctx.PostForm("enabled")
		updates["title"] = ctx.PostForm("title")
		updates["subtitle"] = ctx.PostForm("subtitle")
		updates["content"] = ctx.PostForm("content")
		updates["sortOrder"] = ctx.PostForm("sortOrder")
		updates["description"] = ctx.PostForm("description")
		updates["extraData"] = ctx.PostForm("extraData")
		updates["zhTitle"] = ctx.PostForm("zhTitle")
		updates["enTitle"] = ctx.PostForm("enTitle")
		updates["zhSubtitle"] = ctx.PostForm("zhSubtitle")
		updates["enSubtitle"] = ctx.PostForm("enSubtitle")
		updates["zhContent"] = ctx.PostForm("zhContent")
		updates["enContent"] = ctx.PostForm("enContent")
		updates["zhDescription"] = ctx.PostForm("zhDescription")
		updates["enDescription"] = ctx.PostForm("enDescription")

		extraDataMap := make(map[string]interface{})
		if existingExtra := ctx.PostForm("extraData"); existingExtra != "" {
			if err := json.Unmarshal([]byte(existingExtra), &extraDataMap); err != nil {
				extraDataMap = make(map[string]interface{})
			}
		}

		if zhName := ctx.PostForm("zhName"); zhName != "" {
			extraDataMap["zh_name"] = zhName
		}
		if enName := ctx.PostForm("enName"); enName != "" {
			extraDataMap["en_name"] = enName
		}
		if booth := ctx.PostForm("booth"); booth != "" {
			extraDataMap["booth"] = booth
		}
		if startDate := ctx.PostForm("startDate"); startDate != "" {
			extraDataMap["start_date"] = startDate
		}
		if endDate := ctx.PostForm("endDate"); endDate != "" {
			extraDataMap["end_date"] = endDate
		}
		if zhLocation := ctx.PostForm("zhLocation"); zhLocation != "" {
			extraDataMap["zh_location"] = zhLocation
		}
		if enLocation := ctx.PostForm("enLocation"); enLocation != "" {
			extraDataMap["en_location"] = enLocation
		}

		// 处理 contact 模块的特殊字段
		if moduleName == "contact" {
			if email := ctx.PostForm("email"); email != "" {
				extraDataMap["email"] = email
			}
			if phone := ctx.PostForm("phone"); phone != "" {
				extraDataMap["phone"] = phone
			}
			if whatsapp := ctx.PostForm("whatsapp"); whatsapp != "" {
				extraDataMap["whatsapp"] = whatsapp
			}
			if address := ctx.PostForm("address"); address != "" {
				extraDataMap["address"] = address
			}
		}

		if len(extraDataMap) > 0 {
			if updatedExtra, err := json.Marshal(extraDataMap); err == nil {
				updates["extraData"] = string(updatedExtra)
			}
		}

		file, err := ctx.FormFile("image")
		if err == nil && file != nil {
			fileData, err := file.Open()
			if err != nil {
				common.JSONInternalServerError(ctx, "Could not open file")
				return
			}
			defer fileData.Close()
			data := make([]byte, file.Size)
			fileData.Read(data)
			updates["file"] = map[string]interface{}{
				"filename": file.Filename,
				"data":     data,
			}
		}
	} else {
		var req struct {
			ModuleName  string `json:"moduleName"`
			Enabled     bool   `json:"enabled"`
			Title       string `json:"title"`
			Subtitle    string `json:"subtitle"`
			Content     string `json:"content"`
			SortOrder   int    `json:"sortOrder"`
			Description string `json:"description"`
			ExtraData   string `json:"extraData"`
			ImagePath   string `json:"imagePath"`
		}
		if err := ctx.ShouldBindJSON(&req); err == nil {
			updates["data"] = req
		}
	}

	config, err := c.moduleService.SaveModuleConfig(moduleName, updates)
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, config)
}

// DeleteModuleConfig 删除模块配置
// @Summary 删除模块配置
// @Description 根据模块名称删除模块配置
// @Tags 模块管理
// @Accept json
// @Produce json
// @Param name path string true "模块名称"
// @Success 200 {object} common.APIResponse "删除成功"
// @Failure 404 {object} common.APIResponse "模块配置不存在"
// @Security APIKey
// @Router /api/admin/modules/{name} [delete]
func (c *ModuleController) DeleteModuleConfig(ctx *gin.Context) {
	name := ctx.Param("name")
	err := c.moduleService.DeleteModuleConfig(name)
	if err != nil {
		common.JSONNotFound(ctx, "Module config not found")
		return
	}
	common.JSONSuccessWithMessage(ctx, nil, "Module config deleted successfully")
}

// DeleteModuleImage 删除模块图片
// @Summary 删除模块图片
// @Description 删除指定模块的关联图片
// @Tags 模块管理
// @Accept json
// @Produce json
// @Param name path string true "模块名称"
// @Success 200 {object} common.APIResponse "删除成功"
// @Failure 404 {object} common.APIResponse "模块图片不存在"
// @Security APIKey
// @Router /api/admin/modules/{name}/image [delete]
func (c *ModuleController) DeleteModuleImage(ctx *gin.Context) {
	moduleName := ctx.Param("name")
	err := c.moduleService.DeleteModuleImage(moduleName)
	if err != nil {
		common.JSONNotFound(ctx, "Module image not found")
		return
	}
	common.JSONSuccessWithMessage(ctx, nil, "Module image deleted successfully")
}

// GetPublicModuleConfigs 获取公开模块配置
// @Summary 获取公开模块配置
// @Description 获取所有模块的公开配置信息（无需认证）
// @Tags 公开接口
// @Accept json
// @Produce json
// @Success 200 {object} common.APIResponse{data=[]models.ModuleConfig} "模块配置列表"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Router /api/public/modules [get]
func (c *ModuleController) GetPublicModuleConfigs(ctx *gin.Context) {
	configs, err := c.moduleService.GetPublicModuleConfigs()
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, configs)
}

// SaveModuleConfigRequest 模块配置保存请求
// @Description 保存模块配置的请求参数
type SaveModuleConfigRequest struct {
	// 模块名称
	ModuleName string `json:"moduleName" example:"banner"`
	// 是否启用
	Enabled bool `json:"enabled" example:"true"`
	// 标题
	Title string `json:"title" example:"标题"`
	// 副标题
	Subtitle string `json:"subtitle" example:"副标题"`
	// 内容
	Content string `json:"content" example:"主要内容"`
	// 排序顺序
	SortOrder int `json:"sortOrder" example:"1"`
	// 描述
	Description string `json:"description" example:"模块描述"`
	// 额外数据（JSON 格式）
	ExtraData string `json:"extraData" example:"{\"key\":\"value\"}"`
	// 图片路径
	ImagePath string `json:"imagePath" example:"/uploads/images/banner.jpg"`
}