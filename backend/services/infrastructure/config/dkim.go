package config

import (
	"os"
	"strings"
)

const defaultDKIMStateDir = ".state/infrastructure/dkim"

type DKIM struct {
	StateDir string
}

func DKIMFromEnvironment() DKIM {
	value := strings.TrimSpace(os.Getenv("MAIL_INFRASTRUCTURE_DKIM_STATE_DIR"))
	if value != "" {
		return DKIM{StateDir: value}
	}
	return DKIM{StateDir: defaultDKIMStateDir}
}
