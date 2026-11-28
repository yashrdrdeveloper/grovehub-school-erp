package identity

import (
	"gorm.io/gorm"
	"grovehub/services/identity/client/invitation"
	repository "grovehub/services/identity/internal/app/repository/invitation/impl"
)

// NewInvitationDeliveryReader constructs an operator diagnostic read boundary
// without starting Identity workers or exposing its repository internals.
func NewInvitationDeliveryReader(db *gorm.DB) invitation.DeliveryReader {
	return repository.NewDeliveryReader(db)
}
