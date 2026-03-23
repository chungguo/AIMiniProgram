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

// Paper 表示论文 - 严格匹配数据库表 arxiv_cs_ai 结构
type Paper struct {
	ID          string    `json:"id" db:"id"`                   // varchar(50) PRIMARY KEY
	Title       string    `json:"title" db:"title"`             // text
	Author      string    `json:"author" db:"author"`           // text (逗号分隔的作者字符串)
	Abstract    string    `json:"abstract" db:"abstract"`       // text
	TitleCn     string    `json:"title_cn" db:"title_cn"`       // text (可为null)
	AbstractCn  string    `json:"abstract_cn" db:"abstract_cn"` // text (可为null)
	SubmitAt    string    `json:"submit_at" db:"submit_at"`     // date (ISO格式)
	CreatedAt   time.Time `json:"created_at" db:"created_at"`   // timestamp
	UpdateAt    time.Time `json:"update_at" db:"update_at"`     // timestamp (注意: 数据库是update_at)
}

// ArtificialAnalysis 表示 ArtificialAnalysis.ai 的评测数据 - 匹配 PostgreSQL 表结构
type ArtificialAnalysis struct {
	ID                                  string    `json:"id" db:"id"`                                                 // uuid PRIMARY KEY
	Slug                                string    `json:"slug" db:"slug"`                                             // varchar(100) UNIQUE, 关联 model.name
	ModelCreator                        string    `json:"model_creator" db:"model_creator"`                           // varchar(100)
	ArtificialAnalysisIntelligenceIndex float64   `json:"artificial_analysis_intelligence_index" db:"artificial_analysis_intelligence_index"` // numeric(5,2)
	ArtificialAnalysisCodingIndex       float64   `json:"artificial_analysis_coding_index" db:"artificial_analysis_coding_index"`             // numeric(5,2)
	ArtificialAnalysisMathIndex         float64   `json:"artificial_analysis_math_index" db:"artificial_analysis_math_index"`                 // numeric(5,2)
	MmluPro                             float64   `json:"mmlu_pro" db:"mmlu_pro"`                                     // numeric(5,3)
	Gpqa                                float64   `json:"gpqa" db:"gpqa"`                                             // numeric(5,3)
	Hle                                 float64   `json:"hle" db:"hle"`                                               // numeric(5,3)
	Livecodebench                       float64   `json:"livecodebench" db:"livecodebench"`                           // numeric(5,3)
	Scicode                             float64   `json:"scicode" db:"scicode"`                                       // numeric(5,3)
	Math500                             float64   `json:"math_500" db:"math_500"`                                     // numeric(5,3)
	Aime                                float64   `json:"aime" db:"aime"`                                             // numeric(5,3)
	Aime25                              float64   `json:"aime_25" db:"aime_25"`                                       // numeric(5,3)
	Ifbench                             float64   `json:"ifbench" db:"ifbench"`                                       // numeric(5,3)
	Lcr                                 float64   `json:"lcr" db:"lcr"`                                               // numeric(5,3)
	TerminalbenchHard                   float64   `json:"terminalbench_hard" db:"terminalbench_hard"`                 // numeric(5,3)
	Tau2                                float64   `json:"tau2" db:"tau2"`                                             // numeric(5,3)
	Price1mBlended31                    float64   `json:"price_1m_blended_3_to_1" db:"price_1m_blended_3_to_1"`       // numeric(10,6)
	Price1mInputTokens                  float64   `json:"price_1m_input_tokens" db:"price_1m_input_tokens"`           // numeric(10,6)
	Price1mOutputTokens                 float64   `json:"price_1m_output_tokens" db:"price_1m_output_tokens"`         // numeric(10,6)
	MedianOutputTokensPerSecond         float64   `json:"median_output_tokens_per_second" db:"median_output_tokens_per_second"` // numeric(10,3)
	MedianTimeToFirstTokenSeconds       float64   `json:"median_time_to_first_token_seconds" db:"median_time_to_first_token_seconds"` // numeric(10,3)
	MedianTimeToFirstAnswerToken        float64   `json:"median_time_to_first_answer_token" db:"median_time_to_first_answer_token"` // numeric(10,3)
	CreatedAt                           time.Time `json:"created_at" db:"created_at"`                                 // timestamp
	UpdatedAt                           time.Time `json:"updated_at" db:"updated_at"`                                 // timestamp
}

// ModelWithAnalysis 组合 Model 和 ArtificialAnalysis 数据
type ModelWithAnalysis struct {
	Model              Model              `json:"model"`
	ArtificialAnalysis *ArtificialAnalysis `json:"analysis,omitempty"`
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
