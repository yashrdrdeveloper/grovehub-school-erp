package impl

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"grovehub/services/integration/providers"
 "grovehub/services/integration/providers/mailprovider/mailcontract"
)

type Message struct {
	ID                string          `json:"id"`
	Subject           string          `json:"subject"`
	ConversationID    string          `json:"conversationId"`
	InternetMessageID string          `json:"internetMessageId"`
	ReceivedDateTime  string          `json:"receivedDateTime"`
	SentDateTime      string          `json:"sentDateTime"`
	From              json.RawMessage `json:"from,omitempty"`
	ToRecipients      json.RawMessage `json:"toRecipients,omitempty"`
	Body              json.RawMessage `json:"body,omitempty"`
	IsRead            bool            `json:"isRead"`
}
type MessagePage struct {
	Messages []Message `json:"value"`
	NextLink string    `json:"@odata.nextLink,omitempty"`
}
type SendResult struct {
	Accepted bool `json:"accepted"`
}

func (c *Client) Profile(ctx context.Context, accessToken string) (any, error) {
	var result struct {
		Mail              string `json:"mail"`
		UserPrincipalName string `json:"userPrincipalName"`
	}
	err := c.graph(ctx, accessToken, http.MethodGet, "/me?$select=mail,userPrincipalName", nil, "", &result)
	if err != nil {
		return nil, err
	}
	if result.Mail == "" {
		result.Mail = result.UserPrincipalName
	}
	return map[string]string{"emailAddress": strings.ToLower(result.Mail)}, nil
}
func safeNextLink(raw string) (string, error) {
	if raw == "" {
		return "/me/messages?$top=50&$select=id,subject,conversationId,internetMessageId,receivedDateTime,sentDateTime,from,toRecipients,isRead", nil
	}
	if len(raw) > 4096 {
		return "", providers.ErrInvalidActionInput
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", providers.ErrInvalidActionInput
	}
	if parsed.Scheme != "https" || parsed.Host != "graph.microsoft.com" || parsed.User != nil || parsed.Fragment != "" ||
		parsed.Path != "/v1.0/me/messages" {
		return "", providers.ErrInvalidActionInput
	}
	return "/me/messages?" + parsed.RawQuery, nil
}
func (c *Client) Messages(ctx context.Context, accessToken, page string) (MessagePage, error) {
	var result MessagePage
	path, err := safeNextLink(page)
	if err != nil {
		return result, err
	}
	err = c.graph(ctx, accessToken, http.MethodGet, path, nil, "", &result)
	return result, err
}
func (c *Client) Message(ctx context.Context, accessToken, id string) (Message, error) {
	var result Message
	if id == "" || len(id) > 512 || strings.ContainsAny(id, "?#") {
		return result, providers.ErrInvalidActionInput
	}
	err := c.graph(ctx, accessToken, http.MethodGet, "/me/messages/"+url.PathEscape(id), nil, "", &result)
	return result, err
}
// Send converts provider-neutral mail into the Graph MIME sendMail request.
func (c *Client) Send(ctx context.Context, accessToken string, request mailcontract.Send) (SendResult, error) {
 identity, err := c.Profile(ctx, accessToken)
 if err != nil { return SendResult{}, err }
 me, ok := identity.(map[string]string)
 if !ok || me["emailAddress"] == "" { return SendResult{}, providers.ErrProviderUnavailable }
 raw, err := mailcontract.BuildMIME(me["emailAddress"], request)
 if err != nil { return SendResult{}, providers.ErrInvalidActionInput }
 return c.SendRaw(ctx, accessToken, base64.RawURLEncoding.EncodeToString(raw))
}

func (c *Client) SendRaw(ctx context.Context, accessToken, raw string) (SendResult, error) {
	var result SendResult
	if len(raw) == 0 || len(raw) > 8*1024*1024 {
		return result, providers.ErrInvalidActionInput
	}
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return result, providers.ErrInvalidActionInput
	}
	if len(decoded) == 0 || len(decoded) > 6*1024*1024 {
		return result, providers.ErrInvalidActionInput
	}
	// Microsoft Graph's MIME form expects base64 (not Gmail's base64url)
	// supplied as text/plain; response is HTTP 202 with no message identifier.
	encoded := base64.StdEncoding.EncodeToString(decoded)
	err = c.graph(ctx, accessToken, http.MethodPost, "/me/sendMail", strings.NewReader(encoded), "text/plain", nil)
	if err != nil {
		return result, err
	}
	result.Accepted = true
	return result, nil
}

type DeltaPage struct {
	Messages  []Message `json:"value"`
	NextLink  string    `json:"@odata.nextLink,omitempty"`
	DeltaLink string    `json:"@odata.deltaLink,omitempty"`
}

// Delta reads changes without persisting the cursor.
// Communication must commit messages before advancing its delta cursor.
func (c *Client) Delta(ctx context.Context, accessToken, cursor string) (DeltaPage, error) {
	var result DeltaPage
	path := "/me/mailFolders/inbox/messages/delta?$select=id,subject,conversationId,internetMessageId,receivedDateTime,sentDateTime,from,toRecipients,isRead"
	if cursor != "" {
		if len(cursor) > 4096 {
			return result, providers.ErrInvalidActionInput
		}
		parsed, err := url.Parse(cursor)
		if err != nil || parsed.Scheme != "https" || parsed.Host != "graph.microsoft.com" || parsed.User != nil ||
			parsed.Fragment != "" || !allowedDeltaPath(parsed.Path) {
			return result, providers.ErrInvalidActionInput
		}
		path = "/me/mailFolders/inbox/messages/delta?" + parsed.RawQuery
	}
	err := c.graph(ctx, accessToken, http.MethodGet, path, nil, "", &result)
	return result, err
}

// Allow only a signed-in user's mailbox-folder delta path. Microsoft Graph
// can substitute the folder ID for the "inbox" well-known name in nextLink.
var deltaResourcePath = regexp.MustCompile(`(?i)^/v1[.]0/me/mailfolders(?:/[^/]{1,512}|[(]'[^/]{1,512}'[)])/messages/delta$`)

func allowedDeltaPath(path string) bool {
	return deltaResourcePath.MatchString(path)
}

// Attachment reads one Graph message attachment. Data and metadata remain
// transient; Communication will own any attachment persistence.
func (c *Client) Attachment(ctx context.Context, accessToken, messageID, attachmentID string) (json.RawMessage, error) {
	if messageID == "" || attachmentID == "" || len(messageID) > 512 || len(attachmentID) > 512 ||
		strings.ContainsAny(messageID, "?#") || strings.ContainsAny(attachmentID, "?#") {
		return nil, providers.ErrInvalidActionInput
	}
	var attachment json.RawMessage
	err := c.graph(ctx, accessToken, http.MethodGet,
		"/me/messages/"+url.PathEscape(messageID)+"/attachments/"+url.PathEscape(attachmentID),
		nil, "", &attachment)
	return attachment, err
}
