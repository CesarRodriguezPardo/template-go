package middleware

import (
	logger "CesarRodriguezPardo/template-go/infra/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLogger is a custom logging middleware that uses Zap
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log after processing
		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
			zap.String("user-agent", c.Request.UserAgent()),
		}

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error(e, nil, fields...)
			}
		} else {
			if status >= 400 && status < 500 {
				logger.Warn(path, fields...)
			} else if status >= 500 {
				logger.Error(path, nil, fields...)
			} else {
				logger.Info(path, fields...)
			}
		}
	}
}
