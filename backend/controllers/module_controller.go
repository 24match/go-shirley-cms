package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"medical-device-cms/backend/config"
	"medical-device-cms/backend/services"
)

type ModuleController struct {
	moduleService *services.ModuleService
}

func NewModuleController() *ModuleController {
	return &ModuleController{
		moduleService: services.NewModuleService(),
	}
}

func (c *ModuleController) GetPageConfig(ctx *gin.Context) {
	pageName := ctx.Query("page")
	if pageName != "" {
		data, err := config.GetPageConfig(pageName)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, data)
	} else {
		configs, err := config.GetAllPageConfigs()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, configs)
	}
}

func (c *ModuleController) UpdatePageConfig(ctx *gin.Context) {
	var req struct {
		PageName   string `json:"pageName" binding:"required"`
		ConfigData string `json:"configData"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Config updated"})
}

func (c *ModuleController) GetPublicPageConfig(ctx *gin.Context) {
	configs, err := config.GetAllPageConfigs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, configs)
}

func (c *ModuleController) GetModuleConfigs(ctx *gin.Context) {
	moduleName := ctx.Query("module")
	configs, err := c.moduleService.GetModuleConfigs(moduleName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, configs)
}

func (c *ModuleController) GetModuleConfig(ctx *gin.Context) {
	name := ctx.Param("name")
	config, err := c.moduleService.GetModuleConfig(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, config)
}

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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "moduleName is required"})
		return
	}

	updates := map[string]interface{}{
		"contentType": ctx.GetHeader("Content-Type"),
	}

	contentType := ctx.GetHeader("Content-Type")
	if contentType == "multipart/form-data" {
		updates["enabled"] = ctx.PostForm("enabled")
		updates["title"] = ctx.PostForm("title")
		updates["subtitle"] = ctx.PostForm("subtitle")
		updates["content"] = ctx.PostForm("content")
		updates["sortOrder"] = ctx.PostForm("sortOrder")
		updates["description"] = ctx.PostForm("description")
		updates["extraData"] = ctx.PostForm("extraData")

		file, err := ctx.FormFile("image")
		if err == nil && file != nil {
			data, _ := os.ReadFile(file.Filename)
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, config)
}

func (c *ModuleController) DeleteModuleConfig(ctx *gin.Context) {
	name := ctx.Param("name")
	err := c.moduleService.DeleteModuleConfig(name)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Module config not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Module config deleted successfully"})
}

func (c *ModuleController) GetPublicModuleConfigs(ctx *gin.Context) {
	configs, err := c.moduleService.GetPublicModuleConfigs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, configs)
}