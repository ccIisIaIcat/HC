package handlers

import (
	"net/http"
	"strconv"

	"backend/models" // 注意：这里需要替换为实际的项目路径

	"github.com/gin-gonic/gin"
)

// CreateItem 创建新物品
func CreateItem(c *gin.Context) {
	var item models.Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 检查物品名称是否已存在
	exists, err := models.ExistsItemByName(item.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "检查物品名称时发生错误"})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该物品名称已存在"})
		return
	}

	if err := models.CreateItem(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建物品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// GetItem 获取单个物品信息
func GetItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的物品ID"})
		return
	}

	item, err := models.GetItemByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "物品不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// GetItems 获取物品列表（支持分页）
func GetItems(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	items, total, err := models.GetAllItems(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取物品列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": items,
		"meta": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// UpdateItem 更新物品信息
func UpdateItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的物品ID"})
		return
	}

	// 先检查物品是否存在
	existingItem, err := models.GetItemByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "物品不存在"})
		return
	}

	// 绑定更新数据
	if err := c.ShouldBindJSON(existingItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := models.UpdateItem(existingItem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新物品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": existingItem})
}

// DeleteItem 删除物品
func DeleteItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的物品ID"})
		return
	}

	if err := models.DeleteItem(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除物品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "物品已成功删除"})
}

// SearchItems 搜索物品
func SearchItems(c *gin.Context) {
	name := c.Query("name")
	source := c.Query("source")

	var items []models.Item
	var err error

	if name != "" {
		items, err = models.GetItemsByName(name)
	} else if source != "" {
		items, err = models.GetItemsBySource(source)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供搜索条件"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索物品失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

// GetItemSources 获取所有物品来源
func GetItemSources(c *gin.Context) {
	sources, err := models.GetItemSources()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取物品来源失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sources})
}
