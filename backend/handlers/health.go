package handlers

import (
	"net/http"
	"time"
)

// HealthCheck 健康检查
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}