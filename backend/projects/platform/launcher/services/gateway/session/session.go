// Deprecated: use grovehub/launcher/session. Kept during consumer migration.
package session
import canonical "grovehub/launcher/session"
type Principal = canonical.Principal
type AccessVerifier = canonical.AccessVerifier
type ProtectedRouteAuthorizer = canonical.ProtectedRouteAuthorizer
type TenantCheck = canonical.TenantCheck
type TenantState = canonical.TenantState
type TenantStateLookup = canonical.TenantStateLookup
const (AccessCookie=canonical.AccessCookie;RefreshCookie=canonical.RefreshCookie;CSRFName=canonical.CSRFName;ResetCookie=canonical.ResetCookie;SetupCookie=canonical.SetupCookie)
var ErrTenantDenied=canonical.ErrTenantDenied
var WithPrincipal=canonical.WithPrincipal
var PrincipalFrom=canonical.PrincipalFrom
var ActiveTenantCheck=canonical.ActiveTenantCheck
var WithTenantCheck=canonical.WithTenantCheck
var TenantGeneration=canonical.TenantGeneration
var CheckTenant=canonical.CheckTenant
var CheckTenantSession=canonical.CheckTenantSession
var CheckAccountTenant=canonical.CheckAccountTenant
var CheckAccountGeneration=canonical.CheckAccountGeneration
