package jlog

import (
	"log/slog"
	"os"

	slogtrace "github.com/DataDog/dd-trace-go/contrib/log/slog/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jenv"
)

func Init(appName, appVersion string) {
	var slogHandler slog.Handler
	options := &slog.HandlerOptions{Level: slog.LevelDebug}
	if jenv.InLocal() {
		slogHandler = slog.NewTextHandler(os.Stdout, options)
	} else {
		slogHandler = slog.NewJSONHandler(os.Stdout, options)
	}

	slogHandler = slogHandler.WithAttrs([]slog.Attr{
		slog.String("dd.env", jenv.Name()),
		slog.String("dd.service", appName),
		slog.String("dd.version", appVersion),
	})

	slogHandler = slogtrace.WrapHandler(slogHandler)
	slogHandler = &timestampHandler{handler: slogHandler}
	slogHandler = &userHandler{handler: slogHandler}

	logger := slog.New(slogHandler)

	slog.SetDefault(logger)
}
