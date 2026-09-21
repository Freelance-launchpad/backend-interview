package jgin

import (
	"context"
	"errors"
	"net/http"
	"slices"

	gintrace "github.com/DataDog/dd-trace-go/contrib/gin-gonic/gin/v2"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Freelance-launchpad/backend-interview/common/jenv"
	"github.com/Freelance-launchpad/backend-interview/common/jgin/middlewares"
)

type Engine struct {
	*gin.Engine
	srv *http.Server
}

func (e *Engine) Start() error {
	err := e.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (e *Engine) Stop(ctx context.Context) error {
	return e.srv.Shutdown(ctx)
}

func New(address, appName, version, buildDate string) *Engine {
	if !jenv.InLocal() {
		gin.SetMode(gin.ReleaseMode)
	}

	checkRoutes := []string{
		"/ping",
		"/health",
		"/metrics",
		"/" + appName,
	}

	e := gin.New()
	e.Use(
		gin.Recovery(),
		middlewares.RequestID(),
		gintrace.Middleware(appName,
			gintrace.WithIgnoreRequest(func(c *gin.Context) bool {
				return slices.Contains(checkRoutes, c.Request.URL.Path)
			}),
			gintrace.WithUseGinErrors(),
		),
		middlewares.Logger(checkRoutes),
		middlewares.Monitoring(appName, version, checkRoutes),
		middlewares.Tracing(),
	)

	e.NoRoute(noRoute)
	e.GET("/health", health(appName, version, buildDate))
	e.GET("/metrics", gin.WrapH(promhttp.Handler()))
	e.GET("/ping", ping)
	e.GET("/"+appName, ping)
	return &Engine{
		Engine: e,
		srv: &http.Server{
			Addr:    address,
			Handler: e,
		},
	}
}
