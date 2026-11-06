package impl

import (
 "strings"
 "grovehub/services/integration/providers"
)
type oauthRegistry struct{strategies map[string]providers.OAuthStrategy}
func NewOAuthRegistry()providers.OAuthRegistry{return &oauthRegistry{strategies:make(map[string]providers.OAuthStrategy)}}
func(r *oauthRegistry)Register(provider string,strategy providers.OAuthStrategy)error{
 provider=strings.ToUpper(strings.TrimSpace(provider))
 if provider==""||strategy==nil{return providers.ErrInvalidActionInput}
 if _,exists:=r.strategies[provider];exists{return providers.ErrDuplicateProvider}
 r.strategies[provider]=strategy
 return nil
}
func(r *oauthRegistry)Find(provider string)(providers.OAuthStrategy,error){
 strategy,ok:=r.strategies[strings.ToUpper(strings.TrimSpace(provider))]
 if !ok{return nil,providers.ErrProviderUnavailable}
 return strategy,nil
}
