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
		providers.ActionMailProfile, providers.ActionMailList, providers.ActionMailGet,
		providers.ActionMailHistory, providers.ActionMailSend, providers.ActionMailAttachment,
	} {
		if err := registry.Register("GMAIL", action, &ActionHandler{client: client, action: action}); err != nil {
			return err
		}
	}
	return nil
}

func (h *ActionHandler) Execute(ctx context.Context, session providers.Session, input json.RawMessage) (any, error) {
	token := Token{AccessToken: session.AccessToken}
	switch h.action {
	case providers.ActionMailProfile:
		if len(input) != 0 && strings.TrimSpace(string(input)) != "{}" {
			return nil, providers.ErrInvalidActionInput
		}
		return h.client.Profile(ctx, token)
	case providers.ActionMailList:
		var request struct {
			PageToken string `json:"pageToken"`
		}
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		if len(request.PageToken) > 2048 {
			return nil, providers.ErrInvalidActionInput
		}
		page,err:=h.client.Messages(ctx,token,request.PageToken)
  if err!=nil{return nil,err}
  return normalizeGmailPage(page),nil
	case providers.ActionMailGet:
		var request struct {
			MessageID string `json:"messageId"`
		}
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		if len(request.MessageID) == 0 || len(request.MessageID) > 512 {
			return nil, providers.ErrInvalidActionInput
		}
		message,err:=h.client.MessageMetadata(ctx,token,request.MessageID)
  if err!=nil{return nil,err}
  return normalizeGmailMessage(message),nil
	case providers.ActionMailAttachment:
		var request struct {
			MessageID    string `json:"messageId"`
			AttachmentID string `json:"attachmentId"`
		}
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		if request.MessageID == "" || request.AttachmentID == "" || len(request.MessageID) > 512 || len(request.AttachmentID) > 512 {
			return nil, providers.ErrInvalidActionInput
		}
		return h.client.Attachment(ctx, token, request.MessageID, request.AttachmentID)
	case providers.ActionMailHistory:
		var request struct {
			StartHistoryID string `json:"startHistoryId"`
			PageToken      string `json:"pageToken"`
		}
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		if request.StartHistoryID == "" || len(request.StartHistoryID) > 128 || len(request.PageToken) > 2048 {
			return nil, providers.ErrInvalidActionInput
		}
		return h.client.History(ctx, token, request.StartHistoryID, request.PageToken)
	case providers.ActionMailSend:
		var request mailcontract.Send
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		if err := request.Validate(); err != nil {
			return nil, providers.ErrInvalidActionInput
		}
		return h.client.Send(ctx, token, request)
	default:
		return nil, providers.ErrUnsupportedAction
	}
}
