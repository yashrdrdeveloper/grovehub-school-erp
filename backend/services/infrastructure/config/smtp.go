package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultSMTPHost      = "127.0.0.1"
	defaultSMTPPort      = 2525
	defaultSubmitTimeout = 15 * time.Second
	defaultHealthTimeout = 3 * time.Second
)

type SMTP struct {
	Host          string
	Port          int
	SubmitTimeout time.Duration
	HealthTimeout time.Duration
}

func SMTPFromEnvironment() SMTP {
	return SMTP{
		Host: firstEnvironmentValue(
			[]string{"MAIL_INFRASTRUCTURE_SMTP_HOST", "MAIL_PROVIDER_SMTP_HOST"},
			defaultSMTPHost,
		),
		Port: environmentPort(
			[]string{"MAIL_INFRASTRUCTURE_SMTP_PORT", "MAIL_PROVIDER_SMTP_PORT"},
			defaultSMTPPort,
		),
		SubmitTimeout: defaultSubmitTimeout,
		HealthTimeout: defaultHealthTimeout,
	}
}

func firstEnvironmentValue(names []string, fallback string) string {
	for _, name := range names {
		value := strings.TrimSpace(os.Getenv(name))
		if value != "" {
			return value
		}
	}
	return fallback
}

func environmentPort(names []string, fallback int) int {
	for _, name := range names {
		value := strings.TrimSpace(os.Getenv(name))
		if value == "" {
			continue
		}
		port, err := strconv.Atoi(value)
		if err == nil && port > 0 && port <= 65535 {
			return port
		}
	}
	return fallback
}
