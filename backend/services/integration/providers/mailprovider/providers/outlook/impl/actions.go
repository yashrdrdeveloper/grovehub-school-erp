package impl

import (
	"context"
	"encoding/json"
	"strings"

	"grovehub/services/integration/providers"
 "grovehub/services/integration/providers/mailprovider/mailcontract"
	providerhelpers "grovehub/services/integration/providers/helpers"
)

type ActionHandler struct {
	client *Client
	action providers.Action
}

func RegisterActions(registry providers.ActionRegistry, client *Client) error {
	for _, action := range []providers.Action{
		providers.ActionMailProfile, providers.ActionMailList,
		providers.ActionMailGet, providers.ActionMailSend, providers.ActionMailDelta, providers.ActionMailAttachment,
	} {
		if err := registry.Register("OUTLOOK", action, &ActionHandler{client: client, action: action}); err != nil {
			return err
		}
	}
	return nil
}
func (h *ActionHandler) Execute(ctx context.Context, session providers.Session, input json.RawMessage) (any, error) {
	switch h.action {
	case providers.ActionMailProfile:
		if len(input) > 0 && strings.TrimSpace(string(input)) != "{}" {
			return nil, providers.ErrInvalidActionInput
		}
		return h.client.Profile(ctx, session.AccessToken)
	case providers.ActionMailList:
		var request struct {
			PageToken string `json:"pageToken"`
		}
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		page,err:=h.client.Messages(ctx,session.AccessToken,request.PageToken)
  if err!=nil{return nil,err}
  return normalizeOutlookPage(page),nil
	case providers.ActionMailGet:
		var request struct {
			MessageID string `json:"messageId"`
		}
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		message,err:=h.client.Message(ctx,session.AccessToken,request.MessageID)
  if err!=nil{return nil,err}
  return normalizeOutlookMessage(message),nil
	case providers.ActionMailAttachment:
		var request struct {
			MessageID    string `json:"messageId"`
			AttachmentID string `json:"attachmentId"`
		}
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		raw,err:=h.client.Attachment(ctx,session.AccessToken,request.MessageID,request.AttachmentID)
  if err!=nil{return nil,err}
  var file struct{
   ID string `json:"id"`
   Name string `json:"name"`
   ContentType string `json:"contentType"`
   Size int64 `json:"size"`
   ContentBytes string `json:"contentBytes"`
  }
  if json.Unmarshal(raw,&file)!=nil{return nil,providers.ErrProviderUnavailable}
  if file.ContentBytes==""{return nil,providers.ErrUnsupportedAction}
  return mailcontract.AttachmentData{
   AttachmentID:file.ID,Filename:file.Name,ContentType:file.ContentType,
   Size:file.Size,ContentBase64:file.ContentBytes,
  },nil
	case providers.ActionMailDelta:
		var request struct {
			Cursor string `json:"cursor"`
		}
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		return h.client.Delta(ctx, session.AccessToken, request.Cursor)
	case providers.ActionMailSend:
		var request mailcontract.Send
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		if err := request.Validate(); err != nil {
			return nil, providers.ErrInvalidActionInput
		}
		return h.client.Send(ctx, session.AccessToken, request)
	default:
		return nil, providers.ErrUnsupportedAction
	}
}
