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

func (c *UserController) Login(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.JSONBadRequest(ctx, "Invalid request parameters")
		return
	}

	user, token, err := c.userService.Login(req.Username, req.Password)
	if err != nil {
		common.JSONUnauthorized(ctx, "Invalid credentials")
		return
	}

	data := gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	}
	common.JSONSuccess(ctx, data)
}