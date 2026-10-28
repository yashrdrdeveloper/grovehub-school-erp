package contracts

import "log/slog"

type Levels struct {
	enabled map[slog.Level]struct{}
}

func NewLevels(levels ...slog.Level) Levels {
	enabled := make(map[slog.Level]struct{}, len(levels))
	for _, level := range levels {
		enabled[level] = struct{}{}
	}
	return Levels{enabled: enabled}
}

func (levels Levels) Enabled(level slog.Level) bool {
	_, ok := levels.enabled[level]
	return ok
}
