package contracts

import "context"

type Shutdown func(context.Context) error

type Service struct {
	Key  string
	Name string
}

var (
	ServiceLauncherEdge = Service{Key: "launcher-edge", Name: "LAUNCHER (EDGE)"}
	ServiceFoundation   = Service{Key: "foundation", Name: "FOUNDATION"}
	ServiceIdentity     = Service{Key: "identity", Name: "IDENTITY"}
	ServiceAudit        = Service{Key: "audit", Name: "AUDIT"}
)
