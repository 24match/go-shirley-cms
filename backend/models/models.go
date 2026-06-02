package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// Claims JWT Token 声明结构
// @Description 用于 JWT 认证的用户声明信息，包含用户 ID、角色和标准声明
type Claims struct {
	// 用户 ID
	UserID uint `json:"user_id" example:"1"`
	// 用户角色（admin 或 user）
	Role string `json:"role" example:"admin"`
	jwt.RegisteredClaims
}

// User 用户模型
// @Description 系统用户信息，用于后台管理认证
type User struct {
	// GORM 自动管理的 ID
	gorm.Model
	// 用户名（唯一）
	Username string `gorm:"uniqueIndex;not null" json:"username" example:"admin"`
	// 密码（加密存储）
	Password string `gorm:"not null" json:"password" example:"***"`
	// 用户角色
	Role string `gorm:"default:'admin'" json:"role" example:"admin"`
}

// Image 图片模型
// @Description 系统上传的图片资源信息
type Image struct {
	// 图片 ID
	ID uint `gorm:"primaryKey" json:"id" example:"1"`
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
// @Description 存储各页面的配置数据
type PageConfig struct {
	// GORM 自动管理的 ID
	gorm.Model
	// 页面名称（唯一）
	PageName string `gorm:"uniqueIndex;not null" json:"pageName" example:"home"`
	// 配置数据（JSON 格式）
	ConfigData string `gorm:"type:text" json:"configData" example:"{\"theme\":\"light\"}"`
}

// ModuleConfig 模块配置模型
// @Description 存储各功能模块的配置信息，支持多语言
type ModuleConfig struct {
	// GORM 自动管理的 ID
	gorm.Model
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
// @Description 存储页面内容项，支持多语言和图标
type ContentItem struct {
	// GORM 自动管理的 ID
	gorm.Model
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
// @Description 存储系统多语言文本，支持版本管理
type LanguageText struct {
	// GORM 自动管理的 ID
	gorm.Model
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
// @Description 存储多语言文本的历史版本，用于版本回溯
type LanguageTextVersion struct {
	// GORM 自动管理的 ID
	gorm.Model
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
// @Description 存储站点配置项，键值对形式
type SiteSetting struct {
	// GORM 自动管理的 ID
	gorm.Model
	// 设置键（唯一）
	Key string `gorm:"uniqueIndex;not null" json:"key" example:"site.name"`
	// 设置值
	Value string `gorm:"type:text" json:"value" example:"Shirley CMS"`
}

// ContactSubmission 联系表单提交模型
// @Description 存储用户通过联系我们表单提交的数据
type ContactSubmission struct {
	// GORM 自动管理的 ID
	gorm.Model
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
