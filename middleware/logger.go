package middleware

import (
	"boilerplate/utils/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		// stop timer
		finish := time.Since(start)
		reqMethod := c.Request.Method
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		if raw != "" {
			path += "?" + raw
		}

		err := c.Errors.Last()
		if err != nil {
			logger.Log.Infow(err.Error(),
				zap.String("method", reqMethod),
				zap.String("uri", path),
				zap.Int("code", statusCode),
				zap.Duration("latency", finish),
				zap.String("Client IP", clientIP),
			)

			return
		}

		logger.Log.Infow("Success Request",
			zap.String("method", reqMethod),
			zap.String("uri", path),
			zap.Int("code", statusCode),
			zap.Duration("latency", finish),
			zap.String("Client IP", clientIP),
		)
	}
}
