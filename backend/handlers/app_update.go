package handlers

import (
	"backend/models"
	"backend/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	androidAppPath = "storage/apps/android"
)

// 确保存储目录存在
func init() {
	// 创建存储目录
	err := os.MkdirAll(androidAppPath, 0755)
	if err != nil {
		panic(fmt.Sprintf("Failed to create app storage directory: %v", err))
	}
}

// UploadAppUpdate 上传新的应用更新
func UploadAppUpdate(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的文件"})
		return
	}

	// 获取其他表单数据
	version := c.PostForm("version")
	platform := c.PostForm("platform")
	forceUpdate := c.PostForm("force_update") == "true"
	updateNotes := c.PostForm("update_notes")

	if version == "" || platform == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "版本号和平台信息不能为空"})
		return
	}

	// 生成文件名
	fileName := fmt.Sprintf("%s_%s.apk", platform, version)
	filePath := filepath.Join(androidAppPath, fileName)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 计算文件大小和MD5
	fileSize := file.Size
	md5, err := utils.CalculateMD5(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "计算文件MD5失败"})
		return
	}

	// 创建更新记录
	appUpdate := models.AppUpdate{
		Version:     version,
		Platform:    platform,
		UpdateURL:   fmt.Sprintf("/api/app-updates/download/%s/%s", platform, version),
		ForceUpdate: forceUpdate,
		UpdateNotes: updateNotes,
		ReleaseDate: time.Now(),
		FileSize:    fileSize,
		MD5:         md5,
	}

	if err := models.DB.Create(&appUpdate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建更新记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "应用更新上传成功",
		"data":    appUpdate,
	})
}

// GetAppUpdates 获取应用更新列表
func GetAppUpdates(c *gin.Context) {
	platform := c.Query("platform")
	query := models.DB.Model(&models.AppUpdate{})

	if platform != "" {
		query = query.Where("platform = ?", platform)
	}

	var updates []models.AppUpdate
	if err := query.Order("release_date desc").Find(&updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取更新列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updates})
}

// CheckUpdate 检查更新
func CheckUpdate(c *gin.Context) {
	currentVersion := c.Query("current_version")
	platform := c.Query("platform")

	if currentVersion == "" || platform == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "当前版本号和平台信息不能为空"})
		return
	}

	var latestUpdate models.AppUpdate
	if err := models.DB.Where("platform = ?", platform).
		Order("release_date desc").
		First(&latestUpdate).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到更新信息"})
		return
	}

	// 比较版本号，判断是否需要更新
	needsUpdate := utils.CompareVersions(currentVersion, latestUpdate.Version) < 0

	c.JSON(http.StatusOK, gin.H{
		"needs_update": needsUpdate,
		"update_info":  latestUpdate,
	})
}

// DownloadUpdate 下载更新包
func DownloadUpdate(c *gin.Context) {
	platform := c.Param("platform")
	version := c.Param("version")

	var update models.AppUpdate
	if err := models.DB.Where("platform = ? AND version = ?", platform, version).
		First(&update).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到更新包"})
		return
	}

	filePath := filepath.Join(androidAppPath, fmt.Sprintf("%s_%s.apk", platform, version))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "更新包文件不存在"})
		return
	}

	c.File(filePath)
}

// DeleteUpdate 删除更新包
func DeleteUpdate(c *gin.Context) {
	id := c.Param("id")

	var update models.AppUpdate
	if err := models.DB.First(&update, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到更新记录"})
		return
	}

	// 删除文件
	filePath := filepath.Join(androidAppPath, fmt.Sprintf("%s_%s.apk", update.Platform, update.Version))
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文件失败"})
		return
	}

	// 删除数据库记录
	if err := models.DB.Delete(&update).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除更新记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新包删除成功"})
}
