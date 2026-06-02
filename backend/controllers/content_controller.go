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

// GetContentItems 获取内容项列表
// @Summary 获取内容项列表
// @Description 获取所有或指定区域的内容项列表
// @Tags 内容管理
// @Accept json
// @Produce json
// @Param section query string false "区域名称，不传则返回所有内容项"
// @Success 200 {object} common.APIResponse{data=[]models.ContentItem} "内容项列表"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Security APIKey
// @Router /api/admin/content [get]
func (c *ContentController) GetContentItems(ctx *gin.Context) {
	section := ctx.Query("section")
	items, err := c.contentService.GetContentItems(section)
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, items)
}

// CreateContentItem 创建内容项
// @Summary 创建内容项
// @Description 创建新的内容项，支持图片上传和多语言
// @Tags 内容管理
// @Accept multipart/form-data
// @Produce json
// @Param section formData string true "区域名称"
// @Param zhTitle formData string false "中文标题"
// @Param enTitle formData string false "英文标题"
// @Param zhDescription formData string false "中文描述"
// @Param enDescription formData string false "英文描述"
// @Param image formData file false "图片文件"
// @Param sort_order formData int false "排序顺序"
// @Success 200 {object} common.APIResponse{data=models.ContentItem} "创建成功"
// @Failure 400 {object} common.APIResponse "请求参数错误"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Security APIKey
// @Router /api/admin/content [post]
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

// UpdateContentItem 更新内容项
// @Summary 更新内容项
// @Description 更新指定内容项的信息，支持图片上传和多语言
// @Tags 内容管理
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "内容项 ID"
// @Param section formData string false "区域名称"
// @Param zhTitle formData string false "中文标题"
// @Param enTitle formData string false "英文标题"
// @Param zhDescription formData string false "中文描述"
// @Param enDescription formData string false "英文描述"
// @Param image formData file false "图片文件"
// @Param sort_order formData int false "排序顺序"
// @Success 200 {object} common.APIResponse{data=models.ContentItem} "更新成功"
// @Failure 404 {object} common.APIResponse "内容项不存在"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Security APIKey
// @Router /api/admin/content/{id} [put]
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

// DeleteContentItem 删除内容项
// @Summary 删除内容项
// @Description 删除指定的内容项
// @Tags 内容管理
// @Accept json
// @Produce json
// @Param id path int true "内容项 ID"
// @Success 200 {object} common.APIResponse "删除成功"
// @Failure 404 {object} common.APIResponse "内容项不存在"
// @Security APIKey
// @Router /api/admin/content/{id} [delete]
func (c *ContentController) DeleteContentItem(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 32)
	err := c.contentService.DeleteContentItem(uint(id))
	if err != nil {
		common.JSONNotFound(ctx, "Item not found")
		return
	}
	common.JSONSuccessWithMessage(ctx, nil, "Item deleted successfully")
}

// GetPublicContentItems 获取公开内容项
// @Summary 获取公开内容项
// @Description 获取所有或指定区域的公开内容项（无需认证）
// @Tags 公开接口
// @Accept json
// @Produce json
// @Param section query string false "区域名称，不传则返回所有内容项"
// @Success 200 {object} common.APIResponse{data=[]models.ContentItem} "内容项列表"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Router /api/public/content [get]
func (c *ContentController) GetPublicContentItems(ctx *gin.Context) {
	section := ctx.Query("section")
	items, err := c.contentService.GetContentItems(section)
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}
	common.JSONSuccess(ctx, items)
}