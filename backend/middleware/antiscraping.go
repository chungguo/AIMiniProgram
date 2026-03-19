package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AntiScrapingConfig 反爬虫配置
type AntiScrapingConfig struct {
	// User-Agent黑名单
	BlockedUserAgents []string
	// 请求路径黑名单
	BlockedPaths []string
	// 最小请求间隔（毫秒）
	MinRequestInterval int64
	// 是否检查Referer
	CheckReferer bool
	// 允许的Referer列表
	AllowedReferers []string
}

// DefaultAntiScrapingConfig 默认反爬虫配置
var DefaultAntiScrapingConfig = AntiScrapingConfig{
	BlockedUserAgents: []string{
		"scrapy", "curl", "wget", "python-requests", "httpx",
		"aiohttp", "httpclient", "java", "bot", "spider", "crawler",
	},
	BlockedPaths: []string{
		".env", ".git", "config", "admin", "wp-admin", "phpmyadmin",
	},
	MinRequestInterval: 50, // 50毫秒
	CheckReferer:       false,
}

// RequestRecord 请求记录
type RequestRecord struct {
	Count     int
	LastTime  int64
	BlockTime int64
}

// AntiScrapingMiddleware 反爬虫中间件
type AntiScrapingMiddleware struct {
	config   AntiScrapingConfig
	records  map[string]*RequestRecord
}

// NewAntiScrapingMiddleware 创建反爬虫中间件
func NewAntiScrapingMiddleware(config AntiScrapingConfig) *AntiScrapingMiddleware {
	return &AntiScrapingMiddleware{
		config:  config,
		records: make(map[string]*RequestRecord),
	}
}

// Middleware 返回Gin中间件
func (as *AntiScrapingMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 检查User-Agent
		userAgent := strings.ToLower(c.GetHeader("User-Agent"))
		
		// 开发环境放行（允许curl测试）
		if c.GetHeader("X-Dev-Mode") == "true" {
			c.Next()
			return
		}
		
		if as.isBlockedUserAgent(userAgent) {
			as.blockRequest(c, "BLOCKED_USER_AGENT")
			return
		}

		// 2. 检查请求路径
		path := strings.ToLower(c.Request.URL.Path)
		if as.isBlockedPath(path) {
			as.blockRequest(c, "BLOCKED_PATH")
			return
		}

		// 3. 检查请求频率
		ip := c.ClientIP()
		if ip == "" {
			ip = c.Request.RemoteAddr
		}
		if as.isRequestTooFrequent(ip) {
			as.blockRequest(c, "REQUEST_TOO_FREQUENT")
			return
		}

		c.Next()
	}
}

// isBlockedUserAgent 检查User-Agent是否在黑名单
func (as *AntiScrapingMiddleware) isBlockedUserAgent(ua string) bool {
	for _, blocked := range as.config.BlockedUserAgents {
		if strings.Contains(ua, blocked) {
			return true
		}
	}
	return false
}

// isBlockedPath 检查路径是否在黑名单
func (as *AntiScrapingMiddleware) isBlockedPath(path string) bool {
	for _, blocked := range as.config.BlockedPaths {
		if strings.Contains(path, blocked) {
			return true
		}
	}
	return false
}

// isRequestTooFrequent 检查请求是否过于频繁
func (as *AntiScrapingMiddleware) isRequestTooFrequent(ip string) bool {
	now := time.Now().UnixMilli()
	
	record, exists := as.records[ip]
	if !exists {
		as.records[ip] = &RequestRecord{
			Count:    1,
			LastTime: now,
		}
		return false
	}

	// 如果被临时封禁
	if record.BlockTime > now {
		return true
	}

	interval := now - record.LastTime
	if interval < as.config.MinRequestInterval {
		record.Count++
		// 如果连续频繁请求，临时封禁
		if record.Count > 10 {
			record.BlockTime = now + 60000 // 封禁1分钟
		}
		return true
	}

	record.LastTime = now
	if interval > 1000 { // 超过1秒，重置计数
		record.Count = 1
	}
	
	return false
}

// blockRequest 阻止请求
func (as *AntiScrapingMiddleware) blockRequest(c *gin.Context, reason string) {
	c.JSON(http.StatusForbidden, gin.H{
		"success": false,
		"message": "访问被拒绝",
		"code":    reason,
	})
	c.Abort()
}

// AntiScraping 默认反爬虫中间件
func AntiScraping() gin.HandlerFunc {
	middleware := NewAntiScrapingMiddleware(DefaultAntiScrapingConfig)
	return middleware.Middleware()
}
