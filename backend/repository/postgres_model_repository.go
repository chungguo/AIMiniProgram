package repository

import (
	"ai-model-papers-backend/models"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

// PostgresModelRepository PostgreSQL 实现的模型仓库
type PostgresModelRepository struct {
	db *sql.DB
}

// NewPostgresModelRepository 创建 PostgreSQL 模型仓库
func NewPostgresModelRepository(db *sql.DB) *PostgresModelRepository {
	return &PostgresModelRepository{db: db}
}

// scanModel 扫描单行数据到 Model 结构体
func scanModel(rows *sql.Rows) (*models.Model, error) {
	var m models.Model
	var modalitiesInput, modalitiesOutput []string

	err := rows.Scan(
		&m.ID,
		&m.Name,
		&m.Family,
		&m.Attachment,
		&m.Reasoning,
		&m.ToolCall,
		&m.StructuredOutput,
		&m.Temperature,
		&m.Knowledge,
		&m.ReleaseDate,
		&m.LastUpdated,
		pq.Array(&modalitiesInput),
		pq.Array(&modalitiesOutput),
		&m.OpenWeights,
		&m.CostInput,
		&m.CostOutput,
		&m.LimitContext,
		&m.LimitOutput,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.InterleavedField,
		&m.LimitInput,
		&m.CostReasoning,
		&m.CostCacheRead,
		&m.CostCacheWrite,
		&m.CostInputAudio,
		&m.CostOutputAudio,
	)
	if err != nil {
		return nil, err
	}

	// 转换字符串数组为 Modality 类型
	m.ModalitiesInput = make([]models.Modality, len(modalitiesInput))
	for i, mod := range modalitiesInput {
		m.ModalitiesInput[i] = models.Modality(mod)
	}
	m.ModalitiesOutput = make([]models.Modality, len(modalitiesOutput))
	for i, mod := range modalitiesOutput {
		m.ModalitiesOutput[i] = models.Modality(mod)
	}

	return &m, nil
}

// GetAll 获取所有模型（支持筛选、排序、分页）
func (r *PostgresModelRepository) GetAll(
	filter *models.ModelFilter,
	sort *models.ModelSort,
	page, limit int,
) ([]models.Model, int, error) {
	// 构建 WHERE 条件
	whereClauses := []string{"1=1"}
	args := []interface{}{}
	argIdx := 1

	if filter != nil {
		if filter.Family != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("family = $%d", argIdx))
			args = append(args, filter.Family)
			argIdx++
		}
		if filter.HasAttachment != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("attachment = $%d", argIdx))
			args = append(args, *filter.HasAttachment)
			argIdx++
		}
		if filter.HasReasoning != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("reasoning = $%d", argIdx))
			args = append(args, *filter.HasReasoning)
			argIdx++
		}
		if filter.HasToolCall != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("tool_call = $%d", argIdx))
			args = append(args, *filter.HasToolCall)
			argIdx++
		}
		if filter.OpenWeights != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("open_weights = $%d", argIdx))
			args = append(args, *filter.OpenWeights)
			argIdx++
		}
		if filter.MinContext > 0 {
			whereClauses = append(whereClauses, fmt.Sprintf("limit_context >= $%d", argIdx))
			args = append(args, filter.MinContext)
			argIdx++
		}
		if filter.MaxCostInput > 0 {
			whereClauses = append(whereClauses, fmt.Sprintf("cost_input <= $%d", argIdx))
			args = append(args, filter.MaxCostInput)
			argIdx++
		}
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	// 获取总数
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM model WHERE %s", whereSQL)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 构建排序
	orderBy := "name ASC"
	if sort != nil && sort.Field != "" {
		validFields := map[string]string{
			"cost_input":    "cost_input",
			"cost_output":   "cost_output",
			"limit_context": "limit_context",
			"release_date":  "release_date",
			"name":          "name",
			"family":        "family",
		}
		if field, ok := validFields[sort.Field]; ok {
			order := "ASC"
			if strings.ToUpper(sort.Order) == "DESC" {
				order = "DESC"
			}
			orderBy = fmt.Sprintf("%s %s", field, order)
		}
	}

	// 构建查询
	columns := `
		id, name, family, attachment, reasoning, tool_call, structured_output, temperature,
		knowledge, release_date, last_updated, modalities_input, modalities_output,
		open_weights, cost_input, cost_output, limit_context, limit_output,
		created_at, updated_at, interleaved_field, limit_input,
		cost_reasoning, cost_cache_read, cost_cache_write, cost_input_audio, cost_output_audio
	`
	query := fmt.Sprintf(
		"SELECT %s FROM model WHERE %s ORDER BY %s LIMIT $%d OFFSET $%d",
		columns, whereSQL, orderBy, argIdx, argIdx+1,
	)

	offset := (page - 1) * limit
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var modelsList []models.Model
	for rows.Next() {
		m, err := scanModel(rows)
		if err != nil {
			return nil, 0, err
		}
		modelsList = append(modelsList, *m)
	}

	return modelsList, total, rows.Err()
}

// GetByID 根据 ID 获取模型
func (r *PostgresModelRepository) GetByID(id string) (*models.Model, error) {
	columns := `
		id, name, family, attachment, reasoning, tool_call, structured_output, temperature,
		knowledge, release_date, last_updated, modalities_input, modalities_output,
		open_weights, cost_input, cost_output, limit_context, limit_output,
		created_at, updated_at, interleaved_field, limit_input,
		cost_reasoning, cost_cache_read, cost_cache_write, cost_input_audio, cost_output_audio
	`
	query := fmt.Sprintf("SELECT %s FROM model WHERE id = $1", columns)

	row := r.db.QueryRow(query, id)

	var m models.Model
	var modalitiesInput, modalitiesOutput []string

	err := row.Scan(
		&m.ID,
		&m.Name,
		&m.Family,
		&m.Attachment,
		&m.Reasoning,
		&m.ToolCall,
		&m.StructuredOutput,
		&m.Temperature,
		&m.Knowledge,
		&m.ReleaseDate,
		&m.LastUpdated,
		pq.Array(&modalitiesInput),
		pq.Array(&modalitiesOutput),
		&m.OpenWeights,
		&m.CostInput,
		&m.CostOutput,
		&m.LimitContext,
		&m.LimitOutput,
		&m.CreatedAt,
		&m.UpdatedAt,
		&m.InterleavedField,
		&m.LimitInput,
		&m.CostReasoning,
		&m.CostCacheRead,
		&m.CostCacheWrite,
		&m.CostInputAudio,
		&m.CostOutputAudio,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// 转换字符串数组为 Modality 类型
	m.ModalitiesInput = make([]models.Modality, len(modalitiesInput))
	for i, mod := range modalitiesInput {
		m.ModalitiesInput[i] = models.Modality(mod)
	}
	m.ModalitiesOutput = make([]models.Modality, len(modalitiesOutput))
	for i, mod := range modalitiesOutput {
		m.ModalitiesOutput[i] = models.Modality(mod)
	}

	return &m, nil
}

// GetByFamily 根据家族获取模型
func (r *PostgresModelRepository) GetByFamily(family string) ([]models.Model, error) {
	columns := `
		id, name, family, attachment, reasoning, tool_call, structured_output, temperature,
		knowledge, release_date, last_updated, modalities_input, modalities_output,
		open_weights, cost_input, cost_output, limit_context, limit_output,
		created_at, updated_at, interleaved_field, limit_input,
		cost_reasoning, cost_cache_read, cost_cache_write, cost_input_audio, cost_output_audio
	`
	query := fmt.Sprintf("SELECT %s FROM model WHERE family = $1 ORDER BY name", columns)

	rows, err := r.db.Query(query, family)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modelsList []models.Model
	for rows.Next() {
		m, err := scanModel(rows)
		if err != nil {
			return nil, err
		}
		modelsList = append(modelsList, *m)
	}

	return modelsList, rows.Err()
}

// Search 搜索模型
func (r *PostgresModelRepository) Search(keyword string) ([]models.Model, error) {
	columns := `
		id, name, family, attachment, reasoning, tool_call, structured_output, temperature,
		knowledge, release_date, last_updated, modalities_input, modalities_output,
		open_weights, cost_input, cost_output, limit_context, limit_output,
		created_at, updated_at, interleaved_field, limit_input,
		cost_reasoning, cost_cache_read, cost_cache_write, cost_input_audio, cost_output_audio
	`
	query := fmt.Sprintf(
		"SELECT %s FROM model WHERE name ILIKE $1 OR family ILIKE $1 ORDER BY name",
		columns,
	)

	rows, err := r.db.Query(query, "%"+keyword+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var modelsList []models.Model
	for rows.Next() {
		m, err := scanModel(rows)
		if err != nil {
			return nil, err
		}
		modelsList = append(modelsList, *m)
	}

	return modelsList, rows.Err()
}

// GetFamilies 获取所有模型家族
func (r *PostgresModelRepository) GetFamilies() ([]string, error) {
	query := `SELECT DISTINCT family FROM model WHERE family IS NOT NULL ORDER BY family`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var families []string
	for rows.Next() {
		var family string
		if err := rows.Scan(&family); err != nil {
			return nil, err
		}
		families = append(families, family)
	}

	return families, rows.Err()
}

// GetComparisonCategories 获取对比类别（静态配置）
func (r *PostgresModelRepository) GetComparisonCategories() ([]models.ComparisonCategory, error) {
	return []models.ComparisonCategory{
		{
			Key:  "basic",
			Name: "基本信息",
			Items: []models.ComparisonItem{
				{Key: "family", Name: "家族", Type: "text"},
				{Key: "releaseDate", Name: "发布日期", Type: "date"},
				{Key: "knowledge", Name: "知识截止", Type: "text"},
				{Key: "openWeights", Name: "开源权重", Type: "boolean"},
			},
		},
		{
			Key:  "capabilities",
			Name: "能力特性",
			Items: []models.ComparisonItem{
				{Key: "reasoning", Name: "推理能力", Type: "boolean"},
				{Key: "toolCall", Name: "工具调用", Type: "boolean"},
				{Key: "attachment", Name: "附件支持", Type: "boolean"},
				{Key: "structuredOutput", Name: "结构化输出", Type: "boolean"},
				{Key: "temperature", Name: "温度调节", Type: "boolean"},
			},
		},
		{
			Key:  "modalities",
			Name: "模态支持",
			Items: []models.ComparisonItem{
				{Key: "modalitiesInput", Name: "输入模态", Type: "array"},
				{Key: "modalitiesOutput", Name: "输出模态", Type: "array"},
			},
		},
		{
			Key:  "limits",
			Name: "限制",
			Items: []models.ComparisonItem{
				{Key: "limitContext", Name: "上下文窗口", Type: "number", Unit: "tokens"},
				{Key: "limitInput", Name: "最大输入", Type: "number", Unit: "tokens"},
				{Key: "limitOutput", Name: "最大输出", Type: "number", Unit: "tokens"},
			},
		},
		{
			Key:  "pricing",
			Name: "定价",
			Items: []models.ComparisonItem{
				{Key: "costInput", Name: "输入价格", Type: "currency", Unit: "$/1M"},
				{Key: "costOutput", Name: "输出价格", Type: "currency", Unit: "$/1M"},
				{Key: "costReasoning", Name: "推理价格", Type: "currency", Unit: "$/1M"},
				{Key: "costCacheRead", Name: "缓存读取", Type: "currency", Unit: "$/1M"},
			},
		},
	}, nil
}
