package handler

import (
	"net/http"
	"strconv"

	"modellens/backend/internal/repository"

	"github.com/gin-gonic/gin"
)

// AnalysisHandler 评测数据处理器
type AnalysisHandler struct {
	repo repository.AnalysisRepository
}

// NewAnalysisHandler 创建处理器
func NewAnalysisHandler(repo repository.AnalysisRepository) *AnalysisHandler {
	return &AnalysisHandler{repo: repo}
}

// ListArtificialAnalysis 获取评测数据列表
func (h *AnalysisHandler) ListArtificialAnalysis(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// 获取数据
	analysis, total, err := h.repo.List(c, int32(page), int32(limit))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "success",
		"data":    analysis,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

// GetArtificialAnalysisBySlug 通过 slug 获取评测数据
func (h *AnalysisHandler) GetArtificialAnalysisBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "slug is required",
		})
		return
	}

	analysis, err := h.repo.GetBySlug(c, slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "success",
		"data":    analysis,
	})
}
