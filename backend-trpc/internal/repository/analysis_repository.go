package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ArtificialAnalysis 评测数据定义
type ArtificialAnalysis struct {
	ID                                  string  `json:"id"`
	Slug                                string  `json:"slug"`
	ModelCreator                        string  `json:"model_creator"`
	ArtificialAnalysisIntelligenceIndex float64 `json:"artificial_analysis_intelligence_index"`
	ArtificialAnalysisCodingIndex       float64 `json:"artificial_analysis_coding_index"`
	ArtificialAnalysisMathIndex         float64 `json:"artificial_analysis_math_index"`
	MmluPro                             float64 `json:"mmlu_pro"`
	Gpqa                                float64 `json:"gpqa"`
	Hle                                 float64 `json:"hle"`
	Livecodebench                       float64 `json:"livecodebench"`
	Scicode                             float64 `json:"scicode"`
	Math500                             float64 `json:"math_500"`
	Aime                                float64 `json:"aime"`
	Aime25                              float64 `json:"aime_25"`
	Ifbench                             float64 `json:"ifbench"`
	Lcr                                 float64 `json:"lcr"`
	TerminalbenchHard                   float64 `json:"terminalbench_hard"`
	Tau2                                float64 `json:"tau2"`
	Price1mBlended31                    float64 `json:"price_1m_blended_3_to_1"`
	Price1mInputTokens                  float64 `json:"price_1m_input_tokens"`
	Price1mOutputTokens                 float64 `json:"price_1m_output_tokens"`
	MedianOutputTokensPerSecond         float64 `json:"median_output_tokens_per_second"`
	MedianTimeToFirstTokenSeconds       float64 `json:"median_time_to_first_token_seconds"`
	MedianTimeToFirstAnswerToken        float64 `json:"median_time_to_first_answer_token"`
}

// AnalysisRepository 接口
type AnalysisRepository interface {
	List(ctx context.Context, page, limit int32) ([]*ArtificialAnalysis, int, error)
	GetBySlug(ctx context.Context, slug string) (*ArtificialAnalysis, error)
}

// analysisRepository 实现
type analysisRepository struct {
	dataPath string
	analysis []*ArtificialAnalysis
}

// NewAnalysisRepository 创建仓库实例
func NewAnalysisRepository() AnalysisRepository {
	repo := &analysisRepository{
		dataPath: filepath.Join("..", "backend", "data", "artificial_analysis.json"),
	}
	repo.loadData()
	return repo
}

// loadData 加载数据
func (r *analysisRepository) loadData() {
	data, err := os.ReadFile(r.dataPath)
	if err != nil {
		r.analysis = getDefaultAnalysis()
		return
	}

	var analysis []*ArtificialAnalysis
	if err := json.Unmarshal(data, &analysis); err != nil {
		r.analysis = getDefaultAnalysis()
		return
	}
	r.analysis = analysis
}

// List 获取列表
func (r *analysisRepository) List(ctx context.Context, page, limit int32) ([]*ArtificialAnalysis, int, error) {
	start := (page - 1) * limit
	if start > int32(len(r.analysis)) {
		return []*ArtificialAnalysis{}, len(r.analysis), nil
	}

	end := start + limit
	if end > int32(len(r.analysis)) {
		end = int32(len(r.analysis))
	}

	return r.analysis[start:end], len(r.analysis), nil
}

// GetBySlug 根据 slug 获取
func (r *analysisRepository) GetBySlug(ctx context.Context, slug string) (*ArtificialAnalysis, error) {
	for _, a := range r.analysis {
		if a.Slug == slug {
			return a, nil
		}
	}
	return nil, fmt.Errorf("analysis not found for slug: %s", slug)
}

// getDefaultAnalysis 获取默认数据
func getDefaultAnalysis() []*ArtificialAnalysis {
	return []*ArtificialAnalysis{
		{
			ID:                                  "1",
			Slug:                                "gpt-4o",
			ModelCreator:                        "OpenAI",
			ArtificialAnalysisIntelligenceIndex: 85.5,
			ArtificialAnalysisCodingIndex:       92.0,
			ArtificialAnalysisMathIndex:         88.5,
			MmluPro:                             0.875,
			Price1mInputTokens:                  2.5,
			Price1mOutputTokens:                 10.0,
		},
	}
}
