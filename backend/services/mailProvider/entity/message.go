package entity

import (
	"time"

	"github.com/lib/pq"
)

type GID string
type MessageId string
type MessageStatus string

const (
	MessageStatusQueued    MessageStatus = "QUEUED"
	MessageStatusSubmitted MessageStatus = "SUBMITTED"
	MessageStatusFailed    MessageStatus = "FAILED"
)

type Message struct {
	MessageId           MessageId      `gorm:"column:message_id;type:uuid;primaryKey"`
	IdempotencyKey      string         `gorm:"column:idempotency_key;not null"`
	RequestHash         string         `gorm:"column:request_hash;not null"`
	GID                 GID            `gorm:"column:gid;type:uuid;not null"`
	SenderId            SenderId       `gorm:"column:sender_id;type:uuid;not null"`
	ThreadId            MessageId      `gorm:"column:thread_id;type:uuid;not null"`
	ParentMessageId     *MessageId     `gorm:"column:parent_message_id;type:uuid"`
	ResendOfMessageId   *MessageId     `gorm:"column:resend_of_message_id;type:uuid"`
	InternetMessageId   string         `gorm:"column:internet_message_id;not null"`
	InReplyTo           *string        `gorm:"column:in_reply_to"`
	ReferenceMessageIds pq.StringArray `gorm:"column:reference_message_ids;type:text[];not null"`
	To                  pq.StringArray `gorm:"column:to_addresses;type:text[];not null"`
	Cc                  pq.StringArray `gorm:"column:cc_addresses;type:text[];not null"`
	Bcc                 pq.StringArray `gorm:"column:bcc_addresses;type:text[];not null"`
	Subject             string         `gorm:"column:subject;not null"`
	TextBody            string         `gorm:"column:text_body;not null"`
	HTMLBody            string         `gorm:"column:html_body;not null"`
	Status              MessageStatus  `gorm:"column:status;not null"`
	FailureMessage      *string        `gorm:"column:failure_message"`
	SubmissionStartedAt *time.Time     `gorm:"column:submission_started_at"`
	Version             int64          `gorm:"column:version;not null"`
	CreatedAt           time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt           time.Time      `gorm:"column:updated_at;not null"`
}

func (Message) TableName() string { return "mail_provider.messages" }

func NewMessage(
	messageId MessageId,
	idempotencyKey string,
	requestHash string,
	gid GID,
	senderId SenderId,
	to []string,
	cc []string,
	bcc []string,
	subject string,
	textBody string,
	htmlBody string,
	status MessageStatus,
	failureMessage string,
	version int64,
	createdAt time.Time,
	updatedAt time.Time,
) Message {
	return Message{
		MessageId:           messageId,
		IdempotencyKey:      idempotencyKey,
		RequestHash:         requestHash,
		GID:                 gid,
		SenderId:            senderId,
		ThreadId:            messageId,
		InternetMessageId:   string(messageId) + "@grovehub.local",
		ReferenceMessageIds: pq.StringArray{},
		To:                  pq.StringArray(copyStrings(to)),
		Cc:                  pq.StringArray(copyStrings(cc)),
		Bcc:                 pq.StringArray(copyStrings(bcc)),
		Subject:             subject,
		TextBody:            textBody,
		HTMLBody:            htmlBody,
		Status:              status,
		FailureMessage:      optionalString(failureMessage),
		Version:             version,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
	}
}

func (message Message) GetMessageId() MessageId   { return message.MessageId }
func (message Message) GetIdempotencyKey() string { return message.IdempotencyKey }
func (message Message) GetRequestHash() string    { return message.RequestHash }
func (message Message) GetGID() GID               { return message.GID }
func (message Message) GetSenderId() SenderId     { return message.SenderId }
func (message Message) GetThreadId() MessageId {
	if message.ThreadId == "" {
		return message.MessageId
	}
	return message.ThreadId
}
func (message Message) GetParentMessageId() MessageId {
	if message.ParentMessageId == nil {
		return ""
	}
	return *message.ParentMessageId
}
func (message Message) GetResendOfMessageId() MessageId {
	if message.ResendOfMessageId == nil {
		return ""
	}
	return *message.ResendOfMessageId
}
func (message Message) GetInternetMessageId() string {
	if message.InternetMessageId == "" {
		return string(message.MessageId) + "@grovehub.local"
	}
	return message.InternetMessageId
}
func (message Message) GetInReplyTo() string {
	if message.InReplyTo == nil {
		return ""
	}
	return *message.InReplyTo
}
func (message Message) GetReferenceMessageIds() []string {
	return copyStrings([]string(message.ReferenceMessageIds))
}
func (message Message) GetTo() []string          { return copyStrings([]string(message.To)) }
func (message Message) GetCc() []string          { return copyStrings([]string(message.Cc)) }
func (message Message) GetBcc() []string         { return copyStrings([]string(message.Bcc)) }
func (message Message) GetSubject() string       { return message.Subject }
func (message Message) GetTextBody() string      { return message.TextBody }
func (message Message) GetHTMLBody() string      { return message.HTMLBody }
func (message Message) GetStatus() MessageStatus { return message.Status }
func (message Message) GetVersion() int64        { return message.Version }
func (message Message) GetCreatedAt() time.Time  { return message.CreatedAt }
func (message Message) GetUpdatedAt() time.Time  { return message.UpdatedAt }
func (message Message) GetSubmissionStartedAt() *time.Time {
	return message.SubmissionStartedAt
}

func (message Message) GetFailureMessage() string {
	if message.FailureMessage == nil {
		return ""
	}
	return *message.FailureMessage
}

func (message Message) SetThreadId(value MessageId) Message {
	message.ThreadId = value
	return message
}
func (message Message) SetParentMessageId(value MessageId) Message {
	message.ParentMessageId = optionalMessageId(value)
	return message
}
func (message Message) SetResendOfMessageId(value MessageId) Message {
	message.ResendOfMessageId = optionalMessageId(value)
	return message
}
func (message Message) SetInternetMessageId(value string) Message {
	message.InternetMessageId = value
	return message
}
func (message Message) SetInReplyTo(value string) Message {
	message.InReplyTo = optionalString(value)
	return message
}
func (message Message) SetReferenceMessageIds(values []string) Message {
	message.ReferenceMessageIds = pq.StringArray(copyStrings(values))
	return message
}
func (message Message) SetStatus(status MessageStatus) Message {
	message.Status = status
	return message
}
func (message Message) SetFailureMessage(value string) Message {
	message.FailureMessage = optionalString(value)
	return message
}
func (message Message) SetSubmissionStartedAt(value time.Time) Message {
	message.SubmissionStartedAt = &value
	return message
}
func (message Message) SetVersion(version int64) Message {
	message.Version = version
	return message
}
func (message Message) SetUpdatedAt(updatedAt time.Time) Message {
	message.UpdatedAt = updatedAt
	return message
}

func optionalMessageId(value MessageId) *MessageId {
	if value == "" {
		return nil
	}
	return &value
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func copyStrings(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	result := make([]string, len(values))
	copy(result, values)
	return result
}
