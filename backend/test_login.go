package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID       uint
	Username string
	Password string
	Role     string
	TenantID uint
}

type Tenant struct {
	ID     uint
	Status string
}

func main() {
	db, err := gorm.Open(sqlite.Open("backend/medical.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败:", err)
		return
	}

	// 测试登录
	username := "admin"
	password := "admin123"

	var user User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		fmt.Println("用户不存在:", err)
		return
	}

	fmt.Printf("找到用户：ID=%d, Username=%s, Role=%s, TenantID=%d\n", user.ID, user.Username, user.Role, user.TenantID)

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println("密码验证失败:", err)
		return
	}
	fmt.Println("密码验证成功")

	// 检查租户状态
	if user.Role != "superadmin" && user.TenantID != 0 {
		var tenant Tenant
		if err := db.First(&tenant, user.TenantID).Error; err != nil {
			fmt.Println("租户不存在:", err)
			return
		}
		fmt.Printf("租户状态：%s\n", tenant.Status)
		if tenant.Status != "active" {
			fmt.Println("租户未激活")
			return
		}
	}

	fmt.Println("登录成功！")
}