package models

import "time"

// Modality 表示模型支持的模态类型
type Modality string

const (
	ModalityText  Modality = "text"
	ModalityImage Modality = "image"
	ModalityAudio Modality = "audio"
	ModalityVideo Modality = "video"
	ModalityFile  Modality = "file"
)

// Model 表示大模型 - 匹配 PostgreSQL 表结构
type Model struct {
	ID                 string     `json:"id" db:"id"`
	Name               string     `json:"name" db:"name"`
	Family             string     `json:"family" db:"family"` // 替代原来的 Provider
	Attachment         bool       `json:"attachment" db:"attachment"`
	Reasoning          bool       `json:"reasoning" db:"reasoning"`
	ToolCall           bool       `json:"toolCall" db:"tool_call"`
	StructuredOutput   bool       `json:"structuredOutput" db:"structured_output"`
	Temperature        bool       `json:"temperature" db:"temperature"`
	Knowledge          string     `json:"knowledge" db:"knowledge"`           // 知识截止日期
	ReleaseDate        string     `json:"releaseDate" db:"release_date"`      // 发布日期
	LastUpdated        string     `json:"lastUpdated" db:"last_updated"`      // 最后更新
	ModalitiesInput    []Modality `json:"modalitiesInput" db:"modalities_input"`   // 输入模态
	ModalitiesOutput   []Modality `json:"modalitiesOutput" db:"modalities_output"` // 输出模态
	OpenWeights        bool       `json:"openWeights" db:"open_weights"`      // 权重是否开源
	// 定价（每百万 tokens，USD）
	CostInput          float64    `json:"costInput" db:"cost_input"`
	CostOutput         float64    `json:"costOutput" db:"cost_output"`
	CostReasoning      float64    `json:"costReasoning" db:"cost_reasoning"`
	CostCacheRead      float64    `json:"costCacheRead" db:"cost_cache_read"`
	CostCacheWrite     float64    `json:"costCacheWrite" db:"cost_cache_write"`
	CostInputAudio     float64    `json:"costInputAudio" db:"cost_input_audio"`
	CostOutputAudio    float64    `json:"costOutputAudio" db:"cost_output_audio"`
	// 限制
	LimitContext       int        `json:"limitContext" db:"limit_context"` // 最大上下文窗口
	LimitOutput        int        `json:"limitOutput" db:"limit_output"`   // 最大输出 tokens
	LimitInput         int        `json:"limitInput" db:"limit_input"`     // 最大输入 tokens
	// 推理内容字段名
	InterleavedField   string     `json:"interleavedField" db:"interleaved_field"` // "reasoning_content" 或 "reasoning_details"
	// 时间戳
	CreatedAt          time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time  `json:"updatedAt" db:"updated_at"`
	// 扩展字段（JSON 数据兼容）
	Description        string     `json:"description,omitempty" db:"-"`
	Features           []string   `json:"features,omitempty" db:"-"`
	Logo               string     `json:"logo,omitempty" db:"-"`
	Architecture       string     `json:"architecture,omitempty" db:"-"`
	// 旧版兼容字段
	Provider           string     `json:"provider,omitempty" db:"-"`     // 兼容旧 JSON
	ProviderID         string     `json:"providerId,omitempty" db:"-"`   // 兼容旧 JSON
	ContextWindow      int        `json:"contextWindow,omitempty" db:"-"` // 兼容旧 JSON
	MaxTokens          int        `json:"maxTokens,omitempty" db:"-"`     // 兼容旧 JSON
}

// ModelFilter 模型筛选条件
type ModelFilter struct {
	Family          string   `json:"family,omitempty"`
	HasAttachment   *bool    `json:"hasAttachment,omitempty"`
	HasReasoning    *bool    `json:"hasReasoning,omitempty"`
	HasToolCall     *bool    `json:"hasToolCall,omitempty"`
	OpenWeights     *bool    `json:"openWeights,omitempty"`
	ModalitiesInput []string `json:"modalitiesInput,omitempty"`
	MinContext      int      `json:"minContext,omitempty"`
	MaxCostInput    float64  `json:"maxCostInput,omitempty"`
}

// ModelSort 模型排序选项
type ModelSort struct {
	Field string `json:"field"` // cost_input, cost_output, limit_context, release_date
	Order string `json:"order"` // asc, desc
}

// Pricing 表示定价信息（前端展示用，向后兼容）
type Pricing struct {
	InputPrice  float64 `json:"inputPrice"`
	OutputPrice float64 `json:"outputPrice"`
	Currency    string  `json:"currency"`
	Unit        string  `json:"unit"`
}

// Capability 表示模型支持的能力（前端展示用，向后兼容）
type Capability struct {
	Text  bool `json:"text"`
	Image bool `json:"image"`
	Audio bool `json:"audio"`
	Video bool `json:"video"`
	File  bool `json:"file"`
}

// ToCapabilities 将 Modalities 转换为 Capability（兼容层）
func (m *Model) ToCapabilities() Capability {
	c := Capability{}
	for _, mod := range m.ModalitiesInput {
		switch mod {
		case ModalityText:
			c.Text = true
		case ModalityImage:
			c.Image = true
		case ModalityAudio:
			c.Audio = true
		case ModalityVideo:
			c.Video = true
		case ModalityFile:
			c.File = true
		}
	}
	return c
}

// ToPricing 转换为兼容的 Pricing 结构
func (m *Model) ToPricing() Pricing {
	return Pricing{
		InputPrice:  m.CostInput,
		OutputPrice: m.CostOutput,
		Currency:    "USD",
		Unit:        "per 1M tokens",
	}
}

// Provider 表示模型提供商/家族
type Provider struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Logo    string `json:"logo"`
	Website string `json:"website"`
	APIURL  string `json:"apiUrl"`
}

// ComparisonItem 表示对比项
type ComparisonItem struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Type string `json:"type"`
	Unit string `json:"unit,omitempty"`
}

// ComparisonCategory 表示对比类别
type ComparisonCategory struct {
	Key   string           `json:"key"`
	Name  string           `json:"name"`
	Items []ComparisonItem `json:"items"`
}

// Paper 表示论文
type Paper struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	TitleCN      string   `json:"titleCN"`
	Abstract     string   `json:"abstract"`
	AbstractCN   string   `json:"abstractCN"`
	Authors      []string `json:"authors"`
	Institutions []string `json:"institutions"`
	PublishDate  string   `json:"publishDate"`
	ArxivURL     string   `json:"arxivUrl"`
	PDFURL       string   `json:"pdfUrl"`
	Categories   []string `json:"categories"`
	Keywords     []string `json:"keywords"`
	ReadTime     int      `json:"readTime"`
	Language     string   `json:"language"`
}

// PaperCategory 表示论文分类
type PaperCategory struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	NameEN string `json:"nameEn"`
}

// APIResponse 统一API响应格式
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// Pagination 分页信息
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// PaginatedResponse 分页响应
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// ModelsData 存储所有模型数据（JSON兼容层）
type ModelsData struct {
	Models               []Model              `json:"models"`
	Providers            []Provider           `json:"providers"`
	ComparisonCategories []ComparisonCategory `json:"comparisonCategories"`
}

// PapersData 存储所有论文数据（JSON兼容层）
type PapersData struct {
	Papers     []Paper         `json:"papers"`
	Categories []PaperCategory `json:"categories"`
}
