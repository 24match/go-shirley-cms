package controllers

import (
	"context"
	"net/http"
	"strconv"

	"medical-device-cms/backend/middleware"
	"medical-device-cms/backend/services"
	"github.com/gin-gonic/gin"
)

// SuperAdminController 超级管理员控制器
type SuperAdminController struct {
	tenantService *services.TenantService
	userService   *services.UserService
}

// NewSuperAdminController 创建超级管理员控制器实例
func NewSuperAdminController() *SuperAdminController {
	return &SuperAdminController{
		tenantService: services.NewTenantService(),
		userService:   services.NewUserService(),
	}
}

// CreateTenant 创建租户
// @Summary 创建新租户
// @Tags 超级管理员 - 租户管理
// @Accept json
// @Produce json
// @Param tenant body services.CreateTenantRequest true "租户信息"
// @Success 200 {object} models.Tenant
// @Router /api/superadmin/tenants [post]
func (ctrl *SuperAdminController) CreateTenant(c *gin.Context) {
	var req services.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_REQUEST",
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	tenant, err := ctrl.tenantService.CreateTenant(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "CREATE_TENANT_FAILED",
			"message": "创建租户失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    "SUCCESS",
		"message": "租户创建成功",
		"data":    tenant,
	})
}

// GetTenant 获取租户详情
// @Summary 获取租户详情
// @Tags 超级管理员 - 租户管理
// @Accept json
// @Produce json
// @Param id path int true "租户 ID"
// @Success 200 {object} models.Tenant
// @Router /api/superadmin/tenants/:id [get]
func (ctrl *SuperAdminController) GetTenant(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_ID",
			"message": "无效的租户 ID",
		})
		return
	}

	tenant, err := ctrl.tenantService.GetTenant(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "TENANT_NOT_FOUND",
			"message": "租户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "获取租户成功",
		"data":    tenant,
	})
}

// ListTenants 获取租户列表
// @Summary 获取租户列表
// @Tags 超级管理员 - 租户管理
// @Accept json
// @Produce json
// @Param status query string false "租户状态"
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} gin.H
// @Router /api/superadmin/tenants [get]
func (ctrl *SuperAdminController) ListTenants(c *gin.Context) {
	var filter services.TenantFilter

	status := c.Query("status")
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	filter.Status = status
	filter.Page = page
	filter.Size = size

	tenants, total, err := ctrl.tenantService.ListTenants(c, &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "LIST_TENANTS_FAILED",
			"message": "获取租户列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "获取租户列表成功",
		"data": gin.H{
			"list":  tenants,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// UpdateTenant 更新租户
// @Summary 更新租户信息
// @Tags 超级管理员 - 租户管理
// @Accept json
// @Produce json
// @Param id path int true "租户 ID"
// @Param tenant body services.UpdateTenantRequest true "租户信息"
// @Success 200 {object} models.Tenant
// @Router /api/superadmin/tenants/:id [put]
func (ctrl *SuperAdminController) UpdateTenant(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_ID",
			"message": "无效的租户 ID",
		})
		return
	}

	var req services.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_REQUEST",
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	tenant, err := ctrl.tenantService.UpdateTenant(c, uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "UPDATE_TENANT_FAILED",
			"message": "更新租户失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "更新租户成功",
		"data":    tenant,
	})
}

// DeleteTenant 删除租户
// @Summary 删除租户
// @Tags 超级管理员 - 租户管理
// @Accept json
// @Produce json
// @Param id path int true "租户 ID"
// @Success 200 {object} gin.H
// @Router /api/superadmin/tenants/:id [delete]
func (ctrl *SuperAdminController) DeleteTenant(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_ID",
			"message": "无效的租户 ID",
		})
		return
	}

	if err := ctrl.tenantService.DeleteTenant(c, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "DELETE_TENANT_FAILED",
			"message": "删除租户失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "删除租户成功",
	})
}

// ActivateTenant 激活租户
// @Summary 激活租户
// @Tags 超级管理员 - 租户管理
// @Accept json
// @Produce json
// @Param id path int true "租户 ID"
// @Success 200 {object} gin.H
// @Router /api/superadmin/tenants/:id/activate [post]
func (ctrl *SuperAdminController) ActivateTenant(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_ID",
			"message": "无效的租户 ID",
		})
		return
	}

	if err := ctrl.tenantService.ActivateTenant(c, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "ACTIVATE_TENANT_FAILED",
			"message": "激活租户失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "激活租户成功",
	})
}

// DisableTenant 禁用租户
// @Summary 禁用租户
// @Tags 超级管理员 - 租户管理
// @Accept json
// @Produce json
// @Param id path int true "租户 ID"
// @Success 200 {object} gin.H
// @Router /api/superadmin/tenants/:id/disable [post]
func (ctrl *SuperAdminController) DisableTenant(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_ID",
			"message": "无效的租户 ID",
		})
		return
	}

	if err := ctrl.tenantService.DisableTenant(c, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "DISABLE_TENANT_FAILED",
			"message": "禁用租户失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "禁用租户成功",
	})
}

// ImpersonateTenant 切换租户上下文
// @Summary 切换租户上下文（以租户身份操作）
// @Tags 超级管理员 - 租户管理
// @Accept json
// @Produce json
// @Param id path int true "租户 ID"
// @Success 200 {object} gin.H
// @Router /api/superadmin/tenants/:id/impersonate [post]
func (ctrl *SuperAdminController) ImpersonateTenant(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_ID",
			"message": "无效的租户 ID",
		})
		return
	}

	tenant, err := ctrl.tenantService.GetTenant(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "TENANT_NOT_FOUND",
			"message": "租户不存在",
		})
		return
	}

	// 设置租户上下文
	middleware.SetTenantContext(c, tenant)

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "切换租户上下文成功",
		"data": gin.H{
			"tenant_id":   tenant.ID,
			"tenant_code": tenant.TenantCode,
		},
	})
}

// GetSystemStats 获取系统统计
// @Summary 获取系统统计信息
// @Tags 超级管理员 - 系统管理
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Router /api/superadmin/stats [get]
func (ctrl *SuperAdminController) GetSystemStats(c *gin.Context) {
	// 这里可以添加更多统计逻辑
	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "获取系统统计成功",
		"data": gin.H{
			"message": "系统统计功能待实现",
		},
	})
}