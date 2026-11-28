package serve

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	coretelemetry "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/telemetry"
	"grovehub/cmd/internal/src/command"
	"grovehub/launcher"
)

type Dependencies struct {
	Logger     func(io.Writer) (*slog.Logger, error)
	Initialize func(context.Context, *slog.Logger) coretelemetry.Shutdown
	Run        func(*slog.Logger) error
}

// @Strategy
type strategy struct {
	dependencies Dependencies
}

func New() command.Command {
	return NewWithDependencies(Dependencies{
		Logger:     launcher.NewLogger,
		Initialize: coretelemetry.Initialize,
		Run:        launcher.Run,
	})
}

func NewWithDependencies(dependencies Dependencies) command.Command {
	return &strategy{dependencies: dependencies}
}

func (s *strategy) Run(args []string, output, diagnostics io.Writer) int {
	if len(args) != 0 {
		_, _ = fmt.Fprintln(diagnostics, "usage: groovehub-platform serve")
		return 2
	}

	logger, err := s.dependencies.Logger(output)
	if err != nil {
		_, _ = fmt.Fprintf(diagnostics, "invalid GROVEHUB_LOG_LEVELS: %v\n", err)
		return 1
	}

	shutdownTelemetry := s.dependencies.Initialize(context.Background(), logger)
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

	if err := s.dependencies.Run(logger); err != nil {
		logger.Error(
			"GroveHub stopped with error",
			"service", coretelemetry.ServiceLauncherEdge.Name,
			"error", err,
		)
		return 1
	}

	return 0
}
