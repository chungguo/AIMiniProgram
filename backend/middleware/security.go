package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders 安全头中间件
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止点击劫持
		c.Header("X-Frame-Options", "DENY")
		// XSS保护
		c.Header("X-XSS-Protection", "1; mode=block")
		// 内容类型嗅探保护
		c.Header("X-Content-Type-Options", "nosniff")
		// 引用策略
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		// 内容安全策略
		c.Header("Content-Security-Policy", "default-src 'self'")
		
		c.Next()
	}
}

// CORS 跨域配置
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// 允许的域名列表
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"https://your-miniprogram-domain.com",
		}
		
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin || allowedOrigin == "*" {
				allowed = true
				c.Header("Access-Control-Allow-Origin", origin)
				break
			}
		}
		
		if !allowed {
			c.Header("Access-Control-Allow-Origin", "http://localhost:3001")
		}
		
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestLogger 请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return ""
	})
}
