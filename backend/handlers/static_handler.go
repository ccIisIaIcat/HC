package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StaticFileHandler struct{}

func NewStaticFileHandler() *StaticFileHandler {
	// 确保静态文件目录存在
	ensureDirectoryExists("./static/images")
	return &StaticFileHandler{}
}

// 确保目录存在
func ensureDirectoryExists(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		panic(fmt.Sprintf("无法创建目录 %s: %v", path, err))
	}
}

// UploadImage 处理图片上传
func (h *StaticFileHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "获取上传文件失败",
		})
		return
	}

	// 获取目录参数，默认为images
	directory := c.DefaultPostForm("directory", "images")

	// 获取自定义文件名，如果没有提供则生成UUID
	customFilename := c.PostForm("filename")

	// 验证目录名称安全性
	if !isValidDirectory(directory) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的目录名称",
		})
		return
	}

	// 验证文件类型
	if !isAllowedImageType(file.Header.Get("Content-Type")) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "不支持的文件类型，仅支持 JPG, PNG, GIF 格式",
		})
		return
	}

	// 验证文件大小（限制为2MB）
	if file.Size > 2*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "文件大小超过限制（最大2MB）",
		})
		return
	}

	var newFileName string
	if customFilename != "" {
		// 使用自定义文件名，但保留原始扩展名
		ext := filepath.Ext(file.Filename)
		baseCustomName := strings.TrimSuffix(customFilename, filepath.Ext(customFilename))
		newFileName = baseCustomName + ext
	} else {
		// 如果没有提供自定义文件名，则使用UUID
		ext := filepath.Ext(file.Filename)
		newFileName = fmt.Sprintf("%s%s", uuid.New().String(), ext)
	}

	// 构建目录路径
	dirPath := filepath.Join("static", directory)

	// 确保目录存在
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建目录失败",
		})
		return
	}

	// 构建保存路径
	savePath := filepath.Join(dirPath, newFileName)

	// 检查文件是否已存在
	if _, err := os.Stat(savePath); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "文件已存在",
		})
		return
	}

	// 保存文件
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "保存文件失败",
		})
		return
	}

	// 返回可访问的URL
	imageURL := fmt.Sprintf("/static/%s/%s", directory, newFileName)

	c.JSON(http.StatusOK, gin.H{
		"url":     imageURL,
		"message": "上传成功",
	})
}

// 验证目录名称安全性
func isValidDirectory(directory string) bool {
	// 不允许空目录名
	if directory == "" {
		return false
	}

	// 清理路径（移除 .. 和多余的斜杠）
	cleanPath := filepath.Clean(directory)

	// 不允许以点开头（避免访问上级目录）
	if strings.HasPrefix(cleanPath, ".") {
		return false
	}

	// 不允许绝对路径
	if filepath.IsAbs(cleanPath) {
		return false
	}

	// 检查每个目录部分
	parts := strings.Split(cleanPath, string(filepath.Separator))
	for _, part := range parts {
		// 不允许空目录名
		if part == "" {
			return false
		}
		// 不允许 . 或 ..
		if part == "." || part == ".." {
			return false
		}
		// 不允许特殊字符
		if strings.ContainsAny(part, `<>:"|?*`) {
			return false
		}
	}

	return true
}

// 验证文件类型
func isAllowedImageType(contentType string) bool {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
	return allowedTypes[contentType]
}

// GetImageInfo 获取图片信息
func (h *StaticFileHandler) GetImageInfo(c *gin.Context) {
	filename := c.Param("filename")

	// 构建文件路径
	filePath := filepath.Join("static", "images", filename)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "图片不存在",
		})
		return
	}

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取文件信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filename":   filename,
		"size":       fileInfo.Size(),
		"url":        fmt.Sprintf("/static/images/%s", filename),
		"created_at": fileInfo.ModTime(),
	})
}
