package gateway

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func LoggingMiddleware(logger *zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger.Info().
			Str("method", method).
			Str("path", path).
			Int("status", status).
			Dur("latency", latency).
			Msg("HTTP request")
	}
}

func JWTAuthMiddleware(validator *AuthValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			validator.logger.Error().Msg("No authorization token provided")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token required"})
			c.Abort()
			return
		}

		if err := validator.ValidateAdminJWT(token); err != nil {
			validator.logger.Error().Err(err).Msg("JWT validation failed")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or unauthorized token"})
			c.Abort()
			return
		}

		// Store token in context for downstream gRPC calls
		c.Set("jwt_token", token)
		c.Next()
	}
}

func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				switch err.Type {
				case gin.ErrorTypeBind:
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
				case gin.ErrorTypePublic:
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				}
			}
		}
	}
}
