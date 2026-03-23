package repository

import (
	"ai-model-papers-backend/models"
	"database/sql"
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
	err := r.db.QueryRow("SELECT COUNT(*) FROM arxiv_cs_ai").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 分页查询 - 严格匹配数据库表 arxiv_cs_ai
	query := `
		SELECT id, title, author, abstract, title_cn, abstract_cn, 
		       submit_at, created_at, update_at
		FROM arxiv_cs_ai
		ORDER BY submit_at DESC
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
		SELECT id, title, author, abstract, title_cn, abstract_cn, 
		       submit_at, created_at, update_at
		FROM arxiv_cs_ai
		WHERE id = $1
	`

	var p models.Paper
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.Title, &p.Author, &p.Abstract, &p.TitleCn, &p.AbstractCn,
		&p.SubmitAt, &p.CreatedAt, &p.UpdateAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// Search 搜索论文
func (r *PostgresPaperRepository) Search(keyword string) ([]models.Paper, error) {
	query := `
		SELECT id, title, author, abstract, title_cn, abstract_cn, 
		       submit_at, created_at, update_at
		FROM arxiv_cs_ai
		WHERE title ILIKE $1 
		   OR title_cn ILIKE $1 
		   OR abstract ILIKE $1 
		   OR abstract_cn ILIKE $1
		   OR author ILIKE $1
		ORDER BY submit_at DESC
	`

	rows, err := r.db.Query(query, "%"+keyword+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	papers, _, err := r.scanPapers(rows, 0)
	return papers, err
}

// GetLatest 获取最新论文
func (r *PostgresPaperRepository) GetLatest(limit int) ([]models.Paper, error) {
	query := `
		SELECT id, title, author, abstract, title_cn, abstract_cn, 
		       submit_at, created_at, update_at
		FROM arxiv_cs_ai
		ORDER BY submit_at DESC
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

		err := rows.Scan(
			&p.ID, &p.Title, &p.Author, &p.Abstract, &p.TitleCn, &p.AbstractCn,
			&p.SubmitAt, &p.CreatedAt, &p.UpdateAt,
		)
		if err != nil {
			return nil, 0, err
		}

		papers = append(papers, p)
	}

	return papers, total, rows.Err()
}
