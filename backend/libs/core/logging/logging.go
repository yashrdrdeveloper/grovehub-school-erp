// Package logging exposes domain-independent structured logging configuration.
package logging

import (
	"log/slog"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/logging/internal/src/config"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/logging/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/logging/internal/src/operations/filter"
)

type Levels = contracts.Levels

func NewLevels(levels ...slog.Level) Levels {
	return contracts.NewLevels(levels...)
}

func DefaultLevels() Levels {
	return config.DefaultLevels()
}

func ParseLevels(value string) (Levels, error) {
	return config.ParseLevels(value)
}

func NewLevelFilterHandler(next slog.Handler, levels Levels) slog.Handler {
	return filter.NewLevelFilterHandler(next, levels)
}
