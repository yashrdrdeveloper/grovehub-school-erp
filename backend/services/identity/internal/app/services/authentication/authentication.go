package authentication
import ("context";"grovehub/services/identity/internal/app/dto/authentication";"grovehub/services/identity/internal/app/dto/account")
type Cookies struct{Access,Refresh string}
type SessionResult struct{Data authentication.LoginResult;Cookies Cookies}
type RefreshResult struct{Data authentication.RefreshResult;Cookies Cookies}
type Grant struct{Secret string;Data authentication.GrantResult}
type Service interface{
 Login(context.Context,authentication.Login,string)(SessionResult,error)
 Refresh(context.Context,string)(RefreshResult,error)
 Forgot(context.Context,string)error
 Exchange(context.Context,string,string)(Grant,error)
 ResetPassword(context.Context,string,string)error
 Reauthenticate(context.Context,string)error
 ChangePassword(context.Context,string)error
 Logout(context.Context)error
 LogoutAll(context.Context)(int64,error)
 ListSessions(context.Context)([]map[string]any,error)
 RevokeSession(context.Context,string)error
 Me(context.Context)(account.Account,error)
}
