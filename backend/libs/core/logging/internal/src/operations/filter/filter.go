package filter

import (
	"context"
	"log/slog"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/logging/internal/src/contracts"
)

type levelFilterHandler struct {
	next   slog.Handler
	levels contracts.Levels
}

var _ slog.Handler = (*levelFilterHandler)(nil)

func NewLevelFilterHandler(next slog.Handler, levels contracts.Levels) slog.Handler {
	return &levelFilterHandler{next: next, levels: levels}
}

func (handler *levelFilterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return handler.levels.Enabled(level) && handler.next.Enabled(ctx, level)
}

func (handler *levelFilterHandler) Handle(ctx context.Context, record slog.Record) error {
	return handler.next.Handle(ctx, record)
}

func (handler *levelFilterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &levelFilterHandler{next: handler.next.WithAttrs(attrs), levels: handler.levels}
}

func (handler *levelFilterHandler) WithGroup(name string) slog.Handler {
	return &levelFilterHandler{next: handler.next.WithGroup(name), levels: handler.levels}
}
