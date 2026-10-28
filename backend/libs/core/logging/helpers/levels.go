package helpers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"grovehub/libs/core/logging/contracts"
)

func DefaultLevels() contracts.Levels {
	return contracts.NewLevels(slog.LevelInfo, slog.LevelWarn, slog.LevelError)
}

func ParseLevels(value string) (contracts.Levels, error) {
	if strings.TrimSpace(value) == "" {
		return DefaultLevels(), nil
	}

	enabled := make([]slog.Level, 0, 4)
	seen := make(map[slog.Level]struct{}, 4)
	for _, rawLevel := range strings.Split(value, ",") {
		levelName := strings.ToUpper(strings.TrimSpace(rawLevel))
		if levelName == "" {
			continue
		}

		level, ok := standardLevel(levelName)
		if !ok {
			return contracts.Levels{}, fmt.Errorf(
				"unsupported log level %q; allowed levels: DEBUG, INFO, WARN, ERROR",
				rawLevel,
			)
		}
		if _, exists := seen[level]; exists {
			continue
		}
		seen[level] = struct{}{}
		enabled = append(enabled, level)
	}

	if len(enabled) == 0 {
		return contracts.Levels{}, fmt.Errorf("at least one log level is required")
	}
	return contracts.NewLevels(enabled...), nil
}

func NewLevelFilterHandler(next slog.Handler, levels contracts.Levels) slog.Handler {
	return &levelFilterHandler{next: next, levels: levels}
}

type levelFilterHandler struct {
	next   slog.Handler
	levels contracts.Levels
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

func standardLevel(name string) (slog.Level, bool) {
	switch name {
	case "DEBUG":
		return slog.LevelDebug, true
	case "INFO":
		return slog.LevelInfo, true
	case "WARN":
		return slog.LevelWarn, true
	case "ERROR":
		return slog.LevelError, true
	default:
		return 0, false
	}
}
