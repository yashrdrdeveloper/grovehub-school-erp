package entity

import "time"

type MessageId string
type MessageStatus string

const (
	MessageStatusQueued    MessageStatus = "QUEUED"
	MessageStatusSubmitted MessageStatus = "SUBMITTED"
	MessageStatusFailed    MessageStatus = "FAILED"
)

type Message struct {
	messageId      MessageId
	senderId       SenderId
	to             []string
	cc             []string
	bcc            []string
	subject        string
	textBody       string
	htmlBody       string
	status         MessageStatus
	failureMessage string
	version        int64
	createdAt      time.Time
	updatedAt      time.Time
}

func NewMessage(
	messageId MessageId,
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
		messageId:      messageId,
		senderId:       senderId,
		to:             copyStrings(to),
		cc:             copyStrings(cc),
		bcc:            copyStrings(bcc),
		subject:        subject,
		textBody:       textBody,
		htmlBody:       htmlBody,
		status:         status,
		failureMessage: failureMessage,
		version:        version,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

func (message Message) GetMessageId() MessageId { return message.messageId }
func (message Message) GetSenderId() SenderId { return message.senderId }
func (message Message) GetTo() []string { return copyStrings(message.to) }
func (message Message) GetCc() []string { return copyStrings(message.cc) }
func (message Message) GetBcc() []string { return copyStrings(message.bcc) }
func (message Message) GetSubject() string { return message.subject }
func (message Message) GetTextBody() string { return message.textBody }
func (message Message) GetHTMLBody() string { return message.htmlBody }
func (message Message) GetStatus() MessageStatus { return message.status }
func (message Message) GetFailureMessage() string { return message.failureMessage }
func (message Message) GetVersion() int64 { return message.version }
func (message Message) GetCreatedAt() time.Time { return message.createdAt }
func (message Message) GetUpdatedAt() time.Time { return message.updatedAt }

func (message Message) SetStatus(status MessageStatus) Message {
	message.status = status
	return message
}

func (message Message) SetFailureMessage(value string) Message {
	message.failureMessage = value
	return message
}

func (message Message) SetVersion(version int64) Message {
	message.version = version
	return message
}

func (message Message) SetUpdatedAt(updatedAt time.Time) Message {
	message.updatedAt = updatedAt
	return message
}

func copyStrings(values []string) []string {
	result := make([]string, len(values))
	copy(result, values)
	return result
}
