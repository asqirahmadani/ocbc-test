package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"boilerplate/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLoggerMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Should log successful requests", func(t *testing.T) {
		core, recorded := observer.New(zap.InfoLevel)
		logger.Log = zap.New(core).Sugar()

		recorder := httptest.NewRecorder()
		_, r := gin.CreateTestContext(recorder)

		r.Use(LoggerMiddleware())
		r.GET("/test-logger", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/test-logger?query=1", nil)
		r.ServeHTTP(recorder, req)

		assert.Equal(t, 1, recorded.Len())
		logEntry := recorded.All()[0]

		assert.Equal(t, "Success Request", logEntry.Message)
		
		fields := logEntry.ContextMap()
		assert.Equal(t, "GET", fields["method"])
		assert.Equal(t, "/test-logger?query=1", fields["uri"])
		assert.Equal(t, int64(http.StatusOK), fields["code"])
		assert.NotNil(t, fields["latency"])
	})

	t.Run("Should log errors when they exist in context", func(t *testing.T) {
		core, recorded := observer.New(zap.InfoLevel)
		logger.Log = zap.New(core).Sugar()

		recorder := httptest.NewRecorder()
		_, r := gin.CreateTestContext(recorder)

		r.Use(LoggerMiddleware())
		r.GET("/error-logger", func(c *gin.Context) {
			c.Error(errors.New("database connection failed"))
			c.Status(http.StatusInternalServerError)
		})

		req, _ := http.NewRequest("GET", "/error-logger", nil)
		r.ServeHTTP(recorder, req)

		assert.Equal(t, 1, recorded.Len())
		logEntry := recorded.All()[0]

		assert.Equal(t, "database connection failed", logEntry.Message)
		
		fields := logEntry.ContextMap()
		assert.Equal(t, int64(http.StatusInternalServerError), fields["code"])
	})
}