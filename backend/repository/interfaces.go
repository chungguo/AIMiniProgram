package repository

import "ai-model-papers-backend/models"

// ModelRepository 模型数据访问接口 - 更新以支持新表结构
type ModelRepository interface {
	GetAll(filter *models.ModelFilter, sort *models.ModelSort, page, limit int) ([]models.Model, int, error)
	GetByID(id string) (*models.Model, error)
	GetByFamily(family string) ([]models.Model, error)
	Search(keyword string) ([]models.Model, error)
	GetFamilies() ([]string, error)
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
