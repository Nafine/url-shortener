package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

func Logger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Info(fmt.Sprintf(
			"%s %s with time of %d miliseconds",
			c.Request.Method,
			c.Request.URL.Path,
			time.Since(start).Milliseconds(),
		))
	}
}
