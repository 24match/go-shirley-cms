package routes

import (
	"net/http"

	"medical-device-cms/backend/controllers"
	"medical-device-cms/backend/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SuperAdminMiddleware 超级管理员认证中间件
func SuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 调用现有的 AuthMiddleware 进行基础认证
		middleware.AuthMiddleware()(c)
		if c.IsAborted() {
			return
		}

		// 从上下文中获取用户角色并检查是否为超级管理员
		role, exists := c.Get("role")
		if !exists || role != "superadmin" {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    "FORBIDDEN",
				"message": "需要超级管理员权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RegisterRoutes 注册所有路由
// @description 初始化 CORS、错误处理中间件，注册静态文件路由、API 路由和 Swagger 文档路由
func RegisterRoutes(r *gin.Engine) {
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorHandler())

	r.Static("/uploads", "./uploads")
	r.Static("/frontend", "./frontend")
	r.Static("/admin", "./admin")
	r.StaticFile("/", "./frontend/index.html")
	r.StaticFile("/i18n.js", "./i18n.js")

	// 注册 Swagger UI 路由（仅开发环境）
	// @description 访问 /swagger/index.html 查看交互式 API 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	registerSuperAdminRoutes(api)
	registerTenantRoutes(api)
	registerPublicRoutes(api)
	registerAdminRoutes(api)
}

func registerPublicRoutes(api *gin.RouterGroup) {
	api.POST("/login", controllers.NewUserController().Login)

	// 添加租户识别中间件到公开 API
	public := api.Group("/public")
	public.Use(middleware.TenantMiddleware())

	public.GET("/config", controllers.NewModuleController().GetPublicPageConfig)
	public.GET("/modules", controllers.NewModuleController().GetPublicModuleConfigs)
	public.GET("/images", controllers.NewImageController().GetImages)
	public.GET("/content", controllers.NewContentController().GetPublicContentItems)
	public.GET("/lang", controllers.NewLanguageController().GetPublicLanguageTexts)
	public.GET("/site-settings", controllers.NewSiteSettingController().GetPublicSettings)

	// 联系表单提交（公开接口）
	api.POST("/contact/submit", controllers.NewContactSubmissionController().SubmitContactForm)
}

func registerSuperAdminRoutes(api *gin.RouterGroup) {
	superadmin := api.Group("/superadmin")
	superadmin.Use(middleware.AuthMiddleware())
	superadmin.Use(SuperAdminMiddleware())

	// 租户管理
	superadmin.POST("/tenants", controllers.NewSuperAdminController().CreateTenant)
	superadmin.GET("/tenants", controllers.NewSuperAdminController().ListTenants)
	superadmin.GET("/tenants/:id", controllers.NewSuperAdminController().GetTenant)
	superadmin.PUT("/tenants/:id", controllers.NewSuperAdminController().UpdateTenant)
	superadmin.DELETE("/tenants/:id", controllers.NewSuperAdminController().DeleteTenant)
	superadmin.POST("/tenants/:id/activate", controllers.NewSuperAdminController().ActivateTenant)
	superadmin.POST("/tenants/:id/disable", controllers.NewSuperAdminController().DisableTenant)
	superadmin.POST("/tenants/:id/impersonate", controllers.NewSuperAdminController().ImpersonateTenant)

	// 租户配置管理
	superadmin.GET("/tenants/:id/config", controllers.NewSuperAdminController().GetTenantConfig)
	superadmin.PUT("/tenants/:id/config", controllers.NewSuperAdminController().UpdateTenantConfig)
	superadmin.POST("/tenants/:id/quota/reset", controllers.NewSuperAdminController().ResetQuota)
	superadmin.GET("/tenants/:id/quota/usage", controllers.NewSuperAdminController().GetQuotaUsage)

	// 系统统计
	superadmin.GET("/stats", controllers.NewSuperAdminController().GetSystemStats)
}

func registerTenantRoutes(api *gin.RouterGroup) {
	tenant := api.Group("/tenant")
	tenant.Use(middleware.AuthMiddleware())
	tenant.Use(middleware.TenantMiddleware())

	// 用户管理
	tenant.POST("/users", controllers.NewTenantController().CreateUser)
	tenant.GET("/users", controllers.NewTenantController().ListUsers)
	tenant.GET("/users/:id", controllers.NewTenantController().GetUser)
	tenant.PUT("/users/:id", controllers.NewTenantController().UpdateUser)
	tenant.DELETE("/users/:id", controllers.NewTenantController().DeleteUser)

	// 域名管理
	tenant.GET("/domain", controllers.NewTenantController().GetDomainConfig)
	tenant.PUT("/domain", controllers.NewTenantController().UpdateDomainConfig)

	// 租户配置管理
	tenant.GET("/config", controllers.NewTenantController().GetTenantConfig)
	tenant.GET("/features", controllers.NewTenantController().GetFeatures)
	tenant.GET("/quota", controllers.NewTenantController().GetQuota)
}

func registerAdminRoutes(api *gin.RouterGroup) {
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware())

	admin.GET("/images", controllers.NewImageController().GetImages)
	admin.POST("/images", controllers.NewImageController().UploadImage)
	admin.POST("/images/batch", controllers.NewImageController().UploadMultipleImages)
	admin.PUT("/images/:id", controllers.NewImageController().UpdateImage)
	admin.DELETE("/images/by-filename/:filename", controllers.NewImageController().DeleteImageByFilename)
	admin.DELETE("/images/:id", controllers.NewImageController().DeleteImage)

	admin.GET("/config", controllers.NewModuleController().GetPageConfig)
	admin.PUT("/config", controllers.NewModuleController().UpdatePageConfig)

	admin.GET("/modules", controllers.NewModuleController().GetModuleConfigs)
	admin.GET("/modules/:name", controllers.NewModuleController().GetModuleConfig)
	admin.POST("/modules", controllers.NewModuleController().SaveModuleConfig)
	admin.PUT("/modules/:name", controllers.NewModuleController().SaveModuleConfig)
	admin.DELETE("/modules/:name", controllers.NewModuleController().DeleteModuleConfig)
	admin.DELETE("/modules/:name/image", controllers.NewModuleController().DeleteModuleImage)

	admin.GET("/content", controllers.NewContentController().GetContentItems)
	admin.POST("/content", controllers.NewContentController().CreateContentItem)
	admin.PUT("/content/:id", controllers.NewContentController().UpdateContentItem)
	admin.DELETE("/content/:id", controllers.NewContentController().DeleteContentItem)

	admin.GET("/lang", controllers.NewLanguageController().GetLanguageTexts)
	admin.GET("/lang/:id", controllers.NewLanguageController().GetLanguageText)
	admin.POST("/lang", controllers.NewLanguageController().CreateLanguageText)
	admin.PUT("/lang/:id", controllers.NewLanguageController().UpdateLanguageText)
	admin.DELETE("/lang/:id", controllers.NewLanguageController().DeleteLanguageText)
	admin.GET("/lang/:id/versions", controllers.NewLanguageController().GetLanguageTextVersions)
	admin.POST("/lang/:id/restore/:version", controllers.NewLanguageController().RestoreLanguageTextVersion)

	admin.GET("/site-settings", controllers.NewSiteSettingController().GetSettings)
	admin.POST("/site-settings", controllers.NewSiteSettingController().SaveSettings)

	// 联系表单提交管理
	admin.GET("/contact-submissions", controllers.NewContactSubmissionController().GetContactSubmissions)
	admin.DELETE("/contact-submissions/:id", controllers.NewContactSubmissionController().DeleteContactSubmission)
}