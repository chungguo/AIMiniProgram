package repository

import (
	"ai-model-papers-backend/models"
	"database/sql"
	// _ "github.com/lib/pq" // PostgreSQL驱动，需要时取消注释
)

// PostgresModelRepository PostgreSQL实现的模型仓库
// 这是预留的PostgreSQL实现，实际使用时需要：
// 1. 取消注释 pq 导入
// 2. 实现所有方法
// 3. 在 main.go 中替换 JSONRepository
type PostgresModelRepository struct {
	db *sql.DB
}

// NewPostgresModelRepository 创建PostgreSQL模型仓库
// 使用示例：
//   db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=secret dbname=ai_models sslmode=disable")
//   repo, err := NewPostgresModelRepository(db)
func NewPostgresModelRepository(db *sql.DB) *PostgresModelRepository {
	return &PostgresModelRepository{db: db}
}

// GetAll 获取所有模型（预留实现）
func (r *PostgresModelRepository) GetAll() ([]models.Model, error) {
	// TODO: 实现SQL查询
	// query := `
	//     SELECT id, name, provider, provider_id, description, 
	//            capabilities, context_window, max_tokens, 
	//            pricing, speed, quality, release_date, 
	//            model_type, architecture, features
	//     FROM models
	// `
	return nil, nil
}

// GetByID 根据ID获取模型（预留实现）
func (r *PostgresModelRepository) GetByID(id string) (*models.Model, error) {
	// TODO: 实现SQL查询
	// query := `SELECT * FROM models WHERE id = $1`
	return nil, nil
}

// GetByProvider 根据提供商筛选（预留实现）
func (r *PostgresModelRepository) GetByProvider(providerID string) ([]models.Model, error) {
	// TODO: 实现SQL查询
	// query := `SELECT * FROM models WHERE provider_id = $1`
	return nil, nil
}

// GetByCapability 根据能力筛选（预留实现）
func (r *PostgresModelRepository) GetByCapability(capability string) ([]models.Model, error) {
	// TODO: 实现SQL查询
	// query := `SELECT * FROM models WHERE capabilities @> $1`
	return nil, nil
}

// Search 搜索模型（预留实现）
func (r *PostgresModelRepository) Search(keyword string) ([]models.Model, error) {
	// TODO: 实现SQL查询
	// query := `
	//     SELECT * FROM models 
	//     WHERE name ILIKE $1 OR provider ILIKE $1 OR description ILIKE $1
	// `
	return nil, nil
}

// GetProviders 获取所有提供商（预留实现）
func (r *PostgresModelRepository) GetProviders() ([]models.Provider, error) {
	// TODO: 实现SQL查询
	// query := `SELECT * FROM providers`
	return nil, nil
}

// GetComparisonCategories 获取对比类别（预留实现）
func (r *PostgresModelRepository) GetComparisonCategories() ([]models.ComparisonCategory, error) {
	// TODO: 可以从配置文件或数据库获取
	// query := `SELECT * FROM comparison_categories`
	return nil, nil
}
