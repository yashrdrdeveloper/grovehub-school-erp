package mappers

import (
	"strings"

	"grovehub/libs/core/telemetry/contracts"
)

func ResolveService(value string) contracts.Service {
	normalized := strings.ToLower(strings.TrimSpace(value))

	switch normalized {
	case "", "edge", "launcher", "launcher-edge", "launcher (edge)", "launcher(edge)":
		return contracts.ServiceLauncherEdge
	case "foundation":
		return contracts.ServiceFoundation
	case "identity":
		return contracts.ServiceIdentity
	case "mail-provider", "mailprovider", "mail provider":
		return contracts.ServiceMailProvider
	default:
		return contracts.Service{
			Key:  normalized,
			Name: strings.ToUpper(strings.TrimSpace(value)),
		}
	}
}
