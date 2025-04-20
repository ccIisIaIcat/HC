package handlers

import (
	"net/http"
	"strconv"
	"time"

	"backend/models"

	"github.com/gin-gonic/gin"
)

// CreateUserItem 创建用户物品关联
func CreateUserItem(c *gin.Context) {
	// 从URL参数中获取用户ID
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var userItem models.UserItem
	if err := c.ShouldBindJSON(&userItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 设置用户ID和获得时间
	userItem.UserID = uint(userID)
	userItem.ObtainedAt = time.Now()

	// 如果没有提供获得来源，设置默认值
	if userItem.ObtainedFrom == "" {
		userItem.ObtainedFrom = "未知来源"
	}

	// 检查物品是否存在
	_, err = models.GetItemByID(userItem.ItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "物品不存在"})
		return
	}

	// 创建用户物品关联
	if err := models.CreateUserItem(&userItem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户物品关联失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": userItem})
}

// GetUserItems 获取用户的所有物品
func GetUserItems(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	items, err := models.GetUserItems(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户物品列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetUserItemDetails 获取用户特定物品的详细信息
func GetUserItemDetails(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的物品ID"})
		return
	}

	userItem, err := models.GetUserItemDetails(uint(userID), uint(itemID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到用户物品信息"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": userItem})
}

// UpdateUserItemQuantity 更新用户物品数量
func UpdateUserItemQuantity(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的物品ID"})
		return
	}

	var request struct {
		Quantity int `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的数量数据"})
		return
	}

	if request.Quantity < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "数量不能为负数"})
		return
	}

	err = models.UpdateUserItemQuantity(uint(userID), uint(itemID), request.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新物品数量失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "物品数量更新成功"})
}

// DeleteUserItem 删除用户物品
func DeleteUserItem(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的物品ID"})
		return
	}

	err = models.DeleteUserItem(uint(userID), uint(itemID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户物品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户物品删除成功"})
}

// GetUserItemsBySource 获取用户特定来源的物品
func GetUserItemsBySource(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	source := c.Query("source")
	if source == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供物品来源"})
		return
	}

	items, err := models.GetUserItemsBySource(uint(userID), source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户物品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetUserItemCount 获取用户物品总数
func GetUserItemCount(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	count, err := models.GetUserItemCount(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户物品总数失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": count})
}

// CheckUserHasItem 检查用户是否拥有特定物品
func CheckUserHasItem(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	itemID, err := strconv.ParseUint(c.Param("item_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的物品ID"})
		return
	}

	hasItem, err := models.HasUserItem(uint(userID), uint(itemID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查用户物品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"has_item": hasItem})
}
