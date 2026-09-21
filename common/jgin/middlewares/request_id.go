package middlewares

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContextKey int64

const CtxRequestID ContextKey = iota

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-Id")
		if id == "" {
			id = uuid.Must(uuid.NewUUID()).String()
		}

		c.Request = c.Request.WithContext(
			context.WithValue(c.Request.Context(), CtxRequestID, id),
		)
		c.Next()
	}
}

func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(CtxRequestID).(string); ok {
		return id
	}
	return ""
}
