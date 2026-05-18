package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"medical-device-cms/backend/services"
)

type ImageController struct {
	imageService *services.ImageService
}

func NewImageController() *ImageController {
	return &ImageController{
		imageService: services.NewImageService(),
	}
}

func (c *ImageController) GetImages(ctx *gin.Context) {
	category := ctx.Query("category")
	images, err := c.imageService.GetImages(category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, images)
}

func (c *ImageController) UploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	description := ctx.PostForm("description")
	category := ctx.PostForm("category")
	sortOrder := 0
	fmt.Sscanf(ctx.PostForm("sort_order"), "%d", &sortOrder)

	os.MkdirAll("uploads", 0755)
	filename := fmt.Sprintf("%d_%s", os.Getpid(), file.Filename)
	filepathStr := filepath.Join("uploads", filename)

	if err := ctx.SaveUploadedFile(file, filepathStr); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
		return
	}

	image, err := c.imageService.UploadImage(filename, filepathStr, description, category, file.Size, sortOrder)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, image)
}

func (c *ImageController) UploadMultipleImages(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		return
	}

	category := ctx.PostForm("category")
	descriptions := form.Value["descriptions"]

	os.MkdirAll("uploads", 0755)
	var uploadedImages []services.ImageUpload

	for _, file := range files {
		filename := fmt.Sprintf("%d_%s", os.Getpid(), file.Filename)
		filepathStr := filepath.Join("uploads", filename)

		if err := ctx.SaveUploadedFile(file, filepathStr); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
			return
		}

		data, _ := os.ReadFile(filepathStr)
		uploadedImages = append(uploadedImages, services.ImageUpload{
			Filename:     filename,
			Data:         data,
			Size:         file.Size,
			Category:     category,
			Descriptions: descriptions,
		})
	}

	images, err := c.imageService.UploadMultipleImages(uploadedImages)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, images)
}

func (c *ImageController) DeleteImage(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	err := c.imageService.DeleteImage(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}

func (c *ImageController) UpdateImage(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	updates := map[string]interface{}{
		"contentType": ctx.GetHeader("Content-Type"),
	}

	contentType := ctx.GetHeader("Content-Type")
	if contentType == "multipart/form-data" {
		file, err := ctx.FormFile("image")
		if err == nil && file != nil {
			data, _ := os.ReadFile(file.Filename)
			updates["file"] = map[string]interface{}{
				"filename": file.Filename,
				"data":     data,
			}
		}
		updates["description"] = ctx.PostForm("description")
		updates["category"] = ctx.PostForm("category")
		updates["sort_order"] = ctx.PostForm("sort_order")
	} else {
		var req struct {
			Description     string `json:"description"`
			LongDescription string `json:"long_description"`
			Category        string `json:"category"`
			SortOrder       int    `json:"sort_order"`
		}
		if err := ctx.ShouldBindJSON(&req); err == nil {
			updates["description"] = req.Description
			updates["long_description"] = req.LongDescription
			updates["category"] = req.Category
			updates["sort_order"] = req.SortOrder
		}
	}

	image, err := c.imageService.UpdateImage(uint(id), updates)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}
	ctx.JSON(http.StatusOK, image)
}