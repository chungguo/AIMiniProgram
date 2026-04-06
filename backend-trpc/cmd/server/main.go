package main

import (
	"log"
	"net/http"
	"os"

	"aiminiprogram/backend-trpc/internal/handler"
	"aiminiprogram/backend-trpc/internal/repository"
	"aiminiprogram/backend-trpc/internal/wechat"

	"github.com/gin-gonic/gin"
)

func main() {
	// 设置生产模式
	// gin.SetMode(gin.ReleaseMode)

	// 创建仓库
	modelRepo := repository.NewModelRepository()
	paperRepo := repository.NewPaperRepository()
	analysisRepo := repository.NewAnalysisRepository()

	// 创建 handlers
	modelHandler := handler.NewModelHandler(modelRepo)
	paperHandler := handler.NewPaperHandler(paperRepo)
	analysisHandler := handler.NewAnalysisHandler(analysisRepo)

	// 初始化微信客户端（如果配置了环境变量）
	var wechatHandler *handler.WeChatHandler
	if os.Getenv("WECHAT_APPID") != "" && os.Getenv("WECHAT_SECRET") != "" {
		wechatClient, err := wechat.NewClientFromEnv()
		if err != nil {
			log.Printf("Warning: WeChat client init failed: %v", err)
		} else {
			wechatHandler = handler.NewWeChatHandler(wechatClient)
			log.Println("✅ WeChat client initialized")
		}
	}

	// 创建引擎
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(CORS())

	// API 路由组
	api := r.Group("/api")
	{
		// 健康检查
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// 模型相关路由
		api.GET("/models", modelHandler.ListModels)
		api.GET("/models/families", modelHandler.ListFamilies)
		api.GET("/models/family/:family", modelHandler.GetFamilyModels)
		api.POST("/models/compare", modelHandler.CompareModels)
		api.GET("/models/meta/comparison-categories", modelHandler.GetComparisonCategories)
		api.GET("/models/detail/:id", modelHandler.GetModel)

		// 论文相关路由
		api.GET("/papers", paperHandler.ListPapers)
		api.GET("/papers/featured/latest", paperHandler.GetLatestPapers)
		api.GET("/papers/detail/:id", paperHandler.GetPaper)

		// 评测数据路由
		api.GET("/analysis/artificialanalysis", analysisHandler.ListArtificialAnalysis)
		api.GET("/analysis/artificialanalysis/:slug", analysisHandler.GetArtificialAnalysisBySlug)

		// ==================== PowerWeChat 路由 ====================
		if wechatHandler != nil {
			// 1. 小程序登录（高优先级）
			api.POST("/wechat/login", wechatHandler.MiniProgramLogin)

			// 2. 解密手机号（高优先级）
			api.POST("/wechat/phone", wechatHandler.DecryptPhoneNumber)

			// 3. 订阅消息（中优先级）
			api.POST("/wechat/subscribe-message", wechatHandler.SendSubscribeMessage)

			// 4. 微信支付（中优先级）
			api.POST("/wechat/pay/order", wechatHandler.CreateOrder)
			api.POST("/wechat/pay/notify", wechatHandler.PayNotify)

			// 5. 内容安全（低优先级）
			api.POST("/wechat/security/check-content", wechatHandler.CheckContent)
			api.POST("/wechat/security/check-image", wechatHandler.CheckImage)

			// 6. 二维码（低优先级）
			api.POST("/wechat/qrcode", wechatHandler.CreateQRCode)
		}
	}

	port := ":8000"
	log.Printf("🚀 AIMiniProgram Server running on port %s", port)
	log.Println("")
	log.Println("📚 Available endpoints:")
	log.Println("  GET  /api/health")
	log.Println("  GET  /api/models")
	log.Println("  GET  /api/papers")
	log.Println("  GET  /api/analysis/artificialanalysis")
	if wechatHandler != nil {
		log.Println("")
		log.Println("🔑 WeChat endpoints:")
		log.Println("  POST /api/wechat/login              - 小程序登录")
		log.Println("  POST /api/wechat/phone              - 解密手机号")
		log.Println("  POST /api/wechat/subscribe-message  - 发送订阅消息")
		log.Println("  POST /api/wechat/pay/order          - 创建支付订单")
		log.Println("  POST /api/wechat/security/check-content - 内容安全检查")
	}

	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// CORS 中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
