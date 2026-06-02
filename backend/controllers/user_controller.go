package controllers

import (
	"github.com/gin-gonic/gin"
	"medical-device-cms/backend/common"
	"medical-device-cms/backend/services"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户通过用户名和密码登录，返回 JWT Token 和用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录请求"
// @Success 200 {object} common.APIResponse{data=LoginResponse} "登录成功"
// @Failure 400 {object} common.APIResponse "请求参数错误"
// @Failure 401 {object} common.APIResponse "用户名或密码错误"
// @Router /api/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.JSONBadRequest(ctx, "Invalid request parameters")
		return
	}

	user, token, err := c.userService.Login(req.Username, req.Password)
	if err != nil {
		common.JSONUnauthorized(ctx, "Invalid credentials")
		return
	}

	data := LoginResponse{
		Token: token,
		User: UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
		},
	}
	common.JSONSuccess(ctx, data)
}

// LoginRequest 登录请求结构
// @Description 用户登录请求参数
type LoginRequest struct {
	// 用户名
	Username string `json:"username" binding:"required" example:"admin"`
	// 密码
	Password string `json:"password" binding:"required" example:"123456"`
}

// LoginResponse 登录响应结构
// @Description 用户登录成功后的响应数据
type LoginResponse struct {
	// JWT Token
	Token string `json:"token"`
	// 用户信息
	User UserInfo `json:"user"`
}

// UserInfo 用户信息结构
// @Description 登录返回的用户基本信息
type UserInfo struct {
	// 用户 ID
	ID uint `json:"id"`
	// 用户名
	Username string `json:"username"`
	// 用户角色
	Role string `json:"role"`
}