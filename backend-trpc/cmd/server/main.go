package main

import (
	"log"
	"net/http"

	"aiminiprogram/backend-trpc/internal/handler"
	"aiminiprogram/backend-trpc/internal/repository"

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
	}

	port := ":8000"
	log.Printf("tRPC-Go Compatible Server running on port %s", port)
	log.Println("Available endpoints:")
	log.Println("  GET  /api/health")
	log.Println("  GET  /api/models")
	log.Println("  GET  /api/papers")
	log.Println("  GET  /api/analysis/artificialanalysis")

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
