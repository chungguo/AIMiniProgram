package repository

import (
	"fmt"
	"os"
)

// RepositoryFactory 仓库工厂
type RepositoryFactory struct {
	modelRepo ModelRepository
	paperRepo PaperRepository
}

// NewRepositoryFactory 创建仓库工厂
// 根据环境变量或配置决定使用哪种存储实现
func NewRepositoryFactory() (*RepositoryFactory, error) {
	// 检查是否使用PostgreSQL
	dbURL := os.Getenv("DATABASE_URL")
	
	if dbURL != "" {
		// 使用PostgreSQL
		// db, err := sql.Open("postgres", dbURL)
		// if err != nil {
		//     return nil, fmt.Errorf("failed to connect to database: %w", err)
		// }
		// 
		// factory := &RepositoryFactory{
		//     modelRepo: NewPostgresModelRepository(db),
		//     paperRepo: NewPostgresPaperRepository(db),
		// }
		// return factory, nil
		
		// 暂时返回错误，提醒需要实现
		return nil, fmt.Errorf("PostgreSQL not implemented yet, please set up database or unset DATABASE_URL")
	}
	
	// 使用JSON文件（默认）
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	
	modelRepo, err := NewJSONModelRepository(fmt.Sprintf("%s/data/models.json", wd))
	if err != nil {
		return nil, fmt.Errorf("failed to load models: %w", err)
	}
	
	paperRepo, err := NewJSONPaperRepository(fmt.Sprintf("%s/data/papers.json", wd))
	if err != nil {
		return nil, fmt.Errorf("failed to load papers: %w", err)
	}
	
	return &RepositoryFactory{
		modelRepo: modelRepo,
		paperRepo: paperRepo,
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
