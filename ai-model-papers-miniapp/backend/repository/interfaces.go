package repository

import "ai-model-papers-backend/models"

// ModelRepository 模型数据访问接口
type ModelRepository interface {
	GetAll() ([]models.Model, error)
	GetByID(id string) (*models.Model, error)
	GetByProvider(providerID string) ([]models.Model, error)
	GetByCapability(capability string) ([]models.Model, error)
	Search(keyword string) ([]models.Model, error)
	GetProviders() ([]models.Provider, error)
	GetComparisonCategories() ([]models.ComparisonCategory, error)
}

// PaperRepository 论文数据访问接口
type PaperRepository interface {
	GetAll(page, limit int) ([]models.Paper, int, error)
	GetByID(id string) (*models.Paper, error)
	GetByCategory(category string) ([]models.Paper, error)
	Search(keyword string) ([]models.Paper, error)
	GetCategories() ([]models.PaperCategory, error)
	GetLatest(limit int) ([]models.Paper, error)
}
