package handlers

import (
	"backend/models"
	"backend/utils"
	"fmt"
	mathrand "math/rand"
	"net/http"
	"os"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtKey = []byte(os.Getenv("JWT_SECRET"))
)

func init() {
	// 初始化随机数生成器
	mathrand.Seed(time.Now().UnixNano())
}

// 注册请求
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// 验证邮箱请求
type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

// JWT Claims
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Register 处理用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("无效的请求数据: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	log.Printf("尝试注册: %s", req.Email)

	// 检查邮箱是否已存在
	var existingUser models.User
	if result := models.DB.Where("email = ?", req.Email).First(&existingUser); result.Error == nil {
		log.Printf("邮箱已被注册: %s", req.Email)
		c.JSON(http.StatusConflict, gin.H{"error": "邮箱已被注册"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("密码加密失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 生成验证码
	verificationCode, err := generateVerificationCode(req.Email)
	if err != nil {
		log.Printf("生成验证码失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成验证码失败"})
		return
	}

	// 保存用户信息到验证码记录中
	verificationCode.Name = req.Name
	verificationCode.Password = string(hashedPassword)
	if result := models.DB.Save(verificationCode); result.Error != nil {
		log.Printf("保存验证码记录失败: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存验证码记录失败"})
		return
	}

	// 发送验证邮件
	emailBody := fmt.Sprintf(`
		<h2>欢迎注册</h2>
		<p>您的验证码是：<strong>%s</strong></p>
		<p>验证码有效期为24小时。</p>
	`, verificationCode.Code)

	if err := utils.SendEmail(req.Email, "邮箱验证", emailBody); err != nil {
		log.Printf("发送验证邮件失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送验证邮件失败"})
		return
	}

	log.Printf("注册请求已处理，等待邮箱验证: %s", req.Email)
	c.JSON(http.StatusOK, gin.H{
		"message": "请查收验证邮件完成注册",
		"email":   req.Email,
	})
}

// Login 处理用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("无效的请求数据: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	log.Printf("尝试登录: %s", req.Email)

	var user models.User
	if result := models.DB.Where("email = ?", req.Email).First(&user); result.Error != nil {
		log.Printf("用户不存在: %s", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码错误"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Printf("密码错误: %s", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "邮箱或密码错误"})
		return
	}

	// 检查邮箱是否已验证
	if !user.Verified {
		log.Printf("邮箱未验证: %s", req.Email)
		c.JSON(http.StatusForbidden, gin.H{
			"error":             "邮箱未验证",
			"need_verification": true,
		})
		return
	}

	// 生成JWT
	token, err := generateToken(user)
	if err != nil {
		log.Printf("生成token失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	log.Printf("登录成功: %s", req.Email)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// VerifyEmail 处理邮箱验证
func VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("无效的请求数据: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	log.Printf("尝试验证邮箱: %s", req.Email)

	// 查找验证码记录
	var verificationCode models.VerificationCode
	result := models.DB.Where("email = ? AND code = ? AND used = ? AND expires_at > ?",
		req.Email, req.Code, false, time.Now()).First(&verificationCode)

	if result.Error != nil || result.RowsAffected == 0 {
		log.Printf("验证码无效或已过期: %s", req.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码无效或已过期"})
		return
	}

	// 创建新用户
	user := models.User{
		Name:     verificationCode.Name,
		Email:    req.Email,
		Password: verificationCode.Password, // 使用之前加密的密码
		Role:     "user",
		Verified: true,
	}

	if result := models.DB.Create(&user); result.Error != nil {
		log.Printf("创建用户失败: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	// 标记验证码为已使用
	models.DB.Model(&verificationCode).Updates(map[string]interface{}{
		"used": true,
	})

	log.Printf("邮箱验证成功，用户已创建: %s", req.Email)
	c.JSON(http.StatusOK, gin.H{
		"message": "邮箱验证成功，请登录",
		"user":    user,
	})
}

// ResendVerification 重新发送验证邮件
func ResendVerification(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 检查用户是否已经验证
	var user models.User
	if err := models.DB.Where("email = ? AND verified = ?", req.Email, true).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该邮箱已经验证，可直接登录"})
		return
	}

	// 发送验证邮件
	if err := sendVerificationEmail(req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送验证邮件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "验证邮件已发送"})
}

// GetCurrentUser 获取当前用户信息
func GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var user models.User
	if result := models.DB.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// GetUsers 获取所有用户列表（管理员专用）
func GetUsers(c *gin.Context) {
	var users []models.User
	if result := models.DB.Find(&users); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetStats 获取统计信息（管理员专用）
func GetStats(c *gin.Context) {
	var totalUsers int64
	var verifiedUsers int64

	models.DB.Model(&models.User{}).Count(&totalUsers)
	models.DB.Model(&models.User{}).Where("verified = ?", true).Count(&verifiedUsers)

	c.JSON(http.StatusOK, gin.H{
		"total_users":      totalUsers,
		"verified_users":   verifiedUsers,
		"unverified_users": totalUsers - verifiedUsers,
	})
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供token"})
			c.Abort()
			return
		}

		// 移除 "Bearer " 前缀
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// AdminAuthMiddleware 管理员认证中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 生成JWT token
func generateToken(user models.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// 生成验证码
func generateVerificationCode(email string) (*models.VerificationCode, error) {
	// 生成6位纯数字验证码
	code := ""
	for i := 0; i < 6; i++ {
		digit := mathrand.Intn(10) // 生成0-9之间的随机数
		code += fmt.Sprintf("%d", digit)
	}

	// 设置过期时间为24小时后
	expiresAt := time.Now().Add(24 * time.Hour)

	// 创建验证码记录
	verificationCode := &models.VerificationCode{
		Email:     email,
		Code:      code,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		Used:      false,
	}

	// 删除该邮箱之前的验证码
	models.DB.Where("email = ?", email).Delete(&models.VerificationCode{})

	// 保存新验证码
	if result := models.DB.Create(verificationCode); result.Error != nil {
		return nil, result.Error
	}

	return verificationCode, nil
}

// 发送验证邮件
func sendVerificationEmail(email string) error {
	// 生成新的验证码
	verificationCode, err := generateVerificationCode(email)
	if err != nil {
		return err
	}

	// 构建邮件内容
	emailBody := fmt.Sprintf(`
		<h2>邮箱验证</h2>
		<p>您的验证码是：<strong>%s</strong></p>
		<p>验证码有效期为24小时。</p>
	`, verificationCode.Code)

	// 发送邮件
	return utils.SendEmail(email, "邮箱验证", emailBody)
}
