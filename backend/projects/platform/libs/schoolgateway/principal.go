// Package schoolgateway is a temporary compatibility bridge during the gateway migration.
package schoolgateway

import "grovehub/launcher/services/gateway/session"

type Principal = session.Principal
type AccessVerifier = session.AccessVerifier
type ProtectedRouteAuthorizer = session.ProtectedRouteAuthorizer

var WithPrincipal = session.WithPrincipal
var PrincipalFrom = session.PrincipalFrom
