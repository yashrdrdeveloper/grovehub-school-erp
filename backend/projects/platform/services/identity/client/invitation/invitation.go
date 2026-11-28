// Package invitation exposes the read-only Identity delivery projection.
package invitation

import (
	"context"

	dto "grovehub/services/identity/internal/app/dto/invitation"
)

type Delivery = dto.Delivery

// @Client invitationClient
type Client interface {
	// @Proxy reader.ListDeliveryPage
	ListDeliveryPage(ctx context.Context, tenantID, after string, limit int) ([]dto.Delivery, error)
}

// DeliveryReader is the existing read-only invitation contract.
type DeliveryReader = Client
