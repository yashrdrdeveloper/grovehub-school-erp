package gmail

import (
 "context"
 "encoding/json"
 "errors"
 "strings"

 "grovehub/services/integration/providers"
)

type ActionHandler struct {
 client *Client
 action providers.Action
}

func RegisterActions(registry *providers.ActionRegistry,client *Client)error{
 for _,action:=range []providers.Action{
  providers.ActionMailProfile,providers.ActionMailList,providers.ActionMailGet,
  providers.ActionMailHistory,providers.ActionMailSend,
 } {
  if err:=registry.Register("GMAIL",action,&ActionHandler{client:client,action:action});err!=nil{return err}
 }
 return nil
}

func(h *ActionHandler)Execute(ctx context.Context,session providers.Session,input json.RawMessage)(any,error){
 token:=Token{AccessToken:session.AccessToken}
 switch h.action{
 case providers.ActionMailProfile:
  if len(input)!=0 && strings.TrimSpace(string(input))!="{}" {return nil,providers.ErrInvalidActionInput}
  return h.client.Profile(ctx,token)
 case providers.ActionMailList:
  var request struct{PageToken string `json:"pageToken"`}
  if err:=decodeActionInput(input,&request);err!=nil{return nil,err}
  if len(request.PageToken)>2048{return nil,providers.ErrInvalidActionInput}
  return h.client.Messages(ctx,token,request.PageToken)
 case providers.ActionMailGet:
  var request struct{MessageID string `json:"messageId"`}
  if err:=decodeActionInput(input,&request);err!=nil{return nil,err}
  if len(request.MessageID)==0||len(request.MessageID)>512{return nil,providers.ErrInvalidActionInput}
  return h.client.MessageMetadata(ctx,token,request.MessageID)
 case providers.ActionMailHistory:
  var request struct{
   StartHistoryID string `json:"startHistoryId"`
   PageToken string `json:"pageToken"`
  }
  if err:=decodeActionInput(input,&request);err!=nil{return nil,err}
  if request.StartHistoryID==""||len(request.StartHistoryID)>128||len(request.PageToken)>2048{
   return nil,providers.ErrInvalidActionInput
  }
  return h.client.History(ctx,token,request.StartHistoryID,request.PageToken)
 case providers.ActionMailSend:
  var request struct{Raw string `json:"raw"`}
  if err:=decodeActionInput(input,&request);err!=nil{return nil,err}
  if request.Raw==""{return nil,providers.ErrInvalidActionInput}
  return h.client.SendRaw(ctx,token,request.Raw)
 default:
  return nil,providers.ErrUnsupportedAction
 }
}

func decodeActionInput(input json.RawMessage,target any)error{
 if len(input)==0 {input=[]byte("{}")}
 decoder:=json.NewDecoder(strings.NewReader(string(input)))
 decoder.DisallowUnknownFields()
 if err:=decoder.Decode(target);err!=nil{return errors.Join(providers.ErrInvalidActionInput,err)}
 var trailing any
 if decoder.Decode(&trailing)==nil{return providers.ErrInvalidActionInput}
 return nil
}
