package routes

import (
    "net/http"
    protectedcontroller "grovehub/services/identity/controller/protected"
    publiccontroller "grovehub/services/identity/controller/public"
)

// RegisterHTTP keeps public authentication and protected admin registration
// on explicit, non-overlapping routes. Registration authorization is also
// enforced inside the service; a route name alone is not an access control.
func RegisterHTTP(mux *http.ServeMux, protected protectedcontroller.AuthController, public publiccontroller.AuthController, protectedMiddleware func(http.Handler) http.Handler) {
    mux.HandleFunc("POST /identity/public/auth/login", public.Login)
    mux.HandleFunc("POST /identity/public/auth/set-password", public.SetPassword)
    mux.Handle("POST /identity/protected/auth/register", protectedMiddleware(http.HandlerFunc(protected.Register)))
}
