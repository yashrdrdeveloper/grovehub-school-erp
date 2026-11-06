package impl

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"grovehub/services/integration/providers"
	"net/http"
	"net/url"
)

type MailProfile struct {
	EmailAddress string `json:"emailAddress"`
	HistoryID    string `json:"historyId"`
}
type SentMessage struct {
	ID       string `json:"id"`
	ThreadID string `json:"threadId"`
}

type Message struct {
	ID           string          `json:"id"`
	ThreadID     string          `json:"threadId"`
	HistoryID    string          `json:"historyId"`
	LabelIDs     []string        `json:"labelIds"`
	Snippet      string          `json:"snippet,omitempty"`
	SizeEstimate int64           `json:"sizeEstimate,omitempty"`
	Payload      json.RawMessage `json:"payload,omitempty"`
}
type MessagePage struct {
	Messages      []Message `json:"messages"`
	NextPageToken string    `json:"nextPageToken"`
}
type HistoryChange struct {
	MessagesAdded []struct {
		Message Message `json:"message"`
	} `json:"messagesAdded"`
	MessagesDeleted []struct {
		Message Message `json:"message"`
	} `json:"messagesDeleted"`
	LabelsAdded []struct {
		Message Message `json:"message"`
	} `json:"labelsAdded"`
	LabelsRemoved []struct {
		Message Message `json:"message"`
	} `json:"labelsRemoved"`
}
type HistoryPage struct {
	History       []HistoryChange `json:"history"`
	NextPageToken string          `json:"nextPageToken"`
	HistoryID     string          `json:"historyId"`
}

var ErrExpiredHistory = errors.New("gmail history cursor expired")

func (c *Client) Profile(ctx context.Context, t Token) (MailProfile, error) {
	var profile MailProfile
	if err := c.do(ctx, t, http.MethodGet, "/profile", nil, &profile); err != nil {
		return MailProfile{}, err
	}
	if profile.EmailAddress == "" || profile.HistoryID == "" {
		return MailProfile{}, errors.New("gmail profile missing mailbox or history cursor")
	}
	return profile, nil
}

func (c *Client) HistoryID(ctx context.Context, t Token) (string, error) {
	profile, err := c.Profile(ctx, t)
	return profile.HistoryID, err
}
func (c *Client) Messages(ctx context.Context, t Token, page string) (MessagePage, error) {
	var p MessagePage
	query := url.Values{"maxResults": {"100"}}
	if page != "" {
		query.Set("pageToken", page)
	}
	err := c.do(ctx, t, http.MethodGet, "/messages?"+query.Encode(), nil, &p)
	return p, err
}
func (c *Client) MessageMetadata(ctx context.Context, t Token, id string) (Message, error) {
	var m Message
	// Full format is an ephemeral provider response, not an Integration record.
	err := c.do(ctx, t, http.MethodGet, "/messages/"+url.PathEscape(id)+"?format=full", nil, &m)
	return m, err
}
func (c *Client) History(ctx context.Context, t Token, startID, page string) (HistoryPage, error) {
	var p HistoryPage
	if startID == "" {
		return p, errors.New("missing Gmail startHistoryId")
	}
	query := url.Values{"startHistoryId": {startID}, "maxResults": {"100"}}
	if page != "" {
		query.Set("pageToken", page)
	}
	path := "/history?" + query.Encode()
	err := c.do(ctx, t, http.MethodGet, path, nil, &p)
	if err != nil {
		if providers.HTTPStatus(err) == http.StatusNotFound {
			return p, ErrExpiredHistory
		}
		return p, err
	}
	return p, nil
}

// SendRaw submits an RFC 5322 message encoded as base64url.
// The caller (eventually Communication) owns its message and delivery state.
func (c *Client) SendRaw(ctx context.Context, t Token, raw string) (SentMessage, error) {
	var sent SentMessage
	if len(raw) == 0 || len(raw) > 8*1024*1024 {
		return sent, errors.New("invalid raw email size")
	}
	if _, err := base64.RawURLEncoding.DecodeString(raw); err != nil {
		return sent, errors.New("invalid base64url email")
	}
	err := c.do(ctx, t, http.MethodPost, "/messages/send", map[string]string{"raw": raw}, &sent)
	return sent, err
}

// Attachment returns Gmail's base64url-encoded attachment bytes without
// creating any Integration-owned attachment record.
func (c *Client) Attachment(ctx context.Context, t Token, messageID, attachmentID string) (any, error) {
	if messageID == "" || attachmentID == "" || len(messageID) > 512 || len(attachmentID) > 512 {
		return nil, errors.New("invalid Gmail attachment reference")
	}
	var result struct {
		AttachmentID string `json:"attachmentId,omitempty"`
		Size         int64  `json:"size"`
		Data         string `json:"data"`
	}
	err := c.do(ctx, t, http.MethodGet,
		"/messages/"+url.PathEscape(messageID)+"/attachments/"+url.PathEscape(attachmentID),
		nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
