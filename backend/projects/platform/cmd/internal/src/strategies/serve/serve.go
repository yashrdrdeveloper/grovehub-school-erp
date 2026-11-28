package serve

import (
	"context"
	"fmt"
	"io"
	"time"

	coretelemetry "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry"
	"grovehub/cmd/internal/src/command"
	"grovehub/launcher"
)

// @Strategy
type strategy struct{}

func New() command.Command {
	return strategy{}
}

func (strategy) Run(args []string, output, diagnostics io.Writer) int {
	if len(args) != 0 {
		_, _ = fmt.Fprintln(diagnostics, "usage: groovehub-platform serve")
		return 2
	}

	logger, err := launcher.NewLogger(output)
	if err != nil {
		_, _ = fmt.Fprintf(diagnostics, "invalid GROVEHUB_LOG_LEVELS: %v\n", err)
		return 1
	}

	shutdownTelemetry := coretelemetry.Initialize(context.Background(), logger)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := shutdownTelemetry(ctx); err != nil {
			logger.Error(
				"OpenTelemetry shutdown failed",
				"service", coretelemetry.ServiceLauncherEdge.Name,
				"error", err,
			)
		}
	}()

	logger.Info(
		"starting GroveHub",
		"service", coretelemetry.ServiceLauncherEdge.Name,
		"serviceKey", coretelemetry.ServiceLauncherEdge.Key,
		"phase", "0.1.1",
	)

	if err := launcher.Run(logger); err != nil {
		logger.Error(
			"GroveHub stopped with error",
			"service", coretelemetry.ServiceLauncherEdge.Name,
			"error", err,
		)
		return 1
	}

	return 0
}
