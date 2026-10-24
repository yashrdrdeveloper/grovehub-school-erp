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
	domainId          DomainId
	name              string
	verificationToken string
	status            DomainStatus
	version           int64
	createdAt         time.Time
	updatedAt         time.Time
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
		domainId:          domainId,
		name:              name,
		verificationToken: verificationToken,
		status:            status,
		version:           version,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
	}
}

func (domain Domain) GetDomainId() DomainId { return domain.domainId }
func (domain Domain) GetName() string { return domain.name }
func (domain Domain) GetVerificationToken() string { return domain.verificationToken }
func (domain Domain) GetStatus() DomainStatus { return domain.status }
func (domain Domain) GetVersion() int64 { return domain.version }
func (domain Domain) GetCreatedAt() time.Time { return domain.createdAt }
func (domain Domain) GetUpdatedAt() time.Time { return domain.updatedAt }

func (domain Domain) SetStatus(status DomainStatus) Domain {
	domain.status = status
	return domain
}

func (domain Domain) SetVersion(version int64) Domain {
	domain.version = version
	return domain
}

func (domain Domain) SetUpdatedAt(updatedAt time.Time) Domain {
	domain.updatedAt = updatedAt
	return domain
}
