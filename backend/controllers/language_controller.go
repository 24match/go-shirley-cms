package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"medical-device-cms/backend/common"
	"medical-device-cms/backend/services"
)

type LanguageController struct {
	languageService *services.LanguageService
}

func NewLanguageController() *LanguageController {
	return &LanguageController{
		languageService: services.NewLanguageService(),
	}
}

// GetPublicLanguageTexts 获取公开多语言文本
// @Summary 获取公开多语言文本
// @Description 获取指定模块或所有模块的公开多语言文本（无需认证）
// @Tags 公开接口
// @Accept json
// @Produce json
// @Param module query string false "模块名称，不传则返回所有模块"
// @Success 200 {object} common.APIResponse{data=[]models.LanguageText} "多语言文本列表"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Router /api/public/lang [get]
func (c *LanguageController) GetPublicLanguageTexts(ctx *gin.Context) {
	module := ctx.Query("module")
	texts, err := c.languageService.GetPublicLanguageTextsByModule(module)
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, texts)
}

// GetLanguageTexts 获取多语言文本列表
// @Summary 获取多语言文本列表
// @Description 获取指定模块或所有模块的多语言文本列表
// @Tags 多语言管理
// @Accept json
// @Produce json
// @Param module query string false "模块名称，不传则返回所有模块"
// @Success 200 {object} common.APIResponse{data=[]models.LanguageText} "多语言文本列表"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Security APIKey
// @Router /api/admin/lang [get]
func (c *LanguageController) GetLanguageTexts(ctx *gin.Context) {
	module := ctx.Query("module")
	texts, err := c.languageService.GetLanguageTexts(module)
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, texts)
}

// GetLanguageText 获取单个多语言文本
// @Summary 获取单个多语言文本
// @Description 根据 ID 获取多语言文本详情
// @Tags 多语言管理
// @Accept json
// @Produce json
// @Param id path int true "多语言文本 ID"
// @Success 200 {object} common.APIResponse{data=models.LanguageText} "多语言文本详情"
// @Failure 404 {object} common.APIResponse "多语言文本不存在"
// @Security APIKey
// @Router /api/admin/lang/{id} [get]
func (c *LanguageController) GetLanguageText(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	text, err := c.languageService.GetLanguageText(uint(id))
	if err != nil {
		common.JSONNotFound(ctx, "Language text not found")
		return
	}
	common.JSONSuccess(ctx, text)
}

// CreateLanguageText 创建多语言文本
// @Summary 创建多语言文本
// @Description 创建新的多语言文本条目
// @Tags 多语言管理
// @Accept json
// @Produce json
// @Param request body CreateLanguageTextRequest true "多语言文本创建请求"
// @Success 200 {object} common.APIResponse{data=models.LanguageText} "创建成功"
// @Failure 400 {object} common.APIResponse "请求参数错误"
// @Failure 409 {object} common.APIResponse "键已存在"
// @Security APIKey
// @Router /api/admin/lang [post]
func (c *LanguageController) CreateLanguageText(ctx *gin.Context) {
	var req CreateLanguageTextRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.JSONBadRequest(ctx, "Invalid request parameters")
		return
	}

	text, err := c.languageService.CreateLanguageText(req.Key, req.Module, req.EnText, req.ZhText, req.Description)
	if err != nil {
		common.JSONError(ctx, 409, common.ErrConflict, "Key already exists")
		return
	}
	common.JSONSuccess(ctx, text)
}

// CreateLanguageTextRequest 创建多语言文本请求
// @Description 创建多语言文本的请求参数
type CreateLanguageTextRequest struct {
	// 文本键（唯一）
	Key string `json:"key" binding:"required" example:"home.welcome"`
	// 所属模块
	Module string `json:"module" binding:"required" example:"home"`
	// 英文文本
	EnText string `json:"enText" example:"Welcome"`
	// 中文文本
	ZhText string `json:"zhText" example:"欢迎"`
	// 描述说明
	Description string `json:"description" example:"首页欢迎语"`
}

// UpdateLanguageText 更新多语言文本
// @Summary 更新多语言文本
// @Description 更新指定多语言文本的内容
// @Tags 多语言管理
// @Accept json
// @Produce json
// @Param id path int true "多语言文本 ID"
// @Param request body UpdateLanguageTextRequest true "多语言文本更新请求"
// @Success 200 {object} common.APIResponse{data=models.LanguageText} "更新成功"
// @Failure 400 {object} common.APIResponse "请求参数错误"
// @Failure 404 {object} common.APIResponse "多语言文本不存在"
// @Security APIKey
// @Router /api/admin/lang/{id} [put]
func (c *LanguageController) UpdateLanguageText(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	var req UpdateLanguageTextRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.JSONBadRequest(ctx, "Invalid request parameters")
		return
	}

	text, err := c.languageService.UpdateLanguageText(uint(id), req.EnText, req.ZhText, req.Description)
	if err != nil {
		common.JSONNotFound(ctx, "Language text not found")
		return
	}
	common.JSONSuccess(ctx, text)
}

// UpdateLanguageTextRequest 更新多语言文本请求
// @Description 更新多语言文本的请求参数
type UpdateLanguageTextRequest struct {
	// 英文文本
	EnText string `json:"enText" example:"Welcome"`
	// 中文文本
	ZhText string `json:"zhText" example:"欢迎"`
	// 描述说明
	Description string `json:"description" example:"首页欢迎语"`
}

// DeleteLanguageText 删除多语言文本
// @Summary 删除多语言文本
// @Description 删除指定的多语言文本
// @Tags 多语言管理
// @Accept json
// @Produce json
// @Param id path int true "多语言文本 ID"
// @Success 200 {object} common.APIResponse "删除成功"
// @Failure 404 {object} common.APIResponse "多语言文本不存在"
// @Security APIKey
// @Router /api/admin/lang/{id} [delete]
func (c *LanguageController) DeleteLanguageText(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	err := c.languageService.DeleteLanguageText(uint(id))
	if err != nil {
		common.JSONNotFound(ctx, "Language text not found")
		return
	}
	common.JSONSuccessWithMessage(ctx, nil, "Language text deleted successfully")
}

// GetLanguageTextVersions 获取多语言文本版本历史
// @Summary 获取多语言文本版本历史
// @Description 获取指定多语言文本的所有历史版本
// @Tags 多语言管理
// @Accept json
// @Produce json
// @Param id path int true "多语言文本 ID"
// @Success 200 {object} common.APIResponse{data=[]models.LanguageTextVersion} "版本历史列表"
// @Failure 404 {object} common.APIResponse "多语言文本不存在"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Security APIKey
// @Router /api/admin/lang/{id}/versions [get]
func (c *LanguageController) GetLanguageTextVersions(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	versions, err := c.languageService.GetLanguageTextVersions(uint(id))
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, versions)
}

// RestoreLanguageTextVersion 恢复多语言文本版本
// @Summary 恢复多语言文本版本
// @Description 将指定多语言文本恢复到某个历史版本
// @Tags 多语言管理
// @Accept json
// @Produce json
// @Param id path int true "多语言文本 ID"
// @Param version path int true "版本号"
// @Success 200 {object} common.APIResponse{data=models.LanguageText} "恢复成功"
// @Failure 404 {object} common.APIResponse "版本不存在"
// @Security APIKey
// @Router /api/admin/lang/{id}/restore/{version} [post]
func (c *LanguageController) RestoreLanguageTextVersion(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	version, _ := strconv.ParseUint(ctx.Param("version"), 10, 32)

	text, err := c.languageService.RestoreLanguageTextVersion(uint(id), uint(version))
	if err != nil {
		common.JSONNotFound(ctx, "Version not found")
		return
	}
	common.JSONSuccess(ctx, text)
}