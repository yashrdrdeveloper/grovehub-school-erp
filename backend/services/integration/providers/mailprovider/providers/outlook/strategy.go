package outlook

import (
 "grovehub/services/integration/providers"
 "grovehub/services/integration/providers/mailprovider"
)

type strategy struct{client *Client}

func NewStrategy(client *Client)mailprovider.Strategy{return &strategy{client:client}}
func(s *strategy)Descriptor()providers.Descriptor{return New().Descriptor()}
func(s *strategy)Register(actions providers.ActionRegistry,oauth providers.OAuthRegistry)error{
 if s.client==nil{return nil}
 if err:=RegisterActions(actions,s.client);err!=nil{return err}
 return oauth.Register("OUTLOOK",s.client)
}
