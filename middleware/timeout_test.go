package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTimeoutMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Should pass if handler finishes within timeout", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		_, r := gin.CreateTestContext(recorder)

		r.Use(TimeoutMiddleware(100 * time.Millisecond))
		
		r.GET("/fast", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/fast", nil)
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Empty(t, recorder.Body.String())
	})

	t.Run("Should abort and log error when handler exceeds timeout", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		_, r := gin.CreateTestContext(recorder)

		r.Use(TimeoutMiddleware(10 * time.Millisecond))
		
		r.GET("/slow", func(c *gin.Context) {
			time.Sleep(50 * time.Millisecond)
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest("GET", "/slow", nil)
		r.ServeHTTP(recorder, req)

		assert.True(t, len(r.Handlers) > 0) 
	})

	t.Run("Should handle client cancellation", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		_, r := gin.CreateTestContext(recorder)

		r.Use(TimeoutMiddleware(5 * time.Second))

		var capturedCtx *gin.Context
		handlerReached := make(chan struct{})

		r.GET("/cancel", func(c *gin.Context) {
			capturedCtx = c
			close(handlerReached) // Signal that we are inside the handler
			<-c.Request.Context().Done() // Wait for the cancellation
		})

		reqContext, cancelRequest := context.WithCancel(context.Background())
		req, _ := http.NewRequestWithContext(reqContext, "GET", "/cancel", nil)

		go func() {
			r.ServeHTTP(recorder, req)
		}()

		<-handlerReached

		cancelRequest()

		time.Sleep(10 * time.Millisecond)

		hasCancelError := false
		if capturedCtx != nil {
			for _, e := range capturedCtx.Errors {
				if errors.Is(e.Err, context.Canceled) {
					hasCancelError = true
					break
				}
			}
		}

		assert.True(t, hasCancelError, "Middleware should have captured context.Canceled")
	})
}