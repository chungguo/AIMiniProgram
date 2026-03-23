package routers

import (
	"ai-model-papers-backend/handlers"
	"ai-model-papers-backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine) {
	// 全局中间件
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS())
	r.Use(middleware.AntiScraping())
	r.Use(middleware.RateLimitMiddleware(middleware.DefaultRateLimiterConfig))

	// API 路由组
	api := r.Group("/api")
	{
		// 健康检查
		api.GET("/health", gin.WrapF(handlers.HealthCheck))

		// 模型相关路由
		api.GET("/models", gin.WrapF(handlers.GetModels))
		api.GET("/models/families", gin.WrapF(handlers.GetFamilies))
		api.GET("/models/family/:family", gin.WrapF(handlers.GetFamilyModels))
		api.POST("/models/compare", middleware.StrictRateLimitMiddleware(), gin.WrapF(handlers.CompareModels))
		api.GET("/models/meta/comparison-categories", gin.WrapF(handlers.GetComparisonCategories))
		api.GET("/models/detail/:id", gin.WrapF(handlers.GetModelByID))

		// 论文相关路由
		api.GET("/papers", gin.WrapF(handlers.GetPapers))
		api.GET("/papers/featured/latest", gin.WrapF(handlers.GetLatestPapers))
		api.GET("/papers/detail/:id", gin.WrapF(handlers.GetPaperByID))

		// ArtificialAnalysis 评测数据路由
		api.GET("/analysis/artificialanalysis", gin.WrapF(handlers.GetArtificialAnalysis))
		api.GET("/analysis/artificialanalysis/:slug", gin.WrapF(handlers.GetArtificialAnalysisBySlug))
		api.GET("/models/analysis/:id", gin.WrapF(handlers.GetModelWithAnalysis))
	}
}
