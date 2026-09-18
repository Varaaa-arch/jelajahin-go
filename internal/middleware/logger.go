package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLoggerMiddleware mencatat setiap request (method, path, status, durasi,
// dan IP client) untuk keperluan monitoring.
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		status := c.Writer.Status()
		log.Printf("[REQ] %s %s -> %d (%.2fms) client=%s",
			c.Request.Method,
			c.Request.URL.Path,
			status,
			float64(time.Since(start).Microseconds())/1000.0,
			c.ClientIP(),
		)
	}
}