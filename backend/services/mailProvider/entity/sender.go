package entity

import "time"

type SenderId string
type SenderStatus string

const (
	SenderStatusActive      SenderStatus = "ACTIVE"
	SenderStatusDeactivated SenderStatus = "DEACTIVATED"
)

type Sender struct {
	SenderId    SenderId     `gorm:"column:sender_id;type:uuid;primaryKey"`
	DomainId    DomainId     `gorm:"column:domain_id;type:uuid;not null"`
	Email       string       `gorm:"column:email;not null"`
	DisplayName string       `gorm:"column:display_name;not null"`
	Status      SenderStatus `gorm:"column:status;not null"`
	Version     int64        `gorm:"column:version;not null"`
	CreatedAt   time.Time    `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time    `gorm:"column:updated_at;not null"`
}

func (Sender) TableName() string {
	return "mail_provider.senders"
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
		SenderId:    senderId,
		DomainId:    domainId,
		Email:       email,
		DisplayName: displayName,
		Status:      status,
		Version:     version,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func (sender Sender) GetSenderId() SenderId   { return sender.SenderId }
func (sender Sender) GetDomainId() DomainId   { return sender.DomainId }
func (sender Sender) GetEmail() string        { return sender.Email }
func (sender Sender) GetDisplayName() string  { return sender.DisplayName }
func (sender Sender) GetStatus() SenderStatus { return sender.Status }
func (sender Sender) GetVersion() int64       { return sender.Version }
func (sender Sender) GetCreatedAt() time.Time { return sender.CreatedAt }
func (sender Sender) GetUpdatedAt() time.Time { return sender.UpdatedAt }

func (sender Sender) SetStatus(status SenderStatus) Sender {
	sender.Status = status
	return sender
}

func (sender Sender) SetVersion(version int64) Sender {
	sender.Version = version
	return sender
}

func (sender Sender) SetUpdatedAt(updatedAt time.Time) Sender {
	sender.UpdatedAt = updatedAt
	return sender
}
