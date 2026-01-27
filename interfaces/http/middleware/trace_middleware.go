package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kiosanim/pismo-code-assessment/internal/core/contextkeys"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		xTrace := c.GetHeader(contextkeys.TraceIDKey)
		if xTrace == "" {
			xTrace = uuid.NewString()
		}
		c.Writer.Header().Set(contextkeys.TraceIDKey, xTrace)
		ctx := context.WithValue(c.Request.Context(), contextkeys.TraceIDKey, xTrace)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
