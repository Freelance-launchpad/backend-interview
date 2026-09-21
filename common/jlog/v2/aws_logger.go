package jlog

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/smithy-go/logging"
)

type AWSLogger struct {
	context.Context
}

var _ logging.Logger = &AWSLogger{}
var _ logging.ContextLogger = &AWSLogger{}

func (l *AWSLogger) Logf(classification logging.Classification, format string, v ...any) {
	ctx := context.Background()
	if l.Context != nil {
		ctx = l.Context
	}

	level := slog.LevelDebug
	if classification == logging.Warn {
		level = slog.LevelWarn
	}

	slog.Log(ctx, level, fmt.Sprintf(format, v...))
}

func (l *AWSLogger) WithContext(ctx context.Context) logging.Logger {
	return &AWSLogger{ctx}
}
