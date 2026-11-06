package impl

import (
	"context"
 "bytes"
 "net/mail"
 "encoding/base64"
	"encoding/json"
	"strings"

	"grovehub/services/integration/providers"
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
		return h.client.Messages(ctx, token, request.PageToken)
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
		return h.client.MessageMetadata(ctx, token, request.MessageID)
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
		var request struct {
			Raw string `json:"raw"`
		}
		if err := providerhelpers.DecodeStrictActionInput(input, &request); err != nil {
			return nil, err
		}
		if request.Raw == "" || len(request.Raw) > 8*1024*1024 {
   return nil, providers.ErrInvalidActionInput
  }
  decoded,err:=base64.RawURLEncoding.DecodeString(request.Raw)
  if err!=nil || len(decoded)==0 {
   return nil,providers.ErrInvalidActionInput
  }
  // Go's RFC 5322 reader accepts both CRLF and LF line endings.
  // Require actual message headers rather than substrings in the body.
  parsed,parseErr:=mail.ReadMessage(bytes.NewReader(decoded))
  if parseErr!=nil || parsed.Header.Get("To")=="" || parsed.Header.Get("From")=="" {
   return nil,providers.ErrInvalidActionInput
  }
  return h.client.SendRaw(ctx, token, request.Raw)
	default:
		return nil, providers.ErrUnsupportedAction
	}
}
