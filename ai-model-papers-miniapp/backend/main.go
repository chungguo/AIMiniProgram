package main

import (
	"ai-model-papers-backend/handlers"
	"ai-model-papers-backend/repository"
	"ai-model-papers-backend/routers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 设置生产模式
	// gin.SetMode(gin.ReleaseMode)

	// 创建仓库工厂（默认使用JSON，可通过 DATABASE_URL 环境变量切换到PostgreSQL）
	factory, err := repository.NewRepositoryFactory()
	if err != nil {
		log.Fatalf("Failed to initialize repositories: %v", err)
	}
	
	// 初始化 handlers
	handlers.InitRepositories(factory)
	log.Println("Repositories initialized successfully (JSON mode)")

	// 创建引擎
	r := gin.New()
	r.Use(gin.Recovery())

	// 注册路由（包含所有中间件）
	routers.SetupRoutes(r)

	port := ":3000"
	log.Printf("Server running on port %s", port)
	log.Println("Middleware enabled: SecurityHeaders, CORS, AntiScraping, RateLimit")
	log.Println("Data storage: JSON files (set DATABASE_URL env to use PostgreSQL)")
	
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
