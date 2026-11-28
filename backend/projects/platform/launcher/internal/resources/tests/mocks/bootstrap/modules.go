package bootstrap

import (
	"context"
	"net/http"

	"grovehub/libs/platform/session"
	person "grovehub/services/foundation/client/person"
	account "grovehub/services/identity/client/account"
	"grovehub/services/infrastructure/client/messaging"
	connection "grovehub/services/integration/client/connection"
)

// @Mock
type InfrastructureModule interface {
	MessagingClient() messaging.MessagingClient
	EventStore() messaging.EventStore
	Close() error
}

// @Mock
type FoundationModule interface {
	RegisterHTTP(*http.ServeMux)
	Client() person.Client
	Close()
}

// @Mock
type IntegrationModule interface {
	RegisterHTTP(*http.ServeMux)
	Client() connection.Client
	VerifyWorkload(context.Context, *http.Request) (context.Context, error)
	LocalOnly() bool
	Close()
}

// @Mock
type IdentityModule interface {
	RegisterHTTP(*http.ServeMux)
	VerifyAccess(context.Context, string) (session.Principal, error)
	VerifyDelegatedSession(context.Context, string) (session.Principal, error)
	AuthorizeProtectedRoute(context.Context, string) error
	AccountClient() account.Client
	RouteLimits() map[string]int
	Close()
}

// @Mock
type AuditModule interface {
	RegisterHTTP(*http.ServeMux)
	Ready() bool
	Close()
}

// @Mock
type MessagingClient interface {
	Publish(context.Context, string, []byte, map[string]string) error
	PublishWithOptions(context.Context, string, []byte, map[string]string, messaging.PublishOptions) error
	Consume(context.Context, string, messaging.Handler) error
	Health(context.Context) error
}

// @Mock
type Server interface {
	Run(context.Context) error
}
