package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"backend/models"

	"github.com/gin-gonic/gin"
)

type FoodAnalysisHandler struct {
	OpenAIKey string
	Model     string
}

func NewFoodAnalysisHandler(openAIKey string, model string) *FoodAnalysisHandler {
	if openAIKey == "" {
		log.Fatal("错误：OPENAI_API_KEY未设置")
	}
	if model == "" {
		model = "gpt-4-vision-preview"
		log.Printf("警告：OPENAI_API_MODEL未设置，使用默认模型：%s", model)
	}
	return &FoodAnalysisHandler{
		OpenAIKey: openAIKey,
		Model:     model,
	}
}

// UploadAndAnalyze 处理图片上传和分析
func (h *FoodAnalysisHandler) UploadAndAnalyze(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("获取上传文件失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无法获取上传的图片"})
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		log.Printf("不支持的文件类型: %s", ext)
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持JPG和PNG格式的图片"})
		return
	}

	// 检查文件大小（限制为10MB）
	if file.Size > 10*1024*1024 {
		log.Printf("文件太大: %d bytes", file.Size)
		c.JSON(http.StatusBadRequest, gin.H{"error": "图片大小不能超过10MB"})
		return
	}

	// 读取文件内容
	openedFile, err := file.Open()
	if err != nil {
		log.Printf("打开文件失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取图片"})
		return
	}
	defer openedFile.Close()

	// 读取文件内容到内存
	imageBytes, err := io.ReadAll(openedFile)
	if err != nil {
		log.Printf("读取文件内容失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法读取图片内容"})
		return
	}

	// 获取图片描述（如果有）
	imageDescription := c.PostForm("image_description")

	// 调用GPT-4 Vision API分析图片
	var analysis *models.FoodAnalysis
	if imageDescription != "" {
		analysis, err = h.analyzeImageWithDescription(imageBytes, imageDescription)
	} else {
		analysis, err = h.analyzeImage(imageBytes)
	}

	if err != nil {
		log.Printf("分析图片失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("图片分析失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// analyzeImage 使用GPT-4 Vision API分析图片
func (h *FoodAnalysisHandler) analyzeImage(imageBytes []byte) (*models.FoodAnalysis, error) {
	// 构建API请求
	url := "https://api.openai-proxy.org/v1/chat/completions"
	prompt := `分析这张图片中的食物。请提供以下信息，必须严格按照指定的JSON格式返回：
{
    "hasFood": true,  // 布尔值，表示是否包含食物
    "foodType": "食物名称",  // 字符串，食物的类型
    "weight": 100,  // 数字，估计重量（克）
    "nutrition": {
        "calories": 0,  // 数字，热量（卡路里）
        "protein": 0,  // 数字，蛋白质（克）
        "totalFat": 0,  // 数字，总脂肪（克）
        "saturatedFat": 0,  // 数字，饱和脂肪（克）
        "transFat": 0,  // 数字，反式脂肪（克）
        "unsaturatedFat": 0,  // 数字，不饱和脂肪（克）
        "carbohydrates": 0,  // 数字，碳水化合物（克）
        "sugar": 0,  // 数字，糖分（克）
        "fiber": 0,  // 数字，膳食纤维（克）
        "vitamins": {
            "vitaminA": 0,  // 数字，维生素A（毫克）
            "vitaminC": 0,  // 数字，维生素C（毫克）
            "vitaminD": 0,  // 数字，维生素D（微克）
            "vitaminE": 0,  // 数字，维生素E（毫克）
            "vitaminK": 0,  // 数字，维生素K（微克）
            "vitaminB": {
                "b1": 0,  // 数字，维生素B1（毫克）
                "b2": 0,  // 数字，维生素B2（毫克）
                "b6": 0,  // 数字，维生素B6（毫克）
                "b12": 0  // 数字，维生素B12（微克）
            }
        },
        "minerals": {
            "calcium": 0,  // 数字，钙（毫克）
            "iron": 0,  // 数字，铁（毫克）
            "sodium": 0,  // 数字，钠（毫克）
            "potassium": 0,  // 数字，钾（毫克）
            "zinc": 0,  // 数字，锌（毫克）
            "magnesium": 0  // 数字，镁（毫克）
        }
    }
}`

	// 将图片转换为base64
	base64Image := base64.StdEncoding.EncodeToString(imageBytes)

	// 构建请求体
	requestBody := map[string]interface{}{
		"model": h.Model,
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": prompt,
					},
					{
						"type": "image_url",
						"image_url": map[string]interface{}{
							"url": fmt.Sprintf("data:image/jpeg;base64,%s", base64Image),
						},
					},
				},
			},
		},
		"max_tokens": 4000,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("构建请求体失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.OpenAIKey))

	log.Printf("发送请求到 OpenAI API，URL: %s", url)
	log.Printf("使用模型: %s", h.Model)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	log.Printf("API响应状态码: %d", resp.StatusCode)

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d，响应: %s", resp.StatusCode, string(bodyBytes))
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
		return nil, fmt.Errorf("解析响应JSON失败: %v, 响应内容: %s", err, string(bodyBytes))
	}

	// 检查API错误
	if response.Error != nil {
		return nil, fmt.Errorf("API返回错误: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("API响应中没有choices")
	}

	// 解析GPT返回的JSON到FoodAnalysis结构体
	content := response.Choices[0].Message.Content
	log.Printf("API返回的原始内容: %s", content)

	// 如果内容被Markdown代码块包裹，提取其中的JSON
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
		if idx := strings.Index(content, "```"); idx != -1 {
			content = content[:idx]
		}
	} else if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
		if idx := strings.Index(content, "```"); idx != -1 {
			content = content[:idx]
		}
	}
	content = strings.TrimSpace(content)

	var analysis models.FoodAnalysis
	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		return nil, fmt.Errorf("解析分析结果失败: %v, Content: %s", err, content)
	}

	return &analysis, nil
}

// analyzeImageWithDescription 使用GPT-4 Vision API结合图片描述分析图片
func (h *FoodAnalysisHandler) analyzeImageWithDescription(imageBytes []byte, description string) (*models.FoodAnalysis, error) {
	// 构建API请求
	url := "https://api.openai-proxy.org/v1/chat/completions"
	prompt := fmt.Sprintf(`分析这张图片中的食物，
	1、分析时尤其注意热量和质量的比例关系，
	2、可以先考虑食物中的主要成分和高热量成分的质量，再按照比例推算热量，
	3、用户描述: %s
	4、用户可能在上传的图片包含多个食物，在名称中需一一列举，在考虑营养成分时整体考虑，考虑食物的比例，各自的营养成分，给出整体评价
	请提供以下信息，必须严格按照指定的JSON格式返回：
{
    "hasFood": true,  // 布尔值，表示是否包含食物
    "foodType": "食物名称",  // 字符串，食物的类型
    "weight": 100,  // 数字，估计重量（克）
    "nutrition": {
        "calories": 0,  // 数字，热量（卡路里）
        "protein": 0,  // 数字，蛋白质（克）
        "totalFat": 0,  // 数字，总脂肪（克）
        "saturatedFat": 0,  // 数字，饱和脂肪（克）
        "transFat": 0,  // 数字，反式脂肪（克）
        "unsaturatedFat": 0,  // 数字，不饱和脂肪（克）
        "carbohydrates": 0,  // 数字，碳水化合物（克）
        "sugar": 0,  // 数字，糖分（克）
        "fiber": 0,  // 数字，膳食纤维（克）
        "vitamins": {
            "vitaminA": 0,  // 数字，维生素A（毫克）
            "vitaminC": 0,  // 数字，维生素C（毫克）
            "vitaminD": 0,  // 数字，维生素D（微克）
            "vitaminE": 0,  // 数字，维生素E（毫克）
            "vitaminK": 0,  // 数字，维生素K（微克）
            "vitaminB": {
                "b1": 0,  // 数字，维生素B1（毫克）
                "b2": 0,  // 数字，维生素B2（毫克）
                "b6": 0,  // 数字，维生素B6（毫克）
                "b12": 0  // 数字，维生素B12（微克）
            }
        },
        "minerals": {
            "calcium": 0,  // 数字，钙（毫克）
            "iron": 0,  // 数字，铁（毫克）
            "sodium": 0,  // 数字，钠（毫克）
            "potassium": 0,  // 数字，钾（毫克）
            "zinc": 0,  // 数字，锌（毫克）
            "magnesium": 0  // 数字，镁（毫克）
        }
    }
}`, description)
	fmt.Println(prompt)
	// 将图片转换为base64
	base64Image := base64.StdEncoding.EncodeToString(imageBytes)

	// 构建请求体
	requestBody := map[string]interface{}{
		"model": h.Model,
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": prompt,
					},
					{
						"type": "image_url",
						"image_url": map[string]interface{}{
							"url": fmt.Sprintf("data:image/jpeg;base64,%s", base64Image),
						},
					},
				},
			},
		},
		"max_tokens": 4000,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("构建请求体失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.OpenAIKey))

	log.Printf("发送请求到 OpenAI API，URL: %s", url)
	log.Printf("使用模型: %s，带有图片描述", h.Model)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	log.Printf("API响应状态码: %d", resp.StatusCode)

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d，响应: %s", resp.StatusCode, string(bodyBytes))
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
		return nil, fmt.Errorf("解析响应JSON失败: %v, 响应内容: %s", err, string(bodyBytes))
	}

	// 检查API错误
	if response.Error != nil {
		return nil, fmt.Errorf("API返回错误: %s", response.Error.Message)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("API响应中没有choices")
	}

	// 解析GPT返回的JSON到FoodAnalysis结构体
	content := response.Choices[0].Message.Content
	log.Printf("API返回的原始内容: %s", content)

	// 如果内容被Markdown代码块包裹，提取其中的JSON
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
		if idx := strings.Index(content, "```"); idx != -1 {
			content = content[:idx]
		}
	} else if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
		if idx := strings.Index(content, "```"); idx != -1 {
			content = content[:idx]
		}
	}
	content = strings.TrimSpace(content)

	var analysis models.FoodAnalysis
	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		return nil, fmt.Errorf("解析分析结果失败: %v, Content: %s", err, content)
	}

	return &analysis, nil
}

// AnalyzeImageBytes 提供公共接口来分析图片字节数据
func (h *FoodAnalysisHandler) AnalyzeImageBytes(imageBytes []byte) (*models.FoodAnalysis, error) {
	return h.analyzeImage(imageBytes)
}

// AnalyzeImageBytesWithDescription 提供公共接口来分析图片字节数据，并结合描述信息
func (h *FoodAnalysisHandler) AnalyzeImageBytesWithDescription(imageBytes []byte, description string) (*models.FoodAnalysis, error) {
	return h.analyzeImageWithDescription(imageBytes, description)
}
