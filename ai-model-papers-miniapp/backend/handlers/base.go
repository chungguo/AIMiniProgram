package handlers

import (
	"ai-model-papers-backend/models"
	"ai-model-papers-backend/repository"
	"encoding/json"
	"net/http"
)

var (
	modelRepo repository.ModelRepository
	paperRepo repository.PaperRepository
)

// InitRepositories 初始化仓库（在 main.go 中调用）
func InitRepositories(factory *repository.RepositoryFactory) {
	modelRepo = factory.GetModelRepository()
	paperRepo = factory.GetPaperRepository()
}

// JSONResponse 返回JSON响应
func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// LoadData 兼容旧代码，现在数据通过 Repository 加载
// 这个函数保留用于向后兼容，实际数据加载在 RepositoryFactory 中完成
func LoadData() error {
	// 数据现在通过 InitRepositories 注入
	return nil
}

// GetModelsData 兼容旧代码
func GetModelsData() *models.ModelsData {
	return nil // 不再使用
}

// GetPapersData 兼容旧代码
func GetPapersData() *models.PapersData {
	return nil // 不再使用
}
