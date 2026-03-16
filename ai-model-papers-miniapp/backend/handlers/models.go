package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

// GetModels 获取所有模型
func GetModels(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	capability := r.URL.Query().Get("capability")
	search := strings.ToLower(r.URL.Query().Get("search"))

	var result interface{}
	var err error

	// 根据参数选择查询方式
	switch {
	case provider != "":
		result, err = modelRepo.GetByProvider(provider)
	case capability != "":
		result, err = modelRepo.GetByCapability(capability)
	case search != "":
		result, err = modelRepo.Search(search)
	default:
		result, err = modelRepo.GetAll()
	}

	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    result,
	})
}

// GetModelByID 获取单个模型
func GetModelByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/models/detail/"):]

	model, err := modelRepo.GetByID(id)
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if model == nil {
		JSONResponse(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": "Model not found",
		})
		return
	}

	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    model,
	})
}

// GetProviders 获取所有提供商
func GetProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := modelRepo.GetProviders()
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    providers,
	})
}

// CompareRequest 对比请求
type CompareRequest struct {
	IDs []string `json:"ids"`
}

// CompareModels 对比多个模型
func CompareModels(w http.ResponseWriter, r *http.Request) {
	var req CompareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	if len(req.IDs) < 2 {
		JSONResponse(w, http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Please provide at least 2 model IDs",
		})
		return
	}

	result := make([]interface{}, 0, len(req.IDs))
	for _, id := range req.IDs {
		model, err := modelRepo.GetByID(id)
		if err != nil {
			JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		if model != nil {
			result = append(result, model)
		}
	}

	categories, err := modelRepo.GetComparisonCategories()
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"models":               result,
			"comparisonCategories": categories,
		},
	})
}

// GetComparisonCategories 获取对比类别
func GetComparisonCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := modelRepo.GetComparisonCategories()
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    categories,
	})
}
