package middlewares

import (
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/Freelance-launchpad/backend-interview/common/jgin/keys"
	"github.com/gin-gonic/gin"
)

// Tracing returns a middleware that adds userID and offerID (if found) in the trace.
func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		ctx := c.Request.Context()

		span, ok := tracer.SpanFromContext(ctx)
		if !ok {
			return
		}

		if userID, ok := ctx.Value(keys.UserIDKey).(string); ok && userID != "" {
			span.SetTag("request.user_id", userID)
		}

		if offerID := c.GetHeader(keys.HeaderOfferID); offerID != "" {
			span.SetTag("request.offer_id", offerID)
		}
	}
}
