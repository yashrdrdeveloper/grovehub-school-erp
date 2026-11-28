package config

import (
	"fmt"
	"log/slog"
	"strings"

	"grovehub/libs/core/logging/internal/src/constants"
	"grovehub/libs/core/logging/internal/src/contracts"
	"grovehub/libs/core/logging/internal/src/helpers"
)

func DefaultLevels() contracts.Levels {
	return contracts.NewLevels(slog.LevelInfo, slog.LevelWarn, slog.LevelError)
}

func ParseLevels(value string) (contracts.Levels, error) {
	if strings.TrimSpace(value) == "" {
		return DefaultLevels(), nil
	}

	enabled := make([]slog.Level, 0, constants.StandardLevelCount)
	seen := make(map[slog.Level]struct{}, constants.StandardLevelCount)
	for _, rawLevel := range strings.Split(value, ",") {
		levelName := strings.ToUpper(strings.TrimSpace(rawLevel))
		if levelName == "" {
			continue
		}

		level, ok := helpers.StandardLevel(levelName)
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
