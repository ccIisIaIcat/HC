package handlers

import (
	"backend/models"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 注册用户健康状态相关路由
func RegisterUserHealthStateRoutes(router *gin.Engine) {
	// 需要认证的API组
	healthStateGroup := router.Group("/api")
	healthStateGroup.Use(AuthMiddleware()) // 应用认证中间件

	// 用户健康状态相关路由
	healthStateGroup.POST("/health-states", CreateUserHealthStateHandler)
	healthStateGroup.GET("/health-states", GetUserHealthStatesHandler)
	healthStateGroup.GET("/health-states/:id", GetUserHealthStateHandler)
	healthStateGroup.PUT("/health-states/:id", UpdateUserHealthStateHandler)
	healthStateGroup.DELETE("/health-states/:id", DeleteUserHealthStateHandler)
	healthStateGroup.GET("/health-states/latest", GetLatestUserHealthStateHandler)
}

// 创建用户健康状态处理函数
func CreateUserHealthStateHandler(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 解析请求体
	var record models.UserHealthState
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
	if err := models.CreateUserHealthState(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存健康状态记录失败"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "健康状态记录创建成功",
		"record":  record,
	})
}

// 获取用户的健康状态记录处理函数
func GetUserHealthStatesHandler(c *gin.Context) {
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
	records, err := models.GetUserHealthStates(userID.(uint), startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取健康状态记录失败"})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"records": records,
		"total":   len(records),
	})
}

// 获取单个健康状态记录处理函数
func GetUserHealthStateHandler(c *gin.Context) {
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
	record, err := models.GetUserHealthStateByID(uint(recordID))
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

// 获取用户最新的健康状态记录处理函数
func GetLatestUserHealthStateHandler(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	// 获取最新记录
	record, err := models.GetLatestUserHealthState(userID.(uint))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "没有找到健康状态记录"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取记录失败"})
		}
		return
	}

	// 返回记录
	c.JSON(http.StatusOK, gin.H{"record": record})
}

// 更新健康状态记录处理函数
func UpdateUserHealthStateHandler(c *gin.Context) {
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
	existingRecord, err := models.GetUserHealthStateByID(uint(recordID))
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
	var updatedRecord models.UserHealthState
	if err := c.ShouldBindJSON(&updatedRecord); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式错误"})
		return
	}

	// 保留ID和UserID不变
	updatedRecord.ID = existingRecord.ID
	updatedRecord.UserID = existingRecord.UserID

	// 保存更新
	if err := models.UpdateUserHealthState(&updatedRecord); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新记录失败"})
		return
	}

	// 返回更新后的记录
	c.JSON(http.StatusOK, gin.H{
		"message": "记录更新成功",
		"record":  updatedRecord,
	})
}

// 删除健康状态记录处理函数
func DeleteUserHealthStateHandler(c *gin.Context) {
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
	existingRecord, err := models.GetUserHealthStateByID(uint(recordID))
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
	if err := models.DeleteUserHealthState(uint(recordID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除记录失败"})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "记录已成功删除"})
}
