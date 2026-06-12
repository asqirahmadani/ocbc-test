package middleware

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan bool)

		go func() {
			c.Next()

			done <- true
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.Canceled) {
				c.Error(context.Canceled)
				c.Abort()
			} else if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				c.Error(context.DeadlineExceeded)
				c.Abort()
			}

			return
		}
	}
}
