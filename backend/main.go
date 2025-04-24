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

	// 创建静态文件处理器
	staticFileHandler := handlers.NewStaticFileHandler()

	// 设置食物分析处理器用于其他处理器
	handlers.SetFoodAnalysisHandler(foodAnalysisHandler)

	// 设置静态文件服务
	r.Static("/static", "./static")

	// 设置下载页面
	r.Static("/download", "./download")

	// 设置下载页面根路径
	r.GET("/download", func(c *gin.Context) {
		c.File("./download/index.html")
	})

	// API路由
	api := r.Group("/api")
	{
		// 认证相关路由
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
		api.POST("/verify-email", handlers.VerifyEmail)
		api.POST("/resend-verification", handlers.ResendVerification)

		// 静态文件路由（公开访问）
		api.GET("/image/:filename", staticFileHandler.GetImageInfo)

		// 用户物品相关路由（公开访问）
		api.GET("/user-items/:user_id", handlers.GetUserItems)
		api.GET("/user-items/:user_id/:item_id", handlers.GetUserItemDetails)
		api.GET("/user-items/:user_id/source", handlers.GetUserItemsBySource)
		api.GET("/user-items/:user_id/count", handlers.GetUserItemCount)
		api.GET("/user-items/:user_id/check/:item_id", handlers.CheckUserHasItem)
		api.POST("/user-items/:user_id", handlers.CreateUserItem)
		api.PUT("/user-items/:user_id/:item_id/quantity", handlers.UpdateUserItemQuantity)
		api.DELETE("/user-items/:user_id/:item_id", handlers.DeleteUserItem)

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

			// 用户健康状态路由
			authorized.POST("/health-states", handlers.CreateUserHealthStateHandler)
			authorized.GET("/health-states", handlers.GetUserHealthStatesHandler)
			authorized.GET("/health-states/:id", handlers.GetUserHealthStateHandler)
			authorized.PUT("/health-states/:id", handlers.UpdateUserHealthStateHandler)
			authorized.DELETE("/health-states/:id", handlers.DeleteUserHealthStateHandler)
			authorized.GET("/health-states/latest", handlers.GetLatestUserHealthStateHandler)

			// 打卡相关路由
			authorized.POST("/check-in", handlers.HandleCheckIn)        // 用户打卡
			authorized.GET("/check-in/today", handlers.GetTodayCheckIn) // 获取今日打卡状态
			authorized.GET("/check-ins", handlers.GetUserCheckIns)      // 获取用户所有打卡记录

			// 物品查询相关路由（所有用户可访问）
			authorized.GET("/items", handlers.GetItems)
			authorized.GET("/items/:id", handlers.GetItem)
			authorized.GET("/items/search", handlers.SearchItems)
			authorized.GET("/items/sources", handlers.GetItemSources)
		}

		// 管理员路由
		admin := authorized.Group("/admin")
		admin.Use(handlers.AdminAuthMiddleware())
		{
			// 静态文件上传路由（仅管理员可访问）
			admin.POST("/upload/image", staticFileHandler.UploadImage)

			admin.GET("/users", handlers.GetUsers)
			admin.GET("/stats", handlers.GetStats)

			// 物品管理路由（仅管理员可访问）
			admin.POST("/items", handlers.CreateItem)
			admin.PUT("/items/:id", handlers.UpdateItem)
			admin.DELETE("/items/:id", handlers.DeleteItem)
		}

		// 应用更新相关路由
		appUpdates := api.Group("/app-updates")
		{
			// 需要管理员权限的路由
			admin := appUpdates.Group("")
			admin.Use(handlers.AuthMiddleware(), handlers.AdminAuthMiddleware())
			{
				admin.POST("/upload", handlers.UploadAppUpdate)
				admin.DELETE("/:id", handlers.DeleteUpdate)
			}

			// 公开路由
			appUpdates.GET("", handlers.GetAppUpdates)
			appUpdates.GET("/check", handlers.CheckUpdate)
			appUpdates.GET("/download/:platform/:version", handlers.DownloadUpdate)
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
