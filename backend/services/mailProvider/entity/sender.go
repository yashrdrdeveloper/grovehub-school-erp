package entity

import "time"

type SenderId string
type SenderStatus string

const (
	SenderStatusActive      SenderStatus = "ACTIVE"
	SenderStatusDeactivated SenderStatus = "DEACTIVATED"
)

type Sender struct {
	senderId    SenderId
	domainId    DomainId
	email       string
	displayName string
	status      SenderStatus
	version     int64
	createdAt   time.Time
	updatedAt   time.Time
}

func NewSender(
	senderId SenderId,
	domainId DomainId,
	email string,
	displayName string,
	status SenderStatus,
	version int64,
	createdAt time.Time,
	updatedAt time.Time,
) Sender {
	return Sender{
		senderId:    senderId,
		domainId:    domainId,
		email:       email,
		displayName: displayName,
		status:      status,
		version:     version,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (sender Sender) GetSenderId() SenderId { return sender.senderId }
func (sender Sender) GetDomainId() DomainId { return sender.domainId }
func (sender Sender) GetEmail() string { return sender.email }
func (sender Sender) GetDisplayName() string { return sender.displayName }
func (sender Sender) GetStatus() SenderStatus { return sender.status }
func (sender Sender) GetVersion() int64 { return sender.version }
func (sender Sender) GetCreatedAt() time.Time { return sender.createdAt }
func (sender Sender) GetUpdatedAt() time.Time { return sender.updatedAt }

func (sender Sender) SetStatus(status SenderStatus) Sender {
	sender.status = status
	return sender
}

func (sender Sender) SetVersion(version int64) Sender {
	sender.version = version
	return sender
}

func (sender Sender) SetUpdatedAt(updatedAt time.Time) Sender {
	sender.updatedAt = updatedAt
	return sender
}
