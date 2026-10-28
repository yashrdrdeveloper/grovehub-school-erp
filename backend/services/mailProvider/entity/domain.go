package entity

import "time"

type DomainId string
type DomainStatus string

const (
	DomainStatusPendingVerification DomainStatus = "PENDING_VERIFICATION"
	DomainStatusVerified            DomainStatus = "VERIFIED"
	DomainStatusDeactivated         DomainStatus = "DEACTIVATED"
)

type Domain struct {
	DomainId          DomainId     `gorm:"column:domain_id;type:uuid;primaryKey"`
	Name              string       `gorm:"column:name;not null"`
	VerificationToken string       `gorm:"column:verification_token;not null"`
	Status            DomainStatus `gorm:"column:status;not null"`
	Version           int64        `gorm:"column:version;not null"`
	CreatedAt         time.Time    `gorm:"column:created_at;not null"`
	UpdatedAt         time.Time    `gorm:"column:updated_at;not null"`
}

func (Domain) TableName() string {
	return "mail_provider.domains"
}

func NewDomain(
	domainId DomainId,
	name string,
	verificationToken string,
	status DomainStatus,
	version int64,
	createdAt time.Time,
	updatedAt time.Time,
) Domain {
	return Domain{
		DomainId:          domainId,
		Name:              name,
		VerificationToken: verificationToken,
		Status:            status,
		Version:           version,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}

func (domain Domain) GetDomainId() DomainId        { return domain.DomainId }
func (domain Domain) GetName() string              { return domain.Name }
func (domain Domain) GetVerificationToken() string { return domain.VerificationToken }
func (domain Domain) GetStatus() DomainStatus      { return domain.Status }
func (domain Domain) GetVersion() int64            { return domain.Version }
func (domain Domain) GetCreatedAt() time.Time      { return domain.CreatedAt }
func (domain Domain) GetUpdatedAt() time.Time      { return domain.UpdatedAt }

func (domain Domain) SetStatus(status DomainStatus) Domain {
	domain.Status = status
	return domain
}

func (domain Domain) SetVersion(version int64) Domain {
	domain.Version = version
	return domain
}

func (domain Domain) SetUpdatedAt(updatedAt time.Time) Domain {
	domain.UpdatedAt = updatedAt
	return domain
}
