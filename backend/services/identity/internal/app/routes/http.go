package routes
import(
 "net/http"
 "time"
 protectedAuth "grovehub/services/identity/internal/app/controller/authentication/protected"
 publicAuth "grovehub/services/identity/internal/app/controller/authentication/public"
 protectedAccount "grovehub/services/identity/internal/app/controller/account/protected"
 publicAccount "grovehub/services/identity/internal/app/controller/account/public"
 "grovehub/services/identity/internal/app/middleware"
)
// @Routes /identity
func RegisterHTTP(mux *http.ServeMux,authPublic publicAuth.AuthenticationController,authProtected protectedAuth.AuthenticationController,accountPublic publicAccount.AccountController,accountProtected protectedAccount.AccountController,m middleware.Authentication){
 limiter:=middleware.NewLimiter()
 public:=m.Origin
 tokenProtected:=func(h http.Handler)http.Handler{return m.Origin(m.CSRF(h))}
 protected:=m.Protected
 mux.Handle("POST /identity/public/auth/login",public(limiter.Guard("Login",12,time.Minute,http.HandlerFunc(authPublic.Login))))
 mux.Handle("POST /identity/public/auth/forgot-password",public(limiter.Guard("ForgotPassword",5,time.Minute,http.HandlerFunc(authPublic.ForgotPassword))))
 mux.Handle("POST /identity/public/auth/recovery/exchange",public(limiter.Guard("RecoveryExchange",8,time.Minute,http.HandlerFunc(authPublic.RecoveryExchange))))
 mux.Handle("POST /identity/public/auth/invitations/exchange",public(limiter.Guard("InvitationExchange",8,time.Minute,http.HandlerFunc(authPublic.InvitationExchange))))
 mux.Handle("POST /identity/public/auth/refresh",tokenProtected(limiter.Guard("Refresh",30,time.Minute,http.HandlerFunc(authPublic.Refresh))))
 mux.Handle("POST /identity/public/auth/reset-password",tokenProtected(http.HandlerFunc(authPublic.ResetPassword)))
 mux.Handle("POST /identity/public/accounts/setup",tokenProtected(http.HandlerFunc(accountPublic.Setup)))
 mux.Handle("POST /identity/protected/auth/logout",protected(http.HandlerFunc(authProtected.Logout)))
 mux.Handle("POST /identity/protected/auth/logout-all",protected(http.HandlerFunc(authProtected.LogoutAll)))
 mux.Handle("POST /identity/protected/auth/change-password",protected(http.HandlerFunc(authProtected.ChangePassword)))
 mux.Handle("POST /identity/protected/auth/reauthenticate",protected(http.HandlerFunc(authProtected.Reauthenticate)))
 mux.Handle("GET /identity/protected/auth/sessions",protected(http.HandlerFunc(authProtected.Sessions)))
 mux.Handle("DELETE /identity/protected/auth/sessions/{sessionId}",protected(http.HandlerFunc(authProtected.RevokeSession)))
 mux.Handle("GET /identity/protected/auth/me",protected(http.HandlerFunc(authProtected.Me)))
 mux.Handle("POST /identity/protected/accounts",protected(http.HandlerFunc(accountProtected.Create)))
 mux.Handle("POST /identity/protected/accounts/search",protected(http.HandlerFunc(accountProtected.Search)))
 mux.Handle("GET /identity/protected/accounts/{accountId}",protected(http.HandlerFunc(accountProtected.Get)))
 mux.Handle("PATCH /identity/protected/accounts/{accountId}",protected(http.HandlerFunc(accountProtected.Update)))
 mux.Handle("PATCH /identity/protected/accounts/{accountId}/suspend",protected(http.HandlerFunc(accountProtected.Suspend)))
 mux.Handle("PATCH /identity/protected/accounts/{accountId}/activate",protected(http.HandlerFunc(accountProtected.Activate)))
 mux.Handle("PATCH /identity/protected/accounts/{accountId}/deactivate",protected(http.HandlerFunc(accountProtected.Deactivate)))
 mux.Handle("POST /identity/protected/accounts/{accountId}/resend-invitation",protected(http.HandlerFunc(accountProtected.Resend)))
}
