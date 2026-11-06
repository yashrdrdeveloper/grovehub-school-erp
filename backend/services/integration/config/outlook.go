package config

import (
 "errors"
 "os"
 "strings"
)

type Outlook struct {
 ClientID string
 ClientSecret string
 RedirectURI string
 Tenant string
}
func OutlookFromEnvironment()(Outlook,error){
 c:=Outlook{
  ClientID:strings.TrimSpace(os.Getenv("GROVEHUB_OUTLOOK_CLIENT_ID")),
  ClientSecret:strings.TrimSpace(os.Getenv("GROVEHUB_OUTLOOK_CLIENT_SECRET")),
  RedirectURI:strings.TrimSpace(os.Getenv("GROVEHUB_OUTLOOK_REDIRECT_URI")),
  Tenant:strings.TrimSpace(os.Getenv("GROVEHUB_OUTLOOK_TENANT")),
 }
 if c.Tenant==""{c.Tenant="common"}
 if c.ClientID==""||c.ClientSecret==""||c.RedirectURI==""{
  return Outlook{},errors.New("Outlook OAuth environment configuration is incomplete")
 }
 if !strings.HasPrefix(c.RedirectURI,"https://")&&!strings.HasPrefix(c.RedirectURI,"http://localhost:"){
  return Outlook{},errors.New("Outlook OAuth redirect URI must use HTTPS")
 }
 if strings.ContainsAny(c.Tenant,"/?#\\ \t\r\n"){
  return Outlook{},errors.New("invalid Outlook OAuth tenant")
 }
 return c,nil
}
