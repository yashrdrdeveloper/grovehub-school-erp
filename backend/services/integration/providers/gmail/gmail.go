package gmail

import (
 "io"

 integrationconfig "grovehub/services/integration/config"
 "grovehub/services/integration/providers"
 impl "grovehub/services/integration/providers/gmail/impl"
)

// Public Gmail adapter contracts and factories. Concrete provider
// operations, OAuth and transport behavior live in gmail/impl.
type Client = impl.Client
type Token = impl.Token
type MailProfile = impl.MailProfile
type SentMessage = impl.SentMessage
type Message = impl.Message
type MessagePage = impl.MessagePage
type HistoryChange = impl.HistoryChange
type HistoryPage = impl.HistoryPage
type PushNotification = impl.PushNotification
type PushVerifier = impl.PushVerifier
type Provider = impl.Provider
type StatusError = impl.StatusError

var ErrExpiredHistory = impl.ErrExpiredHistory

func New()providers.Provider{return impl.New()}
func NewClient(cfg integrationconfig.Gmail)*Client{return impl.NewClient(cfg)}
func RegisterActions(registry providers.ActionRegistry,client *Client)error{
 return impl.RegisterActions(registry,client)
}
func NewPushVerifier(audience,email string)*PushVerifier{
 return impl.NewPushVerifier(audience,email)
}
func DecodePush(body io.Reader)(PushNotification,error){return impl.DecodePush(body)}
