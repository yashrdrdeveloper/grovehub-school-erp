package impl
import(
 "grovehub/services/identity/internal/app/config"
 mapper "grovehub/services/identity/internal/app/mappers/account/impl"
 validator "grovehub/services/identity/internal/app/validators/account/impl"
 account "grovehub/services/identity/internal/app/services/account"
 accountimpl "grovehub/services/identity/internal/app/services/account/impl"
 auth "grovehub/services/identity/internal/app/services/authentication"
 authimpl "grovehub/services/identity/internal/app/services/authentication/impl"
 credentialimpl "grovehub/services/identity/internal/app/services/credential/impl"
 invitation "grovehub/services/identity/internal/app/services/registration/operations/invitation"
)
type serviceGraph struct{account account.Service;auth auth.Service;worker *invitation.Worker}
func(b *builder)newServices(r repositories,c config.Config)serviceGraph{
 credentials:=credentialimpl.New(r.credentials)
 invite:=invitation.New(r.tokens,r.mail,c)
 events:=b.newEvents()
 accounts:=accountimpl.New(r.accounts,r.tokens,r.sessions,b.foundation,credentials,invite,mapper.New(),validator.New(),r.unit,events)
 authentication:=authimpl.New(r.accounts,r.sessions,r.tokens,credentials,invite,mapper.New(),r.unit)
 var worker *invitation.Worker
 if b.integration!=nil&&len(c.EncryptionKey)==32&&c.MailTenantID!=""&&c.MailConnectionID!=""{worker=invitation.NewWorker(r.mail,b.integration,c)}
 return serviceGraph{account:accounts,auth:authentication,worker:worker}
}
