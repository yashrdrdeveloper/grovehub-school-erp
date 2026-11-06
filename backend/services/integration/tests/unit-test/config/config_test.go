package config_test

import (
 "testing"

 integrationconfig "grovehub/services/integration/config/impl"
)

func TestOAuthCallbackValidation(t *testing.T){
 tests:=[]struct{provider,envName,uri string;valid bool}{
  {"gmail","GROVEHUB_GMAIL_REDIRECT_URI","https://example.test/integrations/connections/oauth/gmail/callback",true},
  {"gmail","GROVEHUB_GMAIL_REDIRECT_URI","https://example.test/integrations/oauth/gmail/callback",false},
  {"outlook","GROVEHUB_OUTLOOK_REDIRECT_URI","https://example.test/integrations/connections/oauth/outlook/callback",true},
  {"outlook","GROVEHUB_OUTLOOK_REDIRECT_URI","http://example.test/integrations/connections/oauth/outlook/callback",false},
 }
 for _,tc:=range tests{
  t.Run(tc.provider+"/"+tc.uri,func(t *testing.T){
   t.Setenv(tc.envName,tc.uri)
   var err error
   if tc.provider=="gmail"{
    t.Setenv("GROVEHUB_GMAIL_CLIENT_ID","id")
    t.Setenv("GROVEHUB_GMAIL_CLIENT_SECRET","secret")
    _,err=integrationconfig.NewEnvironmentLoader().Gmail()
   }else{
    t.Setenv("GROVEHUB_OUTLOOK_CLIENT_ID","id")
    t.Setenv("GROVEHUB_OUTLOOK_CLIENT_SECRET","secret")
    _,err=integrationconfig.NewEnvironmentLoader().Outlook()
   }
   if (err==nil)!=tc.valid{
    t.Fatalf("callback valid=%t, expected %t, error=%v",err==nil,tc.valid,err)
   }
  })
 }
}
