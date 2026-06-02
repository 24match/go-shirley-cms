package routes

import (
	"medical-device-cms/backend/controllers"
	"medical-device-cms/backend/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes 注册所有路由
// @description 初始化 CORS、错误处理中间件，注册静态文件路由、API 路由和 Swagger 文档路由
func RegisterRoutes(r *gin.Engine) {
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorHandler())

	r.Static("/uploads", "./uploads")
	r.Static("/frontend", "./frontend")
	r.StaticFile("/", "./inde.html")
	r.StaticFile("/i18n.js", "./i18n.js")
	r.Static("/admin", "./admin")

	// 注册 Swagger UI 路由（仅开发环境）
	// @description 访问 /swagger/index.html 查看交互式 API 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	registerPublicRoutes(api)
	registerAdminRoutes(api)
}

func registerPublicRoutes(api *gin.RouterGroup) {
	api.POST("/login", controllers.NewUserController().Login)
	api.GET("/public/config", controllers.NewModuleController().GetPublicPageConfig)
	api.GET("/public/modules", controllers.NewModuleController().GetPublicModuleConfigs)
	api.GET("/public/images", controllers.NewImageController().GetImages)
	api.GET("/public/content", controllers.NewContentController().GetPublicContentItems)
	api.GET("/public/lang", controllers.NewLanguageController().GetPublicLanguageTexts)
	api.GET("/public/site-settings", controllers.NewSiteSettingController().GetPublicSettings)
	
	// 联系表单提交（公开接口）
	api.POST("/contact/submit", controllers.NewContactSubmissionController().SubmitContactForm)
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
