package main

import (
	"log"
	"os"

	"backend/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Printf("警告：未找到.env文件，将使用系统环境变量")
	}

	// 设置为生产模式
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库
	handlers.InitDB()

	// 获取OpenAI API密钥
	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		log.Fatal("错误：OPENAI_API_KEY 未设置")
	}

	// 获取OpenAI模型名称
	model := os.Getenv("OPENAI_API_MODEL")
	if model == "" {
		model = "gpt-4o" // 默认模型
		log.Printf("警告：OPENAI_API_MODEL 未设置，使用默认模型：%s", model)
	}

	// 创建Gin路由
	r := gin.Default()

	// 设置CORS中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 设置受信任的代理
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// 创建处理器
	foodAnalysisHandler := handlers.NewFoodAnalysisHandler(openAIKey, model)

	// 创建健康分析处理器
	healthAnalysisHandler := handlers.NewHealthAnalysisHandler(handlers.GetDB(), openAIKey, model)

	// 设置食物分析处理器用于其他处理器
	handlers.SetFoodAnalysisHandler(foodAnalysisHandler)

	// 设置静态文件服务
	r.Static("/static", "./frontend/build/static")
	r.StaticFile("/", "./frontend/build/index.html")
	r.StaticFile("/favicon.ico", "./frontend/build/favicon.ico")

	// API路由
	api := r.Group("/api")
	{
		// 认证相关路由
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
		api.POST("/verify-email", handlers.VerifyEmail)
		api.POST("/resend-verification", handlers.ResendVerification)

		// 需要认证的路由
		authorized := api.Group("/")
		authorized.Use(handlers.AuthMiddleware())
		{
			authorized.POST("/analyze-food", foodAnalysisHandler.UploadAndAnalyze)
			authorized.GET("/me", handlers.GetCurrentUser)

			// 食物记录路由
			authorized.POST("/food-records", handlers.CreateFoodRecordHandler)
			authorized.GET("/food-records", handlers.GetFoodRecordsHandler)
			authorized.GET("/food-records/:id", handlers.GetFoodRecordHandler)
			authorized.PUT("/food-records/:id", handlers.UpdateFoodRecordHandler)
			authorized.DELETE("/food-records/:id", handlers.DeleteFoodRecordHandler)
			authorized.POST("/analyze-and-save", handlers.AnalyzeAndSaveFoodHandler)

			// 健康分析路由
			authorized.POST("/health-analysis", healthAnalysisHandler.AnalyzeHealth)
		}

		// 管理员路由
		admin := authorized.Group("/admin")
		admin.Use(handlers.AdminAuthMiddleware())
		{
			admin.GET("/users", handlers.GetUsers)
			admin.GET("/stats", handlers.GetStats)
		}
	}

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("服务器启动在 http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
