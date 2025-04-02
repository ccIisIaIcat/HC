package handlers

import (
	"backend/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// InitDB 初始化数据库连接
func InitDB() {
	// 调用models包中的初始化函数，确保DB变量被初始化
	models.InitDB()

	// 检查管理员用户
	ensureAdminUser()
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return models.DB
}

// ensureAdminUser 确保系统中存在管理员用户
func ensureAdminUser() {
	// 创建默认管理员账户
	var adminCount int64
	models.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		// 生成密码哈希
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("生成管理员密码哈希失败：", err)
		}

		adminUser := models.User{
			Name:     "管理员",
			Email:    "admin@example.com",
			Password: string(hashedPassword),
			Role:     "admin",
			Verified: true,
		}
		if result := models.DB.Create(&adminUser); result.Error != nil {
			log.Fatal("创建管理员账户失败：", result.Error)
		}
		log.Println("默认管理员账户已创建：admin@example.com / admin123")
	}
}
