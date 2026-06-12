package middleware

import (
	"boilerplate/config"
	"boilerplate/utils"
	"boilerplate/utils/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(utils.ErrUnauthorize("unauthorized"))
			c.Abort()

			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Error(utils.ErrUnauthorize("unauthorized"))
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		var claims jwt.MapClaims

		_, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (any, error) {
			return []byte(cfg.SecretKey), nil
		})
		if err != nil {
			logger.Log.Debugf("error: %v", err)
			c.Error(utils.ErrUnauthorize("unauthorized"))
			c.Abort()

			return
		}

		c.Set(utils.UserIDKey, int64(claims["sub"].(float64)))
		c.Set(utils.RoleKey, claims["role"].(string))
		c.Set(utils.VerifiedKey, claims["is_verified"].(bool))

		c.Next()
	}
}

func RoleAccess(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.Keys[utils.RoleKey]

		if userRole != role {
			c.Error(utils.ErrUnauthorize("unauthorized"))
			c.Abort()

			return
		}

		c.Next()
	}
}
