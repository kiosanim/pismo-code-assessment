package contextutils

import (
	"context"
	"github.com/kiosanim/pismo-code-assessment/internal/core/contextkeys"
)

func GetTraceID(ctx context.Context) string {
	key := ctx.Value(contextkeys.TraceIDKey)
	if key != nil {
		traceID, ok := key.(string)
		if ok {
			return traceID
		}
	}
	return ""
}
