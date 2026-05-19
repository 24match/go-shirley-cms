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

func (c *LanguageController) GetPublicLanguageTexts(ctx *gin.Context) {
	module := ctx.Query("module")
	texts, err := c.languageService.GetPublicLanguageTextsByModule(module)
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, texts)
}

func (c *LanguageController) GetLanguageTexts(ctx *gin.Context) {
	module := ctx.Query("module")
	texts, err := c.languageService.GetLanguageTexts(module)
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, texts)
}

func (c *LanguageController) GetLanguageText(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	text, err := c.languageService.GetLanguageText(uint(id))
	if err != nil {
		common.JSONNotFound(ctx, "Language text not found")
		return
	}
	common.JSONSuccess(ctx, text)
}

func (c *LanguageController) CreateLanguageText(ctx *gin.Context) {
	var req struct {
		Key         string `json:"key" binding:"required"`
		Module      string `json:"module" binding:"required"`
		EnText      string `json:"enText"`
		ZhText      string `json:"zhText"`
		Description string `json:"description"`
	}

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

func (c *LanguageController) UpdateLanguageText(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	var req struct {
		EnText      string `json:"enText"`
		ZhText      string `json:"zhText"`
		Description string `json:"description"`
	}

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

func (c *LanguageController) DeleteLanguageText(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	err := c.languageService.DeleteLanguageText(uint(id))
	if err != nil {
		common.JSONNotFound(ctx, "Language text not found")
		return
	}
	common.JSONSuccessWithMessage(ctx, nil, "Language text deleted successfully")
}

func (c *LanguageController) GetLanguageTextVersions(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	versions, err := c.languageService.GetLanguageTextVersions(uint(id))
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, versions)
}

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