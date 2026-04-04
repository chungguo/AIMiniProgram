package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Paper 论文定义
type Paper struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	TitleCn    string `json:"title_cn"`
	Author     string `json:"author"`
	Abstract   string `json:"abstract"`
	AbstractCn string `json:"abstract_cn"`
	SubmitAt   string `json:"submit_at"`
}

// PaperFilter 筛选条件
type PaperFilter struct {
	Search   string
	Author   string
	DateFrom string
	DateTo   string
}

// PaperRepository 接口
type PaperRepository interface {
	List(ctx context.Context, filter *PaperFilter, page, limit int32) ([]*Paper, int, error)
	GetByID(ctx context.Context, id string) (*Paper, error)
	GetLatest(ctx context.Context, limit int32) ([]*Paper, error)
}

// paperRepository 实现
type paperRepository struct {
	dataPath string
	papers   []*Paper
}

// NewPaperRepository 创建仓库实例
func NewPaperRepository() PaperRepository {
	repo := &paperRepository{
		dataPath: filepath.Join("..", "backend", "data", "papers.json"),
	}
	repo.loadData()
	return repo
}

// loadData 加载数据
func (r *paperRepository) loadData() {
	data, err := os.ReadFile(r.dataPath)
	if err != nil {
		r.papers = getDefaultPapers()
		return
	}

	var papers []*Paper
	if err := json.Unmarshal(data, &papers); err != nil {
		r.papers = getDefaultPapers()
		return
	}
	r.papers = papers
}

// List 获取列表
func (r *paperRepository) List(ctx context.Context, filter *PaperFilter, page, limit int32) ([]*Paper, int, error) {
	filtered := r.filterPapers(filter)

	start := (page - 1) * limit
	if start > int32(len(filtered)) {
		return []*Paper{}, len(filtered), nil
	}

	end := start + limit
	if end > int32(len(filtered)) {
		end = int32(len(filtered))
	}

	return filtered[start:end], len(filtered), nil
}

// GetByID 根据 ID 获取
func (r *paperRepository) GetByID(ctx context.Context, id string) (*Paper, error) {
	for _, p := range r.papers {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, fmt.Errorf("paper not found: %s", id)
}

// GetLatest 获取最新论文
func (r *paperRepository) GetLatest(ctx context.Context, limit int32) ([]*Paper, error) {
	if limit > int32(len(r.papers)) {
		return r.papers, nil
	}
	return r.papers[:limit], nil
}

// filterPapers 筛选论文
func (r *paperRepository) filterPapers(filter *PaperFilter) []*Paper {
	if filter == nil {
		return r.papers
	}

	var result []*Paper
	for _, p := range r.papers {
		if filter.Search != "" {
			keyword := strings.ToLower(filter.Search)
			if !strings.Contains(strings.ToLower(p.Title), keyword) &&
				!strings.Contains(strings.ToLower(p.TitleCn), keyword) &&
				!strings.Contains(strings.ToLower(p.Author), keyword) {
				continue
			}
		}

		if filter.Author != "" && !strings.Contains(p.Author, filter.Author) {
			continue
		}

		result = append(result, p)
	}

	return result
}

// getDefaultPapers 获取默认数据
func getDefaultPapers() []*Paper {
	return []*Paper{
		{
			ID:         "2401.00001",
			Title:      "Attention Is All You Need",
			TitleCn:    "注意力机制就是你所需要的",
			Author:     "Ashish Vaswani, Noam Shazeer, Niki Parmar",
			Abstract:   "We propose a new simple network architecture, the Transformer.",
			AbstractCn: "我们提出了一种新的简单网络架构，Transformer。",
			SubmitAt:   "2024-01-01",
		},
	}
}
