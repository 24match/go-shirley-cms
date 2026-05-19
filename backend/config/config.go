package config

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	JWTSecret = []byte("medical-device-cms-secret-key-2024")
	DB        *gorm.DB
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	return r
}
