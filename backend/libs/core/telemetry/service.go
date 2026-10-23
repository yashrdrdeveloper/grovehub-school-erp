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
	default:
		return Service{
			Key:  normalized,
			Name: strings.ToUpper(strings.TrimSpace(value)),
		}
	}
}
