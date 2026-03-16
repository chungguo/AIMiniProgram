package repository

import (
	"ai-model-papers-backend/models"
	"encoding/json"
	"os"
	"strings"
	"sync"
)

// JSONPaperRepository JSON文件实现的论文仓库
type JSONPaperRepository struct {
	data     models.PapersData
	dataPath string
	mu       sync.RWMutex
}

// NewJSONPaperRepository 创建JSON论文仓库
func NewJSONPaperRepository(dataPath string) (*JSONPaperRepository, error) {
	repo := &JSONPaperRepository{
		dataPath: dataPath,
	}
	
	if err := repo.load(); err != nil {
		return nil, err
	}
	
	return repo, nil
}

// load 从JSON文件加载数据
func (r *JSONPaperRepository) load() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	data, err := os.ReadFile(r.dataPath)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, &r.data)
}

// GetAll 分页获取所有论文
func (r *JSONPaperRepository) GetAll(page, limit int) ([]models.Paper, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	total := len(r.data.Papers)
	
	start := (page - 1) * limit
	if start > total {
		start = total
	}
	
	end := start + limit
	if end > total {
		end = total
	}
	
	papers := make([]models.Paper, end-start)
	copy(papers, r.data.Papers[start:end])
	
	return papers, total, nil
}

// GetByID 根据ID获取论文
func (r *JSONPaperRepository) GetByID(id string) (*models.Paper, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	for _, paper := range r.data.Papers {
		if paper.ID == id {
			return &paper, nil
		}
	}
	return nil, nil
}

// GetByCategory 根据分类筛选
func (r *JSONPaperRepository) GetByCategory(category string) ([]models.Paper, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var result []models.Paper
	for _, paper := range r.data.Papers {
		for _, cat := range paper.Categories {
			if cat == category {
				result = append(result, paper)
				break
			}
		}
	}
	return result, nil
}

// Search 搜索论文
func (r *JSONPaperRepository) Search(keyword string) ([]models.Paper, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	keyword = strings.ToLower(keyword)
	var result []models.Paper
	
	for _, paper := range r.data.Papers {
		if strings.Contains(strings.ToLower(paper.Title), keyword) ||
			strings.Contains(strings.ToLower(paper.TitleCN), keyword) ||
			strings.Contains(strings.ToLower(paper.Abstract), keyword) ||
			strings.Contains(strings.ToLower(paper.AbstractCN), keyword) {
			result = append(result, paper)
			continue
		}
		
		for _, author := range paper.Authors {
			if strings.Contains(strings.ToLower(author), keyword) {
				result = append(result, paper)
				break
			}
		}
		
		for _, keyword := range paper.Keywords {
			if strings.Contains(strings.ToLower(keyword), keyword) {
				result = append(result, paper)
				break
			}
		}
	}
	return result, nil
}

// GetCategories 获取所有分类
func (r *JSONPaperRepository) GetCategories() ([]models.PaperCategory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	categories := make([]models.PaperCategory, len(r.data.Categories))
	copy(categories, r.data.Categories)
	return categories, nil
}

// GetLatest 获取最新论文
func (r *JSONPaperRepository) GetLatest(limit int) ([]models.Paper, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	if limit > len(r.data.Papers) {
		limit = len(r.data.Papers)
	}
	
	papers := make([]models.Paper, limit)
	copy(papers, r.data.Papers[:limit])
	return papers, nil
}
