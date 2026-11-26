package routes

import (
	"net/http"
	auditprotected "grovehub/services/identity/internal/app/controller/audit/protected"
	capabilityprotected "grovehub/services/identity/internal/app/controller/capability/protected"
	roleprotected "grovehub/services/identity/internal/app/controller/role/protected"
	accountroleprotected "grovehub/services/identity/internal/app/controller/accountrole/protected"
	approvalprotected "grovehub/services/identity/internal/app/controller/approval/protected"

	accountprotected "grovehub/services/identity/internal/app/controller/account/protected"
	activationpublic "grovehub/services/identity/internal/app/controller/activation/public"
	authprotected "grovehub/services/identity/internal/app/controller/authentication/protected"
	authpublic "grovehub/services/identity/internal/app/controller/authentication/public"
)

// RegisterHTTP only composes annotated, generated route groups.
// Authentication and authorization are enforced by the Launcher gateway.
// @Routes /identity
func RegisterHTTP(mux *http.ServeMux, publicAuth authpublic.AuthenticationController, protectedAuth authprotected.AuthenticationController, publicActivation activationpublic.ActivationController, protectedAccount accountprotected.AccountController, roleController roleprotected.RoleController, accountroleController accountroleprotected.AccountRoleController, approvalController approvalprotected.ApprovalController, capabilityController capabilityprotected.CapabilityController, auditController auditprotected.AuditController) {
	registerPublicAuthenticationRoutes(mux, publicAuth)
	registerProtectedAuthenticationRoutes(mux, protectedAuth)
	registerPublicActivationRoutes(mux, publicActivation)
	registerProtectedAccountRoutes(mux, protectedAccount)
	registerProtectedRoleRoutes(mux, roleController)
	registerProtectedAccountRoleRoutes(mux, accountroleController)
	registerProtectedApprovalRoutes(mux, approvalController)
	registerProtectedCapabilityRoutes(mux, capabilityController)
	registerProtectedAuditRoutes(mux, auditController)
}
