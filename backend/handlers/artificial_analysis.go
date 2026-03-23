package handlers

import (
	"net/http"
)

// GetArtificialAnalysis 获取所有 ArtificialAnalysis 评测数据
func GetArtificialAnalysis(w http.ResponseWriter, r *http.Request) {
	analyses, err := artificialAnalysisRepo.GetAll()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, analyses)
}

// GetArtificialAnalysisBySlug 根据 slug 获取单个评测数据
func GetArtificialAnalysisBySlug(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Path[len("/api/analysis/artificialanalysis/"):]

	analysis, err := artificialAnalysisRepo.GetBySlug(slug)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if analysis == nil {
		ErrorResponse(w, http.StatusNotFound, "Analysis not found")
		return
	}

	JSONResponse(w, http.StatusOK, analysis)
}

// GetModelWithAnalysis 获取模型及其评测数据（通过 slug 关联）
func GetModelWithAnalysis(w http.ResponseWriter, r *http.Request) {
	modelID := r.URL.Path[len("/api/models/analysis/"):]

	// 先获取模型信息
	model, err := modelRepo.GetByID(modelID)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if model == nil {
		ErrorResponse(w, http.StatusNotFound, "Model not found")
		return
	}

	// 使用 model.name 作为 slug 查询评测数据
	analysis, err := artificialAnalysisRepo.GetBySlug(model.Name)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"model":    model,
		"analysis": analysis,
	})
}
