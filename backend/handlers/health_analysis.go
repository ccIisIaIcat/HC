package handlers

import (
	"backend/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// HealthAnalysisHandler 处理健康分析相关请求
type HealthAnalysisHandler struct {
	DB        *gorm.DB
	OpenAIKey string
	Model     string
}

// NewHealthAnalysisHandler 创建健康分析处理器
func NewHealthAnalysisHandler(db *gorm.DB, openAIKey, model string) *HealthAnalysisHandler {
	return &HealthAnalysisHandler{
		DB:        db,
		OpenAIKey: openAIKey,
		Model:     model,
	}
}

// HealthAnalysisRequest 健康分析请求结构
type HealthAnalysisRequest struct {
	StartDate    string `json:"start_date" binding:"required"`
	EndDate      string `json:"end_date" binding:"required"`
	AnalysisType string `json:"analysis_type" binding:"required"`
	Description  string `json:"description"`
}

// HealthAnalysisResponse 健康分析响应结构
type HealthAnalysisResponse struct {
	Analysis string `json:"analysis"`
}

// RegisterHealthAnalysisRoutes 注册健康分析相关路由
func RegisterHealthAnalysisRoutes(router *gin.Engine, handler *HealthAnalysisHandler) {
	healthAnalysisGroup := router.Group("/api")
	healthAnalysisGroup.Use(AuthMiddleware())

	// 健康分析接口
	healthAnalysisGroup.POST("/health-analysis", handler.AnalyzeHealth)
}

// AnalyzeHealth 处理健康分析请求
func (h *HealthAnalysisHandler) AnalyzeHealth(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 解析请求参数
	var req HealthAnalysisRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("请求格式错误: %v", err)})
		return
	}

	// 验证日期格式
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "起始日期格式错误，应为YYYY-MM-DD"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束日期格式错误，应为YYYY-MM-DD"})
		return
	}

	// 结束日期需要包含当天的所有记录，所以增加一天
	endDate = endDate.AddDate(0, 0, 1)

	// 验证日期范围
	if startDate.After(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "起始日期不能晚于结束日期"})
		return
	}

	// 从数据库获取该用户在指定时间范围内的食物记录
	var foodRecords []models.FoodRecord
	result := h.DB.Where("user_id = ? AND record_time >= ? AND record_time < ?", userID, startDate, endDate).Find(&foodRecords)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取食物记录失败: %v", result.Error)})
		return
	}

	// 检查是否有记录
	if len(foodRecords) == 0 {
		c.JSON(http.StatusOK, gin.H{"analysis": "在所选时间范围内没有发现食物记录。请尝试扩大时间范围或添加新的食物记录。"})
		return
	}

	// 将食物记录转换为字符串
	recordsStr := formatFoodRecordsToString(foodRecords)

	// 获取分析类型名称
	analysisTypeName := getAnalysisTypeName(req.AnalysisType)

	// 构建提示词
	prompt := constructAnalysisPrompt(recordsStr, analysisTypeName, req.Description)

	// 调用OpenAI获取分析结果
	analysis, err := h.callOpenAI(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("调用AI分析失败: %v", err)})
		return
	}

	// 返回分析结果
	c.JSON(http.StatusOK, gin.H{"analysis": analysis})
}

// 将食物记录格式化为字符串
func formatFoodRecordsToString(records []models.FoodRecord) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("食物记录总数: %d\n\n", len(records)))

	for i, record := range records {
		builder.WriteString(fmt.Sprintf("记录 #%d:\n", i+1))
		builder.WriteString(fmt.Sprintf("- 时间: %s\n", record.RecordTime.Format("2006-01-02 15:04:05")))
		builder.WriteString(fmt.Sprintf("- 食物名称: %s\n", record.FoodName))
		builder.WriteString(fmt.Sprintf("- 重量: %.1f克\n", record.Weight))
		builder.WriteString(fmt.Sprintf("- 热量: %.1f卡路里\n", record.Calories))
		builder.WriteString(fmt.Sprintf("- 蛋白质: %.1f克\n", record.Protein))
		builder.WriteString(fmt.Sprintf("- 总脂肪: %.1f克\n", record.TotalFat))
		builder.WriteString(fmt.Sprintf("- 碳水化合物: %.1f克\n", record.Carbohydrates))
		builder.WriteString(fmt.Sprintf("- 膳食纤维: %.1f克\n", record.Fiber))
		builder.WriteString(fmt.Sprintf("- 糖: %.1f克\n", record.Sugar))
		builder.WriteString(fmt.Sprintf("- 餐食类型: %s\n", record.MealType))

		// 添加其他详细营养素信息
		builder.WriteString("- 脂肪详情: ")
		builder.WriteString(fmt.Sprintf("饱和脂肪 %.1f克, ", record.SaturatedFat))
		builder.WriteString(fmt.Sprintf("反式脂肪 %.1f克, ", record.TransFat))
		builder.WriteString(fmt.Sprintf("不饱和脂肪 %.1f克\n", record.UnsaturatedFat))

		// 添加维生素信息
		builder.WriteString("- 维生素: ")
		builder.WriteString(fmt.Sprintf("A %.1fμg, ", record.VitaminA))
		builder.WriteString(fmt.Sprintf("C %.1fmg, ", record.VitaminC))
		builder.WriteString(fmt.Sprintf("D %.1fμg, ", record.VitaminD))
		builder.WriteString(fmt.Sprintf("B1 %.1fmg, ", record.VitaminB1))
		builder.WriteString(fmt.Sprintf("B2 %.1fmg\n", record.VitaminB2))

		// 添加矿物质信息
		builder.WriteString("- 矿物质: ")
		builder.WriteString(fmt.Sprintf("钙 %.1fmg, ", record.Calcium))
		builder.WriteString(fmt.Sprintf("铁 %.1fmg, ", record.Iron))
		builder.WriteString(fmt.Sprintf("钠 %.1fmg, ", record.Sodium))
		builder.WriteString(fmt.Sprintf("钾 %.1fmg\n", record.Potassium))

		// 添加备注（如果有）
		if record.Notes != "" {
			builder.WriteString(fmt.Sprintf("- 备注: %s\n", record.Notes))
		}

		builder.WriteString("\n")
	}

	return builder.String()
}

// 获取分析类型的中文名称
func getAnalysisTypeName(analysisType string) string {
	switch analysisType {
	case "comprehensive":
		return "全面分析"
	case "nutrition":
		return "营养均衡分析"
	case "calories":
		return "热量摄入分析"
	case "macros":
		return "宏量营养素分析"
	default:
		return "全面分析"
	}
}

// 构建分析提示词
func constructAnalysisPrompt(recordsStr, analysisType, userDescription string) string {
	// 获取用户最新的健康状态
	var latestHealthState models.UserHealthState
	result := models.DB.Order("created_at desc").First(&latestHealthState)
	healthStateStr := ""
	if result.Error == nil {
		healthStateStr = fmt.Sprintf(`
用户最新健康状态：
- 体重: %.1f kg
- 身高: %.1f cm
- BMI: %.1f
- 心率: %d 次/分
- 血糖: %.1f mmol/L
- 记录时间: %s
`,
			latestHealthState.Weight,
			latestHealthState.Height,
			latestHealthState.BMI,
			latestHealthState.HeartRate,
			latestHealthState.FastingGlucose,
			latestHealthState.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	prompt := fmt.Sprintf(`你是一名专业的营养学家和健康顾问。请基于以下用户的饮食记录数据，进行%s。

用户描述：%s

用户当前基本状况：%s

饮食记录数据：
%s

请提供分析结果，遵循以下要求：
1. 首先给出一句20字以内的幽默、略带戏谑或鼓励的话。
2. 然后基于数据进行专业、客观的分析，重点关注%s方面。
3. 必要时可以提供一些改进建议，但不要过于严厉，保持积极鼓励的态度。
4. 整体分析不超过400字。
`, analysisType, userDescription, healthStateStr, recordsStr, getAnalysisTypeDetails(analysisType))

	return prompt
}

// 获取分析类型的详细关注点
func getAnalysisTypeDetails(analysisType string) string {
	switch analysisType {
	case "全面分析":
		return "整体营养平衡、热量摄入、饮食多样性和规律性"
	case "营养均衡分析":
		return "各种维生素、矿物质等微量元素的摄入是否充足、均衡"
	case "热量摄入分析":
		return "每日热量摄入总量、分布以及与标准需求的对比"
	case "宏量营养素分析":
		return "蛋白质、脂肪和碳水化合物的摄入比例和质量"
	default:
		return "整体营养平衡、热量摄入、饮食多样性和规律性"
	}
}

// 调用OpenAI进行分析
func (h *HealthAnalysisHandler) callOpenAI(prompt string) (string, error) {
	// 构建API请求
	url := "https://api.openai-proxy.org/v1/chat/completions"

	// 构建请求体
	requestBody := map[string]interface{}{
		"model": h.Model,
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": "你是一名专业的营养学家和健康顾问，会提供简洁、准确、有用的健康分析建议。",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens": 1000,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("构建请求体失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.OpenAIKey))

	log.Printf("发送健康分析请求到OpenAI，模型: %s", h.Model)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %v", err)
	}

	log.Printf("OpenAI响应状态码: %d", resp.StatusCode)

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败，状态码: %d，响应: %s", resp.StatusCode, string(bodyBytes))
	}

	// 解析响应
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return "", fmt.Errorf("解析响应JSON失败: %v, 响应内容: %s", err, string(bodyBytes))
	}

	// 检查API错误
	if response.Error != nil {
		return "", fmt.Errorf("API返回错误: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("API响应中没有choices")
	}

	// 返回分析结果
	return response.Choices[0].Message.Content, nil
}
