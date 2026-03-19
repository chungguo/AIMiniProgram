package repository

import (
	"ai-model-papers-backend/models"
	"database/sql"

	"github.com/lib/pq"
)

// PostgresPaperRepository PostgreSQL 实现的论文仓库
type PostgresPaperRepository struct {
	db *sql.DB
}

// NewPostgresPaperRepository 创建 PostgreSQL 论文仓库
func NewPostgresPaperRepository(db *sql.DB) *PostgresPaperRepository {
	return &PostgresPaperRepository{db: db}
}

// GetAll 分页获取所有论文
func (r *PostgresPaperRepository) GetAll(page, limit int) ([]models.Paper, int, error) {
	// 获取总数
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM paper").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	query := `
		SELECT id, title, title_cn, abstract, abstract_cn, 
		       authors, institutions, publish_date, 
		       arxiv_url, pdf_url, categories, keywords, read_time, language
		FROM paper
		ORDER BY publish_date DESC
		LIMIT $1 OFFSET $2
	`

	offset := (page - 1) * limit
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	return r.scanPapers(rows, total)
}

// GetByID 根据ID获取论文
func (r *PostgresPaperRepository) GetByID(id string) (*models.Paper, error) {
	query := `
		SELECT id, title, title_cn, abstract, abstract_cn, 
		       authors, institutions, publish_date, 
		       arxiv_url, pdf_url, categories, keywords, read_time, language
		FROM paper
		WHERE id = $1
	`

	var p models.Paper
	var authors, institutions, categories, keywords []string

	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.Title, &p.TitleCN, &p.Abstract, &p.AbstractCN,
		pq.Array(&authors), pq.Array(&institutions), &p.PublishDate,
		&p.ArxivURL, &p.PDFURL, pq.Array(&categories), pq.Array(&keywords),
		&p.ReadTime, &p.Language,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	p.Authors = authors
	p.Institutions = institutions
	p.Categories = categories
	p.Keywords = keywords

	return &p, nil
}

// GetByCategory 根据分类筛选
func (r *PostgresPaperRepository) GetByCategory(category string) ([]models.Paper, error) {
	query := `
		SELECT id, title, title_cn, abstract, abstract_cn, 
		       authors, institutions, publish_date, 
		       arxiv_url, pdf_url, categories, keywords, read_time, language
		FROM paper
		WHERE $1 = ANY(categories)
		ORDER BY publish_date DESC
	`

	rows, err := r.db.Query(query, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	papers, _, err := r.scanPapers(rows, 0)
	return papers, err
}

// Search 搜索论文
func (r *PostgresPaperRepository) Search(keyword string) ([]models.Paper, error) {
	query := `
		SELECT id, title, title_cn, abstract, abstract_cn, 
		       authors, institutions, publish_date, 
		       arxiv_url, pdf_url, categories, keywords, read_time, language
		FROM paper
		WHERE title ILIKE $1 
		   OR title_cn ILIKE $1 
		   OR abstract ILIKE $1 
		   OR abstract_cn ILIKE $1
		   OR EXISTS (
		       SELECT 1 FROM unnest(authors) a 
		       WHERE a ILIKE $1
		   )
		ORDER BY publish_date DESC
	`

	rows, err := r.db.Query(query, "%"+keyword+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	papers, _, err := r.scanPapers(rows, 0)
	return papers, err
}

// GetCategories 获取所有分类
func (r *PostgresPaperRepository) GetCategories() ([]models.PaperCategory, error) {
	// 从 paper 表中提取唯一分类
	query := `
		SELECT DISTINCT UNNEST(categories) as category
		FROM paper
		ORDER BY category
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.PaperCategory
	for rows.Next() {
		var cat string
		if err := rows.Scan(&cat); err != nil {
			return nil, err
		}
		categories = append(categories, models.PaperCategory{
			ID:     cat,
			Name:   cat,
			NameEN: cat,
		})
	}

	return categories, rows.Err()
}

// GetLatest 获取最新论文
func (r *PostgresPaperRepository) GetLatest(limit int) ([]models.Paper, error) {
	query := `
		SELECT id, title, title_cn, abstract, abstract_cn, 
		       authors, institutions, publish_date, 
		       arxiv_url, pdf_url, categories, keywords, read_time, language
		FROM paper
		ORDER BY publish_date DESC
		LIMIT $1
	`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	papers, _, err := r.scanPapers(rows, 0)
	return papers, err
}

// scanPapers 扫描论文结果集
func (r *PostgresPaperRepository) scanPapers(rows *sql.Rows, total int) ([]models.Paper, int, error) {
	var papers []models.Paper

	for rows.Next() {
		var p models.Paper
		var authors, institutions, categories, keywords []string

		err := rows.Scan(
			&p.ID, &p.Title, &p.TitleCN, &p.Abstract, &p.AbstractCN,
			pq.Array(&authors), pq.Array(&institutions), &p.PublishDate,
			&p.ArxivURL, &p.PDFURL, pq.Array(&categories), pq.Array(&keywords),
			&p.ReadTime, &p.Language,
		)
		if err != nil {
			return nil, 0, err
		}

		p.Authors = authors
		p.Institutions = institutions
		p.Categories = categories
		p.Keywords = keywords

		papers = append(papers, p)
	}

	return papers, total, rows.Err()
}
