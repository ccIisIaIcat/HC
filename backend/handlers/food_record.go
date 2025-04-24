package handlers

import (
	"backend/achievement"
	"backend/models"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 全局食物分析处理器
var foodAnalysisHandler *FoodAnalysisHandler

// 设置食物分析处理器
func SetFoodAnalysisHandler(handler *FoodAnalysisHandler) {
	foodAnalysisHandler = handler
}

// AnalyzeFood 分析食物图片并返回分析结果
func AnalyzeFood(c *gin.Context, file *multipart.FileHeader, imageDescription string) (*models.FoodAnalysis, error) {
	// 检查食物分析处理器是否已设置
	if foodAnalysisHandler == nil {
		return nil, fmt.Errorf("食物分析处理器未初始化")
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return nil, fmt.Errorf("只支持JPG和PNG格式的图片")
	}

	// 检查文件大小（限制为10MB）
	if file.Size > 10*1024*1024 {
		return nil, fmt.Errorf("图片大小不能超过10MB")
	}

	// 读取文件内容
	openedFile, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer openedFile.Close()

	// 读取文件内容到内存
	imageBytes, err := io.ReadAll(openedFile)
	if err != nil {
		return nil, fmt.Errorf("读取文件内容失败: %v", err)
	}

	// 根据是否有图片描述调用不同的方法
	if imageDescription != "" {
		return foodAnalysisHandler.AnalyzeImageBytesWithDescription(imageBytes, imageDescription)
	}

	// 无描述时调用原始方法
	return foodAnalysisHandler.AnalyzeImageBytes(imageBytes)
}

// 注册食物记录相关路由
func RegisterFoodRecordRoutes(router *gin.Engine) {
	// 需要认证的API组
	foodRecordGroup := router.Group("/api")
	foodRecordGroup.Use(AuthMiddleware()) // 应用认证中间件

	// 食物记录相关路由
	foodRecordGroup.POST("/food-records", CreateFoodRecordHandler)
	foodRecordGroup.GET("/food-records", GetFoodRecordsHandler)
	foodRecordGroup.GET("/food-records/:id", GetFoodRecordHandler)
	foodRecordGroup.PUT("/food-records/:id", UpdateFoodRecordHandler)
	foodRecordGroup.DELETE("/food-records/:id", DeleteFoodRecordHandler)

	// 从食物分析创建食物记录
	foodRecordGroup.POST("/analyze-and-save", AnalyzeAndSaveFoodHandler)
}

// 创建食物记录处理函数
func CreateFoodRecordHandler(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 解析请求体
	var record models.FoodRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式错误"})
		return
	}

	// 设置用户ID
	record.UserID = userID.(uint)

	// 如果没有提供记录时间，默认为当前时间
	if record.RecordTime.IsZero() {
		record.RecordTime = time.Now()
	}

	// 保存到数据库
	if err := models.CreateFoodRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存食物记录失败"})
		return
	}
	achievement.CheckAchievements(userID.(uint))
	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "食物记录创建成功",
		"record":  record,
	})
}

// 获取用户的食物记录处理函数
func GetFoodRecordsHandler(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取查询参数
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startTime, endTime time.Time
	var err error

	// 如果提供了开始日期，解析它
	if startDateStr != "" {
		startTime, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "开始日期格式错误"})
			return
		}
	} else {
		// 默认为30天前
		startTime = time.Now().AddDate(0, 0, -30)
	}

	// 如果提供了结束日期，解析它
	if endDateStr != "" {
		endTime, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "结束日期格式错误"})
			return
		}
		// 设置为当天的结束时间
		endTime = endTime.Add(24*time.Hour - time.Second)
	} else {
		// 默认为当前时间
		endTime = time.Now()
	}

	// 获取记录
	var records []models.FoodRecord
	var getErr error

	if startDateStr == "" && endDateStr == "" {
		// 如果没有日期过滤，获取所有记录
		records, getErr = models.GetAllUserFoodRecords(userID.(uint))
	} else {
		// 获取指定日期范围的记录
		records, getErr = models.GetUserFoodRecords(userID.(uint), startTime, endTime)
	}

	if getErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取食物记录失败"})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"records": records,
		"total":   len(records),
	})
}

// 获取单个食物记录处理函数
func GetFoodRecordHandler(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取记录ID
	recordIDStr := c.Param("id")
	recordID, err := strconv.ParseUint(recordIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	// 查询记录
	record, err := models.GetFoodRecordByID(uint(recordID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取记录失败"})
		}
		return
	}

	// 验证记录是否属于当前用户
	if record.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此记录"})
		return
	}

	// 返回记录
	c.JSON(http.StatusOK, gin.H{"record": record})
}

// 更新食物记录处理函数
func UpdateFoodRecordHandler(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取记录ID
	recordIDStr := c.Param("id")
	recordID, err := strconv.ParseUint(recordIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	// 查询现有记录
	existingRecord, err := models.GetFoodRecordByID(uint(recordID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取记录失败"})
		}
		return
	}

	// 验证记录是否属于当前用户
	if existingRecord.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权修改此记录"})
		return
	}

	// 解析请求体
	var updatedRecord models.FoodRecord
	if err := c.ShouldBindJSON(&updatedRecord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式错误"})
		return
	}

	// 保留ID和UserID不变
	updatedRecord.ID = existingRecord.ID
	updatedRecord.UserID = existingRecord.UserID

	// 保存更新
	if err := models.UpdateFoodRecord(&updatedRecord); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新记录失败"})
		return
	}
	achievement.CheckAchievements(userID.(uint))
	// 返回更新后的记录
	c.JSON(http.StatusOK, gin.H{
		"message": "记录更新成功",
		"record":  updatedRecord,
	})
}

// 删除食物记录处理函数
func DeleteFoodRecordHandler(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取记录ID
	recordIDStr := c.Param("id")
	recordID, err := strconv.ParseUint(recordIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	// 查询现有记录
	existingRecord, err := models.GetFoodRecordByID(uint(recordID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取记录失败"})
		}
		return
	}

	// 验证记录是否属于当前用户
	if existingRecord.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此记录"})
		return
	}

	// 删除记录
	if err := models.DeleteFoodRecord(uint(recordID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除记录失败"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "记录已成功删除"})
}

// 分析并保存食物记录处理函数
func AnalyzeAndSaveFoodHandler(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取食物图片
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供食物图片"})
		return
	}

	// 获取其他表单数据
	mealType := c.PostForm("meal_type")
	notes := c.PostForm("notes")
	imageDescription := c.PostForm("image_description")

	// 调用食物分析API（利用现有的分析函数）
	analysis, err := AnalyzeFood(c, file, imageDescription)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "分析食物失败: " + err.Error()})
		return
	}

	if !analysis.HasFood {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未能检测到食物"})
		return
	}

	// 创建食物记录
	record := models.CreateFoodRecordFromAnalysis(userID.(uint), analysis, mealType, notes)

	// 保存图片路径（如果需要）
	// ... 这里可以添加保存图片的代码 ...

	// 保存记录到数据库
	if err := models.CreateFoodRecord(record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存食物记录失败"})
		return
	}
	achievement.CheckAchievements(userID.(uint))
	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message":  "食物分析和记录保存成功",
		"record":   record,
		"analysis": analysis,
	})
}
