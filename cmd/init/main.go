package main

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"medical-device-cms/backend/config"
	"medical-device-cms/backend/models"
)

func main() {
	// 初始化数据库连接
	config.InitDB()

	// 检查是否已存在超级管理员
	var superAdmin models.User
	err := config.DB.Where("username = ?", "admin@system.com").First(&superAdmin).Error
	if err == nil {
		fmt.Println("超级管理员已存在")
	} else {
		// 创建超级管理员
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		superAdmin = models.User{
			Username: "admin@system.com",
			Password: string(hashedPassword),
			Role:     "superadmin",
		}
		if err := config.DB.Create(&superAdmin).Error; err != nil {
			log.Printf("创建超级管理员失败：%v", err)
		} else {
			fmt.Println("✅ 超级管理员已创建：admin@system.com / admin123")
		}
	}

	// 检查是否已存在租户管理员
	var tenant models.User
	err = config.DB.Where("username = ?", "tenant@example.com").First(&tenant).Error
	if err == nil {
		fmt.Println("租户管理员已存在")
	} else {
		// 创建租户管理员
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		tenant = models.User{
			Username: "tenant@example.com",
			Password: string(hashedPassword),
			Role:     "tenant",
		}
		if err := config.DB.Create(&tenant).Error; err != nil {
			log.Printf("创建租户管理员失败：%v", err)
		} else {
			fmt.Println("✅ 租户管理员已创建：tenant@example.com / password123")
		}
	}

	fmt.Println("\n🎉 默认用户初始化完成")
	fmt.Println("\n现在可以使用以下账号登录管理后台：")
	fmt.Println("  超级管理员：admin@system.com / admin123")
	fmt.Println("  租户管理员：tenant@example.com / password123")
}