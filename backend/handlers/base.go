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

// JSONResponse 返回JSON响应
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
