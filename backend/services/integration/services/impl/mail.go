package impl

import (
 "context"
 "errors"
 "strings"

 "grovehub/services/integration/entity"
 "grovehub/services/integration/providers/gmail"
)

var ErrUnsupportedMailProvider=errors.New("provider does not support Gmail mail capabilities")

func(s *ConnectionService)mailToken(ctx context.Context,tenantID,id string)(gmail.Token,error){
 if s.gmail==nil{return gmail.Token{},ErrOAuthUnavailable}
 connection,err:=s.repository.Get(ctx,tenantID,id)
 if err!=nil{return gmail.Token{},err}
 if connection.Provider!="GMAIL"{return gmail.Token{},ErrUnsupportedMailProvider}
 if connection.Status!=entity.StatusActive{return gmail.Token{},ErrConnectionInactive}
 return s.validToken(ctx,id)
}
func(s *ConnectionService)MailProfile(ctx context.Context,tenantID,id string)(gmail.MailProfile,error){
 token,err:=s.mailToken(ctx,tenantID,id)
 if err!=nil{return gmail.MailProfile{},err}
 return s.gmail.Profile(ctx,token)
}
func(s *ConnectionService)MailMessages(ctx context.Context,tenantID,id,page string)(gmail.MessagePage,error){
 if len(page)>2048{return gmail.MessagePage{},ErrInvalidConnection}
 token,err:=s.mailToken(ctx,tenantID,id)
 if err!=nil{return gmail.MessagePage{},err}
 return s.gmail.Messages(ctx,token,page)
}
func(s *ConnectionService)MailMessage(ctx context.Context,tenantID,id,messageID string)(gmail.Message,error){
 if strings.TrimSpace(messageID)==""||len(messageID)>512{return gmail.Message{},ErrInvalidConnection}
 token,err:=s.mailToken(ctx,tenantID,id)
 if err!=nil{return gmail.Message{},err}
 return s.gmail.MessageMetadata(ctx,token,messageID)
}
func(s *ConnectionService)MailHistory(ctx context.Context,tenantID,id,start,page string)(gmail.HistoryPage,error){
 if start==""||len(start)>128||len(page)>2048{return gmail.HistoryPage{},ErrInvalidConnection}
 token,err:=s.mailToken(ctx,tenantID,id)
 if err!=nil{return gmail.HistoryPage{},err}
 return s.gmail.History(ctx,token,start,page)
}
func(s *ConnectionService)MailSend(ctx context.Context,tenantID,id,raw string)(gmail.SentMessage,error){
 token,err:=s.mailToken(ctx,tenantID,id)
 if err!=nil{return gmail.SentMessage{},err}
 return s.gmail.SendRaw(ctx,token,raw)
}
