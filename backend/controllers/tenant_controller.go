package controllers

import (
	"net/http"
	"strconv"

	"medical-device-cms/backend/middleware"
	"medical-device-cms/backend/services"
	"github.com/gin-gonic/gin"
)

// TenantController 租户管理控制器
type TenantController struct {
	userService         *services.UserService
	tenantConfigService *services.TenantConfigService
	quotaService        *services.QuotaService
}

// NewTenantController 创建租户管理控制器实例
func NewTenantController() *TenantController {
	return &TenantController{
		userService:         services.NewUserService(),
		tenantConfigService: services.NewTenantConfigService(),
		quotaService:        services.NewQuotaService(),
	}
}

// CreateUser 创建租户用户
// @Summary 创建租户用户
// @Tags 租户管理 - 用户管理
// @Accept json
// @Produce json
// @Param user body services.CreateUserRequest true "用户信息"
// @Success 200 {object} gin.H
// @Router /api/tenant/users [post]
func (ctrl *TenantController) CreateUser(c *gin.Context) {
	// 获取当前租户 ID
	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    "TENANT_REQUIRED",
			"message": "租户上下文不存在",
		})
		return
	}

	var req services.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_REQUEST",
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	user, err := ctrl.userService.CreateUser(c, tenantID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "CREATE_USER_FAILED",
			"message": "创建用户失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    "SUCCESS",
		"message": "用户创建成功",
		"data":    user,
	})
}

// ListUsers 获取租户用户列表
// @Summary 获取租户用户列表
// @Tags 租户管理 - 用户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} gin.H
// @Router /api/tenant/users [get]
func (ctrl *TenantController) ListUsers(c *gin.Context) {
	// 获取当前租户 ID
	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    "TENANT_REQUIRED",
			"message": "租户上下文不存在",
		})
		return
	}

	var filter services.UserFilter

	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)

	filter.Page = page
	filter.Size = size

	users, total, err := ctrl.userService.ListUsersByTenant(c, tenantID, &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "LIST_USERS_FAILED",
			"message": "获取用户列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "获取用户列表成功",
		"data": gin.H{
			"list":  users,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Tags 租户管理 - 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户 ID"
// @Success 200 {object} gin.H
// @Router /api/tenant/users/:id [get]
func (ctrl *TenantController) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_ID",
			"message": "无效的用户 ID",
		})
		return
	}

	user, err := ctrl.userService.GetUserByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "USER_NOT_FOUND",
			"message": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "获取用户成功",
		"data":    user,
	})
}

// UpdateUser 更新用户
// @Summary 更新用户信息
// @Tags 租户管理 - 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户 ID"
// @Param user body services.CreateUserRequest true "用户信息"
// @Success 200 {object} gin.H
// @Router /api/tenant/users/:id [put]
func (ctrl *TenantController) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_ID",
			"message": "无效的用户 ID",
		})
		return
	}

	var req services.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_REQUEST",
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	user, err := ctrl.userService.UpdateUser(c, uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "UPDATE_USER_FAILED",
			"message": "更新用户失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "更新用户成功",
		"data":    user,
	})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Tags 租户管理 - 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户 ID"
// @Success 200 {object} gin.H
// @Router /api/tenant/users/:id [delete]
func (ctrl *TenantController) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "INVALID_ID",
			"message": "无效的用户 ID",
		})
		return
	}

	if err := ctrl.userService.DeleteUser(c, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "DELETE_USER_FAILED",
			"message": "删除用户失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "删除用户成功",
	})
}

// GetDomainConfig 获取域名配置
// @Summary 获取租户域名配置
// @Tags 租户管理 - 域名管理
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Router /api/tenant/domain [get]
func (ctrl *TenantController) GetDomainConfig(c *gin.Context) {
	tenant, exists := middleware.GetTenantFromContext(c)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    "TENANT_REQUIRED",
			"message": "租户上下文不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "获取域名配置成功",
		"data": gin.H{
			"sub_domain":    tenant.SubDomain,
			"custom_domain": tenant.CustomDomain,
		},
	})
}

// UpdateDomainConfig 更新域名配置
// @Summary 更新租户域名配置
// @Tags 租户管理 - 域名管理
// @Accept json
// @Produce json
// @Param domain body gin.H true "域名配置"
// @Success 200 {object} gin.H
// @Router /api/tenant/domain [put]
func (ctrl *TenantController) UpdateDomainConfig(c *gin.Context) {
	// 此功能需要调用租户服务更新租户信息
	// 简化实现，实际项目中应添加更新逻辑
	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "域名配置更新功能待实现",
	})
}

// GetTenantConfig 获取当前租户配置
// @Summary 获取当前租户配置信息
// @Tags 租户管理 - 配置管理
// @Accept json
// @Produce json
// @Success 200 {object} services.TenantConfigDTO
// @Router /api/tenant/config [get]
func (ctrl *TenantController) GetTenantConfig(c *gin.Context) {
	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    "TENANT_REQUIRED",
			"message": "租户上下文不存在",
		})
		return
	}

	config, err := ctrl.tenantConfigService.GetTenantConfig(c, tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "CONFIG_NOT_FOUND",
			"message": "租户配置不存在",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "获取租户配置成功",
		"data":    config,
	})
}

// GetFeatures 获取当前租户功能模块状态
// @Summary 获取当前租户功能模块状态
// @Tags 租户管理 - 配置管理
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Router /api/tenant/features [get]
func (ctrl *TenantController) GetFeatures(c *gin.Context) {
	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    "TENANT_REQUIRED",
			"message": "租户上下文不存在",
		})
		return
	}

	features, err := ctrl.tenantConfigService.GetFeatures(c, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "GET_FEATURES_FAILED",
			"message": "获取功能模块状态失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "获取功能模块状态成功",
		"data":    features,
	})
}

// GetQuota 获取当前租户配额使用情况
// @Summary 获取当前租户配额使用情况
// @Tags 租户管理 - 配置管理
// @Accept json
// @Produce json
// @Success 200 {object} gin.H
// @Router /api/tenant/quota [get]
func (ctrl *TenantController) GetQuota(c *gin.Context) {
	tenantID, exists := middleware.GetTenantIDFromContext(c)
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{
			"code":    "TENANT_REQUIRED",
			"message": "租户上下文不存在",
		})
		return
	}

	quota, usage, err := ctrl.tenantConfigService.GetQuotaUsage(c, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "GET_QUOTA_FAILED",
			"message": "获取配额使用情况失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "SUCCESS",
		"message": "获取配额使用情况成功",
		"data": gin.H{
			"quota": quota,
			"usage": usage,
		},
	})
}