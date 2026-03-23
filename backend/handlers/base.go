package handlers

import (
	"ai-model-papers-backend/repository"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"unicode"
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
// 如果消息包含非 ASCII 字符，会进行 base64 编码并在 X-Error-Message-Encoding 头中标记
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")

	// 使用 strings.IndexFunc 检查是否包含非 ASCII 字符
	// 如果找到第一个非 ASCII 字符，返回其索引；否则返回 -1
	if strings.IndexFunc(message, func(r rune) bool { return r > unicode.MaxASCII }) >= 0 {
		// base64 编码
		encoded := base64.StdEncoding.EncodeToString([]byte(message))
		w.Header().Set("X-Error-Message", encoded)
		w.Header().Set("X-Error-Message-Encoding", "base64")
	} else {
		// 纯 ASCII，直接传递
		w.Header().Set("X-Error-Message", message)
	}

	w.WriteHeader(status)
	// 错误时返回空对象或简单结构
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
