package repository

import (
	"ai-model-papers-backend/models"
	"database/sql"
)

// PostgresArtificialAnalysisRepository PostgreSQL 实现的 ArtificialAnalysis 仓库
type PostgresArtificialAnalysisRepository struct {
	db *sql.DB
}

// NewPostgresArtificialAnalysisRepository 创建 PostgreSQL ArtificialAnalysis 仓库
func NewPostgresArtificialAnalysisRepository(db *sql.DB) *PostgresArtificialAnalysisRepository {
	return &PostgresArtificialAnalysisRepository{db: db}
}

// GetBySlug 根据 slug 获取评测数据
func (r *PostgresArtificialAnalysisRepository) GetBySlug(slug string) (*models.ArtificialAnalysis, error) {
	query := `
		SELECT id, slug, model_creator, 
		       artificial_analysis_intelligence_index, 
		       artificial_analysis_coding_index, 
		       artificial_analysis_math_index,
		       mmlu_pro, gpqa, hle, livecodebench, scicode, math_500, aime, aime_25,
		       ifbench, lcr, terminalbench_hard, tau2,
		       price_1m_blended_3_to_1, price_1m_input_tokens, price_1m_output_tokens,
		       median_output_tokens_per_second, median_time_to_first_token_seconds, 
		       median_time_to_first_answer_token, created_at, updated_at
		FROM artificialanalysis
		WHERE slug = $1
	`

	var a models.ArtificialAnalysis
	err := r.db.QueryRow(query, slug).Scan(
		&a.ID, &a.Slug, &a.ModelCreator,
		&a.ArtificialAnalysisIntelligenceIndex, &a.ArtificialAnalysisCodingIndex, &a.ArtificialAnalysisMathIndex,
		&a.MmluPro, &a.Gpqa, &a.Hle, &a.Livecodebench, &a.Scicode, &a.Math500, &a.Aime, &a.Aime25,
		&a.Ifbench, &a.Lcr, &a.TerminalbenchHard, &a.Tau2,
		&a.Price1mBlended31, &a.Price1mInputTokens, &a.Price1mOutputTokens,
		&a.MedianOutputTokensPerSecond, &a.MedianTimeToFirstTokenSeconds,
		&a.MedianTimeToFirstAnswerToken, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// GetBySlugs 批量根据 slugs 获取评测数据
func (r *PostgresArtificialAnalysisRepository) GetBySlugs(slugs []string) (map[string]*models.ArtificialAnalysis, error) {
	if len(slugs) == 0 {
		return make(map[string]*models.ArtificialAnalysis), nil
	}

	query := `
		SELECT id, slug, model_creator, 
		       artificial_analysis_intelligence_index, 
		       artificial_analysis_coding_index, 
		       artificial_analysis_math_index,
		       mmlu_pro, gpqa, hle, livecodebench, scicode, math_500, aime, aime_25,
		       ifbench, lcr, terminalbench_hard, tau2,
		       price_1m_blended_3_to_1, price_1m_input_tokens, price_1m_output_tokens,
		       median_output_tokens_per_second, median_time_to_first_token_seconds, 
		       median_time_to_first_answer_token, created_at, updated_at
		FROM artificialanalysis
		WHERE slug = ANY($1)
	`

	rows, err := r.db.Query(query, slugs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*models.ArtificialAnalysis)
	for rows.Next() {
		var a models.ArtificialAnalysis
		err := rows.Scan(
			&a.ID, &a.Slug, &a.ModelCreator,
			&a.ArtificialAnalysisIntelligenceIndex, &a.ArtificialAnalysisCodingIndex, &a.ArtificialAnalysisMathIndex,
			&a.MmluPro, &a.Gpqa, &a.Hle, &a.Livecodebench, &a.Scicode, &a.Math500, &a.Aime, &a.Aime25,
			&a.Ifbench, &a.Lcr, &a.TerminalbenchHard, &a.Tau2,
			&a.Price1mBlended31, &a.Price1mInputTokens, &a.Price1mOutputTokens,
			&a.MedianOutputTokensPerSecond, &a.MedianTimeToFirstTokenSeconds,
			&a.MedianTimeToFirstAnswerToken, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result[a.Slug] = &a
	}

	return result, rows.Err()
}

// GetAll 获取所有评测数据
func (r *PostgresArtificialAnalysisRepository) GetAll() ([]models.ArtificialAnalysis, error) {
	query := `
		SELECT id, slug, model_creator, 
		       artificial_analysis_intelligence_index, 
		       artificial_analysis_coding_index, 
		       artificial_analysis_math_index,
		       mmlu_pro, gpqa, hle, livecodebench, scicode, math_500, aime, aime_25,
		       ifbench, lcr, terminalbench_hard, tau2,
		       price_1m_blended_3_to_1, price_1m_input_tokens, price_1m_output_tokens,
		       median_output_tokens_per_second, median_time_to_first_token_seconds, 
		       median_time_to_first_answer_token, created_at, updated_at
		FROM artificialanalysis
		ORDER BY artificial_analysis_intelligence_index DESC NULLS LAST
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyses []models.ArtificialAnalysis
	for rows.Next() {
		var a models.ArtificialAnalysis
		err := rows.Scan(
			&a.ID, &a.Slug, &a.ModelCreator,
			&a.ArtificialAnalysisIntelligenceIndex, &a.ArtificialAnalysisCodingIndex, &a.ArtificialAnalysisMathIndex,
			&a.MmluPro, &a.Gpqa, &a.Hle, &a.Livecodebench, &a.Scicode, &a.Math500, &a.Aime, &a.Aime25,
			&a.Ifbench, &a.Lcr, &a.TerminalbenchHard, &a.Tau2,
			&a.Price1mBlended31, &a.Price1mInputTokens, &a.Price1mOutputTokens,
			&a.MedianOutputTokensPerSecond, &a.MedianTimeToFirstTokenSeconds,
			&a.MedianTimeToFirstAnswerToken, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		analyses = append(analyses, a)
	}

	return analyses, rows.Err()
}
