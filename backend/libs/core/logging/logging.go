package logging

import (
	"log/slog"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/logging/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/logging/helpers"
)

type Levels = contracts.Levels

func DefaultLevels() Levels {
	return helpers.DefaultLevels()
}

func ParseLevels(value string) (Levels, error) {
	return helpers.ParseLevels(value)
}

func NewLevelFilterHandler(next slog.Handler, levels Levels) slog.Handler {
	return helpers.NewLevelFilterHandler(next, levels)
}
