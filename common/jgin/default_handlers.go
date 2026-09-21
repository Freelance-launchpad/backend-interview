package jgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ping(ctx *gin.Context) {
	ctx.String(http.StatusOK, "pong\n")
}

type healthResponse struct {
	Application string `json:"application"`
	Version     string `json:"version"`
	BuildDate   string `json:"build_date"`
}

func health(appName, version, buildDate string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		response := healthResponse{
			Application: appName,
			Version:     version,
			BuildDate:   buildDate,
		}
		ctx.JSON(http.StatusOK, response)
	}
}

func noRoute(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
