package logging

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

type Levels struct {
	enabled map[slog.Level]struct{}
}

func DefaultLevels() Levels {
	return newLevels(slog.LevelInfo, slog.LevelWarn, slog.LevelError)
}

func ParseLevels(value string) (Levels, error) {
	if strings.TrimSpace(value) == "" {
		return DefaultLevels(), nil
	}

	enabled := make(map[slog.Level]struct{})
	for _, rawLevel := range strings.Split(value, ",") {
		levelName := strings.ToUpper(strings.TrimSpace(rawLevel))
		if levelName == "" {
			continue
		}

		level, ok := standardLevel(levelName)
		if !ok {
			return Levels{}, fmt.Errorf(
				"unsupported log level %q; allowed levels: DEBUG, INFO, WARN, ERROR",
				rawLevel,
			)
		}

		enabled[level] = struct{}{}
	}

	if len(enabled) == 0 {
		return Levels{}, fmt.Errorf("at least one log level is required")
	}

	return Levels{enabled: enabled}, nil
}

func (levels Levels) Enabled(level slog.Level) bool {
	_, ok := levels.enabled[level]
	return ok
}

func NewLevelFilterHandler(next slog.Handler, levels Levels) slog.Handler {
	return &levelFilterHandler{
		next:   next,
		levels: levels,
	}
}

type levelFilterHandler struct {
	next   slog.Handler
	levels Levels
}

func (handler *levelFilterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return handler.levels.Enabled(level) && handler.next.Enabled(ctx, level)
}

func (handler *levelFilterHandler) Handle(ctx context.Context, record slog.Record) error {
	return handler.next.Handle(ctx, record)
}

func (handler *levelFilterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &levelFilterHandler{
		next:   handler.next.WithAttrs(attrs),
		levels: handler.levels,
	}
}

func (handler *levelFilterHandler) WithGroup(name string) slog.Handler {
	return &levelFilterHandler{
		next:   handler.next.WithGroup(name),
		levels: handler.levels,
	}
}

func newLevels(levels ...slog.Level) Levels {
	enabled := make(map[slog.Level]struct{}, len(levels))
	for _, level := range levels {
		enabled[level] = struct{}{}
	}
	return Levels{enabled: enabled}
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
