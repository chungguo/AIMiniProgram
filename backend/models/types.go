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
	Family             string     `json:"family" db:"family"`
	Attachment         bool       `json:"attachment" db:"attachment"`
	Reasoning          bool       `json:"reasoning" db:"reasoning"`
	ToolCall           bool       `json:"toolCall" db:"tool_call"`
	StructuredOutput   bool       `json:"structuredOutput" db:"structured_output"`
	Temperature        bool       `json:"temperature" db:"temperature"`
	Knowledge          string     `json:"knowledge" db:"knowledge"`
	ReleaseDate        string     `json:"releaseDate" db:"release_date"`
	LastUpdated        string     `json:"lastUpdated" db:"last_updated"`
	ModalitiesInput    []Modality `json:"modalitiesInput" db:"modalities_input"`
	ModalitiesOutput   []Modality `json:"modalitiesOutput" db:"modalities_output"`
	OpenWeights        bool       `json:"openWeights" db:"open_weights"`
	CostInput          float64    `json:"costInput" db:"cost_input"`
	CostOutput         float64    `json:"costOutput" db:"cost_output"`
	CostReasoning      float64    `json:"costReasoning" db:"cost_reasoning"`
	CostCacheRead      float64    `json:"costCacheRead" db:"cost_cache_read"`
	CostCacheWrite     float64    `json:"costCacheWrite" db:"cost_cache_write"`
	CostInputAudio     float64    `json:"costInputAudio" db:"cost_input_audio"`
	CostOutputAudio    float64    `json:"costOutputAudio" db:"cost_output_audio"`
	LimitContext       int        `json:"limitContext" db:"limit_context"`
	LimitOutput        int        `json:"limitOutput" db:"limit_output"`
	LimitInput         int        `json:"limitInput" db:"limit_input"`
	InterleavedField   string     `json:"interleavedField" db:"interleaved_field"`
	CreatedAt          time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time  `json:"updatedAt" db:"updated_at"`
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
	Field string `json:"field"`
	Order string `json:"order"`
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
