package contracts

import "context"

const ServiceNamespace = "grovehub"

type Shutdown func(context.Context) error

type Service struct {
	Key  string
	Name string
}

var (
	ServiceLauncherEdge = Service{Key: "launcher-edge", Name: "LAUNCHER (EDGE)"}
	ServiceFoundation   = Service{Key: "foundation", Name: "FOUNDATION"}
	ServiceIdentity     = Service{Key: "identity", Name: "IDENTITY"}
	ServiceMailProvider = Service{Key: "mail-provider", Name: "MAIL PROVIDER"}
)
