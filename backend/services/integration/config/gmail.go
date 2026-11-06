package config

import (
 "encoding/base64"
 "errors"
 "os"
 "strings"
)

type Gmail struct {
 ClientID string
 ClientSecret string
 RedirectURI string
 EncryptionKey []byte
}

func GmailFromEnvironment()(Gmail,error){
 c:=Gmail{
 ClientID:strings.TrimSpace(os.Getenv("GROVEHUB_GMAIL_CLIENT_ID")),
 ClientSecret:strings.TrimSpace(os.Getenv("GROVEHUB_GMAIL_CLIENT_SECRET")),
 RedirectURI:strings.TrimSpace(os.Getenv("GROVEHUB_GMAIL_REDIRECT_URI")),
 }
 if c.ClientID==""||c.ClientSecret==""||c.RedirectURI=="" {return c,errors.New("Gmail OAuth env credentials and redirect URI required")}
 key,err:=base64.StdEncoding.DecodeString(os.Getenv("GROVEHUB_INTEGRATION_ENCRYPTION_KEY"))
 if err!=nil||len(key)!=32{return c,errors.New("GROVEHUB_INTEGRATION_ENCRYPTION_KEY must be base64-encoded 32 bytes")}
 c.EncryptionKey=key
 if !strings.HasPrefix(c.RedirectURI,"https://")&&!strings.HasPrefix(c.RedirectURI,"http://localhost:"){
  return c,errors.New("Gmail OAuth redirect URI must use HTTPS")
 }
 return c,nil
}
