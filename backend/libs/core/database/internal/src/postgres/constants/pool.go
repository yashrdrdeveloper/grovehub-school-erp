package constants

import "time"

const (
	ConnectionTimeout = 5 * time.Second
	MaxOpenConns      = 10
	MaxIdleConns      = 2
	ConnMaxIdleTime   = 5 * time.Minute
	ConnMaxLifetime   = 30 * time.Minute
)
