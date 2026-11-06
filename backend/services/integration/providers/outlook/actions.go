package outlook

import (
 "context"
 "encoding/json"
 "strings"
 "io"
 "errors"

 "grovehub/services/integration/providers"
)

type ActionHandler struct {
 client *Client
 action providers.Action
}

func RegisterActions(registry *providers.ActionRegistry,client *Client)error{
 for _,action:=range []providers.Action{
  providers.ActionMailProfile,providers.ActionMailList,
  providers.ActionMailGet,providers.ActionMailSend,providers.ActionMailDelta,
 } {
  if err:=registry.Register("OUTLOOK",action,&ActionHandler{client:client,action:action});err!=nil{return err}
 }
 return nil
}
func(h *ActionHandler)Execute(ctx context.Context,session providers.Session,input json.RawMessage)(any,error){
 switch h.action{
 case providers.ActionMailProfile:
  if len(input)>0&&strings.TrimSpace(string(input))!="{}"{return nil,providers.ErrInvalidActionInput}
  return h.client.Profile(ctx,session.AccessToken)
 case providers.ActionMailList:
  var request struct{PageToken string `json:"pageToken"`}
  if err:=decodeInput(input,&request);err!=nil{return nil,err}
  return h.client.Messages(ctx,session.AccessToken,request.PageToken)
 case providers.ActionMailGet:
  var request struct{MessageID string `json:"messageId"`}
  if err:=decodeInput(input,&request);err!=nil{return nil,err}
  return h.client.Message(ctx,session.AccessToken,request.MessageID)
 case providers.ActionMailDelta:
  var request struct{Cursor string `json:"cursor"`}
  if err:=decodeInput(input,&request);err!=nil{return nil,err}
  return h.client.Delta(ctx,session.AccessToken,request.Cursor)
 case providers.ActionMailSend:
  var request struct{Raw string `json:"raw"`}
  if err:=decodeInput(input,&request);err!=nil{return nil,err}
  return h.client.SendRaw(ctx,session.AccessToken,request.Raw)
 default:
  return nil,providers.ErrUnsupportedAction
 }
}
func decodeInput(raw json.RawMessage,out any)error{
 if len(raw)==0{raw=[]byte("{}")}
 dec:=json.NewDecoder(strings.NewReader(string(raw)))
 dec.DisallowUnknownFields()
 if err:=dec.Decode(out);err!=nil{return errors.Join(providers.ErrInvalidActionInput,err)}
 var trailing any
 if err:=dec.Decode(&trailing);err!=io.EOF{return providers.ErrInvalidActionInput}
 return nil
}
