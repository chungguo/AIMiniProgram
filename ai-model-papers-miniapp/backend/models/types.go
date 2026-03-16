package models

// Capability 表示模型支持的能力
type Capability struct {
	Text  bool `json:"text"`
	Image bool `json:"image"`
	Audio bool `json:"audio"`
	Video bool `json:"video"`
	File  bool `json:"file"`
}

// Pricing 表示定价信息
type Pricing struct {
	InputPrice  float64 `json:"inputPrice"`
	OutputPrice float64 `json:"outputPrice"`
	Currency    string  `json:"currency"`
	Unit        string  `json:"unit"`
}

// Speed 表示速度指标
type Speed struct {
	Latency         string `json:"latency"`
	TokensPerSecond int    `json:"tokensPerSecond"`
}

// Quality 表示质量基准
type Quality struct {
	Mmlu      float64 `json:"mmlu,omitempty"`
	HumanEval float64 `json:"humanEval,omitempty"`
	MtBench   float64 `json:"mtBench,omitempty"`
	Math      float64 `json:"math,omitempty"`
}

// Model 表示大模型
type Model struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Provider       string     `json:"provider"`
	ProviderID     string     `json:"providerId"`
	Logo           string     `json:"logo"`
	Description    string     `json:"description"`
	Capabilities   Capability `json:"capabilities"`
	ContextWindow  int        `json:"contextWindow"`
	MaxTokens      int        `json:"maxTokens"`
	Pricing        Pricing    `json:"pricing"`
	Speed          Speed      `json:"speed"`
	Quality        Quality    `json:"quality"`
	ReleaseDate    string     `json:"releaseDate"`
	ModelType      string     `json:"modelType"`
	Architecture   string     `json:"architecture"`
	Features       []string   `json:"features"`
}

// Provider 表示模型提供商
type Provider struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Logo    string `json:"logo"`
	Website string `json:"website"`
	APIURL  string `json:"apiUrl"`
}

// ComparisonItem 表示对比项
type ComparisonItem struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Unit  string `json:"unit,omitempty"`
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

// ModelsData 存储所有模型数据
type ModelsData struct {
	Models               []Model              `json:"models"`
	Providers            []Provider           `json:"providers"`
	ComparisonCategories []ComparisonCategory `json:"comparisonCategories"`
}

// PapersData 存储所有论文数据
type PapersData struct {
	Papers     []Paper         `json:"papers"`
	Categories []PaperCategory `json:"categories"`
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