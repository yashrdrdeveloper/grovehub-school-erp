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
		mailconstants.MailProfile, mailconstants.MailList, mailconstants.MailGet,
		mailconstants.MailHistory, mailconstants.MailSend, mailconstants.MailAttachment,
	} {
		if err := registry.Register("GMAIL", action, &ActionHandler{client: client, action: action}); err != nil {
			return err
		}
	}
	return nil
}

func (h *ActionHandler) Execute(ctx context.Context, session providerdto.Session, input json.RawMessage) (any, error) {
	token := Token{AccessToken: session.AccessToken}
	switch h.action {
	case mailconstants.MailProfile:
		if len(input) != 0 && strings.TrimSpace(string(input)) != "{}" {
			return nil, providerhelpers.ErrInvalidActionInput
		}
		return h.client.Profile(ctx, token)
	case mailconstants.MailList:
		var request struct {
			PageToken string `json:"pageToken"`
		}
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		if len(request.PageToken) > 2048 {
			return nil, providerhelpers.ErrInvalidActionInput
		}
		page,err:=h.client.Messages(ctx,token,request.PageToken)
  if err!=nil{return nil,err}
  return normalizeGmailPage(page),nil
	case mailconstants.MailGet:
		var request struct {
			MessageID string `json:"messageId"`
		}
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		if len(request.MessageID) == 0 || len(request.MessageID) > 512 {
			return nil, providerhelpers.ErrInvalidActionInput
		}
		message,err:=h.client.MessageMetadata(ctx,token,request.MessageID)
  if err!=nil{return nil,err}
  return normalizeGmailMessage(message),nil
	case mailconstants.MailAttachment:
		var request struct {
			MessageID    string `json:"messageId"`
			AttachmentID string `json:"attachmentId"`
		}
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		if request.MessageID == "" || request.AttachmentID == "" || len(request.MessageID) > 512 || len(request.AttachmentID) > 512 {
			return nil, providerhelpers.ErrInvalidActionInput
		}
		return h.client.Attachment(ctx, token, request.MessageID, request.AttachmentID)
	case mailconstants.MailHistory:
		var request struct {
			StartHistoryID string `json:"startHistoryId"`
			PageToken      string `json:"pageToken"`
		}
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		if request.StartHistoryID == "" || len(request.StartHistoryID) > 128 || len(request.PageToken) > 2048 {
			return nil, providerhelpers.ErrInvalidActionInput
		}
		return h.client.History(ctx, token, request.StartHistoryID, request.PageToken)
	case mailconstants.MailSend:
		var request maildto.Send
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		if err := request.Validate(); err != nil {
			return nil, providerhelpers.ErrInvalidActionInput
		}
		return h.client.Send(ctx, token, request)
	default:
		return nil, providerhelpers.ErrUnsupportedAction
	}
}
