package handlers

import (
	"backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 定义中国时区
var cst *time.Location

func init() {
	var err error
	cst, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// 如果加载失败，手动设置UTC+8
		cst = time.FixedZone("CST", 8*3600)
	}
}

// CheckInRequest 打卡请求结构
type CheckInRequest struct {
	Content string `json:"content" binding:"-"` // 打卡内容（可选）
}

// CheckInResponse 打卡响应结构
type CheckInResponse struct {
	Success          bool            `json:"success"`                 // 是否成功
	Message          string          `json:"message"`                 // 消息
	HasFoodRecord    bool            `json:"has_food_record"`         // 是否有食物记录
	AlreadyCheckedIn bool            `json:"already_checked_in"`      // 是否已经打卡
	CheckInData      *models.CheckIn `json:"check_in_data,omitempty"` // 打卡数据
}

// HandleCheckIn 处理用户打卡请求
func HandleCheckIn(c *gin.Context) {
	// 1. 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, CheckInResponse{
			Success: false,
			Message: "用户未登录",
		})
		return
	}

	// 2. 检查用户今日是否已打卡
	hasCheckedIn, err := models.HasUserCheckedInToday(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, CheckInResponse{
			Success: false,
			Message: "检查打卡记录失败",
		})
		return
	}

	if hasCheckedIn {
		c.JSON(http.StatusBadRequest, CheckInResponse{
			Success:          false,
			Message:          "今日已打卡",
			AlreadyCheckedIn: true,
		})
		return
	}

	// 3. 检查用户今日是否有食物记录
	today := time.Now().In(cst)
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, cst)
	endOfDay := startOfDay.Add(24 * time.Hour)

	foodRecords, err := models.GetUserFoodRecords(userID.(uint), startOfDay, endOfDay)
	if err != nil {
		c.JSON(http.StatusInternalServerError, CheckInResponse{
			Success: false,
			Message: "检查食物记录失败",
		})
		return
	}

	if len(foodRecords) == 0 {
		c.JSON(http.StatusBadRequest, CheckInResponse{
			Success:       false,
			Message:       "今日还未记录任何食物，请先记录饮食后再打卡",
			HasFoodRecord: false,
		})
		return
	}

	// 4. 解析请求体
	var req CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CheckInResponse{
			Success: false,
			Message: "无效的请求数据",
		})
		return
	}

	// 5. 创建打卡记录
	checkIn := &models.CheckIn{
		UserID:    userID.(uint),
		CheckInAt: time.Now().In(cst), // 使用中国时区
		Content:   req.Content,
	}

	if err := models.CreateCheckIn(checkIn); err != nil {
		c.JSON(http.StatusInternalServerError, CheckInResponse{
			Success: false,
			Message: "创建打卡记录失败",
		})
		return
	}

	// 6. 返回成功响应
	c.JSON(http.StatusOK, CheckInResponse{
		Success:          true,
		Message:          "打卡成功",
		HasFoodRecord:    true,
		AlreadyCheckedIn: false,
		CheckInData:      checkIn,
	})
}

// GetTodayCheckIn 获取用户今日打卡状态
func GetTodayCheckIn(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户未登录",
		})
		return
	}

	// 检查今日是否已打卡
	hasCheckedIn, err := models.HasUserCheckedInToday(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "检查打卡记录失败",
		})
		return
	}

	// 检查今日是否有食物记录
	today := time.Now().In(cst)
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, cst)
	endOfDay := startOfDay.Add(24 * time.Hour)

	foodRecords, err := models.GetUserFoodRecords(userID.(uint), startOfDay, endOfDay)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "检查食物记录失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"has_checked_in":  hasCheckedIn,
			"has_food_record": len(foodRecords) > 0,
		},
	})
}

// GetUserCheckIns 获取用户所有打卡记录
func GetUserCheckIns(c *gin.Context) {
	// 1. 获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户未登录",
		})
		return
	}

	// 2. 从查询参数获取分页信息（可选）
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	// 3. 调用模型层获取打卡记录
	checkIns, err := models.GetUserCheckIns(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取打卡记录失败",
		})
		return
	}

	// 4. 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取打卡记录成功",
		"data": gin.H{
			"total":     len(checkIns),
			"records":   checkIns,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
