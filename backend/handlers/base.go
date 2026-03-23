package handlers

import (
	"ai-model-papers-backend/repository"
	"encoding/json"
	"net/http"
)

var (
	modelRepo              repository.ModelRepository
	paperRepo              repository.PaperRepository
	artificialAnalysisRepo repository.ArtificialAnalysisRepository
)

// InitRepositories 初始化仓库（在 main.go 中调用）
func InitRepositories(factory *repository.RepositoryFactory) {
	modelRepo = factory.GetModelRepository()
	paperRepo = factory.GetPaperRepository()
	artificialAnalysisRepo = factory.GetArtificialAnalysisRepository()
}

// JSONResponse 返回JSON响应（成功时使用）
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// ErrorResponse 返回错误响应，错误信息放入响应头
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Error-Message", message)
	w.WriteHeader(status)
	// 错误时返回空对象或简单结构
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
