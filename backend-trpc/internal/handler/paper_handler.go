package handler

import (
	"net/http"
	"strconv"

	"aiminiprogram/backend-trpc/internal/repository"

	"github.com/gin-gonic/gin"
)

// PaperHandler 论文处理器
type PaperHandler struct {
	repo repository.PaperRepository
}

// NewPaperHandler 创建处理器
func NewPaperHandler(repo repository.PaperRepository) *PaperHandler {
	return &PaperHandler{repo: repo}
}

// ListPapers 获取论文列表
func (h *PaperHandler) ListPapers(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// 构建筛选条件
	filter := &repository.PaperFilter{
		Search: c.Query("search"),
		Author: c.Query("author"),
	}

	// 获取数据
	papers, total, err := h.repo.List(c, filter, int32(page), int32(limit))
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
		"data":    papers,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

// GetPaper 获取单个论文
func (h *PaperHandler) GetPaper(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "id is required",
		})
		return
	}

	paper, err := h.repo.GetByID(c, id)
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
		"data":    paper,
	})
}

// GetLatestPapers 获取最新论文
func (h *PaperHandler) GetLatestPapers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "3"))
	if limit < 1 || limit > 20 {
		limit = 3
	}

	papers, err := h.repo.GetLatest(c, int32(limit))
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
		"data":    papers,
	})
}
