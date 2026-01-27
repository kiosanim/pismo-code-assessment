package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kiosanim/pismo-code-assessment/internal/core/contextutils"
	"github.com/kiosanim/pismo-code-assessment/internal/core/logger"
	"time"
)

func LoggerMiddleware(sLogger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		method := c.Request.Method
		uri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		durationMs := fmt.Sprintf("%dms", time.Since(start))
		clientIP := c.ClientIP()
		host := c.Request.Host
		xTraceID := contextutils.GetTraceID(c.Request.Context())
		sLogger.Info(
			"HTTP Request",
			"method", method,
			"uri", uri,
			"status", statusCode,
			"duration_ms", durationMs,
			"client_ip", clientIP,
			"host", host,
			"x_trace_id", xTraceID,
		)
	}
}
