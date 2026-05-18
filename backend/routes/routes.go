package routes

import (
	"github.com/gin-gonic/gin"
	"medical-device-cms/backend/controllers"
	"medical-device-cms/backend/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	r.Use(middleware.CORSMiddleware())

	r.Static("/uploads", "./uploads")
	r.Static("/frontend", "./frontend")
	r.StaticFile("/", "./inde.html")
	r.StaticFile("/i18n.js", "./i18n.js")
	r.Static("/admin", "./admin")

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
}

func registerAdminRoutes(api *gin.RouterGroup) {
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware())

	admin.GET("/images", controllers.NewImageController().GetImages)
	admin.POST("/images", controllers.NewImageController().UploadImage)
	admin.POST("/images/batch", controllers.NewImageController().UploadMultipleImages)
	admin.PUT("/images/:id", controllers.NewImageController().UpdateImage)
	admin.DELETE("/images/:id", controllers.NewImageController().DeleteImage)

	admin.GET("/config", controllers.NewModuleController().GetPageConfig)
	admin.PUT("/config", controllers.NewModuleController().UpdatePageConfig)

	admin.GET("/modules", controllers.NewModuleController().GetModuleConfigs)
	admin.GET("/modules/:name", controllers.NewModuleController().GetModuleConfig)
	admin.POST("/modules", controllers.NewModuleController().SaveModuleConfig)
	admin.PUT("/modules/:name", controllers.NewModuleController().SaveModuleConfig)
	admin.DELETE("/modules/:name", controllers.NewModuleController().DeleteModuleConfig)

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
}