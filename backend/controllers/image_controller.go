package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"medical-device-cms/backend/common"
	"medical-device-cms/backend/services"

	"github.com/gin-gonic/gin"
)

type ImageController struct {
	imageService *services.ImageService
}

func NewImageController() *ImageController {
	return &ImageController{
		imageService: services.NewImageService(),
	}
}

// GetImages 获取图片列表
// @Summary 获取图片列表
// @Description 获取所有或指定分类的图片列表
// @Tags 图片管理
// @Accept json
// @Produce json
// @Param category query string false "图片分类，不传则返回所有图片"
// @Success 200 {object} common.APIResponse{data=[]models.Image} "图片列表"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Security APIKey
// @Router /api/admin/images [get]
func (c *ImageController) GetImages(ctx *gin.Context) {
	category := ctx.Query("category")
	images, err := c.imageService.GetImages(category)
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, images)
}

// UploadImage 上传单张图片
// @Summary 上传单张图片
// @Description 上传单张图片到服务器并记录到数据库
// @Tags 图片管理
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "图片文件"
// @Param description formData string false "图片描述"
// @Param category formData string false "图片分类"
// @Param sort_order formData int false "排序顺序"
// @Success 200 {object} common.APIResponse{data=models.Image} "上传成功"
// @Failure 400 {object} common.APIResponse "未上传文件"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Security APIKey
// @Router /api/admin/images [post]
func (c *ImageController) UploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		common.JSONBadRequest(ctx, "No file uploaded")
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
		common.JSONInternalServerError(ctx, "Could not save file")
		return
	}

	image, err := c.imageService.UploadImage(filename, filepathStr, description, category, file.Size, sortOrder)
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, image)
}

// UploadMultipleImages 上传多张图片
// @Summary 上传多张图片
// @Description 批量上传多张图片到服务器
// @Tags 图片管理
// @Accept multipart/form-data
// @Produce json
// @Param images formData file true "图片文件数组"
// @Param category formData string false "图片分类"
// @Param descriptions formData []string false "图片描述数组"
// @Success 200 {object} common.APIResponse{data=[]models.Image} "上传成功"
// @Failure 400 {object} common.APIResponse "未上传文件"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Security APIKey
// @Router /api/admin/images/batch [post]
func (c *ImageController) UploadMultipleImages(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		common.JSONBadRequest(ctx, "Failed to parse multipart form")
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		common.JSONBadRequest(ctx, "No files uploaded")
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
			common.JSONInternalServerError(ctx, "Could not save file")
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
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, images)
}

// DeleteImage 删除图片
// @Summary 删除图片
// @Description 根据图片 ID 删除图片
// @Tags 图片管理
// @Accept json
// @Produce json
// @Param id path int true "图片 ID"
// @Success 200 {object} common.APIResponse "删除成功"
// @Failure 404 {object} common.APIResponse "图片不存在"
// @Security APIKey
// @Router /api/admin/images/{id} [delete]
func (c *ImageController) DeleteImage(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	err := c.imageService.DeleteImage(uint(id))
	if err != nil {
		common.JSONNotFound(ctx, "Image not found")
		return
	}
	common.JSONSuccessWithMessage(ctx, nil, "Image deleted successfully")
}

// DeleteImageByFilename 根据文件名删除图片
// @Summary 根据文件名删除图片
// @Description 根据图片文件名删除图片
// @Tags 图片管理
// @Accept json
// @Produce json
// @Param filename path string true "图片文件名"
// @Success 200 {object} common.APIResponse "删除成功"
// @Failure 404 {object} common.APIResponse "图片不存在"
// @Security APIKey
// @Router /api/admin/images/by-filename/{filename} [delete]
func (c *ImageController) DeleteImageByFilename(ctx *gin.Context) {
	filename := ctx.Param("filename")
	err := c.imageService.DeleteImageByFilename(filename)
	if err != nil {
		common.JSONNotFound(ctx, "Image not found")
		return
	}
	common.JSONSuccessWithMessage(ctx, nil, "Image deleted successfully")
}

// UpdateImage 更新图片信息
// @Summary 更新图片信息
// @Description 更新指定图片的描述、分类等信息，支持重新上传图片
// @Tags 图片管理
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "图片 ID"
// @Param image formData file false "新的图片文件"
// @Param description formData string false "图片描述"
// @Param long_description formData string false "图片详细描述"
// @Param category formData string false "图片分类"
// @Param sort_order formData int false "排序顺序"
// @Success 200 {object} common.APIResponse{data=models.Image} "更新成功"
// @Failure 404 {object} common.APIResponse "图片不存在"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Security APIKey
// @Router /api/admin/images/{id} [put]
func (c *ImageController) UpdateImage(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)

	updates := map[string]interface{}{
		"contentType": ctx.GetHeader("Content-Type"),
	}

	contentType := ctx.GetHeader("Content-Type")
	if contentType == "multipart/form-data" {
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
		updates["description"] = ctx.PostForm("description")
		updates["category"] = ctx.PostForm("category")
		updates["sort_order"] = ctx.PostForm("sort_order")
	} else {
		var req UpdateImageRequest
		if err := ctx.ShouldBindJSON(&req); err == nil {
			updates["description"] = req.Description
			updates["long_description"] = req.LongDescription
			updates["category"] = req.Category
			updates["sort_order"] = req.SortOrder
		}
	}

	image, err := c.imageService.UpdateImage(uint(id), updates)
	if err != nil {
		common.JSONNotFound(ctx, "Image not found")
		return
	}
	common.JSONSuccess(ctx, image)
}

// UpdateImageRequest 图片更新请求
// @Description 更新图片信息的请求参数
type UpdateImageRequest struct {
	// 图片描述
	Description string `json:"description" example:"图片描述"`
	// 图片详细描述
	LongDescription string `json:"long_description" example:"这是图片的详细描述"`
	// 图片分类
	Category string `json:"category" example:"banner"`
	// 排序顺序
	SortOrder int `json:"sort_order" example:"1"`
}