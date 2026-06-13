package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// 租户状态常量
const (
	// TenantStatusActive 活跃状态
	TenantStatusActive = "active"
	// TenantStatusDisabled 禁用状态
	TenantStatusDisabled = "disabled"
	// TenantStatusPending 待激活状态
	TenantStatusPending = "pending"
	// TenantStatusSuspended 暂停状态
	TenantStatusSuspended = "suspended"
)

// 订阅计划常量
const (
	// SubscriptionPlanFree 免费版
	SubscriptionPlanFree = "free"
	// SubscriptionPlanPro 专业版
	SubscriptionPlanPro = "pro"
	// SubscriptionPlanEnterprise 企业版
	SubscriptionPlanEnterprise = "enterprise"
)

// Tenant 租户模型
// @Description SaaS 系统中的租户（多租户）信息，每个租户拥有独立的数据空间
type Tenant struct {
	// GORM 自动管理的 ID
	gorm.Model
	// 租户代码（唯一标识符，用于域名和路径识别）
	TenantCode string `gorm:"uniqueIndex;not null" json:"tenant_code" example:"acme-corp"`
	// 租户名称（显示名称）
	Name string `gorm:"not null" json:"name" example:"Acme Corporation"`
	// 租户状态（active/disabled/pending/suspended）
	Status string `gorm:"default:'active'" json:"status" example:"active"`
	// 子域名（用于租户访问）
	SubDomain string `gorm:"index" json:"sub_domain" example:"acme"`
	// 自定义域名
	CustomDomain string `gorm:"index" json:"custom_domain" example:"www.acme.com"`
	// 订阅计划（free/pro/enterprise）
	SubscriptionPlan string `gorm:"default:'free'" json:"subscription_plan" example:"free"`
	// 订阅过期时间
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at,omitempty" example:"2025-12-31T23:59:59Z"`
}

// AuditLog 审计日志模型
// @Description 记录系统关键操作日志，用于审计和追踪
type AuditLog struct {
	// GORM 自动管理的 ID
	gorm.Model
	// 操作人用户 ID
	OperatorID uint `gorm:"not null" json:"operator_id" example:"1"`
	// 租户 ID
	TenantID uint `gorm:"index;not null" json:"tenant_id" example:"1"`
	// 操作类型（CREATE/UPDATE/DELETE/LOGIN/LOGOUT 等）
	Action string `gorm:"not null" json:"action" example:"CREATE"`
	// 资源类型（User/Image/ModuleConfig 等）
	ResourceType string `json:"resource_type" example:"Image"`
	// 资源 ID
	ResourceID uint `json:"resource_id" example:"123"`
	// 操作前的值（JSON 格式）
	BeforeValue string `gorm:"type:text" json:"before_value,omitempty"`
	// 操作后的值（JSON 格式）
	AfterValue string `gorm:"type:text" json:"after_value,omitempty"`
}

// Claims JWT Token 声明结构
// @Description 用于 JWT 认证的用户声明信息，包含用户 ID、角色、租户信息和标准声明
type Claims struct {
	// 用户 ID
	UserID uint `json:"user_id" example:"1"`
	// 用户角色（superadmin/tenant_admin/user）
	Role string `json:"role" example:"admin"`
	// 租户 ID（超级管理员时为 0）
	TenantID uint `json:"tenant_id,omitempty" example:"1"`
	// 租户代码（超级管理员时为空）
	TenantCode string `json:"tenant_code,omitempty" example:"acme-corp"`
	jwt.RegisteredClaims
}

// User 用户模型
// @Description 系统用户信息，用于后台管理认证，支持多租户
type User struct {
	// GORM 自动管理的 ID
	gorm.Model
	// 租户 ID（超级管理员时为 0）
	TenantID uint `gorm:"index;default:0" json:"tenant_id" example:"1"`
	// 用户名（唯一）
	Username string `gorm:"uniqueIndex;not null" json:"username" example:"admin"`
	// 密码（加密存储）
	Password string `gorm:"not null" json:"password" example:"***"`
	// 用户角色（superadmin/tenant_admin/user）
	Role string `gorm:"default:'tenant_admin'" json:"role" example:"tenant_admin"`
}

// Image 图片模型
// @Description 系统上传的图片资源信息，支持多租户隔离
type Image struct {
	// 图片 ID
	ID uint `gorm:"primaryKey" json:"id" example:"1"`
	// 租户 ID
	TenantID uint `gorm:"index;not null" json:"tenant_id" example:"1"`
	// 创建时间
	CreatedAt time.Time `json:"createdAt" example:"2024-01-01T00:00:00Z"`
	// 更新时间
	UpdatedAt time.Time `json:"updatedAt" example:"2024-01-01T00:00:00Z"`
	// 删除时间（软删除）
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	// 文件名
	Filename string `gorm:"not null" json:"filename" example:"banner.jpg"`
	// 文件路径
	FilePath string `gorm:"not null" json:"filePath" example:"/uploads/images/banner.jpg"`
	// 文件大小（字节）
	FileSize int64 `json:"fileSize" example:"102400"`
	// 图片描述
	Description string `json:"description" example:"首页横幅图片"`
	// 详细描述
	LongDescription string `json:"longDescription" example:"这是首页展示的横幅图片"`
	// 分类
	Category string `json:"category" example:"banner"`
	// 排序顺序
	SortOrder int `gorm:"default:0" json:"sortOrder" example:"1"`
}

// PageConfig 页面配置模型
// @Description 存储各页面的配置数据，支持多租户隔离
type PageConfig struct {
	// GORM 自动管理的 ID
	gorm.Model
	// 租户 ID
	TenantID uint `gorm:"index;not null" json:"tenant_id" example:"1"`
	// 页面名称（唯一）
	PageName string `gorm:"uniqueIndex;not null" json:"pageName" example:"home"`
	// 配置数据（JSON 格式）
	ConfigData string `gorm:"type:text" json:"configData" example:"{\"theme\":\"light\"}"`
}

// ModuleConfig 模块配置模型
// @Description 存储各功能模块的配置信息，支持多语言和多租户隔离
type ModuleConfig struct {
	// GORM 自动管理的 ID
	gorm.Model
	// 租户 ID
	TenantID uint `gorm:"index;not null" json:"tenant_id" example:"1"`
	// 模块名称（唯一）
	ModuleName string `gorm:"uniqueIndex;not null" json:"moduleName" example:"banner"`
	// 是否启用
	Enabled bool `gorm:"default:true" json:"enabled" example:"true"`
	// 中文标题
	ZhTitle string `json:"zhTitle" example:"欢迎来到我们的网站"`
	// 英文标题
	EnTitle string `json:"enTitle" example:"Welcome to our website"`
	// 中文副标题
	ZhSubtitle string `json:"zhSubtitle" example:"专业医疗器械供应商"`
	// 英文副标题
	EnSubtitle string `json:"enSubtitle" example:"Professional medical device supplier"`
	// 中文内容
	ZhContent string `gorm:"type:text" json:"zhContent" example:"这里是主要内容..."`
	// 英文内容
	EnContent string `gorm:"type:text" json:"enContent" example:"Main content here..."`
	// 标题（向后兼容）
	Title string `json:"title" example:"Title"`
	// 副标题（向后兼容）
	Subtitle string `json:"subtitle" example:"Subtitle"`
	// 内容（向后兼容）
	Content string `gorm:"type:text" json:"content" example:"Content..."`
	// 图片路径
	ImagePath string `json:"imagePath" example:"/uploads/images/module.jpg"`
	// 排序顺序
	SortOrder int `gorm:"default:0" json:"sortOrder" example:"1"`
	// 额外数据（JSON 格式）
	ExtraData string `gorm:"type:text" json:"extraData" example:"{\"key\":\"value\"}"`
	// 中文描述
	ZhDescription string `json:"zhDescription" example:"模块描述"`
	// 英文描述
	EnDescription string `json:"enDescription" example:"Module description"`
	// 描述（向后兼容）
	Description string `json:"description" example:"Description"`
}

// ContentItem 内容项模型
// @Description 存储页面内容项，支持多语言和图标，支持多租户隔离
type ContentItem struct {
	// GORM 自动管理的 ID
	gorm.Model
	// 租户 ID
	TenantID uint `gorm:"index;not null" json:"tenant_id" example:"1"`
	// 所属区域/板块
	Section string `gorm:"not null" json:"section" example:"advantages"`
	// 中文标题
	ZhTitle string `json:"zhTitle" example:"专业团队"`
	// 英文标题
	EnTitle string `json:"enTitle" example:"Professional Team"`
	// 中文描述
	ZhDescription string `json:"zhDescription" example:"我们拥有专业的技术团队"`
	// 英文描述
	EnDescription string `json:"enDescription" example:"We have a professional technical team"`
	// 标题（向后兼容）
	Title string `json:"title" example:"Title"`
	// 描述（向后兼容）
	Description string `json:"description" example:"Description"`
	// 图标类名
	Icon string `json:"icon" example:"fa-users"`
	// 图片路径
	ImagePath string `json:"imagePath" example:"/uploads/images/team.jpg"`
	// 排序顺序
	SortOrder int `gorm:"default:0" json:"sortOrder" example:"1"`
}

// LanguageText 多语言文本模型
// @Description 存储系统多语言文本，支持版本管理和多租户隔离
type LanguageText struct {
	// GORM 自动管理的 ID
	gorm.Model
	// 租户 ID
	TenantID uint `gorm:"index;not null" json:"tenant_id" example:"1"`
	// 文本键（唯一）
	Key string `gorm:"uniqueIndex;not null" json:"key" example:"home.welcome"`
	// 所属模块
	Module string `gorm:"not null" json:"module" example:"home"`
	// 英文文本
	EnText string `gorm:"type:text" json:"enText" example:"Welcome"`
	// 中文文本
	ZhText string `gorm:"type:text" json:"zhText" example:"欢迎"`
	// 描述说明
	Description string `json:"description" example:"首页欢迎语"`
	// 版本号
	Version int `gorm:"default:1" json:"version" example:"1"`
}

// LanguageTextVersion 多语言版本文本模型
// @Description 存储多语言文本的历史版本，用于版本回溯，支持多租户隔离
type LanguageTextVersion struct {
	// GORM 自动管理的 ID
	gorm.Model
	// 租户 ID
	TenantID uint `gorm:"index;not null" json:"tenant_id" example:"1"`
	// 关联的 LanguageText ID
	LanguageTextID uint `gorm:"not null" json:"languageTextId" example:"1"`
	// 文本键
	Key string `gorm:"not null" json:"key" example:"home.welcome"`
	// 所属模块
	Module string `gorm:"not null" json:"module" example:"home"`
	// 英文文本
	EnText string `gorm:"type:text" json:"enText" example:"Welcome"`
	// 中文文本
	ZhText string `gorm:"type:text" json:"zhText" example:"欢迎"`
	// 描述说明
	Description string `json:"description" example:"首页欢迎语"`
	// 版本号
	Version int `gorm:"not null" json:"version" example:"1"`
	// 更新时间
	UpdatedAt time.Time `json:"updatedAt" example:"2024-01-01T00:00:00Z"`
}

// SiteSetting 站点设置模型
// @Description 存储站点配置项，键值对形式，支持多租户隔离
type SiteSetting struct {
	// GORM 自动管理的 ID
	gorm.Model
	// 租户 ID
	TenantID uint `gorm:"index;not null" json:"tenant_id" example:"1"`
	// 设置键（唯一）
	Key string `gorm:"uniqueIndex;not null" json:"key" example:"site.name"`
	// 设置值
	Value string `gorm:"type:text" json:"value" example:"Shirley CMS"`
}

// ContactSubmission 联系表单提交模型
// @Description 存储用户通过联系我们表单提交的数据，支持多租户隔离
type ContactSubmission struct {
	// GORM 自动管理的 ID
	gorm.Model
	// 租户 ID
	TenantID uint `gorm:"index;not null" json:"tenant_id" example:"1"`
	// 提交人姓名
	Name string `gorm:"not null" json:"name" example:"John Doe"`
	// 提交人邮箱
	Email string `gorm:"not null" json:"email" example:"john@example.com"`
	// 国家/公司名称
	Company string `json:"company" example:"ABC Company"`
	// 询盘内容
	Inquiry string `gorm:"type:text" json:"inquiry" example:"I want to know more about your products"`
	// IP 地址
	IPAddress string `json:"ipAddress" example:"192.168.1.1"`
	// 是否已读
	IsRead bool `gorm:"default:false" json:"isRead" example:"false"`
}

// TenantConfig 租户配置模型
// @Description 存储租户级别的配置信息，包括功能模块开关、资源配额、订阅计划
// @ID TenantConfig
// @Tags Models
type TenantConfig struct {
	// GORM 自动管理的 ID
	gorm.Model
	// 租户 ID（唯一关联）
	TenantID uint `gorm:"uniqueIndex;not null" json:"tenant_id" example:"1"`
	// 功能模块配置（JSON 格式）
	// 结构：{"image_management": true, "page_config": true, "multi_language": false, ...}
	FeatureFlags string `gorm:"type:text" json:"feature_flags" example:"{\"image_management\":true,\"page_config\":true}"`
	// 资源配额配置（JSON 格式）
	// 结构：{"max_images": 100, "max_storage_mb": 1024, "max_content_items": 50, "max_users": 5}
	ResourceQuota string `gorm:"type:text" json:"resource_quota" example:"{\"max_images\":100,\"max_storage_mb\":1024}"`
	// 资源使用统计（JSON 格式）
	// 结构：{"used_images": 45, "used_storage_mb": 512, "used_content_items": 20, "used_users": 3}
	ResourceUsage string `gorm:"type:text" json:"resource_usage" example:"{\"used_images\":45,\"used_storage_mb\":512}"`
	// 订阅计划（free/pro/enterprise）
	SubscriptionPlan string `gorm:"default:'free'" json:"subscription_plan" example:"free"`
	// 订阅过期时间
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at,omitempty" example:"2025-12-31T23:59:59Z"`
	// 租户关联（可选，用于预加载）
	Tenant Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
}

// TenantConfigDTO 租户配置数据传输对象
// @Description 租户配置数据传输对象，用于 API 响应
// @ID TenantConfigDTO
// @Tags Models
type TenantConfigDTO struct {
	// 配置 ID
	ID uint `json:"id" example:"1"`
	// 租户 ID
	TenantID uint `json:"tenant_id" example:"1"`
	// 功能模块配置
	FeatureFlags map[string]bool `json:"feature_flags" example:"{\"image_management\":true,\"page_config\":true}"`
	// 资源配额配置
	ResourceQuota map[string]int `json:"resource_quota" example:"{\"max_images\":100,\"max_storage_mb\":1024}"`
	// 资源使用统计
	ResourceUsage map[string]int `json:"resource_usage" example:"{\"used_images\":45,\"used_storage_mb\":512}"`
	// 订阅计划
	SubscriptionPlan string `json:"subscription_plan" example:"free"`
	// 订阅过期时间
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at,omitempty"`
}

// UpdateTenantConfigRequest 更新租户配置请求
// @Description 更新租户配置请求体
// @ID UpdateTenantConfigRequest
// @Tags Models
type UpdateTenantConfigRequest struct {
	// 功能模块配置
	FeatureFlags map[string]bool `json:"feature_flags"`
	// 资源配额配置
	ResourceQuota map[string]int `json:"resource_quota"`
	// 订阅计划
	SubscriptionPlan string `json:"subscription_plan"`
	// 订阅过期时间
	SubscriptionExpiresAt *time.Time `json:"subscription_expires_at,omitempty"`
}

// GetFeatureFlags 解析功能模块配置为 map
func (tc *TenantConfig) GetFeatureFlags() map[string]bool {
	flags := make(map[string]bool)
	if tc.FeatureFlags == "" {
		return flags
	}
	// 简单 JSON 解析，实际项目中可使用 encoding/json
	// 这里使用占位实现，实际服务层会正确处理
	return flags
}

// GetResourceQuota 解析资源配额为 map
func (tc *TenantConfig) GetResourceQuota() map[string]int {
	quota := make(map[string]int)
	if tc.ResourceQuota == "" {
		return quota
	}
	return quota
}

// GetResourceUsage 解析资源使用统计为 map
func (tc *TenantConfig) GetResourceUsage() map[string]int {
	usage := make(map[string]int)
	if tc.ResourceUsage == "" {
		return usage
	}
	return usage
}