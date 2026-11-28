// Package producer is the public Audit event publication contract.
package producer

import (
	"context"

	"grovehub/services/audit/client/producer/constants"
	dto "grovehub/services/audit/internal/app/dto/producer"
)

// Client-facing aliases preserve the event envelope and payload contracts.
type (
	Event               = dto.Event
	Activity            = dto.Activity
	InvitationEvent     = dto.InvitationEvent
	OperatorActivity    = dto.OperatorActivity
	SchoolRenameReceipt = dto.SchoolRenameReceipt
	SchoolStatusReceipt = dto.SchoolStatusReceipt
)

const (
	ActivityOccurred   = constants.ActivityOccurred
	InvitationDelivery = constants.InvitationDelivery
)

// @Client producerClient
type Client interface {
	// @Proxy service.Produce
	Produce(ctx context.Context, event dto.Event) error

	// @Proxy service.Record
	Record(ctx context.Context, activity dto.Activity) error
}
