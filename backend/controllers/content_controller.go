package controllers

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"medical-device-cms/backend/common"
	"medical-device-cms/backend/services"
)

type ContentController struct {
	contentService *services.ContentService
}

func NewContentController() *ContentController {
	return &ContentController{
		contentService: services.NewContentService(),
	}
}

func (c *ContentController) GetContentItems(ctx *gin.Context) {
	section := ctx.Query("section")
	items, err := c.contentService.GetContentItems(section)
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, items)
}

func (c *ContentController) CreateContentItem(ctx *gin.Context) {
	section := ctx.PostForm("section")
	title := ctx.PostForm("title")
	description := ctx.PostForm("description")
	sortOrder, _ := strconv.Atoi(ctx.PostForm("sort_order"))
	zhTitle := ctx.PostForm("zhTitle")
	enTitle := ctx.PostForm("enTitle")
	zhDescription := ctx.PostForm("zhDescription")
	enDescription := ctx.PostForm("enDescription")

	var imagePath string
	file, err := ctx.FormFile("image")
	if err == nil {
		filename := strconv.Itoa(os.Getpid()) + "_" + file.Filename
		if err := ctx.SaveUploadedFile(file, "./uploads/"+filename); err != nil {
			common.JSONInternalServerError(ctx, "Failed to save image")
			return
		}
		imagePath = filename
	} else {
		imagePath = ctx.PostForm("image_path")
	}

	item, err := c.contentService.CreateContentItem(section, title, description, imagePath, sortOrder, "", zhTitle, enTitle, zhDescription, enDescription)
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, item)
}

func (c *ContentController) UpdateContentItem(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	updates := map[string]interface{}{}

	if section := ctx.PostForm("section"); section != "" {
		updates["section"] = section
	}
	if title := ctx.PostForm("title"); title != "" {
		updates["title"] = title
	}
	if description := ctx.PostForm("description"); description != "" {
		updates["description"] = description
	}
	// 多语言字段
	if zhTitle := ctx.PostForm("zhTitle"); zhTitle != "" {
		updates["zhTitle"] = zhTitle
	}
	if enTitle := ctx.PostForm("enTitle"); enTitle != "" {
		updates["enTitle"] = enTitle
	}
	if zhDescription := ctx.PostForm("zhDescription"); zhDescription != "" {
		updates["zhDescription"] = zhDescription
	}
	if enDescription := ctx.PostForm("enDescription"); enDescription != "" {
		updates["enDescription"] = enDescription
	}
	if sortOrderStr := ctx.PostForm("sort_order"); sortOrderStr != "" {
		updates["sort_order"] = sortOrderStr
	}

	file, err := ctx.FormFile("image")
	if err == nil {
		data, _ := os.ReadFile(file.Filename)
		updates["file"] = map[string]interface{}{
			"filename": file.Filename,
			"data":     data,
		}
	} else {
		if imagePath := ctx.PostForm("image_path"); imagePath != "" {
			updates["image_path"] = imagePath
		}
	}

	item, err := c.contentService.UpdateContentItem(uint(id), updates)
	if err != nil {
		common.JSONNotFound(ctx, "Item not found")
		return
	}
	common.JSONSuccess(ctx, item)
}

func (c *ContentController) DeleteContentItem(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	err := c.contentService.DeleteContentItem(uint(id))
	if err != nil {
		common.JSONNotFound(ctx, "Item not found")
		return
	}
	common.JSONSuccessWithMessage(ctx, nil, "Item deleted successfully")
}

func (c *ContentController) GetPublicContentItems(ctx *gin.Context) {
	section := ctx.Query("section")
	items, err := c.contentService.GetContentItems(section)
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, items)
}