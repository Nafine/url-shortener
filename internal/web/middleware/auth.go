package middleware

import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

func APIKeyAuth(log *slog.Logger, validKeys map[string]string) gin.HandlerFunc {
	const op = "handler.Delete"
	return func(c *gin.Context) {
		log := log.With("operation", op)

		key := c.GetHeader("X-API-Key")
		if user, ok := validKeys[key]; ok {
			c.Set("user", user)
			c.Next()
			return
		}

		log.Info("Unsuccessful authentication", "Key", key)
		c.AbortWithStatus(401)
	}
}
