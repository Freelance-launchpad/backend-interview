package middlewares

import (
	"fmt"
	"slices"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Freelance-launchpad/backend-interview/common/jmonitoring"
)

// Monitoring returns a middleware that sends metric to the global statsd client
func Monitoring(appName, appVersion string, checkRoutes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if slices.Contains(checkRoutes, c.Request.URL.Path) {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()
		duration := time.Since(start)
		dataLength := max(c.Writer.Size(), 0)

		tags := []string{
			fmt.Sprintf("application:%s", appName),
			fmt.Sprintf("version:%s", appVersion),
			fmt.Sprintf("route:%s", c.FullPath()),
			fmt.Sprintf("method:%s", c.Request.Method),
			fmt.Sprintf("status_code:%d", c.Writer.Status()),
			fmt.Sprintf("status_class:%dxx", c.Writer.Status()/100),
		}
		jmonitoring.Count("http.request.count", 1, tags, 1)
		jmonitoring.Timing("http.request.response.latency", duration, tags, 1)
		jmonitoring.Histogram("http.request.response.size", float64(dataLength), tags, 1)
	}
}
