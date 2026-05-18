package controllers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

func (c *ContentController) CreateContentItem(ctx *gin.Context) {
	section := ctx.PostForm("section")
	title := ctx.PostForm("title")
	description := ctx.PostForm("description")
	sortOrder, _ := strconv.Atoi(ctx.PostForm("sort_order"))

	var imagePath string
	file, err := ctx.FormFile("image")
	if err == nil {
		filename := strconv.Itoa(os.Getpid()) + "_" + file.Filename
		if err := ctx.SaveUploadedFile(file, "./uploads/"+filename); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		imagePath = filename
	} else {
		imagePath = ctx.PostForm("image_path")
	}

	item, err := c.contentService.CreateContentItem(section, title, description, imagePath, sortOrder, "")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, item)
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
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *ContentController) DeleteContentItem(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	err := c.contentService.DeleteContentItem(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}

func (c *ContentController) GetPublicContentItems(ctx *gin.Context) {
	section := ctx.Query("section")
	items, err := c.contentService.GetContentItems(section)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, items)
}