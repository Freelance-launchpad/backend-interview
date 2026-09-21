package jlog

import (
	"context"
	"log/slog"

	"github.com/Freelance-launchpad/backend-interview/common/jgin/keys"
)

type userHandler struct {
	handler slog.Handler
}

func (h *userHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *userHandler) Handle(ctx context.Context, record slog.Record) error {
	if userID := ctx.Value(keys.UserIDKey); userID != nil {
		record.AddAttrs(slog.Any("user.id", userID))
	}

	if offerID := ctx.Value(keys.OfferIDKey); offerID != nil {
		record.AddAttrs(slog.Any("user.offer_id", offerID))
	}

	if entity := ctx.Value(keys.EntityKey); entity != nil {
		record.AddAttrs(slog.Any("user.entity", entity))
	}

	return h.handler.Handle(ctx, record)
}

func (h *userHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &userHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *userHandler) WithGroup(name string) slog.Handler {
	return &userHandler{handler: h.handler.WithGroup(name)}
}
