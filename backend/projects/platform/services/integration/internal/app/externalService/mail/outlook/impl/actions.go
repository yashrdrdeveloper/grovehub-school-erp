package impl

import (
    providerdto "grovehub/services/integration/internal/app/providers/dto"
    providercontracts "grovehub/services/integration/internal/app/providers/contracts"
    mailconstants "grovehub/services/integration/internal/app/providers/mailprovider/constants"
	"context"
	"encoding/json"
	"strings"

 maildto "grovehub/services/integration/internal/app/providers/mailprovider/dto"
	providerhelpers "grovehub/services/integration/internal/app/providers/helpers"
)

type ActionHandler struct {
	client *Client
	action providerdto.Action
}

func RegisterActions(registry providercontracts.ActionRegistry, client *Client) error {
	for _, action := range []providerdto.Action{
		mailconstants.MailProfile, mailconstants.MailList,
		mailconstants.MailGet, mailconstants.MailSend, mailconstants.MailDelta, mailconstants.MailAttachment,
	} {
		if err := registry.Register("OUTLOOK", action, &ActionHandler{client: client, action: action}); err != nil {
			return err
		}
	}
	return nil
}

func (h *ActionHandler) Execute(ctx context.Context, session providerdto.Session, input json.RawMessage) (any, error) {
	switch h.action {
	case mailconstants.MailProfile:
		if len(input) > 0 && strings.TrimSpace(string(input)) != "{}" {
			return nil, providerhelpers.ErrInvalidActionInput
		}
		return h.client.Profile(ctx, session.AccessToken)
	case mailconstants.MailList:
		var request struct {
			PageToken string `json:"pageToken"`
		}
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		page,err:=h.client.Messages(ctx,session.AccessToken,request.PageToken)
  if err!=nil{return nil,err}
  return normalizeOutlookPage(page),nil
	case mailconstants.MailGet:
		var request struct {
			MessageID string `json:"messageId"`
		}
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		message,err:=h.client.Message(ctx,session.AccessToken,request.MessageID)
  if err!=nil{return nil,err}
  return normalizeOutlookMessage(message),nil
	case mailconstants.MailAttachment:
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
  if json.Unmarshal(raw,&file)!=nil{return nil,providerhelpers.ErrProviderUnavailable}
  if file.ContentBytes==""{return nil,providerhelpers.ErrUnsupportedAction}
  return maildto.AttachmentData{
   AttachmentID:file.ID,Filename:file.Name,ContentType:file.ContentType,
   Size:file.Size,ContentBase64:file.ContentBytes,
  },nil
	case mailconstants.MailDelta:
		var request struct {
			Cursor string `json:"cursor"`
		}
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		return h.client.Delta(ctx, session.AccessToken, request.Cursor)
	case mailconstants.MailSend:
		var request maildto.Send
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		if err := request.Validate(); err != nil {
			return nil, providerhelpers.ErrInvalidActionInput
		}
		return h.client.Send(ctx, session.AccessToken, request)
	default:
		return nil, providerhelpers.ErrUnsupportedAction
	}
}
