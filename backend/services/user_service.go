package services

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"medical-device-cms/backend/config"
	"medical-device-cms/backend/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// Login 用户登录
// 验证用户名密码，检查租户状态，生成 JWT Token
func (s *UserService) Login(username, password string) (*models.User, string, error) {
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", err
	}

	// 检查租户状态（如果不是超级管理员）
	if user.Role != "superadmin" && user.TenantID != 0 {
		var tenant models.Tenant
		if err := config.DB.First(&tenant, user.TenantID).Error; err != nil {
			return nil, "", fmt.Errorf("tenant not found")
		}
		if tenant.Status != models.TenantStatusActive {
			return nil, "", fmt.Errorf("tenant is not active: %s", tenant.Status)
		}
	}

	// 准备 JWT Claims
	claims := &models.Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	// 如果是租户用户，添加租户信息到 Claims
	if user.TenantID != 0 {
		claims.TenantID = user.TenantID
		var tenant models.Tenant
		if err := config.DB.First(&tenant, user.TenantID).Error; err == nil {
			claims.TenantCode = tenant.TenantCode
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		return nil, "", err
	}

	return &user, tokenString, nil
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=tenant_admin user"`
}

// UserFilter 用户查询过滤器
type UserFilter struct {
	Page int `form:"page" binding:"min=1"`
	Size int `form:"size" binding:"min=1,max=100"`
}

// CreateUser 创建租户用户
func (s *UserService) CreateUser(ctx context.Context, tenantID uint, req *CreateUserRequest) (*models.User, error) {
	// 检查用户名是否已存在
	var existing models.User
	if err := config.DB.Where("username = ? AND tenant_id = ?", req.Username, tenantID).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("username already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := models.User{
		TenantID: tenantID,
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// ListUsersByTenant 查询租户用户列表
func (s *UserService) ListUsersByTenant(ctx context.Context, tenantID uint, filter *UserFilter) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	query := config.DB.Model(&models.User{}).Where("tenant_id = ?", tenantID)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	if filter != nil && filter.Page > 0 && filter.Size > 0 {
		offset := (filter.Page - 1) * filter.Size
		query = query.Offset(offset).Limit(filter.Size)
	}

	// 按创建时间倒序
	query = query.Order("created_at DESC")

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserByID 根据 ID 获取用户
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, id uint, req *CreateUserRequest) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	// 更新字段
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUser 删除用户（软删除）
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return err
	}

	// 使用 GORM 的软删除
	if err := config.DB.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}