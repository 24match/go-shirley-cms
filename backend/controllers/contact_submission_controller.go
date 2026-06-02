package controllers

import (
	"medical-device-cms/backend/common"
	"medical-device-cms/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContactSubmissionController struct {
	service *services.ContactSubmissionService
}

func NewContactSubmissionController() *ContactSubmissionController {
	return &ContactSubmissionController{
		service: services.NewContactSubmissionService(),
	}
}

// SubmitContactForm 提交联系表单
// @Summary 提交联系表单
// @Description 用户提交联系我们表单数据
// @Tags 联系表单
// @Accept json
// @Produce json
// @Param request body SubmitContactRequest true "联系表单数据"
// @Success 200 {object} common.APIResponse "提交成功"
// @Failure 400 {object} common.APIResponse "请求参数错误"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Router /api/contact/submit [post]
func (c *ContactSubmissionController) SubmitContactForm(ctx *gin.Context) {
	var req SubmitContactRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.JSONBadRequest(ctx, "Invalid request parameters")
		return
	}

	// 验证必填字段
	if req.Name == "" {
		common.JSONBadRequest(ctx, "Name is required")
		return
	}
	if req.Email == "" {
		common.JSONBadRequest(ctx, "Email is required")
		return
	}

	// 获取客户端 IP
	ipAddress := ctx.ClientIP()

	// 创建提交记录
	_, err := c.service.CreateSubmission(req.Name, req.Email, req.Company, req.Inquiry, ipAddress)
	if err != nil {
		common.JSONInternalServerError(ctx, "Failed to submit form")
		return
	}

	common.JSONSuccessWithMessage(ctx, nil, "Form submitted successfully")
}

// GetContactSubmissions 获取联系表单提交列表
// @Summary 获取联系表单提交列表
// @Description 分页获取用户提交的联系表单数据
// @Tags 联系表单
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} common.APIResponse{data=PaginationResponse} "提交列表"
// @Failure 500 {object} common.APIResponse "服务器错误"
// @Security APIKey
// @Router /api/admin/contact-submissions [get]
func (c *ContactSubmissionController) GetContactSubmissions(ctx *gin.Context) {
	page := 1
	pageSize := 10

	if p := ctx.Query("page"); p != "" {
		if _, err := strconv.Atoi(p); err == nil {
			page, _ = strconv.Atoi(p)
		}
	}
	if ps := ctx.Query("page_size"); ps != "" {
		if _, err := strconv.Atoi(ps); err == nil {
			pageSize, _ = strconv.Atoi(ps)
		}
	}

	submissions, total, err := c.service.GetSubmissions(page, pageSize)
	if err != nil {
		common.JSONInternalServerError(ctx, err.Error())
		return
	}

	response := PaginationResponse{
		List:      submissions,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: (int(total) + pageSize - 1) / pageSize,
	}

	common.JSONSuccess(ctx, response)
}

// DeleteContactSubmission 删除联系表单提交记录
// @Summary 删除联系表单提交记录
// @Description 根据 ID 删除联系表单提交记录
// @Tags 联系表单
// @Accept json
// @Produce json
// @Param id path int true "提交记录 ID"
// @Success 200 {object} common.APIResponse "删除成功"
// @Failure 404 {object} common.APIResponse "记录不存在"
// @Security APIKey
// @Router /api/admin/contact-submissions/{id} [delete]
func (c *ContactSubmissionController) DeleteContactSubmission(ctx *gin.Context) {
	id := ctx.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		common.JSONBadRequest(ctx, "Invalid ID")
		return
	}

	err = c.service.DeleteSubmission(uint(idUint))
	if err != nil {
		common.JSONNotFound(ctx, "Record not found")
		return
	}

	common.JSONSuccessWithMessage(ctx, nil, "Record deleted successfully")
}

// SubmitContactRequest 联系表单提交请求
// @Description 用户提交联系表单的请求参数
type SubmitContactRequest struct {
	// 姓名（必填）
	Name string `json:"name" binding:"required" example:"John Doe"`
	// 邮箱（必填）
	Email string `json:"email" binding:"required,email" example:"john@example.com"`
	// 国家/公司名称
	Company string `json:"company" example:"ABC Company"`
	// 询盘内容
	Inquiry string `json:"inquiry" example:"I want to know more about your products"`
}

// PaginationResponse 分页响应
// @Description 分页查询响应数据结构
type PaginationResponse struct {
	// 数据列表
	List interface{} `json:"list"`
	// 总记录数
	Total int64 `json:"total"`
	// 当前页码
	Page int `json:"page"`
	// 每页数量
	PageSize int `json:"page_size"`
	// 总页数
	TotalPage int `json:"total_page"`
}