package repository

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// RepositoryFactory 仓库工厂
type RepositoryFactory struct {
	modelRepo ModelRepository
	paperRepo PaperRepository
}

// NewRepositoryFactory 创建仓库工厂
// 必须使用 PostgreSQL，通过 DATABASE_URL 环境变量配置
func NewRepositoryFactory() (*RepositoryFactory, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &RepositoryFactory{
		modelRepo: NewPostgresModelRepository(db),
		paperRepo: NewPostgresPaperRepository(db),
	}, nil
}

// GetModelRepository 获取模型仓库
func (f *RepositoryFactory) GetModelRepository() ModelRepository {
	return f.modelRepo
}

// GetPaperRepository 获取论文仓库
func (f *RepositoryFactory) GetPaperRepository() PaperRepository {
	return f.paperRepo
}
