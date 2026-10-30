package config

import (
	"os"
	"strings"
)

const defaultDKIMStateDir = ".state/mailProvider/dkim"

type DKIM struct {
	StateDir string
}

func DKIMFromEnvironment() DKIM {
	for _, name := range []string{"MAIL_INFRASTRUCTURE_DKIM_STATE_DIR", "MAIL_PROVIDER_DKIM_STATE_DIR"} {
		value := strings.TrimSpace(os.Getenv(name))
		if value != "" {
			return DKIM{StateDir: value}
		}
	}
	return DKIM{StateDir: defaultDKIMStateDir}
}
