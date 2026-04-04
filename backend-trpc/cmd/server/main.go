package main

import (
	"context"
	"fmt"
	"log"

	"aiminiprogram/backend-trpc/internal/service"

	// 导入生成的 proto 代码
	pbmodel "aiminiprogram/proto/model/v1"
	pbpaper "aiminiprogram/proto/paper/v1"
	pbanalysis "aiminiprogram/proto/analysis/v1"

	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/filter"
	"trpc.group/trpc-go/trpc-go/log"
	"trpc.group/trpc-go/trpc-go/recovery"
	"trpc.group/trpc-go/trpc-go/server"
)

func main() {
	// 创建 trpc 服务
	s, err := trpc.NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// 注册全局拦截器
	filter.Register(
		"recovery",
		recovery.ServerFilter(),
	)

	// 初始化服务实现
	modelSvc := service.NewModelService()
	paperSvc := service.NewPaperService()
	analysisSvc := service.NewAnalysisService()

	// 注册服务
	pbmodel.RegisterModelService(s, modelSvc)
	pbpaper.RegisterPaperService(s, paperSvc)
	pbanalysis.RegisterAnalysisService(s, analysisSvc)

	log.Info("Server starting...")
	log.Info("Registered services:")
	log.Info("  - ModelService")
	log.Info("  - PaperService")
	log.Info("  - AnalysisService")

	// 启动服务
	if err := s.Serve(); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
