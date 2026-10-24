package telemetry

import "strings"

const ServiceNamespace = "grovehub"

type Service struct {
	Key  string
	Name string
}

var (
	ServiceLauncherEdge = Service{
		Key:  "launcher-edge",
		Name: "LAUNCHER (EDGE)",
	}
	ServiceFoundation = Service{
		Key:  "foundation",
		Name: "FOUNDATION",
	}
	ServiceIdentity = Service{
		Key:  "identity",
		Name: "IDENTITY",
	}
	ServiceMailProvider = Service{
		Key:  "mail-provider",
		Name: "MAIL PROVIDER",
	}
)

func ResolveService(value string) Service {
	normalized := strings.ToLower(strings.TrimSpace(value))

	switch normalized {
	case "", "edge", "launcher", "launcher-edge", "launcher (edge)", "launcher(edge)":
		return ServiceLauncherEdge
	case "foundation":
		return ServiceFoundation
	case "identity":
		return ServiceIdentity
	case "mail-provider", "mailprovider", "mail provider":
		return ServiceMailProvider
	default:
		return Service{
			Key:  normalized,
			Name: strings.ToUpper(strings.TrimSpace(value)),
		}
	}
}
