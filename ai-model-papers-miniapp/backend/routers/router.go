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

		// 模型相关路由 - 注意顺序：具体路由在前，通配路由在后
		api.GET("/models", gin.WrapF(handlers.GetModels))
		api.GET("/models/providers", gin.WrapF(handlers.GetProviders))
		api.POST("/models/compare", middleware.StrictRateLimitMiddleware(), gin.WrapF(handlers.CompareModels))
		api.GET("/models/meta/comparison-categories", gin.WrapF(handlers.GetComparisonCategories))
		api.GET("/models/detail/:id", gin.WrapF(handlers.GetModelByID))

		// 论文相关路由
		api.GET("/papers", gin.WrapF(handlers.GetPapers))
		api.GET("/papers/categories", gin.WrapF(handlers.GetPaperCategories))
		api.GET("/papers/featured/latest", gin.WrapF(handlers.GetLatestPapers))
		api.GET("/papers/detail/:id", gin.WrapF(handlers.GetPaperByID))
	}
}
