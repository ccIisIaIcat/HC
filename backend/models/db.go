package models

import (
	"backend/config"
	"fmt"
	mathrand "math/rand"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局数据库连接实例
var DB *gorm.DB

func init() {
	// 初始化随机数生成器
	mathrand.Seed(time.Now().UnixNano())
}

// InitDB 初始化数据库连接和表结构
func InitDB() *gorm.DB {
	conf := config.GetConfig()

	// 先尝试连接MySQL服务器
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBPort)), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MySQL server: %v", err))
	}

	// 创建数据库（如果不存在）
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", conf.DBName)
	if err := db.Exec(sql).Error; err != nil {
		panic(fmt.Sprintf("Failed to create database: %v", err))
	}

	// 连接到特定数据库
	db, err = gorm.Open(mysql.Open(conf.GetDSN()), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	// 自动迁移数据库表
	db.AutoMigrate(&User{}, &VerificationCode{}, &FoodRecord{})

	// 创建管理员账户（如果不存在）
	var adminCount int64
	db.Model(&User{}).Where("role = ?", "admin").Count(&adminCount)
	if adminCount == 0 {
		adminUser := User{
			Name:     "超级管理员",
			Email:    "admin@yourdomain.com",
			Password: "admin888",
			Role:     "admin",
			Verified: true, // 管理员账户默认已验证
		}
		db.Create(&adminUser)
		fmt.Println("默认管理员账户已创建：admin@yourdomain.com / admin888")
	}

	// 设置全局DB变量
	DB = db
	return db
}

// 用户角色常量
const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// User 模型
type User struct {
	gorm.Model
	Name     string `json:"Name" gorm:"column:name"`
	Email    string `json:"Email" gorm:"column:email;unique"`
	Password string `json:"Password,omitempty" gorm:"column:password"`     // 返回时省略
	Role     string `json:"Role" gorm:"column:role;default:user"`          // 用户角色：user 或 admin
	Verified bool   `json:"Verified" gorm:"column:verified;default:false"` // 邮箱是否已验证
}

// VerificationCode 模型 - 用于邮箱验证
type VerificationCode struct {
	ID        uint      `gorm:"primarykey"`
	Email     string    `gorm:"column:email;index"` // 用户邮箱
	Code      string    `gorm:"column:code"`        // 验证码
	ExpiresAt time.Time `gorm:"column:expires_at"`  // 过期时间
	CreatedAt time.Time
	Used      bool   `gorm:"column:used;default:false"` // 是否已使用
	Name      string `gorm:"column:name"`               // 用户注册时的姓名
	Password  string `gorm:"column:password"`           // 用户注册时的密码
}

// 生成随机验证码
func GenerateVerificationCode(email string) (*VerificationCode, error) {
	// 生成6位纯数字验证码
	code := ""
	// 使用math/rand包生成数字
	for i := 0; i < 6; i++ {
		digit := mathrand.Intn(10) // 生成0-9之间的随机数
		code += fmt.Sprintf("%d", digit)
	}

	// 设置过期时间为24小时后
	expiresAt := time.Now().Add(24 * time.Hour)

	// 创建验证码记录
	verificationCode := &VerificationCode{
		Email:     email,
		Code:      code,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		Used:      false,
	}

	// 删除该邮箱之前的验证码
	DB.Where("email = ?", email).Delete(&VerificationCode{})

	// 保存新验证码
	result := DB.Create(verificationCode)
	if result.Error != nil {
		return nil, result.Error
	}

	return verificationCode, nil
}

// 验证验证码
func VerifyCode(email, code string) (*VerificationCode, bool) {
	var verificationCode VerificationCode

	// 查找验证码记录
	result := DB.Where("email = ? AND code = ? AND used = ? AND expires_at > ?",
		email, code, false, time.Now()).First(&verificationCode)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, false
	}

	// 标记验证码为已使用
	DB.Model(&verificationCode).Updates(map[string]interface{}{
		"used": true,
	})

	return &verificationCode, true
}
