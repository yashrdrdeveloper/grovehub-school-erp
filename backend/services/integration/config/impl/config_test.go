package impl

import (
 "testing"
)

func TestValidCallback(t *testing.T){
 cases:=[]struct{
  name string
  uri string
  provider string
  valid bool
 }{
  {"gmail","https://app.example/integrations/connections/oauth/gmail/callback","gmail",true},
  {"outlook","https://app.example/integrations/connections/oauth/outlook/callback","outlook",true},
  {"local","http://localhost:8080/integrations/connections/oauth/gmail/callback","gmail",true},
  {"old route","https://app.example/integrations/oauth/gmail/callback","gmail",false},
  {"wrong provider","https://app.example/integrations/connections/oauth/gmail/callback","outlook",false},
  {"http remote","http://app.example/integrations/connections/oauth/gmail/callback","gmail",false},
  {"query","https://app.example/integrations/connections/oauth/gmail/callback?token=secret","gmail",false},
 }
 for _,test:=range cases{
  t.Run(test.name,func(t *testing.T){
   if actual:=validCallback(test.uri,test.provider);actual!=test.valid{
    t.Fatalf("validCallback() = %v, want %v",actual,test.valid)
   }
  })
 }
}
