package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Should set CORS headers on a standard GET request", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		_, r := gin.CreateTestContext(recorder)

		r.Use(CORSMiddleware())
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "*", recorder.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", recorder.Header().Get("Access-Control-Allow-Credentials"))
		assert.Contains(t, recorder.Header().Get("Access-Control-Allow-Methods"), "GET")
	})

	t.Run("Should handle OPTIONS preflight request", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		_, r := gin.CreateTestContext(recorder)

		r.Use(CORSMiddleware())
		r.GET("/test", func(c *gin.Context) {})

		req, _ := http.NewRequest("OPTIONS", "/test", nil)
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusNoContent, recorder.Code) // 204 status
		assert.Equal(t, "*", recorder.Header().Get("Access-Control-Allow-Origin"))
	})
}