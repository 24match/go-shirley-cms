package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	texts, err := c.languageService.GetPublicLanguageTexts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, texts)
}

func (c *LanguageController) GetLanguageTexts(ctx *gin.Context) {
	module := ctx.Query("module")
	texts, err := c.languageService.GetLanguageTexts(module)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, texts)
}

func (c *LanguageController) GetLanguageText(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	text, err := c.languageService.GetLanguageText(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Language text not found"})
		return
	}
	ctx.JSON(http.StatusOK, text)
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	text, err := c.languageService.CreateLanguageText(req.Key, req.Module, req.EnText, req.ZhText, req.Description)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Key already exists"})
		return
	}
	ctx.JSON(http.StatusCreated, text)
}

func (c *LanguageController) UpdateLanguageText(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	var req struct {
		EnText      string `json:"enText"`
		ZhText      string `json:"zhText"`
		Description string `json:"description"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	text, err := c.languageService.UpdateLanguageText(uint(id), req.EnText, req.ZhText, req.Description)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Language text not found"})
		return
	}
	ctx.JSON(http.StatusOK, text)
}

func (c *LanguageController) DeleteLanguageText(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	err := c.languageService.DeleteLanguageText(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Language text not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Language text deleted successfully"})
}

func (c *LanguageController) GetLanguageTextVersions(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	versions, err := c.languageService.GetLanguageTextVersions(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, versions)
}

func (c *LanguageController) RestoreLanguageTextVersion(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	version, _ := strconv.ParseUint(ctx.Param("version"), 10, 32)

	text, err := c.languageService.RestoreLanguageTextVersion(uint(id), uint(version))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}
	ctx.JSON(http.StatusOK, text)
}