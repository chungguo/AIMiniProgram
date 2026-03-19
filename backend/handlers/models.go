package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"ai-model-papers-backend/models"
)

// GetModels 获取所有模型（支持筛选、排序、分页）
func GetModels(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// 解析筛选参数
	filter := &models.ModelFilter{}
	if family := query.Get("family"); family != "" {
		filter.Family = family
	}
	if v := query.Get("hasAttachment"); v != "" {
		b := v == "true"
		filter.HasAttachment = &b
	}
	if v := query.Get("hasReasoning"); v != "" {
		b := v == "true"
		filter.HasReasoning = &b
	}
	if v := query.Get("hasToolCall"); v != "" {
		b := v == "true"
		filter.HasToolCall = &b
	}
	if v := query.Get("openWeights"); v != "" {
		b := v == "true"
		filter.OpenWeights = &b
	}
	if v := query.Get("minContext"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.MinContext = n
		}
	}
	if v := query.Get("maxCostInput"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filter.MaxCostInput = f
		}
	}

	// 解析排序参数
	sort := &models.ModelSort{
		Field: query.Get("sortBy"),
		Order: query.Get("sortOrder"),
	}
	if sort.Field == "" {
		sort.Field = "name"
	}
	if sort.Order == "" {
		sort.Order = "asc"
	}

	// 解析分页参数
	page := 1
	limit := 20
	if v := query.Get("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}
	if v := query.Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}

	// 处理搜索（覆盖其他筛选）
	search := strings.TrimSpace(query.Get("search"))
	if search != "" {
		modelsList, err := modelRepo.Search(strings.ToLower(search))
		if err != nil {
			JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
			return
		}
		JSONResponse(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"data":    modelsList,
		})
		return
	}

	// 获取分页数据
	modelsList, total, err := modelRepo.GetAll(filter, sort, page, limit)
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	totalPages := (total + limit - 1) / limit

	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    modelsList,
		"pagination": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
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

// GetFamilies 获取所有模型家族（替代 GetProviders）
func GetFamilies(w http.ResponseWriter, r *http.Request) {
	families, err := modelRepo.GetFamilies()
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    families,
	})
}

// GetFamilyModels 获取指定家族的所有模型
func GetFamilyModels(w http.ResponseWriter, r *http.Request) {
	family := r.URL.Path[len("/api/models/family/"):]

	modelsList, err := modelRepo.GetByFamily(family)
	if err != nil {
		JSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    modelsList,
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

// GetProviders 兼容旧接口 - 返回家族列表作为提供商
func GetProviders(w http.ResponseWriter, r *http.Request) {
	GetFamilies(w, r)
}

// GetProviderModels 兼容旧接口
func GetProviderModels(w http.ResponseWriter, r *http.Request) {
	GetFamilyModels(w, r)
}
