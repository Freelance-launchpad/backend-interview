package jlog

import (
	"context"
	"log/slog"
	"time"
)

type timestampHandler struct {
	handler slog.Handler
}

func (h *timestampHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *timestampHandler) Handle(ctx context.Context, record slog.Record) error {
	record.AddAttrs(slog.Int64("@timestamp", time.Now().UnixMilli()))
	return h.handler.Handle(ctx, record)
}

func (h *timestampHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &timestampHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *timestampHandler) WithGroup(name string) slog.Handler {
	return &timestampHandler{handler: h.handler.WithGroup(name)}
}
