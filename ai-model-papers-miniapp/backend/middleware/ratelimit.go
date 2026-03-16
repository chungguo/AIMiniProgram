package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiterConfig 限流配置
type RateLimiterConfig struct {
	RequestsPerSecond float64
	BurstSize         int
}

// DefaultRateLimiterConfig 默认限流配置
var DefaultRateLimiterConfig = RateLimiterConfig{
	RequestsPerSecond: 10, // 每秒10个请求
	BurstSize:         20, // 突发20个请求
}

// IPRateLimiter 基于IP的限流器
type IPRateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	config   RateLimiterConfig
}

// NewIPRateLimiter 创建新的IP限流器
func NewIPRateLimiter(config RateLimiterConfig) *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		config:   config,
	}
}

// GetLimiter 获取指定IP的限流器
func (rl *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[ip]
	rl.mu.RUnlock()

	if exists {
		return limiter
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	// 双重检查
	if limiter, exists := rl.limiters[ip]; exists {
		return limiter
	}

	limiter = rate.NewLimiter(rate.Limit(rl.config.RequestsPerSecond), rl.config.BurstSize)
	rl.limiters[ip] = limiter
	return limiter
}

// Cleanup 清理过期的限流器（可定时调用）
func (rl *IPRateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	// 实际项目中可以实现过期清理逻辑
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(config RateLimiterConfig) gin.HandlerFunc {
	limiter := NewIPRateLimiter(config)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			ip = c.Request.RemoteAddr
		}

		if !limiter.GetLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "请求过于频繁，请稍后再试",
				"code":    "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// StrictRateLimitMiddleware 严格的限流（用于敏感接口）
func StrictRateLimitMiddleware() gin.HandlerFunc {
	return RateLimitMiddleware(RateLimiterConfig{
		RequestsPerSecond: 1, // 每秒1个请求
		BurstSize:         3, // 突发3个请求
	})
}
