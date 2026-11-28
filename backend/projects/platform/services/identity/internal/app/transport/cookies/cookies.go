// Package cookies defines the presentation-only HTTP cookie contract.
// Session verification, Origin/CSRF policy, and rate limits belong to Launcher.
package cookies

import (
	"net/http"

	schoolgateway "grovehub/launcher/session"
)

const (
	AccessCookie  = schoolgateway.AccessCookie
	RefreshCookie = schoolgateway.RefreshCookie
	CSRFName      = schoolgateway.CSRFName
	ResetCookie   = schoolgateway.ResetCookie
	SetupCookie   = schoolgateway.SetupCookie
)

// Manager is the resource transport interface implemented in impl/.
// Use this contract in controller constructors; no package-level helpers.
type Manager interface {
	Issue(http.ResponseWriter, string, string, string)
	ClearSession(http.ResponseWriter)
	Read(*http.Request, string) string
	SetGrant(http.ResponseWriter, string, string, string)
	ClearGrant(http.ResponseWriter, string, string)
	Proof(*http.Request, string) string
}
