package handler

import (
	"net/http"
	"strconv"

	"aiminiprogram/backend-trpc/internal/repository"

	"github.com/gin-gonic/gin"
)

// ModelHandler 模型处理器
type ModelHandler struct {
	repo repository.ModelRepository
}

// NewModelHandler 创建处理器
func NewModelHandler(repo repository.ModelRepository) *ModelHandler {
	return &ModelHandler{repo: repo}
}

// ListModels 获取模型列表
func (h *ModelHandler) ListModels(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// 构建筛选条件
	filter := &repository.ModelFilter{
		Family: c.Query("family"),
		Search: c.Query("search"),
	}

	if v := c.Query("hasAttachment"); v != "" {
		b := v == "true"
		filter.HasAttachment = &b
	}
	if v := c.Query("hasReasoning"); v != "" {
		b := v == "true"
		filter.HasReasoning = &b
	}
	if v := c.Query("hasToolCall"); v != "" {
		b := v == "true"
		filter.HasToolCall = &b
	}
	if v := c.Query("openWeights"); v != "" {
		b := v == "true"
		filter.OpenWeights = &b
	}
	if v := c.Query("minContext"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.MinContext = int32(n)
		}
	}

	// 获取数据
	models, total, err := h.repo.List(c, filter, int32(page), int32(limit))
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
		"data":    models,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + limit - 1) / limit,
		},
	})
}

// GetModel 获取单个模型
func (h *ModelHandler) GetModel(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "id is required",
		})
		return
	}

	model, err := h.repo.GetByID(c, id)
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
		"data":    model,
	})
}

// ListFamilies 获取家族列表
func (h *ModelHandler) ListFamilies(c *gin.Context) {
	families, err := h.repo.GetFamilies(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "success",
		"families": families,
	})
}

// GetFamilyModels 获取家族模型
func (h *ModelHandler) GetFamilyModels(c *gin.Context) {
	family := c.Param("family")
	if family == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "family is required",
		})
		return
	}

	models, err := h.repo.GetByFamily(c, family)
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
		"data":    models,
	})
}

// GetComparisonCategories 获取对比类别
func (h *ModelHandler) GetComparisonCategories(c *gin.Context) {
	categories := getComparisonCategories()
	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "success",
		"categories": categories,
	})
}

// CompareModels 对比模型
func (h *ModelHandler) CompareModels(c *gin.Context) {
	var req struct {
		Ids []string `json:"ids" binding:"required,min=2,max=5"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ids is required (2-5 models)",
		})
		return
	}

	models := make([]*repository.Model, 0, len(req.Ids))
	for _, id := range req.Ids {
		model, err := h.repo.GetByID(c, id)
		if err != nil {
			continue
		}
		models = append(models, model)
	}

	if len(models) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "at least 2 valid models required",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":             true,
		"message":             "success",
		"models":              models,
		"comparisonCategories": getComparisonCategories(),
	})
}

// ComparisonItem 对比项
type ComparisonItem struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Type string `json:"type"`
	Unit string `json:"unit,omitempty"`
}

// ComparisonCategory 对比类别
type ComparisonCategory struct {
	Key   string           `json:"key"`
	Name  string           `json:"name"`
	Items []ComparisonItem `json:"items"`
}

func getComparisonCategories() []ComparisonCategory {
	return []ComparisonCategory{
		{
			Key:  "basic",
			Name: "基本信息",
			Items: []ComparisonItem{
				{Key: "name", Name: "名称", Type: "string"},
				{Key: "family", Name: "家族", Type: "string"},
				{Key: "releaseDate", Name: "发布日期", Type: "string"},
			},
		},
		{
			Key:  "capabilities",
			Name: "能力特性",
			Items: []ComparisonItem{
				{Key: "reasoning", Name: "推理", Type: "boolean"},
				{Key: "toolCall", Name: "工具调用", Type: "boolean"},
				{Key: "attachment", Name: "附件", Type: "boolean"},
				{Key: "openWeights", Name: "开源", Type: "boolean"},
			},
		},
		{
			Key:  "limits",
			Name: "限制",
			Items: []ComparisonItem{
				{Key: "limitContext", Name: "上下文", Type: "number", Unit: "tokens"},
				{Key: "limitInput", Name: "最大输入", Type: "number", Unit: "tokens"},
				{Key: "limitOutput", Name: "最大输出", Type: "number", Unit: "tokens"},
			},
		},
		{
			Key:  "pricing",
			Name: "定价",
			Items: []ComparisonItem{
				{Key: "costInput", Name: "输入", Type: "number", Unit: "$"},
				{Key: "costOutput", Name: "输出", Type: "number", Unit: "$"},
			},
		},
	}
}
