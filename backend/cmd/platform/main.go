package main

import (
	"os"

	"github.com/yashrdrdeveloper/grovehub-school-erp/backend/launcher"
	"github.com/yashrdrdeveloper/grovehub-school-erp/backend/launcher/bootstrap"
)

func main() {
	logger := bootstrap.NewLogger(os.Stdout)
	logger.Info("starting GroveHub", "service", "platform", "phase", "0.1.1")

	if err := launcher.Run(logger); err != nil {
		logger.Error("GroveHub stopped with error", "service", "platform", "error", err)
		os.Exit(1)
	}
}
