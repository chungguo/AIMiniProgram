package handlers

import (
	"net/http"
	"strconv"
	"strings"
)

// GetPapers 获取所有论文
func GetPapers(w http.ResponseWriter, r *http.Request) {
	search := strings.ToLower(r.URL.Query().Get("search"))
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 10
	}

	var papers interface{}
	var total int
	var err error

	if search != "" {
		papers, err = paperRepo.Search(search)
		total = len(papers.([]interface{}))
	} else {
		papers, total, err = paperRepo.GetAll(page, limit)
	}

	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	totalPages := (total + limit - 1) / limit

	// 直接返回带分页的数据结构
	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"data": papers,
		"pagination": map[string]interface{}{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
		},
	})
}

// GetPaperByID 获取单篇论文
func GetPaperByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/papers/detail/"):]

	paper, err := paperRepo.GetByID(id)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if paper == nil {
		ErrorResponse(w, http.StatusNotFound, "Paper not found")
		return
	}

	JSONResponse(w, http.StatusOK, paper)
}

// GetLatestPapers 获取最新论文
func GetLatestPapers(w http.ResponseWriter, r *http.Request) {
	limit := 5
	papers, err := paperRepo.GetLatest(limit)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	JSONResponse(w, http.StatusOK, papers)
}
