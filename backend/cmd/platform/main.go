package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"grovehub/launcher"
	"grovehub/launcher/bootstrap"
	coretelemetry "grovehub/libs/core/telemetry"
)

func main() {
	os.Exit(run())
}

func run() int {
	logger, err := bootstrap.NewLogger(os.Stdout)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "invalid GROVEHUB_LOG_LEVELS: %v\n", err)
		return 1
	}

	shutdownTelemetry := coretelemetry.Initialize(context.Background(), logger)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := shutdownTelemetry(ctx); err != nil {
			logger.Error("OpenTelemetry shutdown failed", "service", "platform", "error", err)
		}
	}()

	logger.Info("starting GroveHub", "service", "platform", "phase", "0.1.1")

	if err := launcher.Run(logger); err != nil {
		logger.Error("GroveHub stopped with error", "service", "platform", "error", err)
		return 1
	}

	return 0
}
