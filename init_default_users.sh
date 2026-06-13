#!/bin/bash

# 初始化默认用户脚本
# 用于在数据库中创建默认的超级管理员和租户管理员账号

echo "🔧 初始化默认用户..."

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo "❌ 未检测到 Go 环境，请先安装 Go"
    exit 1
fi

# 检查数据库是否存在
if [ ! -f "cms.db" ]; then
    echo "⚠️ 数据库文件不存在，先启动程序初始化数据库..."
    timeout 5 go run main.go || true
fi

# 使用 Go 程序创建默认用户
echo "📝 创建默认用户..."

# 创建初始化用户的 Go 程序
cat > init_users.go << 'EOF'
package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"medical-device-cms/backend/config"
	"medical-device-cms/backend/models"
)

func main() {
	// 初始化数据库连接
	config.InitDB()

	// 检查是否已存在用户
	var count int64
	config.DB.Model(&models.User{}).Count(&count)
	
	if count > 0 {
		fmt.Println("数据库中已存在用户，跳过初始化")
		os.Exit(0)
	}

	// 创建超级管理员
	superAdminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	superAdmin := models.User{
		Username: "admin@system.com",
		Password: string(superAdminPassword),
		Email:    "admin@system.com",
		Role:     "superadmin",
	}
	if err := config.DB.Create(&superAdmin).Error; err != nil {
		log.Printf("创建超级管理员失败：%v", err)
	} else {
		fmt.Println("✅ 超级管理员已创建：admin@system.com / admin123")
	}

	// 创建租户管理员
	tenantPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	tenant := models.User{
		Username: "tenant@example.com",
		Password: string(tenantPassword),
		Email:    "tenant@example.com",
		Role:     "tenant",
	}
	if err := config.DB.Create(&tenant).Error; err != nil {
		log.Printf("创建租户管理员失败：%v", err)
	} else {
		fmt.Println("✅ 租户管理员已创建：tenant@example.com / password123")
	}

	fmt.Println("🎉 默认用户初始化完成")
}
EOF

# 运行初始化程序
go run init_users.go

# 清理临时文件
rm -f init_users.go

echo "✨ 初始化完成！"
echo ""
echo "默认账号："
echo "  超级管理员：admin@system.com / admin123"
echo "  租户管理员：tenant@example.com / password123"