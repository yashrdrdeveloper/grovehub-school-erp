package impl
import(
 ap "grovehub/services/identity/internal/app/controller/account/protected"
 api "grovehub/services/identity/internal/app/controller/account/protected/impl"
 au "grovehub/services/identity/internal/app/controller/account/public"
 aui "grovehub/services/identity/internal/app/controller/account/public/impl"
 hp "grovehub/services/identity/internal/app/controller/authentication/protected"
 hpi "grovehub/services/identity/internal/app/controller/authentication/protected/impl"
 hu "grovehub/services/identity/internal/app/controller/authentication/public"
 hui "grovehub/services/identity/internal/app/controller/authentication/public/impl"
)
func(b *builder)newControllers(s serviceGraph)(hu.AuthenticationController,hp.AuthenticationController,au.AccountController,ap.AccountController){
 return hui.New(s.auth),hpi.New(s.auth),aui.New(s.account),api.New(s.account)
}
