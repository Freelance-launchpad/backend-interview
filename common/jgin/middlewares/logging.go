package middlewares

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	jerrorv2 "github.com/Freelance-launchpad/backend-interview/common/jerror/v2"

	"github.com/gin-gonic/gin"
)

// Logger returns a logging middleware.
func Logger(skipPaths []string) gin.HandlerFunc {
	skip := make(map[string]struct{}, len(skipPaths))
	for _, path := range skipPaths {
		skip[path] = struct{}{}
	}

	return func(c *gin.Context) {
		if _, ok := skip[c.Request.URL.Path]; ok {
			return
		}

		requestID := c.Request.Context().Value(CtxRequestID).(string)

		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := stop.Milliseconds()
		status := c.Writer.Status()
		if status == http.StatusNotImplemented {
			// 501 are returned on routes that don't exist, so no need to log them.
			return
		}
		dataLength := max(c.Writer.Size(), 0)

		attrs := []slog.Attr{
			slog.Int("status", status),
			slog.Int64("latency", latency),
			slog.String("clientIP", c.ClientIP()),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("dataLength", dataLength),
			slog.String("userAgent", c.Request.UserAgent()),
			slog.String("requestID", requestID),
		}

		msg := "request"
		if len(c.Errors) > 0 {
			msg = strings.Join(c.Errors.Errors(), ",")

			for i, ginErr := range c.Errors {
				var details []string
				if err, ok := errors.AsType[*jerrorv2.Error](ginErr.Err); ok {
					for _, detail := range err.Details {
						details = append(details, detail.String())
					}
				}
				attrs = append(attrs, slog.Group(fmt.Sprintf("error%02d", i+1), "msg", ginErr.Error(), "details", strings.Join(details, ", ")))
			}
		}

		if status >= http.StatusBadRequest && status < http.StatusInternalServerError {
			slog.LogAttrs(c.Request.Context(), slog.LevelWarn, msg, attrs...)
		}
		if status >= http.StatusInternalServerError {
			slog.LogAttrs(c.Request.Context(), slog.LevelError, msg, attrs...)
		}
	}
}
