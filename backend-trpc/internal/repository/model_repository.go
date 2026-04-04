package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Modality 模态类型
type Modality string

const (
	ModalityText  Modality = "text"
	ModalityImage Modality = "image"
	ModalityAudio Modality = "audio"
	ModalityVideo Modality = "video"
	ModalityFile  Modality = "file"
)

// Model 模型定义
type Model struct {
	ID                string     `json:"id"`
	Name              string     `json:"name"`
	Family            string     `json:"family"`
	Provider          string     `json:"provider"`
	Description       string     `json:"description"`
	Architecture      string     `json:"architecture"`
	Knowledge         string     `json:"knowledge"`
	ReleaseDate       string     `json:"releaseDate"`
	LastUpdated       string     `json:"lastUpdated"`
	Attachment        bool       `json:"attachment"`
	Reasoning         bool       `json:"reasoning"`
	ToolCall          bool       `json:"toolCall"`
	StructuredOutput  bool       `json:"structuredOutput"`
	Temperature       bool       `json:"temperature"`
	OpenWeights       bool       `json:"openWeights"`
	ModalitiesInput   []Modality `json:"modalitiesInput"`
	ModalitiesOutput  []string   `json:"modalitiesOutput"`
	CostInput         float64    `json:"costInput"`
	CostOutput        float64    `json:"costOutput"`
	CostReasoning     float64    `json:"costReasoning"`
	CostCacheRead     float64    `json:"costCacheRead"`
	CostCacheWrite    float64    `json:"costCacheWrite"`
	CostInputAudio    float64    `json:"costInputAudio"`
	CostOutputAudio   float64    `json:"costOutputAudio"`
	LimitContext      int        `json:"limitContext"`
	LimitInput        int        `json:"limitInput"`
	LimitOutput       int        `json:"limitOutput"`
	InterleavedField  string     `json:"interleavedField"`
}

// ModelFilter 筛选条件
type ModelFilter struct {
	Family          string
	HasAttachment   *bool
	HasReasoning    *bool
	HasToolCall     *bool
	OpenWeights     *bool
	ModalitiesInput []string
	MinContext      int32
	MaxCostInput    float64
	Search          string
}

// ModelRepository 接口
type ModelRepository interface {
	List(ctx context.Context, filter *ModelFilter, page, limit int32) ([]*Model, int, error)
	GetByID(ctx context.Context, id string) (*Model, error)
	GetByFamily(ctx context.Context, family string) ([]*Model, error)
	GetFamilies(ctx context.Context) ([]string, error)
	Search(ctx context.Context, keyword string) ([]*Model, error)
}

// modelRepository 实现
type modelRepository struct {
	dataPath string
	models   []*Model
}

// NewModelRepository 创建仓库实例
func NewModelRepository() ModelRepository {
	repo := &modelRepository{
		dataPath: filepath.Join("..", "backend", "data", "models.json"),
	}
	repo.loadData()
	return repo
}

// loadData 加载数据
func (r *modelRepository) loadData() {
	data, err := os.ReadFile(r.dataPath)
	if err != nil {
		// 使用默认数据
		r.models = getDefaultModels()
		return
	}

	var models []*Model
	if err := json.Unmarshal(data, &models); err != nil {
		r.models = getDefaultModels()
		return
	}
	r.models = models
}

// List 获取列表
func (r *modelRepository) List(ctx context.Context, filter *ModelFilter, page, limit int32) ([]*Model, int, error) {
	filtered := r.filterModels(filter)
	
	// 分页
	start := (page - 1) * limit
	if start > int32(len(filtered)) {
		return []*Model{}, len(filtered), nil
	}
	
	end := start + limit
	if end > int32(len(filtered)) {
		end = int32(len(filtered))
	}
	
	return filtered[start:end], len(filtered), nil
}

// GetByID 根据 ID 获取
func (r *modelRepository) GetByID(ctx context.Context, id string) (*Model, error) {
	for _, m := range r.models {
		if m.ID == id {
			return m, nil
		}
	}
	return nil, fmt.Errorf("model not found: %s", id)
}

// GetByFamily 根据家族获取
func (r *modelRepository) GetByFamily(ctx context.Context, family string) ([]*Model, error) {
	var result []*Model
	for _, m := range r.models {
		if m.Family == family {
			result = append(result, m)
		}
	}
	return result, nil
}

// GetFamilies 获取所有家族
func (r *modelRepository) GetFamilies(ctx context.Context) ([]string, error) {
	familyMap := make(map[string]bool)
	for _, m := range r.models {
		familyMap[m.Family] = true
	}
	
	families := make([]string, 0, len(familyMap))
	for f := range familyMap {
		families = append(families, f)
	}
	return families, nil
}

// Search 搜索
func (r *modelRepository) Search(ctx context.Context, keyword string) ([]*Model, error) {
	keyword = strings.ToLower(keyword)
	var result []*Model
	for _, m := range r.models {
		if strings.Contains(strings.ToLower(m.Name), keyword) ||
			strings.Contains(strings.ToLower(m.Family), keyword) ||
			strings.Contains(strings.ToLower(m.Description), keyword) {
			result = append(result, m)
		}
	}
	return result, nil
}

// filterModels 筛选模型
func (r *modelRepository) filterModels(filter *ModelFilter) []*Model {
	if filter == nil {
		return r.models
	}

	var result []*Model
	for _, m := range r.models {
		// 搜索关键词
		if filter.Search != "" {
			keyword := strings.ToLower(filter.Search)
			if !strings.Contains(strings.ToLower(m.Name), keyword) &&
				!strings.Contains(strings.ToLower(m.Family), keyword) {
				continue
			}
		}

		// 家族筛选
		if filter.Family != "" && m.Family != filter.Family {
			continue
		}

		// 布尔特性筛选
		if filter.HasAttachment != nil && m.Attachment != *filter.HasAttachment {
			continue
		}
		if filter.HasReasoning != nil && m.Reasoning != *filter.HasReasoning {
			continue
		}
		if filter.HasToolCall != nil && m.ToolCall != *filter.HasToolCall {
			continue
		}
		if filter.OpenWeights != nil && m.OpenWeights != *filter.OpenWeights {
			continue
		}

		// 上下文限制筛选
		if filter.MinContext > 0 && m.LimitContext < int(filter.MinContext) {
			continue
		}

		// 成本筛选
		if filter.MaxCostInput > 0 && m.CostInput > filter.MaxCostInput {
			continue
		}

		result = append(result, m)
	}

	return result
}

// getDefaultModels 获取默认数据
func getDefaultModels() []*Model {
	return []*Model{
		{
			ID:               "gpt-4o",
			Name:             "GPT-4o",
			Family:           "GPT-4",
			Provider:         "OpenAI",
			Description:      "OpenAI 的旗舰模型",
			Attachment:       true,
			Reasoning:        false,
			ToolCall:         true,
			StructuredOutput: true,
			Temperature:      true,
			ModalitiesInput:  []Modality{ModalityText, ModalityImage},
			ModalitiesOutput: []string{"text"},
			CostInput:        2.5,
			CostOutput:       10.0,
			LimitContext:     128000,
			LimitInput:       128000,
			LimitOutput:      4096,
		},
		{
			ID:               "claude-3-5-sonnet",
			Name:             "Claude 3.5 Sonnet",
			Family:           "Claude 3.5",
			Provider:         "Anthropic",
			Description:      "Anthropic 的智能模型",
			Attachment:       true,
			Reasoning:        true,
			ToolCall:         true,
			StructuredOutput: true,
			Temperature:      true,
			ModalitiesInput:  []Modality{ModalityText, ModalityImage},
			ModalitiesOutput: []string{"text"},
			CostInput:        3.0,
			CostOutput:       15.0,
			LimitContext:     200000,
			LimitInput:       200000,
			LimitOutput:      8192,
		},
	}
}
