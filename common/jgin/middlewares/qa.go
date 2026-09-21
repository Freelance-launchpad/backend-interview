package middlewares

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"

	"github.com/Freelance-launchpad/backend-interview/common/jenv"
)

func OnlyInEnv(environments ...jenv.Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !slices.Contains(environments, jenv.Get()) {
			c.AbortWithStatus(http.StatusNotImplemented)
		}
	}
}
