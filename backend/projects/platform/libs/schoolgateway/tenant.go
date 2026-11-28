package schoolgateway

import "grovehub/launcher/services/gateway/session"

type TenantCheck = session.TenantCheck
type TenantState = session.TenantState
type TenantStateLookup = session.TenantStateLookup

var ErrTenantDenied = session.ErrTenantDenied
var ActiveTenantCheck = session.ActiveTenantCheck
var WithTenantCheck = session.WithTenantCheck
var TenantGeneration = session.TenantGeneration
var CheckTenant = session.CheckTenant
var CheckTenantSession = session.CheckTenantSession
var CheckAccountTenant = session.CheckAccountTenant
var CheckAccountGeneration = session.CheckAccountGeneration
