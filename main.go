// @title Shirley CMS API
// @version 1.0
// @description 医疗器械内容管理系统 API 接口文档
// @description 提供内容管理、模块配置、图片管理、多语言支持等功能

// @contact.name API Support
// @contact.email support@example.com

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey APIKey
// @in header
// @name Authorization
// @description JWT Token 认证，格式：Bearer {token}

package main

import (
	"fmt"

	_ "medical-device-cms/docs"

	"medical-device-cms/backend/config"
	"medical-device-cms/backend/routes"
)

// @summary 启动应用
// @description 初始化数据库、配置路由并启动 HTTP 服务器
func main() {
	if err := config.InitDB(); err != nil {
		panic(err)
	}

	r := config.SetupRouter()
	routes.RegisterRoutes(r)

	fmt.Println("Server starting on :8080")
	r.Run(":8080")
}