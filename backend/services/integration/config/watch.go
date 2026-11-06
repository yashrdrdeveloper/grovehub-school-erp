package config
import ("os";"strings")
type GmailPush struct{Topic,Audience,ServiceAccount string}
func GmailPushFromEnvironment()GmailPush{
 return GmailPush{
 Topic:strings.TrimSpace(os.Getenv("GROVEHUB_GMAIL_PUBSUB_TOPIC")),
 Audience:strings.TrimSpace(os.Getenv("GROVEHUB_GMAIL_PUBSUB_AUDIENCE")),
 ServiceAccount:strings.TrimSpace(os.Getenv("GROVEHUB_GMAIL_PUBSUB_SERVICE_ACCOUNT")),
 }
}
func(c GmailPush)Enabled()bool{return c.Topic!=""&&c.Audience!=""&&c.ServiceAccount!=""}
