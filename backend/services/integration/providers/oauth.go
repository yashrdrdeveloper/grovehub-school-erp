package providers

import (
 "context"
 "strings"
)

// OAuthStrategy is an optional provider authentication capability, separate
// from provider actions and independent of any specific vendor.
type OAuthStrategy interface {
 AuthorizationURL(state string)string
 Exchange(ctx context.Context,code string)(OAuthToken,error)
 Mailbox(ctx context.Context,token OAuthToken)(string,error)
 Refresh(ctx context.Context,token OAuthToken)(OAuthToken,error)
}
type OAuthRegistry struct{strategies map[string]OAuthStrategy}
func NewOAuthRegistry()*OAuthRegistry{
 return &OAuthRegistry{strategies:make(map[string]OAuthStrategy)}
}
func(r *OAuthRegistry)Register(provider string,strategy OAuthStrategy)error{
 provider=strings.ToUpper(strings.TrimSpace(provider))
 if provider==""||strategy==nil{return ErrInvalidActionInput}
 if _,exists:=r.strategies[provider];exists{return ErrDuplicateProvider}
 r.strategies[provider]=strategy
 return nil
}
func(r *OAuthRegistry)Find(provider string)(OAuthStrategy,error){
 if r==nil{return nil,ErrProviderUnavailable}
 strategy,ok:=r.strategies[strings.ToUpper(provider)]
 if !ok{return nil,ErrProviderUnavailable}
 return strategy,nil
}
